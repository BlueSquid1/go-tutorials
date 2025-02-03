package main

import "fmt"

type Person struct {
	firstName string
	lastName  string
	age       int
}

type PersonAgain struct {
	firstName string
	lastName  string
	age       int
}

func main() {
	clint := Person{"Clinton", "Page", 29}
	fmt.Println(clint.firstName, "is", clint.age, "years old")

	clint1 := Person{"Clinton", "Page", 29}

	fmt.Println(clint1 == clint)

	// Can convert between types that have the same fields
	clintAgain := PersonAgain(clint)
	fmt.Println(clintAgain.firstName, "is", clintAgain.age, "years old")
}
