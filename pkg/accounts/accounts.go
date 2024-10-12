package accounts

import "errors"

// BankAccount struct to represent a bank account
type BankAccount struct {
	owner   string
	balance int
}

// NewAccount function to create a new bank account
func NewAccount(owner string) *BankAccount {
	account := BankAccount{owner: owner, balance: 0}

	return &account
}

// Deposit method to deposit money into the bank account
func (ba *BankAccount) Deposit(amount int) {
	ba.balance += amount
}

// Withdraw method to withdraw money from the bank account
func (ba *BankAccount) Withdraw(amount int) error {
	if ba.balance < amount {
		return errors.New("Can't withdraw you are poor")
	}
	ba.balance -= amount
	return nil
}

// Balance method to get the balance of the bank account
func (ba BankAccount) Balance() int {
	return ba.balance
}

// Owner method to get the owner of the bank account
func (ba BankAccount) Owner() string {
	return ba.owner
}

// ChangeOwner method to change the owner of the bank account
func (ba *BankAccount) ChangeOwner(newOwner string) {
	ba.owner = newOwner
}

// String method to represent the bank account as a string
func (ba BankAccount) String() string {
	return "Owner: " + ba.owner + ", Balance: " + string(ba.balance)
}
