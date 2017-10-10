package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

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
var version string
var sessionStore *sessions.FilesystemStore

func parseEnvVariables() {
	prodSession, _ = strconv.ParseBool(os.Getenv("prodSession"))
	sessionStore = sessions.NewFilesystemStore("/tmp", []byte(os.Getenv("bfSecret")))

	if !prodSession {
		sessionStore.Options = &sessions.Options{
			Path:   "/",
			Domain: "localhost",
			MaxAge: 30 * 89280, // one month
		}
	} else {
		sessionStore.Options = &sessions.Options{
			Path:     "/",
			Domain:   "battlefleet.ackerson.de",
			MaxAge:   30 * 89280, // one month
			Secure:   true,
			HttpOnly: true,
		}
	}

	mongoDBUser = os.Getenv("mongoDBUser")
	mongoDBPass = os.Getenv("mongoDBPass")
	mongoDBHost = os.Getenv("mongoDBHost")
	version = os.Getenv("CIRCLE_BUILD_NUM")
	if version == "" {
		version = "TEST"
	}
}

func main() {
	parseEnvVariables()

	router := mux.NewRouter()
	setUpMuxHandlers(router)
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(httpPort, n)
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

	if session.Values[gameUUIDKey] != nil {
		// TODO search for gameID in mongoDB and reload state
		// perhaps even redirect to /games/{gameUUID} ?
	} else {
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

// VersionHandler now commenteds
func versionHandler(w http.ResponseWriter, req *http.Request) {
	buildURL := "https://circleci.com/gh/danackerson/battlefleet/" + version
	v := map[string]string{"version": buildURL, "build": version}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
