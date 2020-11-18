// ********************************************************************
// * message_test.go                                                  *
// *                                                                  *
// * 2020-03-15 First Version, JR                                     *
// *                                                                  *
// * This file contains all unit test related with message operations.*
// *                                                                  *
// * Usage: go test -v ./executer/message                             *
// ********************************************************************

package message

import (
	"authorizer/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Max int code value allowed to violations references.
const maxValidCode = 6

// Test Messages to cover our unit testing.
var tmsgs = map[string]*Message{

	"Sample": {
		Violations: []string{},
	},

	"Account": {
		Account: &AccountMessage{
			Active: true,
			Limit:  100,
		},
	},

	"Transaction": {
		Transaction: &account.Transaction{
			Merchant: "Fulanito",
			Amount:   100,
			Time:     "2019-02-13T10:00:00.000Z",
		},
	},
}

// Test message creation.
func TestNewMessage(t *testing.T) {
	msg := New(nil, nil, []string{})
	assert := assert.New(t)
	assert.Equal(
		[]string{},
		msg.Violations,
		"Expected a new message with no violations.",
	)
}

// Test if a message is of type "account".
func TestAccountMessage(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		Account,
		tmsgs["Account"].Type(),
		"Expected an account type message.",
	)
}

// Test if a message is of type "transaction".
func TestTransactionMessage(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		Transaction,
		tmsgs["Transaction"].Type(),
		"Expected a transaction type message.",
	)
}

// Test if not valid violation code is added to the message's list.
func TestAddNotValidViolationCode(t *testing.T) {
	tmsgs["Sample"].AddViolation(maxValidCode + 1)
	assert := assert.New(t)
	assert.Equal(
		[]string{},
		tmsgs["Sample"].Violations,
		"Expected array with no violations.",
	)
}

// Test if a violation code is added successfully to the message's list.
func TestAddValidViolationCode(t *testing.T) {
	tmsgs["Sample"].AddViolation(maxValidCode)
	assert := assert.New(t)
	assert.Equal(
		[]string{account.Violations[maxValidCode]},
		tmsgs["Sample"].Violations,
		"Expected array with violation associated string message.",
	)
}
