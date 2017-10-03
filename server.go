package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var prodSession = false
var httpPort = ":8083"
var sessionID = "battlefleetID"
var gameUUIDKey = "gameUUID"
var newGameUUID = "__new__"
var mongoDBUser string
var mongoDBPass string
var mongoDBHost string
var version string
var secret string

func parseEnvVariables() {
	prodSession, _ = strconv.ParseBool(os.Getenv("prodSession"))
	log.Printf("%t", prodSession)
	mongoDBUser = os.Getenv("mongoDBUser")
	mongoDBPass = os.Getenv("mongoDBPass")
	mongoDBHost = os.Getenv("mongoDBHost")
	version = os.Getenv("CIRCLE_BUILD_NUM")
	secret = os.Getenv("bfSecret")
}

func main() {
	parseEnvVariables()

	router := mux.NewRouter()
	setUpMuxHandlers(router)
	n := negroni.Classic()

	store := cookiestore.New([]byte(secret))

	if !prodSession {
		store.Options(sessions.Options{
			Path:   "/",
			Domain: "localhost",
			MaxAge: 30 * 89280, // one month
		})
		log.Println("NOT prodSession")
	} else {
		store.Options(sessions.Options{
			Path:     "/",
			Domain:   "battlefleet.ackerson.de",
			MaxAge:   30 * 89280, // one month
			Secure:   true,
			HTTPOnly: true,
		})

		log.Println("prodSession")
	}

	n.Use(sessions.Sessions(sessionID, store))
	n.UseHandler(router)

	http.ListenAndServe(httpPort, n)
}

func setUpMuxHandlers(router *mux.Router) {
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/games/{id}", gameHandler).Name("games")
	router.HandleFunc("/wsInit", serveWebSocket)
	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			VersionHandler(w, r)
		}
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Cache-Control", "max-age=10800")

	session := sessions.GetSession(r)

	cmdrName := "n/a"
	gameUUID := session.Get(gameUUIDKey)
	if gameUUID != nil {
		// TODO search for gameID in mongoDB and reload state
		// perhaps even redirect to /games/{gameUUID} ?
		if session.Get("cmdrName") == nil {
			cmdrName = "Captain Janeway"
		}
	} else {
		gameUUID = newGameUUID
	}
	session.Set(gameUUIDKey, gameUUID.(string))

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})

	gameInfo := GameInfo{
		GameUUID:      template.JS(gameUUID.(string)),
		Timestamp:     strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		CommanderName: cmdrName,
	}
	render.HTML(w, http.StatusOK, "home", gameInfo)
}

// VersionHandler now commenteds
func VersionHandler(w http.ResponseWriter, req *http.Request) {
	buildURL := "https://circleci.com/gh/danackerson/battlefleet/" + version
	v := map[string]string{"version": buildURL, "build": version}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
