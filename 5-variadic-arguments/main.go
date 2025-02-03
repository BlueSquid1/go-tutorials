package main

import "fmt"

func addTo(base int, vals ...int) []int {
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		out = append(out, base+v)
	}
	return out
}

func main() {
	a := []int{5, 6, 7}
	b := addTo(3, a...)
	fmt.Println(b)
}
