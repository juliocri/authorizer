// ********************************************************************
// * account.go                                                       *
// *                                                                  *
// * 2020-03-13 First Version, JR                                     *
// * 2020-11-17 Adds Blocked merchant feature                         *
// * This package holds all bussiness logic related with an account.  *
// *                                                                  *
// * Usage:                                                           *
// * acn: = account.Init(active, limit)                               *
// * acn.ApplyTransaction(transation)                                 *
// ********************************************************************

package account

import (
	"math"
	"time"
)

// Number of minutes - time window for transation checks.
const maxTimeDiff = 2

// Max number of transactions allowed in $maxTimeDiff (minutes);
const maxTransactions = 3

// Violations - Account violation's codes and its meaning.
var Violations = map[int]string{
	0: "account-not-initialized",
	1: "card-not-active",
	2: "account-already-initialized",
	3: "insufficient-limit",
	4: "doubled-transaction",
	5: "high-frequency-small-interval",
	6: "blocked-merchant",
}

var blockedlist = []string{
	"Burger King",
}

// Account - stores limit, state and authorized transactions,
// for working account.
type Account struct {
	active       bool
	limit        int
	transactions []*Transaction
}

// Transaction - represents transaction fields gotten from json input.
type Transaction struct {
	Merchant string `json:"merchant"`
	Amount   int    `json:"amount"`
	Time     string `json:"time"`
}

// Init - Initializes an account, and return a violation code,
// if account is already initialized.
func (acn *Account) Init(a bool, l int) (*Account, int) {
	// If account is already initialized, we return same account,
	// and violation code.
	if acn.Initialized() {
		return acn, 2
	}

	// Otherwise we prepare a new account.
	account := &Account{
		active:       a,
		limit:        l,
		transactions: []*Transaction{},
	}
	// and return the account and not violations int representation.
	return account, -1
}

// Limit - Returns current limit from account.
func (acn *Account) Limit() int {
	return acn.limit
}

// Active - Returns if account is active or not.
func (acn *Account) Active() bool {
	return acn.active
}

// Initialized - returns if the accoun is wheather or not initialized.
func (acn *Account) Initialized() bool {
	if acn != nil {
		return true
	}

	return false
}

// ApplyTransaction - updates the account's limit if no violations found,
// and registert the transation in the account's history transactions.
// otherwise returns an integer array with violation codes found.
func (acn *Account) ApplyTransaction(tsn *Transaction) []int {
	violations := []int{}
	// Only if the account is initialized, is worthy to look if more,
	// violations are detected for the transaction.
	if acn.Initialized() {
		// If account is active we still continue with validations.
		// Does not make sense try to apply a transaction with an account,
		// that is not active.
		if acn.Active() {
			// Check if transaction is not duplicated.
			duplicated := acn.duplicatedTransaction(tsn)
			if duplicated {
				violations = append(violations, 4)
			}

			// Or if we don't reach the frequency limit.
			overpassed := acn.frenquencyOverpass(tsn)
			if overpassed {
				violations = append(violations, 5)
			}

			// Finally if  not duplicated nor overpassed, we check if the account,
			// has enoght limit to execute the transation.
			if !duplicated && !overpassed {
				if acn.limit < tsn.Amount {
					violations = append(violations, 3)
				}
			}

			// Check if the merchant is not blocked
			merchantBlocked := acn.merchantBlocked(tsn)
			if merchantBlocked {
				violations = append(violations, 6)
			}

		} else {
			// Card not active.
			violations = append(violations, 1)
		}
	} else {
		// Account not initialized.
		violations = append(violations, 0)
	}

	// If no violations found apply the transaction and register in history.
	if len(violations) == 0 {
		acn.limit = acn.limit - tsn.Amount
		acn.registryTransaction(tsn)
	}

	return violations
}

// merchantBlocked - return a boolean if the transaction merchant is within
// the blockedlist.
func (acn *Account) merchantBlocked(tsn *Transaction) bool {
	for _, val := range blockedlist {
		if tsn.Merchant == val {
			return true
		}
	}

	return false
}

// registryTransaction - Adds a new transation to the account authored,
// transactions history.
func (acn *Account) registryTransaction(tsn *Transaction) {
	acn.transactions = append(acn.transactions, tsn)
}

// duplicatedTransaction - check if a transaction with same amount and merchant,
// does not exist in a timeframe of 2 minutes diff.
func (acn *Account) duplicatedTransaction(tsn *Transaction) bool {
	for _, val := range acn.transactions {
		// Convert strings to time golang objects.
		t1, _ := time.Parse(time.RFC3339, val.Time)
		t2, _ := time.Parse(time.RFC3339, tsn.Time)
		// Get diff in minutes (absolute val);
		diff := math.Abs(t1.Sub(t2).Minutes())
		// Check if the transaction in progress have similar values,
		// in authored account's transactions.
		if val.Amount == tsn.Amount &&
			val.Merchant == tsn.Merchant &&
			diff <= maxTimeDiff {
			return true
		}
	}

	return false
}

// frenquencyOverpass - Returns true if the current transaction breaks the,
// number of transactions allowed ($maxTransactions) in the frequency,
// of $maxTimeDiff minutes.
func (acn *Account) frenquencyOverpass(tsn *Transaction) bool {
	allowed := 0
	for _, val := range acn.transactions {
		t1, _ := time.Parse(time.RFC3339, val.Time)
		t2, _ := time.Parse(time.RFC3339, tsn.Time)
		diff := math.Abs(t1.Sub(t2).Minutes())
		if diff <= maxTimeDiff {
			allowed++
		}
		// If we reach the maxnumber of transactions allowed, we can't go for it,
		// and the violation is reported.
		if allowed == maxTransactions {
			return true
		}
	}

	return false
}
