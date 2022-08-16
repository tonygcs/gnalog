# gnalog

Logging library for go applications.

## Install

```sh
go get github.com/tonygcs/gnalog
```

## Usage

```go
l := gnalog.New()
l := l.With("<field_name>", "<field_value>")
l.Debug("debug message")
```

### Formatter

The formatter must implement the `Formatter` interface:

```go
type Formatter interface {
    Format(l *gnalog.Log) ([]byte, error)
}

type formatter struct {
}

func (f *formatter) Format(l *gnalog.Log) ([]byte, error) {
    // Format log
}

func main() {
    gnalog.SetFormatter(&formatter{})
}
```
