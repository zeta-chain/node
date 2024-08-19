package types

import "fmt"

type ErrInvalidAddr struct {
	Got string
}

func (e ErrInvalidAddr) Error() string {
	return fmt.Sprintf("invalid address %s", e.Got)
}

type ErrInvalidNumberOfArgs struct {
	Got, Expect int
}

func (e ErrInvalidNumberOfArgs) Error() string {
	return fmt.Sprintf("invalid number of arguments; expected %d; got: %d", e.Expect, e.Got)
}

type ErrInvalidMethod struct {
	Method string
}

func (e ErrInvalidMethod) Error() string {
	return fmt.Sprintf("invalid method: %s", e.Method)
}
