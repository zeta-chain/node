package keeper

import (
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/systemcontract.sol"

	"github.com/zeta-chain/node/x/fungible/types"
)

// SetGasPrice sets gas price on the system contract in zEVM; return the gasUsed and error code
func (k Keeper) SetGasPrice(ctx sdk.Context, chainid *big.Int, gasPrice *big.Int) (uint64, error) {
	if gasPrice == nil {
		return 0, cosmoserrors.Wrapf(types.ErrNilGasPrice, "gas price param should be set")
	}
	system, found := k.GetSystemContract(ctx)
	if !found {
		return 0, cosmoserrors.Wrapf(types.ErrContractNotFound, "system contract state variable not found")
	}
	oracle := ethcommon.HexToAddress(system.SystemContract)
	if oracle == (ethcommon.Address{}) {
		return 0, cosmoserrors.Wrapf(types.ErrContractNotFound, "system contract invalid address")
	}
	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return 0, cosmoserrors.Wrapf(types.ErrABIGet, "SystemContractMetaData")
	}
	res, err := k.CallEVM(
		ctx,
		*abi,
		types.ModuleAddressEVM,
		oracle,
		BigIntZero,
		DefaultGasLimit,
		true,
		false,
		"setGasPrice",
		chainid,
		gasPrice,
	)
	if err != nil {
		return 0, cosmoserrors.Wrap(types.ErrContractCall, err.Error())
	}
	if res.Failed() {
		return res.GasUsed, cosmoserrors.Wrapf(types.ErrContractCall, "setGasPrice tx failed")
	}

	return res.GasUsed, nil
}

func (k Keeper) SetGasCoin(ctx sdk.Context, chainid *big.Int, address ethcommon.Address) error {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return cosmoserrors.Wrapf(types.ErrContractNotFound, "system contract state variable not found")
	}
	oracle := ethcommon.HexToAddress(system.SystemContract)
	if oracle == (ethcommon.Address{}) {
		return cosmoserrors.Wrapf(types.ErrContractNotFound, "system contract invalid address")
	}
	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrABIGet, "SystemContractMetaData")
	}
	res, err := k.CallEVM(
		ctx,
		*abi,
		types.ModuleAddressEVM,
		oracle,
		BigIntZero,
		DefaultGasLimit,
		true,
		false,
		"setGasCoinZRC20",
		chainid,
		address,
	)
	if err != nil {
		return cosmoserrors.Wrap(types.ErrContractCall, err.Error())
	}
	if res.Failed() {
		return cosmoserrors.Wrapf(types.ErrContractCall, "setGasCoinZRC20 tx failed")
	}

	return nil
}

func (k Keeper) SetGasZetaPool(ctx sdk.Context, chainid *big.Int, pool ethcommon.Address) error {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return cosmoserrors.Wrapf(types.ErrContractNotFound, "system contract state variable not found")
	}
	oracle := ethcommon.HexToAddress(system.SystemContract)
	if oracle == (ethcommon.Address{}) {
		return cosmoserrors.Wrapf(types.ErrContractNotFound, "system contract invalid address")
	}
	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrABIGet, "SystemContractMetaData")
	}
	res, err := k.CallEVM(
		ctx,
		*abi,
		types.ModuleAddressEVM,
		oracle,
		BigIntZero,
		DefaultGasLimit,
		true,
		false,
		"setGasZetaPool",
		chainid,
		pool,
	)
	if err != nil {
		return cosmoserrors.Wrap(types.ErrContractCall, err.Error())
	}
	if res.Failed() {
		return cosmoserrors.Wrapf(types.ErrContractCall, "setGasZetaPool tx failed")
	}

	return nil
}
