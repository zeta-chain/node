package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
)

func (k Keeper) SetGasPrice(ctx sdk.Context, chainid *big.Int, gasPrice *big.Int) error {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "system contract not found")
	}
	oracle := common.HexToAddress(system.GasPriceOracleContract)
	if oracle == common.HexToAddress("0x0") {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "gas price oracle contract wrong address")
	}
	abi, err := contracts.GasPriceOracleMetaData.GetAbi()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIGet, "GasPriceOracleMetaData")
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, oracle, true, "setGasPrice", chainid, gasPrice)
	if err != nil || res.Failed() {
		return sdkerrors.Wrapf(types.ErrContractCall, "setGasPrice")
	}

	return nil
}

func (k Keeper) SetGasCoin(ctx sdk.Context, chainid *big.Int, address common.Address) error {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "system contract not found")
	}
	oracle := common.HexToAddress(system.GasPriceOracleContract)
	if oracle == common.HexToAddress("0x0") {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "gas price oracle contract wrong address")
	}
	abi, err := contracts.GasPriceOracleMetaData.GetAbi()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIGet, "GasPriceOracleMetaData")
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, oracle, true, "setGasCoinERC4", chainid, address)
	if err != nil || res.Failed() {
		return sdkerrors.Wrapf(types.ErrContractCall, "SetGasCoin")
	}

	return nil
}

func (k Keeper) SetGasZetaPool(ctx sdk.Context, chainid *big.Int, pool common.Address) error {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "system contract not found")
	}
	oracle := common.HexToAddress(system.GasPriceOracleContract)
	if oracle == common.HexToAddress("0x0") {
		return sdkerrors.Wrapf(types.ErrContractNotFound, "gas price oracle contract wrong address")
	}
	abi, err := contracts.GasPriceOracleMetaData.GetAbi()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIGet, "GasPriceOracleMetaData")
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, oracle, true, "setGasZetaPool", chainid, pool)
	if err != nil || res.Failed() {
		return sdkerrors.Wrapf(types.ErrContractCall, "setGasZetaPool")
	}

	return nil
}
