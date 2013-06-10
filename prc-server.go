package main

import (
	"database/sql"
	"encoding/json"
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
	type bs struct {
		User        string
		Collections string
	}

	var u User
	var cs []Collection
	var err error
	if c.IsLoggedIn() {
		// Check errors, this is going to bite someday
		up, _ := GetUser(c.GetUsername())
		u = *up
		cs, err = CollectionListByUser(u.Username)
		log.Println(err)
	} else {
		u = User{}
	}
	uj, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
	}
	log.Println("cs:")
	log.Println(cs)
	cj, err := json.Marshal(cs)
	if err != nil {
		log.Println(err)
	}

	var bsc = bs{
		User:        string(uj),
		Collections: string(cj),
	}
	w.Header().Set("Content-Type", "text/javascript")
	tmpl.ExecuteTemplate(w, "bootstrap.js", bsc)
}

func setupHandlers() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v0.1").Subrouter()
	api.Handle("/parcels/{id:[0-9]+}", CtxHandler(ParcelDetailsHandler))
	api.Handle("/collections/{cid:[0-9]+}", CtxHandler(CollectionHandler))

	collPar := "/collections/{cid:[0-9]+}/parcels/{pid:[0-9]+}"
	api.Handle(collPar,
		CtxHandler(AddParcelToCollecditonHandler)).Methods("PUT")
	api.Handle(collPar,
		CtxHandler(RemoveParcelFromCollecditonHandler)).Methods("DELETE")

	api.Handle("/collections/", CtxHandler(UserCollectionHandler)).Methods("GET")
	api.Handle("/collections/", CtxHandler(NewCollectionHandler)).Methods("PUT")

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
