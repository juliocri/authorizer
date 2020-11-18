// ********************************************************************
// * account_test.go                                                  *
// *                                                                  *
// * 2020-03-15 First Version, JR                                     *
// * 2020-03-18 Adds multiple violation scnario, JR                   *
// *                                                                  *
// * This file contains all unit testing related with executer.       *                                                    *
// *                                                                  *
// * Usage: go test -v ./executer                                     *
// ********************************************************************

package executer

import (
	"authorizer/account"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test cases scenarios to cover our executer unit testing.
var tinputs = map[string]map[string][]string{

	"notInitialized": {
		"in": []string{
			`{"transaction": {"merchant": "Burger Queen", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
		},
		"out": []string{
			fmt.Sprintf(`{"account": {}, "violations": ["%v"]}`, account.Violations[0]),
		},
	},

	"NotActive": {
		"in": []string{
			`{"account": {"active-card": false, "available-limit": 100}}`,
			`{"transaction": {"merchant": "Burger Queen", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
		},
		"out": []string{
			`{"account": {"active-card": false, "available-limit": 100}, "violations": []}`,
			fmt.Sprintf(`{"account": {"active-card": false, "available-limit": 100}, "violations": ["%v"]}`, account.Violations[1]),
		},
	},

	"AlreadyInitialized": {
		"in": []string{
			`{"account": {"active-card": true, "available-limit": 100}}`,
			`{"account": {"active-card": true, "available-limit": 350}}`,
		},
		"out": []string{
			`{"account": {"active-card": true, "available-limit": 100}, "violations": []}`,
			fmt.Sprintf(`{"account": {"active-card": true, "available-limit": 100}, "violations": ["%v"]}`, account.Violations[2]),
		},
	},

	"blockedMerchant": {
		"in": []string{
			`{"account": {"active-card": true, "available-limit": 100}}`,
			`{"transaction": {"merchant": "Burger King", "amount": 100, "time": "2019-02-13T10:00:00.000Z"}}`,
		},
		"out": []string{
			`{"account": {"active-card": true, "available-limit": 100}, "violations": []}`,
			fmt.Sprintf(`{"account": {"active-card": true, "available-limit": 100}, "violations": ["%v"]}`, account.Violations[6]),
		},
	},

	"insufficientLimit": {
		"in": []string{
			`{"account": {"active-card": true, "available-limit": 100}}`,
			`{"transaction": {"merchant": "Burger Queen", "amount": 200, "time": "2019-02-13T10:00:00.000Z"}}`,
		},
		"out": []string{
			`{"account": {"active-card": true, "available-limit": 100}, "violations": []}`,
			fmt.Sprintf(`{"account": {"active-card": true, "available-limit": 100}, "violations": ["%v"]}`, account.Violations[3]),
		},
	},

	"DoubledTransaction": {
		"in": []string{
			`{"account": {"active-card": true, "available-limit": 100}}`,
			`{"transaction": {"merchant": "Burger Queen", "amount": 50, "time": "2019-02-13T10:00:00.000Z"}}`,
			`{"transaction": {"merchant": "Burger Queen", "amount": 50, "time": "2019-02-13T10:01:00.000Z"}}`,
		},
		"out": []string{
			`{"account": {"active-card": true, "available-limit": 100}, "violations": []}`,
			`{"account": {"active-card": true, "available-limit": 50}, "violations": []}`,
			fmt.Sprintf(`{"account": {"active-card": true, "available-limit": 50}, "violations": ["%v"]}`, account.Violations[4]),
		},
	},

	"HighFrequency": {
		"in": []string{
			`{"account": {"active-card": true, "available-limit": 100}}`,
			`{"transaction": {"merchant": "Burger Queen1", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
			`{"transaction": {"merchant": "Burger Queen2", "amount": 30, "time": "2019-02-13T10:00:30.000Z"}}`,
			`{"transaction": {"merchant": "Burger Queen3", "amount": 10, "time": "2019-02-13T10:01:00.000Z"}}`,
			`{"transaction": {"merchant": "Burger Queen", "amount": 10, "time": "2019-02-13T10:02:00.000Z"}}`,
		},
		"out": []string{
			`{"account": {"active-card": true, "available-limit": 100}, "violations": []}`,
			`{"account": {"active-card": true, "available-limit": 80}, "violations": []}`,
			`{"account": {"active-card": true, "available-limit": 50}, "violations": []}`,
			`{"account": {"active-card": true, "available-limit": 40}, "violations": []}`,
			fmt.Sprintf(`{"account": {"active-card": true, "available-limit": 40}, "violations": ["%v"]}`, account.Violations[5]),
		},
	},

	"HighFrequencyDoubled": {
		"in": []string{
			`{"account": {"active-card": true, "available-limit": 100}}`,
			`{"transaction": {"merchant": "Burger Queen1", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
			`{"transaction": {"merchant": "Burger Queen2", "amount": 30, "time": "2019-02-13T10:00:30.000Z"}}`,
			`{"transaction": {"merchant": "Burger Queen", "amount": 10, "time": "2019-02-13T10:01:00.000Z"}}`,
			`{"transaction": {"merchant": "Burger Queen", "amount": 10, "time": "2019-02-13T10:01:00.000Z"}}`,
		},
		"out": []string{
			`{"account": {"active-card": true, "available-limit": 100}, "violations": []}`,
			`{"account": {"active-card": true, "available-limit": 80}, "violations": []}`,
			`{"account": {"active-card": true, "available-limit": 50}, "violations": []}`,
			`{"account": {"active-card": true, "available-limit": 40}, "violations": []}`,
			fmt.Sprintf(`{"account": {"active-card": true, "available-limit": 40}, "violations": ["%v", "%v"]}`, account.Violations[4], account.Violations[5]),
		},
	},
}

// Test executer initializtion.
func TestInitExecuter(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	assert.Equal(
		&Executer{},
		exe,
		"Expected a new executer.",
	)
}

// Test transaction with account not initialized.
func TestNotAccountInitializedTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["notInitialized"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["notInitialized"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}

// Test account not active transaction.
func TestNotActiveAccountTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["NotActive"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["NotActive"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}

// Test account operation with account already initialized.
func TestAlreadyinitializedAccountTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["AlreadyInitialized"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["AlreadyInitialized"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}

// Test operation with account limit not sufficient.
func TestInsuficientLimitTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["insufficientLimit"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["insufficientLimit"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}

// Test a doubled transation.
func TestDoubledTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["DoubledTransaction"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["DoubledTransaction"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}

// Test a Blocked merchant transation.
func TestBlockedMerchantTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["blockedMerchant"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["blockedMerchant"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}

// Test a high frequency limit break transaction.
func TestHighFrequencyTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["HighFrequency"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["HighFrequency"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}

// Test a high frequency limit break and doubled transaction.
func TestHighFrequencyDoubledTransaction(t *testing.T) {
	exe := Init()
	assert := assert.New(t)
	for index, value := range tinputs["HighFrequencyDoubled"]["out"] {
		assert.Equal(
			value,
			exe.Exec(tinputs["HighFrequencyDoubled"]["in"][index]),
			"Expected same output from execution.",
		)
	}
}
