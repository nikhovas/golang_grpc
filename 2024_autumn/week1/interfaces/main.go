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

// Perimeter method for Rectangle
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Length + r.Width)
}

type Circle struct {
	Radius float64
}

func (r *Circle) Area() float64 {
	return math.Pi * r.Radius * r.Radius
}

func (r *Circle) Perimeter() float64 {
	return 2 * math.Pi * r.Radius
}

func PrintShapeDetails(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())
	fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
	rect := &Rectangle{Length: 10, Width: 5}

	fmt.Println("Rectangle details:")
	PrintShapeDetails(rect)

	c := &Circle{Radius: 5}

	fmt.Println("Circle details:")
	PrintShapeDetails(c)
}
