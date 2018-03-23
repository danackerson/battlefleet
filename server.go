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
		// and API calls to battlefleet are redirected to api.battlefleet.online
	} else if strings.HasPrefix(r.URL.RequestURI(), "/version") ||
		strings.HasPrefix(r.URL.RequestURI(), "/api") {
		http.Redirect(w, r, "https://api.battlefleet.online/", http.StatusMovedPermanently)
	}

	next(w, r)
}

func main() {
	n := negroni.New()
	n.Use(negroni.HandlerFunc(checkAPICall))

	logger := negroni.NewLogger()
	logger.SetFormat("{{.StartTime}} {{.Request.RemoteAddr}} {{.Method}} {{.Path}} @ {{.Hostname}} - [{{.Status}} {{.Duration}}] - \"{{.Request.UserAgent}}\"\n")
	n.Use(logger)

	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewStatic(http.Dir("public")))

	isMainExec := true // main() is only called from full application start
	app.Init(isMainExec)
	n.UseHandler(routes.SetUpMuxHandlers(isMainExec))

	http.ListenAndServe(app.HTTPPort, n)
}
