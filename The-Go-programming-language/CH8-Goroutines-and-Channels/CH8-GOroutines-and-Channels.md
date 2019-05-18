# CH8 Go routines and Channels

***Go enables two styles of concurrent programming.*** This chapter presents goroutines and channels, which support communicating sequential processes or CSP, a model of concurrency in which values are passed between independent activities (goroutines) but variables are for the most part confined to a single activity.

## 8.1 Goroutines

In Go, each concurrently executing activity is called a goroutine. Consider a program that has two functions, one that does some computation and one that writes some output, and assume that neither function calls the other. A sequential program may call one function and then call the other, ***but in a concurrent program with two or more goroutines, calls to both functions can be active at the same time. We’ll see such a program in a moment.***

> The differences between threads and goroutines are essentially quantitative , not qualitative, and will be describ ed in Section 9.8.
>
> quantitative  : 양적인
>
> Syntactically : 구문론의

When a program starts, ***its only goroutine is the one that calls the main function,*** so we call it the ***main goroutine.*** New goroutines are created by the go statement. Syntactically, a go statement is an ordinary function or method call prefixed by the key word go. ***A go statement causes the function to be called in a newly created goroutine. The go statement itself completes immediately:***

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n)
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

```

## 8.2 Example: Concurrent Clock Server

In this section, we’ll introduce the net package, which provides the components for building networked client and server programs that communicate over TCP, UDP, or Unix domain sockets. 

```go
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		handleConn(conn)
	}
}
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
```

The Listen function creates a net.Listener, an object that listens for incoming connections on a network port, in this case TCP port localhost:8000

The handleConn function handles one complete client connection. In a loop, it writes the current time, time.Now(), to the client. Since net.Conn satisfies the io.Writer interface, we can write directly to it. ***The loop ends when the write fails, most likely because the client has disconnected, at which point handleConn closes its side of the connection using a deferred call to Close and goes back to waiting for another connection request.***

```go
handleConn(conn)
go handleConn(conn)
```

> handleConn(conn) : It is not worked concurrently
>
> go handleConn(conn) : it is worked concurrently
>
> if run *netcat* in multi-terminal, in handleConn(conn), after finishing first one, second one start
>
> But in go handleConn(conn), multi-terminal work concurrently

## 8.3 Example:Concurrent Echo Server