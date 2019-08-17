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