// ********************************************************************
// * message.go                                                       *
// *                                                                  *
// * 2020-03-15 First Version, JR                                     *
// *                                                                  *
// * This package serves as container in memory for json strings,     *
// * the message struct contains all the fields required to keep,     *
// * and parse json input and output messages.                        *                                                                *
// *                                                                  *
// * Usage:                                                           *
// * msg := message.New(Account, Transaction, Violations)             *
// * msg.Type()                                                       *
// * msg.AddViolation(Code)                                           *
// ********************************************************************

package message

import (
	"authorizer/account"
)

// Type of messages.
const (
	Account     = "account"
	Transaction = "transaction"
)

// Message - Represents json output line while is in memory.
type Message struct {
	Account     *AccountMessage      `json:"account"`
	Transaction *account.Transaction `json:"transaction,omitempty"`
	Violations  []string             `json:"violations"`
}

// AccountMessage - represents account fields gotten from json input.
type AccountMessage struct {
	Active bool `json:"active-card"`
	Limit  int  `json:"available-limit"`
}

// New - Returns a new empty ready to be constructed with executer process.
func New(a *AccountMessage, t *account.Transaction, v []string) *Message {
	return &Message{a, t, v}
}

// Type returns a string to identify the operation type
// "account" or "transaction"
func (msg *Message) Type() string {
	if msg.Account != nil {
		return Account
	}

	return Transaction
}

// AddViolation - adds a new violation message to the array, accordignly the,
// account violation code.
func (msg *Message) AddViolation(code int) {
	if account.Violations[code] != "" {
		msg.Violations = append(msg.Violations, account.Violations[code])
	}
}
