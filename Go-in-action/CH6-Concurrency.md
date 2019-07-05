# CH6 Concurrency

Concurrency in Go is the ability for functions to run independent of each other. When a function is created as a goroutine, ***it's treated as an independent unit of work that gets scheduled and then executed on an available logical processor.

The Go runtime scheduler is a sophisticated piece of software that manages all the goroutines that are created and need processor time. The scheduler sits on top of the operating system, binding operating system's threads to logical processors which, in turn, execute goroutines. The scheduler controls everything related to which goroutines are running on which logical processors at any given time. Concurrency synchronization comes from a paradigm called communicating sequential processes or CSP. CSP is a message-passing model that works by communicating data between goroutines instead of locking data to synchronize access. The key data type for synchronizing and passing messages between goroutines is called a channel. 

## 6.1 Concurrency versus parallelism

Let's start by learning at a high level what operating system `processes` and `threads` are. This will help you understand later on how the Go runtime scheduler works with the operating system to run goroutines concurrently. ***When you run an application, such as an IDE or editor, the operating system starts a process for the application.*** You can think of a process like a container that holds all the resources an application uses and maintains as it runs. These resources include but are not limited to a memory address space, handles to files, devices, and threads.

***A thread is a path of execution that's scheduled by the operating system to run the code that you write in your functions.*** Each process contains at least one thread, and the initial thread for each process is called the main thread. When the main thread terminates, the application terminates, because this path of the execution is the origin for the application. The operating system schedules threads to run against processors regardless of the process they belong to. 

> The process maintains a memory address space, handles to files, and devices and threads for a running application. The OS scheduler decides which threads will receive time on any given CPU.

> The Go runtime schedules goroutines to run in a logical processor that is bound to a single operating system thread. When goroutines are runnable, they are added to a logical processor's run queue.

> When a goroutine makes a blocking syscall, the scheduler will detach the thread from the processor and create
> a new thread to service that processor.

Sometimes a running goroutine may need to perform a blocking syscall, such as opening a file. When this happens, the thread and goroutine are detached from the logical processor and the thread continues to block waiting for the syscall to return. In the meantime, there's a logical processor without a thread. So the scheduler creates a new thread and attaches it to the logical processor. Then the scheduler will choose another goroutine from the local run queue for execution. ***Once the syscall returns, the goroutine is placed back into a local run queue, and the thread is put aside for future use.***

There's no restriction built into the scheduler for the number of logical processors that can be created. But the runtime limits each program to a maximum of 10,000 threads by default. This value can be changed by calling the `SetMaxThreads` function from the `runtime/debug` package. If any program attempts to use more threads, the program crashes.

Concurrency is not parallelism. Parallelism can only be achieved when multiple pieces of code are executing simultaneously against different physical processors.

- Parallelism is about doing a lot of things at once.
- Concurrency is about managing a lot of things at once. In many cases, concurrency can outperform(능가하다) parallelism, because the strain on the operating system and hardware is much less, which allows the system to do more.

 It's not recommended to blindly change the runtime default for a logical processor. The scheduler contains intelligent algorithms that are updated and improved with every release of Go. If you're seeing performance issues that you believe could be resolved by changing the number of logical processors, you have the ability to do so.

## 6.2 Goroutines

```go
package main
import (
	"fmt"
    "runtime"
    "sync"
)

func main() {
    runtime.GOMAXPROCS(1)
    var wg sync.WaitGroup
    wg.Add(2)
    
    fmt.Println("Start Goroutine")
    
    go func() {
        defer wg.Done()
        for count := 0; count < 3; count++ {
            for char := 'a'; char < 'a' + 26; char++ {
                fmt.Printf("%c", char)
            }
        }
    }()
    
    go func() {
        defer wg.Done()        
        for count := 0; count < 3; count++ {
            for char := 'A'; char < 'A'+26; char++ {
                fmt.Printf("%c", char)
            }
        }
    }()
    
    fmt.Println("Waiting To Finish")
    wg.Wait()
    fmt.Println("\nTerminating Program")
}
```

The amount of time it takes the first goroutine to finish displaying the alphabet is so small that it can complete its work before the scheduler swaps it out for the second goroutine. This is why you see the entire alphabet in capital letters first and then in lowercase letters second. The two goroutines we created ran concurrently, one after the other, performing their individual task of displaying the alphabet.

Once the two anonymous functions are created as goroutines, the code in main keeps running. This means that the ***main function can return before the goroutines complete their work.*** If this happens, the program will terminate before the goroutines have a chance to run.

The keyword `defer` is used to schedule other functions from inside the executing function to be called when the function returns. In the case of our sample program, we use the keyword `defer` to guarantee that the method call to Done is made once each goroutine is finished with its work.

***Based on the internal algorithms of the scheduler, a running goroutine can be stopped and rescheduled to run again before it finishes its work.*** The scheduler does this to prevent any single goroutine from holding the logical processor hostage. It will stop the currently running goroutine and give another runnable goroutine a chance to run

> That is why capital letter is printed first.

```go
//Listing 6.4

```

