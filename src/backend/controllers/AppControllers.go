package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"t_webdevapp/middlewares"
	"t_webdevapp/models"
)

type AppController struct {
	endpoint  string
	mw        middlewares.Middleware
	userModel *models.Users
}

var endpoint = os.Getenv("ENDPOINT")

func (a AppController) toJson(w http.ResponseWriter, b interface{}) {
	jsonMarshal, _ := json.MarshalIndent(b, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonMarshal))
}

func (a AppController) handleError(w http.ResponseWriter, code int, message string) {
	errObj := models.ErrorObject{
		Code:    code,
		Message: message,
	}

	jsonMarshal, _ := json.MarshalIndent(errObj, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errObj.Code)
	w.Write([]byte(jsonMarshal))
}
