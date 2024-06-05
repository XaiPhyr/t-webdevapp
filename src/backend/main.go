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
	log.Printf("-> Local:   http://localhost:8200")
	log.Printf("-> Version: %s", cfg.Env)
	fmt.Println()

	log.Printf("%s", http.ListenAndServe(":8200", r))
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
