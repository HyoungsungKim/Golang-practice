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

In this section, we’ll build an echo server that uses multiple goroutines per connection. Most echo servers merely write whatever they read, which can be done with this trivial version of handleConn:

## 8.4 Channels

***If goroutines are the activities of a concurrent Go program, channels are the connections between them.***

To create a channel, we use the built-in make function:

```go
ch := make(chan int)	//ch has type 'chan int'
```

A channel is a reference to the data structure created by make.

When we copy a channel or pass one as an argument to a function, ***we are copying a reference, so caller and callee refer to the same data structure.*** As with other reference types, ***the zero value of a channel is nil.***

A channel has two principal operations collectively known as communications.

> principal : 주요한, 주된

- send

***A send statement transmits a value from one goroutine, through the channel,*** to another goroutine executing a corresponding receive expression. Both operations are written using the <- operator.

```go
ch <- x // a send statement
x = <- ch // a receive expression in an assignment statment
<- ch // a receive statement; result is discarded
```

- receive

Channels support a third operation, close, which sets a flag indicating that no more values will ever be sent on this channel; subsequent attempts to send will panic. 

> subsequent : 그 다음의, 차후의

***closed channel yield the values that have been sent until no more values are left;*** any receive operations there after complete immediately and yield the zero value of the channel’s element type.

```go
close(ch) //To close channel
ch = make(chan int)		//unbuffered channel
ch = make(chan int, 0)	//unbuffered channel
ch = make(chan int, 3)	// buffered channel with capacity 3
```

### 8.4.1 Unbuffered Channels

A send operation on an unbuffered channel ***blocks the sending goroutine until another goroutine executes a corresponding receive on the same channel,*** at which point the value is transmitted and both goroutines may continue.

Conversely, ***if the receive operation was attempted first, the receiving goroutine is blocked until another goroutine performs a send on the same channel***

Communication over an unbuffered channel causes the sending and receiving goroutines to synchronize. Because of this, ***unbuffered channels are sometimes called synchronous channels.*** 

it’s necessary to order certain events during the program’s execution ***to avoid the problems that arise when two goroutines access the same variable concurrently.***

### 8.4.2 Pipelines

Channels can be used to connect goroutines together ***so that the output of one is the input to another. This is called a pipeline.***

> Output is input of others -> pipeline

```go
package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	for x := range squares {
		fmt.Println(x)
	}
}
```

***You needn’t close every channel when you've finished with it.*** It’s only necessary to ***close a channel when it is important to tell the receiving goroutines that all data have been sent.***

### 8.4.3 Unidirectional Channel Types

This arrangement is typical. When a channel is supplied as a function parameter, it is nearly always with the intent that it be used exclusively for sending or exclusively for receiving.
To document this intent and prevent misuse, ***the Go type system provides unidirectional channel types that expose only one or the other of the send and receive operations.***

> intent  : 강한 관심을 보이는

- The type *chan<- int*, a send-only channel of int, allows sends but not receives.
- Conversely, the type *<-chan int*, a receive-only channel of int, allows receives but not sends.