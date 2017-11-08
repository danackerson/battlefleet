package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"golang.org/x/oauth2"

	"github.com/danackerson/battlefleet/structures"
	"github.com/danackerson/battlefleet/websockets"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
)

// https://github.com/riscie/websocket-tic-tac-toe/ <= cool ideas

// TemplateVars as helper for rendering pages
type TemplateVars struct {
	Account *structures.Account
	Data    Auth0Data
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
    <br/><br/><h3>You've encountered a wormhole - <a href="/">return to base</a> immediately!</h3><br/></br>
		<hr><br/><br/>
		If this problem persists, please forward this msg to admin@ackerson.de:<br/><br/>
		{{ . }}<br/><br/>
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

		requestParams := r.URL.Query()
		if len(requestParams) > 0 {
			if requestParams["action"][0] == "delete" {
				account.DeleteAccount(db)
				session.Options.MaxAge = -1
				if e := session.Save(r, w); e != nil {
					t, _ := template.New("errorPage").Parse(errorPage)
					t.Execute(w, "saveSession1: "+e.Error())
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}
				// Session Flash msg "Account deleted"
				http.Redirect(w, r, "/?account=deleted", http.StatusTemporaryRedirect)
				return
			} else if requestParams["action"][0] == "logout" {
				account.EndSession(db)
				session.Options.MaxAge = -1
				if e := session.Save(r, w); e != nil {
					t, _ := template.New("errorPage").Parse(errorPage)
					t.Execute(w, "saveSession2: "+e.Error())
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}

				scheme := "http"
				if prodSession {
					scheme = "https"
				}
				returnHost := strings.Replace(r.Host, ":", "%3A", 1)
				returnTo := "?returnTo=" + scheme + "%3A%2F%2F" + returnHost + "&client_id=" + auth0data.Auth0ClientID
				http.Redirect(w, r, "https://battlefleet.eu.auth0.com/v2/logout"+returnTo, http.StatusTemporaryRedirect)
				return
			}

		}

		render := render.New(render.Options{
			Layout:        "content",
			IsDevelopment: !prodSession,
			Funcs:         []template.FuncMap{funcMap},
		})

		render.HTML(w, http.StatusOK, "account", &TemplateVars{account, auth0data})
	} else {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "You are currently not logged in!")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieKey)

	if session.Values[cmdrNameKey] == nil {
		session.Values[cmdrNameKey] = defaultCmdrName
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
		IsDevelopment: !prodSession,
		Funcs:         []template.FuncMap{funcMap},
	})

	render.HTML(w, http.StatusOK, "home", &TemplateVars{account, auth0data})
}

// VersionHandler now commented
func versionHandler(w http.ResponseWriter, req *http.Request) {
	buildURL := "https://circleci.com/gh/danackerson/battlefleet/" + version
	v := map[string]string{"version": buildURL, "build": version}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	parseErr := r.ParseForm()
	if parseErr != nil {
		log.Println(parseErr)
	}

	session, sessionErr := sessionStore.Get(r, sessionCookieKey)
	if sessionErr != nil {
		if strings.Contains(sessionErr.Error(), "no such file or directory") {
			// don't panic
			log.Printf("Creating new session: %s", session.ID)
		} else {
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, "getSession: "+sessionErr.Error())
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
	}

	account := getAccount(r, w, session)
	if account != nil {
		gameUUID := requestParams["gameid"]
		action, ok := r.URL.Query()["action"]
		if ok && len(action) > 0 && action[0] == "delete" {
			account.DeleteGame(gameUUID)
			if session.Values[gameUUIDKey] == gameUUID {
				session.Values[gameUUIDKey] = ""
			}
			// remember, the session *is* the persistence store
			// a new request will fetch the account from the session on disk
			// so deleting a game is not really Deleted until the session is saved!
			if e := session.Save(r, w); e != nil {
				t, _ := template.New("errorPage").Parse(errorPage)
				t.Execute(w, "saveSession: "+e.Error())
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/account/", http.StatusTemporaryRedirect)
			return
		}

		redirected := setupGame(r, w, session, account, gameUUID)
		if !redirected {
			render := render.New(render.Options{
				Layout:        "content",
				IsDevelopment: !prodSession,
				Funcs:         []template.FuncMap{funcMap},
			})

			render.HTML(w, http.StatusOK, "game", &TemplateVars{account, auth0data})
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
				t, _ := template.New("errorPage").Parse(errorPage)
				t.Execute(w, "saveSession1: "+e.Error())
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				redirected = true
				return redirected
			}
		} else {
			t, _ := template.New("errorPage").Parse(errorPage)
			log.Println("hello?!")
			errorString := "You neither own Game ID:<span style='color:orange;'>" + gameUUID + "</span> nor have you been invited to join."
			t.Execute(w, template.JS(errorString))
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
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, "saveSession2: "+e.Error())
			http.Redirect(w, r, "/", http.StatusPreconditionRequired)
			redirected = true
			return redirected
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
		if r.FormValue("cmdrName") == "" || r.FormValue("cmdrName") == defaultCmdrName {
			// new accounts require a CommanderName and 'stranger!' is reserved ;)
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, "New accounts require a Commander name and '"+defaultCmdrName+"' is not allowed.")
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

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	domain := auth0data.Auth0Domain
	conf := &oauth2.Config{
		ClientID:     auth0data.Auth0ClientID,
		ClientSecret: auth0data.Auth0ClientSecret,
		RedirectURL:  auth0data.Auth0CallbackURLString,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "Auth0 token: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "Auth0 client: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err = json.Unmarshal(raw, &profile); err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	session, err := sessionStore.Get(r, sessionCookieKey)
	if err != nil {
		origError := err.Error()
		if strings.Contains(err.Error(), "no such file or directory") {
			// probably have an old session cookie so recreate
			session, err = sessionStore.New(r, sessionCookieKey)
			if err != nil {
				t, _ := template.New("errorPage").Parse(errorPage)
				t.Execute(w, "recreateSession: "+err.Error()+" (after '"+origError+"')")
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				return
			}
		} else {
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, "getSession: "+err.Error())
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
	}

	var accountFound *structures.Account
	mongoSession := db.Copy()
	defer mongoSession.Close()
	c := mongoSession.DB("fleetbattle").C("sessions")
	err = c.Find(bson.M{"auth0profile.sub": profile["sub"]}).One(&accountFound)
	if err != nil && err.Error() == "not found" {
		log.Printf("didn't find '%s' in mongoDB\n", profile["name"])
		accountFound = structures.NewAccount(profile["nickname"].(string))
	} else if err != nil && err.Error() != "not found" {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "MongoDB: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	var accountPlaying *structures.Account
	if session.Values[accountKey] != nil {
		accountPlaying = session.Values[accountKey].(*structures.Account)
		if accountPlaying.Commander != defaultCmdrName {
			accountFound.Commander = accountPlaying.Commander
		}
		accountFound.Games = append(accountFound.Games, accountPlaying.Games...)
	}

	accountFound.Auth0Token = token
	accountFound.Auth0Profile = profile
	accountFound.LastLogin = time.Now()

	session.Values[accountKey] = accountFound
	err = session.Save(r, w)
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "saveSession: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/account/", http.StatusSeeOther)
}
