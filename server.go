package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/negroni"
)

var prodSession = false
var httpPort = ":8083"
var sessionCookieKey = "battlefleetID"
var accountIDKey = "accountID"
var cmdrNameKey = "cmdrName"
var gameUUIDKey = "gameUUID"
var newGameUUID = "__new__"
var mongoDBUser string
var mongoDBPass string
var mongoDBHost string
var mongoDBName = "fleetbattle"
var mongoCollection = "battlefleetSessions"
var sessionStore *sessions.FilesystemStore
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
	sessionStore = sessions.NewFilesystemStore("/tmp", []byte(os.Getenv("bfSecret")))

	if !prodSession {
		maxAge := 7 * 24 * 3600 // 1 week
		sessionStore.Options = &sessions.Options{
			Path:   "/",
			Domain: "localhost",
			MaxAge: maxAge, // one week
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

	version = os.Getenv("CIRCLE_BUILD_NUM")
	if version == "" {
		version = "TEST"
	}
}

// Issue was connectionTimeouts. Fallback to FileSystemStore and a sane
// method of periodically storing these to mongoDB
// (e.g. on user manually quiting/saving or browser close)
func setupMongoDBSessionStore(mongoDBUser string, mongoDBPass string, mongoDBHost string, mongoMaxAge int) {
	// test: setupMongoDBSessionStore("testDBBF", "testDBBF123", "ds113915.mlab.com:13915", maxAge)
	// prod: setupMongoDBSessionStore(os.Getenv("mongoDBUser"), os.Getenv("mongoDBPass"), os.Getenv("mongoDBHost"), maxAge)
	/*mongoURL := "mongodb://" + mongoDBUser + ":" + mongoDBPass + "@" + mongoDBHost + "/" + mongoDBName
	mongoSession, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Printf("Mongo Session ERR: %v\n", err)
	}
	mongoCollection := mongoSession.DB(mongoDBName).C(mongoCollection)
	return mongostore.NewMongoStore(mongoCollection, mongoMaxAge, true, []byte(os.Getenv("bfSecret")))*/
}

func setUpMuxHandlers(router *mux.Router) {
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/games/{gameid}", gameHandler).Name("games")
	router.HandleFunc("/wsInit", serveWebSocket)
	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			versionHandler(w, r)
		}
	})
}
