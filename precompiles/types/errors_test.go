package types

import "testing"

func Test_ErrInvalidAddr(t *testing.T) {
	e := ErrInvalidAddr{
		Got:    "foo",
		Reason: "bar",
	}
	got := e.Error()
	expect := "invalid address foo, reason: bar"
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

func Test_ErrInvalidCoin(t *testing.T) {
	e := ErrInvalidCoin{
		Got:      "foo",
		Negative: true,
		Nil:      false,
	}
	got := e.Error()
	expect := "invalid coin: denom: foo, is negative: true, is nil: false"
	if got != expect {
		t.Errorf("Expected %v, got %v", expect, got)
	}
}

func Test_ErrInvalidAmount(t *testing.T) {
	e := ErrInvalidAmount{
		Got: "foo",
	}
	got := e.Error()
	expect := "invalid token amount: foo"
	if got != expect {
		t.Errorf("Expected %v, got %v", expect, got)
	}
}

func Test_ErrUnexpected(t *testing.T) {
	e := ErrUnexpected{
		When: "foo",
		Got:  "bar",
	}
	got := e.Error()
	expect := "unexpected foo, got: bar"
	if got != expect {
		t.Errorf("Expected %v, got %v", expect, got)
	}
}
