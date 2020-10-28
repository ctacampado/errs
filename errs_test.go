package errs

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func TestSingleStack(t *testing.T) {
	op := Op("TestSingleStack")
	eans, e := intDiv(100, 20, 0)
	if e != nil {
		log.Println(fmt.Errorf("division error: %w", errors.New(e.Err)))
	} else {
		fmt.Printf("[%s] answer: %d\n", op, eans)
	}

	want := "can't divide by zero"
	got := e.Err
	if want != got {
		t.Errorf("want %s got %s | errs %#v | fail", want, got, e)
	}
}

type Args []int

func (a Args) Len() int {
	return len(a)
}

func intDiv(a, b, c int) (int, *Error) {
	if b == 0 || c == 0 {
		err := New(
			"intDiv",
			"can't divide by zero",
			"DIVERR",
			Args{a, b, c},
		)
		return 0, err
	}
	return a / b / c, nil
}
