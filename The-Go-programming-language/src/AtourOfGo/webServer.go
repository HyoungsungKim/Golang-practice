package main

import (
    "fmt"
    "net/http"
)

type Hello struct{}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello!")
}

func main() {
    var h Hello
    var i Hello
    http.ListenAndServe("localhost:8000", h)
    http.ListenAndServe("localhost:8000", i)
}
