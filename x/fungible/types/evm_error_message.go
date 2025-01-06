package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func EvmErrorMessage(method string, contract common.Address, args interface{}) string {
	return fmt.Sprintf(
		"contract call failed: method '%s', contract '%s', args: %v",
		method,
		contract.Hex(),
		args,
	)
}

func EvmErrorMessageWithRevertError(errorMessage string, reason interface{}) string {
	return fmt.Sprintf(
		"%s, reason: %v",
		errorMessage,
		reason,
	)
}
