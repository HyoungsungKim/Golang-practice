package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type configuration struct {
	Enabled bool
	Path    string
}

func main() {
	var checkDefault bool
	fmt.Println(checkDefault)

	temp := configuration{}
	fmt.Println(temp.Enabled)

	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := configuration{}
	fmt.Println(conf.Enabled)
	err := decoder.Decode(&conf)
	fmt.Println(conf.Enabled)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(conf.Enabled)
	fmt.Println(conf.Path)
}
