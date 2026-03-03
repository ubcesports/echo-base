package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ubcesports/echo-base/internal/database/sqlc"
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
	queries := sqlc.New(r.db)
	rows, err := queries.GetGamerActivity(ctx, studentNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query activities: %w", err)
	}
	return toGamerActivities(rows), nil
}

func (r *GamerActivityRepository) GetTodayActivitiesByStudent(ctx context.Context, studentNumber string) ([]models.GamerActivity, error) {
	queries := sqlc.New(r.db)
	rows, err := queries.GetTodayActivitiesByStudent(ctx, studentNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query today activities: %w", err)
	}
	return toGamerActivitiesBasic(rows), nil
}

func (r *GamerActivityRepository) GetRecentActivities(ctx context.Context, page, limit int, search string) ([]models.GamerActivity, error) {
	queries := sqlc.New(r.db)
	offset := (page - 1) * limit

	if search != "" {
		rows, err := queries.GetRecentActivitiesWithSearch(ctx, sqlc.GetRecentActivitiesWithSearchParams{
			Limit:   int64(limit),
			Offset:  int64(offset),
			Column3: sql.NullString{String: search, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query recent activities: %w", err)
		}
		return toGamerActivitiesWithNameFromSearch(rows), nil
	}

	rows, err := queries.GetRecentActivities(ctx, sqlc.GetRecentActivitiesParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query recent activities: %w", err)
	}
	return toGamerActivitiesWithName(rows), nil
}

func (r *GamerActivityRepository) Create(ctx context.Context, activity *models.GamerActivity) (*models.GamerActivity, error) {
	var activityID uuid.UUID
	var err error

	if activity.ID == "" {
		activityID = uuid.New()
	} else {
		activityID, err = uuid.Parse(activity.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID: %w", err)
		}
	}

	queries := sqlc.New(r.db)
	row, err := queries.CreateGamerActivity(ctx, sqlc.CreateGamerActivityParams{
		ID:            activityID,
		StudentNumber: activity.StudentNumber,
		PcNumber:      sql.NullInt32{Int32: int32(activity.PCNumber), Valid: true},
		Game:          sql.NullString{String: activity.Game, Valid: true},
		StartedAt:     sql.NullTime{Time: activity.StartedAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return toGamerActivityFromCreate(row), nil
}

func (r *GamerActivityRepository) UpdateEndTime(ctx context.Context, studentNumber string, pcNumber int, endedAt time.Time, execName string) (*models.GamerActivity, error) {
	queries := sqlc.New(r.db)
	row, err := queries.UpdateActivityEndTime(ctx, sqlc.UpdateActivityEndTimeParams{
		EndedAt:       sql.NullTime{Time: endedAt, Valid: true},
		ExecName:      sql.NullString{String: execName, Valid: true},
		StudentNumber: studentNumber,
		PcNumber:      sql.NullInt32{Int32: int32(pcNumber), Valid: true},
	})

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active session found for student %s on PC %d", studentNumber, pcNumber)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	return toGamerActivityFromUpdate(row), nil
}

func (r *GamerActivityRepository) GetActiveSessions(ctx context.Context) ([]models.GamerActivity, error) {
	queries := sqlc.New(r.db)
	rows, err := queries.GetActiveSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query active sessions: %w", err)
	}
	return toGamerActivitiesFromActiveSessions(rows), nil
}

/*
sqlc model conversion helpers
*/
func toGamerActivityFromCreate(row sqlc.CreateGamerActivityRow) *models.GamerActivity {
	activity := &models.GamerActivity{
		ID:            row.ID.String(),
		StudentNumber: row.StudentNumber,
		PCNumber:      int(row.PcNumber.Int32),
		Game:          row.Game.String,
		StartedAt:     row.StartedAt.Time,
	}
	if row.EndedAt.Valid {
		activity.EndedAt = &row.EndedAt.Time
	}
	if row.ExecName.Valid {
		activity.ExecName = &row.ExecName.String
	}
	return activity
}

func toGamerActivityFromUpdate(row sqlc.UpdateActivityEndTimeRow) *models.GamerActivity {
	activity := &models.GamerActivity{
		ID:            row.ID.String(),
		StudentNumber: row.StudentNumber,
		PCNumber:      int(row.PcNumber.Int32),
		Game:          row.Game.String,
		StartedAt:     row.StartedAt.Time,
	}
	if row.EndedAt.Valid {
		activity.EndedAt = &row.EndedAt.Time
	}
	if row.ExecName.Valid {
		activity.ExecName = &row.ExecName.String
	}
	return activity
}

func toGamerActivities(rows []sqlc.GetGamerActivityRow) []models.GamerActivity {
	activities := make([]models.GamerActivity, len(rows))
	for i, row := range rows {
		activities[i] = models.GamerActivity{
			ID:            row.ID.String(),
			StudentNumber: row.StudentNumber,
			PCNumber:      int(row.PcNumber.Int32),
			Game:          row.Game.String,
			StartedAt:     row.StartedAt.Time,
		}
		if row.EndedAt.Valid {
			activities[i].EndedAt = &row.EndedAt.Time
		}
		if row.ExecName.Valid {
			activities[i].ExecName = &row.ExecName.String
		}
	}
	return activities
}

func toGamerActivitiesBasic(rows []sqlc.GetTodayActivitiesByStudentRow) []models.GamerActivity {
	activities := make([]models.GamerActivity, len(rows))
	for i, row := range rows {
		activities[i] = models.GamerActivity{
			ID:            row.ID.String(),
			StudentNumber: row.StudentNumber,
			PCNumber:      int(row.PcNumber.Int32),
			Game:          row.Game.String,
			StartedAt:     row.StartedAt.Time,
		}
		if row.EndedAt.Valid {
			activities[i].EndedAt = &row.EndedAt.Time
		}
		if row.ExecName.Valid {
			activities[i].ExecName = &row.ExecName.String
		}
	}
	return activities
}

func toGamerActivitiesWithName(rows []sqlc.GetRecentActivitiesRow) []models.GamerActivity {
	activities := make([]models.GamerActivity, len(rows))
	for i, row := range rows {
		activities[i] = models.GamerActivity{
			ID:            row.ID.String(),
			StudentNumber: row.StudentNumber,
			PCNumber:      int(row.PcNumber.Int32),
			Game:          row.Game.String,
			StartedAt:     row.StartedAt.Time,
			FirstName:     &row.FirstName,
			LastName:      &row.LastName,
		}
		if row.EndedAt.Valid {
			activities[i].EndedAt = &row.EndedAt.Time
		}
		if row.ExecName.Valid {
			activities[i].ExecName = &row.ExecName.String
		}
	}
	return activities
}

func toGamerActivitiesWithNameFromSearch(rows []sqlc.GetRecentActivitiesWithSearchRow) []models.GamerActivity {
	activities := make([]models.GamerActivity, len(rows))
	for i, row := range rows {
		activities[i] = models.GamerActivity{
			ID:            row.ID.String(),
			StudentNumber: row.StudentNumber,
			PCNumber:      int(row.PcNumber.Int32),
			Game:          row.Game.String,
			StartedAt:     row.StartedAt.Time,
			FirstName:     &row.FirstName,
			LastName:      &row.LastName,
		}
		if row.EndedAt.Valid {
			activities[i].EndedAt = &row.EndedAt.Time
		}
		if row.ExecName.Valid {
			activities[i].ExecName = &row.ExecName.String
		}
	}
	return activities
}

func toGamerActivitiesFromActiveSessions(rows []sqlc.GetActiveSessionsRow) []models.GamerActivity {
	activities := make([]models.GamerActivity, len(rows))
	for i, row := range rows {
		activities[i] = models.GamerActivity{
			ID:            row.ID.String(),
			StudentNumber: row.StudentNumber,
			PCNumber:      int(row.PcNumber.Int32),
			Game:          row.Game.String,
			StartedAt:     row.StartedAt.Time,
			FirstName:     &row.FirstName,
			LastName:      &row.LastName,
		}
		if row.EndedAt.Valid {
			activities[i].EndedAt = &row.EndedAt.Time
		}
		if row.ExecName.Valid {
			activities[i].ExecName = &row.ExecName.String
		}
	}
	return activities
}
