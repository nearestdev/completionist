package services

import (
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
)

type StreakService struct {
	streakRepo *repository.StreakRepository
	badgeSvc   *BadgeService
}

func NewStreakService(streakRepo *repository.StreakRepository, badgeSvc *BadgeService) *StreakService {
	return &StreakService{streakRepo: streakRepo, badgeSvc: badgeSvc}
}

func (s *StreakService) TouchStreak(userID int64) (*models.UserStreak, error) {
	streak, err := s.streakRepo.TouchStreak(userID, time.Now())
	if err != nil {
		return nil, err
	}
	if s.badgeSvc != nil {
		_ = s.badgeSvc.CheckStreakBadges(userID, streak.CurrentStreak)
	}
	return streak, nil
}
