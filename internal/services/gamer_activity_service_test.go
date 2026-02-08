package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ubcesports/echo-base/internal/models"
)

type mockGamerActivityRepository struct {
	activities []models.GamerActivity
}

func (m *mockGamerActivityRepository) GetByStudentNumber(ctx context.Context, studentNumber string) ([]models.GamerActivity, error) {
	var result []models.GamerActivity
	for _, a := range m.activities {
		if a.StudentNumber == studentNumber {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockGamerActivityRepository) GetTodayActivitiesByStudent(ctx context.Context, studentNumber string) ([]models.GamerActivity, error) {
	var result []models.GamerActivity
	for _, a := range m.activities {
		if a.StudentNumber == studentNumber {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockGamerActivityRepository) GetRecentActivities(ctx context.Context, page, limit int, search string) ([]models.GamerActivity, error) {
	return m.activities, nil
}

func (m *mockGamerActivityRepository) Create(ctx context.Context, activity *models.GamerActivity) (*models.GamerActivity, error) {
	m.activities = append(m.activities, *activity)
	return activity, nil
}

func (m *mockGamerActivityRepository) UpdateEndTime(ctx context.Context, studentNumber string, pcNumber int, endedAt time.Time, execName string) (*models.GamerActivity, error) {
	for i, a := range m.activities {
		if a.StudentNumber == studentNumber && a.PCNumber == pcNumber && a.EndedAt == nil {
			m.activities[i].EndedAt = &endedAt
			m.activities[i].ExecName = &execName
			return &m.activities[i], nil
		}
	}
	return nil, nil
}

func (m *mockGamerActivityRepository) GetActiveSessions(ctx context.Context) ([]models.GamerActivity, error) {
	var result []models.GamerActivity
	for _, a := range m.activities {
		if a.EndedAt == nil {
			result = append(result, a)
		}
	}
	return result, nil
}

func TestStartActivity(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name          string
		req           *models.CreateActivityRequest
		tierNum       int
		expiryDate    *time.Time
		profileExists bool
		wantErr       bool
		errContains   string
	}{
		{
			name: "valid check-in with tier 1",
			req: &models.CreateActivityRequest{
				StudentNumber: "12345678",
				PCNumber:      1,
				Game:          "League of Legends",
			},
			tierNum:       1,
			expiryDate:    &tomorrow,
			profileExists: true,
			wantErr:       false,
		},
		{
			name: "expired membership",
			req: &models.CreateActivityRequest{
				StudentNumber: "12345678",
				PCNumber:      1,
				Game:          "League of Legends",
			},
			tierNum:       1,
			expiryDate:    &yesterday,
			profileExists: true,
			wantErr:       true,
			errContains:   "expired",
		},
		{
			name: "student not found",
			req: &models.CreateActivityRequest{
				StudentNumber: "99999999",
				PCNumber:      1,
				Game:          "League of Legends",
			},
			profileExists: false,
			wantErr:       true,
			errContains:   "Foreign key",
		},
		{
			name: "invalid student number",
			req: &models.CreateActivityRequest{
				StudentNumber: "123",
				PCNumber:      1,
				Game:          "League of Legends",
			},
			wantErr:     true,
			errContains: "8 digits",
		},
		{
			name: "missing game",
			req: &models.CreateActivityRequest{
				StudentNumber: "12345678",
				PCNumber:      1,
				Game:          "",
			},
			tierNum:       1,
			expiryDate:    &tomorrow,
			profileExists: true,
			wantErr:       true,
			errContains:   "game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActivityRepo := &mockGamerActivityRepository{
				activities: []models.GamerActivity{},
			}

			mockProfileRepo := &mockGamerProfileRepository{
				profiles: make(map[string]*models.GamerProfile),
			}

			if tt.profileExists {
				mockProfileRepo.profiles[tt.req.StudentNumber] = &models.GamerProfile{
					StudentNumber:        tt.req.StudentNumber,
					MembershipTier:       tt.tierNum,
					MembershipExpiryDate: tt.expiryDate,
				}
			}

			service := NewGamerActivityService(mockActivityRepo, mockProfileRepo)

			activity, err := service.StartActivity(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("StartActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %v", tt.errContains, err)
				}
			}

			if !tt.wantErr && activity == nil {
				t.Errorf("expected activity, got nil")
			}
		})
	}
}

func TestGetRecentActivities(t *testing.T) {
	tests := []struct {
		name    string
		page    int
		limit   int
		wantErr bool
	}{
		{"valid pagination", 1, 10, false},
		{"invalid page", 0, 10, true},
		{"invalid limit - too small", 1, 0, true},
		{"invalid limit - too large", 1, 101, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActivityRepo := &mockGamerActivityRepository{}
			mockProfileRepo := &mockGamerProfileRepository{
				profiles: make(map[string]*models.GamerProfile),
			}
			service := NewGamerActivityService(mockActivityRepo, mockProfileRepo)

			_, err := service.GetRecentActivities(context.Background(), tt.page, tt.limit, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecentActivities() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEndActivity(t *testing.T) {
	tests := []struct {
		name          string
		studentNumber string
		req           *models.UpdateActivityRequest
		wantErr       bool
	}{
		{
			name:          "valid end activity",
			studentNumber: "12345678",
			req: &models.UpdateActivityRequest{
				PCNumber: 1,
				ExecName: "Admin",
			},
			wantErr: false,
		},
		{
			name:          "invalid student number",
			studentNumber: "123",
			req: &models.UpdateActivityRequest{
				PCNumber: 1,
				ExecName: "Admin",
			},
			wantErr: true,
		},
		{
			name:          "missing exec name",
			studentNumber: "12345678",
			req: &models.UpdateActivityRequest{
				PCNumber: 1,
				ExecName: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActivityRepo := &mockGamerActivityRepository{
				activities: []models.GamerActivity{
					{
						StudentNumber: "12345678",
						PCNumber:      1,
						Game:          "Test Game",
						StartedAt:     time.Now(),
						EndedAt:       nil,
					},
				},
			}
			mockProfileRepo := &mockGamerProfileRepository{
				profiles: make(map[string]*models.GamerProfile),
			}
			service := NewGamerActivityService(mockActivityRepo, mockProfileRepo)

			_, err := service.EndActivity(context.Background(), tt.studentNumber, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
