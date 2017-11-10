package routes

import (
	"html/template"
	"net/http"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/structures"
	"github.com/unrolled/render"
)

// HomeHandler for handling index page requests
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.SessionStore.Get(r, app.SessionCookieKey)

	if session.Values[app.CmdrNameKey] == nil {
		session.Values[app.CmdrNameKey] = app.DefaultCmdrName
	}

	if session.Values[app.GameUUIDKey] == nil {
		session.Values[app.GameUUIDKey] = structures.NewGameUUID
	}

	account := structures.NewAccount(session.Values[app.CmdrNameKey].(string))
	// retrieve account
	if session.Values[app.AccountKey] != nil {
		account = session.Values[app.AccountKey].(*structures.Account)
	}

	render := render.New(render.Options{
		Layout:        "content",
		IsDevelopment: !app.ProdSession,
		Funcs:         []template.FuncMap{FuncMap},
	})

	render.HTML(w, http.StatusOK, "home",
		map[string]interface{}{"Account": account, "Data": app.AuthZeroData})
}
