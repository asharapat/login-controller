package main

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var db *sql.DB
var cache redis.Conn


func main(){
	initCache()

	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("signup",SignUp)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)

	initDb()
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	cache = conn
}

func initDb(){
	var err error
	db, err = sql.Open("postgres", "dbname=mysystem sslmode=disable")
	if err != nil {
		panic(err)
	}
}
