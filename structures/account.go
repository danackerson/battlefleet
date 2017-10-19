package structures

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

var onlineAccounts []*Account

// Account object representing a user account
type Account struct {
	ID         string // unique
	Auth0ID    string // if nil, not a registered user: no persistence
	Commander  string
	Games      []Game
	LastLogout time.Time
	LastLogin  time.Time
}

// NewAccount and session
func NewAccount(username string) *Account {
	account := Account{
		ID:        uuid.NewV5(uuid.NamespaceOID, username).String(),
		Commander: username,
		LastLogin: time.Now(),
	}

	onlineAccounts = append(onlineAccounts, &account)

	return &account
}

// AddGame to an account
func (account *Account) AddGame(game *Game) {
	// TODO: max 3 games per account!
	account.Games = append(account.Games, *game)
}

// GetAccount finds the account among the online sessions
func GetAccount(accountID string) *Account {
	var account *Account
	for _, accountToCheck := range onlineAccounts {
		if accountToCheck.ID == accountID {
			account = accountToCheck
			break
		}
	}

	if account == nil {
		// TODO: go search in MongoDB!
		// & add to online sessions
	}

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

	// TODO
	if account.Auth0ID != "" {
		// persist games owned by this account!
	}
}

// DeleteAccount and all games from sessions AND persistence store
func (account Account) DeleteAccount() {
	removeAccountFromActiveSessions(account)

	// TODO
	if account.Auth0ID != "" {
		// delete Auth0 account
		// delete MongoDB games owned by this account
	}
}
