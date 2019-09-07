package errors

type (
	Op   string
	Kind int
)

type Error struct {
	Op   Op
	Kind Kind
	Err  error
}

const (
	KindUnknown Kind = iota
	KindDatabaseFailure
)

func New(op Op, kind Kind, err error) *Error {
	return &Error{
		Op:   op,
		Kind: kind,
		Err:  err,
	}
}

func NilOrNew(op Op, kind Kind, err error) *Error {
	if err != nil {
		return New(op, kind, err)
	}
	return nil
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func Ops(e *Error) []Op {
	o := []Op{e.Op}

	subErr, ok := e.Err.(*Error)
	if !ok {
		return o
	}

	return append(o, Ops(subErr)...)
}
