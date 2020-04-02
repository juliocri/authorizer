// ********************************************************************
// * executer.go                                                      *
// *                                                                  *
// * 2020-03-16 First Version, JR                                     *
// *                                                                  *
// * Package responsible of build an output json line                 *
// * based in another input json message.                             *
// *                                                                  *
// * Usage:                                                           *
// * e:= executer.Init()                                              *
// * e.Exec(string)                                                   *
// ********************************************************************

package executer

import (
	"authorizer/account"
	"authorizer/executer/message"
	"encoding/json"
	"strings"
)

// Executer - Holds the reference to the working account.
type Executer struct {
	account *account.Account
}

// Init Returns a new Executer.
func Init() *Executer {
	return &Executer{}
}

// Exec - Returns a json line string build based in a json operation line.
func (exe *Executer) Exec(op string) string {
	// Transform json string to Message struct.
	msg := message.New(nil, nil, []string{})
	json.Unmarshal([]byte(op), msg)

	// Check operation type is "account".
	if msg.Type() == message.Account {
		exe.initAccount(msg)
	}

	// If is a "transaction" message then:
	if msg.Type() == message.Transaction {
		exe.processTransaction(msg)
	}

	// Clean temporary data from our output structure message,
	// in order to be converted to json.
	msg.Transaction = nil
	//msg.Account = &message.AccountMessage{}
	// TODO: Check with nubank how we going to handle json message for "account",
	// When account is not initialized.
	// Currently assuming `{"account": {}, "violations":[]}`
	if exe.account != nil {
		msg.Account = &message.AccountMessage{
			Active: exe.account.Active(),
			Limit:  exe.account.Limit(),
		}
	}

	// Converting to json.
	output, _ := json.Marshal(msg)
	// Prettify with no null refs and spaces.
	result := strings.Replace(string(output), "null", "{}", -1)
	result = strings.Replace(result, ":", ": ", -1)
	result = strings.Replace(result, ",", ", ", -1)
	return result
}

// initAccount - Create a new account and add the reference to the executioner,
// additionaly adds violations if there was found.
func (exe *Executer) initAccount(msg *message.Message) {
	// Try init account.
	acn, v := exe.account.Init(msg.Account.Active, msg.Account.Limit)
	// If violation found.
	if v != -1 {
		// Then add to message.
		msg.AddViolation(v)
	} else {
		exe.account = acn
	}
}

// processTransaction - Execute a transaction and add violations if there was,
// found in the process.
func (exe *Executer) processTransaction(msg *message.Message) {
	// check if account is initialized.
	violations := exe.account.ApplyTransaction(msg.Transaction)
	if len(violations) > 0 {
		for _, v := range violations {
			msg.AddViolation(v)
		}
	}
}
