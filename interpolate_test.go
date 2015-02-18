package interpolate

import "github.com/bmizerany/assert"
import "testing"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestParseBroken(t *testing.T) {
	_, err := New(`Hello {name`)
	assert.Equal(t, `missing '}'`, err.Error())
}

func TestEval(t *testing.T) {
	tmpl, err := New(`Hello {name}`)
	assert.Equal(t, nil, err)

	s, err := tmpl.Eval(map[string]interface{}{"name": "Tobi"})
	assert.Equal(t, nil, err)
	assert.Equal(t, `Hello Tobi`, s)
}

func TestEvalMany(t *testing.T) {
	tmpl, err := New(`Hello {first} {last}`)
	assert.Equal(t, nil, err)

	s, err := tmpl.Eval(map[string]interface{}{
		"first": "Tobi",
		"last":  "Ferret",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, `Hello Tobi Ferret`, s)
}

func TestEvalManyNested(t *testing.T) {
	tmpl, err := New(`Hello {name.first} {name.last} you are {color}`)
	assert.Equal(t, nil, err)

	s, err := tmpl.Eval(map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Tobi",
			"last":  "Ferret",
		},
		"color": "Albino",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, `Hello Tobi Ferret you are Albino`, s)
}

func TestEvalTooShort(t *testing.T) {
	tmpl, err := New(`Hello {name.first} {name.last} you are {color}`)
	assert.Equal(t, nil, err)

	_, err = tmpl.Eval(map[string]interface{}{
		"name":  "Tobi Ferret",
		"color": "Albino",
	})

	assert.Equal(t, `invalid value at path name.first`, err.Error())
}

func TestEvalTooLong(t *testing.T) {
	tmpl, err := New(`Hello {name.first.whatever.stuff.here.whoop} {name.last} you are {color}`)
	assert.Equal(t, nil, err)

	s, err := tmpl.Eval(map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Tobi",
			"last":  "Ferret",
		},
		"color": "Albino",
	})

	assert.Equal(t, "", s)
	assert.Equal(t, `invalid value at path name.first.whatever.stuff.here.whoop`, err.Error())
}

func TestEvalTrailingLit(t *testing.T) {
	tmpl, err := New(`stream:project:{projectId}:ingress`)
	assert.Equal(t, nil, err)

	s, err := tmpl.Eval(map[string]interface{}{"projectId": "1234"})
	assert.Equal(t, nil, err)
	assert.Equal(t, `stream:project:1234:ingress`, s)
}

func BenchmarkParse(t *testing.B) {
	for i := 0; i < t.N; i++ {
		New(`Hello {name.first} {name.last} you are {color}`)
	}
}

func BenchmarkEval(t *testing.B) {
	tmpl, _ := New(`Hello {name.first} {name.last} you are {color}`)

	v := map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Tobi",
			"last":  "Ferret",
		},
		"color": "Albino",
	}

	for i := 0; i < t.N; i++ {
		tmpl.Eval(v)
	}
}
