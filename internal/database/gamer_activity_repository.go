package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ubcesports/echo-base/internal/interfaces/gamer"
	"github.com/ubcesports/echo-base/internal/models"
)

type GamerActivityRepository struct {
	db *sql.DB
}

func NewGamerActivityRepository(db *sql.DB) gamer.GamerActivityRepository {
	return &GamerActivityRepository{db: db}
}

func (r *GamerActivityRepository) GetByStudentNumber(ctx context.Context, studentNumber string) ([]models.GamerActivity, error) {
	query := `
		SELECT id, student_number, pc_number, game, started_at, ended_at, exec_name
		FROM gamer_activity
		WHERE student_number = $1
		ORDER BY started_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, studentNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query activities: %w", err)
	}
	defer rows.Close()

	activities := []models.GamerActivity{}
	for rows.Next() {
		var activity models.GamerActivity
		err := rows.Scan(
			&activity.ID,
			&activity.StudentNumber,
			&activity.PCNumber,
			&activity.Game,
			&activity.StartedAt,
			&activity.EndedAt,
			&activity.ExecName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activities: %w", err)
	}

	return activities, nil
}

func (r *GamerActivityRepository) GetTodayActivitiesByStudent(ctx context.Context, studentNumber string) ([]models.GamerActivity, error) {
	query := `
		SELECT ga.id, ga.student_number, ga.pc_number, ga.game, ga.started_at, ga.ended_at, ga.exec_name
		FROM gamer_activity ga
		JOIN gamer_profile gp ON ga.student_number = gp.student_number
		WHERE ga.student_number = $1
		AND gp.membership_tier = 1
		AND DATE(ga.started_at AT TIME ZONE 'America/Los_Angeles') = DATE(NOW() AT TIME ZONE 'America/Los_Angeles')
	`

	rows, err := r.db.QueryContext(ctx, query, studentNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query today activities: %w", err)
	}
	defer rows.Close()

	activities := []models.GamerActivity{}
	for rows.Next() {
		var activity models.GamerActivity
		err := rows.Scan(
			&activity.ID,
			&activity.StudentNumber,
			&activity.PCNumber,
			&activity.Game,
			&activity.StartedAt,
			&activity.EndedAt,
			&activity.ExecName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activities: %w", err)
	}

	return activities, nil
}

func (r *GamerActivityRepository) GetRecentActivities(ctx context.Context, page, limit int, search string) ([]models.GamerActivity, error) {
	offset := (page - 1) * limit

	query := `
		SELECT ga.id, ga.student_number, ga.pc_number, ga.game, ga.started_at, ga.ended_at, ga.exec_name,
		       gp.first_name, gp.last_name
		FROM gamer_activity ga
		JOIN gamer_profile gp ON ga.student_number = gp.student_number
	`

	args := []interface{}{limit, offset}

	if search != "" {
		query += `
			WHERE ga.student_number % $3
			   OR gp.first_name % $3
			   OR gp.last_name % $3
			   OR ga.game % $3
			   OR ga.exec_name % $3
		`
		args = append(args, search)
	}

	query += `
		ORDER BY ga.started_at DESC NULLS LAST
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent activities: %w", err)
	}
	defer rows.Close()

	activities := []models.GamerActivity{}
	for rows.Next() {
		var activity models.GamerActivity
		err := rows.Scan(
			&activity.ID,
			&activity.StudentNumber,
			&activity.PCNumber,
			&activity.Game,
			&activity.StartedAt,
			&activity.EndedAt,
			&activity.ExecName,
			&activity.FirstName,
			&activity.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activities: %w", err)
	}

	return activities, nil
}

func (r *GamerActivityRepository) Create(ctx context.Context, activity *models.GamerActivity) (*models.GamerActivity, error) {
	if activity.ID == "" {
		activity.ID = uuid.New().String()
	}

	query := `
		INSERT INTO gamer_activity (id, student_number, pc_number, game, started_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, student_number, pc_number, game, started_at, ended_at, exec_name
	`

	result := &models.GamerActivity{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		activity.ID,
		activity.StudentNumber,
		activity.PCNumber,
		activity.Game,
		activity.StartedAt,
	).Scan(
		&result.ID,
		&result.StudentNumber,
		&result.PCNumber,
		&result.Game,
		&result.StartedAt,
		&result.EndedAt,
		&result.ExecName,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return result, nil
}

func (r *GamerActivityRepository) UpdateEndTime(ctx context.Context, studentNumber string, pcNumber int, endedAt time.Time, execName string) (*models.GamerActivity, error) {
	query := `
		UPDATE gamer_activity
		SET ended_at = $1, exec_name = $2
		WHERE student_number = $3
		AND pc_number = $4
		AND ended_at IS NULL
		RETURNING id, student_number, pc_number, game, started_at, ended_at, exec_name
	`

	result := &models.GamerActivity{}
	err := r.db.QueryRowContext(ctx, query, endedAt, execName, studentNumber, pcNumber).Scan(
		&result.ID,
		&result.StudentNumber,
		&result.PCNumber,
		&result.Game,
		&result.StartedAt,
		&result.EndedAt,
		&result.ExecName,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active session found for student %s on PC %d", studentNumber, pcNumber)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	return result, nil
}

func (r *GamerActivityRepository) GetActiveSessions(ctx context.Context) ([]models.GamerActivity, error) {
	query := `
		SELECT ga.id, ga.student_number, ga.pc_number, ga.game, ga.started_at, ga.ended_at, ga.exec_name,
		       gp.first_name, gp.last_name
		FROM gamer_activity ga
		JOIN gamer_profile gp ON ga.student_number = gp.student_number
		WHERE ga.ended_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active sessions: %w", err)
	}
	defer rows.Close()

	activities := []models.GamerActivity{}
	for rows.Next() {
		var activity models.GamerActivity
		err := rows.Scan(
			&activity.ID,
			&activity.StudentNumber,
			&activity.PCNumber,
			&activity.Game,
			&activity.StartedAt,
			&activity.EndedAt,
			&activity.ExecName,
			&activity.FirstName,
			&activity.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activities: %w", err)
	}

	return activities, nil
}
