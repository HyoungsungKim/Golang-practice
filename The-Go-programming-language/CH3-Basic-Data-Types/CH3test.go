package main

import (
	"fmt"
)

func main() {
	var str string
	str = "Hello"
	fmt.Printf(str)
	str[:3] = "a"
	fmt.Printf(str)
}
