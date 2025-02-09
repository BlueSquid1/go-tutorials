package main

import "fmt"

func modifySlice(x []int) {
	// because slice passed by reference this changes the callees slice
	x[1] = 1
	fmt.Printf("before append: %p\n", x)

	// append can changes the reference. If that happens changes to x afterwards will not effect the callee's x
	x = append(x, 6)
	fmt.Printf("after append: %p\n", x)
	x[2] = 2
}

func main() {
	// create size of length 3 and capacity 3
	x := make([]int, 3, 4)
	x[0] = 1
	x[1] = 4
	x[2] = 5

	fmt.Printf("before call: %p\n", x)
	modifySlice(x)
	fmt.Println(x)

	// shallow copy like python
	y := x

	// deep copy
	z := make([]int, len(x))
	copy(z, x)

	y[0] = 3

	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(z)

	a := "abc"
	// print out string
	fmt.Println(a[0:1])
	fmt.Println(a[0])
}
