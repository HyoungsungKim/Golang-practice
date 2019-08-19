package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	var file io.ReadCloser

	file, err := OpenCSV("data.csv")
	defer file.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("File is closed?")
}

func OpenCSV(filename string) (file *os.File, err error) {
	defer func() {
		if r := recover(); r != nil {
			file.Close()
			err = r.(error)
		}
	}()

	file, err = os.Open(filename)

	if err != nil {
		fmt.Printf("Failed to open file\n")
		return file, err
	}

	RemoveEmptyLines(file)
	fmt.Println("Touch here")
	return file, err
}

func RemoveEmptyLines(f *os.File) {
	panic(errors.New("Failed parse"))
}
