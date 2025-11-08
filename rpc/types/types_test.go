package types_test

import (
	"maps"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/cosmos/evm/x/vm/statedb"
	"github.com/cosmos/evm/x/vm/types/mocks"
	rpc "github.com/zeta-chain/node/rpc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type precompileContract struct{}

func (p *precompileContract) Address() common.Address { return common.Address{} }

func (p *precompileContract) RequiredGas(input []byte) uint64 { return 0 }

func (p *precompileContract) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	return nil, nil
}

func TestApply(t *testing.T) {
	emptyTxConfig := statedb.NewEmptyTxConfig()
	db := statedb.New(sdk.Context{}, mocks.NewEVMKeeper(), emptyTxConfig)
	precompiles := map[common.Address]vm.PrecompiledContract{
		common.BytesToAddress([]byte{0x1}): &precompileContract{},
		common.BytesToAddress([]byte{0x2}): &precompileContract{},
	}
	bytes2Addr := func(b []byte) *common.Address {
		a := common.BytesToAddress(b)
		return &a
	}
	testCases := map[string]struct {
		overrides           *rpc.StateOverride
		expectedPrecompiles map[common.Address]struct{}
		fail                bool
	}{
		"move to already touched precompile": {
			overrides: &rpc.StateOverride{
				common.BytesToAddress([]byte{0x1}): {
					Code:             &hexutil.Bytes{0xff},
					MovePrecompileTo: bytes2Addr([]byte{0x2}),
				},
				common.BytesToAddress([]byte{0x2}): {
					Code: &hexutil.Bytes{0x00},
				},
			},
			fail: true,
		},
		"move non-precompile": {
			overrides: &rpc.StateOverride{
				common.BytesToAddress([]byte{0x1}): {
					Code:             &hexutil.Bytes{0xff},
					MovePrecompileTo: bytes2Addr([]byte{0xff}),
				},
				common.BytesToAddress([]byte{0x3}): {
					Code:             &hexutil.Bytes{0x00},
					MovePrecompileTo: bytes2Addr([]byte{0xfe}),
				},
			},
			fail: true,
		},
		"move two precompiles": {
			overrides: &rpc.StateOverride{
				common.BytesToAddress([]byte{0x1}): {
					Code:             &hexutil.Bytes{0xff},
					MovePrecompileTo: bytes2Addr([]byte{0xff}),
				},
				common.BytesToAddress([]byte{0x2}): {
					Code:             &hexutil.Bytes{0x00},
					MovePrecompileTo: bytes2Addr([]byte{0xfe}),
				},
			},
			expectedPrecompiles: map[common.Address]struct{}{
				common.BytesToAddress([]byte{0xfe}): {},
				common.BytesToAddress([]byte{0xff}): {},
			},
			fail: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cpy := maps.Clone(precompiles)
			err := tc.overrides.Apply(db, cpy)
			if tc.fail {
				if err == nil {
					t.Errorf("%s: want error, have nothing", name)
				}
				return
			}
			if err != nil {
				t.Errorf("%s: want no error, have %v", name, err)
				return
			}
			if len(cpy) != len(tc.expectedPrecompiles) {
				t.Errorf("%s: precompile mismatch, want %d, have %d", name, len(tc.expectedPrecompiles), len(cpy))
			}
			for k := range tc.expectedPrecompiles {
				if _, ok := cpy[k]; !ok {
					t.Errorf("%s: precompile not found: %s", name, k.String())
				}
			}
		})
	}
}
