package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/pkg/systemcontract.sol"

	"github.com/zeta-chain/node/pkg/coin"
)

// ProcessDeposit handles a deposit from an inbound tx with protocol version 2
// returns [txResponse, isContractCall, error]
// isContractCall is true if the message is non empty
func (k Keeper) ProcessDeposit(
	ctx sdk.Context,
	from []byte,
	senderChainID int64,
	zrc20Addr ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int,
	message []byte,
	coinType coin.CoinType,
	isCrossChainCall bool,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	if coinType == coin.CoinType_Zeta {
		return nil, false, errors.New("ZETA asset is currently unsupported for deposit with V2 protocol contracts")
	}

	context := systemcontract.ZContext{
		Origin:  []byte{},
		Sender:  ethcommon.BytesToAddress(from),
		ChainID: big.NewInt(senderChainID),
	}

	if coinType == coin.CoinType_NoAssetCall {
		// simple call
		res, err := k.CallExecute(ctx, context, zrc20Addr, amount, to, message)
		return res, true, err
	} else if isCrossChainCall {
		// call with asset
		res, err := k.CallDepositAndCallZRC20(ctx, context, zrc20Addr, amount, to, message)
		return res, true, err
	}

	// simple deposit
	res, err := k.DepositZRC20(ctx, zrc20Addr, to, amount)
	return res, false, err
}

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
	if coinType == coin.CoinType_Zeta {
		return nil, errors.New("ZETA asset is currently unsupported for revert with V2 protocol contracts")
	}

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
		} else {
			// no asset, no call, do nothing
			return nil, nil
		}
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
		} else {
			// simply deposit back to the revert address
			res, err := k.DepositZRC20(ctx, zrc20Addr, revertAddress, amount)
			return res, err
		}
	}

	return nil, fmt.Errorf("unsupported coin type for revert %s", coinType)
}

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
	if coinType == coin.CoinType_Zeta {
		return nil, errors.New("ZETA asset is currently unsupported for abort with V2 protocol contracts")
	}

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

	// if the cctx contains asset, the asset is first deposited to the abort address, separately from onAbort call
	if coinType == coin.CoinType_ERC20 || coinType == coin.CoinType_NoAssetCall {
		// simply deposit back to the revert address
		// if the deposit fails, processing the abort entirely fails
		// MsgRefundAbort can still be used to retry the operation later on
		if _, err := k.DepositZRC20(ctx, zrc20Addr, abortAddress, amount); err != nil {
			return nil, err
		}
	}

	// call onAbort
	return k.CallExecuteAbort(
		ctx,
		inboundSender,
		zrc20Addr,
		amount,
		outgoing,
		big.NewInt(chainID),
		abortAddress,
		revertMessage,
	)
}
