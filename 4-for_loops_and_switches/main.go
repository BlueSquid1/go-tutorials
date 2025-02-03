package main

import (
	"fmt"
	"math/rand"
)

func main() {
	list := make([]int, 100)
	for i := 0; i < len(list); i++ {
		list[i] = rand.Intn(101)
	}

	for _, v := range list {
		switch {
		case v%2 == 0 && v%3 == 0:
			fmt.Println("Six!")
		case v%2 == 0:
			fmt.Println("Two!")
		case v%3 == 0:
			fmt.Println("Three!")
		default:
			fmt.Println("Never mind")
		}
	}
}
