package models

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jameskeane/bcrypt"
	"gitlab.com/sirinibin/go-mysql-rest/config"
) //import "encoding/json"

type User struct {
	ID        uint64    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Username  string    `bson:"username" json:"username"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func FindUserByUsername(username string) (user User, err error) {
	err = config.DB.QueryRow("SELECT id,username,password from user where username=?", username).Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

func FindUserByID(id uint64) (*User, error) {

	var createdAt string
	var updatedAt string

	var user User
	err := config.DB.QueryRow("SELECT id,name,username,email,password,created_at,updated_at from user where id=?", id).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02 15:04:05"

	user.CreatedAt, err = time.Parse(layout, createdAt)
	if err != nil {
		return nil, err
	}

	user.UpdatedAt, err = time.Parse(layout, updatedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (user *User) IsEmailExists() (exists bool, err error) {

	var id uint64

	if user.ID != 0 {
		//Old Record
		err = config.DB.QueryRow("SELECT id from user where email=? and id!=?", user.Email, user.ID).Scan(&id)
	} else {
		//New Record
		err = config.DB.QueryRow("SELECT id from user where email=?", user.Email).Scan(&id)
	}
	return id != 0, err
}

func (user *User) IsUsernameExists() (exists bool, err error) {

	var id uint64

	if user.ID != 0 {
		//Old Record
		err = config.DB.QueryRow("SELECT id from user where username=? and id!=?", user.Username, user.ID).Scan(&id)
	} else {
		//New Record
		err = config.DB.QueryRow("SELECT id from user where username=?", user.Username).Scan(&id)
	}

	return id != 0, err
}

func (user *User) Insert() error {

	res, err := config.DB.Exec("INSERT INTO user(name, username, email, password,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?)", user.Name, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
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
	user.ID = uint64(id)
	log.Print("user.ID:")
	log.Print(user.ID)
	log.Printf("%d user created ", rows)

	return nil
}

func (user *User) Validate(w http.ResponseWriter, r *http.Request) (errs map[string]string) {

	errs = make(map[string]string)

	if govalidator.IsNull(user.Name) {
		errs["name"] = "Name is required"
	}

	if govalidator.IsNull(user.Username) {

		errs["username"] = "Username is required"
	}

	if govalidator.IsNull(user.Email) {
		errs["email"] = "E-mail is required"
	}

	if govalidator.IsNull(user.Password) {
		errs["password"] = "Password is required"
	}

	emailExists, err := user.IsEmailExists()
	if err != nil && err != sql.ErrNoRows {
		errs["email"] = err.Error()
	}

	if emailExists {
		errs["email"] = "E-mail is Already in use"
	}

	usernameExists, err := user.IsUsernameExists()
	if err != nil && err != sql.ErrNoRows {
		errs["username"] = err.Error()
	}

	if usernameExists {
		errs["username"] = "Username is Already in use"
	}

	if usernameExists || emailExists {
		w.WriteHeader(http.StatusConflict)
	} else if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	return errs
}

func HashPassword(password string) string {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(password, salt)
	return hash
}
