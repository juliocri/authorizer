// ********************************************************************
// * account_test.go                                                  *
// *                                                                  *
// * 2020-03-15 First Version, JR                                     *
// * 2020-03-18 Adds multiple violation scnario, JR                   *
// *                                                                  *
// * This file contains all unit-test representations related         *
// * with the Account struct.                                         *
// *                                                                  *                                                               *
// * Usage: go test -v ./account                                      *
// ********************************************************************

package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test Account limit to all our testcases.
const tlimit = 100

// Test Accounts to cover all our test scenarios.
var taccounts = map[string]*Account{

	"NotInitialzed": nil,

	"NotActive": {
		active: false,
		limit:  tlimit,
	},

	"Active": {
		active: true,
		limit:  tlimit,
	},

	"WithTransaction": {
		active: true,
		limit:  tlimit,
		transactions: []*Transaction{
			ttransactions["Valid"],
		},
	},

	"WithHighFrequencyLimit": {
		active: true,
		limit:  tlimit,
		transactions: []*Transaction{
			{
				Merchant: "Fulanito1",
				Amount:   20,
				Time:     "2019-02-13T10:01:00.000Z",
			},

			{
				Merchant: "Fulanito2",
				Amount:   10,
				Time:     "2019-02-13T10:00:00.000Z",
			},

			{
				Merchant: "Fulanito3",
				Amount:   30,
				Time:     "2019-02-13T10:02:00.000Z",
			},
		},
	},

	"WithHighFrequencyDoubled": {
		active: true,
		limit:  tlimit,
		transactions: []*Transaction{
			{
				Merchant: "Fulanito1",
				Amount:   20,
				Time:     "2019-02-13T10:01:00.000Z",
			},

			{
				Merchant: "Fulanito2",
				Amount:   10,
				Time:     "2019-02-13T10:00:00.000Z",
			},

			ttransactions["Valid"],
		},
	},
}

// Test Transctions to cover our test scenarios.
var ttransactions = map[string]*Transaction{
	"Valid": {
		Merchant: "Fulanito",
		Amount:   tlimit,
		Time:     "2019-02-13T10:00:00.000Z",
	},

	"Insufficient": {
		Merchant: "",
		Amount:   tlimit + 1,
		Time:     "2019-02-13T10:00:00.000Z",
	},
}

// Test if account is not initialized.
func TestAccountNotInitialized(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		false,
		taccounts["NotInitialzed"].Initialized(),
		"Expected a not initialized account.",
	)
}

// Test if account is already initialized.
func TestAccountInitialized(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		true,
		taccounts["NotActive"].Initialized(),
		"Expected an initialized account.",
	)
}

// Test if account is active.
func TestAccountActive(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		true,
		taccounts["Active"].Active(),
		"Expected an active account.",
	)
}

// Test if account is not active.
func TestAccountNotactive(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		false,
		taccounts["NotActive"].Active(),
		"Expected a not active account.",
	)
}

// Test if account retrives its limit value.
func TestAccountLimit(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		tlimit,
		taccounts["Active"].Limit(),
		"Expected same number in account limit.",
	)
}

// Test transaction over not initialized account.
func TestAccountNotInitializedTransaction(t *testing.T) {
	violations := taccounts["NotInitialzed"].ApplyTransaction(ttransactions["Valid"])
	assert := assert.New(t)
	assert.Equal(
		[]int{0},
		violations,
		"Expected array with violation code 0.",
	)
}

// Test transaction over not active account.
func TestAccountNotAcctiveTransaction(t *testing.T) {
	violations := taccounts["NotActive"].ApplyTransaction(ttransactions["Valid"])
	assert := assert.New(t)
	assert.Equal(
		[]int{1},
		violations,
		"Expected array with violation code 1.",
	)
}

// Test transaction with amount higher than the account limit.
func TestAccountInsuficientLimitTransaction(t *testing.T) {
	violations := taccounts["Active"].ApplyTransaction(ttransactions["Insufficient"])
	assert := assert.New(t)
	assert.Equal(
		[]int{3},
		violations,
		"Expected array with violation code 3.",
	)
}

// Test a duplicated transaction.
func TestAccountDoubledTransaction(t *testing.T) {
	violations := taccounts["WithTransaction"].ApplyTransaction(ttransactions["Valid"])
	assert := assert.New(t)
	assert.Equal(
		[]int{4},
		violations,
		"Expected array with violation code 4.",
	)
}

// Test a high frequency limit break transaction.
func TestAccountHighFrequencyTransaction(t *testing.T) {
	violations := taccounts["WithHighFrequencyLimit"].ApplyTransaction(ttransactions["Valid"])
	assert := assert.New(t)
	assert.Equal(
		[]int{5},
		violations,
		"Expected array with violation code 5.",
	)
}

// Test a high frequency limit break and double transaction.
func TestAccountHighFrequencyDoubledTransaction(t *testing.T) {
	violations := taccounts["WithHighFrequencyDoubled"].ApplyTransaction(ttransactions["Valid"])
	assert := assert.New(t)
	assert.Equal(
		[]int{4, 5},
		violations,
		"Expected array with violation code 4 and 5.",
	)
}

// Test a successful transaction.
func TestAccountSuccessfullTransaction(t *testing.T) {
	limit := taccounts["Active"].Limit()
	violations := taccounts["Active"].ApplyTransaction(ttransactions["Valid"])
	assert := assert.New(t)
	assert.Equal(
		limit-ttransactions["Valid"].Amount,
		taccounts["Active"].Limit(),
		"Expected account limit was reduced.",
	)
	assert.Equal(
		[]int{},
		violations,
		"Expected no violations.",
	)
	assert.Equal(
		ttransactions["Valid"],
		taccounts["Active"].transactions[0],
		"Expected transaction was registered.",
	)
}

// Test account initializaton.
func TestAccountInitialization(t *testing.T) {
	acn, v := taccounts["NotInitialzed"].Init(false, 100)
	assert := assert.New(t)
	assert.Equal(
		true,
		acn.Initialized(),
		"Expected a new account initialized.",
	)
	assert.Equal(-1, v, "No violations expected.")
}

// Test account initializaton over account already initialized.
func TestAccountInitializationAccountInitialized(t *testing.T) {
	_, v := taccounts["Active"].Init(false, 100)
	assert := assert.New(t)
	assert.Equal(2, v, "Violations code for initializaton expected.")
}
