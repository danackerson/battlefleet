package main

import (
	"encoding/json"
	"html/template"
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

var httpPort = ":8083"
var sessionID = "battlefleetID"
var gameUUIDKey = "gameUUID"
var mongoDBUser string
var mongoDBPass string
var mongoDBHost string
var version string
var secret string

func main() {
	parseEnvVariables()

	router := mux.NewRouter()
	setUpMuxHandlers(router)
	n := negroni.Classic()

	store := cookiestore.New([]byte(secret))
	store.Options(sessions.Options{
		Path:   "/",
		Domain: "battlefleet.ackerson.de",
		MaxAge: 2678400, // one month
	})

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
		gameUUID = "__new__"
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

func parseEnvVariables() {
	mongoDBUser = os.Getenv("mongoDBUser")
	mongoDBPass = os.Getenv("mongoDBPass")
	mongoDBHost = os.Getenv("mongoDBHost")
	version = os.Getenv("CIRCLE_BUILD_NUM")
	secret = os.Getenv("bfSecret")
}

// VersionHandler now commenteds
func VersionHandler(w http.ResponseWriter, req *http.Request) {
	buildURL := "https://circleci.com/gh/danackerson/battlefleet/" + version
	v := map[string]string{"version": buildURL, "build": version}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
