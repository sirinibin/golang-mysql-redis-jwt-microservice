package models

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"gitlab.com/sirinibin/go-mysql-rest/config"
)

// Employee : struct for Employee model
type Employee struct {
	ID        uint64    `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedBy uint64    `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy uint64    `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

type SearchCriterias struct {
	Page     uint32                 `bson:"page,omitempty" json:"page,omitempty"`
	Size     uint32                 `bson:"size,omitempty" json:"size,omitempty"`
	SearchBy map[string]interface{} `bson:"search_by,omitempty" json:"search_by,omitempty"`
	SortBy   string                 `bson:"sort_by,omitempty" json:"sort_by,omitempty"`
}

func FindEmployees(criterias SearchCriterias) (*[]Employee, error) {

	var employees []Employee

	offset := (criterias.Page - 1) * criterias.Size

	searchString := ""

	args := []interface{}{}

	if len(criterias.SearchBy) > 0 {
		searchString += " WHERE "
		i := 1
		for field, v := range criterias.SearchBy {

			value, ok := v.(string)
			if ok {
				if _, err := strconv.Atoi(value); err == nil {
					//Integer
					searchString += field + " = ? "
					args = append(args, value)
				} else {
					//string
					searchString += field + " like ? "
					args = append(args, value+"%")
				}

				if i < len(criterias.SearchBy) {
					searchString += " AND "
				}

			}

			i++
		}
	}

	sortString := ""
	if !govalidator.IsNull(criterias.SortBy) {
		sortString = " ORDER BY  ? "
		args = append(args, criterias.SortBy)
	}
	args = append(args, offset)
	args = append(args, criterias.Size)

	query := "SELECT id,name,email,created_by,updated_by,created_At,updated_at FROM employee " + searchString + sortString + " limit ?,?"
	res, err := config.DB.Query(query, args...)
	defer res.Close()

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var employee Employee

		var createdAt string
		var updatedAt string
		err := res.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.CreatedBy, &employee.UpdatedBy, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		layout := "2006-01-02 15:04:05"

		employee.CreatedAt, err = time.Parse(layout, createdAt)
		if err != nil {
			return nil, err
		}

		employee.UpdatedAt, err = time.Parse(layout, updatedAt)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)

	}
	return &employees, nil

}

func DeleteEmployee(employeeID uint64) (int64, error) {

	res, err := config.DB.Exec("DELETE from employee where id=?", employeeID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()

}

func IsEmployeeExists(employeeID uint64) (exists bool, err error) {

	var id uint64

	err = config.DB.QueryRow("SELECT id from employee where id=?", employeeID).Scan(&id)

	return id != 0, err
}

func FindEmployeeByID(id uint64) (*Employee, error) {

	var createdAt string
	var updatedAt string
	var employee Employee

	err := config.DB.QueryRow("SELECT id,created_by,updated_by,name,email,created_at,updated_at from employee where id=?", id).Scan(&employee.ID, &employee.CreatedBy, &employee.UpdatedBy, &employee.Name, &employee.Email, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02 15:04:05"

	employee.CreatedAt, err = time.Parse(layout, createdAt)
	if err != nil {
		return &employee, err
	}

	employee.UpdatedAt, err = time.Parse(layout, updatedAt)
	if err != nil {
		return &employee, err
	}

	return &employee, err
}

func (employee *Employee) IsEmailExists() (exists bool, err error) {

	var id uint64

	if employee.ID != 0 {
		//Old Record
		err = config.DB.QueryRow("SELECT id from employee where email=? and id!=?", employee.Email, employee.ID).Scan(&id)
	} else {
		//New Record
		err = config.DB.QueryRow("SELECT id from employee where email=?", employee.Email).Scan(&id)
	}
	return id != 0, err
}

func (employee *Employee) Validate(w http.ResponseWriter, r *http.Request, scenario string) (errs map[string]string) {

	errs = make(map[string]string)

	if scenario == "update" {
		if employee.ID == 0 {
			errs["id"] = "ID is required"
			return errs
		}
		exists, err := IsEmployeeExists(employee.ID)
		if err != nil || !exists {
			errs["id"] = err.Error()
			return errs
		}

	}

	if govalidator.IsNull(employee.Name) {
		errs["name"] = "Name is required"
	}

	if govalidator.IsNull(employee.Email) {

		errs["username"] = "E-mail is required"
	}

	emailExists, err := employee.IsEmailExists()
	if err != nil && err != sql.ErrNoRows {
		errs["email"] = err.Error()
	}

	if emailExists {
		errs["email"] = "E-mail is Already in use"
	}

	if emailExists {
		w.WriteHeader(http.StatusConflict)
	} else if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	return errs
}

func (employee *Employee) Insert() error {

	res, err := config.DB.Exec("insert into employee (name,created_by,updated_by, email,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?)", employee.Name, employee.CreatedBy, employee.UpdatedBy, employee.Email, employee.CreatedAt, employee.UpdatedAt)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error %s when finding last insert Id", err)
		return err
	}
	employee.ID = uint64(id)
	log.Print("user.ID:")
	log.Print(employee.ID)
	log.Printf("%d employee created ", rows)

	return nil
}

func (employee *Employee) Update() (*Employee, error) {

	res, err := config.DB.Exec("UPDATE employee SET name=?, updated_by=? ,email=?, updated_at=? WHERE id=?", employee.Name, employee.UpdatedBy, employee.Email, employee.UpdatedAt, employee.ID)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return nil, err
	}

	employee, err = FindEmployeeByID(employee.ID)
	if err != nil {
		return nil, err
	}

	log.Print("user.ID:")
	log.Print(employee.ID)
	log.Printf("%d employee updated ", rows)

	return employee, nil
}
