# lexoffice API client

This is a fork of <https://github.com/hostwithquantum/golexoffice>.

It was customized to bit a bit more ergonomic and to use contexts.

Error handling and code was also made more idiomatic. See [pkg.go.dev](https://pkg.go.dev/github.com/karitham/go-lexoffice) for docs.

The library is not 100% feature complete. I'm adding things I need. Feel free to contribute, it is very easy to add new endpoints.

## Install

```console
go get github.com/karitham/go-lexoffice
```

## Useage

```go
import (
    "fmt"
    "log"
    "os"

    lexoffice "github.com/karitham/go-lexoffice"
)

lc := lexoffice.NewClient(os.Getenv("LEXOFFICE_API_KEY"))
contacts, err := lc.GetContacts(context.TODO(), lexoffice.GetContactsParams{
    Name: "hitori gotoh"
})
if err != nil {
    log.Println(err.Error())
}

fmt.Println(contacts)
```
