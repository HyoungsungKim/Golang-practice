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



