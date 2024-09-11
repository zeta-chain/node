package bank

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	ptypes "github.com/zeta-chain/node/precompiles/types"
)

// CallContract calls a given contract on behalf of the precompiled contract.
// Note that the precompile contract address is hardcoded.
func (c *Contract) CallContract(
	ctx sdk.Context,
	abi *abi.ABI,
	dst common.Address,
	method string,
	args []interface{},
) ([]interface{}, error) {
	res, err := c.fungibleKeeper.CallEVM(
		ctx,             // ctx
		*abi,            // abi
		ContractAddress, // from
		dst,             // to
		big.NewInt(0),   // value
		nil,             // gasLimit
		true,            // commit
		false,           // noEthereumTxEvent
		method,          // method
		args...,         // args
	)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "CallEVM " + method,
			Got:  err.Error(),
		}
	}

	if res.VmError != "" {
		return nil, &ptypes.ErrUnexpected{
			When: "VmError " + method,
			Got:  res.VmError,
		}
	}

	ret, err := abi.Methods[method].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "Unpack " + method,
			Got:  err.Error(),
		}
	}

	return ret, nil
}
