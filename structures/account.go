package structures

import (
	"html/template"
	"time"

	uuid "github.com/satori/go.uuid"
)

var onlineAccounts []*Account

// Account object representing a user account
type Account struct {
	ID            string // unique
	Auth0ID       string // if nil, not a registered user: no persistence
	Commander     string
	Games         []*Game // max 3
	CurrentGameID string
	LastLogout    time.Time
	LastLogin     time.Time
	ClickableID   template.JS
}

// NewAccount and session
func NewAccount(username string) *Account {
	id := uuid.NewV5(uuid.NamespaceOID, username+time.Now().String()).String()
	account := &Account{
		ID:            id,
		Auth0ID:       "",
		CurrentGameID: NewGameUUID,
		Commander:     username,
		LastLogin:     time.Now(),
		LastLogout:    time.Now(),
		ClickableID:   template.JS(id),
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
func GetAccount(accountID string) *Account {
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
func (account Account) EndSession() {
	account.LastLogout = time.Now()
	removeAccountFromActiveSessions(account)

	if account.Auth0ID != "" {
		// persist games owned by this account!
	}
}

// DeleteAccount and all games from sessions AND persistence store
func (account Account) DeleteAccount() {
	removeAccountFromActiveSessions(account)

	if account.Auth0ID != "" {
		// delete Auth0 account
		// delete MongoDB games owned by this account
	}
}
