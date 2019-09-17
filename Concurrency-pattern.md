Go concurrency pattern

Based on Concurrency in go tools and techniques for developers

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
  - Channels are composable(구성할 수 있는). Therefore, It is preferable way to unblock multiple goroutines at the same time.

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

- Goroutine with sending to goroutine is wait until buffered channel get a space

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

```go
var dataStream chan interface{}
<-dataStream
// Deadlock!
```

- Reading from a nil channel will block(although not necessarily deadlock) a program.

```go
var dataStream chan interface{}
dataStream <- struct{}{}
// Deadlock!
```

- Write is deadlock too

```go
// Closing deadlock
var dataStream chan interface{}
close(dataStream)
// Panic
```

### Channel ownership

#### About closing channel...

[Is it OK to leave a channel open?](https://stackoverflow.com/questions/8593645/is-it-ok-to-leave-a-channel-open)

***It is important to clarify which goroutine owns a channel in order to reason about our programs logically.***

- Unidirectional channel declarations are the tool that will allow us to distinguish between goroutines that own channels and those that only utilize them
  - Channel owners have a write-access view into the channel ( chan or chan<- )
  - And channel utilizers only have a read-only view into the channel ( <-chan ).
  - Once we make this distinction between channel owners and nonchannel owners, the results from the preceding table follow naturally, and we can begin to assign responsibilities to goroutines that own channels and those that do not.

> 채널 소유자 : 쓰기 가능
>
> 채널 사용자 : 쓰기 불가능, 읽기만 가능

#### Channel owner

1. Instantiate the channel.
2. Perform writes, or pass ownership to another goroutine.
   - It prevents deadlocking
3. ***Close the channel.***
   - It prevents risk of panicing by closing a nil channel
   - It prevents risk of panicing by writing to a closed channel
   - It prevents risk of panicing by closing a channel more than once
4. Encapsulate the previous three things in this list and expose them via a reader channel.

#### Channel utilizer(consumer)

Two things have to worry

- Knowing when a channel is closed
- Responsibly handling blocking for any reason
  - It is very hard to define. Because it depends on algorithms

Example)

```go
chanOwner := func() <- chan int {
    resultStream := make(chan int, 5)
    go func() {
        defer close(resultStream) 
        for i := 0; i <= 5; i++ {
            resultStream <- i
        }
    }()
    return resultStream
}

resultStream := chanOwner()
for result := range resultStream {
    fmt.Printf("received : %d\n", result)
}
fmt.Println("Done receiving")
```

***Keep the scope of channel ownership small so that these things remain obvious***

- If you have a channel as a member variable of a struct with numerous methods on it, it is going to quickly become unclear how the channel will behave

### The select Statement

The `select` statement is the glue that binds channels together

- ***It is how we are able to compose channels together in a program to from larger abstractions.***
- `select` statements can help safely bring channels together with concepts like cancellations, timeouts, waiting, and default values.

Example)

```go
var c1, c2 <-chan interface{}
var c3 <-chan intergace{}
select {
case <- c1:
    // Do something
case <- c2:
    // Do somthing
case c3 <- struct{}{}
    // Do somthing
}
```

- Instead, all channel reads and writes are considered simultaneously to see if any of them are ready:
  - populated or closed channels in the case of reads, and channels that are not at capacity in the case of writes.
  - ***If none of the channels are ready, the entire select statement blocks.***
    - Then when one the channels is ready, that operation will proceed, and its corresponding statements will execute.

```go
start := time.Now()
c := make(chan interface{})
go func() {
    time.Sleep(5*time.Second)
    close(c)
}()

fmt.Println("Blocking on read...")
select {
case <-c:
    fmt.Printf("Unblocked %v later.\n", time.Since(start))
}
```

- Read multiple channels simultaneously

```go
c1 := make(chan interface{}); close(c1)
c2 := make(chan interface{}); close(c2)

var c1Count, c2Count int
for i := 1000; i >= 0; i-- {
    select {
    // Check channel is closed or not
    // If channel is closed, then choose one of <-c1 or <-c2
    case <-c1:
        c1Count++
    case <-c2:
        c2Count++
    }
}
fmt.Printf("c1Count : %v, c2Count : %v", c1Count, c2Count)
// c1Count : 528, c2Count : 472
```

It means that of your set of case statements, ***each has an equal chance of being selected as all the others***

***What happens if there are never any channels that become ready?***

- If there is nothing useful you can do when all the channels are blocked, but you also cannot block forever, you may want to time out.
- Go's `time` package provides an elegant way to do this with channels that fits nice;y within the paradigm of `select` statements.

```go
var c <-chan int
select {
    case <-c:
    case <-time.After(1*time.Second):
    fmt.Println("Timed out.")
}
```

***What happens when no channel is ready and we need to do something in the meantime?***

- `select` statement also allows for a default clause in case you would like to do something if all the channels you are selecting against are blocking
- This allows you to exit a select block without blocking.

```go
start := time.Now()
var c1, c2 <-chan int
select {
case <- c1:
case <- c2:
default:
    fmt.Printf("In default after %v\n\n", time.Since(start))
}
```

Block forever

```go
select{}
```

## Concurrency Patterns in Go