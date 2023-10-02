package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	systemcontract "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// SetGasPrice sets gas price on the system contract in zEVM; return the gasUsed and error code
func (k Keeper) SetGasPrice(ctx sdk.Context, chainid *big.Int, gasPrice *big.Int) (uint64, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return 0, sdkerrors.Wrapf(types.ErrContractNotFound, "system contract state variable not found")
	}
	oracle := common.HexToAddress(system.SystemContract)
	if oracle == common.HexToAddress("0x0") {
		return 0, sdkerrors.Wrapf(types.ErrContractNotFound, "system contract invalid address")
	}
	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return 0, sdkerrors.Wrapf(types.ErrABIGet, "SystemContractMetaData")
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, oracle, BigIntZero, big.NewInt(50_000), true, false, "setGasPrice", chainid, gasPrice)
	if err != nil {
		return 0, sdkerrors.Wrapf(types.ErrABIGet, err.Error())
	}
	if res.Failed() {
		return res.GasUsed, sdkerrors.Wrapf(types.ErrContractCall, "setGasPrice")
	}

	return res.GasUsed, nil
}

func (k Keeper) SetGasCoin(ctx sdk.Context, chainid *big.Int, address common.Address) error {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "system contract state variable not found")
	}
	oracle := common.HexToAddress(system.SystemContract)
	if oracle == common.HexToAddress("0x0") {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "system contract invalid address")
	}
	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIGet, "SystemContractMetaData")
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, oracle, BigIntZero, nil, true, false, "setGasCoinZRC20", chainid, address)
	if err != nil || res.Failed() {
		return sdkerrors.Wrapf(types.ErrContractCall, "setGasCoinZRC20")
	}

	return nil
}

func (k Keeper) SetGasZetaPool(ctx sdk.Context, chainid *big.Int, pool common.Address) error {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "system contract state variable not found")
	}
	oracle := common.HexToAddress(system.SystemContract)
	if oracle == common.HexToAddress("0x0") {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "system contract invalid address")
	}
	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIGet, "SystemContractMetaData")
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, oracle, BigIntZero, nil, true, false, "setGasZetaPool", chainid, pool)
	if err != nil || res.Failed() {
		return sdkerrors.Wrapf(types.ErrContractCall, "setGasZetaPool")
	}

	return nil
}
