package bank

import "sync"

var balance int
var mu sync.Mutex

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func Balance() int {
	mu.Lock()
	defer mu.Unlock()
	return getBalance()
}

func Withdraw(amount int) bool {
	mu.Lock()
	defer mu.Unlock()
	return withdraw(amount)
}

// Private function to seperate logic from thread safe controls so have ability to reuse it without holding a mutex.
func deposit(amount int) {
	balance += amount
}

func getBalance() int {
	return balance
}

func withdraw(amount int) bool {
	if balance-amount < 0 {
		return false
	}
	balance -= amount
	return true
}
