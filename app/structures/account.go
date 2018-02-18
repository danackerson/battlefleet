package structures

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"golang.org/x/oauth2"
)

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

// GetGame returns the current active game from chosen account
func (account *Account) GetGame() *Game {
	currentGame := &Game{}

	for _, game := range account.Games {
		if game.ID == account.CurrentGameID {
			currentGame = game
			break
		}
	}

	return currentGame
}

func removeAccountFromActiveSessions(account Account) {
	RemoveOnlineAccount(account.ID)
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
