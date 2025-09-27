package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/models"
	"github.com/ubcesports/echo-base/internal/repositories"
)

type Handler struct {
	Config *config.Config
	DB     *pgxpool.Pool
}

func (h *Handler) GetGamerActivityByStudent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	studentNumber := r.PathValue("student_number")

	var activities []models.GamerActivity
	query := repositories.BuildGamerActivityByStudentQuery()

	if err := pgxscan.Select(ctx, h.DB, &activities, query, studentNumber); err != nil {
		return NewHTTPError(http.StatusInternalServerError, "database query failed", err)
	}

	if len(activities) == 0 {
		return NewHTTPError(http.StatusNotFound, "student not found")
	}

	return WriteJSON(w, http.StatusOK, activities)
}

func (h *Handler) GetGamerActivityByTierOneStudentToday(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	studentNumber := r.PathValue("student_number")
	now := time.Now().In(h.Config.Location)

	var activities []models.GamerActivity
	query := repositories.BuildGamerActivityByTierOneStudentTodayQuery()

	if err := pgxscan.Select(ctx, h.DB, &activities, query, studentNumber, now); err != nil {
		return NewHTTPError(http.StatusInternalServerError, "database query failed", err)
	}

	return WriteJSON(w, http.StatusOK, activities)
}

func (h *Handler) GetGamerActivity(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Pagination parameters
	page := QueryStringToInt(r, "page", 1)
	limit := QueryStringToInt(r, "limit", 10)
	search := r.URL.Query().Get("search")
	offset := (page - 1) * limit

	var activities []models.GamerActivityWithName
	query, args := repositories.BuildGamerActivityRecentQuery(limit, offset, search)

	if err := pgxscan.Select(ctx, h.DB, &activities, query, args...); err != nil {
		return NewHTTPError(http.StatusInternalServerError, "database query failed", err)
	}

	return WriteJSON(w, http.StatusOK, activities)
}

func (h *Handler) AddGamerActivity(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	started_at := time.Now().In(h.Config.Location)
	query := repositories.BuildInsertGamerActivityQuery()

	if !RequireJSONContentType(r) {
		return NewHTTPError(http.StatusBadRequest, "invalid content type")
	}

	var ga models.GamerActivityWithName
	if err := DecodeJSONBody(r, &ga); err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid request body", err)
	}

	var m models.MembershipResult
	checkMemberQuery := repositories.BuildCheckMemberQuery()
	if err := pgxscan.Get(ctx, h.DB, &m, checkMemberQuery, ga.StudentNumber); err != nil {
		msg := fmt.Sprintf("foreign key %s not found", ga.StudentNumber)
		return NewHTTPError(http.StatusNotFound, msg, err)
	}

	today := time.Now().In(h.Config.Location).Truncate(24 * time.Hour)
	expiryDate := m.MembershipExpiryDate.In(h.Config.Location).Truncate(24 * time.Hour)

	if today.After(expiryDate) {
		msg := fmt.Sprintf(
			"Membership expired on %s. Please ask the user to purchase a new membership. "+
				"If the member has already purchased a new membership for this year please verify via Showpass then create a new profile for them.",
			expiryDate.Format("2006-01-02"),
		)
		return NewHTTPError(http.StatusForbidden, msg)
	}

	var gamerActivities []models.GamerActivity
	if err := pgxscan.Select(ctx, h.DB, &gamerActivities, query, ga.StudentNumber, ga.PCNumber, ga.Game, started_at); err != nil {
		return NewHTTPError(http.StatusInternalServerError, "database query failed", err)
	}

	return WriteJSON(w, http.StatusCreated, gamerActivities[0])
}

func (h *Handler) UpdateGamerActivity(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ended_at := time.Now().In(h.Config.Location)

	studentNumber := r.PathValue("student_number")

	if !RequireJSONContentType(r) {
		return NewHTTPError(http.StatusBadRequest, "invalid content type")
	}

	var rb models.UpdateGamerActivityRequest
	if err := DecodeJSONBody(r, &rb); err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid request body", err)
	}

	var activities []models.GamerActivity
	query := repositories.BuildUpdateGamerActivityQuery()
	if err := pgxscan.Select(ctx, h.DB, &activities, query, ended_at, studentNumber, rb.PCNumber, rb.ExecName); err != nil {
		return NewHTTPError(http.StatusInternalServerError, "database query failed", err)
	}

	if len(activities) == 0 {
		return NewHTTPError(http.StatusNotFound, "student not active")
	}

	return WriteJSON(w, http.StatusCreated, activities[0])
}

func (h *Handler) GetAllActivePCs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	query := repositories.BuildGetAllActivePCsQuery()
	var activePCs []models.ActivePC

	if err := pgxscan.Select(ctx, h.DB, &activePCs, query); err != nil {
		return NewHTTPError(http.StatusInternalServerError, "database query failed", err)
	}

	return WriteJSON(w, http.StatusOK, activePCs)
}
