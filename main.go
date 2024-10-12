package main

import (
	"fmt"
	"log"

	"github.com/coding-for-fun-org/go-playground/pkg/accounts"
)

func main() {
	account := accounts.BankAccount{}
	account.Deposit(100)
	fmt.Println(account)
	err := account.Withdraw(50)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(account)
	account.ChangeOwner("new owner")
	fmt.Println(account)
}
