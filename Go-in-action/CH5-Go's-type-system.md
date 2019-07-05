# CH5 Go's type system

## 5.1 User-defined types

The most common way is to use the keyword struct , which allows you to create a composite type.

```go
type user struct {
    name string
    email string
    ext int
    privilieged bool
}
var bill user
/*
bill = {
	name = ""
	email = ""
	ext = 0
	privilieged = false
}
*/

Lisa := user{
    name: "Lisa",
    email: "Hello@world.com",
    ext: 123,
    privileged: true,
}
//or
Lisa := {"Lisa", "Hello@world.com", 123, true}

type admin struct {
    person user
    level string
}

fred := admin {
    person: user{
        name: "Lisa",
        email: "Hello@world.com",
        ext: 123,
        priviliege: true
    }
    level: "super",
}
```

A second way to declare a user-defined type is by taking an existing type and using it as the type specification for the new type.

```go
package main

type Duration int64

func main() {
    var dur Duration
    dur = int64(1000)
}
//Compile Error
```

## 5.2 Methods

Methods provide a way to add behavior to user-defined types. Methods are really functions that contain an extra parameter that's declared between the keyword func and the function name

```go
package main

import "fmt"

type user struct {
    name string
    email string
}

func (u user)notify(){
    fmt.Printf("name : %s, email : %s \n", u.name, u.email)
}

func (u &user)changeEmail(email string){
    u.email = email
}

func main() {
    bill := user{
        name: "bill",
        email: "bill@email.com",
    }    
    bill.notify()
    
    lisa := &user{
        name: "lisa",
        email: "lisa@email.com",
    }
    list.notify()
    
    bill.changeEmail("bill@newEmail.com")
    bill.notify()
    
    lisa.changeEmail("lisa@newEmail.com")
    list.notify()
}
```

The parameter between the keyword func and the function name is called a receiver and binds the function to
the specified type. ***When a function has a `receiver`, that function is called a `method`.***

There are two types of receivers in Go:

- value receivers
- pointer receivers

***Value receivers operate on a copy of the value used to make the method call and pointer receivers operate on the actual value.***

You can also call methods that are declared with a pointer receiver using a value.

## 5.3 The Nature of types

The idea is to not focus on what the method is doing with the value, but to focus on what the nature of the value is.

### 5.3.1 Built-in types

string

```go
func Trim(s string, cutset string) string {
    if s == "" || cutset == "" {
        return s
    }
    return TrimFunc(s, makeCutsetFunc(cutset))
}
```

### 5.3.2 Reference types

Reference types in Go are the set of slice, map, channel, interface, and function types. When you declare a variable from one of these types, the value that's created is called a `header` value. You never share reference type values because the header value is designed to be copied.

### 5.3.3 Struct types

Struct types can represent data values that could have either a primitive or nonprimitive nature.

## 5.4 Interface

### 5.4.1 Standard library

```go
import (
	"fmt"
    "io"
    "net/http"
    "os"
)

func init() {
    if len(os.Args) != 2 {
        fmt.Println("Usage : ./example2 <url>")
        os.Exit(-1)
    }
}

func main() {
    r, err := http.Get(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }
    
    io.Copy(os.Stdout, r.Body)
    if err := r.Body.Close(); err != nil {
        fmt.Println(err)
    }
}
```

In a few lines of code, we have a curl program by leveraging two functions that work with interface values. 

```go
package main
import (
	"bytes"
    "fmt"
    "io"
    "os"
)

func main() {
    var b bytes.Buffer
    b.Write([]byte("Hello"))
    fmt.Fprintf(&b, "World!")
    io.Copy(os.Stdout, &b)
}
```

### 5.4.2 Implementation

Interfaces are types that just declare behavior. ***This behavior is never implemented by the interface type directly but instead by user-defined types via methods.*** 

- When a user-defined type implements the set of methods declared by an interface type, values of the user-defined type can be assigned to values of the interface type.

If a method call is made against an interface value, the equivalent method for the stored user-defined value is executed. There are rules around whether values or pointers of a user-defined type satisfy the implementation of an interface. Not all values are created equal. These rules come from the specification under the section called method sets. Before you begin to investigate the details of method sets, it helps to understand what interface type values look like and how user-defined type values are stored inside them.

Interface values are two-word data structures.

- The first word contains a pointer to an internal table called an iTable, which contains type information about the stored value. The iTable contains the type of value that has been stored and a list of methods associated with the value.
- The second word is a pointer to the stored value. The combination of type information and pointer binds the relationship between the two values.

### 5.4.3 Method sets

```go
package main

import (
	"fmt"	
)

type notifier interface {
    notify()
}

type user struct {
    name string
    email string
}

func (u *user) notify() {
    fmt.Printf("Sending user email to %s<%s>\n", u.name, u.email)
}

func main() {
    u := user{"Bill", "bill@email.com"}
    sendNotification(u)
}

func sendNotification(n notifier) {
    n.notify()
}

/*
./prog.go:22:18: cannot use u (type user) as type notifier in argument to sendNotification:
	user does not implement notifier (notify method has pointer receiver)
*/
```

To understand why values of type user don't implement the interface when an interface is implemented with a pointer receiver, you need to understand what method sets are.

Method sets define the set of methods that are associated with values or pointers of a given type. The type of receiver used will determine whether a method is associated with a value, pointer, or both.

Let's start with explaining the rules for method sets as it's documented by the Go specification.

| Values | Methods Receivers |
| ------ | ----------------- |
| T      | (t T)             |
| *T     | (t T) and (t *T)  |



| Methods Receivers | values   |
| ----------------- | -------- |
| (t T)             | T and *T |
| (t *T)            | *T       |

The question now is why the restriction? The answer comes from the fact that it’s not always possible to get the address of a value.

```go
package main
import "fmt"
type duration int

func (d *duration) pretty() string {
    return fmt.Sprintf("Duration: %d", *d)
}

func main() {
    duration(42).pretty()
}
/*
./prog.go:11:14: cannot call pointer method on duration(42)
./prog.go:11:14: cannot take the address of duration(42)
*/
```

This shows that it's not always possible to get the address of a value.

### 5.4.4 Polymorphism

```go
package main

import "fmt"

type notifier interface {
    notify()
}

type user struct {
    name string
    email string
}

func (u *user) notify() {
    fmt.Printf("Sending user email to %s<%s>\n", u.name, u.email)
}

type admin struct {
    name string
    email string
}

func (a *admin) notify() {
    fmt.Printf("Sending admin email to %s<%s>\n", a.name, a.email)
}

func main() {
    bill := user{"Bill", "bill@email.com"}
    sendNotification(&bill)
    
    lisa := admin("Lisa", "lisa@email.com")
    sendNotification(&lisa)
}

func sendNotification(n notifier) {
    n.notify()
}
```

## 5.5 Type embedding

Go allows you to take existing types and both extend and change their behavior. This capability is important for code reuse and for changing the behavior of an existing type to suit a new need. ***This is accomplished through type embedding.*** It works by taking an existing type and declaring that type within the declaration of a new struct type. The type that is embedded is then called an inner type of the new outer type.

- ***Through inner type promotion, identifiers from the inner type are promoted up to the outer type.*** These promoted identifiers become part of the outer type as if they were declared explicitly by the type itself. The outer type is then composed of everything the inner type contains, and new fields and methods can be added.

- The outer type can also declare the same identifiers as the inner type and override any fields or methods it needs to. This is how an existing type can be both extended and changed.

```go
package main
import "fmt"

type user struct {
    name string
    email string
}

func (u *user)notify() {
    fmt.Printf("Name : %s, Email : %s", u.name, u.email)
}

type admin struct {
    user	//Embedding Type
    level string
}

func main() {
    ad := admin{
        user: user{
            name: "john",
            email: "john@email.com",
        },
        level : "super",
    }
    
    ad.user.notify()
    ad.notify()
    //thanks to inner type promotion, the notify method can also be accessed directly from the ad variable.
}
/*
Result
Name : john, Email : john@email.com
Name : john, Email : john@email.com
*/
```

Once we embed the user type inside of admin , we can say that ***`user` is an inner type*** of the ***outer type `admin`*** . The concept of having an inner and outer type makes it easier to understand the relationship between the two.

```go
package main
import "fmt"

type user struct {
    name string
    email string
}

func (u *user)notify() {
    fmt.Printf("Name : %s, Email : %s\n", u.name, u.email)
}

type admin struct {
    user	//Embedding Type
    level string
}

type notifies interface {
    notify()
} 

func main() {
    ad := admin{
        user: user{
            name: "john",
            email: "john@email.com",
        },
        level : "super",
    }
    
    ad.user.notify()
    ad.notify()
    //thanks to inner type promotion, the notify method can also be accessed directly from the ad variable.
    sendNotification(&ad)
}


func sendNotification(n notifies) {
    n.notify()
}
/*
result
Name : john, Email : john@email.com
Name : john, Email : john@email.com
Name : john, Email : john@email.com
*/
```

***Thanks to inner type promotion, the implementation of the interface by the inner type has been promoted up to the outer type.*** That means the outer type now implements the interface, thanks to the inner type's implementation. When we run this sample program, we get the following output.

## 5.6 Exporting and unexporting identifiers

The ability to apply visibility rules to the identifiers you declare is critical for good API design. Go supports the exporting and unexporting of identifiers from a package to provide this functionality.

Sometimes, ***you may not want identifiers such as types, functions, or methods to be a part of the public API for a package.*** In these cases, you need a way to declare those identifiers so they're unknown outside the package. You need to declare them to be unexported.

```go
package counters
type alertCounter int
```

```go
package main
import (
    "fmt"
    "counters"
)

func main() {
    counters := counters.alertCounter(10)
    fmt.Printf("Counter: %d\n", counter)
}
```

When an identifier starts with a lowercase letter, the identifier is unexported or unknown to code outside the package. When an identifier starts with an uppercase letter, it's exported or known to code outside the package. Let's look at the code that imports this package.