package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	c "t_webdevapp/controllers"
	"t_webdevapp/models"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func InitUsersTest(req *http.Request) *httptest.ResponseRecorder {
	os.Setenv("APP_ENVIRONMENT", "test")

	a := &c.Users{}
	r := chi.NewRouter()

	var mux = models.MuxServer{
		Mux:      r,
		Endpoint: "/api",
	}

	a.InitUsers(mux)

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	return res
}

func TestGetAllUsers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer 1")
	res := InitUsersTest(req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestUpdateUser(t *testing.T) {
	uuid := "09a0211f-1284-458b-82d7-2c2d7ece39f0"
	url := "/api/users/" + uuid + "/update"

	jsonBody := map[string]interface{}{
		"id":       52,
		"username": "rdev6update1",
		"email":    "upu@example.com",
	}

	b, _ := json.Marshal(jsonBody)

	req, _ := http.NewRequest("PUT", url, strings.NewReader(string(b)))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IiIsImV4cCI6MTcxNjc5MTM0NX0.czTEb3XKPKOlVXLftigu1wlpt0rYeM6Zb8VnPZKVSrY")
	res := InitUsersTest(req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestDeleteUser(t *testing.T) {
	uuid := "8378840f-a785-4685-8dcf-98ca5e8c7bff"
	url := "/api/users/" + uuid + "/delete"

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IiIsImV4cCI6MTcxNjg5MjE1NH0.Fbs5DKfBGLGsPb-5jCypA-d27Lqz8MOquDTjQghiWyE")
	res := InitUsersTest(req)

	require.Equal(t, http.StatusOK, res.Code)
}
