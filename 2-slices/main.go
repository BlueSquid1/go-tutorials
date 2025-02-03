package main

import "fmt"

func main() {
	x := []int{1, 4, 5}

	// shallow copy like python
	y := x

	// deep copy
	z := make([]int, len(x))
	copy(z, x)

	y[0] = 3

	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(z)

	//a := []string{"abc", "def", "ghi"}
	a := "abc"
	// print out string
	fmt.Println(a[0:1])
	fmt.Println(a[0])
}
