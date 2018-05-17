package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/danackerson/battlefleet/app"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	render "github.com/unrolled/render"
)

var renderer *render.Render

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

// ErrorType is a JSON obj for error msgs sent to browser
type errorType struct {
	HTTPCode int
	Error    string
}

type infoType struct {
	HTTPCode int
	Message  string
}

// AccountType is a JSON obj for minimal Account info sent to browser
type accountType struct {
	ID        string
	Commander string
	Auth0     string
}

// GameType is a JSON obj for minimal Game info sent to browser
type gameType struct {
	ID       string
	Account  accountType
	GridSize float64
	Message  string
}

func setupCORSOptions(w http.ResponseWriter) {
	if !app.ProdSession {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8443")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
}

func sendSuccess(w http.ResponseWriter, httpCode int, errorMsg string) {
	var infoJSON infoType
	infoJSON.HTTPCode = httpCode
	infoJSON.Message = errorMsg
	json.NewEncoder(w).Encode(infoJSON)
}

func sendError(w http.ResponseWriter, httpCode int, errorMsg string) {
	var errorJSON errorType
	errorJSON.HTTPCode = httpCode
	errorJSON.Error = errorMsg
	json.NewEncoder(w).Encode(errorJSON)
}

// RetrieveSession fetches existing session store, or, if unavailable, recreates
func RetrieveSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := app.SessionStore.Get(r, app.SessionCookieKey)
	if err != nil || session.ID == "" {
		origError := err
		log.Println(origError)
		session.Options.MaxAge = -1
		session.Save(r, w) // destroys the session

		session.Options.MaxAge = app.SessionMaxAge
		session, err = app.SessionStore.New(r, app.SessionCookieKey)
		session.Save(r, w) // renews the session
		if err != nil {
			log.Println("recreateSession: " + err.Error() + " (after '" + origError.Error() + "')")
		}

		//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		log.Println("GOT Session: " + session.ID)
	}

	return session
}

// SetUpMuxHandlers sets up the router
func SetUpMuxHandlers(isMainExec bool) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api", APIAccountsHandler)
	//router.HandleFunc("/callback", CallbackHandler)
	router.HandleFunc("/games/{gameid}", GameHandler).Name("games")
	router.HandleFunc("/account/", AccountHandler)
	router.HandleFunc("/wsInit", ServeWebSocket)
	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			VersionHandler(w, r)
		}
	})

	templateDir := os.Getenv("TEMPLATE_DIR")
	if !isMainExec {
		templateDir = "templates"
	}
	renderer = render.New(render.Options{
		Layout:        "content",
		IsDevelopment: !app.ProdSession,
		Funcs:         []template.FuncMap{FuncMap},
		Directory:     templateDir,
	})

	return router
}

// FuncMap for common functions in templates
var FuncMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"curr_time": func() int64 {
		return time.Now().Unix()
	},
	"to_string": func(value interface{}) string {
		switch v := value.(type) {
		case string:
			return v
		case int:
			return strconv.Itoa(v)
			// Add whatever other types you need
		case bson.ObjectId:
			return v.Hex()
		default:
			return ""
		}
	},
}
