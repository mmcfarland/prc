package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/bmizerany/pq"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"text/template"
)

var (
	port      = flag.Int("port", 7878, "Port")
	DbConn    = setupDb()
	indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, r.Host)
}

func dbHandler(fn func(w http.ResponseWriter, r *http.Request, db *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func setupHandlers() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v0.1").Subrouter()
	api.HandleFunc("/parcels/{id:[0-9]+}", ParcelDetailsHandler)
	api.HandleFunc("/collections/{cid:[0-9]+}", CollectionDetailsHandler)
	api.HandleFunc("/parcels/", ParcelLocationHandler).Queries("lat", "", "lon", "")

	r.HandleFunc("/", indexHandler)
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("client"))))
	http.Handle("/", r)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	defer DbConn.Close()

	setupHandlers()
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), nil); err != nil {
		fmt.Println("Failed to start server: %v", err)
	} else {
		fmt.Println("Listening on port: " + strconv.Itoa(*port))
	}
}
