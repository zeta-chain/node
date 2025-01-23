package errors

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// CCTXErrorMessage is used to provide a detailed error message for a cctx

// The message and error fields together provide a view on what cause the outbound or revert to fail
// Message : Internal message from the protocol to indicate what went wrong.
// Error : This is the error from the protocol

// Zevm specific fields. These fields are only avalaible if the outbound or revert was involved in a ZEVM transaction, as a deposit
// Contract: Contract called
// Args: Arguments provided to the call
// Method: Contract method
// RevertReason: Reason for the revert if available
type CCTXErrorMessage struct {
	Message      string `json:"message"`
	Error        string `json:"error"`
	Method       string `json:"method"`
	Contract     string `json:"contract"`
	Args         string `json:"args"`
	RevertReason string `json:"revert_reason"`
}

// NewZEVMErrorMessage creates a new CCTXErrorMessage , which contains additional fields which are specific to ZEVM calls only
func NewZEVMErrorMessage(method string, contract common.Address, args interface{}, message string, err error) CCTXErrorMessage {
	return CCTXErrorMessage{
		Method:   method,
		Contract: contract.String(),
		Args:     fmt.Sprintf("%v", args),
		Message:  message,
		Error:    err.Error(),
	}
}

// NewCCTXErrorJsonMessage creates a new CCTXErrorMessage struct and returns it as a JSON string.
// The error field can be of the following types
// 1 - Nil : If the error is nil, we don't need to do anything
// 2 - ErrorString : If the error is a string, we add the error and Pack it into a CCTXErrorMessage struct
// 3 - CCTXErrorMessage : If the error is already a CCTXErrorMessage, we should unpack it and return the JSON string,
// This function does not return an error, if the marshalling fails, it returns a string representation of the CCTXErrorMessage
func NewCCTXErrorJsonMessage(message string, newError error) string {
	m := &CCTXErrorMessage{
		Message: message,
	}

	if newError != nil {
		parsed, err := ParseCCTXErrorMessage(newError.Error())
		switch {
		// parsing failed so this is not a CCTXErrorMessage
		case err != nil:
			{
				m.Error = newError.Error()
			}
		// parsing succeeded, this is a CCTXErrorMessage unpack it
		default:
			{
				m.Message = parsed.Message
				m.Error = parsed.Error
				m.Method = parsed.Method
				m.Contract = parsed.Contract
				m.Args = parsed.Args
				m.RevertReason = parsed.RevertReason
			}
		}
	}

	jsonString, err := m.ToJSON()
	if err != nil {
		return fmt.Sprintf("Json Marshalling failed: %s", m.String())
	}
	return jsonString
}

// NewCCTXError creates a new CCTXErrorMessage struct wraps it in an error
func NewCCTXError(message string, err error) error {
	m := &CCTXErrorMessage{
		Message: message,
	}
	if err != nil {
		m.Error = err.Error()
	}
	jsonString, err := m.ToJSON()
	if err != nil {
		return errors.New(fmt.Sprintf("Json Marshalling failed,message: %s error: %s", m.Message, m.Error))
	}
	return errors.New(jsonString)
}

// AddRevertReason adds a revert reason to the CCTXErrorMessage struct
func (e *CCTXErrorMessage) AddRevertReason(revertReason interface{}) {
	e.RevertReason = fmt.Sprintf("%v", revertReason)
}

func (e *CCTXErrorMessage) Wrap(message string) {
	newMessage := fmt.Sprintf("%s , %s", e.Message, message)
	e.Message = newMessage
}

// ToJSON marshals an CCTXErrorMessage struct into a JSON string
func (e *CCTXErrorMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("error marshalling CCTXErrorMessage to JSON: %v", err)
	}
	return string(jsonData), nil
}

// String returns a string representation of an CCTXErrorMessage struct
func (e *CCTXErrorMessage) String() string {
	return fmt.Sprintf("Message: %s, Error: %s, Method: %s, Contract: %s, Args: %s, RevertReason: %s",
		e.Message, e.Error, e.Method, e.Contract, e.Args, e.RevertReason)
}

// ParseCCTXErrorMessage parses a JSON string into an CCTXErrorMessage struct
func ParseCCTXErrorMessage(jsonData string) (CCTXErrorMessage, error) {
	var evmErrorMessage CCTXErrorMessage
	err := json.Unmarshal([]byte(jsonData), &evmErrorMessage)
	if err != nil {
		return CCTXErrorMessage{}, fmt.Errorf("error unmarshalling JSON to CCTXErrorMessage: %v", err)
	}
	return evmErrorMessage, nil
}
