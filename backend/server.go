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
	"github.com/volatiletech/authboss"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	_ "github.com/volatiletech/authboss/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/volatiletech/authboss/defaults"

	"gopkg.in/mgo.v2"
)

var (
	ab           = authboss.New()
	database     = NewMemStorer()
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

	password := "super"
	pass, err := bcrypt.GenerateFromPassword([]byte(password), ab.Config.Modules.BCryptCost)
	if err != nil {
		panic(err)
	}

	database.Save(nil, &User{
		Name:     "Batman",
		Email:    "kek@mail.ru",
		Password: string(pass)})

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
	ID int

	// Non-authboss related field
	Name string

	// Auth
	Email    string
	Password string
}

func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"rick@councilofricks.com": {
				ID:       1,
				Name:     "Rick",
				Password: "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Email:    "rick@councilofricks.com",
			},
		},
		Tokens: make(map[string][]string),
	}
}

type MemStorer struct {
	Users  map[string]User
	Tokens map[string][]string
}

func (m MemStorer) Save(_ context.Context, user authboss.User) error {
	u := user.(*User)
	m.Users[u.Email] = *u

	debugln("Saved user:", u.Name)
	return nil
}

func (m MemStorer) Load(_ context.Context, key string) (user authboss.User, err error) {
	u, ok := m.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	debugln("Loaded user:", u.Name)
	return &u, nil
}

// New user creation
func (m MemStorer) New(_ context.Context) authboss.User {
	return &User{}
}

// Create the user
func (m MemStorer) Create(_ context.Context, user authboss.User) error {
	u := user.(*User)

	if _, ok := m.Users[u.Email]; ok {
		return authboss.ErrUserFound
	}

	debugln("Created new user:", u.Name)
	m.Users[u.Email] = *u
	return nil
}

// GetPID from user
func (u User) GetPID() string { return u.Email }

// PutPID into user
func (u *User) PutPID(pid string) { u.Email = pid }

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
