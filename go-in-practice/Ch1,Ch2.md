# Ch1

## 1.2 Noteworthy aspects of Go

### 1.2.2 A modern standard library

#### Networking And Http

The Go standard library makes this easy, whether you’re working with HTTP or dealing directly with Transmission Control Protocol ( TCP ), User Datagram Protocol ( UDP ), or other common setups.

```go
package main
import (
    "bufio"
    "fmt"
    "net"
)
func main() {
    conn, _ := net.Dial("tcp", "golang.org:80")
    fmt.Fprintf(conn, "GET/ HTTP/1.0\r\n\r\n")
    status, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Println(status)
}
```

HTTP, Representational State Transfer (REST), and web servers are incredibly common. ***To handle this common case, Go has the http package for providing both a client and a server (see the following listing).*** The client is simple enough to use that it meets the needs of the common everyday cases and extensible enough to use for the complex cases.

```go
package main
import (
    "fmt"
    "io/ioutil"
    "net/http"
)
func main() {
    resp, _ := http.Get("http://example.com/")
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
    resp.Body.Close()
}
```

### 1.2.4 Go the toolchain - more than a language

#### Testing

Go’s naming convention for test files is that they end in _test.go . This suffix tells Go that this is a file to be run when tests execute, and excluded when the application is built, as shown in the next listing.

# Ch2 A solid foundation

```go
package main

import (
	"flag"
    "fmt"
)
var name = flag.String("name", "world", "A name to say hello to.")
var spanish bool

func init() {
    flag.BoolVar(&spanish, "spanish", false, "Use Spanish language.")
    flag.BoolVar(&spanish, "s", false, "Use Spanish language.")
}

func main() {
    flag.Parse()
    if spanish == true{
        fmt.Printf("Hola %s!\n", *name)
    }else {
        fmt.Printf("Hello %s!\n", *name)
    }
}
```

### 2.2 Handling configuration

```go
package main
import (
	"encoding/json"
    "fmt"
    "os"
)
type configuration struct {
    Enabled bool
    Path	string
}

func main() {
    file, _  := os.Open("conf.json")
    defer file.Close()
    
    decoder := json.NewDecoder(file)
    conf := configuration{}
    err := decoder.Decode(&conf)
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Println(conf.Path)
}
```

#### PROBLEM

A program requires configuration above and beyond that which you can provide with a few command-line arguments. And it would be great to use a standard file format.

#### Solution 1

using JSON

#### Solution 2

YAML, a recursive acronym meaning YAML Ain’t Markup Language, is a ***human-readable data serialization format.*** YAML is easy to read, can contain comments, and is fairly easy to work with. Using YAML for application configuration is common and a method we, the authors, recommend and practice. Although Go doesn’t ship with a YAML processor, several third-party libraries are readily available. You’ll look at one here.

#### SOLUTION 3

INI files are a format in wide use and have been around for decades. This is another format your Go applications can potentially use. Although the Go developers didn’t include a processor in the language, once again libraries are readily available to meet
your needs.

### 2.3 Working With real-world web servers

Although the Go standard library provides a great foundation for building web servers, it has some options you may want to change and some tolerance you may want to add. Two common areas, which we cover in this section, are

- Matching URL paths to callback functions 
- Starting and stopping servers with an interest in gracefully shutting down.

Web servers are a core feature of the http package. This package uses the foundation for handling TCP connections from the net package. Because web servers are a core part of the standard library and are commonly in use, simple web servers were introduced in chapter 1. This section moves beyond base web servers and covers some practical gotchas that come up when building applications.

#### 2.3.1 Starting up and shutting down a server

A COMMON ANTIPATTERN : A CALLBACK URL

A simple pattern (or rather antipattern) for development is to have a URL such as
/kill or /shutdown, that will shut down the server when called. The following listing
showcases a simple version of this method.

Graceful shutdowns using manners

When a server shuts down, you’ll often want to stop receiving new requests, save any data to disk, and cleanly end connections with existing open connections. The http package in the standard library shuts down immediately and doesn’t provide an opportunity to handle any of these situations. In the worst cases, this results in lost or corrupted data.

problem

To avoid data loss and unexpected behavior, a server may need to do some cleanup on
shutdown.

SOLUTION
To handle these, you’ll need to implement your own logic or use a package such as
github.com/braintree/manners .

#### 2.3.2 Routing web requests

One of the fundamental tasks of any HTTP server is to receive a given request and map
it to an internal function that can then return a result to the client. ***This routing of a request to a handler is important;*** do it well, and you can build web services that are easily maintainable and flexible enough to fit future needs. This section presents various routing scenarios and solutions for each.

PROBLEM
To correctly route requests, a web server needs to be able to quickly and efficiently
parse the path portion of a URL.

SOLUTION: MULTIPLE HANDLERS
To expand on the method used in listing 1.16, this technique uses a handler function
for each path. This technique, presented in the guide “Writing Web Applications” http://golang.org/doc/articles/wiki/, uses a simple pattern that can be great for web apps with only a few simple paths. This technique has nuances that you’ll see in a moment that may make you consider one of the techniques following it.