# Effective Go

## Commentary

The program—and web server—`godoc` processes Go source files to extract documentation about the contents of the package. Comments that appear before top-level declarations, with no intervening newlines, are extracted along with the declaration to serve as explanatory text for the item. The nature and style of these comments determines the quality of the documentation `godoc` produces.

Every package should have a *package comment*, a block comment preceding the package clause. For multi-file packages, the package comment only needs to be present in one file, and any one will do. The package comment should introduce the package and provide information relevant to the package as a whole. It will appear first on the `godoc` page and should set up the detailed documentation that follows.

## Names

Names are as important in Go as in any other language. ***They even have semantic effect: the visibility of a name outside a package is determined by whether its first character is upper case.*** It's therefore worth spending a little time talking about naming conventions in Go programs.

### Package names

By convention, packages are given lower case, single-word names; there should be no need for underscores or mixedCaps. Err on the side of brevity, since everyone using your package will be typing that name. And don't worry about collisions *a priori*. ***The package name is only the default name for imports; it need not be unique across all source code, and in the rare case of a collision the importing package can choose a different name to use locally.*** In any case, confusion is rare because the file name in the import determines just which package is being used.

The importer of a package will use the name to refer to its contents, so exported names in the package can use that fact
to avoid stutter

Similarly, the function to make new instances of `ring.Ring`—which is the definition of a *constructor* in Go—would
normally be called `NewRing`, but since `Ring` is the only type exported by the package, and ***since the package is called `ring`, it's called just `New`, which clients of the package see as `ring.New`. Use the package structure to help you choose good names.***

### Interface names

By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification
to construct an agent noun: `Reader`, `Writer`, `Formatter`, `CloseNotifier` etc.

## Semicolons

Like C, Go's formal grammar uses semicolons to terminate statements, but unlike in C, those semicolons do not appear in the source. Instead the lexer uses a simple rule to insert semicolons automatically as it scans, so the input text is mostly free of them.

```go
if i < f() {
    g()
}

if i < f()  // wrong!
{           // wrong!
    g()
}
```

## Functions

### Defer

Go's `defer` statement schedules a function call (the *deferred* function) to be run immediately before the function executing the `defer` returns.  It's an unusual but effective way to deal with situations such as resources that must be released regardless of which path a function takes to return. The canonical examples are unlocking a mutex or closing a file.

## Data

### Allocation with `new`

Go has two allocation primitives, the built-in functions `new` and `make`. They do different things and apply to different types, which can be confusing, but the rules are simple.

Let's talk about `new` first. It's a built-in function that allocates memory, but unlike its namesakes in some other languages it does not *initialize* the memory, it only *zeros* it. ***That is, `new(T)` allocates zeroed storage*** for a new item of type `T` and returns its address, a value of type `*T`. In Go terminology, it returns a pointer to a newly allocated zero value of type `T`.

### Constructors and composite literals

 Sometimes the zero value isn't good enough and an initializing constructor is necessary, as in this example derived from package `os`. 

Note that, unlike in C, it's perfectly OK to return the address of a local variable;

```go
return &File{fd, name, nil, 0}
```

### Allocation with `make`

Back to allocation. The built-in function `make(T, *args*)` serves a purpose different from `new(T)`. ***It creates slices, maps, and channels only,*** and it returns an ***initialized (not zeroed) value of type*** `T` (not `*T`). The reason for the distinction is that these three types represent, under the covers, references to data structures that must be initialized before use. A slice, for example, is a three-item descriptor containing a pointer to the data (inside an array), the length, and the capacity, and until those items are initialized, the slice is `nil`. For slices, maps, and channels, `make` initializes the internal data structure and prepares the value for use.