package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
	"unicode"
	"unicode/utf8"
)

type Employee struct {
	Id        int
	Name      int
	Address   string
	DoB       time.Time
	Position  string
	Salary    int
	managerID int
}

func main() {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune \t count \n")
	for c, n := range counts {
		fmt.Printf("%q \t %d \n", c, n)
	}
	fmt.Print("\n len \t count \n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d \t %d \n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters \n", invalid)
	}
}
