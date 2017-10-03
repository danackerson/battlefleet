package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danackerson/battlefleet/websockets"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	uuid "gopkg.in/myesui/uuid.v1"
)

// https://github.com/riscie/websocket-tic-tac-toe/ <= cool ideas
// TODO store an array of sessions

// GameInfo is now commented
type GameInfo struct {
	Timestamp     string
	GameUUID      template.JS
	CommanderName string
}

const errorPage = `
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8">
    <title>Wormhole detected!</title>
    <meta name="robots" content="noindex, nofollow">
    <meta name="googlebot" content="noindex, nofollow">
    <link rel="stylesheet" href="/css/bf.css"/>
		<link rel="icon" href="/images/wormhole.png">
  </head>
  <body>
    Invalid game ID. Double check the ID or make a new game.<br/>
`

func gameHandler(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	sessionIDCookie, _ := r.Cookie(sessionID)
	session := sessions.GetSession(r)
	gameUUID := session.Get(gameUUIDKey)

	// they come in without a cookie or request a gameID that doesn't match their own
	if sessionIDCookie == nil || (gameUUID != requestParams["id"] && requestParams["id"] != newGameUUID) {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, nil)
		http.Redirect(w, r, "/", http.StatusPreconditionRequired)
		return
	}

	if requestParams["id"] == newGameUUID || gameUUID == nil {
		// TODO: if __new__ but they have a gameUUID perhaps warn they are about to lose their game?
		sessionIDString := strings.Split(sessionIDCookie.String(), "=")[1]
		newlyGameUUID := uuid.NewV5(uuid.NameSpaceOID, sessionIDString)
		session.Set(gameUUIDKey, newlyGameUUID.String())
		http.Redirect(w, r, "/games/"+newlyGameUUID.String(), http.StatusMovedPermanently)
		return
	}

	log.Printf("gameUUID: %v+\n", gameUUID)

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})

	// TODO: use for saving the gameState to MongoDB (mlab.com)
	gameInfo := GameInfo{
		Timestamp:     strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		GameUUID:      template.JS(gameUUID.(string)),
		CommanderName: "Janeway",
	}

	render.HTML(w, http.StatusOK, "game", gameInfo)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	scheme := strings.Split(r.Header.Get("Origin"), ":")[0]

	if r.Header.Get("Origin") != scheme+"://"+r.Host {
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
