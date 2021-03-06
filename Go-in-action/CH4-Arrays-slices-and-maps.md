# CH4 Arrays, slices, and maps

## 4.1 Arrays internals and fundamentals

### 4.1.1 Internals

Arrays are valuable data structures because the memory is allocated sequentially.

### 4.1.2 Declaring and initializing

```go
var array [5]int
array := [5]int{10, 20, 30, 40, 50}
array := [...]int{10, 20, 30, 40, 50}
array := [...]{1:10, 2:20}	//0:0, 1:10, 2:20, 3:0, 4:0
array [2] = 35	//0:0, 1:10, 2:35, 3:0, 4:0
```

If you need more elements, you need to create a new array with the length needed and then copy the values from one array to the other.

### 4.1.3 Working with arrays

```go
array := [5]*int{0: new(int), 1:new(int)}
*array[0] = 10
*array[1] = 20
//array : 0:addr -> 10, 1:addr -> 20, 2:nil, 3:nil, 4:nil

var array1 [3]*string
var array2 := [3]*string{new(string, new(string), new(string))}
*array2[0] = "Red"
*array2[1] = "Blue"
*array2[2] = "Green"
array1 = array2
//After the copy, you have two arrays pointing to the same strings
//주소는 다름, 가르키는 값은 같음
```

### 4.1.4 Multidimensional arrays

```go
var array [4][2]int
array := [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}
array := [4][2]int{1: {20, 21}, 3: {40,41}}
array[0][0] = 10
```

### 4.1.5 Passing arrays between functions

Passing an array between functions can be an expensive operation in terms of memory and performance.

```go
var array [1e6]int

foo(&array)
func foo(array *[1e6]array) {
    ...
}
```

## 4.2 Slice internals and fundamentals

A slice is a data structure that provides a way for you to work with and manage collections of data. Slices are built around the concept of dynamic arrays that can grow and shrink as you see fit. They’re flexible in terms of growth because they have their own built-in function called `append`

### 4.2.2 Creating and initilazing

```go
slice := make([]string, 5)
slice := make([]string, 3, 5)	//(type, size, capacity)
slice := []string{"Red", "Blue", "Green", "Yellow", "Pink"}
slice := []int{10, 20, 30}
slice := []string{99:""}
```

Remember, if you specify a value inside the [ ] operator, you’re creating an array. If you don’t specify a value, you’re creating a slice.

#### Nil and empty slices

Sometimes in your programs you may need to declare a `nil` slice. A `nil` slice is created by declaring a slice without any initialization.

```go
var slice []int
slice := make([]int,0)
slice := []int{}
```

### 4.2.3 Working with slices

```go
slice := []int{10, 20, 30, 40, 50}
slice[1] = 25
newSlice := slice[1:3]	//newSlice = {20, 30, 40}
```

***Not copy but allocate same pointer.*** therefore if elements of new slices are changed, original elements also changed.

#### Growing slices

Go takes care of all the operational details when you use the built-in function append. 

```go
slice := []int{10, 20, 30, 40, 50}
newSlice := slice[1:3]
newSlice := append(newSlice, 60)
```

```go
package main

import (
	"fmt"
)

func main() {
	slice :=[]int{10, 20, 30, 40, 50}
	newSlice :=slice[1:3]
	newSlice = append(newSlice, 100)
	fmt.Print(slice)
}
//output : 10, 20, 30, 100, 50
```

#### Three index slices

The purpose is not to increase capacity, but to restrict the capacity. 

```go
source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}
slice := source[2:3:4]
```

```go
source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}
slice := source[2:3:3]
slice = append(slice, "Kiwi")
```

```go
s1 := []int{1, 2}
s2 := []int{3, 4}
s4 := append(s1, s2)
fmt.Printf("%v\n", append(s1, s2...))
//슬라이스에 슬라이스 붙일때 ... 씀
//[1 2 3 4]
```

#### Iterating over slices

Since a slice is a collection, you can iterate over the elements.

```go
slice := []int{10, 20, 30, 40}
for index, value := range slice {
    fmt.Printf ("Index: %d Value: %d\n", index, value)
}
```

***It’s important to know that range is making a copy of the value, not returning a reference.***

```go
slice := []int{10, 20, 30, 40}

for index, value := range slice {
    fmt.Printf("Value: %d Value-Addr: %X ElemAddr: %X\n", value, &value, &slice[index])
}

/*
output: address is different
Value: 10 Value-Addr: 10500168 ElemAddr: 1052E100
.
.
.
*/
```

If you don’t need the index value, you can use the underscore character to discard the value.

```go
slice := []int{10, 20, 30, 40}

for _, value := range slice {
    fmt.Println("Value : %d", value)
}
```

***If you need more control iterating over a slice, you can always use a traditional for loop.***

```go
slice := []int{10, 20, 30, 40}

for index := 2; index < len(slice); index++ {
    fmt.Printf("Index %d Value: %d\n", index, slice[index])
}
```

### 4.2.4 Multidimensional slices

```go
slice := [][]int{{10}, {100, 200}}
slice[0] = append(slice[0], 20)
```

### 4.2.5 Passing slices between functions

```go
slice := make([]int, 1e6)
slice = foo(slice)

func foo(slice []int) []int {
    ...
    return slice
}
```

***On a 64-bit architecture, a slice requires 24 bytes of memory.*** The pointer field requires 8 bytes, and the length and capacity fields require 8 bytes respectively.

You don’t need to pass pointers around and deal with complicated syntax. You just create copies of your slices, make the changes you need, and then pass a new copy back.

> slice 복사해서 함수 매개변수로 넣으면 알아서 주소로 전달 됨(와... 좋네)

## 4.3 Map internals and fundamentals

A map is a data structure that provides you with an unordered collection of key/value pairs.

### 4.3.1 Internals

Maps are unordered collections, and ***there’s no way to predict the order in which the key/value pairs will be returned.*** Even if you store your key/value pairs in the same order, every iteration over a map could return a different order. This is because a map is implemented using a hash table.

***Just remember one thing: a map is an unordered collection of key/value pairs.***

### 4.3.2 Creating and Initializing

```go
dict := make(map[string]int)
dict := map[string]string("Red": "#da1337", "Orange": "#e95a22")
```

map[a]b -> a is key, b is value

### 4.3.3 Working with maps

```go
colors := map[string]string()
colors["Red"] = "#da1337"
```

***A `nil` map can't be used to store key/value pairs.***

```go
var colors map[string]string
colors["Red"] = "#da1337"

//Runtime error
```

```go
value, exists := colors["Blue"]

if exists {
    fmt.Println(value)
}
```

hen you index a map in Go, it will always return a value, even when the key doesn't exist. In this case, the zero value for the value's type is returned.

```go
value, exists := colors["Blue"]

if value != "" {
    fmt.Println(value)
}
// value == "" : true
// value != "" : false
```

iterating maps

```go
colors := map[string]string {
    "Red":"red",
    "Blue":"blue",
    "Green":"green",
}

for key, value := range colors {
    fmt.Println("key :%s value:%s",key, value)
}
```

Removing an item from a map

```go
delete(colors, "Red")
```

### 4.3.4 passing maps between functions

***Passing a map between two functions doesn’t make a copy of the map.***

```go
func map() {
    colors := map[string]string {
        "Red":"red"
        "Blue":"blue"
        "Green":"green"
    }
    
    for key, value := range colors {
        fmt.Printlf("Key: %s Value: %s ", key, value)
    }
    removeColor(colors, "Blue")
    
    for key, value := range colors {
        fmt.Printf("key: %s Value: %s", key, value)
    }
}

func removeColor(colors map[string]string, key string ) {
    delete(colors, key)
}
```

