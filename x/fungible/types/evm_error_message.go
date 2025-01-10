package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

const (
	MethodKey       = "method"
	ContractKey     = "contract"
	ArgsKey         = "args"
	ErrorKey        = "error"
	RevertReasonKey = "revertReason"
)

func EvmErrorMessage(method string, contract common.Address, args interface{}) string {
	return fmt.Sprintf(
		"%s:%s,%s:%s,%s:%v",
		MethodKey,
		method,
		ContractKey,
		contract.Hex(),
		ArgsKey,
		args)
}

func EvmErrorMessageAddErrorString(errorMessage string, error string) string {
	return fmt.Sprintf(
		"%s,%s:%s",
		errorMessage,
		ErrorKey,
		error)
}

func EvmErrorMessageAddRevertReason(errorMessage string, revertReason interface{}) string {
	return fmt.Sprintf(
		"%s,%s:%v",
		errorMessage,
		RevertReasonKey,
		revertReason)
}
