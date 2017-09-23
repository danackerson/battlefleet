package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/danackerson/battlefleet/websockets"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var httpPort = ":8083"

func main() {
	parseEnvVariables()

	mux := http.NewServeMux()
	setUpMuxHandlers(mux)
	n := negroni.Classic()

	store := cookiestore.New([]byte(secret))
	store.Options(sessions.Options{
		//Path:   "battlefleet",
		//Domain: "ackerson.de",
		MaxAge: 3600,
	})

	n.Use(sessions.Sessions("gurkherpadab", store))
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

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if session != nil && session.Get("hello") != nil {
		log.Println("ws session: " + session.Get("hello").(string))
	}

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

	// TODO get/set cookie & find/create session

	go websockets.ServerTime(ws)
}
