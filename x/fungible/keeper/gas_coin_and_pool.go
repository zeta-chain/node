package keeper

import (
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	systemcontract "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	uniswapv2router02 "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// SetupChainGasCoinAndPool setup gas ZRC20, and ZETA/gas pool for a chain
// add 0.1gas/0.1wzeta to the pool
// FIXME: add cointype and use proper gas limit based on cointype/chain
func (k Keeper) SetupChainGasCoinAndPool(
	ctx sdk.Context,
	chainID int64,
	gasAssetName string,
	symbol string,
	decimals uint8,
	gasLimit *big.Int,
) (ethcommon.Address, error) {
	// additional on-chain static chain information
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)

	chain, found := chains.GetChainFromChainID(chainID, additionalChains)
	if !found {
		return ethcommon.Address{}, zetaObserverTypes.ErrSupportedChains
	}
	name := fmt.Sprintf("%s-%s", gasAssetName, chain.ChainName)

	transferGasLimit := gasLimit

	// Check if gas coin already exists
	_, found = k.GetGasCoinForForeignCoin(ctx, chainID)
	if found {
		return ethcommon.Address{}, types.ErrForeignCoinAlreadyExist
	}

	// default values
	if transferGasLimit == nil {
		transferGasLimit = big.NewInt(21_000)
		if chains.IsBitcoinChain(chain.ChainId, additionalChains) {
			transferGasLimit = big.NewInt(100) // 100B for a typical tx
		}
	}

	zrc20Addr, err := k.DeployZRC20Contract(
		ctx,
		name,
		symbol,
		decimals,
		chain.ChainId,
		coin.CoinType_Gas,
		"",
		transferGasLimit,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to DeployZRC20Contract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(name, zrc20Addr.String()),
		),
	)
	err = k.SetGasCoin(ctx, big.NewInt(chain.ChainId), zrc20Addr)
	if err != nil {
		return ethcommon.Address{}, err
	}
	amount := big.NewInt(10)
	// #nosec G115 always in range
	amount.Exp(amount, big.NewInt(int64(decimals-1)), nil)
	amountAZeta := big.NewInt(1e17)

	_, err = k.DepositZRC20(ctx, zrc20Addr, types.ModuleAddressEVM, amount)
	if err != nil {
		return ethcommon.Address{}, err
	}
	err = k.bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin("azeta", sdk.NewIntFromBigInt(amountAZeta))),
	)
	if err != nil {
		return ethcommon.Address{}, err
	}
	systemContractAddress, err := k.GetSystemContractAddress(ctx)
	if err != nil || systemContractAddress == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractNotFound,
			"system contract address invalid: %s",
			systemContractAddress,
		)
	}
	systemABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}
	_, err = k.CallEVM(
		ctx,
		*systemABI,
		types.ModuleAddressEVM,
		systemContractAddress,
		BigIntZero,
		nil,
		true,
		false,
		"setGasZetaPool",
		big.NewInt(chain.ChainId),
		zrc20Addr,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			err,
			"failed to CallEVM method setGasZetaPool(%d, %s)",
			chain.ChainId,
			zrc20Addr.String(),
		)
	}

	// setup uniswap v2 pools gas/zeta
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get uniswap router abi")
	}
	ZRC20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to GetAbi zrc20")
	}
	_, err = k.CallEVM(
		ctx,
		*ZRC20ABI,
		types.ModuleAddressEVM,
		zrc20Addr,
		BigIntZero,
		nil,
		true,
		false,
		"approve",
		routerAddress,
		amount,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			err,
			"failed to CallEVM method approve(%s, %d)",
			routerAddress.String(),
			amount,
		)
	}

	//function addLiquidityETH(
	//	address token,
	//	uint amountTokenDesired,
	//	uint amountTokenMin,
	//	uint amountETHMin,
	//	address to,
	//	uint deadline
	//) external payable returns (uint amountToken, uint amountETH, uint liquidity);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		amountAZeta,
		big.NewInt(5_000_000),
		true,
		false,
		"addLiquidityETH",
		zrc20Addr,
		amount,
		BigIntZero,
		BigIntZero,
		types.ModuleAddressEVM,
		amountAZeta,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			err,
			"failed to CallEVM method addLiquidityETH(%s, %s)",
			zrc20Addr.String(),
			amountAZeta.String(),
		)
	}
	AmountToken := new(*big.Int)
	AmountETH := new(*big.Int)
	Liquidity := new(*big.Int)
	err = routerABI.UnpackIntoInterface(&[]interface{}{AmountToken, AmountETH, Liquidity}, "addLiquidityETH", res.Ret)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface addLiquidityETH")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("function", "addLiquidityETH"),
			sdk.NewAttribute("amountToken", (*AmountToken).String()),
			sdk.NewAttribute("amountETH", (*AmountETH).String()),
			sdk.NewAttribute("liquidity", (*Liquidity).String()),
		),
	)
	return zrc20Addr, nil
}
