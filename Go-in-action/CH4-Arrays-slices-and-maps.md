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

