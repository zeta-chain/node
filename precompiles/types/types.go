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
	CallContract(
		ctx sdk.Context,
		fungibleKeeper *fungiblekeeper.Keeper,
		abi *abi.ABI,
		dst common.Address,
		method string,
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

// CallContract calls a contract method on behalf of a precompiled contract.
//   - noEtherumTxEvent is set to true because we don't want to emit EthereumTxEvent,
//     as any MsgEthereumTx with more than one ethereum_tx will fail and the receipt
//     won't be able to be retrieved.
//   - from is set always to the precompiled contract address.
func (c *baseContract) CallContract(
	ctx sdk.Context,
	fungibleKeeper *fungiblekeeper.Keeper,
	abi *abi.ABI,
	dst common.Address,
	method string,
	args []interface{},
) ([]interface{}, error) {
	res, err := fungibleKeeper.CallEVM(
		ctx,             // ctx
		*abi,            // abi
		c.RegistryKey(), // from
		dst,             // to
		big.NewInt(0),   // value
		nil,             // gasLimit
		true,            // commit
		true,            // noEthereumTxEvent
		method,          // method
		args...,         // args
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
