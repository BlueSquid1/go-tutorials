package main

import (
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p1 Point) Distance(p2 Point) float64 {
	return math.Sqrt(math.Pow(math.Abs(p1.X-p2.X), 2.0) + math.Pow(math.Abs(p1.Y-p2.Y), 2.0))
}

type Circle struct {
	Point
	Radius float64
}

func (c Circle) Distance(p2 Point) float64 {
	return math.Abs(c.Radius - c.Point.Distance(p2))
}

type Wheel struct {
	Circle
	Spokes int
}

func main() {
	w1 := Wheel{}
	w1.Radius = 1
	w1.X = 3 // same as w1.Circle.Point.X = 3
	w1.Y = 0

	w2 := Wheel{}
	w2.X = 5
	w2.Y = 0

	// Can use Point's methods since it's embedded into Wheel (via Circle)
	fmt.Println(w1.Distance(w2.Point))       // 1
	fmt.Println(w1.Point.Distance(w2.Point)) // 2
}
