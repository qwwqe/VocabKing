package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Op string

type Kind int

type Meta map[string]string

const (
	KindUnknown Kind = iota
	KindDatabaseFailure
	KindUnauthorized
	KindBadRequest
	KindBadForm
)

var kindName = map[Kind]string{
	KindUnknown:         "unknown",
	KindDatabaseFailure: "database failure",
	KindUnauthorized:    "unauthorized",
	KindBadRequest:      "bad request",
	KindBadForm:         "bad form",
}

var kindStatusCode = map[Kind]int{
	KindUnknown:         http.StatusInternalServerError,
	KindDatabaseFailure: http.StatusInternalServerError,
	KindUnauthorized:    http.StatusUnauthorized,
	KindBadRequest:      http.StatusBadRequest,
	KindBadForm:         http.StatusBadRequest,
}

type opErr struct {
	op   Op
	kind Kind
	meta Meta
	err  error
}

func (m Meta) String() string {
	if m == nil || len(m) == 0 {
		return ""
	}
	s := ""
	for k, v := range m {
		s += k + ": " + v + ", "
	}
	return s[:len(s)-2]
}

func (m Meta) Set(k, v string) {
	if m == nil {
		return
	}
	m[k] = v
}

func (k Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Name())
}

func (k Kind) Name() string {
	return kindName[k]
}

func (k Kind) StatusCode() int {
	return kindStatusCode[k]
}

func New(o Op, k Kind, e error, m Meta) *opErr {
	return &opErr{
		op:   o,
		kind: k,
		err:  e,
		meta: m,
	}
}

func NewFromString(o Op, k Kind, msg string, m Meta) *opErr {
	return New(o, k, errors.New(msg), m)
}

func NilOrNew(o Op, k Kind, err error, m Meta) *opErr {
	if err != nil {
		return New(o, k, err, m)
	}
	return nil
}

func (e *opErr) Error() string {
	meta := ""
	for k, v := range e.meta {
		meta += k + ": " + v + ", "
	}
	return string(e.op) + ": " + e.kind.Name() + ": " + e.err.Error() + " (" + e.meta.String() + ")"
}

func (e *opErr) Kind() Kind {
	return e.kind
}

func (e *opErr) Meta() Meta {
	return e.meta
}

func (e *opErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Error string   `json:"error"`
		Ops   []Op     `json:"ops"`
		Stack []string `json:"stack"`
	}{
		Error: e.Kind().Name(),
		Ops:   Ops(e),
		Stack: Stack(e),
	})
}

func Ops(e *opErr) []Op {
	ops := []Op{e.op}
	subErr, ok := e.err.(*opErr)
	if !ok {
		return ops
	}
	return append(ops, Ops(subErr)...)
}

func Stack(e *opErr) []string {
	stack := []string{e.Error()}
	subErr, ok := e.err.(*opErr)
	if !ok {
		return stack
	}
	return append(stack, Stack(subErr)...)
}
