package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/danackerson/battlefleet/structures"
	"github.com/danackerson/battlefleet/websockets"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
)

// https://github.com/riscie/websocket-tic-tac-toe/ <= cool ideas

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
    Invalid game ID or Commander name.<br/>
		Double check the ID or choose a different name to make a <a href="/">new game</a>.
`

var funcMap template.FuncMap

func setUpFuncMaps() {
	funcMap = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"curr_time": func() int64 {
			return time.Now().Unix()
		},
	}
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	var account *structures.Account
	session, _ := sessionStore.Get(r, sessionCookieKey)
	if session.Values[accountKey] != nil {
		account = session.Values[accountKey].(*structures.Account)
	}

	requestParams := r.URL.Query()
	if len(requestParams) > 0 && requestParams["action"][0] == "delete" {
		account.DeleteAccount()
		session.Options.MaxAge = -1
		if e := session.Save(r, w); e != nil {
			panic(e) // for now
		}
		// Session Flash msg "Account deleted"
		http.Redirect(w, r, "/?account=deleted", http.StatusTemporaryRedirect)
		return
	}
	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
		Funcs:         []template.FuncMap{funcMap},
	})

	render.HTML(w, http.StatusOK, "account", *account)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieKey)

	if session.Values[cmdrNameKey] == nil {
		session.Values[cmdrNameKey] = "stranger!"
	}

	if session.Values[gameUUIDKey] == nil {
		session.Values[gameUUIDKey] = structures.NewGameUUID
	}

	account := structures.NewAccount(session.Values[cmdrNameKey].(string))
	// retrieve account
	if session.Values[accountKey] != nil {
		account = session.Values[accountKey].(*structures.Account)
	}

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: true,
		Funcs:         []template.FuncMap{funcMap},
	})

	render.HTML(w, http.StatusOK, "home", account)
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

	account := getAccount(r, w, session)
	if account != nil {
		gameUUID := requestParams["gameid"]
		redirected := setupGame(r, w, session, account, gameUUID)

		if !redirected {
			queryParams := r.URL.Query()
			if len(queryParams) > 0 && queryParams["action"][0] == "delete" {
				account.DeleteGame(gameUUID)
				if session.Values[gameUUIDKey] == gameUUID {
					session.Values[gameUUIDKey] = ""
				}
				// remember, the session *is* the persistence store
				// a new request will fetch the account from the session on disk
				// so deleting a game is not really Deleted until the session is saved!
				if e := session.Save(r, w); e != nil {
					panic(e) // for now
				}
				http.Redirect(w, r, "/account/", http.StatusTemporaryRedirect)
				return
			}

			render := render.New(render.Options{
				Layout:        "content",
				IsDevelopment: true,
				Funcs:         []template.FuncMap{funcMap},
			})

			render.HTML(w, http.StatusOK, "game", account)
		}
	}
}

func setupGame(r *http.Request, w http.ResponseWriter,
	session *sessions.Session, account *structures.Account, gameUUID string) bool {
	redirected := false

	// they come in without a cookie or request a gameID that doesn't match their own
	if gameUUID != structures.NewGameUUID {
		if account.OwnsGame(gameUUID) {
			account.CurrentGameID = gameUUID
			session.Values[gameUUIDKey] = gameUUID
			if e := session.Save(r, w); e != nil {
				panic(e) // for now
			}
		} else {
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, nil)
			http.Redirect(w, r, "/", http.StatusPreconditionRequired)
			redirected = true
			return redirected
		}
	} else if gameUUID == structures.NewGameUUID {
		sessionIDHash := session.ID + time.Now().String()
		gameUUID = uuid.NewV5(uuid.NamespaceOID, sessionIDHash).String()
		newGame := structures.NewGame(gameUUID, account.ID)
		account.AddGame(newGame)
		session.Values[accountKey] = account
		session.Values[gameUUIDKey] = gameUUID
		if e := session.Save(r, w); e != nil {
			panic(e) // for now
		}
		http.Redirect(w, r, "/games/"+gameUUID, http.StatusMovedPermanently)
		redirected = true
		return redirected
	}

	return redirected
}

func getAccount(r *http.Request, w http.ResponseWriter, session *sessions.Session) *structures.Account {
	var account *structures.Account

	if session.Values[accountKey] == nil {
		if r.FormValue("cmdrName") == "" || r.FormValue("cmdrName") == "stranger!" {
			// new accounts require a CommanderName and 'stranger!' is reserved ;)
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, nil)
			http.Redirect(w, r, "/", http.StatusPreconditionRequired)
			return nil
		}

		account = structures.NewAccount(r.FormValue("cmdrName"))
		session.Values[accountKey] = account
	} else {
		account = session.Values[accountKey].(*structures.Account)
	}

	return account
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
