package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	_ "github.com/volatiletech/authboss/auth"
	"golang.org/x/crypto/bcrypt"

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
)

const schema = `CREATE TABLE users (
    username text PRIMARY KEY NOT NULL,
    password text NOT NULL)`

func registerUser() {
	password := "super"
	pass, err := bcrypt.GenerateFromPassword([]byte(password), ab.Config.Modules.BCryptCost)
	if err != nil {
		panic(err)
	}

	err = database.Save(nil, &User{
		Username: "kelan007",
		Password: string(pass)})

	if err != nil {
		log.Fatal(err)
	}
}

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

	newsCollection := InitNewsCollection(databaseSite)
	newsHandler := &NewsHandler{Collection: newsCollection}
	setNewsRouter(mux, newsHandler)

	publicationCollection := InitPublicationsCollection(databaseSite)
	publicationHandler := &PublicationHandler{
		Collection: publicationCollection,
		Path: "publications/",
	}
	setPublicationsRouter(mux, publicationHandler)

	studenWorkCollection := InitStudentWorksCollection(databaseSite)
	studenWorkHandler := &StudentWorkHandler{Collection: studenWorkCollection}
	setStudentWorksRouter(mux, studenWorkHandler)

	disciplineCollection := InitDisciplineCollection(databaseSite)
	disciplineHandler := &DisciplineHandler{
		Collection: disciplineCollection,
		Path:       "disciplines/",
	}
	setDisciplineRouter(mux, disciplineHandler)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Printf("Listening on localhost: %s", port)
	log.Println(http.ListenAndServe("localhost:"+port, mux))
}

func setupAuthboss() {
	ab.Config.Paths.RootURL = "http://localhost:8080"

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	// ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}
	ab.Paths.AuthLoginOK = "http://localhost:3000/"

	defaults.SetCore(&ab.Config, false, false)

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		UseUsername: true,
	}

	if err := ab.Init(); err != nil {
		panic(err)
	}
}

func setNewsRouter(r *chi.Mux, newsHandler *NewsHandler) {
	r.Route("/news", func(r chi.Router) {
		r.Get("/", newsHandler.GetNewsList)

		r.Group(func(r chi.Router) {
			r.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
			r.Post("/", newsHandler.CreateNews) // POST /articles
			r.With().Route("/{id}", func(r chi.Router) {
				r.Get("/", newsHandler.GetOneNews)    // GET /news/123
				r.Put("/", newsHandler.UpdateNews)    // PUT /news/123
				r.Delete("/", newsHandler.DeleteNews) // DELETE /news/123
			})
		})
	})
}

func setPublicationsRouter(r *chi.Mux, publicationsHandler *PublicationHandler) {
	r.Route("/publications", func(r chi.Router) {
		r.Get("/", publicationsHandler.GetPublicationsList)

		r.Group(func(r chi.Router) {
			r.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
			r.Post("/", publicationsHandler.CreatePublication) // POST /publications
			r.With().Route("/{id}", func(r chi.Router) {
				r.Get("/", publicationsHandler.GetOnePublication)    // GET /publications/123
				r.Put("/", publicationsHandler.UpdatePublication)    // PUT /publications/123
				r.Delete("/", publicationsHandler.DeletePublication) // DELETE /publications/123
				r.Post("/files", publicationsHandler.UploadPublicationFile) // PUT /disciplines/123
			})
		})
	})
}

func setStudentWorksRouter(r *chi.Mux, studentWorkHandler *StudentWorkHandler) {
	r.Route("/student-works", func(r chi.Router) {
		r.Get("/", studentWorkHandler.GetStudentWorkList)

		r.Group(func(r chi.Router) {
			r.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
			r.Post("/", studentWorkHandler.CreateStudentWork) // POST /student-works
			r.With().Route("/{id}", func(r chi.Router) {
				r.Get("/", studentWorkHandler.GetOneStudentWork)    // GET /student-works/123
				r.Put("/", studentWorkHandler.UpdateStudentWork)    // PUT /student-works/123
				r.Delete("/", studentWorkHandler.DeleteStudentWork) // DELETE /student-works/123
			})
		})
	})
}

func setDisciplineRouter(r *chi.Mux, disciplinesHandler *DisciplineHandler) {
	r.Route("/disciplines", func(r chi.Router) {
		r.Get("/", disciplinesHandler.GetDisciplineList)

		r.Group(func(r chi.Router) {
			r.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
			r.Post("/", disciplinesHandler.CreateDiscipline) // POST /disciplines
			r.With().Route("/{id}", func(r chi.Router) {
				r.Get("/", disciplinesHandler.GetOneDiscipline)            // GET /disciplines/123
				r.Put("/", disciplinesHandler.UpdateDiscipline)            // PUT /disciplines/123
				r.Delete("/", disciplinesHandler.DeleteDiscipline)         // DELETE /disciplines/123
				r.Post("/files", disciplinesHandler.UploadDisciplineFiles) // PUT /disciplines/123
			})
		})
	})
}

type StorageOpts struct {
	ConnString         string
	MaxConnections     int
	MaxIdleConnections int
	ConnLifetime       time.Duration
}
