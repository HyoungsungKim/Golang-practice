# Ch2 program Structure

In this chapter, we’ll go into more
det ail about the basic structural elements of a Go program. The example programs are intentionally
simple, so we can focus on the language without getting sidetracked by complicated
algorithms or data structures.



## 2.1 Names

***If the name begins with an upper-case letter, it is extorted, which means that it is visible and accessible outside of its own package and may be referred to by other parts of the program, as with Printf in the fmt package.*** Package names themselves are always in lower case. Stylistically, GO programmers use "camel case" when forming names by combining words, that is, interior capital letters are preferred over interior underscores.



## 2.2 Declarations

A declaration names a program entity and specifies some or all of its properties. There are four major kinds of declarations:var, const, type, and func.

```go
package main

import"fmt"

const boilingF = 212.0

func main() {
    var f = boilingF
    var c = (f - 32)* 5/ 9
    fmt.Print("boiling point = %g F or %g C \n", f, c)
}
```

***The constant boilingF is a package-level declaration(as is main)***, whereas the valuables f and c are local to the function main. ***The name of each package-level entity is visible not only throughout the source file that contains its declaration, but throughout all the files of the package.



## 2.3 Variables

A var declaration creates a variable of a particular type, attaches name to it, and sets its initial value.

```go
var name type = expression
//If expression is omitted then initialized value is the zero(0)
```

In Go there is no such thing as an uninitialized d variable.

### 2.3.1 Short Variable Declarations

Within a function, an alternate form called a short variable declaration may be used to declare and initialize local variables. 

```go
name := expression

i := 100
var boiling float64 = 100
var names []string
var err error
var p Point
```

***:= is a declaration, whereas = is an assignment***

```go
i, j = j ,i
//Swap values of i and j

f, err := os.Open(name)
if err != nil{
    return err
}

f.Close()
```

One subtle but important point: a short variable declaration does not necessarily declare all the variables on its left-hand side. *If some of them were already declared in the same lexical block , then the short variable declaration acts like an assignment to those variables.*

***A short variable declaration must declare at least one new variable***

```go
f, err := os.Open(infile)
f, err := os.Create(outfile)	// compile erroe: no new variable
```

### 2.3.2 Pointers

***A variable is a piece of storage containing a value.*** A pointer is the address of a variable. A pointer is thus the location at which a value is stored. ***not every value has am address, but every variable does.*** If a variable is declared var x int the expression ***&x ("address of x" )*** yields a pointer to an integer variable, that is, a value of type *int, which pronounced "pointer to int." 

```go
x := 1
p := &x
fmt.Println(*p)
*p = 2
fmt.Println(x)
```

***The zero value for a pointer of any type is nil.*** the test p != nil is true if p points to a variable. Pointers are comparable; ***two pointers are equal if and only if they point to the same variable or both are nil.***

```go
var x, y int
fmt.Println(&x == &x, &x == &y, &x == nil) //"true false false"
```

It is perfectly safe foe a function to return the address of a local variable.

```go
var p = f()
func f() *int {
    v := 1
    return &v
}

fmt.Println(f() == f()) // false!
```

***Because a pointer contains the address of a variable,*** passing a pointer argument to a function makes it possible for the function to update the variable that was indirectly passed.

```go
func incr(p *int) int{
    *p++
    return *p
}

v := 1
incr(&v)				// v is 2
fmt.Println(incr&v))	// v is 3
```

Pointers are key to flag package, which uses a program's command-line arguments to set the values of ***certain variables distributed throughout the program.*** To illustrate, this variation on the earlier echo command takes two optional flags: n causes echo to omit the trailing newline that would normally be printed, and -s sep causes it to separate the output arguments by the contents of the string sep instead of the default single space.

```go
package main

import (
	"fmt"
    "flag"
    "strings"
)

var n = flag.Bool("n", false, "omit trailing newline")
//The name of the flag "n", The variable's dafault value "false", and message "omit~"
// Message will be printed If the user provides an invalid argument, an invalid flag, or -h or -help.

var sep = flag.String("s", " ", "separator")

func main(){
    flag.Parse()
    fmt.Print(strings.Join(flag.Args(), *sep))
    if !*n{
        fmt.Println()
    }
}


////////////////////
$ go build gopl.io/ch2/echo4
$ ./echo4 a bc def
a bc def

$ ./echo4 -s / a bc def	//call var sep = flag.String
a/bc/def

$ ./echo4 -n a bc def	//call var n = flag.Bool 
a bc def$

$ ./echo4 help
Usage of ./echo4:
-n omit trailing newline
-s string
	separator (default " ")
////////////////////
```

### 2.3.3 The new Function

Another way to create a variable is to use the built-in function new. The expression new(T) creates an unnamed variable of type T, initializes it to the zero value of T, and returns its address, which is a value of type *T

```go
p := new(int)
fmt.Print(*p)
*p = 2
fmt.println(*p)
```

Implementation of new function

```go
func newInt() *int{
    return new(int)
}

func newInt() *int{
    var dummy int
    return *dummy
}
```

```go
p := new(int)
q := new(int)
fmt.Println(p == q) //false
```

There is one exception to this rule: two variables whose type carries no information and is therefore of size zero, such as struct{} or [0]int, may, depending on the implementation, have the same address.

The new function is relatively rarely used because the most common unnamed variables are of struct types, for which the ***struct literal syntax is more flexible.***

### 2.3.4 Lifetime of Variable

The lifetime of a variable is the interval of time during which it exists as the program executes.

> package-level variable is similar global variable in C/C++

How does the garbage collector know that a variable's storage can be reclaimed? The basic idea is that ***every package-level variable, and every local variable of each currently active function, can potentially be the start or root of a path to the variable in question, following pointers and other kinds of references that ultimately lead to the variable.***  If no such path exists, the variable has become unreachable, so it can no longer affect the rest of the computation.

***Example of deallocation***

```go
var global *int
func f() {
    var x int
    x = 1
    global = &x
}

//////////////////////

func g() {
    y := new(int)
    *y = 1
}
```

> In function f() variable x escapes function because of variable global so it cannot be recycled.
>
> But in function g() y can be recycled since y does not escape the function g()

In Go, programmers don't need to free memory but have to aware of the lifetime of variables.

## 2.4. Assignments

### 2.4.1 Tuple assignment

Another form of assignment, known as tuple assignment, allows several variables to be assigned at once. All of the right-hand side expression are evaluated before any of the variables are updated

```go
x, y = y, x
a[i], a[j] = a[j], a[i]

/////

func gcd(x, y int) int{
    for y != 0 {
        x, y, = y, x%y
    }
    return x
}
```

Tuple assignment can also make a sequence of trivial assignments more compact,

```go
i, j, k = 2, 3, 5

///
f, err = os.Open("foo.txt")	//Function call returns two values

v, ok = m[key]				//map lookup
v. ok = x.(T)				//type assertion
v, ok = <-ch				//channel receive

_, err = io.Copy9dst, src)	//discard byte count
_, ok = x.(T)				// check type but discard
```

### 2.4.2. Assignability

 Assignment statements are an explicit form of assignment, but there are many places in a program where an assignment occurs implicitly: a function call implicitly assigns the argument values to the corresponding parameter variables; a *return* statement implicitly assigns the return operands to the corresponding result variable; and a literal expression for a composite type.

```go
//slice
medals := []string{"gold", "silver", "bronze"}

medals[0] = "gold"
medals[1] = "silver"
medals[2] - "bronze"
```

An assignment, explicit or implicit, is always legal if the left-hand side (the variable) and the
right-hand side (the value) have the same type.

```go
variable = value		//legal!
```



## 2.5 Type Declarations

The type of a variable or expression defines the characteristics of the values it may take on

- such as their size
- how they are represented internally 
- the intrinsic operations that can be performed on them
- The methods associated with them.

A type declaration defines a new named type that has the same underlying type as an existing type. The names type provides a way to separate different and perhaps incompatible uses of the underlying type so that they cannot be mixed unintentionally.

```go
type name underlying-type
```

Type declarations most often appear at package level, where the named type is visible throughout the package. and if the name is exported, it is accessible from other package as well.

```go
package tempconv
import "fmt"

type Celsius float64
type Fahrenheit float64
//Celsius and Fahrenheit is float64 but their type is different because of Type Declarations
const (
    AbsoluteZeorC	Celsius = -273.15
    FreezingC		Celsius = 0
    BoilingC	 	Celsius = 100
)

func CToF(c Celsius) Fahreheit { return Fahrenheit(c*9/5 + 32)}
func FToC(f Fahrenheit) Celsius{ return Celsius((f - 32) * 5/9)}
```

> ***Similar with typedef in C/C++***

***For every type T, there is a corresponding conversion operation T(x) that converts the value x to type T.*** A conversion from one type to another is ***allowed if both have the same underlying type, or if both are unnamed pointer types that point to variables of the same underlying type***

> underlying type of Celsius and Fahrenheit is float64 -> same!

-> These conversions change the type but not the representation of the value. If x is assignable to T, a ***conversion is permitted but is usually redundant***

- Conversions are also allowed between numeric types and between string and some slice types.

- The underlying type of a named type determines its structures and representation

  -> The set of intrinsic operations it supports

> Celsius + Celsius -> result is Celsius too
>
> Celsius + Fahrenheit -> compile error:Type mismatching

```go
var c Celcius
var f Fahrenheit
fmt.Println(c == Celcius(f)) //true!
```

```go
func (c Celsius) String() string { return fmt.Sprintf("%g C", c)}                       // String() : function name
// There is (c Celsius) in front of function name
//-> return string value and in string there is Celsius c

c := FToC(212.0)
fmt.Println(c.String()) // "100°C"
fmt.Printf("%v\n", c) // "100°C"; no need to call String explicitly
fmt.Printf("%s\n", c) // "100°C"
fmt.Println(c) // "100°C"
fmt.Printf("%g\n", c) // "100"; does not call String
fmt.Println(float64(c)) // "100"; does not call String

//for more information about print verb
//https://golang.org/pkg/fmt/
```

