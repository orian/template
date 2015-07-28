# template
A Golang html/template wrapper to allow lazy loading of templates.

It's a drop in replacement for a simple use cases of `html/template` package. Just replace import by:
``` go
import "github.com/orian/template"
```

and later turn the debug mode:
```go
template.Debug = true
```

It was a quick experiment when I was playing with some simple webapp written in Go.
