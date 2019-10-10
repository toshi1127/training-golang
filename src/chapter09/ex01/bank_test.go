package bank_test

import (
	"fmt"
	"sync"
	"testing"

	"./bank"
)

func TestWithdraw(t *testing.T) {
	bank.Deposit(300)

	ts := []struct {
		amount  int
		rest    int
		success bool
	}{
		{200, 100, true},
		{100, 0, true},
		{50, 0, false},
	}

	for _, tc := range ts {
		ok := bank.Withdraw(tc.amount)
		if ok != tc.success {
			t.Errorf("got %v, want %v", ok, tc.success)
		}
		b := bank.Balance()
		if b != tc.rest {
			t.Errorf("balance %d, want %d", b, tc.rest)
		}
	}
}