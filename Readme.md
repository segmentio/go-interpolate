# interpolate

 Simple interpolation templates ex: `Hello {name.first}`.

## Usage

#### type Node

```go
type Node interface {
	Value(v interface{}) (string, error)
}
```

Node represents a literal or interpolated node.

#### type Template

```go
type Template struct {
	nodes []Node
}
```

Template represents a series of literal and interpolated nodes.

#### func  New

```go
func New(s string) (*Template, error)
```
New template from the given string.

#### func (*Template) Eval

```go
func (t *Template) Eval(v interface{}) (string, error)
```
Eval evalutes the given value against the template, returning an error if
there's a path mismatch.
