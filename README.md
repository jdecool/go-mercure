go-mercure
==========

go-mercure is a Go client library for [Mercure](https://mercure.rocks).

## Usage

```go
import "github.com/jdecool/go-mercure/mercure"
```

### Publish an event

```go
u, _ := url.Parse("https://localhost:3000/.well-known/mercure")
jwt := "my-jwt"

m := mercure.Message{
    Topics: []string{
        "http://localhost/test",
    },
    Data:  []byte("Hello World"),
}

p := mercure.NewPublisher(u, jwt, nil)
evtId, err := p.Publish(m)
```
