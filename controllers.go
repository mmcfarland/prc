package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func ResponseWithError(id int, err error, w http.ResponseWriter, name string) {
	if err == sql.ErrNoRows {
		http.Error(w, "Parcel Not Found", 404)
		log.Printf("Could not find parcel: %d", id)
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
