package challenges

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
	"github.com/GATEOPENERZ/completionist-api/internal/services"
)

type Service struct {
	chalRepo    *repository.ChallengeRepository
	userRepo    *repository.UserRepository
	rankService *services.RankService
}

func NewService(chalRepo *repository.ChallengeRepository, userRepo *repository.UserRepository, rankService *services.RankService) *Service {
	return &Service{
		chalRepo:    chalRepo,
		userRepo:    userRepo,
		rankService: rankService,
	}
}

func (s *Service) NotifyAction(ctx context.Context, userID int64, actionType string, metadata map[string]interface{}) error {
	season, _ := s.chalRepo.GetActiveSeason()

	challenges, err := s.chalRepo.GetChallengesByFrequency(models.FrequencyDaily)
	if err != nil {
		return err
	}
	weekly, _ := s.chalRepo.GetChallengesByFrequency(models.FrequencyWeekly)
	challenges = append(challenges, weekly...)

	if season != nil {
		seasonal, _ := s.chalRepo.GetChallengesByFrequency(models.FrequencySeasonal)
		challenges = append(challenges, seasonal...)
	}

	for _, chal := range challenges {
		if err := s.evaluateChallenge(userID, chal, season, actionType, metadata); err != nil {
			log.Printf("Error evaluating challenge %d: %v", chal.ID, err)
		}
	}
	return nil
}

func (s *Service) evaluateChallenge(userID int64, chal models.Challenge, season *models.Season, actionType string, metadata map[string]interface{}) error {
	var seasonID *int64
	if season != nil {
		seasonID = &season.ID
	}

	uc, err := s.chalRepo.FindUserChallenge(userID, chal.ID, seasonID)
	if err != nil {
		target := 1
		switch chal.CriteriaType {
		case models.CriteriaCountItems:
			var cm models.CriteriaCountMetadata
			_ = json.Unmarshal(chal.CriteriaMetadata, &cm)
			if cm.Target > 0 {
				target = cm.Target
			}
		case models.CriteriaGenreCount:
			var cm models.CriteriaGenreCountMetadata
			_ = json.Unmarshal(chal.CriteriaMetadata, &cm)
			if cm.Target > 0 {
				target = cm.Target
			}
		}

		uc = &models.UserChallenge{
			UserID:          userID,
			ChallengeID:     chal.ID,
			SeasonID:        seasonID,
			CurrentProgress: 0,
			TargetProgress:  target,
			IsCompleted:     false,
		}
	}

	if uc.IsCompleted {
		return nil
	}

	matched := false

	switch chal.CriteriaType {
	case models.CriteriaCountItems:
		matched = s.checkCountCriteria(chal.CriteriaMetadata, actionType, metadata)
	case models.CriteriaSpecificItem:
		matched = s.checkSpecificItemCriteria(chal.CriteriaMetadata, actionType, metadata)
	case models.CriteriaGenreCount:
		matched = s.checkGenreCriteria(chal.CriteriaMetadata, actionType, metadata)
	}

	if matched {
		uc.CurrentProgress++
		if uc.CurrentProgress >= uc.TargetProgress {
			uc.CurrentProgress = uc.TargetProgress
			uc.IsCompleted = true
			now := time.Now()
			uc.CompletedAt = &now

			desc := fmt.Sprintf("Completed challenge: %s", chal.Title)
			sid := fmt.Sprintf("%d", chal.ID)
			_ = s.userRepo.AddXP(userID, "challenge_complete", &sid, &desc)

			if seasonID != nil {
				_, _ = s.rankService.UpdateElo(userID, *seasonID, chal.XPReward)
			}
		}
		return s.chalRepo.UpsertUserChallenge(uc)
	}

	return nil
}

func (s *Service) checkCountCriteria(rawMeta json.RawMessage, actionType string, meta map[string]interface{}) bool {
	if actionType != "complete_item" {
		return false
	}
	var criteria models.CriteriaCountMetadata
	if err := json.Unmarshal(rawMeta, &criteria); err != nil {
		return false
	}

	if criteria.ItemType != "" {
		if val, ok := meta["item_type"]; ok {
			if strVal, ok := val.(string); ok && strVal != criteria.ItemType {
				return false
			}
		}
	}
	return true
}

func (s *Service) checkSpecificItemCriteria(rawMeta json.RawMessage, actionType string, meta map[string]interface{}) bool {
	var criteria models.CriteriaSpecificItemMetadata
	if err := json.Unmarshal(rawMeta, &criteria); err != nil {
		return false
	}

	if criteria.Action != "" && criteria.Action != actionType {
		return false
	}

	if val, ok := meta["media_item_id"]; ok {
		if strVal, ok := val.(string); ok && strVal == criteria.MediaItemID {
			return true
		}
	}
	return false
}

func (s *Service) checkGenreCriteria(rawMeta json.RawMessage, actionType string, meta map[string]interface{}) bool {
	if actionType != "complete_item" {
		return false
	}
	var criteria models.CriteriaGenreCountMetadata
	if err := json.Unmarshal(rawMeta, &criteria); err != nil {
		return false
	}

	if val, ok := meta["genres"]; ok {
		if genres, ok := val.([]string); ok {
			for _, g := range genres {
				if strings.EqualFold(g, criteria.Genre) {
					return true
				}
			}
		}
	}
	return false
}