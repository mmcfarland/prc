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
	u = &User{}
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

func (u *User) Register() (err error) {
	exists := "SELECT username FROM users where username = $1;"
	var un string
	if s, e := DbConn.Prepare(exists); e != nil {
		log.Println(e)
		err = e
	} else {
		serr := s.QueryRow(u.Username).Scan(&un)
		if serr == sql.ErrNoRows {
			log.Println("here")
			u.SetPassword(u.UnhashedPassword)
			err = addUser(u)
			// TODO: queue email
		} else {
			err = fmt.Errorf("Username %q already exists", u.Username)
		}
	}
	return
}

func addUser(u *User) (err error) {
	NoCreateError := fmt.Errorf("Could not create user: %q", u.Username)
	insert := `INSERT INTO users (username, password, email)
                VALUES ($1, $2, $3);`
	if s, serr := DbConn.Prepare(insert); serr != nil {
		log.Println(serr)
		err = NoCreateError
	} else {
		if _, serr := s.Exec(u.Username, u.password, u.Email); serr != nil {
			log.Println(serr)
			err = NoCreateError
		}
	}
	return
}
