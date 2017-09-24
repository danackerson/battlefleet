package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	uuid "gopkg.in/myesui/uuid.v1"

	"github.com/danackerson/battlefleet/websockets"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var httpPort = ":8083"
var sessionID = "battlefleetID"

func main() {
	parseEnvVariables()

	mux := http.NewServeMux()
	setUpMuxHandlers(mux)
	n := negroni.Classic()

	store := cookiestore.New([]byte(secret))
	store.Options(sessions.Options{
		//Path: "battlefleet",
		//Domain: "ackerson.de",
		MaxAge: 2678400, // one month
	})

	n.Use(sessions.Sessions(sessionID, store))
	n.UseHandler(mux)

	http.ListenAndServe(httpPort, n)

}

var mongoDBUser string
var mongoDBPass string
var mongoDBHost string
var version string
var secret string

func parseEnvVariables() {
	mongoDBUser = os.Getenv("mongoDBUser")
	mongoDBPass = os.Getenv("mongoDBPass")
	mongoDBHost = os.Getenv("mongoDBHost")
	version = os.Getenv("CIRCLE_BUILD_NUM")
	secret = os.Getenv("bfSecret")
}

func setUpMuxHandlers(mux *http.ServeMux) {
	post := "POST"

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == post {
			VersionHandler(w, r)
		}
	})

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/wsInit", serveWebSocket)
}

// VersionHandler now commenteds
func VersionHandler(w http.ResponseWriter, req *http.Request) {
	buildURL := "https://circleci.com/gh/danackerson/battlefleet/" + version
	v := map[string]string{"version": buildURL, "build": version}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=10800")

	session := sessions.GetSession(r)
	session.Set("hello", "world") // e.g. login, gameID details?

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})

	render.HTML(w, http.StatusOK, "home", "")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// https://github.com/riscie/websocket-tic-tac-toe/ <= cool ideas

// TODO store an array of sessions
func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if session.Get("gameUUID") == nil {
		sessionIDCookie, _ := r.Cookie(sessionID)
		sessionIDString := strings.Split(sessionIDCookie.String(), "=")[1]
		gameUUID := uuid.NewV5(uuid.NameSpaceOID, sessionIDString)

		session.Set("gameUUID", gameUUID.String())
	}
	// TODO: use for saving the gameState to MongoDB (mlab.com)
	gameUUIDString := session.Get("gameUUID").(string)

	// TODO: provide to end-users for "bookmarking/sharing" their game
	encodedGameUUID := base64.StdEncoding.EncodeToString([]byte(gameUUIDString))

	log.Println("Game UUID: " + gameUUIDString)
	log.Println("Game UUID encoded: " + encodedGameUUID)

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	ws, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println("WS handshake + " + err.Error())
		}
		return
	}

	go websockets.ServerTime(ws)
}
