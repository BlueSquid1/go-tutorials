package main

import (
	"fmt"
	"main/bank"
)

func main() {

	bank.Deposit(1)
	bank.Deposit(2)
	r := bank.Withdraw(3)
	fmt.Println(r)
	fmt.Println(bank.Balance())
}
