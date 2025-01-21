package types

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type EvmErrorMessage struct {
	Message      string `json:"message"`
	Method       string `json:"method"`
	Contract     string `json:"contract"`
	Args         string `json:"args"`
	Error        string `json:"error"`
	RevertReason string `json:"revert_reason"`
}

func NewEvmErrorMessage(method string, contract common.Address, args interface{}, message string) EvmErrorMessage {
	return EvmErrorMessage{
		Method:   method,
		Contract: contract.String(),
		Args:     fmt.Sprintf("%v", args),
		Message:  message,
	}
}

func (e *EvmErrorMessage) AddError(error string) {
	e.Error = error
}

func (e *EvmErrorMessage) AddRevertReason(revertReason interface{}) {
	e.RevertReason = fmt.Sprintf("%v", revertReason)
}

func (e *EvmErrorMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("error marshalling EvmErrorMessage to JSON: %v", err)
	}
	return string(jsonData), nil
}

// Write a function to parse Json string to EvmErrorMessage
func ParseEvmErrorMessage(jsonData string) (EvmErrorMessage, error) {
	var evmErrorMessage EvmErrorMessage
	err := json.Unmarshal([]byte(jsonData), &evmErrorMessage)
	if err != nil {
		return EvmErrorMessage{}, fmt.Errorf("error unmarshalling JSON to EvmErrorMessage: %v", err)
	}
	return evmErrorMessage, nil
}
