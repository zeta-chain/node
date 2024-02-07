package keeper

import (
	"errors"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) RefundAbortedAmountOnZetaChainForEvmChain(ctx sdk.Context, cctx types.CrossChainTx) error {
	coinType := cctx.InboundTxParams.CoinType
	switch coinType {
	case common.CoinType_Gas:
		return k.RefundAmountOnZetaChainGas(ctx, cctx)
	case common.CoinType_Zeta:
		return k.RefundAmountOnZetaChainZeta(ctx, cctx)
	case common.CoinType_ERC20:
		return k.RefundAmountOnZetaChainERC20(ctx, cctx)
	default:
		return errors.New("unsupported coin type for refund on ZetaChain")
	}
}

func (k Keeper) RefundAbortedAmountOnZetaChainForBitcoinChain(ctx sdk.Context, cctx types.CrossChainTx, evmAddressForBtcRefund string) error {
	refundTo := ethcommon.HexToAddress(evmAddressForBtcRefund)
	if refundTo == (ethcommon.Address{}) {
		return errors.New("invalid address for refund")
	}
	// Set TxOrigin to the supplied address so that the refund is made to the evm address
	cctx.InboundTxParams.TxOrigin = refundTo.String()
	return k.RefundAmountOnZetaChainGas(ctx, cctx)
}

// RefundAmountOnZetaChainGas refunds the amount of the cctx on ZetaChain in case of aborted cctx with cointype gas
func (k Keeper) RefundAmountOnZetaChainGas(ctx sdk.Context, cctx types.CrossChainTx) error {
	// refund in gas token of a sender chain to the tx origin
	chainID := cctx.InboundTxParams.SenderChainId
	amountOfGasTokenLocked := cctx.InboundTxParams.Amount
	refundTo := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
	if refundTo == (ethcommon.Address{}) {
		return errors.New("invalid address for refund")
	}
	if chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, chainID); chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}
	fcSenderChain, found := k.fungibleKeeper.GetGasCoinForForeignCoin(ctx, chainID)
	if !found {
		return types.ErrForeignCoinNotFound
	}
	zrc20 := ethcommon.HexToAddress(fcSenderChain.Zrc20ContractAddress)
	if zrc20 == (ethcommon.Address{}) {
		return cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "zrc20 contract address not found for chain %d", chainID)
	}
	// deposit the amount to the tx origin instead of receiver as this is a refund
	if _, err := k.fungibleKeeper.DepositZRC20(ctx, zrc20, refundTo, amountOfGasTokenLocked.BigInt()); err != nil {
		return errors.New("failed to refund zeta on ZetaChain" + err.Error())
	}
	return nil
}

// RefundAmountOnZetaChainGas refunds the amount of the cctx on ZetaChain in case of aborted cctx with cointype zeta
func (k Keeper) RefundAmountOnZetaChainZeta(ctx sdk.Context, cctx types.CrossChainTx) error {
	// if coin type is Zeta, handle this as a deposit ZETA to zEVM.
	// deposit the amount to the tx orgin instead of receiver as this is a refund
	to := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
	if to == (ethcommon.Address{}) {
		return errors.New("invalid receiver address")
	}
	if cctx.InboundTxParams.Amount.IsNil() || cctx.InboundTxParams.Amount.IsZero() {
		return errors.New("no amount to refund")
	}
	if err := k.fungibleKeeper.DepositCoinZeta(ctx, to, cctx.InboundTxParams.Amount.BigInt()); err != nil {
		return errors.New("failed to refund zeta on ZetaChain" + err.Error())
	}
	return nil
}

// RefundAmountOnZetaChainERC20 refunds the amount of the cctx on ZetaChain in case of aborted cctx
// NOTE: GetCurrentOutTxParam should contain the last up to date cctx amount
func (k Keeper) RefundAmountOnZetaChainERC20(ctx sdk.Context, cctx types.CrossChainTx) error {
	inputAmount := cctx.InboundTxParams.Amount
	// preliminary checks
	if cctx.InboundTxParams.CoinType != common.CoinType_ERC20 {
		return errors.New("unsupported coin type for refund on ZetaChain")
	}
	if !common.IsEVMChain(cctx.InboundTxParams.SenderChainId) {
		return errors.New("only EVM chains are supported for refund on ZetaChain")
	}
	sender := ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
	if sender == (ethcommon.Address{}) {
		return errors.New("invalid sender address")
	}
	if inputAmount.IsNil() || inputAmount.IsZero() {
		return errors.New("no amount to refund")
	}

	// get address of the zrc20
	fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundTxParams.Asset, cctx.InboundTxParams.SenderChainId)
	if !found {
		return fmt.Errorf("asset %s zrc not found", cctx.InboundTxParams.Asset)
	}
	zrc20 := ethcommon.HexToAddress(fc.Zrc20ContractAddress)
	if zrc20 == (ethcommon.Address{}) {
		return fmt.Errorf("asset %s invalid zrc address", cctx.InboundTxParams.Asset)
	}

	// deposit the amount to the sender
	if _, err := k.fungibleKeeper.DepositZRC20(ctx, zrc20, sender, inputAmount.BigInt()); err != nil {
		return errors.New("failed to deposit zrc20 on ZetaChain" + err.Error())
	}

	return nil
}
