package interpolate

import (
	"bytes"
	"fmt"
	"strings"
)

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
func (n *literal) Value(v interface{}) (string, error) {
	return n.text, nil
}

// Variable node.
type variable struct {
	path []string
}

// Get value from interface at the variable's path.
func (n *variable) get(v interface{}) interface{} {
	for _, key := range n.path {
		if m, ok := v.(map[string]interface{}); ok {
			v = m[key]
		} else {
			return nil
		}
	}
	return v
}

// Value of interpolated variable.
func (n *variable) Value(v interface{}) (string, error) {
	v = n.get(v)
	switch v.(type) {
	case string:
		return v.(string), nil
	default:
		path := strings.Join(n.path, ".")
		return "", fmt.Errorf("invalid value at path %v", path)
	}
}

// Node represents a literal or interpolated node.
type Node interface {
	Value(v interface{}) (string, error)
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
				path := strings.Split(buf.String(), ".")
				tmpl.nodes = append(tmpl.nodes, &variable{path})
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
	var buf bytes.Buffer

	for _, node := range t.nodes {
		ret, err := node.Value(v)
		if err != nil {
			return "", err
		}
		buf.WriteString(ret)
	}

	return buf.String(), nil
}
