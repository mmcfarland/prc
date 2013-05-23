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
