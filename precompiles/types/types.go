package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/zeta-chain/ethermint/x/evm/statedb"

	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
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

type ContractCaller interface {
	CallContract(ctx sdk.Context,
		fungibleKeeper *fungiblekeeper.Keeper,
		abi *abi.ABI,
		from common.Address,
		dst common.Address,
		method string,
		noEthereumTxEvent bool,
		args []interface{}) ([]interface{}, error)
}

type BaseContract interface {
	Registrable
	ContractCaller
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

func (c *baseContract) CallContract(
	ctx sdk.Context,
	fungibleKeeper *fungiblekeeper.Keeper,
	abi *abi.ABI,
	from common.Address,
	dst common.Address,
	method string,
	noEthereumTxEvent bool,
	args []interface{},
) ([]interface{}, error) {
	res, err := fungibleKeeper.CallEVM(
		ctx,               // ctx
		*abi,              // abi
		from,              // from
		dst,               // to
		big.NewInt(0),     // value
		nil,               // gasLimit
		true,              // commit
		noEthereumTxEvent, // noEthereumTxEvent
		method,            // method
		args...,           // args
	)
	if err != nil {
		return nil, &ErrUnexpected{
			When: "CallEVM " + method,
			Got:  err.Error(),
		}
	}

	if res.VmError != "" {
		return nil, &ErrUnexpected{
			When: "VmError " + method,
			Got:  res.VmError,
		}
	}

	ret, err := abi.Methods[method].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, &ErrUnexpected{
			When: "Unpack " + method,
			Got:  err.Error(),
		}
	}

	return ret, nil
}
