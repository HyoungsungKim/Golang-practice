# CH9 Concurrency with Shared Variables

## 9.1 Race Conditions

But in general, we don't know whether an event x in one goroutine happens before an event y in another goroutine, or happens after it, or is simultaneous with it.

We avoid concurrent access to most variables either by confining them to a single goroutine or by maintaining a higher-level invariant of mutual exclusion.

There are many reasons a function might not work when called concurrently, including deadlock, live lock, and resource starvation.

***A race condition is a situation in which the program does not give the correct result for some interleavings of the operations of multiple goroutines.*** Race conditions are pernicious because they may remain latent in a program and appear infrequently, perhaps only under heavy load or when using certain compilers, platforms, or architectures. This makes them hard to reproduce and diagnose.

> pernicious  : 치명적인

how do we avoid data races in our programs?

A data race occurs whenever two goroutines ***access the same variable concurrently and at least one of the accesses is a write.*** It follows from this definition that there are three ways to avoid a data race.

The first way is not to write the variable

```go
var icons = make(map[string]image.Image)
func loadIcon(name string) image.Image
//Note: not concurrency-safe!
func Icon(name string) image.Image {
    icon, ok := icons[name]
    if !ok {
        icon = loadIcons(name)
        icons[name] = icon
    }
    return icon
}

//Safe way
var icons = map[string]image.Image{
    "spades.png": loadIcon("spades.png"),
    "hearts.png": loadIcon("hearts.png"),
    "diamonds.png": loadIcon("diamonds.png")
    "clubs.png": loadIcon("clubs.png")
}
func Icon(name string) image.Image {return icons[name]}
```

The second way to avoid a data race is to avoid accessing the variable from multiple goroutines.

Since other goroutines cannot access the variable directly, they must use a ***channel*** to send the confining goroutine a request to query or update the variable.

```go
package bank

var deposits = make(chan int)
var balances = make(chan int)

func Deposit(amount int) {deposit <- amout}
func Balance() int		{ return <- balance}
func teller() {
    var balance int
    for {
        select {
        case amount := <- deposits:
            balance += amount
        case balances <- balance:
        }
    }
    func init() {
        go teller()
    }
}
```

The third way to avoid a data race is to allow many goroutines to access the variable, but only one at a time.

## 9.2 Mutual Exclusion:sync.Mutex

A semaphore that counts only 1 is called a binary semaphore

```go
var(
    sema = make(chan struct{}, 1)
    balance int
)

func Deposit(amount int) {
    sema <- struct{}{}	//acquire token
    balance = balance + amount
    <-sema	//release token
}
func Balance() int {
    sema <- struct{}{}	//acquire token
    b := balance
    <-sema	//release token
    return b
}
```

```go
import "sync"
var(
	mu	sync.Mutex
    balance int
)
func Deposit(amount int) {
    mu.Lock()
    balance = balance + amount
    mu.Unlock()
}
func Balance() int {
    mu.Lock()
    b := balance
    mu.Unlock()
    return b
}
```

Each time a goroutine accesses the variables of the bank (just balance here), it must call the mutex’s Lock method to acquire an exclusive lock.

If some other goroutine has acquired the lock, this operation will block until the other goroutine calls Unlock and the lock becomes available again. 

This arrangement of functions, mutex lock, and variables is called a monitor.

Since the critical sections in the Deposit and Balance functions are so short calling Unlock at the end is straightforward. In more complex critical sections, especially those in which errors must be dealt with by returning early, it can be hard to tell that calls to Lock and Unlock are strictly paired on all paths. Go’s `defer` statement comes to the rescue: by deferring a call to Unlock, the critical section implicitly extends to the end of the current function, ***freeing us from having to remember to insert Unlock calls in one or more places far from the call to Lock.***

```go
func Balance() int {
    mu.Lock()
    defer mu.Unlock()
    return balance
}
```

In the example above, the Unlock executes after the return statement has read the value of balance, so the Balance function is concurrency-safe.

***it’s not possible to lock a mutex that’s already locked—this leads to a deadlock where nothing can proceed***

> defer 를 쓸려면 다른 함수에서도 다 defer 써야됨.(특히 defer 쓰는 함수 안에서 다른 함수를 호출하는데 그 함수에서 mutex 사용시 defer 써야됨)

## 9.3 Read/Write Mutexes:sync.RWMutex

