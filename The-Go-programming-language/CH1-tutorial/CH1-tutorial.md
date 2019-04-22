# Ch1-tutorial

## Section 1.3

A map holds a set of key/value pairs and provides constant-t ime operations to store, retrieve,
or test for an item in the set. The key may be of any typ e whose values can compared with ==,
strings being the most common example; the value may be of any typ e at all. In this example,
the keys are strings and the values are ints.

The scanner reads from the program’s standard input. Each cal l to input. Scan() reads the next line and removes the newline character from the end; the result can be retrieved by calling input. Text(). The Scan function returns true if there is a line and false when there is no more input.

Printf, whereas those whose names end in ln follow Println, formatting their arguments as if by %v, followed by a newline. 

```go
// dup2.go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d \t %s \n", n, line)
		}
	}
}
func countLines(f *os.File, counts map[string]int) {
    //function receive copy of reference
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[inputs.Text()]++
	}
}
```



```go
package main

import(
	"fmt"
    "io/ioutil"
    "os"
    "strings"
)

fun main(){
    counts := make(map[string]int)
    for _, filename := range os.Args[1:] {
        data, err := ioutil;.ReadFile(filename)
        if err!- nil {
            fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
            continue
        }
        for _, line := range string.Split(string(data), "\n"){
            counts[line]++
        }
    }
    for line, n := range counts {
        if n > 1 {
            fmt.Printf("%d\t%s\n", n, line)
        }
    }
}
```

