# CH5 Functions

## 5.1. Function Declarations

```go
func name(parameter-list) (result-list) {
    body
}
```

If the function returns one unnamed result or no results at all, parentheses are optional and usually omitted.

```go
func hypot(x, y float64) float64 {
    return math.sqrt(x*x + y*y)
}
fmt.Println(hypot(3,4)) // 5

//Equivalent declaration
func f(i, j, k int, s, t string) { /* ... */ }
func f(i int, j int, k int, s string, t string) { /* ... */ }
```

The type of a function is sometimes called its signature.

Go has no concept of default parameter values, nor any way to specif y arguments by name, so the names of parameters and results don’t matter to the caller except as documentation.

***Arguments are passed by value,*** so the function receives a copy of each argument; modifications to the copy do not affect the caller. ***However, if the argument contains some kind of reference, like a pointer, slice, map, function, or channel, then the caller may be affected by any modifications the function makes to variables indirectly referred to by the argument.***

## 5.2. Recursion

```
> fetch http://golang.org | findlinks1
/
/
#
/doc/
/pkg/
/project/
/help/
/blog/
#
#
//tour.golang.org/
/dl/
//blog.golang.org/
https://developers.google.com/site-policies#restrictions
https://creativecommons.org/licenses/by/3.0/
/LICENSE
/doc/tos.html
http://www.google.com/intl/en/policies/privacy/
```

As you can see by experimenting with outline, most HTML documents can be processed with only a few levels of recursion, but it’s not hard to construct pathological web pages that require extremely deep recursion.

***Many programming language implementations use a fixed-size function call stack; sizes from 64KB to 2MB are typical.***  Fixed-size stacks impose a limit on the depth of recursion, so one must be careful to avoid a stack overflow when traversing large data structures recursively;

***In contrast, typical Go implementations use variable-size stacks that start small and grow as needed up to a limit on the order of a gigabyte. This lets us use recursion safely and without worrying about overflow.***

## 5.3 Multiple Return Values

A function can return more than one result.

Go’s garbage collector recycles unused memory, ***but do not assume it will release unused operating system resources like open files and net work connections. They should be closed explicitly.***

```go
//To ignore one of the values, use underscore
links, err := findLinks(url)
links, _ := findLinks(url) //err ignored
```

A multi-value d call may appear as the sole argument when calling a function of multiple parameters. Although rarely used in production code, this feature is sometimes convenient during debugging since it lets us print all the results of a call using a single statement.

```go
log.Println(findLinks(url))
//is same with
links, err := findLinks(url)
log.println(links, err)
```

## 5.4 Errors

some functions always succeed so long as their preconditions are met. The built-in type error is an interface type.

By contrast, Go programs use ordinary control-flow mechanisms like if and return to respond to errors. This style undeniably demands that more attention be paid to error-handling logic, but that is precisely the point.

### 5.4.1 Error-Handling Startegies

- *propagate* the error

```go
resp, err := http.Get(url)
if err != nil {
    return nil, err
}
```

if the cal l to html.Parse fails, findLinks do es not return the HTML parser’s error directly because it lacks two crucial pieces of information:

```go
doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
    return nil, fmt.Errorf("parsing 5s as HTML: %v", url, err)
}
```

When designing error messages, be deliberate, so that each one is a meaningful description of the problem with sufficient and relevant detail, and be consistent, so that errors returned by the same function or by a group of functions in the same package are similar in form and can be dealt with in the same way.

- retry the failed operation

```go
// CH5/wait
```

- if progress is impossible, the caller can print the error and stop the program gracefully, ***but this course of action should generally be reserved for the main package of a program.***

```go
if err := WaitForServer(url); err != nil {
    fmt.Fprintf(os.Stderr, "site is down: %v\n", err)
    os.Exit(1)
}

//More convenient way
if err := WaitForServer(url); err != nil {
    log.Fatalf("Site is down: %v\n", err)
}
```

- In some cases, it is sufficient just to log the error and then continue.

```go
if err := Ping(); err != nil {
    log.Printf("ping failed: %v;", err)
}

if err:= Ping(); err != nil {
    fmt.Fprintf(os.Stderr, "ping failed %v;", err)
}
```

- We can safely ignore an error entirely

```go
dir, err := ioutil.TempDir("", "scratch")
if err != nil {
    return fmt.Errorf("failed to create temp dir : %v", err)
}
os.RemoveAll(dir) //Ignore errors
```



Get into the habit of considering errors after every function call, and when you deliberately ignore one, document your intention clearly.

### 5.4.2 End of File(EOF)

if the caller repeatedly tries to read fixed-size chunks until the file is exhausted, the caller must respond differently to an end-of-file condition than it does to all other errors. ***For this reason, the io package guarantees that any read failure caused by an end-of-file condition is always reported by a distinguished error, io.EOF,***

```go
package io
import "errors"
var EOF = errors.New("EOF")

in := bufio.NewReader(os.Stdin)
for {
    r, _, err := in.ReadRune()
    if err == io.EOF {
        break
    }
    if err != nil {
        return fmt.Errorf("read failed: %v", err)
    }
}
```

5.5 Function Values

Functions are first-class values in GO:like other values, function values have types, and they may be assigned to variables or passed to or returned from functions. ***A function values may be called like any other function.***

> similar with function pointer or lambda?

***The zero value of a function type is nil.***  

```go
var f func(int) int
if f != nil {
    f(3)
}
```

***But they are not comparable***

```go
Printf("%*s</%s>\n", depth*2, "", n.Data)
```

The * adverb in %*s prints a string padded with a variable number of spaces.

## 5.6 Anonymous Functions

Named functions can be declared only ate the package level, but we can use a function literal to denote a function value within any expression.

```go
//Example
func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}
```

> ***function*** is returned

These hidden variable references are why we classify functions as reference types and why function values are not comparable. ***Function values like these are implemented using a technique called closures, and Go programmers often use this term for function values.***

When an anonymous function requires recursion, ***we must first declare a variable, and then assign the anonymous function to that variable.***

```go
visitAll = func(items []string) {
    for _, item := range items {
        if !seen[item] {
            seen[item] = true
            visitAll(m[item])
            order = append(order, item)
        }
    }
}
//visitAll := func()
//Worng!!!!
```

#### 5.6.1 Caveat: Capturing iteration Variables

We’ll look at a pitfall of Go’s lexical scope rules that can cause surprising results.

```go
//Right version
var rmdirs []func()
for _, d := range tenpDirs)_ {
    dir := d
    os.MkdirAll(dir, 0755)
    rmdirs = append(rmdirs, func() {
        os.RemoveAll(dir)
    })
    //... do some work...
    for _, rmdir := range rmdirs {
        rmdirs()
    }    
}

//Wrong version
var rmdirs []finc()
for _, dir := range tempDirs() {
    os.MkdirAll(dir, 0755)
    rmdirs = append(rmdirs, func() {
        os.RemoveAll(dir)
    })
}
//아니... 이거는 당연히 안되져;;
```

*dir* is defined. The reason is a consequence of the scope rules for loop variables.

The value of dir is updated in successive iterations, so by the time the cleanup functions are called, ***the dir variable has been updated several times by the now-completed for loop.*** Thus dir holds the value from the final iteration, and consequently all calls to os.RemoveAll will attempt to remove the same director y.



## 5.7 Variadic Functions

A variadic function is one that can be called with varying numbers of arguments. The most familiar examples are fmt.Printf and its variants.

***Printf requires one fixed argument at the beginning, then accepts any number of subsequent arguments.***

To declare a variadic function, the type of the final parameter is preceded by an ellipsis, ‘‘...’’, which indicates that the function may be called with any number of arguments of this type.

```go
func sum(cals ...int) int {
    total := 0
    for _, val := reange vals {
        total += val
    }
    return total
}
values := []int{1,2,3,4}
fmt.Println(sum(values...)) // "10"
```

Although the ...int parameter behaves like a slice within the function body, the type of a variadic function is distinct from the type of a function with an ordinary slice parameter.

```go
func f(...int) {}
func g([]int) {}
fmt.Printf("%T\n", f) // "func(...int)"
fmt.Printf("%T\n", g) // "func([]int)"
```



## 5.8 Deferred Function Calls

```go
package main

import (
	"CH5/outline2"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		resp.Body.Close()
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Parsing %s as HTML: %v", url, err)
	}
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	outline2.ForEachNode(doc, visitNode, nil)
	return nil
}
func main() {
	title(os.Args[1])
}

/*
C:\github-repository\Golang-practice\The-Go-programming-language>title1.exe http://gopl.io
The Go Programming Language
*/
```

​	resp.Body.Close() is duplicated -> Errors can be happened

***A defer statement is often used with paired operations like open and close***, connect and disconnect, or lock and unlock to ensure that resources are released in all cases, no matter how complex the control flow.

***The right place for a defer statement that releases a resource is immediately after the resource has been successfully acquired.***

```go
//More elegance way
package main

import (
	"CH5/outline2"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("Parsing %s as HTML : %v", url, err)
	}

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	outline2.ForEachNode(doc, visitNode, nil)
	return nil
}
func main() {
	title(os.Args[1])
}
```

***The defer statement can also be used to pair ‘‘on entry’’ and ‘‘on exit’’ actions when debugging a complex function.***

By deferring a call to the returned function in this way, we can instrument the entry point and all exit points of a function in a single statement and even pass values, like the start time, be ween the two actions.

defer red functions aren't executed until the very end of a function’s execution, a defer statement in a loop deserves extra scrutiny.

> scrutiny : 정밀 조사

The code below could run out of file descriptors since no file will be closed until all files have been processed:

```go
for _, filename := range filenames {
    f, err := os.Open(filename)
    if err != nil {
    return err
}
defer f.Close() // NOTE: risky; could run out of file descriptors
// ...process f...
}
```

## 5.9 Panic

Go’s type system catches many mistakes at compile time, but others, like an out-of-bounds array access or nil pointer dereference, require checks at run time. When the Go runtime detects these mistakes, it ***panics***.

During a typical panic, normal execution stops, all defer red function cal ls in that goroutine are executed, and the program crashes with a log message.

Not all panics come from the runtime. The built-in panic function may be called directly;

```go
func Reset(x *Buffer){
    if x == nil {
        panic("x is nil")
    }
    x.elements = nil
}
```

Although Go’s panic mechanism resembles exceptions in other languages, the situations in which panic is used are quite different.

Since a panic causes the program to crash, it is generally used for grave errors, such as a logical inconsistency in the program

When a panic occurs, all deferred functions are run in reverse order, starting with those of the topmost function on the stack and proceeding up to main, as the program below demonstrates:

```go
func main() {
    f(3)
}
func f(x int) {
    fmt.Printf("f(%d)\n", x + 0/x)
    defer fmt.Printf("defer %d\n",x)
    f(x - 1)
}
// Result
f(3)
f(2)
f(1)
defer 1
defer 2
defer 3
```

## 5.10 Recover

Giving up is usually the rig ht response to a panic, but not always. It might be possible to recover in some way, or at least clean up the mess before quitting.

If the built-in recover function is called within a deferred function and the function containing the defer statement is panicking, recover ends the current state of panic and returns the panic value.

> [For-more-information](https://blog.golang.org/defer-panic-and-recover)



```go
//It returns 2
func c() (i int) {
    defer func() { i++ }()
    return 1
}
```

