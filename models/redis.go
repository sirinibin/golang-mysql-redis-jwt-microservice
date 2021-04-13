package models

import (
	"errors"
	"strconv"
	"time"

	"gitlab.com/sirinibin/go-mysql-rest/config"
)

func (token *Token) SaveToRedis() error {
	expires := time.Unix(token.ExpiresAt, 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	errAccess := config.RedisClient.Set(token.AccessUUID, strconv.Itoa(int(token.UserID)), expires.Sub(now)).Err()

	return errAccess
}

func (token *Token) ExistsInRedis() error {

	userid, err := config.RedisClient.Get(token.AccessUUID).Result()
	if err != nil {
		return err
	}

	userID, err := strconv.ParseUint(userid, 10, 64)
	if err != nil {
		return err
	}

	if token.UserID != userID {
		return errors.New("User id doesn't exist in redis!")
	}

	return nil

}
