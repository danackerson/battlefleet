package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/Don-V/mongostore"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var prodSession = false
var httpPort = ":8083"
var sessionCookieKey = "battlefleetID"
var cmdrNameKey = "cmdrName"
var gameUUIDKey = "gameUUID"
var newGameUUID = "__new__"
var mongoDBUser string
var mongoDBPass string
var mongoDBHost string
var mongoDBName = "fleetbattle"
var mongoCollection = "battlefleetSessions"
var sessionStore *mongostore.MongoStore
var version string

func main() {
	parseEnvVariables()

	router := mux.NewRouter()
	setUpMuxHandlers(router)
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(httpPort, n)
}

func parseEnvVariables() {
	prodSession, _ = strconv.ParseBool(os.Getenv("prodSession"))

	// if performance of mongoDB store sucks:
	// sessionStore = sessions.NewFilesystemStore("/tmp", []byte(os.Getenv("bfSecret")))
	if !prodSession {
		maxAge := 7 * 24 * 3600 // 1 week
		sessionStore = setupMongoDBSessionStore("testDBBF", "testDBBF123",
			"ds113915.mlab.com:13915", maxAge)
		sessionStore.Options = &sessions.Options{
			Path:   "/",
			Domain: "localhost",
			MaxAge: maxAge, // one week
		}
	} else {
		maxAge := 3600 * 24 * 365 // 1 year expiration
		sessionStore = setupMongoDBSessionStore(os.Getenv("mongoDBUser"),
			os.Getenv("mongoDBPass"), os.Getenv("mongoDBHost"), maxAge)
		sessionStore.Options = &sessions.Options{
			Path:     "/",
			Domain:   "battlefleet.ackerson.de",
			MaxAge:   maxAge,
			Secure:   true,
			HttpOnly: true,
		}
	}

	version = os.Getenv("CIRCLE_BUILD_NUM")
	if version == "" {
		version = "TEST"
	}
}

func setupMongoDBSessionStore(mongoDBUser string, mongoDBPass string, mongoDBHost string, mongoMaxAge int) *mongostore.MongoStore {
	mongoURL := "mongodb://" + mongoDBUser + ":" + mongoDBPass + "@" + mongoDBHost + "/" + mongoDBName
	mongoSession, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Printf("Mongo Session ERR: %v\n", err)
	}
	mongoCollection := mongoSession.DB(mongoDBName).C(mongoCollection)
	return mongostore.NewMongoStore(mongoCollection, mongoMaxAge, true, []byte(os.Getenv("bfSecret")))
}

func setUpMuxHandlers(router *mux.Router) {
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/games/{id}", gameHandler).Name("games")
	router.HandleFunc("/wsInit", serveWebSocket)
	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			versionHandler(w, r)
		}
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieKey)

	if session.Values[cmdrNameKey] == nil {
		session.Values[cmdrNameKey] = "stranger!"
	}

	if session.Values[gameUUIDKey] == nil {
		session.Values[gameUUIDKey] = newGameUUID
	}

	if e := session.Save(r, w); e != nil {
		panic(e) // for now
	}

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})

	gameInfo := GameInfo{
		GameUUID:      template.JS(session.Values[gameUUIDKey].(string)),
		Timestamp:     strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		CommanderName: session.Values[cmdrNameKey].(string),
	}
	render.HTML(w, http.StatusOK, "home", gameInfo)
}

// VersionHandler now commented
func versionHandler(w http.ResponseWriter, req *http.Request) {
	buildURL := "https://circleci.com/gh/danackerson/battlefleet/" + version
	v := map[string]string{"version": buildURL, "build": version}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
