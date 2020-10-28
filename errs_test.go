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
		log.Println(E(op, e).Error())
	} else {
		fmt.Printf("[%s] answer: %d\n", op, eans)
	}

	want := errors.New("can't divide by zero")
	got := e.Err
	if want.Error() != got.Error() {
		t.Errorf("want %s got %s | errs %#v | fail", want, got, e)
	}
}

type TestStr struct {
	A int
	B string
}
type TestStrs []TestStr

func (s TestStrs) Len() int {
	return len(s)
}

func TestMultiStack(t *testing.T) {

	tArgs := TestStrs{
		{A: 1, B: "hello"},
		{A: 2, B: "world"},
	}

	fourth := func() *Error {
		err := New("fourth", errors.New("in fourth"), "4th", tArgs)
		log.Println(err.Error())
		return err
	}

	third := func() *Error {
		op := Op("third")
		err := fourth()
		return E(op, err)
	}

	second := func() *Error {
		op := Op("second")
		err := third()
		log.Println(err.Error())
		return E(op, err)
	}

	first := func() *Error {
		op := Op("first")
		err := second()
		log.Println(err.Error())
		return E(op, err)
	}

	op := Op("TestMultiStack")
	err := first()
	e := E(op, err)
	log.Println(e.Error())

	want := errors.New("in fourth")
	got := e.Err

	if want.Error() != got.Error() {
		t.Errorf("want %s got %s | err %#v | fail", want, got, e)
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
			errors.New("can't divide by zero"),
			"DIVERR",
			Args{a, b, c},
		)
		return 0, err
	}
	return a / b / c, nil
}
