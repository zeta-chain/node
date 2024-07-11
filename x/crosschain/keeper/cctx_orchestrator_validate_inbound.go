package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/crypto"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ValidateInbound is the only entry-point to create new CCTX (eg. when observers voting is done or new inbound event is detected).
// It creates new CCTX object and calls InitiateOutbound method.
func (k Keeper) ValidateInbound(
	ctx sdk.Context,
	msg *types.MsgVoteInbound,
	shouldPayGas bool,
) (*types.CrossChainTx, error) {
	err := k.CheckMigration(ctx, msg)
	if err != nil {
		return nil, err
	}

	tss, tssFound := k.zetaObserverKeeper.GetTSS(ctx)
	if !tssFound {
		return nil, types.ErrCannotFindTSSKeys
	}

	// Do not process if inbound is disabled
	if !k.zetaObserverKeeper.IsInboundEnabled(ctx) {
		return nil, observertypes.ErrInboundDisabled
	}

	// create a new CCTX from the inbound message. The status of the new CCTX is set to PendingInbound.
	cctx, err := types.NewCCTX(ctx, *msg, tss.TssPubkey)
	if err != nil {
		return nil, err
	}

	// Initiate outbound, the process function manages the state commit and cctx status change.
	// If the process fails, the changes to the evm state are rolled back.
	_, err = k.InitiateOutbound(ctx, InitiateOutboundConfig{
		CCTX:         &cctx,
		ShouldPayGas: shouldPayGas,
	})
	if err != nil {
		return nil, err
	}

	inCctxIndex, ok := ctx.Value(InCCTXIndexKey).(string)
	if ok {
		cctx.InboundParams.ObservedHash = inCctxIndex
	}
	k.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, cctx)

	return &cctx, nil
}

// CheckMigration checks if the sender is a TSS address and returns an error if it is.
// If the sender is an older TSS address, this means that it is a migration transfer, and we do not need to treat this as a deposit.
func (k Keeper) CheckMigration(ctx sdk.Context, msg *types.MsgVoteInbound) error {
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)
	// Ignore cctx originating from zeta chain/zevm for this check as there is no TSS in zeta chain
	if chains.IsZetaChain(msg.SenderChainId, additionalChains) {
		return nil
	}

	historicalTssList := k.zetaObserverKeeper.GetAllTSS(ctx)
	chain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, msg.SenderChainId)
	if !found {
		return observertypes.ErrSupportedChains.Wrapf("chain not found for chainID %d", msg.SenderChainId)
	}

	for _, tss := range historicalTssList {
		if chains.IsEVMChain(chain.ChainId, additionalChains) {
			ethTssAddress, err := crypto.GetTssAddrEVM(tss.TssPubkey)
			if err != nil {
				return errors.Wrap(types.ErrInvalidAddress, err.Error())
			}
			if ethTssAddress.Hex() == msg.Sender {
				return types.ErrMigrationFromOldTss
			}
		} else if chains.IsBitcoinChain(chain.ChainId, additionalChains) {
			bitcoinParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
			if err != nil {
				return err
			}
			btcTssAddress, err := crypto.GetTssAddrBTC(tss.TssPubkey, bitcoinParams)
			if err != nil {
				return errors.Wrap(types.ErrInvalidAddress, err.Error())
			}
			if btcTssAddress == msg.Sender {
				return types.ErrMigrationFromOldTss
			}
		}
	}
	return nil
}
