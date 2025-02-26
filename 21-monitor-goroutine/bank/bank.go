package bank

type withdraw struct {
	amount int
	result chan bool
}

var deposits = make(chan int)
var balances = make(chan int)
var withdraws = make(chan withdraw)

func Deposit(amount int) {
	deposits <- amount
}

func Balance() int {
	return <-balances
}

func Withdraw(amount int) bool {
	asyncResult := make(chan bool)
	w := withdraw{amount: amount, result: asyncResult}
	withdraws <- w
	return <-asyncResult
}

// For a monitor goroutine a broker is used to read/write to variables shared across goroutines.
func teller() {
	var balance int
loop:
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case w := <-withdraws:
			if balance-w.amount < 0 {
				w.result <- false
				continue loop
			}
			balance -= w.amount
			w.result <- true
		}
	}
}

func init() {
	// Need to start the broker when the package is loaded
	go teller()
}
