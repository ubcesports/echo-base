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
	profile, err := queries.GetStudent(ctx, studentNumber)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("student", studentNumber)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return toGamerProfile(profile), nil
}

func (r *GamerProfileRepository) Upsert(ctx context.Context, profile *models.GamerProfile) (*models.GamerProfile, error) {
	query := `
		INSERT INTO gamer_profile (first_name, last_name, student_number, membership_tier, banned, notes, created_at, membership_expiry_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (student_number)
		DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			membership_tier = EXCLUDED.membership_tier,
			banned = EXCLUDED.banned,
			notes = EXCLUDED.notes,
			created_at = EXCLUDED.created_at,
			membership_expiry_date = EXCLUDED.membership_expiry_date
		RETURNING id, student_number, first_name, last_name, membership_tier, banned, notes, created_at, membership_expiry_date
	`

	result := &models.GamerProfile{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		profile.FirstName,
		profile.LastName,
		profile.StudentNumber,
		profile.MembershipTier,
		profile.Banned,
		profile.Notes,
		profile.CreatedAt,
		profile.MembershipExpiryDate,
	).Scan(
		&result.ID,
		&result.StudentNumber,
		&result.FirstName,
		&result.LastName,
		&result.MembershipTier,
		&result.Banned,
		&result.Notes,
		&result.CreatedAt,
		&result.MembershipExpiryDate,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert profile: %w", err)
	}

	return result, nil
}

func (r *GamerProfileRepository) Delete(ctx context.Context, studentNumber string) error {
	query := `DELETE FROM gamer_profile WHERE student_number = $1`

	result, err := r.db.ExecContext(ctx, query, studentNumber)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("student", studentNumber)
	}

	return nil
}

func (r *GamerProfileRepository) CheckMembershipValidity(ctx context.Context, studentNumber string) (tier int, expiryDate *time.Time, err error) {
	query := `
		SELECT membership_tier, membership_expiry_date
		FROM gamer_profile
		WHERE student_number = $1
	`

	err = r.db.QueryRowContext(ctx, query, studentNumber).Scan(&tier, &expiryDate)
	if err == sql.ErrNoRows {
		return 0, nil, fmt.Errorf("student %s not found", studentNumber)
	}
	if err != nil {
		return 0, nil, fmt.Errorf("failed to check membership: %w", err)
	}

	return tier, expiryDate, nil
}

func toGamerProfile(row sqlc.GetStudentRow) *models.GamerProfile {
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
