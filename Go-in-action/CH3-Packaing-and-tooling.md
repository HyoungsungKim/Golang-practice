# CH3 Packaging and tooling

Packages are a critical concept in Go. The idea is to separate semantic units of functionality into different packages. When you do this, you enable code reuse and control the use of the data inside each package.

## 3.1 Package

All .go files must declare the package that they belong to as the first line of the file excluding whitespace and comments. Packages are contained in a single directory. You may not have multiple packages in the same directory, nor may you split a package across multiple directories. ***This means that all .go files in a single directory must declare the same package name.***

### 3.1.1 Package-naming conventions

The convention for naming your package is to use the name of the directory containing it. When naming
your packages and their directories, you should use short, concise, lowercase names. Your package name is used as the default name when your package is imported, but it can be overridden. 

### 3.1.2 Package main

All of the executable programs you build in Go must have a package called main.

## 3.2 Imports

### 3.2.2 Named imports

What happens when you need to import multiple packages with the same name? When this is the case, both of these packages can be imported by using named imports.

> C++의 namespace 같은?

```go
package main
import (
	"fmt"
    myfmt "mylib/fmt"
)

func main(){
    fmt.Println("Hello world")
    myfmt.Println("mylib/fmt")
}
```

## 3.3 init

Each package has the ability to provide as many init functions as necessary to be invoked at the beginning of execution time. All the `init` functions that are discovered by the compiler are scheduled to be executed prior to the main function being executed.

## 3.5 Going farther with Go developer tools

### 3.5.1 go vet

It won’t write code for you, but once you've written some code, ***the vet command will check your code for common errors.*** Let’s look at the types of errors vet can catch:

- Bad parameters in Printf -style function calls
- Method signature errors for common method definitions
- Bad struct tags
- Unkeyed composite literals

### 3.5.3 go documentation

If you’re working at a command prompt, you can use the `go doc`

```go
go doc tar
go doc fmt
```

The Go documentation is also available in a browsable format.

```go
godoc -http=:6060
```

