package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	isMainExec := true // main() is only called from full application start
	app.Init(isMainExec)

	r := mux.NewRouter()

	// Serve static assets directly
	r.PathPrefix("/wsInit").HandlerFunc(routes.ServeWebSocket)
	r.PathPrefix("/login").HandlerFunc(routes.AccountHandler).Methods("POST", "OPTIONS")
	r.PathPrefix("/logout").HandlerFunc(routes.AccountHandler).Methods("POST", "OPTIONS")
	r.PathPrefix("/updateAccount").HandlerFunc(routes.AccountHandler).Methods("POST", "OPTIONS")

	r.PathPrefix("/newGame").HandlerFunc(routes.GameHandler).Methods("POST", "OPTIONS")

	r.PathPrefix("/version").HandlerFunc(routes.VersionHandler).Methods("POST", "OPTIONS")

	r.PathPrefix("/fonts").Handler(http.FileServer(http.Dir("public")))
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir("public")))
	r.PathPrefix("/statics").Handler(http.FileServer(http.Dir("public")))
	r.PathPrefix(`/{0.(.*).css}`).Handler(http.FileServer(http.Dir("public")))

	// Catch-all: Serve our JavaScript application's entry-point (index.html).
	r.PathPrefix("/").HandlerFunc(IndexHandler("public/index.html"))

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         app.HTTPPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// IndexHandler serves up the index.template.html which bootstraps the VueJS
func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}
