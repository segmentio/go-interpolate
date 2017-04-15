package interpolate

import (
	"bytes"
	"fmt"
	"strings"
)

type Getter interface {
	Get(k string) (string, error)
}

type mapGetter map[string]interface{}

func (m mapGetter) Get(path string) (string, error) {
	var v interface{}
	v = map[string]interface{}(m)

	for _, key := range strings.Split(path, ".") {
		if m, ok := v.(map[string]interface{}); ok {
			v = m[key]
		} else {
			return "", fmt.Errorf("invalid value at path %v", path)
		}
	}

	return v.(string), nil
}

// TODO: unicode
// TODO: more types

// States.
const (
	sLit = iota
	sVar
)

// Literal node.
type literal struct {
	text string
}

// Value of literal node.
func (n *literal) Value(g Getter) (string, error) {
	return n.text, nil
}

// Variable node.
type variable struct {
	path string
}

// Value of interpolated variable.
func (n *variable) Value(g Getter) (string, error) {
	return g.Get(n.path)
}

// Node represents a literal or interpolated node.
type Node interface {
	Value(g Getter) (string, error)
}

// Template represents a series of literal and interpolated nodes.
type Template struct {
	nodes []Node
}

// New template from the given string.
func New(s string) (*Template, error) {
	tmpl := new(Template)
	state := sLit
	buf := new(bytes.Buffer)

	for i := 0; i < len(s); i++ {
		switch state {
		case sLit:
			switch s[i] {
			case '{':
				tmpl.nodes = append(tmpl.nodes, &literal{buf.String()})
				state = sVar
				buf = new(bytes.Buffer)
			default:
				buf.WriteByte(s[i])
			}
		case sVar:
			switch s[i] {
			case '}':
				tmpl.nodes = append(tmpl.nodes, &variable{buf.String()})
				state = sLit
				buf = new(bytes.Buffer)
			default:
				buf.WriteByte(s[i])
			}
		}
	}

	if state == sVar {
		return nil, fmt.Errorf("missing '}'")
	}

	if state == sLit && buf.Len() > 0 {
		tmpl.nodes = append(tmpl.nodes, &literal{buf.String()})
	}

	return tmpl, nil
}

// Eval evalutes the given value against the template,
// returning an error if there's a path mismatch.
func (t *Template) Eval(v interface{}) (string, error) {
	switch v.(type) {
	case map[string]interface{}:
		return t.EvalWith(mapGetter(v.(map[string]interface{})))
	default:
		return t.EvalWith(nil)
	}

}

// Eval evalutes the given value against the template,
// returning an error if there's a path mismatch.
func (t *Template) EvalWith(g Getter) (string, error) {
	var buf bytes.Buffer

	for _, node := range t.nodes {
		ret, err := node.Value(g)
		if err != nil {
			return "", err
		}
		buf.WriteString(ret)
	}

	return buf.String(), nil
}
