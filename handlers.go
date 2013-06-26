package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	decoder = schema.NewDecoder()
)

func ResponseWithError(id int, err error, w http.ResponseWriter, name string) {
	if err == sql.ErrNoRows {
		http.Error(w, name+" Not Found", 404)
		log.Printf("Could not find %s : %d", name, id)
	} else {
		log.Println(err)
		http.Error(w, "", 500)
	}
}

func ParcelDetailsHandler(w http.ResponseWriter, r *http.Request, c *Context) {
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

func CollectionHandler(w http.ResponseWriter, r *http.Request, c *Context) {
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

func UserCollectionHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	if c.IsLoggedIn() {
		if cs, err := CollectionListByUser(c.GetUsername()); err != nil {
			ResponseWithError(0, err, w, "Collections")
		} else {
			b, _ := json.Marshal(cs)
			w.Write(b)
		}
	} else {
		http.Error(w, "", 404)
	}
}

func NewCollectionHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	if c.IsLoggedIn() {
		col := new(Collection)
		j, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(j, &col)
		if err != nil {
			http.Error(w, "Could not add parcel, bad input", 500)
		}

		col.Owner = c.GetUsername()
		if err := AddCollection(col); err != nil {
			log.Println(err)
			http.Error(w, "Could not add parcel", 500)
		} else {
			b, _ := json.Marshal(col)
			w.Write(b)
		}
	} else {
		http.Error(w, "Must be logged in to add collection", 401)
	}
}

type ParcelCollectionAdjuster func(user string, cid, pid int) (*Collection, error)

func RemoveParcelFromCollecditonHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	ParcelCollecditonModifier(w, r, c, false)
}

func AddParcelToCollecditonHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	ParcelCollecditonModifier(w, r, c, true)
}

func ParcelCollecditonModifier(w http.ResponseWriter, r *http.Request, c *Context, add bool) {
	if c.IsLoggedIn() {
		vars := mux.Vars(r)
		cid, _ := strconv.Atoi(vars["cid"])
		pid, _ := strconv.Atoi(vars["pid"])
		var fn ParcelCollectionAdjuster
		if add {
			fn = AddParcelToCollectionById
		} else {
			fn = RemoveParcelFromCollection
		}
		if c, err := fn(c.GetUsername(), cid, pid); err != nil {
			log.Println(err)
			http.Error(w, "", 401)
		} else {
			b, _ := json.Marshal(c)
			w.Write(b)
		}
	} else {
		http.Error(w, "Must be looged in", 401)
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

func LoginHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	u, p := r.FormValue("username"), r.FormValue("password")
	user, err := Login(u, p)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 401)
	}
	finishLogin(r, w, c, user)
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	user := new(User)
	r.ParseForm()
	if err := decoder.Decode(user, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "Bad Form Values", 400)
	}
	if err := user.Register(); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s", err), 500)
	}
	finishLogin(r, w, c, user)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, c *Context) {
	c.Logout(r, w)
}

func finishLogin(r *http.Request, w http.ResponseWriter, c *Context, u *User) {
	b, _ := json.Marshal(u)
	c.Login(u, r, w)
	w.Write(b)
}
