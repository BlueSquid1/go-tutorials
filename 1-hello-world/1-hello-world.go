package main

import (
	"fmt"

	"rsc.io/quote/v4"
)

func main() {
	const age = 29
	var name string = "abc"
	const unused float64 = age
	fmt.Printf("hello %s (%f years old or %d)\n", name, unused, age)
	fmt.Println(quote.Glass())
}
