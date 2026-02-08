package services

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ubcesports/echo-base/internal/errors"
	"github.com/ubcesports/echo-base/internal/interfaces/gamer"
	"github.com/ubcesports/echo-base/internal/models"
	"github.com/ubcesports/echo-base/internal/utils"
)

var studentNumberRegex = regexp.MustCompile(`^\d{8}$`)

type gamerProfileService struct {
	repo gamer.GamerProfileRepository
}

func NewGamerProfileService(repo gamer.GamerProfileRepository) GamerProfileService {
	return &gamerProfileService{repo: repo}
}

func (s *gamerProfileService) GetProfile(ctx context.Context, studentNumber string) (*models.GamerProfile, error) {
	if err := validateStudentNumber(studentNumber); err != nil {
		return nil, err
	}

	return s.repo.GetByStudentNumber(ctx, studentNumber)
}

func (s *gamerProfileService) CreateOrUpdateProfile(ctx context.Context, req *models.CreateGamerProfileRequest) (*models.GamerProfile, error) {
	if err := validateStudentNumber(req.StudentNumber); err != nil {
		return nil, err
	}

	if req.FirstName == "" {
		return nil, errors.NewValidationError("first_name", "is required")
	}

	if req.LastName == "" {
		return nil, errors.NewValidationError("last_name", "is required")
	}

	tier, err := models.NewMembershipTier(req.MembershipTier)
	if err != nil {
		return nil, errors.NewValidationError("membership_tier", err.Error())
	}

	expiryDate, err := tier.GetExpiryDate()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate expiry date: %w", err)
	}

	createdAt, err := utils.NowInPacific()
	if err != nil {
		return nil, fmt.Errorf("failed to get current time: %w", err)
	}

	profile := &models.GamerProfile{
		StudentNumber:        req.StudentNumber,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		MembershipTier:       req.MembershipTier,
		Banned:               req.Banned,
		Notes:                req.Notes,
		CreatedAt:            createdAt,
		MembershipExpiryDate: expiryDate,
	}

	return s.repo.Upsert(ctx, profile)
}

func (s *gamerProfileService) DeleteProfile(ctx context.Context, studentNumber string) error {
	if err := validateStudentNumber(studentNumber); err != nil {
		return err
	}

	return s.repo.Delete(ctx, studentNumber)
}

func validateStudentNumber(studentNumber string) error {
	if !studentNumberRegex.MatchString(studentNumber) {
		return errors.NewValidationError("student_number", "must be exactly 8 digits")
	}
	return nil
}
