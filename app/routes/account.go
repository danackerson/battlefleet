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
	log.Printf("FIELDS: %v %s", fields, fields["CmdrName"])
	if err != nil {
		sendError(w, 502, "failed to decode fields: "+err.Error())
		return
	}

	var account *structures.Account
	session := RetrieveSession(w, r)
	if session.Values[app.AccountKey] != nil {
		account = session.Values[app.AccountKey].(*structures.Account)

		switch r.URL.Path {
		case "/login":
			// todo
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
