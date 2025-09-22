package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Length, Width float64
}

func (r *Rectangle) Area() float64 {
	return r.Length * r.Width
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Length + r.Width)
}

type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func PrintShapeDetails(s Shape) {
	fmt.Printf("Area: %f\n", s.Area())
	fmt.Printf("Perimeter: %f\n", s.Perimeter())
}

func main() {
	rect := &Rectangle{Width: 5.0, Length: 7.0}

	PrintShapeDetails(rect)

	cir := &Circle{Radius: 5.0}

	PrintShapeDetails(cir)
}
