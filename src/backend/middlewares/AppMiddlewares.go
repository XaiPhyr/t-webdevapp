package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"t_webdevapp/models"
	s "t_webdevapp/services"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
)

type Middleware struct{}

var (
	LimiterExpiration = 1 * time.Minute
	RequestID         = middleware.RequestID
	RealIP            = middleware.RealIP
	Logger            = middleware.Logger
	Recoverer         = middleware.Recoverer
	GetHead           = middleware.GetHead
	HttpRate          = httprate.LimitByIP(100, LimiterExpiration)
)

func (m Middleware) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			errObj := &models.ErrorObject{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			}

			jsonMarshal, _ := json.MarshalIndent(errObj, "", "  ")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(jsonMarshal))

			return
		}

		t := strings.Split(auth, " ")[1]

		if err := s.VerifyJWT(t, w); err != nil {
			errObj := &models.ErrorObject{
				Code:    http.StatusUnauthorized,
				Message: "Token Expired",
			}

			jsonMarshal, _ := json.MarshalIndent(errObj, "", "  ")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(jsonMarshal))

			return
		}

		h.ServeHTTP(w, r)
	})
}

func UseMiddlewares(r *chi.Mux) {
	r.Use(RequestID)
	r.Use(RealIP)
	r.Use(Logger)
	r.Use(Recoverer)
	r.Use(GetHead)
	r.Use(HttpRate)
}
