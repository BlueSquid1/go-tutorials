package main

import (
	"fmt"

	tempconv "github.com/BlueSquid1/go-tutorials/12-packages/tempconv"
)

func main() {
	k := tempconv.Kelvin(10.0)
	c := tempconv.KToC(k)
	fmt.Println(c)
}
