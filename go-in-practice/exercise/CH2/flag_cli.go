package main

import (
	"flag"
	"fmt"
)

var name = flag.String("name", "World", "A name to say hello to.")
var spanish bool

func init() {
	//for long name
	flag.BoolVar(&spanish, "spanish", false, "use Spanish language")
	//for short name
	flag.BoolVar(&spanish, "s", false, "use Spanish language")
}

func main() {
	flag.Parse()
	if spanish == true {
		fmt.Printf("Hola %s!\n", *name)
	} else {
		fmt.Printf("Hello %s!\n", *name)
	}
}
