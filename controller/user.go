package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/sirinibin/go-mysql-rest/models"
	"gitlab.com/sirinibin/go-mysql-rest/utils"
)

// Me : handler function for /v1/me call
func Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	response.Errors = make(map[string]string)

	tokenClaims, err := models.AuthenticateByAccessToken(r)
	if err != nil {
		response.Status = false
		response.Errors["access_token"] = "Invalid Access token:" + err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := models.FindUserByID(tokenClaims.UserID)
	if err != nil {
		response.Status = false
		response.Errors["find_user"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	user.Password = ""

	response.Status = true
	response.Result = user

	json.NewEncoder(w).Encode(response)
}

// Register : Register a new user account
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response

	var user *models.User

	// Decode data
	if !utils.Decode(w, r, &user) {
		return
	}

	// Validate data
	if errs := user.Validate(w, r); len(errs) > 0 {
		response.Status = false
		response.Errors = errs
		json.NewEncoder(w).Encode(response)
		return
	}

	// Insert new record
	user.Password = models.HashPassword(user.Password)
	user.CreatedAt = time.Now().Local()
	user.UpdatedAt = time.Now().Local()

	err := user.Insert()
	if err != nil {
		response.Status = false
		response.Errors = make(map[string]string)
		response.Errors["insert"] = "Unable to Insert to db:" + err.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	user.Password = ""
	response.Result = user

	json.NewEncoder(w).Encode(response)

}
