# CH7 Interface

By generalizing, interfaces let us write functions that are more flexible and adaptable because they are not tied to the details of one particular implementation.

There’s no need to declare all the interfaces that a given concrete type satisfies; simply possessing the necessary methods is enough.

## 7.1 Interface as Contents

All the types we've looked at so far have been concrete types.

There is another kind of type in Go called an interface type. An interface is an abstract type.It doesn't expose the representation or internal structure of its values, or the set of basic operations they support; ***it reveals only some of their methods.***

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

