package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	isMainExec := true // main() is only called from full application start
	app.Init(isMainExec)

	r := mux.NewRouter()

	// Serve static assets directly
	//r.PathPrefix("/wsInit").HandlerFunc(routes.ServeWebSocket).Host("{subdomain:api}.localhost")
	r.PathPrefix("/wsInit").HandlerFunc(routes.ServeWebSocket)
	r.PathPrefix("/fonts").Handler(http.FileServer(http.Dir("public")))
	r.PathPrefix("/css").Handler(http.FileServer(http.Dir("public")))
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir("public")))
	r.PathPrefix("/statics").Handler(http.FileServer(http.Dir("public")))

	// Catch-all: Serve our JavaScript application's entry-point (index.html).
	r.PathPrefix("/").HandlerFunc(IndexHandler("public/index.html"))

	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    app.HTTPPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}
