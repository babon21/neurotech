package main

import (
	"context"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss"
	_ "github.com/volatiletech/authboss/auth"
)

type User struct {
	Username string
	Password string
}

func NewDbStorer(db *sqlx.DB) *DbStorer {
	return &DbStorer{Db: db}
}

func (d DbStorer) Save(_ context.Context, user authboss.User) error {
	u := user.(*User)

	userInsert := `INSERT INTO users (username, password) VALUES ($1, $2);`
	_, err := d.Db.Exec(userInsert, u.Username, u.Password)
	if err != nil {
		return err
	}

	debugln("Saved user:", u.Username)
	return nil
}

func (d DbStorer) Load(_ context.Context, key string) (authboss.User, error) {
	u := User{}
	err := d.Db.Get(&u, "SELECT * FROM users WHERE username=$1 LIMIT 1", key)

	if err != nil {
		return nil, authboss.ErrUserNotFound
	}

	return &u, nil
}

// GetPID from user
func (u User) GetPID() string { return u.Username }

// PutPID into user
func (u *User) PutPID(pid string) { u.Username = pid }

// PutPassword into user
func (u *User) PutPassword(password string) { u.Password = password }

// GetPassword from user
func (u User) GetPassword() string { return u.Password }
