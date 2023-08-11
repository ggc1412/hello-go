package accounts

import "errors"

type account struct {
	owner   string
	balance int
}

var errNoMoney = errors.New("Can't withdraw you are poor")

func NewAccount(name string) *account {
	account := account{owner: name, balance: 0}
	return &account
}

func (a *account) Deposit(amount int) {
	a.balance += amount
}

func (a account) Balance() int {
	return a.balance
}

func (a *account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

func (a account) String() string {
	return "Whatever you want"
}
