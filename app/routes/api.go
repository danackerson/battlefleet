package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danackerson/battlefleet/app/structures"
	"gopkg.in/mgo.v2/bson"
)

type onlineAccounts struct {
	AccountID     bson.ObjectId `json:"id"`
	CommanderName string        `json:"name"`
}

// APIAccountsHandler entrypoint to Battlefleet API
func APIAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts := structures.GetOnlineAccounts()

	var jsonAccounts []onlineAccounts
	for _, account := range accounts {
		var tmpAccount onlineAccounts
		tmpAccount.AccountID = account.ID
		tmpAccount.CommanderName = account.Commander
		jsonAccounts = append(jsonAccounts, tmpAccount)
	}
	jsonResponse := map[string]interface{}{"online": len(accounts), "accounts": jsonAccounts}

	data, _ := json.Marshal(jsonResponse)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
