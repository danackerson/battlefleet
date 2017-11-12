package routes

import (
	"net/http"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/structures"
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

	renderer.HTML(w, http.StatusOK, "home",
		map[string]interface{}{"Account": account, "Data": app.AuthZeroData})
}
