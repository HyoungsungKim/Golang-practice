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

