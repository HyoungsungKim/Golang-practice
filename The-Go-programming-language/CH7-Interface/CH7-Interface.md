# CH7 Interface

By generalizing, interfaces let us write functions that are more flexible and adaptable because they are not tied to the details of one particular implementation.

There’s no need to declare all the interfaces that a given concrete type satisfies; simply possessing the necessary methods is enough.

## 7.1 Interface as Contents

All the types we've looked at so far have been concrete types.

There is another kind of type in Go called an interface type. An interface is an abstract type.It doesn't expose the representation or internal structure of its values, or the set of basic operations they support; ***it reveals only some of their methods.***

***When you have a value of an interface type, you know nothing about what it is;*** you know only what it can do, or more precisely, what behaviors are provided by its methods.

Throughout the book, we've been using two similar functions for string formatting: fmt.Printf, which writes the result to the standard output (a file), and fmt.Sprintf, which returns the result as a string. ***It would be unfortunate if the hard part, formatting the result, had to be duplicated because of these superficial differences in how the result is used.*** Thanks to interfaces, it does not. Both of these functions are, in effect, wrappers around a third function, fmt.Fprintf, that is agnostic about what happens to the result it computes:

> Printf와 Sprintf는 비슷하지만 각각의 구현을 완전히 따로 하지 않음(인터페이스 덕분에)

```go
package fmt
func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
func Printf(format string, args ...interface{}) (int, error) {
    return Fprintf(os.Stdout, format, args...)
}
func Sprintf(format string, args ...interface{}) string{
    var buf bytes.Buffer
    Fprintf(&buf, format, args...)
    return buf.String()
}
```

In the Printf case, the argument, os.Stdout, is an *os.File.*

In the Sprintf case, however, the argument is not a file, though it superficially resembles one: &buf is a pointer to a memory buffer to which bytes can be written. The first parameter of Fprintf is not a file either. ***It's an io.Writer, which is an interface type with the following declaration:***

```go
package io
type Writer interface {
    Writer(p []byte) (n int, err error)
}
```

***The io.Writer interface defines the contract between Fprintf and its callers.***

On the one hand, the contract requires that the caller provide a value of a concrete type like *os.File or *bytes.Buffer that has a method called Write with the appropriate signature and behavior.

***On the other hand, the contract guarantees that Fprintf will do its job given any value that satisfies the io.Writer interface. Fprintf may not assume that it is writing to a file or to memory, only that it can call Write.***

***Because fmt.Fprintf assumes nothing about the representation of the value and relies only on the behaviors guaranteed by the io.Writer contract, we can safely pass a value of any concrete type that satisfies io.Writer as the first argument to fmt.Fprintf.*** 

This freedom to substitute one type for another that satisfies the same interface is called substitutability, and is a hallmark of object-oriented programming.

> substitutability : 대용 가능성
>
> hallmark : (전형적인)특징, (귀금속의)품질 보증 마크

Let's test this out using a new type.The Write method of the *Byte Counter type below merely counts the bytes written to it before discarding them.

```go
//Implement Writer interface by using ByteCounter type
type ByteCounter int
//Satisfying contract
func(c *ByteCounter) Write(p []byte) (int, error) {
    *c += ByteCounter(len(p))
    return len(p), nil
}
```

Since `*ByteCounter` satisfies the io.Writer contract, ***we can pass it to Fprintf,*** which does its string formatting oblivious to this change; the ByteCounter correctly accumulates the length of the result.

```go
var c ByteCounter
c.Write([]byte("hello"))
fmt.Println(c)	// "5", = len("hello")
/*
func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
func Printf(format string, args ...interface{}) (int, error) {
    return Fprintf(os.Stdout, format, args...)
}
*/
c = 0
var name = "Dolly"
fmt.Fprintf(&c, "hello, %s", name)
fmt.Printf(c)	//"12" = len("hello, Dolly")
```

> func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)***
>
> ***fmt.Printf(c)가 12인 이유 : fmt.Fprintf에서 &c가 w io.Writer로 사용되면서 c의 io.Writer 호출***
>
> 와... 기가 막히네;;;

Besides io.Writer, there is another interface of great importance to the fmt package. ***Fprintf and Fprintln provide a way for types to control how their values are printed.*** In Section2.5, we defined a String method for the Celsius type so that temperatures would print as "100°C",and in Section 6.5 we equipped *IntSet with a String method so that sets would be rendered using traditional set notation like "{1 2 3}".

```go
func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }
c:=FToC(212.0)
fmt.Println(c.String()) // "100°C"
```

## 7.3 Interface Satisfaction

A type satisfies an interface if it possesses all the methods the interface requires.

Go programmers often say that a concrete type ‘‘is a’’ particular interface type,meaning that it satisfies the interface. 

```go
var w io.Writer
w = os.Student			//OK
w = new(bytes.Buffer)	 //Ok
w = time.Second			//error

var rwc io.ReadWriterCloser
rwc = os.Stdout			//OK
rwc = new(bytes.Buffer)	//error
```

The type interface{}, which is called the empty interface type, is indispensable. Because the empty interface type places no demands on the types that satisfy it, we can assign any value to the empty interface.

```go
var any interface{}
any = true
any = 12.34
any = "hello"
any = map[string]int{"one": 1}
any = new(bytes.Buffer)
```

> Similar with void* in C/C++

Since interface satisfaction depends only on the methods of the two types involved, there is no need to declare the relationship between a concrete type and the interfaces it satisfies.

```go
type Artifact interface {
    Title() String
    Creators() []string
    Created() time.time
}
type Text interface {
    Pages() int
    Words() int
    PageSize() int
}
type Audio interface {
    Stream() (io.ReadCloser, error)
    RunningTime() time.Duration
    Format() string
}
type Video interface {
    Stream() (io.ReadCloser, error)
    RunningTime() time.Duration
    Format() string
    Resolution() (x, y int)
}
```

## 7.4 Parsing Flags with flag.Value

Because duration-valued flags are so useful, this feature is built in to the flag package, but it’s easy to define new flag notations for our own data types. We need only define a type that satisfies the flag.Value interface.

```go
package flag
type Value interface {
    String() string
    Set(string) error
}
```

## 7.5 Interface Values

Conceptually, a value of an interface type,or interface value, has two components, a *concrete type* and a *value of that type.* These are called the interface’s dynamic type and dynamic value. 

```go
var w io.Writer
w = os.Stdout
w = new(bytes.Buffer)
w = nil
```

An interface value is described as nil or non-nil based on its dynamic type, so this is a nil interface value.You can test whether an interface value is nil using w==nil or w!=nil. Calling any method of a nil interface value causes a panic:

```go
w.Write([]byte("hello"))	//panic : nil pointer dereference
```

In general,we cannot know at compile time what the dynamic type of an interface value will be,so a call through an interface must use dynamic dispatch.

An interface value can hold arbitrarily large dynamic values. For example, the time.Time type,which represents an instant in time, is a struct type with several unexported fields.If we create an interface value from it.

```go
var x interface{} = time.now()
```

Other types are either safely comparable (like basic types and pointers)or not comparable at all(like slices, maps, and functions), ***but when comparing interface values or aggregate types that contain interface values, we must be aware of the potential for a panic.*** A similar risk exists when using interfaces as map keys or switch operands. ***Only compare interface values if you are certain that they contain dynamic values of comparable types***

```go
var w io.Wirter
fmt.Println(x == x)	// Panic
```

## 7.5.1 Caveat: An Interface Containing a Nil pointer is Non-Nil

***A nil interface value, which contains no value at all, is not the same as an interface value containing a pointer that happens to be nil.*** This subtle distinction creates a trap into which every Go programmer has stumbled.

```go
Const debug = true
func main() {
    var buf *bytes.Buffer
    if debug {
        buf = new(bytes.Buffer)
    }
    f(buf)
    if debug {
        //...
    }
}
func f(out io.Writer) {
    if out != nil {
        out.Writer([]byte("done!\n"))
    }
}
```

## 7.6 Sorting with sort.Interface

Fortunately, the sort package provides in-place sorting of any sequence according to any ordering function. Its design is rather unusual. In many languages, the sorting algorithm is associated with the sequence data type, while the ordering function is associated with the type of the elements. 

it uses an interface, sort.Interface, to specify the contract between the generic sort algorithm and each sequence type that may be sorted.

To sort any sequence,we need to define a type that implements these three methods, then apply sort. Sort to an instance of that type.

```go
type StringSlice []string
func (p StringSlice) Len() int				{ return len(p)}
func (p StringSliice) Less(i, j int) bool	 { return p[i] < p [j]}
func (p StringSlice) Swap(i j int)			{ p[i], p[j] = p[j], p[i]}

sort.Sort(StringSlice(names))
```

***It will run faster if each element is a pointer***

## 7.7 The http.Handler Interface

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := database{"shoes": 50, "sockes": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}
//In browser, localhost:8000
//shoes: $50.00
//sockes: $5.00
```

Obviously we could keep adding cases to ServeHTTP, but in a realistic application, it’s convenient to define the logic for each case in a separate function or method. Furthermore, related URLs may need similar logic; several image files may have URLs of the for m /images/*.png, for instance. For these reasons, net/http provides ServeMux, a request multiplexer, to simplify the association between URLs and handlers. A ServeMux aggregates a collection of http.Handlers into a single http.Handler. Again, we see that different types satisfying the same interface are substitutable:the web server can dispatch requests to any http.Handler, regardless of which concrete type is behind it.

> aggregates : 합계, 총액

The expression http.HandlerFunc(db.list) is a conversion, not a function call, since http.HandlerFunc is a type.

## 7.8 The error Interface

```go
type error interface {
    Error() string
}
```

> inadvertent : 고의가 아닌

The underlying type of error String is a struct,not a string , to protect its representation from inadvertent(or premeditated) updates. And the reason that the pointer type *errorString, not errorString alone, satisfies the error interface is so that

***every call to New allocates a distinct error instance that is equal to no other.***

```go
fmt.Println(errors.New("EOF") == errors.New("EOF"))	//false
```

```go
var errors = [...]string{
    1: "operation not permitted",	 // EPERM
    2: "no such file or directory",	 // ENOENT
    3: "no such process",			// ESRCH 
    // ...
}
func (e Errno) Error() string {
	if 0 <= int(e) && int(e) < len(errors) {
        return errors[e]
	}
	return fmt.Sprintf("errno %d", e)
}

var err error = syscall.Errno(2)
fmt.Println(err.Error()) // "no such file or directory"
fmt.Println(err)		// "no such file or directory"
```

## 7.9 Example: Expression Evaluator

> How can i get Parse function...

## 7.10 Type Assertion

A type assertion is an operation applied to an interface value. 

```go
//Syntactically, it looks like
x.(T)
```

There are two possibilities. ***First, if the asserted type T is a concrete type,*** then the type assertion checks whether x’s dynamic type is identical to T. ***A type assertion to a concrete type extracts the concrete value from its operand. If the check fails, then the operation panics.***

> assert : 주장하다

***Second, if instead the asserted type T is an interface type,*** then the type assertion checks whether x’s dynamic type satisfies T. A type assertion to an interface type changes the type of the expression, making a different(and usually larger) set of methods accessible, but ***it preserves the dynamic type and value components inside the interface value.***

## 7.11 Discriminating Errors with Type Assertions

Three kinds of failure often must be handled differently:

- file already exists(for create operations)
- file not found (for read operations)
- permission denied.

The os package provides these three helper functions to classify the failure indicated by a given error value"

```go
package os
func IsExist(err error) bool
func IsNotExist(err error) bool
func IsPermission(err error) bool
```

## 7.12 Querying Behaviors with Interface Type Assertions

Can we avoid allocating memory here?

```go
func writeString(w io.Writer, s string) (n int, err error) {
    type stringWriter interface {
        WriteString(string) (n int, err error)
    }
    if sw, ok := w.(stringWriter); ok {
        return sw.WriteString(s)
    }
    return w.Write([]bytes(s))
}

func writeHeader(w io.Writer, contentType string) error {
    if _, err := writeString(w, "Content-Type: "); err != nil {
        return err
    }
    if _, err := writeString(w, contentType); err != nil {
        return err
    }
    //...
}
```



## 7.13 Type Switches

Interfaces are used in two distinct styles. 

- In the first style, exemplified by io.Reader, io.Writer, fmt.Stringer, sort.Interface, http.Handler, and error, an interface’s methods express the similarities of the concrete types ***that satisfy the interface but hide the representation details and intrinsic operations of those concrete types.***

> intrinsic  : 고유한, 본질적인
>
> discriminate  : 식별하다
>
> exploit : 이용하다, 착취하다.

- The second style exploits the ability of an interface value to hold values of a variety of concrete types and considers the interface to be the union of those types. Type assertions are used to discriminate among these types dynamically and treat each case differently. ***In this style, the emphasis is on the concrete types that satisfy the interface, not on the interface’s methods(if indeed it has any), and there is no hiding of information.***

## 7.15 A Few Words of Advice

When designing a new package, novice Go programmers often start by creating a set of interfaces and only later define the concrete types that satisfy them. This approach results in many interfaces, each of which has only a single implementation. ***Don’t do that.*** Such interfaces are unnecessary abstractions; they also have a runtime cost. You can restrict which methods of a type or fields of a struct are visible outside a package using the export mechanism. ***Interfaces are only needed when there are two or more concrete types that must be dealt with in a uniform way.***



