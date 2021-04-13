package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/sirinibin/go-mysql-rest/models"
	"gitlab.com/sirinibin/go-mysql-rest/utils"
)

// CreateEmployee : handler for POST /employees
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
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

	var employee *models.Employee
	// Decode data
	if !utils.Decode(w, r, &employee) {
		return
	}

	employee.CreatedBy = tokenClaims.UserID
	employee.UpdatedBy = tokenClaims.UserID
	employee.CreatedAt = time.Now().Local()
	employee.UpdatedAt = time.Now().Local()

	// Validate data
	if errs := employee.Validate(w, r, "create"); len(errs) > 0 {
		response.Status = false
		response.Errors = errs
		json.NewEncoder(w).Encode(response)
		return
	}

	err = employee.Insert()
	if err != nil {
		response.Status = false
		response.Errors = make(map[string]string)
		response.Errors["insert"] = "Unable to insert to db:" + err.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Result = employee

	json.NewEncoder(w).Encode(response)

}

// UpdateEmployee : handler for PUT /employees
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
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

	var employee *models.Employee
	// Decode data
	if !utils.Decode(w, r, &employee) {
		return
	}

	employee.UpdatedBy = tokenClaims.UserID
	employee.UpdatedAt = time.Now().Local()

	// Validate data
	if errs := employee.Validate(w, r, "update"); len(errs) > 0 {
		response.Status = false
		response.Errors = errs
		json.NewEncoder(w).Encode(response)
		return
	}

	employee, err = employee.Update()
	if err != nil {
		response.Status = false
		response.Errors = make(map[string]string)
		response.Errors["update"] = "Unable to update:" + err.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Result = employee

	json.NewEncoder(w).Encode(response)

}

// DeleteEmployee : handler for DELETE /employees/{id}
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	response.Errors = make(map[string]string)

	_, err := models.AuthenticateByAccessToken(r)
	if err != nil {
		response.Status = false
		response.Errors["access_token"] = "Invalid Access token:" + err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)

	employeeID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		response.Status = false
		response.Errors["id"] = "Invalid record id:" + err.Error()
		json.NewEncoder(w).Encode(response)
		return

	}

	res, err := models.DeleteEmployee(employeeID)
	if err != nil || res == 0 {
		response.Status = false
		if err != nil {
			response.Errors["delete"] = "Unable to delete:" + err.Error()
		} else {
			response.Errors["delete"] = "Unable to delete"
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if res == 1 {
		response.Status = true
		response.Result = "Deleted successfully"
	}

	json.NewEncoder(w).Encode(response)

}

// ViewEmployee : handler for GET /employees/{id}
func ViewEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	response.Errors = make(map[string]string)

	_, err := models.AuthenticateByAccessToken(r)
	if err != nil {
		response.Status = false
		response.Errors["access_token"] = "Invalid Access token:" + err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)

	employeeID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		response.Status = false
		response.Errors["id"] = "Invalid record id:" + err.Error()
		json.NewEncoder(w).Encode(response)
		return

	}

	var employee *models.Employee

	employee, err = models.FindEmployeeByID(employeeID)
	if err != nil {
		response.Status = false
		response.Errors["view"] = "Unable to view:" + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Result = employee

	json.NewEncoder(w).Encode(response)

}

// ListEmployee : handler for GET /employees
func ListEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	response.Errors = make(map[string]string)

	_, err := models.AuthenticateByAccessToken(r)
	if err != nil {
		response.Status = false
		response.Errors["access_token"] = "Invalid Access token:" + err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var employees *[]models.Employee

	criterias := models.SearchCriterias{
		Page:   1,
		Size:   2,
		SortBy: "updated_at desc",
	}
	keys, ok := r.URL.Query()["page"]
	if ok {
		page, err := strconv.ParseUint(keys[0], 10, 64)
		if err != nil {
			response.Status = false
			response.Errors["page"] = "Invalid page value:" + err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
		criterias.Page = uint32(page)
	}

	keys, ok = r.URL.Query()["size"]
	if ok && len(keys[0]) >= 1 {
		size, err := strconv.ParseUint(keys[0], 10, 64)
		if err != nil {
			response.Status = false
			response.Errors["size"] = "Invalid size value:" + err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
		criterias.Size = uint32(size)
	}

	keys, ok = r.URL.Query()["sort"]
	if ok && len(keys[0]) >= 1 {
		criterias.SortBy = keys[0]
	}

	criterias.SearchBy = make(map[string]interface{})
	keys, ok = r.URL.Query()["search[name]"]
	if ok && len(keys[0]) >= 1 {
		criterias.SearchBy["name"] = keys[0]
	}
	keys, ok = r.URL.Query()["search[id]"]
	if ok && len(keys[0]) >= 1 {
		criterias.SearchBy["id"] = keys[0]
	}
	keys, ok = r.URL.Query()["search[email]"]
	if ok && len(keys[0]) >= 1 {
		criterias.SearchBy["email"] = keys[0]
	}

	employees, err = models.FindEmployees(criterias)
	if err != nil {
		response.Status = false
		response.Errors["find"] = "Unable to find employees:" + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Criterias = criterias
	response.Result = employees

	json.NewEncoder(w).Encode(response)

}
