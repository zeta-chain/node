package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) RefundAbortedAmountOnZetaChain(ctx sdk.Context, cctx types.CrossChainTx, refundAddress ethcommon.Address) error {
	coinType := cctx.InboundParams.CoinType
	switch coinType {
	case coin.CoinType_Gas:
		return k.RefundAmountOnZetaChainGas(ctx, cctx, refundAddress)
	case coin.CoinType_Zeta:
		return k.RefundAmountOnZetaChainZeta(ctx, cctx, refundAddress)
	case coin.CoinType_ERC20:
		return k.RefundAmountOnZetaChainERC20(ctx, cctx, refundAddress)
	default:
		return errors.New("unsupported coin type for refund on ZetaChain")
	}
}

// RefundAmountOnZetaChainGas refunds the amount of the cctx on ZetaChain in case of aborted cctx with cointype gas
func (k Keeper) RefundAmountOnZetaChainGas(ctx sdk.Context, cctx types.CrossChainTx, refundAddress ethcommon.Address) error {
	// refund in gas token to refund address
	// Refund the the amount was previously
	refundAmount := GetAbortedAmount(cctx)
	if refundAmount.IsNil() || refundAmount.IsZero() {
		return errors.New("no amount to refund")
	}
	chainID := cctx.InboundParams.SenderChainId
	// get the zrc20 contract address
	fcSenderChain, found := k.fungibleKeeper.GetGasCoinForForeignCoin(ctx, chainID)
	if !found {
		return types.ErrForeignCoinNotFound
	}
	zrc20 := ethcommon.HexToAddress(fcSenderChain.Zrc20ContractAddress)
	if zrc20 == (ethcommon.Address{}) {
		return errorsmod.Wrapf(types.ErrForeignCoinNotFound, "zrc20 contract address not found for chain %d", chainID)
	}
	// deposit the amount to the tx origin instead of receiver as this is a refund
	if _, err := k.fungibleKeeper.DepositZRC20(ctx, zrc20, refundAddress, refundAmount.BigInt()); err != nil {
		return errors.New("failed to refund zeta on ZetaChain" + err.Error())
	}
	return nil
}

// RefundAmountOnZetaChainGas refunds the amount of the cctx on ZetaChain in case of aborted cctx with cointype zeta
func (k Keeper) RefundAmountOnZetaChainZeta(ctx sdk.Context, cctx types.CrossChainTx, refundAddress ethcommon.Address) error {
	// if coin type is Zeta, handle this as a deposit ZETA to zEVM.
	refundAmount := GetAbortedAmount(cctx)
	chainID := cctx.InboundParams.SenderChainId
	// check if chain is an EVM chain
	if !chains.IsEVMChain(chainID) {
		return errors.New("only EVM chains are supported for refund when coin type is Zeta")
	}
	if cctx.InboundParams.Amount.IsNil() || cctx.InboundParams.Amount.IsZero() {
		return errors.New("no amount to refund")
	}
	// deposit the amount to refund address
	if err := k.fungibleKeeper.DepositCoinZeta(ctx, refundAddress, refundAmount.BigInt()); err != nil {
		return errors.New("failed to refund zeta on ZetaChain" + err.Error())
	}
	return nil
}

// RefundAmountOnZetaChainERC20 refunds the amount of the cctx on ZetaChain in case of aborted cctx
// NOTE: GetCurrentOutboundParam should contain the last up to date cctx amount
// Refund address should already be validated before calling this function
func (k Keeper) RefundAmountOnZetaChainERC20(ctx sdk.Context, cctx types.CrossChainTx, refundAddress ethcommon.Address) error {
	refundAmount := GetAbortedAmount(cctx)
	// preliminary checks
	if cctx.InboundParams.CoinType != coin.CoinType_ERC20 {
		return errors.New("unsupported coin type for refund on ZetaChain")
	}
	if !chains.IsEVMChain(cctx.InboundParams.SenderChainId) {
		return errors.New("only EVM chains are supported for refund on ZetaChain")
	}

	if refundAmount.IsNil() || refundAmount.IsZero() {
		return errors.New("no amount to refund")
	}

	// get address of the zrc20
	fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundParams.Asset, cctx.InboundParams.SenderChainId)
	if !found {
		return fmt.Errorf("asset %s zrc not found", cctx.InboundParams.Asset)
	}
	zrc20 := ethcommon.HexToAddress(fc.Zrc20ContractAddress)
	if zrc20 == (ethcommon.Address{}) {
		return fmt.Errorf("asset %s invalid zrc address", cctx.InboundParams.Asset)
	}

	// deposit the amount to the sender
	if _, err := k.fungibleKeeper.DepositZRC20(ctx, zrc20, refundAddress, refundAmount.BigInt()); err != nil {
		return errors.New("failed to deposit zrc20 on ZetaChain" + err.Error())
	}

	return nil
}
