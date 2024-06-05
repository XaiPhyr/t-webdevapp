package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"t_webdevapp/routers"
	"t_webdevapp/utils"

	"github.com/go-chi/chi"
)

var version = os.Getenv("VERSION")
var port = os.Getenv("HTTP_PORT")

func main() {
	r := routers.NewRoutes()
	cfg := utils.InitConfig()

	log.SetFlags(log.Llongfile | log.LstdFlags)

	src := cfg.Frontend.Source
	_, err := os.Stat(src + "/index.html")

	if !os.IsNotExist(err) {
		FileServer(r, src)
	}

	fmt.Println()
	log.Printf("-> Local:   http://localhost:%s", port)
	log.Printf("-> Version: %s", version)
	fmt.Println()

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func FileServer(r chi.Router, src string) {
	fs := http.FileServer(http.Dir(src))
	r.Handle("/", http.StripPrefix("/", fs))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(src + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
			return
		}

		fs.ServeHTTP(w, r)
	})
}
