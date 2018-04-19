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
	/*setupCORSOptions(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var fields map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&fields)
	log.Printf("FIELDS: %v", fields)
	if err != nil {
		sendError(w, 502, "failed to decode fields: "+err.Error())
		return
	}*/

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
			map[string]interface{}{
				"Account": account, "AuthData": app.AuthZeroData,
				"DevEnv": !app.ProdSession, "Version": GetVersionInfo(),
			})
	} else {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "You are currently not logged in!")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
}

func getAccount(session *sessions.Session, cmdrName string) *structures.Account {
	var account *structures.Account

	/*requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))*/

	if session.Values[app.AccountKey] == nil {
		if cmdrName == "" || cmdrName == app.DefaultCmdrName {
			// new accounts require a CommanderName and 'stranger!' is reserved ;)
			return nil
		}

		account = structures.NewAccount(cmdrName)
		session.Values[app.AccountKey] = account
	} else {
		account = session.Values[app.AccountKey].(*structures.Account)
	}

	return account
}
