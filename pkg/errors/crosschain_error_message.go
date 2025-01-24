package errors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
	Message      string `json:"message"`
	Error        string `json:"error"`
	Method       string `json:"method"`
	Contract     string `json:"contract"`
	Args         string `json:"args"`
	RevertReason string `json:"revert_reason"`
}

// NewZEVMErrorMessage creates a new CCTXErrorMessage , which contains additional fields which are specific to ZEVM calls only
func NewZEVMErrorMessage(
	method string,
	contract common.Address,
	args interface{},
	message string,
	err error,
) CCTXErrorMessage {
	return CCTXErrorMessage{
		Method:   method,
		Contract: contract.String(),
		Args:     fmt.Sprintf("%v", args),
		Message:  message,
		Error:    err.Error(),
	}
}

// NewCCTXErrorJSONMessage creates a new CCTXErrorMessage struct and returns it as a JSON string.
// The error field can be of the following types
// 1 - Nil : If the error is nil, we don't need to do anything
// 2 - ErrorString : If the error is a string, we add the error and Pack it into a CCTXErrorMessage struct
// 3 - CCTXErrorMessage : If the error is already a CCTXErrorMessage, we should unpack it and return the JSON string,
// This function does not return an error, if the marshalling fails, it returns a string representation of the CCTXErrorMessage
func NewCCTXErrorJSONMessage(message string, newError error) string {
	m := &CCTXErrorMessage{
		Message: message,
	}

	if newError != nil {
		errorLists := SplitErrorMessage(newError.Error())

		for _, e := range errorLists {
			parsed, err := ParseCCTXErrorMessage(e)
			switch {
			// parsing failed so this is not a CCTXErrorMessage json string
			case err != nil:
				{
					m.WrapError(e)
				}
			// parsing succeeded, this is a CCTXErrorMessage unpack it
			default:
				{
					fmt.Println("Parsing succeeded handle as CCTXErrorMessage")
					m.Message = parsed.Message
					m.Method = parsed.Method
					m.Contract = parsed.Contract
					m.Args = parsed.Args
					m.RevertReason = parsed.RevertReason
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
		if strings.HasPrefix(suffix, ":") {
			suffix = strings.TrimSpace(strings.TrimPrefix(suffix, ":"))
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
