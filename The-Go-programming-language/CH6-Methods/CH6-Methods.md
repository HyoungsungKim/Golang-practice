# CH6 Methods

object is simply a value or variable that has methods, and a ***method is a function associated with a particular type***. An object-oriented program is one that uses methods to express the properties and operations of each data structure so that clients need not access the object’s representation directly.

## 6.1 Method Declarations

A method is declared with a variant of the ordinary ***function declaration in which an extra parameter appears before the function name.***

```go
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.Xp.X, q.Yp.Y)
}
```

In Go, ***we don’t use a special name like this or self for the receiver***; we choose receiver names just as we would for any other parameter. In a method call, the receiver argument appears before the method name.

***The expression p.Distance is called a selector***, because it selects the appropriate Distance method for the receiver p of type Point. Selectors are also used to select fields of struct types, as in p.X.

Go is unlike many other object-oriented languages. It is often convenient to define additional behaviors for simple types such as numbers, strings, slices, maps, and sometimes even functions.

## 6.2 Methods with a Pointer Receiver

Because calling a function makes a copy of each argument value, if a function needs to update a variable, or if an argument is so large that we wish ***to avoid copying it, we must pass the address of the variable using a pointer.***

```go
func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}
//Case 1 :
r := &Point{1,2}
r.ScaleBy(2)
fmt.Println(*r)
//Case 2 :
p := Point{1,2}
pptr := &p
pptr.ScaleBy(2)
fmt.Println(p)
//Case 3 :
p := Point{1,2}
(&p).ScaleBy(2)
fmt.Println(p)
```

But the last two cases are ungainly.

the receiver argument is a variable of type T and the receiver parameter has type *T. The compiler implicitly takes the address of the variable:

```go
p.ScaleBy(2) // implicit (&p)
```

### 6.2.1 Nil is a Valid Receiver Value

Just as some functions allow nil pointers as arguments, so do some methods for their receiver, especially if nil is a meaningful zero value of the type, as with maps and slices.

```go
m := url.Values{"lang": {"en"}}
m.Add("item", "1")
m.Add("item", "2")

fmt.Println(m.Get("lang"))	//"en:"
fmt.Println(m.Get("q"))		//""
fmt.Println(m.Get("item"))	//"[1 2]"

m = nil
fmt.Println(m.Get("item"))	//""
m.Add("item", "3")

```

## 6.3 Composing Types by Struct Embedding

```go
package coloredpoint

import "image/color"
type Point struct{X, Y float64}
type ColoredPoint struct{
	Point
	Color color.RGBA
}

var cp ColoredPoint
cp.X = 1
fmt.Println(cp.Point.X)
cp.Point.Y = 2
fmt.Println(cp.Y)
```

As we saw in Section 4.4.3, embedding lets us take a syntactic shortcut to defining a ColoredPoint that contains all the fields of Point, plus some more.

## 6.4 Method Values and Expressions

Usually we select and cal l a met hod in the same expression, as in p.Distance(), but it’s possible to separate these two operations.

```go
p := Point{1, 2}
q := Point{4, 6}
distanceFromP := p.Distance // method value
fmt.Println(distanceFromP(q)) // "5"
var origin Point // {0, 0}
fmt.Println(distanceFromP(origin)) // "2.23606797749979", ;5

scaleP := p.ScaleBy // method value
scaleP(2) // p becomes (2, 4)
scaleP(3) // then (6, 12)
scaleP(10) // then (60, 120)
```

## 6.6 Encapsulation

A variable or met hod of an object is said to be encapsulated if it is inaccessible to clients of the object. Encapsulation, sometimes called information hiding, is a key aspect of object-oriented programming.

Go has only one mechanism to control the visibility of names: ***capitalize d identifiers are exported from the package in which they are defined, and uncapitalized names are not.***

As a consequence, to encapsulate an object, we must make it a struct. 

```go
type IntSet struct {
	words []uint64
}
//or
type IntSet []uint64
```

Another consequence of this name-based mechanism is that the unit of encapsulation is the package, not the type as in many other languages.

Encapsulation provides three benefits.

- First, because clients cannot directly modify the object’s variables, one need inspect fewer statements to understand the possible values of those variables.
- Second, hiding implementation details prevents clients from depending on things that might change, which gives the designer greater freedom to evolve the implementation without breaking API compatibility

Since Buffer is a struct type, this space takes the form of an extra field of type [64]byte with an uncapitalized name. When this field was added, because it was not exported, clients of Buffer outside the bytes package were unaware of any change except improved performance.

```go
type Buffer struct {
    buf []byte
    initial [64]byte
    /* ... */
}
```

- The third benefit of encapsulation, and in many cases the most important, is that it prevents clients from setting an object’s variables arbitrarily.

