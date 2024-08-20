package types

import "testing"

func Test_ErrInvalidAddr(t *testing.T) {
	e := ErrInvalidAddr{
		Got: "foo",
	}
	got := e.Error()
	expect := "invalid address foo"
	if got != expect {
		t.Errorf("Expected %v, got %v", expect, got)
	}
}

func Test_ErrInvalidNumberOfArgs(t *testing.T) {
	e := ErrInvalidNumberOfArgs{
		Got:    1,
		Expect: 2,
	}
	got := e.Error()
	expect := "invalid number of arguments; expected 2; got: 1"
	if got != expect {
		t.Errorf("Expected %v, got %v", expect, got)
	}
}

func Test_ErrInvalidArgument(t *testing.T) {
	e := ErrInvalidArgument{
		Got: "foo",
	}
	got := e.Error()
	expect := "invalid argument: foo"
	if got != expect {
		t.Errorf("Expected %v, got %v", expect, got)
	}
}

func Test_ErrInvalidMethod(t *testing.T) {
	e := ErrInvalidMethod{
		Method: "foo",
	}
	got := e.Error()
	expect := "invalid method: foo"
	if got != expect {
		t.Errorf("Expected %v, got %v", expect, got)
	}
}
