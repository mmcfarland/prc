package main

import (
	"database/sql"
	"flag"
	_ "github.com/bmizerany/pq"
	"log"
	"net/http"
	"runtime"
	"strconv"
	//"text/template"
)

var (
	port = flag.Int("port", 7878, "Port")
)

const (
	ConnString = "user=postgres dbname=dataviewer"
)

func setupDb() (db *sql.DB) {
	db, err := sql.Open("postgres", ConnString)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	return
}

func dbHandler(fn func(w http.ResponseWriter, r *http.Request, db *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func setupHandlers() {
	//http.Handle("/",  indexHandler)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	db := setupDb()
	defer db.Close()

	if err := http.ListenAndServe(":"+strconv.Itoa(*port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		log.Println("Listening on port: " + strconv.Itoa(*port))
	}
}
