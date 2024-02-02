package types_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestIsContractReverted(t *testing.T) {
	assert.True(t, types.IsContractReverted(nil, vm.ErrExecutionReverted))
	assert.True(t, types.IsContractReverted(nil, fmt.Errorf("foo : %s", vm.ErrExecutionReverted.Error())))
	assert.True(t, types.IsContractReverted(&evmtypes.MsgEthereumTxResponse{VmError: "foo"}, nil))

	assert.False(t, types.IsContractReverted(nil, nil))
	assert.False(t, types.IsContractReverted(nil, errors.New("foo")))
	assert.False(t, types.IsContractReverted(&evmtypes.MsgEthereumTxResponse{VmError: ""}, nil))
}
