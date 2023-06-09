# interpolate

> **Note**  
> Segment has paused maintenance on this project, but may return it to an active status in the future. Issues and pull requests from external contributors are not being considered, although internal contributions may appear from time to time. The project remains available under its open source license for anyone to use.

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
