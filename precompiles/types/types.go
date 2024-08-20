package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/zeta-chain/ethermint/x/evm/statedb"
)

// Interface compliance.
var _ ExtStateDB = (*statedb.StateDB)(nil)
var _ Registrable = (*baseContract)(nil)
var _ BaseContract = (*baseContract)(nil)

// ExtStateDB defines extra methods of statedb to support stateful precompiled contracts.
// It's used to persist changes into the store.
type ExtStateDB interface {
	vm.StateDB
	ExecuteNativeAction(
		contract common.Address,
		converter statedb.EventConverter,
		action func(ctx sdk.Context) error,
	) error
	CacheContext() sdk.Context
}

type Registrable interface {
	RegistryKey() common.Address
}

type BaseContract interface {
	Registrable
}

// A baseContract implements Registrable and BaseContract interfaces.
type baseContract struct {
	address common.Address
}

func NewBaseContract(address common.Address) BaseContract {
	return &baseContract{
		address: address,
	}
}

func (c *baseContract) RegistryKey() common.Address {
	return c.address
}

func BytesToBigInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(data[:])
}
