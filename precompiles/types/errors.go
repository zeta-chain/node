package types

import "fmt"

/*
	Address related errors
*/

type ErrInvalidAddr struct {
	Got string
}

func (e ErrInvalidAddr) Error() string {
	return fmt.Sprintf("invalid address %s", e.Got)
}

/*
	Argument related errors
*/

type ErrInvalidNumberOfArgs struct {
	Got, Expect int
}

func (e ErrInvalidNumberOfArgs) Error() string {
	return fmt.Sprintf("invalid number of arguments; expected %d; got: %d", e.Expect, e.Got)
}

type ErrInvalidArgument struct {
	Got any
}

func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintf("invalid argument: %s", e.Got.(string))
}

/*
	Coin related errors
*/

type ErrInvalidCoin struct {
	Got      string
	Negative bool
	Nil      bool
}

func (e ErrInvalidCoin) Error() string {
	return fmt.Sprintf("invalid coin: denom: %s, is negative: %v, is nil: %v", e.Got, e.Negative, e.Nil)
}

/*
	Method related errors
*/

type ErrInvalidMethod struct {
	Method string
}

func (e ErrInvalidMethod) Error() string {
	return fmt.Sprintf("invalid method: %s", e.Method)
}
