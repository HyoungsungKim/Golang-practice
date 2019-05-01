package main

import (
	"fmt"
)

func main() {
	type Point struct{ X, Y int }
	A := &Point{1, 2}
	B := &Point{1, 2}
	fmt.Println(A == B)
	fmt.Println(A.X == B.X)

}
