package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jameskeane/bcrypt"
)

//Authorize : Authorize structure
type AuthorizeRequest struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type AuthCodeResponse struct {
	Code      string `bson:"code" json:"code"`
	ExpiresAt int64  `bson:"expires_at" json:"expires_at"`
}

func AuthenticateByAuthCode(tokenStr string) (tokenClaims TokenClaims, err error) {

	if govalidator.IsNull(tokenStr) {
		return tokenClaims, errors.New("auth_code is required")
	}

	tokenClaims, err = AuthenticateByJWTToken(tokenStr)
	if err != nil {
		return tokenClaims, err
	}
	if tokenClaims.Type != "auth_code" {
		return tokenClaims, errors.New("Invalid auth code")
	}

	return tokenClaims, nil

}

//GenerateAuthCode : generate and return authcode
func (auth *AuthorizeRequest) GenerateAuthCode() (authCode AuthCodeResponse, err error) {

	// Generate Auth code
	expiresAt := time.Now().Add(time.Hour * 5) // expiry for auth code is 5min
	token, err := generateAndSaveToken(auth.Username, expiresAt, "auth_code")
	authCode.ExpiresAt = token.ExpiresAt
	authCode.Code = token.TokenStr

	return authCode, err
}

//Validate : Validate authorization data
func (auth *AuthorizeRequest) Authenticate() (errs map[string]string) {

	errs = make(map[string]string)

	if govalidator.IsNull(auth.Username) {
		errs["username"] = "Username is required"
	}
	if govalidator.IsNull(auth.Password) {
		errs["password"] = "Password is required"
	}

	if !govalidator.IsNull(auth.Password) && !govalidator.IsNull(auth.Username) {

		user, err := FindUserByUsername(auth.Username)
		if err != nil && err != sql.ErrNoRows {
			errs["password"] = "Error finding user record:" + err.Error()
		}

		if err == sql.ErrNoRows || !bcrypt.Match(auth.Password, user.Password) {
			errs["password"] = "Username or Password is wrong"
		}

	}

	return errs
}
