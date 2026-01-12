package service

import (
	"encoding/json"
	"fmt"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
)

type ProgressValidator struct{}

func NewProgressValidator() *ProgressValidator {
	return &ProgressValidator{}
}

func (pv *ProgressValidator) Validate(progressData *json.RawMessage, itemType models.ItemType) error {
	if progressData == nil || len(*progressData) == 0 || string(*progressData) == "{}" {
		return nil
	}

	var data models.ProgressData
	if err := json.Unmarshal(*progressData, &data); err != nil {
		return fmt.Errorf("invalid progress data format: %w", err)
	}

	switch itemType {
	case models.ItemTypeBook:
		return pv.validateBook(&data)
	case models.ItemTypeManga:
		return pv.validateManga(&data)
	case models.ItemTypeGame:
		return pv.validateGame(&data)
	case models.ItemTypeSeries:
		return pv.validateSeries(&data)
	case models.ItemTypeMovie:
		return pv.validateMovie(&data)
	case models.ItemTypeMusic:
		return pv.validateMusic(&data)
	default:
		return fmt.Errorf("unsupported item type: %s", itemType)
	}
}

func (pv *ProgressValidator) validateBook(data *models.ProgressData) error {
	if data.Current != nil && data.Total != nil {
		if *data.Current < 0 {
			return fmt.Errorf("current progress cannot be negative")
		}
		if *data.Total <= 0 {
			return fmt.Errorf("total must be greater than 0")
		}
		if *data.Current > *data.Total {
			return fmt.Errorf("current progress (%d) cannot exceed total (%d)", *data.Current, *data.Total)
		}
		if data.Unit != nil {
			validUnits := map[string]bool{"pages": true, "chapters": true}
			if !validUnits[*data.Unit] {
				return fmt.Errorf("invalid unit for book: %s (must be 'pages' or 'chapters')", *data.Unit)
			}
		}
	}
	return nil
}

func (pv *ProgressValidator) validateManga(data *models.ProgressData) error {
	if data.Current != nil && data.Total != nil {
		if *data.Current < 0 {
			return fmt.Errorf("current progress cannot be negative")
		}
		if *data.Total <= 0 {
			return fmt.Errorf("total must be greater than 0")
		}
		if *data.Current > *data.Total {
			return fmt.Errorf("current progress (%d) cannot exceed total (%d)", *data.Current, *data.Total)
		}
		if data.Unit != nil {
			validUnits := map[string]bool{"volumes": true, "chapters": true}
			if !validUnits[*data.Unit] {
				return fmt.Errorf("invalid unit for manga: %s (must be 'volumes' or 'chapters')", *data.Unit)
			}
		}
	}
	return nil
}

func (pv *ProgressValidator) validateGame(data *models.ProgressData) error {
	if data.HoursPlayed != nil && *data.HoursPlayed < 0 {
		return fmt.Errorf("hours played cannot be negative")
	}
	if data.AchievementsUnlocked != nil {
		if *data.AchievementsUnlocked < 0 {
			return fmt.Errorf("achievements unlocked cannot be negative")
		}
		if data.AchievementsTotal != nil {
			if *data.AchievementsTotal <= 0 {
				return fmt.Errorf("total achievements must be greater than 0")
			}
			if *data.AchievementsUnlocked > *data.AchievementsTotal {
				return fmt.Errorf("achievements unlocked (%d) cannot exceed total (%d)", 
					*data.AchievementsUnlocked, *data.AchievementsTotal)
			}
		}
	}
	return nil
}

func (pv *ProgressValidator) validateSeries(data *models.ProgressData) error {
	if data.Current != nil && data.Total != nil {
		if *data.Current < 0 {
			return fmt.Errorf("current progress cannot be negative")
		}
		if *data.Total <= 0 {
			return fmt.Errorf("total must be greater than 0")
		}
		if *data.Current > *data.Total {
			return fmt.Errorf("current progress (%d) cannot exceed total (%d)", *data.Current, *data.Total)
		}
		if data.Unit != nil && *data.Unit != "episodes" {
			return fmt.Errorf("invalid unit for series: %s (must be 'episodes')", *data.Unit)
		}
	}
	if data.Season != nil && *data.Season < 1 {
		return fmt.Errorf("season must be greater than 0")
	}
	return nil
}

func (pv *ProgressValidator) validateMovie(data *models.ProgressData) error {
	if data.Watched != nil {
	}
	return nil
}

func (pv *ProgressValidator) validateMusic(data *models.ProgressData) error {
	if data.TracksListened != nil {
		if *data.TracksListened < 0 {
			return fmt.Errorf("tracks listened cannot be negative")
		}
		if data.TotalTracks != nil {
			if *data.TotalTracks <= 0 {
				return fmt.Errorf("total tracks must be greater than 0")
			}
			if *data.TracksListened > *data.TotalTracks {
				return fmt.Errorf("tracks listened (%d) cannot exceed total (%d)", 
					*data.TracksListened, *data.TotalTracks)
			}
		}
	}
	if data.AlbumsListened != nil {
		if *data.AlbumsListened < 0 {
			return fmt.Errorf("albums listened cannot be negative")
		}
		if data.TotalAlbums != nil {
			if *data.TotalAlbums <= 0 {
				return fmt.Errorf("total albums must be greater than 0")
			}
			if *data.AlbumsListened > *data.TotalAlbums {
				return fmt.Errorf("albums listened (%d) cannot exceed total (%d)", 
					*data.AlbumsListened, *data.TotalAlbums)
			}
		}
	}
	return nil
}

func (pv *ProgressValidator) CalculatePercentage(progressData *json.RawMessage, itemType models.ItemType) float64 {
	if progressData == nil || len(*progressData) == 0 || string(*progressData) == "{}" {
		return 0.0
	}

	var data models.ProgressData
	if err := json.Unmarshal(*progressData, &data); err != nil {
		return 0.0
	}

	switch itemType {
	case models.ItemTypeBook, models.ItemTypeManga, models.ItemTypeSeries:
		if data.Current != nil && data.Total != nil && *data.Total > 0 {
			return (float64(*data.Current) / float64(*data.Total)) * 100.0
		}
	case models.ItemTypeGame:
		if data.AchievementsUnlocked != nil && data.AchievementsTotal != nil && *data.AchievementsTotal > 0 {
			return (float64(*data.AchievementsUnlocked) / float64(*data.AchievementsTotal)) * 100.0
		}
	case models.ItemTypeMovie:
		if data.Watched != nil && *data.Watched {
			return 100.0
		}
	case models.ItemTypeMusic:
		if data.TracksListened != nil && data.TotalTracks != nil && *data.TotalTracks > 0 {
			return (float64(*data.TracksListened) / float64(*data.TotalTracks)) * 100.0
		}
		if data.AlbumsListened != nil && data.TotalAlbums != nil && *data.TotalAlbums > 0 {
			return (float64(*data.AlbumsListened) / float64(*data.TotalAlbums)) * 100.0
		}
	}
	return 0.0
}
