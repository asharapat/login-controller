package main

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/exp/errors/fmt"
	"net/http"
)

// sign up new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = Signup(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// sign in and set session token for user
func SignIn(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie, err := Signin(creds)
	if err != nil {
		if err == sql.ErrNoRows || err.Error() == "Password is not valid"{
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)
}

// check user's session token
func Welcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	res, err := WelcomeS(sessionToken)
	if err != nil  {
		if err.Error() == "No sessions found for user" {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", res)))
}

// refresh user's session token
func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie{
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	cookie, err := RefreshS(sessionToken)
	if err !=nil {
		if err.Error() == "No sessions found for user" {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	http.SetCookie(w, cookie)
}
