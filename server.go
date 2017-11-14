package main

import (
	"net/http"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/routes"
	"github.com/urfave/negroni"
)

func main() {
	app.Init()

	router := routes.SetUpMuxHandlers(false)
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(app.HTTPPort, n)
}
