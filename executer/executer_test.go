// ********************************************************************
// * account_test.go                                                  *
// *                                                                  *
// * 2020-03-15 First Version, JR                                     *
// * 2020-03-18 Adds multiple violation scnario, JR                   *
// * 2020-11-18 Simplifies transaction tests in a single function, JR *
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

// Test all transaction types.
func TestTransactions(t *testing.T) {
	for key, data := range tinputs {
		exe := Init()
		assert := assert.New(t)

		t.Run(fmt.Sprintf("%s", key), func(t *testing.T) {
			for index, value := range data["out"] {
				assert.Equal(
					value,
					exe.Exec(data["in"][index]),
					"Expected same output from execution.",
				)
			}
		})
	}
}
