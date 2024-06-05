package routers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	c "t_webdevapp/controllers"
	"t_webdevapp/models"

	mw "t_webdevapp/middlewares"
	utils "t_webdevapp/utils"
	ws "t_webdevapp/websocket"

	"github.com/go-chi/chi"
)

var (
	auth = &c.Authentication{}
	user = &c.Users{}

	websocket = &ws.Websocket{}
)

func NewRoutes() chi.Router {
	endpoint := os.Getenv("ENDPOINT")

	r := chi.NewRouter()
	mw.UseMiddlewares(r)

	var mux = models.MuxServer{
		Mux:      r,
		Endpoint: endpoint,
	}

	// @routes
	auth.InitAuthentication(mux)
	user.InitUsers(mux)

	websocket.InitWebsocket(mux)

	// @status 404, 405
	PageNotFound(r)
	MethodNotAllowed(r)

	return r
}

func PageNotFound(r *chi.Mux) {
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		content, err := utils.ParseHTML("template/404.html", nil)

		if err != nil {
			log.Printf("Error: %s", err)
			return
		}

		w.Write([]byte(string(content)))
	})
}

func MethodNotAllowed(r *chi.Mux) {
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		b := map[string]interface{}{
			"err":  "method not allowed",
			"code": http.StatusMethodNotAllowed,
		}

		jsonMarshal, _ := json.MarshalIndent(b, "", "  ")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(jsonMarshal))
	})
}
