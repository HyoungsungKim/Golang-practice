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

### 8.4.4 Buffered Channels

A buffered channel has a queue of elements. The queue’s maximum size is determined when it is created, by the capacity argument to make.

```go
ch = make(chan string, 3)
```

A send operation on a ***buffered channel inserts an element at the back of the queue***, and ***a receive operation removes an element from the front.*** If the channel is full, the send operation blocks its goroutine until space is made available by another goroutine's receive . Conversely, if the channel is empty, a receive operation blocks until a value is sent by another goroutine.

> Channel is full -> send operation block
>
> Channel is empty -> receive operation block

```go
fmt.Println(<-ch)	//"A'"
```

In this example, the send and receive operations were all performed by the same goroutine, but in real programs they are usually executed by different goroutines. Novices are sometimes tempted to use buffered channels within a single goroutine as a queue, lured by their pleasingly simple syntax, but this is a mistake. ***Channels are deeply connected to goroutine scheduling, and without another goroutine receiving from the channel, a sender risks becoming blocked forever.*** If all you need is a simple queue, make on e using a slice.

Had we used an unbuffered channel, the two slower goroutines would have gotten stuck trying to send their responses on a channel from which no goroutine will ever receive . This situation, called a *goroutine leak*, would be a bug .

***Unlike garbage variables, leaked goroutines are not automatically collected, so it is important to make sure that goroutines terminate themselves when no longer needed.***

The choice between unbuffered and buffered channels, and the choice of a buffered channel’s capacity, may both affect the correctness of a program.

- Unbuffered channels give stronger synchronization guarantees because every send operation is synchronized with its corresponding receive;
- with buffered channels, these operations are decoupled. Also, when we know an upper bound on the number of values that will be sent on a channel

it’s not unusual to create a buffered channel of that size and perform all the sends before the first value is received. ***Failure to allocate sufficient buffer capacity would cause the program to deadlock.***

  >  decoupled : 비결합

Channel buffering may also affect program performance. Imagine three cooks in a cake shop, one baking, one icing, and one inscribing each cake before passing it on to the next cook in the assembly line. ***In a kitchen with little space, each cook that has finished a cake must wait for the next cook to become ready to accept it;*** this rendezvous is analogous to communication over an unbuffered channel.

> analogous : 다양한

## 8.5 Looping in Parallel

In this section, we’ll explore some common concurrency patterns for executing all the iterations of a loop in parallel.

```go
// makeThumbnails makes thumbnails of the specified files.
func makeThumbnails(filenames []string) {
    for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
            log.Println(err)
        }
    }
}
// NOTE: incorrect!
func makeThumbnails2(filenames []string) {
    for _, f := range filenames {
		go thumbnail.ImageFile(f) // NOTE: ignoring errors
    }
}
```

If there’s no parallelism, how can the concurrent version possibly run faster? The answer is that makeThumbnails returns before it has finished doing what it was supposed to do. ***It starts all the goroutines, one per file name, but doesn't wait for them to finish.***

There is no direct way to wait until a goroutine has finished, but we can change the inner goroutine to report its completion to the outer goroutine by sending an event on a shared channel.

```go
// makeThumnails3 makes thumbnails of the specified files in parallel.
func makeThumnails3(filenames []string) {
    ch := make(chan struct{})
    for _, f := range filenames {
        go func(f string) {
            thumnail.ImageFile(f)	// ignore errors
            ch <- struct{}{}
        }(f)
    }
    for range filenames {
        <-ch
    }
}

```

Notice that we passed the value of f as an explicit argument to the literal function instead of using the declaration of *f* from the enclosing for loop. Above, the single variable f is shared by all the anonymous function values and updated by successive loop iterations.

What if we want to return values from each worker goroutine to the main one?

```go
// makeThumnails4 makes thumbnails for the specified files in parallel.
func makeThumbnail4(filenames []string) error {
    go func(f string) {
        _, err := thumbnail.ImageFile(f)
        error <- err
    }(f)
    for range filenames {
        if err := <-errors; err != nil {
            return err
        }
    }
    return nil
}
```

This function has a subtle bug. ***When it encounters the first non-nil error, it returns the error to the caller, leaving no goroutine draining the errors channel.*** Each remaining worker goroutine will block forever when it tries to send a value on that channel, and will never terminate.

The simplest solution is to ***use a buffered channel with sufficient capacity*** that no worker goroutine will block when it sends a message .

```go
func makeThumbails5(filenames []string) (thumbfiles []string, err error) {
    type item struct {
        thumbfile	 string
        err			error
    }
    ch := make(chan item, len(filenames))
    for _, f := range filenames {
        go func(f string) {
            var it item
            it.thumbfile, it.err = thumbnail.ImageFile(f)
            ch <- it
        }(f)
    }
    for range filenames {
        it := <-ch
        if it.err != nil {
            return nil, it.err
        }
        thumbfiles = append(thumbfiles, it.thumbfile)
    }
    return thumbfiles, nil
}
```

To know when the last goroutine has finished (which may not be the last one to start), ***we need to increment a counter before each goroutine starts and decrement it as each goroutine finishes.***

-> This counter type is known as sync.WaitGroup

```go
func makeThumnails6(filenames <- chan string) int64 {
    sizes := make(chan int64)
    var wg sync.WaitFroup
    for f := range filenames {
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            thumb, err := thumbnail.ImageFile(f)
            if err != nil {
                log.Println(err)
                return
            }
            info, _ := os.Stat(thumb)
            sizes <- info.Size()
        }(f)
    }
    
    go func() {
        wg.Wait()
        close(sizes)
    }()
    var total int64
    for size := range sizes {
        total += size
    }
    return total
}
```

## 8.6 Example:Concurrent Web Crawler

The program is to o parallel. ***Unbounded parallelism is rarely a good idea since there is always a limiting factor in the system,*** such as the number of CPU cores for compute-bound work loads, the number of spindles and heads for local disk I/O operations, the bandwidth of the network for streaming downloads, or the serving capacity of a web service.

***The solution is to limit the number of parallel uses of the resource to match the level of parallelism that is available.***

## 8.7 Multiplexing with select