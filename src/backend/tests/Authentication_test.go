package tests

import (
	"encoding/json"
	"fmt"
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

var (
	APIURL = "/api/login"
	Token  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY1Mzg5OTl9.nUjkEfIYzdHX1xV-bR4GrneaMvWqyjdmqW64hBn-De8"
)

func InitAuthenticationTest(req *http.Request) *httptest.ResponseRecorder {
	os.Setenv("APP_ENVIRONMENT", "test")

	a := &c.Authentication{}
	r := chi.NewRouter()

	var mux = models.MuxServer{
		Mux:      r,
		Endpoint: "/api",
	}

	a.InitAuthentication(mux)

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	return res
}

func TestAuthenticationLogin(t *testing.T) {
	jsonBody := map[string]interface{}{
		"username": "rdev",
		"password": "iamsuperadmin",
	}

	b, _ := json.Marshal(jsonBody)

	req, _ := http.NewRequest("GET", APIURL, nil)
	req.Header.Set("Authentication", string(b))
	res := InitAuthenticationTest(req)

	fmt.Printf("\nRESPONSE: %s", res.Body)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestAuthenticationRegister(t *testing.T) {
	jsonBody := map[string]interface{}{
		"email":     "rdev@local",
		"username":  "rdev",
		"password":  "iamsuperadmin",
		"user_type": "superadmin",
	}

	b, _ := json.Marshal(jsonBody)

	req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(string(b)))
	res := InitAuthenticationTest(req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestPageNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/404", nil)
	res := InitAuthenticationTest(req)

	require.Equal(t, 404, res.Code)
}
