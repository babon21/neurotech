package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	_ "github.com/volatiletech/authboss/auth"
	aboauth "github.com/volatiletech/authboss/oauth2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/volatiletech/authboss/defaults"

	"gopkg.in/mgo.v2"
)

type DbStorer struct {
	Db *sqlx.DB
}

var (
	ab           = authboss.New()
	database     DbStorer
	sessionStore SessionStorer
	cookieStore  abclientstate.CookieStorer
)

const (
	sessionCookieName = "ab_blog"
	port = "3000"
)

const schema = `CREATE TABLE users (
    username text PRIMARY KEY NOT NULL,
	password text NULL,
	oauth2_uid text NULL,
	oauth2_provider text NULL)`

// func registerUser() {
// 	password := "super"
// 	pass, err := bcrypt.GenerateFromPassword([]byte(password), ab.Config.Modules.BCryptCost)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = database.Save(nil, &User{
// 		UserID: "kelan007",
// 		Password: string(pass)})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func initSchema(db *sqlx.DB) {
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}

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

	
	// initSchema(db)
	// registerUser()

	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = true

	sessionStore = NewSessionStorer(sessionCookieName, sessionStoreKey, nil)

	setupAuthboss()

	mux := chi.NewRouter()
	mux.Use(ab.LoadClientStateMiddleware)

	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.ModuleListMiddleware(ab))
		mux.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})

	setSiteHandlers(mux, databaseSite)

	// Start the server
	// port := os.Getenv("PORT")
	// if len(port) == 0 {
		// port = "8002"
	// }
	log.Printf("Listening on localhost: %s", port)
	log.Println(http.ListenAndServe("localhost:"+port, mux))
}

func setupAuthboss() {
	ab.Config.Paths.RootURL = "http://localhost:" + port

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")
	// ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}	

	defaults.SetCore(&ab.Config, false, false)

	oauthcreds := struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	}{}

	// Set up Google OAuth2 if we have credentials in the
	// file oauth2.toml for it.
	_, err := toml.DecodeFile("oauth2.toml", &oauthcreds)
	if err == nil && len(oauthcreds.ClientID) != 0 && len(oauthcreds.ClientSecret) != 0 {
		fmt.Println("oauth2.toml exists, configuring google oauth2")
		ab.Config.Modules.OAuth2Providers = map[string]authboss.OAuth2Provider{
			"google": {
				OAuth2Config: &oauth2.Config{
					ClientID:     oauthcreds.ClientID,
					ClientSecret: oauthcreds.ClientSecret,
					Scopes:       []string{`profile`, `email`},
					Endpoint:     google.Endpoint,
				},
				FindUserDetails: aboauth.GoogleUserDetails,
			},
		}
	} else if os.IsNotExist(err) {
		fmt.Println("oauth2.toml doesn't exist, not registering oauth2 handling")
	} else {
		fmt.Println("error loading oauth2.toml:", err)
	}

	if err := ab.Init(); err != nil {
		panic(err)
	}
}

func setSiteHandlers(mux *chi.Mux, databaseSite *mgo.Database) {
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
}

type StorageOpts struct {
	ConnString         string
	MaxConnections     int
	MaxIdleConnections int
	ConnLifetime       time.Duration
}
