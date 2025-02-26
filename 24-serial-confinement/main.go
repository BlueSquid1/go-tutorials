package main

import "fmt"

type Cake struct {
	state string
	id    int
}

func baker(cooked chan<- *Cake, raw <-chan *Cake) {
	for cake := range raw {
		cake.state = "cooked"
		cooked <- cake
	}
}

func icer(iced chan<- *Cake, cooked <-chan *Cake) {
	for cake := range cooked {
		cake.state = "iced"
		iced <- cake
	}
}

func main() {
	icingChan := make(chan *Cake, 30)
	cookingChan := make(chan *Cake)
	rawChan := make(chan *Cake)

	for i := 0; i < 5; i++ {
		go baker(cookingChan, rawChan)
		go icer(icingChan, cookingChan)
	}

	// Create 30 cakes
	for i := 0; i < 30; i++ {
		cake := new(Cake)
		cake.state = "raw"
		cake.id = i
		rawChan <- cake
	}

	// Deliver 30 cakes
	for i := 0; i < 30; i++ {
		cake := <-icingChan
		fmt.Printf("id: %v state: %v\n", cake.id, cake.state)
	}

	// Stop the goroutines
	close(icingChan)
	close(cookingChan)
	close(rawChan)
}
