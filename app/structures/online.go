package structures

import (
	"gopkg.in/mgo.v2/bson"
)

var onlineAccounts []*Account

func AddOnlineAccount(account *Account) {
	onlineAccounts = append(onlineAccounts, account)
}

// GetOnlineAccounts returns all onlineAccounts
func GetOnlineAccounts() []*Account {
	return onlineAccounts
}

// GetAccount finds the account among the online sessions
func GetAccount(accountID bson.ObjectId) *Account {
	for _, accountToCheck := range onlineAccounts {
		if accountToCheck.ID == accountID {
			return accountToCheck
		}
	}

	// if account == nil => TODO: go search in MongoDB?
	return nil
}

// RemoveOnlineAccount removes the account from active online state
func RemoveOnlineAccount(accountID bson.ObjectId) {
	for index, accountToCheck := range onlineAccounts {
		if accountToCheck.ID == accountID {
			// set current node to copy of last node
			onlineAccounts[index] = onlineAccounts[len(onlineAccounts)-1]

			// cut off last node and set as new list
			onlineAccounts = onlineAccounts[:len(onlineAccounts)-1]
			break
		}
	}
}
