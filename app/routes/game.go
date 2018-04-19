package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/structures"
	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"
)

// GameHandler for handling game requests
func GameHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSOptions(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var fields map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&fields)
	if err != nil {
		sendError(w, 502, "failed to decode fields: "+err.Error())
		return
	}

	session := RetrieveSession(w, r)
	account := getAccount(session, fields["cmdrName"])
	if account != nil {
		gameUUID := fields["gameID"]
		if gameUUID == "" {
			gameUUID = account.CurrentGameID
		}
		/*action, ok := r.URL.Query()["action"]
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
		}*/

		setupGame(r, w, session, account, gameUUID)

		var accountJSON accountType
		accountJSON.ID = account.ID.Hex()
		accountJSON.Commander = account.Commander

		var gameJSON gameType
		gameJSON.ID = account.CurrentGameID
		gameJSON.Account = accountJSON
		gameJSON.GridSize = structures.GridSize
		gameJSON.Version = GetVersionInfo()

		json.NewEncoder(w).Encode(gameJSON)
	} else {
		sendError(w, 412, "Session no longer exists. Please login or start over.")
		return
	}
}

func setupGame(r *http.Request, w http.ResponseWriter,
	session *sessions.Session, account *structures.Account, gameUUID string) {

	// they come in without a cookie or request a gameID that doesn't match their own
	if gameUUID != structures.NewGameUUID {
		if gameUUID == "" && account.CurrentGameID == "" {
			gameUUID = structures.NewGameUUID
		} else if gameUUID == "" {
			gameUUID = account.CurrentGameID
		}

		if account.OwnsGame(gameUUID) {
			account.CurrentGameID = gameUUID
			session.Values[app.GameUUIDKey] = gameUUID
			if e := session.Save(r, w); e != nil {
				sendError(w, 401, "Unable to access your game: "+gameUUID)
				return
			}

			structures.AddOnlineAccount(account)
		} else {
			sendError(w, 401, "You neither own Game ID:<span style='color:orange;'>"+gameUUID+"</span> nor have you been invited to join.")
			return
		}
	}

	if gameUUID == structures.NewGameUUID {
		sessionIDHash := session.ID + time.Now().String()
		gameUUID = uuid.NewV5(uuid.NamespaceOID, sessionIDHash).String()
		newGame := structures.NewGame(gameUUID, account.ID)
		account.AddGame(newGame)
		session.Values[app.AccountKey] = account
		session.Values[app.GameUUIDKey] = gameUUID
		if e := session.Save(r, w); e != nil {
			sendError(w, 501, "There was a problem creating your game.")
			return
		}

		structures.AddOnlineAccount(account)
	}
}
