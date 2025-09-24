package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/handlers"
	"github.com/ubcesports/echo-base/internal/models"
)

var (
	testDB      *pgxpool.Pool
	testRouter  *http.ServeMux
	testHandler *handlers.Handler
)

func TestMain(m *testing.M) {
	testDB = InitTestDB()
	testConfig := config.LoadConfig()
	testHandler = &handlers.Handler{DB: testDB, Config: testConfig}

	testRouter = http.NewServeMux()
	testRouter.HandleFunc("/activity/{student_number}", handlers.Wrap(testHandler.GetGamerActivityByStudent))
	testRouter.HandleFunc("/activity/today/{student_number}", handlers.Wrap(testHandler.GetGamerActivityByTierOneStudentToday))
	testRouter.HandleFunc("/activity/all/recent", handlers.Wrap(testHandler.GetGamerActivity))
	testRouter.HandleFunc("/activity", handlers.Wrap(testHandler.AddGamerActivity))
	testRouter.HandleFunc("/activity/update/{student_number}", handlers.Wrap(testHandler.UpdateGamerActivity))

	exitCode := m.Run()
	CloseTestDB(testDB)
	os.Exit(exitCode)
}

func TestAddGamerActivity(t *testing.T) {
	BeforeEachTest(t, testDB)
	t.Cleanup(func() { AfterEachTest(t, testDB) })

	payload := map[string]interface{}{
		"student_number": "11223344",
		"pc_number":      1,
		"game":           "Valorant",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/activity", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {

	}
	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusCreated, res.StatusCode)

	var activity models.GamerActivityWithName
	err = json.NewDecoder(res.Body).Decode(&activity)
	require.NoError(t, err)

	require.Equal(t, "11223344", activity.StudentNumber)
	require.Equal(t, 1, activity.PCNumber)
	require.Equal(t, "Valorant", activity.Game)
	require.Nil(t, activity.ExecName)
}
