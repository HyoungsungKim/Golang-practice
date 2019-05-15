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

