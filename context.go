package main

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("todo-change-to-secret"))

type Context struct {
	Db      *sql.DB
	Session *sessions.Session
}

func NewContext(r *http.Request) (*Context, error) {
	s, err := store.Get(r, "prc")
	return &Context{
		Db:      DbConn,
		Session: s,
	}, err
}

func (c *Context) Login(u *User, r *http.Request, w http.ResponseWriter) {
	c.Session.Values["loggedin"] = true
	c.Session.Values["user"] = u.Username
	c.Session.Save(r, w)
}

func (c *Context) Logout(r *http.Request, w http.ResponseWriter) {
	c.Session.Values["loggedin"] = false
	c.Session.Values["user"] = nil
	c.Session.Save(r, w)
}

func (c *Context) IsLoggedIn() bool {
	return c.Session.Values["loggedin"] == true
}

func (c *Context) GetUsername() string {
	if str, ok := c.Session.Values["use"].(string); ok {
		return str
	} else {
		return ""
	}
}
