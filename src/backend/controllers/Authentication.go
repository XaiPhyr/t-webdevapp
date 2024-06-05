package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"t_webdevapp/models"
	"time"

	s "t_webdevapp/services"
	utils "t_webdevapp/utils"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type Authentication struct {
	AppController
}

func (a Authentication) InitAuthentication(m models.MuxServer) {
	m.Mux.Route(m.Endpoint+"/login", func(r chi.Router) {
		r.Get("/", a.Login)
	})

	m.Mux.Route(m.Endpoint+"/register", func(r chi.Router) {
		r.Post("/", a.Register)
	})

	m.Mux.Route(m.Endpoint+"/notify", func(r chi.Router) {
		r.Get("/", a.Notify)
	})
}

func (a Authentication) Login(w http.ResponseWriter, r *http.Request) {
	var info models.Authentication
	var l models.Login
	var user *models.Users

	token, err := s.GenerateJWT()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := r.Header.Get("Authentication")
	readAuth := strings.NewReader(auth)
	err = json.NewDecoder(readAuth).Decode(&l)

	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	user, err = a.userModel.ReadByUsernameOrEmail(l.Username)

	if err != nil {
		a.handleError(w, http.StatusNotFound, "User not found")
		return
	}

	if err = s.CheckPassword(user.Password, l.Password); err != nil {
		a.handleError(w, http.StatusUnauthorized, "Password Incorrect")
		return
	}

	user.LastLogin = time.Now()
	if _, err := a.userModel.UpdateUser(w, user, a.handleError); err != nil {
		return
	}

	info.Token = token
	info.Username = l.Username
	info.RememberMe = l.RememberMe

	utils.ListenNotify()
	a.toJson(w, info)
}

func (a Authentication) Register(w http.ResponseWriter, r *http.Request) {
	user := a.userModel.NewRegister()

	if err := json.NewDecoder(r.Body).Decode(&user); err == nil {
		hashPassword, _ := s.HashPassword(user.Password)
		user.Password = hashPassword

		if _, err := a.userModel.CreateUser(w, user, a.handleError); err != nil {
			return
		}

		a.toJson(w, user)
	}
}

func (a Authentication) Notify(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	a.reader(ws)
}

func (a Authentication) reader(conn *websocket.Conn) {
	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}

		log.Printf("Websocket: %s", string(p))

		if err := conn.WriteMessage(mt, p); err != nil {
			log.Printf("Error: %s", err)
			return
		}
	}
}
