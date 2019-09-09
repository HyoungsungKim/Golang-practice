# Go concurrency pattern

## The Difference Between Concurrency and Parallelism

Concurrency is a property of the code; parallelism is a property of the running program.

> Classical classification
>
> Concurrency : Many motherboard
>
> Parallel : 1 motherboard

## Sync package

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("First Goroutine is here!")
		time.Sleep(100)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Second Goroutine is here!")
		time.Sleep(200)
	}()

	wg.Wait()
	fmt.Println("Every goroutine is finished!")
}

```

- ***Call Add is done outside of goroutine.*** It is helpful to track/ If we didn't do this, we would have introduced a race condition. Because we don't have a guarantee when goroutine will be scheduled
- ***Goroutine is scheduled randomly!!***

```go
//output of above code
// 1st run
Third Goroutine is here!
Second Goroutine is here!
First Goroutine is here!
Every goroutine is finished!
// 2nd run
Third Goroutine is here!
First Goroutine is here!
Second Goroutine is here!
Every goroutine is finished!
```

```go
func main() {	
    .
    .
    for i := 0; i < 5; i++ {
            wg.Add(1)
            go func(id int) {
                defer wg.Done()
                fmt.Println(id)
            }(i)
        }
	wg.Wait()
}
//output
0
4
1
2
3
// -> It is random!
```

## Mutex and RWMutex

Mutex : Mutual exclusion

- Mutex provides a concurrent-safe way to express exclusive access to a shared resource.
- Channels share memory by communication
- Mutex shares memory by creating a convention developers must follow to synchronize access to the memory.

```go
mutex.Lock()
defer mutex.Unlock()
.
.
```

sync.RWMutex : It guards access to memory, however, RWMutex gives you a little bit more control over the memory.

- Lock for reading
  - Only reading possible
- Lock for writing
  - Only writing possible

## Cond

- Waiting or announcing the occurrence of an event
- Event is any arbitrary signal between two or more goroutines that carries no information other than the fact that it has occurred.
- We need something for waiting event without waiting CPU clock
  - Loop is super inefficient
  - `time.Sleep()` is better than loop, however still inefficient

```go
c := sync.NewCond(&sync.Mutex{})
c.L.Lock()
for conditionTrue() == false {
    c.Wait()
}
c.L.Unlock()
```

- Wait() doesn't waste CPU clock
  - ***It suspends goroutine***
- `Signal()` announce to `Wait()`

## Once

- Only one call

```go
var count int
var once sync.Once
once.Do(increment)
onde.Do(decrement)

//count is 1 not 0
```

## Pool

```go
myPool := &sync.Pool{
    New: func() interface{} {
        fmt.Println("Hi!")
        return struct {}{}
    },
}

myPool.Get()				// Print Hi!
instance := myPool.Get()	// print Hi!
myPool.Put(instance)	
myPool.Get()				// run someting inside in pool, don't run New func of myPool

//output 
//Hi!
//Hi!
```

## Channels

- Read only channel

```go
var dataStream <-chan interface{}
dataStream := make(<-chan interface{})
```

- Send only channel

```go
var dataStream chan <- interface{}
dataStream := make(chan<- interface{})
```

- These will use as function parameters

```go
// Right way
func main() {
	channel := make(chan string)
	go func() { channel <- "Hello world!" }()
	fmt.Println(<-channel)
}
// Wrong way
func main() {
    channel := make(chan string)
    channel <- "Hello world!"
    fmt.Println(<-channel)
}
// Deadlock
// Channel sned to goroutin, but there is no goroutine
// Therefore send is blocked
```

>***Closing a channel is also one of the ways you can signal multiple goroutines simultaneously.*** If you have n goroutines waiting on a single channel, instead of writing n times to the channel to unblock each goroutine, ***you can simply close the channel.***
>
>Since a closed channel can be read from an infinite number of times, it doesn't matter how many goroutines are waiting on it, and closing the channel is both cheaper and faster than performing n writes. Here's an example of unblocking multiple goroutines at once:

```go
begin := make(chan interface{})
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go fun(i int) {
        defer wg.Done()
        // Wait until chanel is closed
        <-begin
        fmt.Printf("%v has begun\n", i)
    }(i)
}

fmt.Println("unblocking goroutines...")
close(begin)
// Output print is started
wg.Wait()
```

***There are 2 ways for waiting signal***

- Using cond
- Using channel

### Buffered Channel

```go
func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- 1
		ch <- 2
		ch <- 3
	}()
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
// Output
// 1
// 2
// 3
```

***When something is pushed in full channel, It is not abandoned, but waiting until channel get a space***

- Goroutine with sending to goroutin is wait until buffered channel get a space

compiler optimization

```go
func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)
		ch <- 10
		ch <- 20
		ch <- 30
		fmt.Println("Push is done!")
	}()

	fmt.Println(len(ch))		// Print 0

	for val := range ch {
		fmt.Println(len(ch))	// Print 2, 1, 0
		fmt.Println(val)		// Print 10, 20, 30
	}

	fmt.Println(len(ch))		// Print 0
}
// Output
0
Push is done!
2
10
1
20
0
30
0
```

- Little bit different with our expectation.
  - Channel size is 2, however 3 integer were pushed to buffered

### nil of Channel