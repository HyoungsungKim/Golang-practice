package main

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type configuration struct {
	Enabled bool
	Path    string
}

func main() {
	file, _ := os.Open("conf.yaml")
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	conf := configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conf.Enabled)
	fmt.Println(conf.Path)
}
