package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/database"
)

const FK_VIOLATION = "23503"

type Handler struct {
	Config *config.Config
}

type GamerActivity struct {
	StudentNumber string     `json:"student_number"`
	PCNumber      int        `json:"pc_number"`
	Game          string     `json:"game"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`  // optional
	ExecName      *string    `json:"exec_name,omitempty"` // optional
}

type GamerActivityWithName struct {
	GamerActivity
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ActivePC struct {
	GamerActivity
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	MembershipTier int       `json:"membership_tier"`
	Banned         bool      `json:"banned"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}

type UpdateGamerActivityRequest struct {
	PCNumber int    `json:"pc_number"`
	ExecName string `json:"exec_name"`
}

func ScanGamerActivity(rows *sql.Rows) ([]GamerActivity, error) {
	var result []GamerActivity
	for rows.Next() {
		var a GamerActivity
		if err := rows.Scan(&a.StudentNumber, &a.PCNumber, &a.Game, &a.StartedAt, &a.EndedAt, &a.ExecName); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func ScanGamerActivityWithName(rows *sql.Rows) ([]GamerActivityWithName, error) {
	var result []GamerActivityWithName
	for rows.Next() {
		var a GamerActivityWithName
		if err := rows.Scan(&a.StudentNumber, &a.PCNumber, &a.Game, &a.StartedAt, &a.EndedAt, &a.ExecName, &a.FirstName, &a.LastName); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func ScanActivePC(rows *sql.Rows) ([]ActivePC, error) {
	var result []ActivePC
	for rows.Next() {
		var a ActivePC
		if err := rows.Scan(&a.StudentNumber, &a.PCNumber, &a.Game, &a.StartedAt, &a.EndedAt, &a.ExecName, &a.FirstName, &a.LastName, &a.MembershipTier, &a.Banned, &a.Notes, &a.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Converts string query param to integer with default value
func QueryStringToInt(r *http.Request, key string, defaultVal int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return i
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, msg string, err error) {
	http.Error(w, fmt.Sprintf("%s: %v", msg, err), status)
}

// Ensures the request has Content-Type application/json
func RequireJSONContentType(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		return true
	}

	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	return mediaType == "application/json"
}

// Decodes JSON body into the provided struct
func DecodeJSONBody(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if dec.More() {
		return fmt.Errorf("body must contain only one JSON object")
	}
	return nil
}

/**
 * @api {get} /activity/:student_number Get Gamer Activity for specific student
 * @apiName GetGamerActivityByStudent
 * @apiGroup Activity
 *
 * @apiParam {String} student_number Student's unique number.
 *
 * @apiSuccess {Object} gamer_activity Gamer activity object.
 * @apiSuccess {String} gamer_activity.student_number Student number, 8 digit integer.
 * @apiSuccess {Number} gamer_activity.pc_number PC number.
 * @apiSuccess {String} gamer_activity.game Game name.
 * @apiSuccess {String} gamer_activity.started_at Datetime when the activity started.
 * @apiSuccess {String} gamer_activity.ended_at Datetime when the activity ended.
 * @apiSuccess {string} gamer_activity.exec_name Exec that ended the activity.
 *
 * @apiError {String} 500 Server error.
 */

func (h *Handler) GetGamerActivityByStudent(w http.ResponseWriter, r *http.Request) {

	// Extract dynamic path param
	studentNumber := r.PathValue("student_number")

	// Build query
	query := fmt.Sprintf(`
		SELECT *
		FROM %[1]s.gamer_activity
		WHERE student_number = $1
	`, h.Config.Schema)

	// Execute query and close connection
	rows, err := database.DB.Query(query, studentNumber)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error finding gamer activity: ", err)
		return
	}
	defer rows.Close()

	// Scan results
	response, err := ScanGamerActivity(rows)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error scanning rows", err)
		return
	}

	if len(response) == 0 {
		WriteError(w, http.StatusNotFound, "Student not found", nil)
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

/**
 * @api {get} /activity/today/:student_number Get Gamer Tier One Member Activity today for specific student
 * @apiName GetGamerActivityByTierOneStudentToday
 * @apiGroup Activity
 *
 * @apiParam {String} student_number Student's unique number.
 *
 * @apiSuccess {Object} gamer_activity Gamer activity object.
 * @apiSuccess {String} gamer_activity.student_number Student number, 8 digit integer.
 * @apiSuccess {Number} gamer_activity.pc_number PC number.
 * @apiSuccess {String} gamer_activity.game Game name.
 * @apiSuccess {String} gamer_activity.started_at Datetime when the activity started.
 * @apiSuccess {String} gamer_activity.ended_at Datetime when the activity ended.
 * @apiSuccess {string} gamer_activity.exec_name Exec that ended the activity.
 *
 * @apiError {String} 500 Server error.
 */
func (h *Handler) GetGamerActivityByTierOneStudentToday(w http.ResponseWriter, r *http.Request) {
	studentNumber := r.PathValue("student_number")

	now := time.Now().In(h.Config.Location)

	query := fmt.Sprintf(`
		SELECT ga.*
		FROM %[1]s.gamer_activity ga
		JOIN %[1]s.gamer_profile gp 
		ON ga.student_number = gp.student_number
		WHERE ga.student_number = $1
		AND gp.membership_tier = 1
		AND DATE(ga.started_at::timestamp) = DATE($2::timestamp)
	`, h.Config.Schema)

	rows, err := database.DB.Query(query, studentNumber, now)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error checking Tier one member sign in: ", err)
		return
	}
	defer rows.Close()

	response, err := ScanGamerActivity(rows)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error scanning rows", err)
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

// Gamer Activity with Profile info
/**
 * @api {get} /activity/all/recent Get Gamer Activity
 * @apiName GetGamerActivity
 * @apiGroup Activity
 *
 * @apiParam {Number} page Page number.
 * @apiParam {Number} limit Limit Number of results per page.
 * @apiParam {String} search Search query.
 *
 * @apiSuccess {Object} gamer_activity Gamer activty object with additional profile fields.
 * @apiSuccess {String} gamer_activity.student_number Student number, 8 digit integer.
 * @apiSuccess {Number} gamer_activity.pc_number PC number.
 * @apiSuccess {String} gamer_activity.game Game name.
 * @apiSuccess {String} gamer_activity.started_at Datetime when the activity started.
 * @apiSuccess {String} gamer_activity.ended_at Datetime when the activity ended.
 * @apiSuccess {string} gamer_activity.exec_name Exec that ended the activity.
 *
 * @apiError {String} 500 Server error.
 */
func (h *Handler) GetGamerActivity(w http.ResponseWriter, r *http.Request) {
	page := QueryStringToInt(r, "page", 1)
	limit := QueryStringToInt(r, "limit", 10)
	search := r.URL.Query().Get("search")

	offset := (page - 1) * limit

	// Base query
	query := fmt.Sprintf(`
		SELECT ga.*, gp.first_name, gp.last_name
		FROM %[1]s.gamer_activity ga
		JOIN %[1]s.gamer_profile gp 
		ON ga.student_number = gp.student_number
	`, h.Config.Schema)

	args := []interface{}{limit, offset}

	if search != "" {
		query += `
			WHERE ga.student_number ILIKE $3
			OR gp.first_name ILIKE $3
			OR gp.last_name ILIKE $3
			OR ga.game ILIKE $3
			OR ga.exec_name ILIKE $3
			OR TO_CHAR(ga.started_at, 'YYYY-MM-DD') ILIKE $3
		`
		args = append(args, "%"+search+"%")
	}
	query += `ORDER BY ga.started_at DESC NULLS LAST LIMIT $1 OFFSET $2`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error finding recent activity: ", err)
		return
	}
	defer rows.Close()

	response, err := ScanGamerActivityWithName(rows)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error scanning rows", err)
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

//Gamer Activity only
/**
 * @api {post} /activity Add Gamer Activity
 * @apiName AddGamerActivity
 * @apiGroup Activity
 *
 * @apiParam {String} student_number Student number, 8 digit integer.
 * @apiParam {String} pc_number PC number.
 * @apiParam {String} game Game name.
 * @apiParam {Number} started_at Date when the activity started.
 *
 * @apiSuccess {Object} gamer_activity Gamer activity object.
 * @apiSuccess {String} gamer_activity.student_number Student number, 8 digit integer.
 * @apiSuccess {Number} gamer_activity.pc_number PC number.
 * @apiSuccess {String} gamer_activity.game Game name.
 * @apiSuccess {String} gamer_activity.started_at Datetime when the activity started.
 * @apiSuccess {Null} gamer_activity.ended_at Datetime will be null.
 * @apiSuccess {Null} gamer_activity.exec_name will be null.
 *
 * @apiError {String} 500 Server error.
 * @apiError {String} 404 Foreign key not found.
 */
func (h *Handler) AddGamerActivity(w http.ResponseWriter, r *http.Request) {

	if !RequireJSONContentType(r) {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var ga GamerActivity
	if err := DecodeJSONBody(r, &ga); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	now := time.Now().In(h.Config.Location)

	query := fmt.Sprintf(`
		INSERT INTO %[1]s.gamer_activity 
		(student_number, pc_number, game, started_at)
		VALUES ($1, $2, $3, $4)
		RETURNING *`, h.Config.Schema)

	rows, err := database.DB.Query(query, ga.StudentNumber, ga.PCNumber, ga.Game, now)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == FK_VIOLATION {
			WriteError(w, http.StatusNotFound, "Foreign key violation: Gamer profile not found", err)
			return
		}
		WriteError(w, http.StatusInternalServerError, "Error creating activity", err)
		return
	}
	defer rows.Close()

	activities, err := ScanGamerActivity(rows)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error scanning rows", err)
		return
	}

	WriteJSON(w, http.StatusCreated, activities[0])
}

//Gamer Activity only
/**
 * @api {patch} /activity/update/:student_number Update Gamer Activity End Time
 * @apiName UpdateGamerActivity
 * @apiGroup Activity
 *
 * @apiParam {String} student_number Student number, 8 digit integer.
 *
 * @apiSuccess {Object} gamer_activity Gamer activity object.
 * @apiSuccess {String} gamer_activity.student_number Student number, 8 digit integer.
 * @apiSuccess {Number} gamer_activity.pc_number PC number.
 * @apiSuccess {String} gamer_activity.game Game name.
 * @apiSuccess {String} gamer_activity.started_at Datetime when the activity started.
 * @apiSuccess {String} gamer_activity.ended_at Datetime when the activity ended.
 * @apiSuccess {string} gamer_activity.exec_name Exec that ended the activity.
 *
 * @apiError {String} 500 Internal server error.
 */
func (h *Handler) UpdateGamerActivity(w http.ResponseWriter, r *http.Request) {
	now := time.Now().In(h.Config.Location)

	studentNumber := r.PathValue("student_number")

	if !RequireJSONContentType(r) {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var rb UpdateGamerActivityRequest
	if err := DecodeJSONBody(r, &rb); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	query := fmt.Sprintf(`
		UPDATE %[1]s.gamer_activity
		SET ended_at = $1, exec_name = $4
		WHERE student_number = $2 AND pc_number = $3 AND ended_at IS NULL
		RETURNING *
	`, h.Config.Schema)

	rows, err := database.DB.Query(query, now, studentNumber, rb.PCNumber, rb.ExecName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error finding recent activity: ", err)
		return
	}
	defer rows.Close()

	activities, err := ScanGamerActivity(rows)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error scanning rows", err)
		return
	}

	if len(activities) == 0 {
		WriteError(w, http.StatusNotFound, "Student not active", nil)
		return
	}

	WriteJSON(w, http.StatusOK, activities[0])
}

/**
 * @api {get} /activity/get-active-pcs Gets all active PCs
 * @apiName GetAllActivePCs
 * @apiGroup Activity
 *
 * @apiParam {Null} No parameters required.
 *
 * @apiSuccess {Object} gamer_activity Gamer activty object with additional profile fields.
 * @apiSuccess {String} gamer_activity.student_number Student number, 8 digit integer.
 * @apiSuccess {Number} gamer_activity.pc_number PC number.
 * @apiSuccess {String} gamer_activity.game Game name.
 * @apiSuccess {String} gamer_activity.started_at Datetime when the activity started.
 * @apiSuccess {String} gamer_activity.ended_at Datetime when the activity ended.
 * @apiSuccess {string} gamer_activity.exec_name Exec that ended the activity.
 *
 * @apiError {String} 500 Internal server error.
 */
func (h *Handler) GetAllActivePCs(w http.ResponseWriter, r *http.Request) {
	query := fmt.Sprintf(`
		SELECT ga.*, gp.first_name, gp.last_name, gp.membership_tier,
			gp.banned, gp.notes, gp.created_at
		FROM %[1]s.gamer_activity ga
		JOIN %[1]s.gamer_profile gp 
		ON ga.student_number = gp.student_number
		WHERE ga.ended_at IS NULL
	`, h.Config.Schema)

	rows, err := database.DB.Query(query)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error finding recent activity: ", err)
		return
	}
	defer rows.Close()

	response, err := ScanActivePC(rows)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error scanning rows", err)
		return
	}

	WriteJSON(w, http.StatusOK, response)
}
