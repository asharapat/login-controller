package main

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/errors/fmt"
	"net/http"
	"time"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}


func Signup(creds *Credentials) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password),8)
	if _, err = db.Query("insert into users values ($1,$2)", creds.Username, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}

func Signin(creds *Credentials) (*http.Cookie, error) {

	result, err := db.Query("select password from users where username =$1", creds.Username)
	if err != nil {
		return nil, err
	}

	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		return nil, errors.New("Password is not valid")
	}

	sessionToken := uuid.New().String()
	_, err = cache.Do("SETEX", sessionToken, "120", creds.Username)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:       "session_token",
		Value:      sessionToken,
		Expires:    time.Now().Add(120*time.Second),
	}
	
	return cookie, nil

}

func WelcomeS(sessionToken string) (interface{}, error) {
	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("No sessions found for user")
	}

	return response, nil
}


func RefreshS(sessionToken string) (*http.Cookie, error) {
	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("No sessions found for user")
	}

	newSessionToken := uuid.New().String()
	_, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s", response))
	if err != nil {
		return nil, err
	}
	_, err = cache.Do("DEL", sessionToken)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:       "session_token",
		Value:      newSessionToken,
		Expires:    time.Now().Add(120*time.Second),
	}

	return cookie, nil

}
