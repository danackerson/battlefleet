package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/structures"
	"github.com/gorilla/sessions"
)

// https://github.com/riscie/websocket-tic-tac-toe/ <= cool ideas

// AccountHandler for handling account requests
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSOptions(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var fields map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&fields)
	if err != nil {
		sendError(w, 502, "failed to decode fields: "+err.Error())
		return
	}

	var account *structures.Account
	session := RetrieveSession(w, r)
	if session.Values[app.AccountKey] != nil {
		account = session.Values[app.AccountKey].(*structures.Account)

		switch r.URL.Path {
		case "/logout":
			account.EndSession(app.DB)
			session.Options.MaxAge = -1
			if e := session.Save(r, w); e != nil {
				sendError(w, 502, e.Error())
				return
			}

			sendSuccess(w, 200, "Thanks for playing")
			return
		case "/updateAccount":
			account.Commander = fields["CmdrName"].(string)
			if e := session.Save(r, w); e != nil {
				sendError(w, 502, e.Error())
				return
			}
			sendSuccess(w, 200, "Successfully updated account")
			return
		case "/deleteAccount":
			account.DeleteAccount(app.DB)
			session.Options.MaxAge = -1
			if e := session.Save(r, w); e != nil {
				t, _ := template.New("errorPage").Parse(errorPage)
				t.Execute(w, "saveSession1: "+e.Error())
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				return
			}
		}
	} else if r.URL.Path == "/login" {
		log.Printf("LOOKing for account by Auth0 ID: %s\n", fields["Auth0Token"].(string))
		account = getAccountByAuth0(fields["Auth0Token"].(string))
		setupGame(r, w, session, account, account.GetGame().ID)
		var accountJSON accountType
		accountJSON.ID = account.ID.Hex()
		accountJSON.Commander = account.Commander

		var gameJSON gameType
		gameJSON.ID = account.CurrentGameID
		gameJSON.Account = accountJSON
		gameJSON.GridSize = structures.GridSize

		json.NewEncoder(w).Encode(gameJSON)
		// store account to session and save session
		// look at game handler how to prep objects for response
		// send response
		return
	}
}

func getAccountByAuth0(auth0Token string) *structures.Account {
	var fetchedAccount *structures.Account
	fetchedAccount = structures.FindAccountByAuth0ProfileSubToken(app.DB, auth0Token)
	return fetchedAccount
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
