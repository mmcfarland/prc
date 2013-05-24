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
	port   = flag.Int("port", 7878, "Port")
	DbConn = setupDb()
	tmpl   = template.Must(template.ParseFiles("templates/index.html", "templates/bootstrap.js"))
)

const (
	ConnString = "user=postgres dbname=dataviewer"
)

type CtxHandler func(http.ResponseWriter, *http.Request, *Context)

func (h CtxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, err := NewContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h(w, r, ctx)
}

func setupDb() (db *sql.DB) {
	db, err := sql.Open("postgres", ConnString)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", r.Host)
}

func bootstrapHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	log.Println("bst")
	var u *User
	if c.IsLoggedIn() {
		u, _ = GetUser(c.GetUsername())
	}
	w.Header().Set("Content-Type", "text/javascript")
	tmpl.ExecuteTemplate(w, "bootstrap.js", u)
}

func setupHandlers() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v0.1").Subrouter()
	api.Handle("/parcels/{id:[0-9]+}", CtxHandler(ParcelDetailsHandler))
	api.Handle("/collections/{cid:[0-9]+}", CtxHandler(CollectionHandler))
	api.HandleFunc("/parcels/", ParcelLocationHandler).Queries("lat", "", "lon", "")
	api.Handle("/login/", CtxHandler(LoginHandler)).Methods("POST")
	api.Handle("/register/", CtxHandler(RegistrationHandler)).Methods("POST")
	api.Handle("/logout/", CtxHandler(LogoutHandler)).Methods("POST")

	r.Handle("/bs.js", CtxHandler(bootstrapHandler))
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
