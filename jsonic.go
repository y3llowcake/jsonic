package jsonic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type (
	Jsonic struct {
		i interface{}
	}
	Map   map[string]Jsonic
	Array []Jsonic
)

var (
	none = Jsonic{}
)

func New(r io.Reader) (Jsonic, error) {
	j := Jsonic{}
	dec := json.NewDecoder(r)
	dec.UseNumber()
	err := dec.Decode(&(j.i))
	if err != nil {
		return none, err
	}
	return j, nil
}

func MustNew(r io.Reader) Jsonic {
	j, e := New(r)
	if e != nil {
		panic(e)
	}
	return j
}

func NewString(s string) (Jsonic, error) {
	return New(strings.NewReader(s))
}

func MustNewString(s string) Jsonic {
	j, e := NewString(s)
	if e != nil {
		panic(e)
	}
	return j
}

func NewBytes(b []byte) (Jsonic, error) {
	return New(bytes.NewReader(b))
}

func MustNewBytes(b []byte) Jsonic {
	j, e := NewBytes(b)
	if e != nil {
		panic(e)
	}
	return j
}

func (j Jsonic) Type() string {
	switch j.i.(type) {
	case string:
		return "string"
	case json.Number:
		return "number"
	case bool:
		return "bool"
	case map[string]interface{}:
		return "map"
	case []interface{}:
		return "array"
	}
	return "unknown"
}

func (j Jsonic) String() (string, bool) {
	ret, ok := j.i.(string)
	return ret, ok
}

func (j Jsonic) MustString() string {
	if v, ok := j.String(); ok {
		return v
	}
	panic(fmt.Errorf("got %s expected string", j.Type()))
}

func (j Jsonic) Number() (json.Number, bool) {
	ret, ok := j.i.(json.Number)
	return ret, ok
}

func (j Jsonic) MustNumber() json.Number {
	if v, ok := j.Number(); ok {
		return v
	}
	panic(fmt.Errorf("got %s expected number", j.Type()))
}

func (j Jsonic) Bool() (bool, bool) {
	ret, ok := j.i.(bool)
	return ret, ok
}

func (j Jsonic) MustBool() bool {
	if v, ok := j.Bool(); ok {
		return v
	}
	panic(fmt.Errorf("got %s expected bool", j.Type()))
}

func (j Jsonic) Array() (Array, bool) {
	if a, ok := j.i.([]interface{}); ok {
		ja := Array{}
		for _, i := range a {
			ja = append(ja, Jsonic{i: i})
		}
		return ja, true
	}
	return nil, false
}

func (j Jsonic) MustArray() Array {
	if v, ok := j.Array(); ok {
		return v
	}
	panic(fmt.Errorf("got %s expected array", j.Type()))
}

func (j Jsonic) Map() (Map, bool) {
	if m, ok := j.i.(map[string]interface{}); ok {
		jm := Map{}
		for k, i := range m {
			jm[k] = Jsonic{i: i}
		}
		return jm, true
	}
	return nil, false
}

func (j Jsonic) MustMap() Map {
	if v, ok := j.Map(); ok {
		return v
	}
	panic(fmt.Errorf("got %s expected map", j.Type()))
}

func (j Jsonic) At(keys ...string) (Jsonic, bool) {
	l := len(keys)
	if l == 0 {
		return none, false
	}

	jm, ok := j.i.(map[string]interface{})
	if !ok {
		return none, ok
	}
	v, ok := jm[keys[0]]
	if !ok {
		return none, ok
	}
	jv := Jsonic{i: v}
	if l == 1 {
		return jv, ok
	}
	return jv.At(keys[1:]...)
}

func (j Jsonic) MustAt(keys ...string) Jsonic {
	if j, ok := j.At(keys...); ok {
		return j
	}
	panic(fmt.Errorf("invalid path"))
}

/*
func (j Jsonic) MarshalJSONToRawMessage() (*json.RawMessage, error) {
	b, err := j.MarshalJSON()
	if err != nil {
		return nil, err
	}
	rm := json.RawMessage(b)
	return &rm, nil
}

func (j Jsonic) MarshalJSONToString() (string, error) {
	b, err := j.MarshalJSON()
	if err == nil {
		return string(b), err
	}
	return "", err
}

func (j Jsonic) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := j.MarshalJSONToWriter(&buf)
	// Fucking golang json encoder is such a PITA. Adds newlines to the end of
	// objects...
	b := buf.Bytes()
	lastNewline := bytes.LastIndexByte(b, '\n')
	return b[:lastNewline], err
}

func (j Jsonic) MarshalJSONToWriter(w io.Writer) error {
	enc := json.NewEncoder(w)
	// Fucking golang json encoder is a PITA. HTML escapes by default.
	enc.SetEscapeHTML(false)
	return enc.Encode(j.i)
}

// Convience
func (jm Map) Jsonic() Jsonic {
	return NewJsonicInterface(jm)
}

func (ja Array) Jsonic() Jsonic {
	return NewJsonicInterface(ja)
}

func (j Jsonic) MapAt(keys ...string) (Map, bool) {
	if v, ok := j.At(keys...); ok {
		if t, ok := v.Map(); ok {
			return t, true
		}
	}
	return nil, false
}

func (j Jsonic) StringAt(keys ...string) (string, bool) {
	if v, ok := j.At(keys...); ok {
		if s, ok := v.String(); ok {
			return s, true
		}
	}
	return "", false
}

func (j Jsonic) NumberAt(keys ...string) (json.Number, bool) {
	if v, ok := j.At(keys...); ok {
		if n, ok := v.Number(); ok {
			return n, true
		}
	}
	return "0", false
}

func (j Jsonic) BoolAt(keys ...string) (bool, bool) {
	if v, ok := j.At(keys...); ok {
		if b, ok := v.Bool(); ok {
			return b, true
		}
	}
	return false, false
}

func (j Jsonic) ArrayAt(keys ...string) (Array, bool) {
	if v, ok := j.At(keys...); ok {
		if n, ok := v.Array(); ok {
			return n, true
		}
	}
	return nil, false
}

// Super useful for tests. TODO consider moving this into unsafe, or a testing
// package?
func (j Jsonic) DeepEquals(o Jsonic) bool {
	jb, err := j.MarshalJSON()
	if err != nil {
		return false
	}
	ob, err := o.MarshalJSON()
	if err != nil {
		return false
	}
	return bytes.Equal(jb, ob)
}

type LessFunc func(i, j Jsonic) bool

type sortableArray struct {
	a Array
	l LessFunc
}

func (s sortableArray) Len() int {
	return len(s.a)
}
func (s sortableArray) Swap(i, j int) {
	s.a[i], s.a[j] = s.a[j], s.a[i]
}
func (s sortableArray) Less(i, j int) bool {
	return s.l(NewJsonicInterface(s.a[i]), NewJsonicInterface(s.a[j]))
}

func (ja Array) Sort(lf LessFunc) {
	sortable := sortableArray{
		a: ja,
		l: lf,
	}
	sort.Sort(sortable)
}
*/
