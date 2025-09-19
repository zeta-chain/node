package errors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Type string

const (
	ContractCallError Type = "contract_call_error"
	InternalError     Type = "internal_error"
)

// CCTXErrorMessage is used to provide a detailed error message for a cctx

// The message and error fields together provide a view on what cause the outbound or revert to fail
// Message : Internal message from the protocol to indicate what went wrong.
// Error : This is the error from the protocol

// Zevm specific fields. These fields are only available if the outbound or revert was involved in a ZEVM transaction, as a deposit
// Contract: Contract called
// Args: Arguments provided to the call
// Method: Contract method
// RevertReason: Reason for the revert if available
type CCTXErrorMessage struct {
	Type         Type   `json:"type,omitempty"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
	Method       string `json:"method,omitempty"`
	Contract     string `json:"contract,omitempty"`
	Args         string `json:"args,omitempty"`
	RevertReason string `json:"revert_reason,omitempty"`
}

// NewCCTXErrorMessage creates a new CCTXErrorMessage struct
func NewCCTXErrorMessage(message string) CCTXErrorMessage {
	return CCTXErrorMessage{
		Message: message,
		Type:    InternalError,
	}
}

// NewZEVMErrorMessage is s special case of NewCCTXErrorMessage which is called by ZEVM specific code.
// This creates a CCTXErrorMessage with ZEVM specific fields like Method, Contract, Args
func NewZEVMErrorMessage(
	method string,
	contract common.Address,
	args interface{},
	message string,
	err error,
) CCTXErrorMessage {
	c := CCTXErrorMessage{
		Method:   method,
		Contract: contract.String(),
		Args:     fmt.Sprintf("%v", args),
		Message:  message,
		Type:     ContractCallError,
	}
	if err != nil {
		c.Error = err.Error()
	}
	return c
}

// NewCCTXErrorJSONMessage creates a new CCTXErrorMessage struct and returns it as a JSON string.
// The error field can be of the following types
// 1 - Nil : If the error is nil, we don't need to do anything
// 2 - ErrorString : If the error is a string, we add the error and Pack it into a CCTXErrorMessage struct
// 3 - CCTXErrorMessage : If the error is already a CCTXErrorMessage, we should unpack it and return the JSON string,
// 4 - If it the error is a chain of errors, we should unpack each error and add it to the CCTXErrorMessage struct error field
// 5 - If the error is a chain of errors and one of the errors is a CCTXErrorMessage, we should unpack the CCTXErrorMessage and return the JSON string and add the errors into the error field
// This function does not return an error, if the marshalling fails, it returns a string representation of the CCTXErrorMessage
func NewCCTXErrorJSONMessage(message string, newError error) string {
	m := NewCCTXErrorMessage(message)

	if newError != nil {
		// 1. Split the error message into parts where each part is a separate error message. Json or a simple string
		errorLists := SplitErrorMessage(newError.Error())

		for _, e := range errorLists {
			//
			parsed, err := ParseCCTXErrorMessage(e)
			switch {
			// 3.  parsing failed so this is not a CCTXErrorMessage json string. Wrap the error into the CCTXErrorMessage.Error field
			case err != nil:
				{
					m.WrapError(e)
				}
			// parsing succeeded, this is a CCTXErrorMessage unpack it. Wrap the error into the CCTXErrorMessage.Error field and assign the other fields
			// The message fields are overwritten , it should contain only a single message describing the error. The actual error chain is added to the error field
			// the method, contract, args, revert_reason fields are present only in ZEVM errors and are directly added to the CCTXErrorMessage struct
			default:
				{
					m.Message = parsed.Message
					m.Method = parsed.Method
					m.Contract = parsed.Contract
					m.Args = parsed.Args
					m.RevertReason = parsed.RevertReason
					m.Type = parsed.Type
					m.WrapError(parsed.Error)
				}
			}
		}
	}

	jsonString, err := m.ToJSON()
	if err != nil {
		return fmt.Sprintf("json marshalling failed %s,%s", err.Error(), m.String())
	}
	return jsonString
}

// SplitErrorMessage splits the error message into parts it treats a JSON formatted error message as a single part

// Example inputs that this function can handle
// 1. "errorString1:errorString2:errorString3" : deposit error: can't call a non-contract address, processing error: withdraw amount 100 is less than dust amount 1000: invalid withdrawal amount

// 2. "errorString1:{jsonObject}" : deposit error: {"message":"contract call failed when calling EVM with data","error":"execution reverted: ret 0x: evm transaction execution failed",
// "method":"depositAndCall0","contract":"0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9,
// args: "[{[] 0xdFb74337c53141bf912101b0Ee770FA8e2DCB921 1337} 0x13A0c5930C028511Dc02665E7285134B6d11A5f4 10000000000000000
// 0xE83192f6301d090EFB2F38824A12E98877F66fe3 [101 53 52 55 51 50 100 55 50 55 57 99 56 97 99 100 48 97 97 101 55 57 98 51 101 100 99 101 50 55 99 99 50 57 97 49 48 51 56 100 48 102 102
// 54 52 57 98 53 101 51 56 101 55 51 53 56 55 101 100 55 56 102 49 48 101 101 97 49 52 56 100 97 55 97 51 97 57 100 98 49 101 101 98 55 51 57 100 49 52 54 53 56 55 99 101 99 48 56 54 55]]","revert_reason":""}
func SplitErrorMessage(input string) []string {
	var result []string

	// Find the JSON object boundaries
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")

	// The error chain does not have a json formatted error message so return the error as is
	if start == -1 || end == -1 {
		return strings.Split(input, ":")
	}

	// We have a JSON formatted error message
	// Extract JSON part and add to result
	jsonStr := input[start : end+1]
	result = append(result, jsonStr)

	// Extract errors before JSON part
	prefix := strings.TrimSpace(input[:start])
	if prefix != "" && strings.HasSuffix(prefix, ":") {
		prefix = strings.TrimSuffix(prefix, ":")
		prefixParts := strings.Split(prefix, ":")

		if len(prefixParts) > 0 {
			prefixErrors := make([]string, len(prefixParts))
			for i, part := range prefixParts {
				prefixErrors[i] = strings.TrimSpace(part)
			}
			result = append(prefixErrors, result...)
		}
	}

	// Extract errors after the json part
	if end+1 < len(input) {
		suffix := strings.TrimSpace(input[end+1:])
		if after, ok := strings.CutPrefix(suffix, ":"); ok {
			suffix = strings.TrimSpace(after)
			suffixParts := strings.Split(suffix, ":")
			if suffix != "" && len(suffixParts) > 0 {
				for _, part := range suffixParts {
					result = append(result, strings.TrimSpace(part))
				}
			}
		}
	}

	return result
}

// WrapError adds a new error to the CCTXErrorMessage struct
func (e *CCTXErrorMessage) WrapError(newError string) {
	existingError := e.Error
	if existingError == "" {
		e.Error = newError
	} else {
		e.Error = fmt.Sprintf("%s:%s", existingError, newError)
	}
}

// AddRevertReason adds a revert reason to the CCTXErrorMessage struct
func (e *CCTXErrorMessage) AddRevertReason(revertReason interface{}) {
	e.RevertReason = fmt.Sprintf("%v", revertReason)
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
