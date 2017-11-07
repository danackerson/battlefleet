package main

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/danackerson/battlefleet/structures"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

var httpPort = ":8083"
var sessionCookieKey = "battlefleetID"
var accountIDKey = "ownerAccountID"
var accountKey = "ownerAccount"
var cmdrNameKey = "cmdrName"
var defaultCmdrName = "stranger!"
var gameUUIDKey = "gameUUID"

var mongoDBUser string
var mongoDBPass string
var mongoDBHost string
var mongoDBName = "fleetbattle"
var mongoDBCollection = "sessions"
var db *mgo.Session
var version string
var sessionStore *sessions.FilesystemStore
var prodSession = false

// Auth0Data for storing Auth0 variables
type Auth0Data struct {
	Auth0ClientID          string
	Auth0ClientSecret      string
	Auth0Domain            string
	Auth0CallbackURLString string
	Auth0CallbackURL       template.URL
}

var auth0data Auth0Data

func main() {
	prepareSessionEnvironment()

	router := mux.NewRouter()
	setUpMuxHandlers(router)
	setUpFuncMaps()
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(httpPort, n)
}

func prepareSessionEnvironment() {
	gob.Register(&structures.Account{})
	gob.Register(&structures.Game{})
	gob.Register(&structures.Ship{})

	prodSession, _ = strconv.ParseBool(os.Getenv("prodSession"))
	sessionStore = sessions.NewFilesystemStore("/tmp", []byte(os.Getenv("bfSecret")))
	sessionStore.MaxLength(32 * 1024) // else securecookie: value too long error

	if !prodSession {
		maxAge := 7 * 24 * 3600 // 1 week
		sessionStore.Options = &sessions.Options{
			Path:   "/",
			Domain: "localhost",
			MaxAge: maxAge, // one week
		}

		// load test vars from .env
		err := godotenv.Load()
		if err != nil {
			log.Printf("env var load failure: %s", err.Error())
		}
	} else {
		maxAge := 3600 * 24 * 365 // 1 year expiration
		sessionStore.Options = &sessions.Options{
			Path:     "/",
			Domain:   "battlefleet.ackerson.de",
			MaxAge:   maxAge,
			Secure:   true,
			HttpOnly: true,
		}
	}

	auth0data = Auth0Data{
		os.Getenv("AUTH0_CLIENT_ID"),
		os.Getenv("AUTH0_CLIENT_SECRET"),
		os.Getenv("AUTH0_DOMAIN"),
		os.Getenv("AUTH0_CALLBACK_URL"),
		template.URL(os.Getenv("AUTH0_CALLBACK_URL")),
	}

	version = os.Getenv("CIRCLE_BUILD_NUM")
	if version == "" {
		version = "TEST"
	}

	setupMongoDBSession()
}

func setupMongoDBSession() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("mongoDBHost")},
		Timeout:  10 * time.Second,
		Database: os.Getenv("mongoDBName"),
		Username: os.Getenv("mongoDBUser"),
		Password: os.Getenv("mongoDBPass"),
	}

	db, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Printf("Cannot Dial Mongo: %s", err.Error())
	}
	db.SetMode(mgo.Monotonic, true)
}

func setUpMuxHandlers(router *mux.Router) {
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/callback", callbackHandler)
	router.HandleFunc("/games/{gameid}", gameHandler).Name("games")
	router.HandleFunc("/account/", accountHandler)
	router.HandleFunc("/wsInit", serveWebSocket)
	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			versionHandler(w, r)
		}
	})
}
