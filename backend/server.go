package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	_ "github.com/volatiletech/authboss/auth"

	"github.com/volatiletech/authboss/defaults"

	"gopkg.in/mgo.v2"
)

type DbStorer struct {
	Db *sqlx.DB
}

var (
	ab           = authboss.New()
	database     DbStorer
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

const (
	sessionCookieName = "ab_blog"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	databaseSite := session.DB("neurotech")

	connStr := "host=localhost port=5432 dbname=neurotech user=neurotech password=neurotech"
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	database = *NewDbStorer(db)

	// password := "super"
	// pass, err := bcrypt.GenerateFromPassword([]byte(password), ab.Config.Modules.BCryptCost)
	// if err != nil {
	// 	panic(err)
	// }

	// err = database.Save(nil, &User{
	// 	Username: "kelan007",
	// 	Password: string(pass)})

	// if err != nil {
	// 	log.Fatal(err)
	// }

	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	// Initialize authboss
	setupAuthboss()

	mux := chi.NewRouter()
	mux.Use(ab.LoadClientStateMiddleware)

	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.ModuleListMiddleware(ab))
		mux.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))

		studenWorkCollection := InitStudentWorksCollection(databaseSite)
		studenWorkHandler := &StudentWorkHandler{Collection: studenWorkCollection}

		publicationCollection := InitPublicationsCollection(databaseSite)
		publicationHandler := &PublicationHandler{Collection: publicationCollection}

		newsCollection := InitNewsCollection(databaseSite)
		newsHandler := &NewsHandler{Collection: newsCollection}

		disciplineHandler := &DisciplineHandler{path: DisciplinePath}
		studyHandler := &StudyMaterialHandler{path: DisciplinePath}

		mux.Handle("/publications", publicationHandler)
		mux.Handle("/student-work", studenWorkHandler)
		mux.Handle("/news", newsHandler)
		mux.Handle("/disciplines", disciplineHandler)
		mux.Handle("/study-materials", studyHandler)
	})

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8002"
	}
	log.Printf("Listening on localhost: %s", port)
	log.Println(http.ListenAndServe("localhost:"+port, mux))
}

func setupAuthboss() {
	ab.Config.Paths.RootURL = "http://localhost:8002"

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")
	// ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}

	defaults.SetCore(&ab.Config, false, false)

	if err := ab.Init(); err != nil {
		panic(err)
	}
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		h.ServeHTTP(w, r)
	})
}

type User struct {
	Username string
	Password string
}

func NewDbStorer(db *sqlx.DB) *DbStorer {
	return &DbStorer{Db: db}
}

type StorageOpts struct {
	ConnString         string
	MaxConnections     int
	MaxIdleConnections int
	ConnLifetime       time.Duration
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
		return nil, err
	}

	return &u, nil

	debugln("Loaded user:", u.Username)
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

func debugln(args ...interface{}) {
	fmt.Println(args...)
}

func debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
