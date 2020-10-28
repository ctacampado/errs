package errs

import (
	"encoding/json"
	"fmt"
)

// Op is the function name where the error happened
type Op string

// OpStack is stack of operations that has happened
// when the error occured
type OpStack []Op

// ErrType is the error type string
type ErrType string

// ErrArgs is the argument involved in the error
type ErrArgs interface {
	Len() int
}

// Error is a custom error type that enables us to
// generate rich error messages for easier debugging,
// logging, and querying.
type Error struct {
	Op      Op
	OpStack OpStack
	Err     error
	Kind    ErrType
	Args    ErrArgs
}

// New returns a new custom error
func New(o Op, e ErrString, k ErrType, a ErrArgs) *Error {
	return &Error{
		Op:      o,
		OpStack: OpStack{o},
		Err:     e,
		Kind:    k,
		Args:    a,
	}
}

// Error returns an error in json string format
func (e Error) Error() error {
	jsonStack, err := json.Marshal(e.OpStack)
	if err != nil {
		jsonStack = []byte(err.Error())
	}

	jsonArgs, _ := json.Marshal(e.Args)
	if err != nil {
		jsonArgs = []byte(err.Error())
	}

	return fmt.Errorf(
		"{\"stack\":%v,\"type\":\"%s\",\"err\":\"%w\",\"args\":%v}",
		string(jsonStack),
		e.Kind,
		e.Err,
		string(jsonArgs),
	)
}

// StackTrace contains operation stack
func (e Error) StackTrace() string {
	return fmt.Sprintf("%+v", e.OpStack)
}

// E returns an errs.Error
// args order should be: OpStack, Op, ErrString, ErrType, *Error, ErrArgs
func E(args ...interface{}) *Error {
	e := &Error{}
	if len(args) == 0 {
		return e
	}
	for _, arg := range args {
		switch a := arg.(type) {
		case OpStack:
			e.OpStack = a
		case Op:
			e.Op = a
			e.OpStack = append(e.OpStack, "_")
			copy(e.OpStack[1:], e.OpStack[0:])
			e.OpStack[0] = e.Op
		case error:
			e.Err = a
		case ErrType:
			e.Kind = a
		case *Error:
			e = E(a.OpStack, e.Op, a.Kind, a.Err, a.Args)
		case ErrArgs:
			e.Args = a
		}
	}
	return e
}
