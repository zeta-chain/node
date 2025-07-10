package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"

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
		if callOnRevert {
			// no asset, call simple revert
			res, err := k.CallExecuteRevert(ctx, inboundSender, zrc20Addr, amount, revertAddress, revertMessage)
			return res, err
		}
		return nil, nil
	case coin.CoinType_ERC20, coin.CoinType_Gas:
		if callOnRevert {
			// revert with a ZRC20 asset
			res, err := k.CallDepositAndRevert(
				ctx,
				inboundSender,
				zrc20Addr,
				amount,
				revertAddress,
				revertMessage,
			)
			return res, err
		}
		// simply deposit back to the revert address
		res, err := k.DepositZRC20(ctx, zrc20Addr, revertAddress, amount)
		return res, err
	case coin.CoinType_Zeta:
		// if the coin type is Zeta, handle this as a deposit ZETA to zEVM.
		if err := k.MintZetaToFungibleModule(ctx, amount); err != nil {
			return nil, err
		}
		if callOnRevert {
			res, err := k.CallZetaDepositAndRevert(
				ctx,
				inboundSender,
				amount,
				revertAddress,
				revertMessage,
			)
			return res, err
		}
		// deposit ZETA to the revert address
		res, err := k.DepositZeta(ctx, revertAddress, amount)
		return res, err
	}
	return nil, fmt.Errorf("unsupported coin type for revert %s", coinType)
}
