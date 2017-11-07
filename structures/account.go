package structures

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"golang.org/x/oauth2"
)

var onlineAccounts []*Account

// Account object representing a user account
type Account struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Auth0Token    *oauth2.Token
	Auth0Profile  map[string]interface{}
	Commander     string
	Games         []*Game
	CurrentGameID string
	LastLogout    time.Time
	LastLogin     time.Time
}

// NewAccount and session
func NewAccount(username string) *Account {
	account := &Account{
		ID:            bson.NewObjectId(),
		Auth0Token:    nil,
		CurrentGameID: NewGameUUID,
		Commander:     username,
		LastLogin:     time.Now(),
		LastLogout:    time.Now(),
	}

	onlineAccounts = append(onlineAccounts, account)

	return account
}

// DeleteGame from the account
func (account *Account) DeleteGame(gameID string) {
	for i, game := range account.Games {
		if game.ID == gameID {
			//copy(account.Games[i:], account.Games[i+1:])
			account.Games[i] = account.Games[len(account.Games)-1]
			account.Games = account.Games[:len(account.Games)-1]
			if account.CurrentGameID == gameID {
				account.CurrentGameID = NewGameUUID
			}
			break
		}
	}
}

// OwnsGame checks to see if the provided gameUUID is owned by the asking account
func (account *Account) OwnsGame(gameID string) bool {
	owns := false
	for _, game := range account.Games {
		if game.ID == gameID {
			owns = true
			break
		}
	}

	return owns
}

// AddGame to an account
func (account *Account) AddGame(game *Game) {
	account.CurrentGameID = game.ID
	account.Games = append(account.Games, game)
}

// GetAccount finds the account among the online sessions
func GetAccount(accountID bson.ObjectId) *Account {
	account := &Account{}
	for _, accountToCheck := range onlineAccounts {
		if accountToCheck.ID == accountID {
			account = accountToCheck
			break
		}
	}

	// if account == nil => TODO: go search in MongoDB?
	return account
}

func removeAccountFromActiveSessions(account Account) {
	for index, accountToCheck := range onlineAccounts {
		if accountToCheck.ID == account.ID {
			// set current node to copy of last node
			onlineAccounts[index] = onlineAccounts[len(onlineAccounts)-1]

			// cut off last node and set as new list
			onlineAccounts = onlineAccounts[:len(onlineAccounts)-1]
			break
		}
	}
}

// EndSession by removing account
func (account Account) EndSession(db *mgo.Session) {
	removeAccountFromActiveSessions(account)

	if account.Auth0Token != nil {
		mongoSession := db.Copy()
		defer mongoSession.Close()
		c := mongoSession.DB("fleetbattle").C("sessions")
		account.LastLogout = time.Now()
		_, err := c.UpsertId(account.ID, account)
		if err != nil {
			panic(err)
		}
	}
}

// DeleteAccount and all games from sessions AND persistence store
func (account Account) DeleteAccount(db *mgo.Session) {
	removeAccountFromActiveSessions(account)

	if account.Auth0Token != nil {
		mongoSession := db.Copy()
		defer mongoSession.Close()
		c := mongoSession.DB("fleetbattle").C("sessions")
		err := c.RemoveId(account.ID)
		if err != nil {
			panic(err)
		}
	}
}
