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

