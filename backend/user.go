package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/volatiletech/authboss"
	_ "github.com/volatiletech/authboss/auth"
	aboauth "github.com/volatiletech/authboss/oauth2"
)

type User struct {
	UserID string `db:"username"`
	Password sql.NullString

	// OAuth2
	OAuth2UID          sql.NullString `db:"oauth2_uid"`
	OAuth2Provider     sql.NullString `db:"oauth2_provider"`
	OAuth2AccessToken  string
	OAuth2RefreshToken string
	OAuth2Expiry       time.Time
}

func NewDbStorer(db *sqlx.DB) *DbStorer {
	return &DbStorer{Db: db}
}

func (d DbStorer) Save(_ context.Context, user authboss.User) error {
	u := user.(*User)

	userInsert := `INSERT INTO users (username, password) VALUES ($1, $2);`
	_, err := d.Db.Exec(userInsert, u.UserID, u.Password)
	if err != nil {
		return err
	}

	debugln("Saved user:", u.UserID)
	return nil
}

func (d DbStorer) Load(_ context.Context, key string) (authboss.User, error) {
	u := User{}

	// Check to see if our key is actually an oauth2 pid
	provider, uid, err := authboss.ParseOAuth2PID(key)
	if err == nil {

		err := d.Db.QueryRowx("SELECT * FROM users WHERE oauth2_provider=$1 AND oauth2_uid=$2 LIMIT 1", provider, uid).StructScan(&u)
		if err != nil {
			log.Fatal(err)
			return nil, authboss.ErrUserNotFound
		}

		debugln("Loaded OAuth2 user:", u.OAuth2UID)
		return &u, nil
	}

	err = d.Db.Get(&u, "SELECT * FROM users WHERE username=$1 LIMIT 1", key)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// NewFromOAuth2 creates an oauth2 user (but not in the database, just a blank one to be saved later)
func (d DbStorer) NewFromOAuth2(_ context.Context, provider string, details map[string]string) (authboss.OAuth2User, error) {
	switch provider {
	case "google":

		user := &User{}
		user.UserID = details[aboauth.OAuth2Email]
		user.OAuth2UID = sql.NullString{String: details[aboauth.OAuth2UID], Valid: true}

		return user, nil
	}

	return nil, errors.Errorf("unknown provider %s", provider)
}

// SaveOAuth2 user
func (d DbStorer) SaveOAuth2(c context.Context, user authboss.OAuth2User) error {
	// u := user.(*User)

	// userInsert := `INSERT INTO users (username, oauth2_uid, oauth2_provider) VALUES ($1, $2, $3);`
	// _, err := d.Db.Exec(userInsert, u.UserID, u.OAuth2UID.String, u.OAuth2Provider.String)
	// if err != nil {
	// 	return err
	// }

	// debugln("Saved OAuth2 user:", u.UserID, u.OAuth2UID)
	return nil
}

// GetPID from user
func (u User) GetPID() string { return u.UserID }

// PutPID into user
func (u *User) PutPID(pid string) { u.UserID = pid }

// PutPassword into user
func (u *User) PutPassword(password string) { u.Password.String = password }

// GetPassword from user
func (u User) GetPassword() string { return u.Password.String }

// PutOAuth2UID into user
func (u *User) PutOAuth2UID(uid string) { u.OAuth2UID.String = uid }

// PutOAuth2Provider into user
func (u *User) PutOAuth2Provider(provider string) { u.OAuth2Provider.String = provider }

// PutOAuth2AccessToken into user
func (u *User) PutOAuth2AccessToken(token string) { u.OAuth2AccessToken = token }

// PutOAuth2RefreshToken into user
func (u *User) PutOAuth2RefreshToken(refreshToken string) { u.OAuth2RefreshToken = refreshToken }

// PutOAuth2Expiry into user
func (u *User) PutOAuth2Expiry(expiry time.Time) { u.OAuth2Expiry = expiry }

// IsOAuth2User returns true if the user was created with oauth2
func (u User) IsOAuth2User() bool { return len(u.OAuth2UID.String) != 0 }

// GetOAuth2UID from user
func (u User) GetOAuth2UID() (uid string) { return u.OAuth2UID.String }

// GetOAuth2Provider from user
func (u User) GetOAuth2Provider() (provider string) { return u.OAuth2Provider.String }

// GetOAuth2AccessToken from user
func (u User) GetOAuth2AccessToken() (token string) { return u.OAuth2AccessToken }

// GetOAuth2RefreshToken from user
func (u User) GetOAuth2RefreshToken() (refreshToken string) { return u.OAuth2RefreshToken }

// GetOAuth2Expiry from user
func (u User) GetOAuth2Expiry() (expiry time.Time) { return u.OAuth2Expiry }
