package gamer

import (
	"context"
	"time"

	"github.com/ubcesports/echo-base/internal/models"
)

type GamerProfileRepository interface {
	GetByStudentNumber(ctx context.Context, studentNumber string) (*models.GamerProfile, error)
	Upsert(ctx context.Context, profile *models.GamerProfile) (*models.GamerProfile, error)
	Delete(ctx context.Context, studentNumber string) error
	CheckMembershipValidity(ctx context.Context, studentNumber string) (tier int, expiryDate *time.Time, err error)
}

type GamerActivityRepository interface {
	GetByStudentNumber(ctx context.Context, studentNumber string) ([]models.GamerActivity, error)
	GetTodayActivitiesByStudent(ctx context.Context, studentNumber string) ([]models.GamerActivity, error)
	GetRecentActivities(ctx context.Context, page, limit int, search string) ([]models.GamerActivity, error)
	Create(ctx context.Context, activity *models.GamerActivity) (*models.GamerActivity, error)
	UpdateEndTime(ctx context.Context, studentNumber string, pcNumber int, endedAt time.Time, execName string) (*models.GamerActivity, error)
	GetActiveSessions(ctx context.Context) ([]models.GamerActivity, error)
}
