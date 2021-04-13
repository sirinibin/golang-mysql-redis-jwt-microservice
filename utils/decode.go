package utils

import (
	"encoding/json"
	"net/http"

	"gitlab.com/sirinibin/go-mysql-rest/models"
)

func Decode(w http.ResponseWriter, r *http.Request, v interface{}) bool {

	var response models.Response

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil || v == nil {

		response.Status = false
		response.Errors = make(map[string]string)
		response.Errors["input"] = "Invalid Data:" + err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return false
	}
	return true
}
