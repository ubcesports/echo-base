package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ubcesports/echo-base/internal/models"
)

type mockGamerProfileRepository struct {
	profiles map[string]*models.GamerProfile
	getErr   error
	upsertErr error
	deleteErr error
}

func (m *mockGamerProfileRepository) GetByStudentNumber(ctx context.Context, studentNumber string) (*models.GamerProfile, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	profile, exists := m.profiles[studentNumber]
	if !exists {
		return nil, nil
	}
	return profile, nil
}

func (m *mockGamerProfileRepository) Upsert(ctx context.Context, profile *models.GamerProfile) (*models.GamerProfile, error) {
	if m.upsertErr != nil {
		return nil, m.upsertErr
	}
	m.profiles[profile.StudentNumber] = profile
	return profile, nil
}

func (m *mockGamerProfileRepository) Delete(ctx context.Context, studentNumber string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.profiles, studentNumber)
	return nil
}

func (m *mockGamerProfileRepository) CheckMembershipValidity(ctx context.Context, studentNumber string) (tier int, expiryDate *time.Time, err error) {
	profile, exists := m.profiles[studentNumber]
	if !exists {
		return 0, nil, fmt.Errorf("student %s not found", studentNumber)
	}
	return profile.MembershipTier, profile.MembershipExpiryDate, nil
}

func TestCreateOrUpdateProfile(t *testing.T) {
	tests := []struct {
		name        string
		req         *models.CreateGamerProfileRequest
		wantErr     bool
		errContains string
	}{
		{
			name: "valid tier 1 profile",
			req: &models.CreateGamerProfileRequest{
				StudentNumber:  "12345678",
				FirstName:      "John",
				LastName:       "Doe",
				MembershipTier: 1,
			},
			wantErr: false,
		},
		{
			name: "valid tier 2 profile",
			req: &models.CreateGamerProfileRequest{
				StudentNumber:  "87654321",
				FirstName:      "Jane",
				LastName:       "Smith",
				MembershipTier: 2,
			},
			wantErr: false,
		},
		{
			name: "invalid student number - too short",
			req: &models.CreateGamerProfileRequest{
				StudentNumber:  "1234567",
				FirstName:      "John",
				LastName:       "Doe",
				MembershipTier: 1,
			},
			wantErr:     true,
			errContains: "8 digits",
		},
		{
			name: "invalid student number - contains letters",
			req: &models.CreateGamerProfileRequest{
				StudentNumber:  "1234567a",
				FirstName:      "John",
				LastName:       "Doe",
				MembershipTier: 1,
			},
			wantErr:     true,
			errContains: "8 digits",
		},
		{
			name: "missing first name",
			req: &models.CreateGamerProfileRequest{
				StudentNumber:  "12345678",
				FirstName:      "",
				LastName:       "Doe",
				MembershipTier: 1,
			},
			wantErr:     true,
			errContains: "first_name",
		},
		{
			name: "missing last name",
			req: &models.CreateGamerProfileRequest{
				StudentNumber:  "12345678",
				FirstName:      "John",
				LastName:       "",
				MembershipTier: 1,
			},
			wantErr:     true,
			errContains: "last_name",
		},
		{
			name: "invalid tier",
			req: &models.CreateGamerProfileRequest{
				StudentNumber:  "12345678",
				FirstName:      "John",
				LastName:       "Doe",
				MembershipTier: 5,
			},
			wantErr:     true,
			errContains: "tier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockGamerProfileRepository{
				profiles: make(map[string]*models.GamerProfile),
			}
			service := NewGamerProfileService(mockRepo)

			profile, err := service.CreateOrUpdateProfile(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdateProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || len(err.Error()) == 0 {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				}
			}

			if !tt.wantErr {
				if profile == nil {
					t.Errorf("expected profile, got nil")
					return
				}
				if profile.StudentNumber != tt.req.StudentNumber {
					t.Errorf("StudentNumber = %v, want %v", profile.StudentNumber, tt.req.StudentNumber)
				}
				if profile.MembershipExpiryDate == nil && tt.req.MembershipTier != 0 {
					t.Errorf("expected non-nil expiry date for tier %d", tt.req.MembershipTier)
				}
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	mockRepo := &mockGamerProfileRepository{
		profiles: map[string]*models.GamerProfile{
			"12345678": {
				StudentNumber:  "12345678",
				FirstName:      "John",
				LastName:       "Doe",
				MembershipTier: 1,
			},
		},
	}
	service := NewGamerProfileService(mockRepo)

	tests := []struct {
		name          string
		studentNumber string
		wantErr       bool
	}{
		{"existing profile", "12345678", false},
		{"invalid student number", "123", true},
		{"invalid student number with letters", "1234567a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetProfile(context.Background(), tt.studentNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	mockRepo := &mockGamerProfileRepository{
		profiles: map[string]*models.GamerProfile{
			"12345678": {
				StudentNumber: "12345678",
			},
		},
	}
	service := NewGamerProfileService(mockRepo)

	tests := []struct {
		name          string
		studentNumber string
		wantErr       bool
	}{
		{"valid deletion", "12345678", false},
		{"invalid student number", "123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteProfile(context.Background(), tt.studentNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
