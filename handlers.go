package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danackerson/battlefleet/structures"
	"github.com/danackerson/battlefleet/websockets"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieKey)

	if session.Values[cmdrNameKey] == nil {
		session.Values[cmdrNameKey] = "stranger!"
	}

	if session.Values[gameUUIDKey] == nil {
		session.Values[gameUUIDKey] = newGameUUID
	}

	accountID := session.Values[accountIDKey]
	if accountID != nil {
		account := structures.GetAccount(session.Values[accountIDKey].(string))
		session.Values[cmdrNameKey] = account.Commander
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

func loginHander(w http.ResponseWriter, r *http.Request) {
	// https://manage.auth0.com/#/ => Github account
	// https://manage.auth0.com/#/clients/Zz58wpKP7ApF0Pw6KGxE35XYecf2sCEO/quickstart => Go
	// https://auth0.com/blog/vuejs2-authentication-tutorial/
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	parseErr := r.ParseForm()
	if parseErr != nil {
		log.Println(parseErr)
	}

	session, sessionErr := sessionStore.Get(r, sessionCookieKey)
	if sessionErr != nil {
		panic(sessionErr)
	}

	accountID := ""
	if session.Values[accountIDKey] != nil {
		accountID = session.Values[accountIDKey].(string)
	}

	var account *structures.Account
	if accountID == "" {
		if r.FormValue("cmdrName") == "" {
			// new accounts require a CommanderName
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, nil)
			http.Redirect(w, r, "/", http.StatusPreconditionRequired)
			return
		}

		account = structures.NewAccount(r.FormValue("cmdrName"))
		session.Values[accountIDKey] = account.ID
		if e := session.Save(r, w); e != nil {
			panic(e) // for now
		}
	} else {
		account = structures.GetAccount(accountID)
	}
	gameUUID := session.Values[gameUUIDKey].(string)

	// they come in without a cookie or request a gameID that doesn't match their own
	if requestParams["gameid"] != newGameUUID && gameUUID != requestParams["gameid"] {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, nil)
		http.Redirect(w, r, "/", http.StatusPreconditionRequired)
		return
	}

	if requestParams["gameid"] == newGameUUID || gameUUID == "" {
		sessionIDHash := session.ID + time.Now().String()
		gameUUID = uuid.NewV5(uuid.NamespaceOID, sessionIDHash).String()
		account.AddGame(structures.NewGame(gameUUID, account.ID))
		session.Values[gameUUIDKey] = gameUUID
		if e := session.Save(r, w); e != nil {
			panic(e) // for now
		}
		http.Redirect(w, r, "/games/"+gameUUID, http.StatusMovedPermanently)
		return
	}

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
	})

	// TODO: use for saving the gameState to MongoDB (mlab.com)
	gameInfo := GameInfo{
		Timestamp:     strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		GameUUID:      template.JS(gameUUID),
		CommanderName: account.Commander,
	}

	render.HTML(w, http.StatusOK, "game", gameInfo)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	serverPort := httpPort
	remoteHostSettings := strings.Split(r.Host, ":")
	if len(remoteHostSettings) > 1 {
		serverPort = remoteHostSettings[1]
	}
	scheme := strings.Split(r.Header.Get("Origin"), ":")[0]

	// test server runs on different port
	if serverPort == httpPort && r.Header.Get("Origin") != scheme+"://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	ws, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println("error: " + err.Error())
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println("WS handshake + " + err.Error())
		}
		return
	}

	go websockets.ServerTime(ws)
}
