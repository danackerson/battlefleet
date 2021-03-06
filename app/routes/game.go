package routes

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/structures"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"
)

// GameHandler for handling game requests
func GameHandler(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	parseErr := r.ParseForm()
	if parseErr != nil {
		log.Println(parseErr)
	}

	session := RetrieveSession(w, r)
	account := getAccount(r, session)
	if account != nil {
		gameUUID := requestParams["gameid"]
		action, ok := r.URL.Query()["action"]
		if ok && len(action) > 0 && action[0] == "delete" {
			account.DeleteGame(gameUUID)
			if session.Values[app.GameUUIDKey] == gameUUID {
				session.Values[app.GameUUIDKey] = ""
			}
			// remember, the session *is* the persistence store
			// a new request will fetch the account from the session on disk
			// so deleting a game is not really Deleted until the session is saved!
			if e := session.Save(r, w); e != nil {
				t, _ := template.New("errorPage").Parse(errorPage)
				t.Execute(w, "saveSession: "+e.Error())
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/account/", http.StatusTemporaryRedirect)
			return
		}

		redirected := setupGame(r, w, session, account, gameUUID)
		if !redirected {
			renderer.HTML(w, http.StatusOK, "game",
				map[string]interface{}{
					"Account": map[string]interface{}{
						"ID":           account.ID,
						"Commander":    account.Commander,
						"Auth0Profile": account.Auth0Profile,
						"Auth0Token":   account.Auth0Token},
					"AuthData": app.AuthZeroData,
					"DevEnv":   !app.ProdSession,
					"Version":  GetVersionInfo(),
					"GridSize": structures.GridSize,
				})
		}
	} else {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "New accounts require a Commander name and '"+app.DefaultCmdrName+"' is not allowed.")
		http.Redirect(w, r, "/", http.StatusPreconditionRequired)
	}
}

func setupGame(r *http.Request, w http.ResponseWriter,
	session *sessions.Session, account *structures.Account, gameUUID string) bool {
	redirected := false

	// they come in without a cookie or request a gameID that doesn't match their own
	if gameUUID != structures.NewGameUUID {
		if account.OwnsGame(gameUUID) {
			account.CurrentGameID = gameUUID
			session.Values[app.GameUUIDKey] = gameUUID
			if e := session.Save(r, w); e != nil {
				t, _ := template.New("errorPage").Parse(errorPage)
				t.Execute(w, "saveSession1: "+e.Error())
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				redirected = true
				return redirected
			}

			structures.AddOnlineAccount(account)
		} else {
			t, _ := template.New("errorPage").Parse(errorPage)
			errorString := "You neither own Game ID:<span style='color:orange;'>" + gameUUID + "</span> nor have you been invited to join."
			t.Execute(w, template.JS(errorString))
			http.Redirect(w, r, "/", http.StatusPreconditionRequired)
			redirected = true
			return redirected
		}
	} else if gameUUID == structures.NewGameUUID {
		sessionIDHash := session.ID + time.Now().String()
		gameUUID = uuid.NewV5(uuid.NamespaceOID, sessionIDHash).String()
		newGame := structures.NewGame(gameUUID, account.ID)
		account.AddGame(newGame)
		session.Values[app.AccountKey] = account
		session.Values[app.GameUUIDKey] = gameUUID
		if e := session.Save(r, w); e != nil {
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, "saveSession2: "+e.Error())
			http.Redirect(w, r, "/", http.StatusPreconditionRequired)
			redirected = true
			return redirected
		}

		structures.AddOnlineAccount(account)

		http.Redirect(w, r, "/games/"+gameUUID, http.StatusMovedPermanently)
		redirected = true
		return redirected
	}

	return redirected
}
