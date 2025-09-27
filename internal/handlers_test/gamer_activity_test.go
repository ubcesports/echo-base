package handlers_test

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/handlers"
	"github.com/ubcesports/echo-base/internal/models"
	"github.com/ubcesports/echo-base/internal/tests"
)

var testRouter http.Handler

var (
	student1 = "11223344"
	student2 = "87654321"
	student3 = "63347439"

	payloadStudent1 = map[string]interface{}{
		"student_number": student1,
		"pc_number":      1,
		"game":           "Valorant",
	}
	payloadStudent2 = map[string]interface{}{
		"student_number": student2,
		"pc_number":      2,
		"game":           "Valorant",
	}
	payloadStudent3 = map[string]interface{}{
		"student_number": student3,
		"pc_number":      3,
		"game":           "CS:GO",
	}
	badPayload = map[string]interface{}{
		"student_number": "09090909",
		"pc_number":      1,
		"game":           "Valorant",
	}
)

func TestMain(m *testing.M) {
	tests.SetupTestDB()
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	testConfig := config.LoadConfig()
	testHandler := &handlers.Handler{DB: database.DB, Config: testConfig}
	mux := http.NewServeMux()

	mux.HandleFunc("/activity/{student_number}", handlers.Wrap(testHandler.GetGamerActivityByStudent))
	mux.HandleFunc("/activity/today/{student_number}", handlers.Wrap(testHandler.GetGamerActivityByTierOneStudentToday))
	mux.HandleFunc("/activity/all/recent", handlers.Wrap(testHandler.GetGamerActivity))
	mux.HandleFunc("/activity", handlers.Wrap(testHandler.AddGamerActivity))
	mux.HandleFunc("/activity/update/{student_number}", handlers.Wrap(testHandler.UpdateGamerActivity))
	mux.HandleFunc("/activity/get-active-pcs", handlers.Wrap(testHandler.GetAllActivePCs))

	testRouter = mux

	exitCode := m.Run()
	os.Exit(exitCode)
}

func BeforeEach(t *testing.T) {
	ctx := context.Background()

	_, err := database.DB.Exec(ctx, "TRUNCATE TABLE gamer_activity;")
	require.NoError(t, err)

	_, err = database.DB.Exec(ctx, "TRUNCATE TABLE gamer_profile CASCADE;")
	require.NoError(t, err)

	_, err = database.DB.Exec(ctx, `
        INSERT INTO gamer_profile (first_name, last_name, student_number, membership_expiry_date, membership_tier)
        VALUES
        ('John','Doe','11223344','2030-09-18',1),
        ('Jane','Doe','87654321','2030-09-18',2),
		('Jeffrey','Doe','63347439','2020-09-18',1);
    `)
	require.NoError(t, err)
}

func AfterEach(t *testing.T) {
	ctx := context.Background()
	_, err := database.DB.Exec(ctx, "TRUNCATE TABLE gamer_activity;")
	require.NoError(t, err)
}

func postActivity(t *testing.T, payload map[string]interface{}, out interface{}) {
	rr := tests.ExecuteTestRequest(t, testRouter, http.MethodPost, "/activity", payload)
	tests.AssertResponse(t, rr, http.StatusCreated, out)
}

func updateActivity(t *testing.T, student string, payload map[string]interface{}, out interface{}, expectedStatus int) {
	url := "/activity/update/" + student
	rr := tests.ExecuteTestRequest(t, testRouter, http.MethodPost, url, payload)
	tests.AssertResponse(t, rr, expectedStatus, out)
}

func getJSON(t *testing.T, url string, out interface{}) {
	rr := tests.ExecuteTestRequest(t, testRouter, http.MethodGet, url, nil)
	tests.AssertResponse(t, rr, http.StatusOK, out)
}

func TestAddGamerActivity(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })

	t.Run("it should add an activity", func(t *testing.T) {
		var activity models.GamerActivityWithName
		postActivity(t, payloadStudent1, &activity)

		require.Equal(t, student1, activity.StudentNumber)
		require.Equal(t, 1, activity.PCNumber)
		require.Equal(t, "Valorant", activity.Game)
		require.Nil(t, activity.ExecName)
	})

	t.Run("it should return 404 if FK user does not exist", func(t *testing.T) {
		rr := tests.ExecuteTestRequest(t, testRouter, http.MethodPost, "/activity", badPayload)
		tests.AssertResponse(t, rr, http.StatusNotFound, nil)
		require.Equal(t, "foreign key 09090909 not found\n", rr.Body.String())
	})

	t.Run("it should complain for adding expired user", func(t *testing.T) {
		rr := tests.ExecuteTestRequest(t, testRouter, http.MethodPost, "/activity", payloadStudent3)
		tests.AssertResponse(t, rr, http.StatusForbidden, nil)
		require.Equal(t,
			"Membership expired on 2020-09-17. Please ask the user to purchase a new membership. If the member has already purchased a new membership for this year please verify via Showpass then create a new profile for them.\n",
			rr.Body.String())
	})
}

func TestUpdateGamerActivity(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })

	t.Run("should patch an existing activity", func(t *testing.T) {
		var activity models.GamerActivityWithName
		postActivity(t, payloadStudent1, &activity)
		require.Nil(t, activity.EndedAt)

		updatePayload := map[string]interface{}{"pc_number": 1, "exec_name": "John"}
		var updated models.GamerActivity
		updateActivity(t, student1, updatePayload, &updated, http.StatusCreated)

		require.Equal(t, student1, updated.StudentNumber)
		require.Equal(t, "Valorant", updated.Game)
		require.Equal(t, "John", *updated.ExecName)
		require.NotNil(t, updated.EndedAt)
	})

	t.Run("should return 404 if student does not have active activity", func(t *testing.T) {
		updatePayload := map[string]interface{}{"pc_number": 1, "exec_name": "John"}
		updateActivity(t, student1, updatePayload, nil, http.StatusNotFound)
	})
}

func TestGetGamerActivity(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })
	t.Run("should get no data", func(t *testing.T) {
		var activities []models.GamerActivityWithName
		getJSON(t, "/activity/all/recent?page=1&limit=2&search=Valorant", &activities)
		require.Len(t, activities, 0)
	})

	t.Run("should get a gamer activities", func(t *testing.T) {
		postActivity(t, payloadStudent1, nil)
		postActivity(t, payloadStudent2, nil)

		var activities []models.GamerActivityWithName
		getJSON(t, "/activity/all/recent?page=1&limit=2&search=Valorant", &activities)

		require.LessOrEqual(t, len(activities), 2)
		for _, a := range activities {
			require.Contains(t, a.Game, "Valorant")
		}
	})

	t.Run("should get no data from invalid params", func(t *testing.T) {
		var activities []models.GamerActivityWithName
		getJSON(t, "/activity/all/recent?page=1&limit=2&search=ClashRoyale", &activities)
		require.Len(t, activities, 0)
	})
}

func TestGetGamerActivityByStudent(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })

	t.Run("should get gamer activites for specific student", func(t *testing.T) {
		postActivity(t, payloadStudent1, nil)
		postActivity(t, payloadStudent2, nil)

		var activities []models.GamerActivity
		getJSON(t, "/activity/"+student2, &activities)

		for _, a := range activities {
			require.Equal(t, student2, a.StudentNumber)
		}
	})
}
func TestGetGamerActivityByTierOneStudentToday(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })

	t.Run("tier 1 should see no activity before check-in, and one after", func(t *testing.T) {
		var activities []models.GamerActivityWithName
		getJSON(t, "/activity/today/"+student1, &activities)
		require.Len(t, activities, 0)

		payload := map[string]interface{}{"student_number": student1, "pc_number": 2, "game": "Valorant"}
		postActivity(t, payload, nil)

		updatePayload := map[string]interface{}{"pc_number": 2, "exec_name": "John"}
		updateActivity(t, student1, updatePayload, nil, http.StatusCreated)

		var updated []models.GamerActivityWithName
		getJSON(t, "/activity/today/"+student1, &updated)

		require.Len(t, updated, 1)
		require.Equal(t, student1, updated[0].StudentNumber)
		require.Equal(t, "John", *updated[0].ExecName)
		require.NotNil(t, updated[0].EndedAt)
	})
}

func TestGetGamerActivityByTierTwoStudentToday(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })

	t.Run("tier 2 should see no activity from today", func(t *testing.T) {
		var activities []models.GamerActivityWithName
		getJSON(t, "/activity/today/"+student2, &activities)
		require.Len(t, activities, 0)

		postActivity(t, payloadStudent2, nil)

		updatePayload := map[string]interface{}{"pc_number": 2, "exec_name": "John"}
		updateActivity(t, student2, updatePayload, nil, http.StatusCreated)

		var updated []models.GamerActivityWithName
		getJSON(t, "/activity/today/"+student2, &updated)
		require.Len(t, updated, 0)
	})
}

func TestTierTwoMultipleCheckIns(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })

	t.Run("tier 2 can check in twice", func(t *testing.T) {
		var first models.GamerActivity
		postActivity(t, payloadStudent2, &first)
		require.Equal(t, student2, first.StudentNumber)
		require.Equal(t, 2, first.PCNumber)
		require.Equal(t, "Valorant", first.Game)

		updatePayload := map[string]interface{}{"pc_number": 2, "exec_name": "Jane"}
		var updated models.GamerActivity
		updateActivity(t, student2, updatePayload, &updated, http.StatusCreated)
		require.NotNil(t, updated.EndedAt)
		require.NotNil(t, updated.ExecName)
		require.Equal(t, "Jane", *updated.ExecName)

		payload2 := map[string]interface{}{"student_number": student2, "pc_number": 1, "game": "CS:GO"}
		var second models.GamerActivity
		postActivity(t, payload2, &second)
		require.Equal(t, student2, second.StudentNumber)
		require.Equal(t, 1, second.PCNumber)
		require.Equal(t, "CS:GO", second.Game)
	})
}

func TestGetAllActivePCs(t *testing.T) {
	BeforeEach(t)
	t.Cleanup(func() { AfterEach(t) })

	t.Run("returns all active PCs", func(t *testing.T) {
		postActivity(t, payloadStudent1, nil)

		var activePCs []models.ActivePC
		getJSON(t, "/activity/get-active-pcs", &activePCs)
		require.Len(t, activePCs, 1)
		require.Equal(t, 1, activePCs[0].PCNumber)
		require.Equal(t, student1, activePCs[0].StudentNumber)
	})

	t.Run("returns empty list if no active PCs", func(t *testing.T) {
		postActivity(t, payloadStudent1, nil)
		updatePayload := map[string]interface{}{"pc_number": 1, "exec_name": "John"}
		updateActivity(t, student1, updatePayload, nil, http.StatusCreated)

		var activePCs []models.ActivePC
		getJSON(t, "/activity/get-active-pcs", &activePCs)
		require.Len(t, activePCs, 0)
	})
}
