package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 0; ; x++ {
			if x > 100 {
				close(naturals)
			}
			naturals <- x
		}
	}()

	go func() {
		for {
			x := <-naturals
			if x > 100 {
				close(squares)
			}
			squares <- x * x
		}
	}()

	for {
		fmt.Println(<-squares)
	}
}
