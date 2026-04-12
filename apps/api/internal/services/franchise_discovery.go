package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
)

type FranchiseDiscoveryService struct {
	franchiseRepo *repository.FranchiseRepository
	tmdbKey       string
}

func NewFranchiseDiscoveryService(franchiseRepo *repository.FranchiseRepository, tmdbKey string) *FranchiseDiscoveryService {
	return &FranchiseDiscoveryService{franchiseRepo: franchiseRepo, tmdbKey: tmdbKey}
}

type tmdbCollection struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

type tmdbCollectionSearchResult struct {
	Results []tmdbCollection `json:"results"`
}

type tmdbTrendingResult struct {
	Results []struct {
		ID                 int    `json:"id"`
		Title              string `json:"title"`
		Name               string `json:"name"`
		BelongsToCollection *struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			PosterPath string `json:"poster_path"`
		} `json:"belongs_to_collection"`
	} `json:"results"`
}

func (s *FranchiseDiscoveryService) DiscoverMovieFranchises() ([]models.Franchise, error) {
	if s.tmdbKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	// Fetch trending movies and extract their collections
	url := fmt.Sprintf("https://api.themoviedb.org/3/trending/movie/week?api_key=%s", s.tmdbKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var trending tmdbTrendingResult
	if err := json.Unmarshal(body, &trending); err != nil {
		return nil, err
	}

	seen := map[int]bool{}
	var franchises []models.Franchise

	for _, movie := range trending.Results {
		col := movie.BelongsToCollection
		if col == nil || seen[col.ID] {
			continue
		}
		seen[col.ID] = true

		externalID := fmt.Sprintf("%d", col.ID)
		source := "tmdb"
		category := "movies"
		var coverURL *string
		if col.PosterPath != "" {
			u := "https://image.tmdb.org/t/p/w500" + col.PosterPath
			coverURL = &u
		}

		f := &models.Franchise{
			Name:          col.Name,
			CoverImageURL: coverURL,
			Source:        &source,
			ExternalID:    &externalID,
			Category:      &category,
		}
		if err := s.franchiseRepo.UpsertFromSource(f); err == nil {
			franchises = append(franchises, *f)
		}
	}

	// Also fetch popular TV shows for series franchises
	tvURL := fmt.Sprintf("https://api.themoviedb.org/3/trending/tv/week?api_key=%s", s.tmdbKey)
	tvResp, err := http.Get(tvURL)
	if err == nil {
		defer tvResp.Body.Close()
		tvBody, _ := io.ReadAll(tvResp.Body)
		var tvTrending struct {
			Results []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				PosterPath string `json:"poster_path"`
			} `json:"results"`
		}
		if json.Unmarshal(tvBody, &tvTrending) == nil {
			for _, show := range tvTrending.Results {
				externalID := fmt.Sprintf("tv-%d", show.ID)
				source := "tmdb"
				category := "series"
				var coverURL *string
				if show.PosterPath != "" {
					u := "https://image.tmdb.org/t/p/w500" + show.PosterPath
					coverURL = &u
				}
				f := &models.Franchise{
					Name:          show.Name,
					CoverImageURL: coverURL,
					Source:        &source,
					ExternalID:    &externalID,
					Category:      &category,
				}
				if err := s.franchiseRepo.UpsertFromSource(f); err == nil {
					franchises = append(franchises, *f)
				}
			}
		}
	}

	return franchises, nil
}

func (s *FranchiseDiscoveryService) DiscoverByCategory(category string) ([]models.Franchise, error) {
	// First check if we already have data
	existing, err := s.franchiseRepo.GetByCategory(category, 50, 0)
	if err == nil && len(existing) > 0 {
		return existing, nil
	}

	// Auto-discover from APIs
	switch category {
	case "movies", "series":
		return s.DiscoverMovieFranchises()
	default:
		return s.franchiseRepo.GetByCategory(category, 50, 0)
	}
}
