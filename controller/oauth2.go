package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/sirinibin/go-mysql-rest/config"
	"gitlab.com/sirinibin/go-mysql-rest/models"
	"gitlab.com/sirinibin/go-mysql-rest/utils"
)

// Authorize : Authorize a user account
func Authorize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	response.Errors = make(map[string]string)

	var auth *models.AuthorizeRequest

	if !utils.Decode(w, r, &auth) {
		return
	}

	// Authenticate
	if errs := auth.Authenticate(); len(errs) > 0 {
		response.Status = false
		response.Errors = errs
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	authCode, err := auth.GenerateAuthCode()
	if err != nil {
		response.Status = false
		response.Errors["auth_code"] = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Result = authCode

	json.NewEncoder(w).Encode(response)

}

// Accesstoken : handler for POST /accesstoken
func Accesstoken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	response.Errors = make(map[string]string)

	var accessTokenRequest *models.AccesstokenRequest

	if !utils.Decode(w, r, &accessTokenRequest) {
		return
	}

	tokenClaims, err := models.AuthenticateByAuthCode(accessTokenRequest.Code)
	if err != nil {
		response.Status = false

		response.Errors["auth_code"] = err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	accessToken, err := models.GenerateAccesstoken(tokenClaims.Username)
	if err != nil {
		response.Status = false

		response.Errors["access_token"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Result = accessToken

	json.NewEncoder(w).Encode(response)

}

// RefreshAccesstoken : handler for POST /refresh
func RefreshAccesstoken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	response.Errors = make(map[string]string)

	var refreshAccessTokenRequest *models.RefreshAccesstokenRequest

	if !utils.Decode(w, r, &refreshAccessTokenRequest) {
		return
	}

	tokenClaims, err := models.AuthenticateByRefreshToken(refreshAccessTokenRequest.RefreshToken)
	if err != nil {
		response.Status = false

		response.Errors["refresh_token"] = err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	accessToken, err := models.GenerateAccesstoken(tokenClaims.Username)
	if err != nil {
		response.Status = false

		response.Errors["access_token"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Result = accessToken

	json.NewEncoder(w).Encode(response)

}

// LogOut : handler for DELETE /logout
func LogOut(w http.ResponseWriter, r *http.Request) {
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

	deleted, err := config.RedisClient.Del(tokenClaims.AccessUUID).Result()
	if err != nil || deleted == 0 {
		response.Status = false
		response.Errors["access_token"] = err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return

	}
	response.Status = true
	response.Result = "Successfully logged out"

	json.NewEncoder(w).Encode(response)

}

// APIInfo : handler function for / call
func APIInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response models.Response
	response.Status = true
	response.Result = "GoLang / MySql Microservice [ OAuth2, JWT and Redis used for security ] "

	json.NewEncoder(w).Encode(response)
}
