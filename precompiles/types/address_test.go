package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetEVMCallerAddress(t *testing.T) {
	mockEVM := vm.EVM{
		TxContext: vm.TxContext{
			Origin: common.Address{},
		},
	}

	mockVMContract := vm.NewContract(
		contractRef{address: common.Address{}},
		contractRef{address: common.Address{}},
		big.NewInt(0),
		0,
	)

	// When contract.CallerAddress == evm.Origin, caller is set to contract.CallerAddress.
	caller, err := GetEVMCallerAddress(&mockEVM, mockVMContract)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, caller, "address shouldn be the same")

	// When contract.CallerAddress != evm.Origin, caller should be set to evm.Origin.
	mockEVM.Origin = sample.EthAddress()
	caller, err = GetEVMCallerAddress(&mockEVM, mockVMContract)
	require.NoError(t, err)
	require.Equal(t, mockEVM.Origin, caller, "address should be evm.Origin")
}

type contractRef struct {
	address common.Address
}

func (c contractRef) Address() common.Address {
	return c.address
}
