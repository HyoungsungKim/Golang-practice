package main

import "fmt"

//getName return "World"
func getName() string {
	return "World!"
}

func main() {
	name := getName()
	fmt.Println(name)
}
