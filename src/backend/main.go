package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sql "app_backend/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"

	_ "github.com/lib/pq"
)

var endpoint = os.Getenv("ENDPOINT")
var version = os.Getenv("VERSION")
var url = func(api string) string {
	return endpoint + api
}

func main() {
	port := os.Getenv("HTTP_PORT")
	log.SetFlags(log.Llongfile | log.LstdFlags)
	sql.MigrateSQL()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Get(fmt.Sprintf("%s/", endpoint), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("WELCOME!"))
	})

	r.Route(url("/hello"), func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("HELLO WORLD"))
		})
	})

	fmt.Println()
	log.Printf("-> Version: %s", version)
	log.Printf("-> Local:   http://localhost:8200")

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
