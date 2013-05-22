package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type User struct {
	Username         string     `json:"username" schema:"username"`
	password         []byte     `json:"-" schema:"-"`
	UnhashedPassword string     `json:"-" schema:"pass"`
	Email            string     `json:"email" schema:"email"`
	Joined           *time.Time `json:"joined"`
}

func (u *User) SetPassword(pw string) {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.password = h
}

func (u *User) CheckPassword(p string) (err error) {
	return bcrypt.CompareHashAndPassword(u.password, []byte(p))
}

func Login(un, p string) (u *User, err error) {
	loginError := errors.New("Login Failed: Unspecified Error")
	userSql := `SELECT username, password, email, joined 
            FROM users
            WHERE username = $1;`
	if s, e := DbConn.Prepare(userSql); e != nil {
		u = nil
		log.Println(e)
		err = loginError
		return
	} else {
		if e = s.QueryRow(un).Scan(&u.Username, &u.password, &u.Email, &u.Joined); e != nil {
			if e == sql.ErrNoRows {
				u = nil
				err = fmt.Errorf("Login Failed: User %q not found", un)
			} else {
				log.Println(e)
				err = loginError
			}
			return
		}
		if err = u.CheckPassword(p); err != nil {
			u = nil
			err = errors.New("Login Failed: Bad Password")
		}
		return
	}
}
