package types

import (
	"testing"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetEVMCallerAddress(t *testing.T) {
	t.Run("should raise error when evm is nil", func(t *testing.T) {
		_, mockVMContract := setupMockEVMAndContract(common.Address{})
		caller, err := GetEVMCallerAddress(nil, &mockVMContract)
		require.Error(t, err)
		require.Equal(t, common.Address{}, caller, "address should be zeroed")
	})

	t.Run("should raise error when contract is nil", func(t *testing.T) {
		mockEVM, _ := setupMockEVMAndContract(common.Address{})
		caller, err := GetEVMCallerAddress(&mockEVM, nil)
		require.Error(t, err)
		require.Equal(t, common.Address{}, caller, "address should be zeroed")
	})

	// When contract.CallerAddress == evm.Origin, caller is set to contract.CallerAddress.
	t.Run("when caller address equals origin", func(t *testing.T) {
		mockEVM, mockVMContract := setupMockEVMAndContract(common.Address{})
		caller, err := GetEVMCallerAddress(&mockEVM, &mockVMContract)
		require.NoError(t, err)
		require.Equal(t, common.Address{}, caller, "address should be the same")
	})

	// When contract.CallerAddress != evm.Origin, caller should be set to evm.Origin.
	t.Run("when caller address equals origin", func(t *testing.T) {
		mockEVM, mockVMContract := setupMockEVMAndContract(sample.EthAddress())
		caller, err := GetEVMCallerAddress(&mockEVM, &mockVMContract)
		require.NoError(t, err)
		require.Equal(t, mockEVM.Origin, caller, "address should be evm.Origin")
	})
}

func setupMockEVMAndContract(address common.Address) (vm.EVM, vm.Contract) {
	mockEVM := vm.EVM{
		TxContext: vm.TxContext{
			Origin: address,
		},
	}

	mockVMContract := vm.NewContract(
		contractRef{address: common.Address{}},
		contractRef{address: common.Address{}},
		uint256.NewInt(0),
		0,
	)

	return mockEVM, *mockVMContract
}

type contractRef struct {
	address common.Address
}

func (c contractRef) Address() common.Address {
	return c.address
}
