# Go Syntax sugar and tricks

I go even if i cannot go

but... hey my life, are you going well? T_T

## Embedding

Based on [Medium](https://travix.io/type-embedding-in-go-ba40dd4264df) &  [Effective Go](https://golang.org/doc/effective_go.html#embedding) 

Embedding is used for inheritance in GO

---

### Medium

#### Basic concept

- The theory behind embedding is pretty straightforward
  - By including a type as a nameless parameter within another type, the exported parameters and methods defined on the embedded type are accessible through the embedding type.
  - ***The compiler decides on this by using a technique called “promotion”:*** the exported properties and methods of the embedded type ***are promoted to the embedding type.***

Example)

[go playground](https://play.golang.org/p/EBUmBfaCHEC)

```go
type Ball struct {
    Radius int
    Material string
}

//inherit
type Football struct {
    Ball
}

func (b Ball) Bounce() {
    fmt.Printf("Bouncing ball %+v\n", b)
}

func main() {
    fb := Football{Ball{Radius : 5, Material : "leather"}}
    fmt.Printf("fb = %+v\n", fb)
    //fb = {Ball:{Radius:0 Materual:}}
    fb.Bounce()
    //Bouncing ball {Radius:0 Materal:}
	fb.Ball.Bounce()	//->Same way
    //Bouncing ball {Radius:0 Materal:}
}
```

#### Embedding interfaces(1)

[go playground](https://play.golang.org/p/gNkUSwo6839)

***If the embedded type implements a particular interface,*** then that too is accessible through the embedding type.

```go
type Bouncer interface {
    Bounce()
}

func BounceIt(b Bouncer) {
    b.Bounce()
}
//Call functing using embedding
BounceIt(fb)
```

#### Embedding Pointers

```go
type Football struct {
    *Ball
}

func (b *Ball) Bounce() {
    fmt.Printf("Bouncing ball %+v\n", b)
}
```

#### Embedding interfaces(2)

[Go playground](https://play.golang.org/p/KFBGxR2N8hJ)

```go
package main

import (
	"fmt"
)

type Ball struct {
	Radius int
	Material string
}

type Bouncer interface {
	Bounce()
}

type Football struct {
	//Interface embedding
	Bouncer
}

func (b *Ball) Bounce() {
    fmt.Printf("Bouncing ball %+v\n", b)
}

func main() {
    //fb := Football{Ball{Radius : 5, Material : "leather"}} -> complie error
    //Cause receiver of Bounce is pointer
    fb := Football{&Ball{Radius : 5, Material : "leather"}}
	fb.Bounce()
    //Bouncing ball &{Radius:5 Material:leather}
}
```

#### Warning

- The embedded struct has no access to the embedding struct

[go playground](https://play.golang.org/p/LkWxSFIpnh4)

```go
package main

import (
	"fmt"
)

type Ball struct {
	Radius int
	Material string
}

type Football struct {
	Ball
	Radius int
}

func (b Ball)Bounce() {
	fmt.Printf("Radius : %d\n", b.Radius)
}

func main() {
    //Football 타입으로 Bounce 구현 없이 ball의 함수 호출 할 수 있음
    //embedding의 장점
    //만약 Football만의 함수 가지고 싶다면 구현 필요
	fb := Football{Ball{Radius : 5, Material : "leather"}, 7}
	fmt.Printf("fb = %+v\n", fb)
	fb.Bounce()
    fb.Ball.Bounce()
 	//output
    //Radius : 5
	//Radius : 5
    //Need concrete implementation to call 7
    /*
    	func (fb Football) Bounce() {
    		fmt.Printf("Radius : %d\n", b.Radius)
    	}
    */
}
```

### Effective Go

Go does not provide the typical, type-driven notion of subclassing, but it does have the ability to “borrow” pieces of an
implementation by *embedding* types within a struct or interface.

Example)

***Only interfaces can be embedded within interfaces.***

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

//Interface embedding
type ReadWriter interface {
    Reader
    Writer
}
```

The same basic idea applies to structs, but with more far-reaching implications.

The `bufio` package has ***two struct types,*** `bufio.Reader` and `bufio.Writer`, each of which of course implements the analogous interfaces from package `io`. ***And `bufio` also implements a buffered reader/writer, which it does by combining a reader and a writer into one struct using embedding:*** it lists the types within the struct but does not give them field names.

```go
//It implements io.ReadWriter
type ReadWriter struct {
    *Reader // *bufio.Reader
    *Writer	// *bufio.Writer
}

//Same declaration
type ReadWriter struct {
    reader *Reader
    writer *Writer
}
```

forwarding methods

```go
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
    return rw.Reader.Read(p)
}
```

***The methods of embedded types come along for free,*** which means that `bufio.ReadWriter` not only has the methods of `bufio.Reader` and `bufio.Writer`, it also satisfies all three interfaces: `io.Reader`, `io.Writer`, and `io.ReadWriter`.

> 상속(Inheritance)와 비슷한 성질 가지게 됨!

***There's an important way in which embedding differs from subclassing.***

- When we embed a type, ***the methods of that type become methods of the outer type, but when they are invoked the receiver of the method is the inner type, not the outer one.***

>***바깥에 있는 type으로 안에 embedded된 type의 함수 호출 가능. 바깥에 있는 type으로 호출하지만 실제로 호출 되는건 embedded된 type의 함수***

Example)

```go
type Football struct {
	Ball
	Radius int
}

func (b Ball)Bounce() {
	fmt.Printf("Radius : %d\n", b.Radius)
}

func main() {
	fb := Football{Ball{Radius : 5, Material : "leather"}, 7}
	fmt.Printf("fb = %+v\n", fb)
    fb.Bounce()	//Output : Radius :5
	fb.Ball.Bounce() ////Output : Radius :5
}
```

- In our example, when the `Read` method of a `bufio.ReadWriter` is invoked, ***it has exactly the same effect as the forwarding method written out above;*** the receiver is the `reader` field of the `ReadWriter`, not the `ReadWriter` itself.

---

## Struct and Interface type

[Effective go](https://golang.org/ref/spec#Struct_types)

### Struct Type

A struct is a sequence of named elements, called fields, each of which has a name and a type. ***Field names may be specified explicitly (IdentifierList) or implicitly (EmbeddedField).***

A field or method `f` of an embedded field in a struct `x` is called *promoted* if `x.f` is a legal selector that denotes that field or method `f`.

> Embedding에서 공부했었음

### Interface Type

An interface type specifies a method set called its *interface*. A variable of interface type can store a value of any type with a method set that is any superset of the interface. Such a type is said to *implement the interface*. The value of an uninitialized variable of interface type is `nil`.

A type implements any interface comprising any subset of its methods and may therefore implement several distinct interfaces. For instance, all types implement the ***empty interface***:

```go
interface{}
```

An interface `T` may use a (possibly qualified) interface type name `E` in place of a method specification. This is called
***embedding interface*** `E` in `T`; it adds all (exported and non-exported) methods of `E` to the interface `T`.

```go
type ReadWriter interface {
    Read(b Buffer) bool
    Write(b buffer) bool
}

type File interface {
    ReadWriter	// same as adding the methods of ReadWriter
    Locker		// same as adding the methods of Locker
    Close()
}

type LockedFile interface {
    Locker
    File		// illegal : Lock, Unlock not unique
    Lock()		// illegal : Lock, not unique
}
```

### Empty struct

[medium](https://medium.com/@l.peppoloni/how-to-improve-your-go-code-with-empty-structs-3bd0c66bc531)

```go
struct{}
```

The cool thing about an empty structure is that it occupies zero bytes of storage.

- *What can I use an empty struct for, if it has no fields?*
- Basically an empty struct can be used every time you are only interested in a property of a data structure rather than its value

#### Semaphores and tokens

Making semaphore

```go
sem := make(chan bool, numberOfSemaphores)
sem <-true
//or
sem := make(chan int, numberOfSemaphores)
sem <- 1
```

Using struct

```go
//Declare array of empty struct
sem := make(chan struct{}, numberOfSemaphores)
sem <- struct{}{}
```

>Semaphore에서는 signal을 보내는거지 channel 안의 내용은 중요하지 않음
>
>따라서 empty struct{}를 통해 효율적인(?) 프로그래밍 가능

