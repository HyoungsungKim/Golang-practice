# CH 4 Composite Types

In this chapter, we’ll take a look at composite types, the molecules created by combining the basic types in various ways. We’ll talk about four such types

- arrays
- slices
- maps
- structs

## 4.1 Arrays

An array is a ***fixed-length*** sequence of zero or more elements of a particular type.

- Because of their fixed length, arrays are rarely used directly in Go.

Slices, ***which can grow and shrink***, are much more versatile, but to understand slices we must understand arrays first.

```go
var a [3]int
fmt.Println(a[0])
//In string, s[0] return ASCII
fmt.Println(a[len(a) - 1])

for i, v := range a {
    fmt.Printf("%d %d\n", i, v)
}

for _, v := range a {
    fmt.Printf(%d\n", v)
}
```

By default, the elements of a new array variable are initially set to the zero value for the element type, which is 0 for numbers. We can use an array literal to initialize an array with a list of values:

```go
var q [3]int = [3]int{1, 2, 3}
var r [3]int = [3]int{1, 2}
fmt.Println(r[2]) // "0"
```

In an array literal, if an ellipsis ‘‘...’’ appears in place of the length, the array length is determined
by the number of initializers.

```go
q := [...]int{1, 2, 3}
fmt.Printf("%T\n", q) // "[3]int"
```

In this form, indices can appear in any order and some may be omitted; as before, unspecified values take on the zero value for the element type. For instance,

```go
r := [...]int{99: -1}
//{0, 0, 0, 0,....,-1}
```

defines an array r with 100 elements, all zero except for the last, which has value −1.

If an array’s element type is comparable then the array type is comparable too

```go
a := [2]int{1, 2}
b := [...]int{1, 2}
c := [2]int{1, 3}
fmt.Println(a == b, a == c, b == c) // "true false false"
d := [3]int{1, 2}
fmt.Println(a == d) // compile error: cannot compare [2]int == [3]int
```

When a function is called, ***a copy of each argument value is assigned*** to the corresponding parameter variable, ***so the function receives a copy, not the original.*** 

***Go treats arrays like any other type, but this behavior is different from languages that implicitly pass arrays by reference.*** 

```go
func zero(ptr *[32]byte) {
    for i := range ptr {
    	ptr[i] = 0
    }
}
```

> Similar with C/C++

## 4.2 Slices

***Slices represent variable-length sequences whose elements all have the same type.*** A slice type
is written []T, where the elements have type T; it looks like an array type without a size.

A slice has three components: a pointer, a length, and a capacity. 

- Pointer : ***points to the first element of the array that is reachable through the slice,***  which is not
  necessarily the array’s first element.

- Length :  the number of slice elements
- Capacity : Capacity is usually the number of elements between the start of the slice and the end
  of the underlying array.

Multiple slices can share the same underlying array and may refer to overlapping parts of that
array. 

```go
months := [...]string{1: "January", /* ... */, 12: "December"}
```

The slice operator s[i:j], where $0 \leq i \leq j \leq cap(s)$, creates a new slice that refers to elements i through j - 1
of the sequence s, which maybe an array variable, a pointer to an array, or another slice. The resulting slice has j- i elements. If i is omitted, it’s 0, and if j is omitted, it’s len(s).

Slicing ***beyond cap(s) causes a panic***, but ***slicing beyond len(s) extends the slice***, so the result may be longer than the original:

Since a slice contains a pointer to an element of an array, passing a slice to a function per mits
the function to modify the underlying array elements.

There are two reasons why deep equivalence is problematic

- Unlike array elements, the elements of a slice are indirect, making it possible for a slice to contain itself.
- Because slice elements are indirect, a fixed slice value may contain different elements at different times as the contents of the underlying array are modified.

>var a [3]int // array
>
>var b[...]int //array
>
>var c[]int // slice

Because a hash table such as Go’s map type makes only shallow copies of its keys, it requires that equality for each key remain the same throughout the lifetime of the hash table. Deep equivalence would thus make slices unsuitable for use as map keys.

***For reference types like pointers and channels, the == operator tests reference identity, that is, whether the two entities refer to the same thing.*** The safest choice is to disallow slice comparisons altogether.

***The only legal slice comparison*** is against nil, as in

```go
if summer == nil { /* ... */ }

///////////////////////////

var s []int // len(s) == 0, s == nil
s = nil // len(s) == 0, s == nil
s = []int(nil) // len(s) == 0, s == nil
s = []int{} // len(s) == 0, s != nil
```

***So, if you need to test whether a slice is empty, use len(s) == 0, not s == nil***. Other than comparing equal to nil, a nil slice behaves like any other zero-length slice; reverse(nil) is perfectly safe.

The built-in function *make* creates a slice of a specified element type, length, and capacity.

```go
make([]T, len)
make([]T, len, cap) // same as make([]T, cap)[:len]
```

***It creates an unnamed array variable*** and returns a slice of it; the array is accessible only through the returned slice. 

### 4.2.1 The append Function

The built-in append function appends items to slices:

```go
var runes []rune
for _, r := range "Hello, 世界" {
	runes = append(runes, r)
}
fmt.Printf("%q\n", runes) // "['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']"
```

The append function is crucial to understanding how slices work.

> crucial : 중대한

```go
func appendInt(x []int, y int) []int {
    var z []int
    zlen := len(x) + 1
    if zlen <= cap(x) {
        z = x[:zlen]
    } else {
        zcap := zlen
        if zcap < 2 * len(x) {
            zcap = 2 * len(x)
        }
        z = make([]int, zlen, zcap)
        copy(z, x)
    }
    z[len(x)] = y
    return z
}
//Expand slice 2 times
```

> It works like the *vector* of C++
>
> But Go pads zeoros

```go
var x []int
x = append(x, 1)
x = append(x, 2, 3)
x = append(x, 4, 5, 6)
x = append(x, x...) // append the slice x
fmt.Println(x) // "[1 2 3 4 5 6 1 2 3 4 5 6]"
```



### 4.2.2 In-place Slice Techniques

Let’s see more examples of functions that, like *rotate* and *reverse*, modify the elements of a slice in place.

```go
// Nonempty is an example of an inplace
slice algorithm.
package main
import "fmt"
// nonempty returns a slice holding only the nonempty strings.
// The underlying array is modified during the call.
func nonempty(strings []string) []string {
    i := 0
    for _, s := range strings {
        if s != "" {
            strings[i] = s
            i++
        }
    }
	return strings[:i]
}
data := []string{"one", "", "three"}
fmt.Printf("%q\n", nonempty(data)) // `["one" "three"]`
fmt.Printf("%q\n", data) // `["one" "three" "three"]`
```

## 4.3 Maps

The hash table is one of the most ingenious and versatile of all data structures.

> ingenious : 기발한
>
> versatile : 다재다능한

In Go, a map is a reference to a hash table, and a map type is written map[K]V, where K and V
are the types of its keys and values.

- All of the keys in a given map are of the same type,
- All of the values are of the same type
- But the keys need not be of the same type as the values.

