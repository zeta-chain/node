package types

import "fmt"

/*
	Address related errors
*/

type ErrInvalidAddr struct {
	Got    string
	Reason string
}

func (e ErrInvalidAddr) Error() string {
	return fmt.Sprintf("invalid address %s, reason: %s", e.Got, e.Reason)
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
	return fmt.Sprintf("invalid argument: got %v (type %T)", e.Got, e.Got)
}

/*
	Token related errors
*/

type ErrInvalidToken struct {
	Got    string
	Reason string
}

func (e ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid token %s: %s", e.Got, e.Reason)
}

type ErrInvalidCoin struct {
	Got      string
	Negative bool
	Nil      bool
	Empty    bool
}

func (e ErrInvalidCoin) Error() string {
	return fmt.Sprintf(
		"invalid coin: denom: %s, is negative: %v, is nil: %v, is empty: %v",
		e.Got,
		e.Negative,
		e.Nil,
		e.Empty,
	)
}

type ErrInvalidAmount struct {
	Got string
}

func (e ErrInvalidAmount) Error() string {
	return fmt.Sprintf("invalid token amount: %s", e.Got)
}

type ErrInsufficientBalance struct {
	Requested string
	Got       string
}

func (e ErrInsufficientBalance) Error() string {
	return fmt.Sprintf("insufficient balance: requested %s, current %s", e.Requested, e.Got)
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

type ErrDisabledMethod struct {
	Method string
}

func (e ErrDisabledMethod) Error() string {
	return fmt.Sprintf("method %s is disabled", e.Method)
}

type ErrWriteMethod struct {
	Method string
}

func (e ErrWriteMethod) Error() string {
	return fmt.Sprintf("method not allowed in read-only mode: %s", e.Method)
}

type ErrUnexpected struct {
	When string
	Got  string
}

func (e ErrUnexpected) Error() string {
	return fmt.Sprintf("unexpected error in %s: %s", e.When, e.Got)
}
