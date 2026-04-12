package importer

import (
	"fmt"
	"log"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
)

type Provider interface {
	Name() string
	AuthURL(redirectURI string) string
	ExchangeCode(code, redirectURI string) (*models.ConnectedAccount, error)
	Import(account *models.ConnectedAccount, mediaRepo *repository.MediaRepository, listRepo *repository.ListRepository) (*ImportResult, error)
}

type ImportResult struct {
	TotalItems    int
	ImportedItems int
	SkippedItems  int
	FailedItems   int
}

type Service struct {
	providers    map[string]Provider
	connRepo     *repository.ConnectedAccountRepository
	mediaRepo    *repository.MediaRepository
	listRepo     *repository.ListRepository
}

func NewService(connRepo *repository.ConnectedAccountRepository, mediaRepo *repository.MediaRepository, listRepo *repository.ListRepository) *Service {
	return &Service{
		providers: make(map[string]Provider),
		connRepo:  connRepo,
		mediaRepo: mediaRepo,
		listRepo:  listRepo,
	}
}

func (s *Service) RegisterProvider(p Provider) {
	s.providers[p.Name()] = p
}

func (s *Service) GetProvider(name string) (Provider, error) {
	p, ok := s.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return p, nil
}

func (s *Service) ListProviders() []string {
	names := make([]string, 0, len(s.providers))
	for name := range s.providers {
		names = append(names, name)
	}
	return names
}

func (s *Service) GetAuthURL(provider, redirectURI string) (string, error) {
	p, err := s.GetProvider(provider)
	if err != nil {
		return "", err
	}
	return p.AuthURL(redirectURI), nil
}

func (s *Service) HandleCallback(userID int64, provider, code, redirectURI string) (*models.ConnectedAccount, error) {
	p, err := s.GetProvider(provider)
	if err != nil {
		return nil, err
	}

	account, err := p.ExchangeCode(code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	account.UserID = userID
	account.Provider = provider

	existing, _ := s.connRepo.GetByProvider(userID, provider)
	if existing != nil {
		existing.AccessToken = account.AccessToken
		existing.RefreshToken = account.RefreshToken
		existing.TokenExpiresAt = account.TokenExpiresAt
		existing.ProviderUserID = account.ProviderUserID
		existing.ProviderUsername = account.ProviderUsername
		if err := s.connRepo.Update(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	if err := s.connRepo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) TriggerSync(userID int64, provider string) (*models.ImportJob, error) {
	account, err := s.connRepo.GetByProvider(userID, provider)
	if err != nil || account == nil {
		return nil, fmt.Errorf("no connected account for provider %s", provider)
	}

	job := &models.ImportJob{
		UserID:             userID,
		ConnectedAccountID: &account.ID,
		Provider:           provider,
	}
	if err := s.connRepo.CreateImportJob(job); err != nil {
		return nil, err
	}

	go s.runImport(job, account)
	return job, nil
}

func (s *Service) runImport(job *models.ImportJob, account *models.ConnectedAccount) {
	p, err := s.GetProvider(job.Provider)
	if err != nil {
		s.failJob(job, err.Error())
		return
	}

	now := time.Now()
	job.Status = models.ImportJobRunning
	job.StartedAt = &now
	_ = s.connRepo.UpdateImportJob(job)

	result, err := p.Import(account, s.mediaRepo, s.listRepo)
	if err != nil {
		s.failJob(job, err.Error())
		return
	}

	completed := time.Now()
	job.Status = models.ImportJobCompleted
	job.TotalItems = result.TotalItems
	job.ImportedItems = result.ImportedItems
	job.SkippedItems = result.SkippedItems
	job.FailedItems = result.FailedItems
	job.CompletedAt = &completed

	account.LastSyncedAt = &completed
	_ = s.connRepo.Update(account)
	_ = s.connRepo.UpdateImportJob(job)

	log.Printf("import job %d completed: %d imported, %d skipped, %d failed", job.ID, result.ImportedItems, result.SkippedItems, result.FailedItems)
}

func (s *Service) failJob(job *models.ImportJob, errMsg string) {
	now := time.Now()
	job.Status = models.ImportJobFailed
	job.CompletedAt = &now
	_ = s.connRepo.UpdateImportJob(job)
	log.Printf("import job %d failed: %s", job.ID, errMsg)
}
