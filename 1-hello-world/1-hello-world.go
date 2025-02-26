package main

import (
	"fmt"

	"rsc.io/quote/v4"
)

func main() {
	// can run with "go run ."
	// this will compile all the files in the current folder and run the binary
	const age = 29
	var name string = "abc"
	const unused float64 = age
	fmt.Printf("hello %s (%f years old or %d)\n", name, unused, age)
	fmt.Println(quote.Glass())
}
