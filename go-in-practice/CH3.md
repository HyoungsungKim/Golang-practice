# CH3 Concurrency in Go

## 3.1 Understanding Go's concurrency model

Go uses the concurrency model called Communicating Sequential Processes ( CSP ).

Two crucial concepts make Go’s concurrency model work:

- Goroutines : A goroutine is a function that runs independently of the function that started it. Sometimes Go developers explain a goroutine as a function that runs as if it were on its own thread.
- Channels : A channel is a pipeline for sending and receiving data. Think of it as a socket that runs inside your program. Channels provide a way for one goroutine to send structured data to another.

## 3.2 Working with goroutines

### Using goroutine closures

Problem
You want to use a one-shot function in a way that doesn’t block the calling function, and you’d like to make sure that it runs. This use case frequently arises when you want to, say, read a file in the background, send messages to a remote log server, or save a current state without pausing the program.

Solution
Use a closure function and give the scheduler opportunity to run.

In fact, if your Go program can use only one processor, you can almost be sure that it won’t run immediately. Instead, the scheduler will continue executing the outer function until a circumstance arises that causes it to switch to another task.

If there is no `runtime.Gosched()` The goroutine never executes. Why? ***The main function returns (terminating the program) before the scheduler has a chance to run the goroutine.*** At best, you can indicate to the scheduler only that the present goroutine is at a point where it can or should pause.

### Waiting for Goroutine

Sometimes you’ll want to start multiple goroutines but not continue working until
those goroutines have completed their objective. Go wait groups are a simple way to
achieve this.

Problem
One goroutine needs to start one or more other goroutines, and then wait for them to finish. In this practical example, you’ll focus on a more specific problem: you want to compress multiple files as fast as possible and then display a summary.

Solution
Run individual tasks inside goroutines. ***Use `sync.WaitGroup` to signal the outer process that the goroutines are done and it can safely continue.*** several workers are started, and work is delegated to the workers. One process delegates the tasks to the workers and then waits for them to complete.

> `sync.WaitGroup` for telling one goroutine to wait until other goroutines complete

Now here’s the trick: you want to compress a bunch of files in parallel, but have the parent goroutine ( main ) wait around until all of the workers are done. You can easily accomplish this with a wait group. you’ll modify the code in such a way that you don’t change the compress function at all. This is generally considered better design because it doesn’t require your worker function ( compress ) to use a wait group in cases where files need to be compressed serially.

```go
...
go func(fileName string) {
    ...
}(file)
...
```

If your loop runs five times, you’ll have five goroutines scheduled, but possibly none of them executed. And on each of those five iterations, the value of file will change. By the time the goroutines execute, they may all have the same (fifth) version of the file string. That isn’t what you want. You want each to be scheduled with that iteration’s value of file, ***so you pass it as a function parameter, which ensures that the value of file is passed to each goroutine as it’s scheduled.***

### Locking with a mutex

Problem
Multiple goroutines need to access or modify the same piece of data

Solution
One simple way to avoid this situation is for each goroutine to place a “lock” on a resource that it’s using, and then unlock the resource when it’s done. For all other goroutines, when they see the lock, they wait until the lock is removed before attempting to lock that resource on their own. Use sync.Mutex to lock and unlock the object.

***Sometimes it’s useful to allow multiple read operations on a piece of data, but to allow only one write (and no reads) during a write operation.*** The `sync.RWLock` provides this functionality.

## 3.3 Working with Channels

Channels provide a way to send messages from one goroutine to another. But unlike network connections, channels are typed and can send structured data. There’s generally no need to marshal data onto a channel.

### Using multiple channels

Sometimes the best way to solve concurrency problems in Go is to communicate more information. And that often translates into using more channels.

Problem
You want to use channels to send data from one goroutine to another, and be able to interrupt that process to exit.

Solution
Use `select` and multiple channels. It’s a common practice in Go to use channels to signal when something is done or ready to close.

The `select` statement can watch multiple channels (zero or more). Until something happens, it’ll wait (or execute a `default` statement, if supplied). When a channel has an event, the `select` statement will execute that event.

if no default is specified, `select` blocks until one of the case statements can send or receive.

### Closing Channels

What happens if you have a sender and receiver goroutine, and the sender finishes sending data? Are the receiver and channel automatically cleaned up? Nope. The memory manager will only clean up values that it can ensure won’t be used again, and in our example, an open channel and a goroutine can’t be safely cleaned.

The question arises: ***how can you correctly and safely clean up when you’re using goroutines and channels?***

Problem
You don’t want leftover channels and goroutines to consume resources and cause leaky applications. You want to safely close channels and exit goroutines.

Solution
The straightforward answer to the question “How do I avoid leaking channels and goroutines?” is “Close your channels and return from your goroutines.” Although that answer is correct, it’s also incomplete.
***Closing channels the wrong way will cause your program to panic or leak goroutines.*** The predominant method for avoiding unsafe channel closing is to use additional channels to notify goroutines when it’s safe to close a channel.

- Improper way : sending on closed channel

***Proper way***

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	msg := make(chan string)
	done := make(chan bool)
	until := time.After(5 * time.Second)

	go send(msg, done)

	for {
		select {
		case m := <-msg:
			fmt.Println(m)
		case <-until:
			done <- true
			time.Sleep(500 * time.Millisecond)
			return
		}
	}
}

func send(ch chan<- string, done <-chan bool) {
	for {
		select {
		case <-done:
			println("Done")
			close(ch)
			return
		default:
			ch <- "Hello"
			time.Sleep(500 * time.Millisecond)
		}
	}
}
```

> send 함수에 채널 닫는 `select` 만듬

### Locking with buffered channels

You’ve looked at channels that contain one value at a time and are created like this: `make(chan TYPE)` . This is called an `unbuffered channel`.

Problem
In a particularly sensitive portion of code, ***you need to lock certain resources.*** Given the frequent use of channels in your code, you’d like to do this with channels instead of the `sync` package.

Solution
***Use a channel with a buffer size of 1,*** and share the channel among the goroutines you want to synchronize.