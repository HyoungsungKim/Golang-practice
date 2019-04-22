package main

import (
	"fmt"
)

func main() {
	var x, y int
	fmt.Println(&x == &x, &x == &y, &x == nil)
}
