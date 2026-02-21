package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ubcesports/echo-base/internal/database/sqlc"
	"github.com/ubcesports/echo-base/internal/errors"
	"github.com/ubcesports/echo-base/internal/interfaces/gamer"
	"github.com/ubcesports/echo-base/internal/models"
)

type GamerProfileRepository struct {
	db *sql.DB
}

func NewGamerProfileRepository(db *sql.DB) gamer.GamerProfileRepository {
	return &GamerProfileRepository{db: db}
}

func (r *GamerProfileRepository) GetByStudentNumber(ctx context.Context, studentNumber string) (*models.GamerProfile, error) {
	queries := sqlc.New(r.db)
	row, err := queries.GetGamerProfile(ctx, studentNumber)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("student", studentNumber)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return toGamerProfile(row), nil
}

func (r *GamerProfileRepository) Upsert(ctx context.Context, profile *models.GamerProfile) (*models.GamerProfile, error) {
	queries := sqlc.New(r.db)

	row, err := queries.UpsertGamerProfile(ctx, toUpsertParams(profile))
	if err != nil {
		return nil, fmt.Errorf("failed to upsert profile: %w", err)
	}

	return toGamerProfile(row), nil
}

func (r *GamerProfileRepository) Delete(ctx context.Context, studentNumber string) error {
	queries := sqlc.New(r.db)
	rows, err := queries.DeleteGamerProfile(ctx, studentNumber)

	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("student", studentNumber)
	}

	return nil
}

func (r *GamerProfileRepository) CheckMembershipValidity(ctx context.Context, studentNumber string) (tier int, expiryDate *time.Time, err error) {
	queries := sqlc.New(r.db)
	result, err := queries.CheckMembershipValidity(ctx, studentNumber)

	if err == sql.ErrNoRows {
		return 0, nil, fmt.Errorf("student %s not found", studentNumber)
	}
	if err != nil {
		return 0, nil, fmt.Errorf("failed to check membership: %w", err)
	}

	return int(result.MembershipTier), &result.MembershipExpiryDate.Time, nil
}

/*
sqlc model conversion helpers
*/
func toUpsertParams(p *models.GamerProfile) sqlc.UpsertGamerProfileParams {
	return sqlc.UpsertGamerProfileParams{
		FirstName:            p.FirstName,
		LastName:             p.LastName,
		StudentNumber:        p.StudentNumber,
		MembershipTier:       int32(p.MembershipTier),
		Banned:               nullBool(p.Banned),
		Notes:                nullString(p.Notes),
		CreatedAt:            sql.NullTime{Valid: true, Time: p.CreatedAt},
		MembershipExpiryDate: nullTime(p.MembershipExpiryDate),
	}
}

func toGamerProfile(row sqlc.GamerProfile) *models.GamerProfile {
	profile := &models.GamerProfile{
		StudentNumber:  row.StudentNumber,
		FirstName:      row.FirstName,
		LastName:       row.LastName,
		MembershipTier: int(row.MembershipTier),
	}
	if row.ID.Valid {
		profile.ID = row.ID.UUID.String()
	}
	if row.Banned.Valid {
		profile.Banned = &row.Banned.Bool
	}
	if row.Notes.Valid {
		profile.Notes = &row.Notes.String
	}
	if row.CreatedAt.Valid {
		profile.CreatedAt = row.CreatedAt.Time
	}
	if row.MembershipExpiryDate.Valid {
		profile.MembershipExpiryDate = &row.MembershipExpiryDate.Time
	}

	return profile
}

func nullBool(v *bool) sql.NullBool {
	if v == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Valid: true, Bool: *v}
}

func nullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{Valid: true, String: *v}
}

func nullTime(v *time.Time) sql.NullTime {
	if v == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true, Time: *v}
}
