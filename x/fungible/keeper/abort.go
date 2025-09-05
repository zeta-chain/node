package keeper

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	errorspkg "github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// ProcessAbort handles an abort deposit from an inbound tx with protocol version 2
func (k Keeper) ProcessAbort(
	ctx sdk.Context,
	inboundSender string,
	amount *big.Int,
	outgoing bool,
	chainID int64,
	coinType coin.CoinType,
	asset string,
	abortAddress ethcommon.Address,
	revertMessage []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	zrc20Addr, _, err := k.getAndCheckZRC20(
		ctx,
		amount,
		chainID,
		coinType,
		asset,
	)
	if err != nil {
		return nil, err
	}

	// if the cctx contains asset, the asset is first deposited to the abort address, separately from onAbort call
	switch coinType {
	case coin.CoinType_ERC20, coin.CoinType_Gas:
		// simply deposit back to the abort address
		// if the deposit fails, processing the abort entirely fails
		// MsgRefundAbort can still be used to retry the operation later on
		if _, err := k.DepositZRC20(ctx, zrc20Addr, abortAddress, amount); err != nil {
			return nil, err
		}
	case coin.CoinType_Zeta:
		// Deposit native zeta to the abort address
		// If the deposit fails do not mint Zeta
		_, _, err := k.ExecuteWithMintedZeta(
			ctx,
			amount,
			func(tmpCtx sdk.Context) (*evmtypes.MsgEthereumTxResponse, bool, error) {
				res, err := k.DepositZeta(tmpCtx, abortAddress, amount)
				return res, false, err
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// call onAbort
	txRes, err := k.CallExecuteAbort(
		ctx,
		inboundSender,
		zrc20Addr,
		amount,
		outgoing,
		big.NewInt(chainID),
		abortAddress,
		revertMessage,
	)
	if err != nil {
		return txRes, errors.Join(err, types.ErrOnAbortFailed)
	}
	return txRes, nil
}

// getAndCheckZRC20 returns the ZRC20 contract address and the foreign coin information
// It handles the logic based on CoinType
// - For Zeta coin type,it returns an empty address and no foreign coin.Zeta is the native token of the chain.
// - For NoAssetCall and Gas coin types, it retrieves the gas coin for the foreign coin and checks if it is paused or has a liquidity cap.
// - For other coin types(ERC20), it retrieves the foreign coin from the asset and checks if it is paused or has a liquidity cap.
func (k Keeper) getAndCheckZRC20(
	ctx sdk.Context,
	amount *big.Int,
	chainID int64,
	coinType coin.CoinType,
	asset string,
) (ethcommon.Address, types.ForeignCoins, error) {
	var zrc20Contract ethcommon.Address
	var foreignCoin types.ForeignCoins
	var found bool

	// get foreign coin
	// retrieve the gas token of the chain for no asset call
	// this simplify the current workflow and allow to pause calls by pausing the gas token
	// TODO: refactor this logic and create specific workflow for no asset call
	// https://github.com/zeta-chain/node/issues/2627
	switch coinType {
	case coin.CoinType_Zeta:
		return ethcommon.Address{}, types.ForeignCoins{}, nil
	case coin.CoinType_NoAssetCall, coin.CoinType_Gas:
		foreignCoin, found = k.GetGasCoinForForeignCoin(ctx, chainID)
		if !found {
			return ethcommon.Address{}, types.ForeignCoins{}, crosschaintypes.ErrGasCoinNotFound
		}
	default:
		foreignCoin, found = k.GetForeignCoinFromAsset(ctx, asset, chainID)
		if !found {
			return ethcommon.Address{}, types.ForeignCoins{}, errorspkg.Wrapf(
				crosschaintypes.ErrForeignCoinNotFound,
				"asset: %s, chainID %d", asset, chainID,
			)
		}
	}

	zrc20Contract = ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress)

	// check if foreign coin is paused
	if foreignCoin.Paused {
		return ethcommon.Address{}, types.ForeignCoins{}, types.ErrPausedZRC20
	}

	// check foreign coins cap if it has a cap
	if !foreignCoin.LiquidityCap.IsNil() && !foreignCoin.LiquidityCap.IsZero() {
		liquidityCap := foreignCoin.LiquidityCap.BigInt()
		totalSupply, err := k.TotalSupplyZRC4(ctx, zrc20Contract)
		if err != nil {
			return ethcommon.Address{}, types.ForeignCoins{}, err
		}
		newSupply := new(big.Int).Add(totalSupply, amount)
		if newSupply.Cmp(liquidityCap) > 0 {
			return ethcommon.Address{}, types.ForeignCoins{}, types.ErrForeignCoinCapReached
		}
	}

	return zrc20Contract, foreignCoin, nil
}
