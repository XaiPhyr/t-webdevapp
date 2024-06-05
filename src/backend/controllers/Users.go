package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"t_webdevapp/models"

	"github.com/go-chi/chi"
)

type Users struct {
	AppController
}

func (u Users) InitUsers(m models.MuxServer) {
	m.Mux.Route(m.Endpoint+"/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(u.mw.Authenticate)

			r.Get("/", u.GetAllUsers)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/detail", u.GetUser)
				r.Put("/update", u.UpdateUser)
				r.Delete("/delete", u.DeleteUser)
			})
		})
	})
}

func (u Users) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limitInt, _ := strconv.Atoi(query.Get("size"))
	offsetInt, _ := strconv.Atoi(query.Get("next"))

	nextPage := (offsetInt - 1) * limitInt

	result, count, err := u.userModel.ReadAllUsers(limitInt, nextPage)
	data := &models.UserResults{Data: result, Count: count}

	if err != nil {
		return
	}

	u.toJson(w, data)
}

func (u Users) GetUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	result, err := u.userModel.ReadOneUser(uuid)

	if err != nil {
		log.Printf("Error: %s", err)
		u.handleError(w, http.StatusNotFound, "User not found")
		return
	}

	u.toJson(w, result)
}

func (u Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := u.userModel.NewRegister()

	if err := json.NewDecoder(r.Body).Decode(&user); err == nil {
		if _, err := u.userModel.UpdateUser(w, user, u.handleError); err != nil {
			return
		}

		u.toJson(w, user)
	}
}

func (u Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	_, err := u.userModel.DeleteUser(w, uuid, u.handleError)

	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	u.toJson(w, nil)
}
