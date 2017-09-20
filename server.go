package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danackerson/battlefleet/websockets"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var httpPort = ":8083"
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	parseEnvVariables()

	mux := http.NewServeMux()
	setUpMuxHandlers(mux)
	n := negroni.Classic()

	store := cookiestore.New([]byte(secret))
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

	mux.HandleFunc("/test", serveTest)
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/ws", serveWS)
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
	hello := "hello"

	w.Header().Set("Cache-Control", "max-age=10800")
	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})

	render.HTML(w, http.StatusOK, "test", hello)
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}

	go websockets.Writer(ws, lastMod)
	websockets.Reader(ws)
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/test" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p, lastMod, err := websockets.ReadFileIfModified(time.Time{})
	if err != nil {
		p = []byte(err.Error())
		lastMod = time.Unix(0, 0)
	}
	var v = struct {
		Host    string
		Data    string
		LastMod string
	}{
		r.Host,
		string(p),
		strconv.FormatInt(lastMod.UnixNano(), 16),
	}
	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})

	render.HTML(w, http.StatusOK, "fileModified", &v)
}
