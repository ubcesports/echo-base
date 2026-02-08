package services

import (
	"context"
	"fmt"

	"github.com/ubcesports/echo-base/internal/interfaces/gamer"
	"github.com/ubcesports/echo-base/internal/models"
	"github.com/ubcesports/echo-base/internal/utils"
)

type gamerActivityService struct {
	activityRepo gamer.GamerActivityRepository
	profileRepo  gamer.GamerProfileRepository
}

func NewGamerActivityService(activityRepo gamer.GamerActivityRepository, profileRepo gamer.GamerProfileRepository) GamerActivityService {
	return &gamerActivityService{
		activityRepo: activityRepo,
		profileRepo:  profileRepo,
	}
}

func (s *gamerActivityService) GetActivitiesByStudent(ctx context.Context, studentNumber string) ([]models.GamerActivity, error) {
	if err := validateStudentNumber(studentNumber); err != nil {
		return nil, err
	}

	return s.activityRepo.GetByStudentNumber(ctx, studentNumber)
}

func (s *gamerActivityService) GetTodayActivities(ctx context.Context, studentNumber string) ([]models.GamerActivity, error) {
	if err := validateStudentNumber(studentNumber); err != nil {
		return nil, err
	}

	return s.activityRepo.GetTodayActivitiesByStudent(ctx, studentNumber)
}

func (s *gamerActivityService) GetRecentActivities(ctx context.Context, page, limit int, search string) ([]models.GamerActivity, error) {
	if page < 1 {
		return nil, fmt.Errorf("page must be >= 1")
	}

	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100")
	}

	return s.activityRepo.GetRecentActivities(ctx, page, limit, search)
}

func (s *gamerActivityService) StartActivity(ctx context.Context, req *models.CreateActivityRequest) (*models.GamerActivity, error) {
	if err := validateStudentNumber(req.StudentNumber); err != nil {
		return nil, err
	}

	if req.Game == "" {
		return nil, fmt.Errorf("game is required")
	}

	tierNum, expiryDate, err := s.profileRepo.CheckMembershipValidity(ctx, req.StudentNumber)
	if err != nil {
		return nil, fmt.Errorf("Foreign key %s not found.", req.StudentNumber)
	}

	tier, err := models.NewMembershipTier(tierNum)
	if err != nil {
		return nil, fmt.Errorf("invalid membership tier: %w", err)
	}

	expired, err := tier.IsExpired(expiryDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check expiry: %w", err)
	}

	if expired {
		expiryDateStr := "unknown"
		if expiryDate != nil {
			expiryDateStr = expiryDate.Format("2006-01-02")
		}
		return nil, fmt.Errorf("%s membership expired on %s. Please ask the user to purchase a new membership. If the member has already purchased a new membership for this year please verify via Showpass then create a new profile for them.", tier.GetName(), expiryDateStr)
	}

	startedAt, err := utils.NowInPacific()
	if err != nil {
		return nil, fmt.Errorf("failed to get current time: %w", err)
	}

	activity := &models.GamerActivity{
		StudentNumber: req.StudentNumber,
		PCNumber:      req.PCNumber,
		Game:          req.Game,
		StartedAt:     startedAt,
	}

	return s.activityRepo.Create(ctx, activity)
}

func (s *gamerActivityService) EndActivity(ctx context.Context, studentNumber string, req *models.UpdateActivityRequest) (*models.GamerActivity, error) {
	if err := validateStudentNumber(studentNumber); err != nil {
		return nil, err
	}

	if req.ExecName == "" {
		return nil, fmt.Errorf("exec_name is required")
	}

	endedAt, err := utils.NowInPacific()
	if err != nil {
		return nil, fmt.Errorf("failed to get current time: %w", err)
	}

	return s.activityRepo.UpdateEndTime(ctx, studentNumber, req.PCNumber, endedAt, req.ExecName)
}

func (s *gamerActivityService) GetActiveSessions(ctx context.Context) ([]models.GamerActivity, error) {
	return s.activityRepo.GetActiveSessions(ctx)
}
