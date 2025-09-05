package types

import (
	"errors"
	"strings"

	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

var (
	// ErrOnAbortFailed is the error message for failed onAbort
	// It is used to handle case where abort amount is refunded but onAbort is not implemented
	ErrOnAbortFailed = errors.New("onAbort failed")
)

// IsRevertError checks if an error is a evm revert error
func IsRevertError(err error) bool {
	return err != nil && strings.Contains(err.Error(), vm.ErrExecutionReverted.Error())
}

// IsContractReverted checks if the contract execution is reverted
func IsContractReverted(res *evmtypes.MsgEthereumTxResponse, err error) bool {
	return IsRevertError(err) || (res != nil && res.Failed())
}
