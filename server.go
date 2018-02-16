package main

import (
	"net/http"
	"strings"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/routes"
	"github.com/urfave/negroni"
)

func checkAPICall(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// verify only API calls coming in on api.battlefleet.online
	if strings.HasPrefix(r.Host, "api") {
		if !strings.HasPrefix(r.URL.RequestURI(), "/version") &&
			!strings.HasPrefix(r.URL.RequestURI(), "/api") {
			http.Redirect(w, r, "https://battlefleet.online/", http.StatusMovedPermanently)
		}
	}

	next(w, r)
}

func main() {
	isUnitTest := false // main() is only called from full application start
	app.Init(isUnitTest)

	router := routes.SetUpMuxHandlers(isUnitTest)
	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(checkAPICall))
	n.UseHandler(router)

	http.ListenAndServe(app.HTTPPort, n)
}
