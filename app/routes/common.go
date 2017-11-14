package routes

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/danackerson/battlefleet/app"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
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

// SetUpMuxHandlers sets up the router
func SetUpMuxHandlers(isTest bool) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/callback", CallbackHandler)
	router.HandleFunc("/games/{gameid}", GameHandler).Name("games")
	router.HandleFunc("/account/", AccountHandler)
	router.HandleFunc("/wsInit", ServeWebSocket)
	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			VersionHandler(w, r)
		}
	})

	templateDir := os.Getenv("TEMPLATE_DIR")
	if isTest {
		templateDir = "templates"
	}
	renderer := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: !app.ProdSession,
		Funcs:         []template.FuncMap{FuncMap},
		Directory:     templateDir,
	})

	spew.Dump(renderer)

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
