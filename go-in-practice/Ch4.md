# Ch4 Habdling errors and panics

Go distinguishes between errors and panics—two types of bad things that can happen during program execution.

- An error indicates that a particular task couldn’t be completed successfully.
- A panic indicates that a severe event occurred, probably as a result of a programmer error.

This chapter presents a thorough look at each category.

## 4.1 Error handling

### Minimize the nils

Problem
Returning nil results along with errors isn’t always the best practice. It puts more work on your library’s users, provides little useful information, and makes recovery harder.

Solution
When it makes sense, avail yourself of Go’s powerful multiple returns and ***send back
not just an error, but also a usable value.***

### Customer error types

In some cases, you may want your errors to contain more information than a simple string. In such cases, you may choose to create a custom error type.

Problem
Your function returns an error. Important details regarding this error might lead users of this function to code differently, depending on these details.

Solution
Create a type that implements the error interface but provides additional functionality.

### Error variable

Problem
One complex function may encounter more than one kind of error. And it’s useful to users to indicate which kind of error was returned so that the ensuing applications can appropriately handle each error case. But although distinct error conditions may occur, none of them needs extra information

Solution
One convention that’s considered good practice in Go (although not in certain other languages) is to ***create package-scoped error variables*** that can be returned whenever a certain error occurs. The best example of this in the Go standard library comes in the io package, which contains errors such as `io.EOF` and `io.ErrNoProgress`.

## 4.2 Differentiating panics from errors

Panics, on the other hand, are unexpected. They occur when a constraint or limitation is unpredictably surpassed. When it comes to declaring a panic in your code, the general rule of thumb is don’t panic unless there’s no clear way to handle the condition within the present context. ***When possible, return errors instead.***

### 4.2.2 Working with panics

#### Issuing panics

The definition of Go’s panic function can be expressed like this: `panic(interface{})`.

Problem
When you raise a panic, what should you pass into the function? Are there ways of panicking that are useful or idiomatic?

Solution
The best thing to pass to a panic is an error. Use the error type to make it easy for the recovery function (if there is one).

The best thing to pass a panic (under normal circumstances, at least) is something that fulfills the error interface. There are two good reasons for this.

- ***The first is that it’s intuitive.***
- ***The second reason is that it eases handling of a panic.***

```go
package main
import "errors"

func main() {
    panic(errors.New("Something bad happened."))
}
```

### 4.2.3 Recovering from panics

Panic recovery in Go depends on a feature of the language called `deferred functions`.

The `defer` statement is a great way to close files or sockets when you’re finished, free up resources such as database handles, or handle panics.

#### Recovering from panics

Problem
A function your application calls is panicking, and as a result your program is crashing

Solution
Use a deferred function and call recover to find out what happened and handle the panic.

Go provides a way of capturing information from a panic and, in so doing, stopping
the panic from unwinding the function stack further. The `recover` function retrieves
the data.

The `recover` function in Go returns a value ( interface{} ) if a panic has been raised, but in all other cases it returns nil. The value returned is whatever value was passed into the panic.