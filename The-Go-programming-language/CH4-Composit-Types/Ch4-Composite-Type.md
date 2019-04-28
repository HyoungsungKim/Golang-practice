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

//For example, make([]int, 0, 10) allocates an underlying array
//of size 10 and returns a slice of length 0 and capacity 10 that is
//backed by this underlying array.
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

***In Go, a map is a reference to a hash table***, and a map type is written map[K]V,

- K keys
- V Values

- All of the *keys in a given map are of the same type,*
- All of the *values are of the same type*
- But the *keys need not be of the same type as the values.*

The key type K must be comparable using ==, so that the map can test whether a given key is equal to one already within it. ***Though floating-point numbers are comparable, it’s a bad idea to compare floats for equality*** and, as we mentioned in Chapter 3, especially bad if NaN is a possible value.

```go
ages := make(map[string]int)
//Other way
ages := mpa[strng]int{
    "alice":	31,
    "charlie":	34,
}
//It is equivalent to
ages := make(map[string]int)
ages["alice"] = 31
ages["charlie"] = 34
fmt.Println(ages["alice"])	// 32
delete(ages, "alice")	//remove element ages["alice"]
ages["Bob"]	// key is 0
//It is possible
ages["Bob"] += 1
ages["Bob"]++
```

A map element is not a variable, and we cannot take its address:

```go
_= &ages["Bob"]	//Compile error: cannot take address of map element
```

To enumerate all the key/value pairs in the map, we use a range-based for loop similar to those we saw for slices.

```go
for name, age := range ages{
    fmt.Printf(%s\t%d\n", name, age)
}
```

The order of map iteration is unspecified, and different implementations might use a different hash function, leading to a different ordering. In practice, the order is random, varying from one execution to the next.

To enumerate the key/value pairs in order, we must sort the keys explicitly, for instance, using the Strings function from the sort package if the keys are strings. This is a common pattern:

```go
import "sort"

var name []string
for name := range ages{
    names = append(names, name)
}
sotr.String(names)
for _, name := range names {
    fmt.Printf("%s\t%d\n", name, ages[name])
}
```

Since we know the final size of names from the outset, it is more efficient to allocate an array of the required size up front. ***The statement below creates a slice that is initially empty but has sufficient capacity to hold all the keys of the ages map:*** 

```go
names := make([]string, 0, len(ages))
//make(type, size, capacity)
```

***In the first range loop above, we require only the keys of the ages map, so we omit the second loop variable.*** In the second loop, we require only the elements of the names slice, so we use the blank identifier _ to ignore the first variable, the index.

The zero value for a map type is nil, that is, a reference to no hash table at all.

```go
var ages map[string]int
fmt.Println(ages == nil) // "true"
fmt.Println(len(ages) == 0) // "true"
```

Most operations on maps, including lookup, *delete*, *len*, and *range* loops, are safe to perform on a nil map reference, since it behaves like an empty map. ***But storing to a nil map causes a panic:***

```go
ages["carol"] = 21	//panic: assignment to entry in nil map
```

***You must allocate the map before you can store into it.***

***Accessing a map element by subscripting always yields a value.*** If the key is present in the map, you get the corresponding value; ***if not, you get the zero value for the element type, as we saw with ages["Bob"].***

> Can call unpresented value,
>
> but cannot store value 

For many purposes that’s fine, ***but sometimes you need to know whether the element was really there or not.*** 

```go
age, ok := ages["Bob"]
if !ok {/*"Bob" is not a key in this map; age == 0.*/}
// You will often see these two statements combined, like this:
if age, ok := ages["Bob"]; !ok{/*...*/}
```

Subscripting a map in this context yields two values; the second is a boolean that reports whether the element was present. ***The boolean variable is often called ok, especially if it is immediately used in an if condition.***

***As with slices, maps cannot be compared to each other; the only legal comparison is with nil.*** 

```go
func equal(x, y map[string]int) bool{
    if len(x) != len(y) {
        return false
    }
    for k, xv := range x{
        //How to use ok and operator ||
        if yv, ok := y[k]; !ok || yv != xv {
            return false
        }
    }
    return true
}

// True if equal is written incorrectly.
equal(map[string]int{"A": 0}, map[string]int{"B": 42})
```

***Go does not provide a set type, but since the keys of a map are distinct, a map can serve this purpose.***

To illustrate, the program *dedup* reads a sequence of lines and ***prints only the first occurrence of each distinct line.*** The *dedup* program uses a map whose keys represent the set of lines that have already appeared to ensure that subsequent occurrences are not printed.

```go
func main() {
    seen := make(map[string]bool)	// a set of string
    input := bufio.NewScanner(os.Stdin)
    for input.Scan() {
        line := input.Text()
        if !seen[line] {
            seen[line] = true
            fmt.Println(line)
        }
    }
    if err := input.Err(); err != nil {
        fmt.Fpeintf(os.Stderr, "dedup: %v\n", err)
        os.Exit(1)
    }
}
```

Sometimes we need a map or set ***whose keys are slices***, but because a map’s keys must be comparable,
this cannot be expressed directly. However it can be done in two steps.

- First, we define a helper function k that maps each keys to a string, with the property that k(x) == k(y) if and only if we consider x and y equivalent.
- Then we create a map whose keys are strings applying the helper function to each key before we access the map.

The example below uses a map to record the number of times *Add* has been called with a given list of strings. ***It uses fmt.Sprintf to convert a slice of string into a single string*** that is a suitable map key, quoting each slice element with %q to record string boundaries faithfully:

```go
var m = make(map[string]int)
func k(list []string) string {return fmt.Sprintf("%q", list)}
func Add(list []string) {m[k(list)]++}
func Count(list []string) int {return m[k(list)]}
```

The same approach can be used for any non-comparable key type, not just slices.



97page 하는 중