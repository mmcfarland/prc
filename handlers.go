package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"strconv"
)

var decoder = schema.NewDecoder()

func ResponseWithError(id int, err error, w http.ResponseWriter, name string) {
	if err == sql.ErrNoRows {
		http.Error(w, name+" Not Found", 404)
		log.Printf("Could not find %s : %d", name, id)
	} else {
		log.Println(err)
		http.Error(w, "", 500)
	}
}

func ParcelDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	log.Printf("Get parcel: %d", id)
	if p, err := ParcelById(id); err != nil {
		ResponseWithError(id, err, w, "Parcel")
	} else {
		b, _ := json.Marshal(p)
		w.Write(b)
	}
}

func CollectionDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, _ := strconv.Atoi(vars["cid"])
	log.Printf("Get collection: %d", cid)
	if c, err := CollectionById(cid); err != nil {
		ResponseWithError(cid, err, w, "Collection")
	} else {
		b, _ := json.Marshal(c)
		w.Write(b)
	}
}
func ParcelLocationHandler(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.FormValue("lat"), 32)
	if err != nil {
		http.Error(w, "Bad latitude value", 500)
		return
	}
	lon, err := strconv.ParseFloat(r.FormValue("lon"), 32)
	if err != nil {
		http.Error(w, "Bad longitude value", 500)
		return
	}
	log.Printf("Search for: %f,%f", lat, lon)

	if p, err := ParcelByLocation(lat, lon); err != nil {
		ResponseWithError(0, err, w, "Parcel")
	} else {
		b, _ := json.Marshal([]Parcel{*p})
		w.Write(b)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	u, p := r.FormValue("username"), r.FormValue("password")
	user, err := Login(u, p)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 401)
	} else {
		b, _ := json.Marshal(user)
		w.Write(b)
	}
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	r.ParseForm()
	log.Println(r.PostForm)
	if err := decoder.Decode(user, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "Bad Form Values", 400)
	}
	log.Println(user)
	log.Println("I am")
	if err := user.Register(); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s", err), 500)
	} else {
		b, _ := json.Marshal(user)
		w.Write(b)
	}
	// Session start
}