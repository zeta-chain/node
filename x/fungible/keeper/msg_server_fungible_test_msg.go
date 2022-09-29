package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	clientconfig "github.com/zeta-chain/zetacore/zetaclient/config"
	"math/big"
)

func (k msgServer) FungibleTestMsg(goCtx context.Context, msg *types.MsgFungibleTestMsg) (*types.MsgFungibleTestMsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("action", "FungibleTestMsg"),
			sdk.NewAttribute("fungibleModuleAddress", types.ModuleAddressEVM.String()),
		),
	)

	// setup uniswap v2 factory
	uniswapV2Factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployUniswapV2Factory")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("UniswapV2Factory", uniswapV2Factory.String()),
		),
	)

	// setup WZETA contract
	wzeta, err := k.DeployWZETA(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployWZetaContract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployWZetaContract", wzeta.String()),
		),
	)

	router, err := k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployUniswapV2Router02")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployUniswapV2Router02", router.String()),
		),
	)

	SystemContractAddress, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to SystemContractAddress")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("SystemContractAddress", SystemContractAddress.String()),
		),
	)

	// set the system contract
	system, _ := k.GetSystemContract(ctx)
	system.SystemContract = SystemContractAddress.String()
	k.SetSystemContract(ctx, system)

	_, err = k.setupChainGasCoinAndPool(ctx, "GOERLI", "ETH", "gETH", 18)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
	}
	_, err = k.setupChainGasCoinAndPool(ctx, "BSCTESTNET", "BNB", "tBNB", 18)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
	}
	_, err = k.setupChainGasCoinAndPool(ctx, "MUMBAI", "MATIC", "tMATIC", 18)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
	}

	return &types.MsgFungibleTestMsgResponse{}, nil
}

// setup gas ERC-4, and ZETA/gas pool for a chain
//
func (k Keeper) setupChainGasCoinAndPool(ctx sdk.Context, chain string, gasAssetName string, symbol string, decimals uint8) (ethcommon.Address, error) {
	name := fmt.Sprintf("%s-%s", gasAssetName, chain)
	transferGasLimit := big.NewInt(21_000)
	zrc4Addr, err := k.DeployZRC4Contract(ctx, name, symbol, decimals, chain, common.CoinType_Gas, "", transferGasLimit)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to DeployZRC4Contract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(name, zrc4Addr.String()),
		),
	)
	chainid := clientconfig.Chains[chain].ChainID
	err = k.SetGasCoin(ctx, chainid, zrc4Addr)
	if err != nil {
		return ethcommon.Address{}, err
	}

	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, fmt.Errorf("system contract not found")
	}
	systemABI, err := contracts.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to SystemContractMetaData.GetAbi")
	}
	systemContractAddress := ethcommon.HexToAddress(system.SystemContract)
	if systemContractAddress == (ethcommon.Address{}) {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrContractNotFound, "system contract address invalid: %s", systemContractAddress)
	}

	_, err = k.CallEVM(ctx, *systemABI, types.ModuleAddressEVM, systemContractAddress, true, "setGasZetaPool", chainid, zrc4Addr)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to CallEVM method setGasZetaPool(%d, %s)", chainid, zrc4Addr.String())
	}
	//res, err := k.CallEVM(ctx, *systemABI, types.ModuleAddressEVM, systemContractAddress, false, "wzetaContractAddress", chainid, zrc4Addr)
	//if err != nil {
	//	return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to CallEVM method wzetaContractAddress")
	//}
	//var addressResponse types.SystemAddressResponse
	//if err := systemABI.UnpackIntoInterface(&addressResponse, "wzetaContractAddress", res.Ret); err != nil {
	//	return ethcommon.Address{}, sdkerrors.Wrapf(
	//		types.ErrABIUnpack, "failed to unpack wzetaContractAddress: %s", err.Error(),
	//	)
	//}
	//wzeta := addressResponse.Value
	//
	//res, err = k.CallEVM(ctx, *systemABI, types.ModuleAddressEVM, systemContractAddress, false, "uniswapv2FactoryAddress", chainid, zrc4Addr)
	//if err != nil {
	//	return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to CallEVM method uniswapv2FactoryAddress")
	//}
	//if err := systemABI.UnpackIntoInterface(&addressResponse, "uniswapv2FactoryAddress", res.Ret); err != nil {
	//	return ethcommon.Address{}, sdkerrors.Wrapf(
	//		types.ErrABIUnpack, "failed to unpack uniswapv2FactoryAddress: %s", err.Error(),
	//	)
	//}
	//factory := addressResponse.Value

	return zrc4Addr, nil
}
