package main

import (
	"net/http"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/routes"
	"github.com/urfave/negroni"
)

func main() {
	isUnitTest := false // main() is only called from full application start
	app.Init(isUnitTest)

	router := routes.SetUpMuxHandlers(isUnitTest)
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(app.HTTPPort, n)
}
