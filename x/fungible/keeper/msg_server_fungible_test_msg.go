package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
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

	ZDCAddress, err := k.DeployZetaDepositAndCall(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployZetaDepositAndCall")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("ZetaDepositAndCallContract", ZDCAddress.String()),
		),
	)

	gasPriceOracle, err := k.DeployGasPriceOracleContract(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployZetaDepositAndCall")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("gasPriceOracle", gasPriceOracle.String()),
		),
	)

	// setup uniswap v2 factory
	uniswapV2Factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployUniswapV2Factory")
	}

	// set the system contract
	system, _ := k.GetSystemContract(ctx)
	system.ZetaDepositAndCallContract = ZDCAddress.String()
	system.GasPriceOracleContract = gasPriceOracle.String()
	system.Uniswapv2FactoryAddress = uniswapV2Factory.String()
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
func (k Keeper) setupChainGasCoinAndPool(ctx sdk.Context, chain string, gasAssetName string, symbol string, decimals uint8) (ethcommon.Address, error) {
	name := fmt.Sprintf("%s-%s", gasAssetName, chain)
	transferGasLimit := big.NewInt(21_000)
	addr, err := k.DeployZRC4Contract(ctx, name, symbol, decimals, chain, common.CoinType_Gas, "", transferGasLimit)
	if err != nil {
		return ethcommon.Address{}, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(name, addr.String()),
		),
	)
	chainid := clientconfig.Chains["BSCTESTNET"].ChainID
	k.SetGasCoin(ctx, chainid, addr)
	if err != nil {
		return ethcommon.Address{}, err
	}
	return addr, nil
}
