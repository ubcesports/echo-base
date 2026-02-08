package services

import (
	"context"

	"github.com/ubcesports/echo-base/internal/models"
)

type GamerProfileService interface {
	GetProfile(ctx context.Context, studentNumber string) (*models.GamerProfile, error)
	CreateOrUpdateProfile(ctx context.Context, req *models.CreateGamerProfileRequest) (*models.GamerProfile, error)
	DeleteProfile(ctx context.Context, studentNumber string) error
}

type GamerActivityService interface {
	GetActivitiesByStudent(ctx context.Context, studentNumber string) ([]models.GamerActivity, error)
	GetTodayActivities(ctx context.Context, studentNumber string) ([]models.GamerActivity, error)
	GetRecentActivities(ctx context.Context, page, limit int, search string) ([]models.GamerActivity, error)
	StartActivity(ctx context.Context, req *models.CreateActivityRequest) (*models.GamerActivity, error)
	EndActivity(ctx context.Context, studentNumber string, req *models.UpdateActivityRequest) (*models.GamerActivity, error)
	GetActiveSessions(ctx context.Context) ([]models.GamerActivity, error)
}
