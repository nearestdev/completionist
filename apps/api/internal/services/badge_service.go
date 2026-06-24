package services

import (
	"encoding/json"
	"log"

	"github.com/nearestdev/completionist/internal/models"
	"github.com/nearestdev/completionist/internal/repository"
)

type BadgeService struct {
	badgeRepo *repository.BadgeRepository
}

func NewBadgeService(badgeRepo *repository.BadgeRepository) *BadgeService {
	return &BadgeService{badgeRepo: badgeRepo}
}

type badgeCriteria struct {
	Target   int    `json:"target"`
	ItemType string `json:"item_type,omitempty"`
}

func (s *BadgeService) CheckAndAwardAll(userID int64) error {
	badges, err := s.badgeRepo.GetAll()
	if err != nil {
		return err
	}

	for _, badge := range badges {
		has, err := s.badgeRepo.HasBadge(userID, badge.ID)
		if err != nil || has {
			continue
		}

		var criteria badgeCriteria
		if err := json.Unmarshal(badge.CriteriaJSON, &criteria); err != nil {
			log.Printf("badge %s: invalid criteria json: %v", badge.Code, err)
			continue
		}

		earned := false
		switch badge.CriteriaType {
		case models.BadgeCriteriaItemCompletionCount:
			count, err := s.badgeRepo.GetCompletionCountByType(userID, nil)
			if err == nil && count >= criteria.Target {
				earned = true
			}
		case models.BadgeCriteriaItemTypeCount:
			count, err := s.badgeRepo.GetCompletionCountByType(userID, &criteria.ItemType)
			if err == nil && count >= criteria.Target {
				earned = true
			}
		}

		if earned {
			_ = s.badgeRepo.AwardBadge(userID, badge.ID)
		}
	}
	return nil
}

func (s *BadgeService) CheckStreakBadges(userID int64, currentStreak int) error {
	badges, err := s.badgeRepo.GetAll()
	if err != nil {
		return err
	}

	for _, badge := range badges {
		if badge.CriteriaType != models.BadgeCriteriaStreakThreshold {
			continue
		}
		has, err := s.badgeRepo.HasBadge(userID, badge.ID)
		if err != nil || has {
			continue
		}

		var criteria badgeCriteria
		if err := json.Unmarshal(badge.CriteriaJSON, &criteria); err != nil {
			continue
		}
		if currentStreak >= criteria.Target {
			_ = s.badgeRepo.AwardBadge(userID, badge.ID)
		}
	}
	return nil
}
