# CH 3. Basic Data Types

Go’s types fall into four categories: basic types, aggregate types, reference types, and interface types.

Basic types, the topic of this chapter, include numbers, strings, and booleans.

- Aggregate types—arrays  and structs—form more complicated data types by combining values of several simpler ones.
- Reference types are a diverse group that includes pointers, slices, maps, functions, and channels, ***but what they have in common is that they all refer to program variables or state indirectly***, so that the effect of an operation applied to one reference is observed by all copies of that reference
- Finally interface

## 3.1 Integers

The type *rune* is an synonym for *int32* and conventionally indicates that a value is a Unicode code point. The two names maybe used interchangeably. Similarly, the type *byte* is an synonym for *uint8*, and emphasizes that the value is a piece of raw data rather than a small numeric quantity.

There is an unsigned integer type *uintptr*, whose width is not specified but is sufficient to ***hold all the bits of a pointer value.*** The *uintptr* type is used only for low-level programming, such as at the boundary of a Go program with a C library or an operating system.

For integers, *+x* is a *shorthand for 0+x* and *-x* is a *shorthand for 0-x*; for *floating-point and complex numbers, +x is just x and x is the negation of x.*

```go
& bitwise AND
| bitwise OR
^ bitwise XOR
&^ bit cle ar (AND NOT)
<< lef t shif t
>> right shif t
```

***

***&^ bit clear (AND NOT)*** in the expression z = x &^ y, each bit of z is 0 if the corresponding bit of y is 1; other wise it equals the corresponding bit of x. 

-> Not of and's result

Arithmetically, ***a left shift x<<n is equivalent to multiplication by 2^n and a right shift x>>n is equivalent to the floor of division by 2^n .***

When printing numbers using the fmt package, we can control the radix and format with the ***%d, %o, and %x verbs,*** as shown in this example.

```go
o := 0666
fmt.Printf("%d %[1]o %#[1]o\n", o) // "438 666 0666"
x := int64(0xdeadbeef)
fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", x)
// Output:
// 3735928559 deadbeef 0xdeadbeef 0XDEADBEEF
```

Usually a Printf format string containing ***multiple % verbs would require the same number of extra operands,*** but ***the [1] ‘‘adverbs’’*** after % tell Printf to ***use the first operand over and over again.*** 

> In this code first operand is "d" of "%d"

Second, ***the # adverb for %o or %x or %X tells Printf to emit a 0 or 0x or 0X prefix respectively.***



## 3.2 Floating-Point Numbers

Go provides two sizes of floating-point numbers, float32 and float64. Their arithmetic properties are governed by the IEEE 754 standard implemented by all moder n CPUs.

***A float32 provides approximately six decimal digits of precision, whereas a float64 provides about 15 dig its*** float64 should be preferred for most purpos es because float32 computations accumulate error rapidly unless one is quite careful, and the smallest positive integer that cannot be exactly represented as a float32 is not large.

The code above prints ***the powers of e*** with three decimal digits of precision, aligned in an eight-character field:

the positive and negative infinities, which represent numbers of excessive magnitude and the result of  division by zero; and NaN (‘‘not a number’’), the result of such mathematically dubious operations as 0/0 or Sqrt(1).

```go
var z float64
fmt.Println(z, z,1/z, 1/z, z/z) // "0 0 +Inf Inf NaN" z/z is 0/0
```

The function math.IsNaN tests whether its argument is a not-a-number value, and math.NaN returns such a value.

It’s tempting to use NaN as a sentinel value in a numeric computation, ***but testing whether a specific computational result is equal to NaN is fraught with peril because any comparison with NaN always yields false***

```go
nan := math.NaN()
fmt.Println(nan == nan, nan < nan, nan > nan) // "false false false"
```



## 3.3 Complex Numbers

complex64 and complex128, whose components are float32 and float64 respectively. The built-in function complex creates a complex number from its real and imaginary components, and the built-in real and imag functions extract those components:

```go
var x complex128 = complex(1, 2) // 1+2i
var y complex128 = complex(3, 4) // 3+4i
fmt.Println(x*y) // "(5+10i)"
fmt.Println(real(x*y)) // "5"
fmt.Println(imag(x*y)) // "10"
```

Complex numbers may be compared for equality *with == and !=*. ***Two complex numbers are equal if their real parts are equal and their imaginary parts are equal.***

## 3.4 Booleans

A value of type bool, or boolean, has only two possible values, true and false.

## 3.5 Strings

A string is an immutable sequence of bytes. Strings may contain arbitrary data, including bytes with value 0, but usually they contain human-readable text.

***the index operation s[i] retrieves the i-t h byte of string s***

```go
s := "hello, world"
fmt.Println(len(s)) 	// "12"
fmt.Println(s[0], s[7]) // "104 119" ('h' and 'w')
//(s[0] == 'h') -> possible!!! 'h' is changed tp ASCII
c := s[len(s)] 			// panic: index out of range
fmt.Println(s[0:5]) 	// "hello"

fmt.Println(s[:5]) 		// "hello"
fmt.Println(s[7:])		// "world"
fmt.Println(s[:]) 		// "hello, world"
fmt.Println("goodbye" + s[5:]) // "goodbye, world"
//[:num] from begin to num
//[num:] from num to end
//[:1] first character
s[0] = 'L' // compile error: cannot assign to s[0]
```

### 3.5.4 Strings and Bytes Slices

Four standard packages are particularly important for manipulating strings: bytes, strings, strconv, and unicode. 

- The ***bytes*** package has similar functions for manipulating slices of bytes, of type []byte, which share some properties with strings. Because strings are immutable, building up strings incrementally can involve a lot of allocation and copying. In such cases, it’s more efficient to use the bytes:***Buffer type***

- The ***strconv*** package provides ***functions for converting boolean, integer, and floating-point
  values to and from their string representations,*** and functions for quoting and unquoting
  strings.
- The ***unicode*** package provides functions like *IsDigit*, *IsLetter*, *IsUpper*, and *IsLower* for classifying runes. Each function takes a single rune argument and returns a boolean. Conversion functions like *ToUpper* and *ToLower* convert a rune into the given case if it is a letter.

The *basename* function below was inspired by the Unix shell utility of the same name. In

```go
fmt.Println(basename("a/b/c.go")) // "c"
fmt.Println(basename("c.d.go")) // "c.d"
fmt.Println(basename("abc")) // "abc"
```

Strings can be converted to byte slices and back again:

```go
s := "abc"
b := []byte(s)
s2 := string(b)
```

To avoid conversions and unnecessary memory allocation, many of the utility functions in the bytes package directly parallel their counterparts in the strings package.

The *bytes* package provides the *Buffer* type for efficient manipulation of byte slices. A *Buffer* starts out empty but grows as data of types like *string*, *byte*, and *[]byte* are written to it. As the example below shows, a *bytes.Buffer* variable requires no initialization because its zero value is usable:

```go
// intsToString is like fmt.Sprintf(values) but adds commas.
func intsToString(values []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
        //insert 1 when i == 0
		if i > 0 {
			buf.WriteString(", ")
		}
	fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}
func main() {
    fmt.Println(intsToString([]int{1, 2, 3})) // "[1, 2, 3]"
}
```

### 3.5.5 Conversion between String and Numbers

```go
x := 123
y := fmt.Sprintf("%d", x)
fmt.Println(y, strconv.Itoa(x)) // "123 123"

fmt.Println(strconv.FormatInt(int64(x), 2)) // "1111011"
//The fmt.Printf verbs %b, %d, %u, and %x are often more convenient than Format functions

s := fmt.Sprintf("x=%b", x) // "x=1111011"
//To parse a string representing an integer, use the strconv functions Atoi or ParseInt, or ParseUint for unsig ned integers:

x, err := strconv.Atoi("123") // x is an int
//Atoi : ASCII to integer
y, err := strconv.ParseInt("123", 10, 64) // base 10, up to 64 bits
```



## 3.6 Constant

***Constants are expressions whose value is known to the compiler and whose evaluation is guaranteed
to occur at compile time, not at run time.*** The underlying type of every constant is a basic type: boolean, string, or number.

```go
const pi = 3.14159 // approximately; math.Pi is a better approximation

const(
    e = 2.71828182845904523536028747135266249775724709369995957496696763
    pi = 3.14159265358979323846264338327950288419716939937510582097494459
)
```

A constant declaration may specify a type as well as a value, but in the absence of an explicit type, the type is inferred from the expression on the right-hand side.

```go 
const pi = 3.14 //pi is float32
```



When a sequence of constants is declared as a group, the right-hand side expression may be omitted for all but the first of the group, ***implying that the previous expression and its type should be used again.*** For example:

```go
const (
    a = 1
    b
    c = 2
    d
)
fmt.Println(a, b, c, d) // "1 1 2 2"
```

### 3.6.1 The Constant Generator iota

A const declaration may use the constant generator iota, ***which is used to create a sequence of related values without spelling out each one explicitly.***

***In a const declaration, the value of iota begins at zero and increments by one for each item in the sequence.***

```go
type Weekday int
const (
    Sunday Weekday = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
)
//Weekday is 0 so sunday is 0 too
//Monday will be 1
```

```go
type Flags uint
const (
    FlagUp Flags = 1 << iota // is up
    FlagBroadcast // supports broadcast access capability
    FlagLoopback // is a loopback interface
    FlagPointToPoint // belongs to a pointtopoint
    link
    FlagMulticast // supports multicast access capability
)
//increment as power of 2
```

```go
//As a more complex example of iota, this declarat ion names the powers of 1024:
const (
_ = 1 << (10 * iota)
    KiB // 1024
    MiB // 1048576
    GiB // 1073741824
    TiB // 1099511627776 (exceeds 1 << 32)
    PiB // 1125899906842624
    EiB // 1152921504606846976
    ZiB // 1180591620717411303424 (exceeds 1 << 64)
    YiB // 1208925819614629174706176
)
```

### 3.6.2 Untyped Constant

By deferring this commitment, untyped constants

- Retain their higher precision until later
- They can participate in many more expressions than committed constants without requiring conversions.

> Deffer : 미루다

```go
var x float32 = math.Pi
var y float64 = math.Pi
var z complex128 = math.Pi
```

If math.Pi had been committed to a ***specific type  such as float64, the result would not be as precise,*** and ***type conversions would be required to use it when a float32 or complex128 value is wanted***

***Only constants can be untyped.*** When an untyped constant is assigned to a variable, as in the first statement below, or appears on the right-hand side of a variable declaration with an explicit type, as in the other three statements, ***the constant is implicitly converted to the type of that variable if possible.***

```go
var f float64 = 3 + 0i // untyped complex > float64
f = 2 // untyped integer > float64
f = 1e123 // untyped floatingpoint > float64
f = 'a' // untyped rune > float64
```

Whether implicit or explicit, converting a constant from one type to another requires that the target type can represent the original value. ***Rounding is allowed for real and complex floating-point numbers***

```go
const (
    deadbeef = 0xdeadbeef // untyped int with value 3735928559
    a = uint32(deadbeef) // uint32 with value 3735928559
    b = float32(deadbeef) // float32 with value 3735928576 (rounded up)
    c = float64(deadbeef) // float64 with value 3735928559 (exact)
    d = int32(deadbeef) // compile error: constant overflows int32
    e = float64(1e309) // compile error: constant overflows float64
    f = uint(1)
    // compile error: constant underflows uint
)
```

```go
//Convert type explictly
var i = int8(0)
var i int8 = 0
```

These defaults are particularly important when converting an untyped constant to an interface value  since they determine its dynamic type.

```go
fmt.Printf("%T\n", 0) // "int"
fmt.Printf("%T\n", 0.0) // "float64"
fmt.Printf("%T\n", 0i) // "complex128"
fmt.Printf("%T\n", '\000') // "int32" (rune)
```

