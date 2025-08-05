package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/coin"
)

// ProcessRevert handles a revert deposit from an inbound tx with protocol version 2
func (k Keeper) ProcessRevert(
	ctx sdk.Context,
	inboundSender string,
	amount *big.Int,
	chainID int64,
	coinType coin.CoinType,
	asset string,
	revertAddress ethcommon.Address,
	callOnRevert bool,
	revertMessage []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	// get the zrc20 contract
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

	switch coinType {
	case coin.CoinType_NoAssetCall:
		return k.processNoAssetCallRevert(
			ctx,
			inboundSender,
			zrc20Addr,
			amount,
			revertAddress,
			revertMessage,
			callOnRevert,
		)

	case coin.CoinType_Zeta:
		return k.processZetaRevert(ctx, inboundSender, amount, revertAddress, revertMessage, callOnRevert)

	case coin.CoinType_ERC20, coin.CoinType_Gas:
		return k.processZRC20Revert(ctx, inboundSender, zrc20Addr, amount, revertAddress, revertMessage, callOnRevert)

	default:
		return nil, fmt.Errorf("unsupported coin type for revert %s", coinType)
	}
}

// processNoAssetCallRevert handles reverts with no asset (simple calls)
func (k Keeper) processNoAssetCallRevert(
	ctx sdk.Context,
	inboundSender string,
	zrc20Addr ethcommon.Address,
	amount *big.Int,
	revertAddress ethcommon.Address,
	revertMessage []byte,
	callOnRevert bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
	if callOnRevert {
		return k.CallExecuteRevert(ctx, inboundSender, zrc20Addr, amount, revertAddress, revertMessage)
	}
	return nil, nil
}

// processZetaRevert handles ZETA coin reverts
func (k Keeper) processZetaRevert(
	ctx sdk.Context,
	inboundSender string,
	amount *big.Int,
	revertAddress ethcommon.Address,
	revertMessage []byte,
	callOnRevert bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
	res, _, err := k.ExecuteWithMintedZeta(
		ctx,
		amount,
		func(tmpCtx sdk.Context) (*evmtypes.MsgEthereumTxResponse, bool, error) {
			if callOnRevert {
				res, err := k.CallZetaDepositAndRevert(tmpCtx, inboundSender, amount, revertAddress, revertMessage)
				return res, false, err
			}

			res, err := k.DepositZeta(tmpCtx, revertAddress, amount)
			return res, false, err
		},
	)

	return res, err
}

// processZRC20Revert handles ZRC20 token reverts [ZRC20 tokens exist for ERC20 and GAS tokens]
func (k Keeper) processZRC20Revert(
	ctx sdk.Context,
	inboundSender string,
	zrc20Addr ethcommon.Address,
	amount *big.Int,
	revertAddress ethcommon.Address,
	revertMessage []byte,
	callOnRevert bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
	if callOnRevert {
		return k.CallDepositAndRevert(ctx, inboundSender, zrc20Addr, amount, revertAddress, revertMessage)
	}

	return k.DepositZRC20(ctx, zrc20Addr, revertAddress, amount)
}
