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
  </head>
  <body>
    Invalid game ID. Double check the ID or make a new game.<br/>
`

func gameHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	sessionIDCookie, _ := r.Cookie(sessionID)
	if sessionIDCookie == nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, nil)
		http.Redirect(w, r, "/", http.StatusPreconditionRequired)
		return
	}

	vars := mux.Vars(r)
	if vars["id"] == "__new__" {
		// TODO: perhaps warn they are about to lose their game?
		session.Delete(gameUUIDKey)
	}

	// No gameUUID in session - so create a new one
	if session.Get(gameUUIDKey) == nil {
		sessionIDString := strings.Split(sessionIDCookie.String(), "=")[1]
		gameUUID := uuid.NewV5(uuid.NameSpaceOID, sessionIDString)
		session.Set(gameUUIDKey, gameUUID.String())
		http.Redirect(w, r, "/games/"+gameUUID.String(), http.StatusMovedPermanently)
	} else {
		gameUUID := session.Get(gameUUIDKey).(string)
		if gameUUID != vars["id"] { // URL doesn't match gameUUID - redirect
			http.Redirect(w, r, "/games/"+gameUUID, http.StatusMovedPermanently)
		}
	}

	// TODO: use for saving the gameState to MongoDB (mlab.com)
	gameUUIDString := session.Get(gameUUIDKey).(string)

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})
	// comment
	gameInfo := GameInfo{
		Timestamp:     strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		GameUUID:      template.JS(gameUUIDString),
		CommanderName: "Janeway",
	}

	render.HTML(w, http.StatusOK, "game", gameInfo)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
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
