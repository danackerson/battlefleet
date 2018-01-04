package routes

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/structures"
	"github.com/gorilla/sessions"
)

// https://github.com/riscie/websocket-tic-tac-toe/ <= cool ideas

// AccountHandler for handling account requests
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	var account *structures.Account
	session := RetrieveSession(w, r)
	if session.Values[app.AccountKey] != nil {
		account = session.Values[app.AccountKey].(*structures.Account)

		requestParams := r.URL.Query()
		if len(requestParams) > 0 {
			if requestParams["action"][0] == "delete" {
				account.DeleteAccount(app.DB)
				session.Options.MaxAge = -1
				if e := session.Save(r, w); e != nil {
					t, _ := template.New("errorPage").Parse(errorPage)
					t.Execute(w, "saveSession1: "+e.Error())
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}
				// Session Flash msg "Account deleted"
				http.Redirect(w, r, "/?account=deleted", http.StatusTemporaryRedirect)
				return
			} else if requestParams["action"][0] == "logout" {
				account.EndSession(app.DB)
				session.Options.MaxAge = -1
				if e := session.Save(r, w); e != nil {
					t, _ := template.New("errorPage").Parse(errorPage)
					t.Execute(w, "saveSession2: "+e.Error())
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}

				returnHost := strings.Replace(r.Host, ":", "%3A", 1)
				returnTo := "?returnTo=" + app.URIScheme + "%3A%2F%2F" + returnHost + "&client_id=" + app.AuthZeroData.Auth0ClientID
				http.Redirect(w, r, "https://battlefleet.eu.auth0.com/v2/logout"+returnTo, http.StatusTemporaryRedirect)
				return
			}

		}

		renderer.HTML(w, http.StatusOK, "account",
			map[string]interface{}{"Account": account, "AuthData": app.AuthZeroData, "DevEnv": !app.ProdSession})
	} else {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "You are currently not logged in!")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
}

func getAccount(r *http.Request, session *sessions.Session) *structures.Account {
	var account *structures.Account

	if session.Values[app.AccountKey] == nil {
		if r.FormValue("cmdrName") == "" || r.FormValue("cmdrName") == app.DefaultCmdrName {
			// new accounts require a CommanderName and 'stranger!' is reserved ;)
			return nil
		}

		account = structures.NewAccount(r.FormValue("cmdrName"))
		session.Values[app.AccountKey] = account
	} else {
		account = session.Values[app.AccountKey].(*structures.Account)
	}

	return account
}
