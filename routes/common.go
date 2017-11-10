package routes

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const ErrorPage = `
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
func SetUpMuxHandlers() *mux.Router {
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
}
