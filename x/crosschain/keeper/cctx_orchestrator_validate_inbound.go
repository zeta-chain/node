package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// ValidateInbound is the only entry-point to create new CCTX (eg. when observers voting is done or new inbound event is detected).
// It creates new CCTX object and calls InitiateOutbound method.
func (k Keeper) ValidateInbound(
	ctx sdk.Context,
	msg *types.MsgVoteInbound,
	shouldPayGas bool,
) (*types.CrossChainTx, error) {
	tss, tssFound := k.zetaObserverKeeper.GetTSS(ctx)
	if !tssFound {
		return nil, types.ErrCannotFindTSSKeys
	}

	err := k.CheckIfTSSMigrationTransfer(ctx, msg)
	if err != nil {
		return nil, err
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
	k.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, cctx, tss.TssPubkey)

	return &cctx, nil
}

// CheckIfTSSMigrationTransfer checks if the sender is a TSS address and returns an error if it is.
// If the sender is an older TSS address, this means that it is a migration transfer, and we do not need to treat this as a deposit and process the CCTX
func (k Keeper) CheckIfTSSMigrationTransfer(ctx sdk.Context, msg *types.MsgVoteInbound) error {
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)

	historicalTSSList := k.zetaObserverKeeper.GetAllTSS(ctx)
	chain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, msg.SenderChainId)
	if !found {
		return observertypes.ErrSupportedChains.Wrapf("chain not found for chainID %d", msg.SenderChainId)
	}

	// the check is only necessary if the inbound is validated from observers from a connected chain
	if chain.CctxGateway != chains.CCTXGateway_observers {
		return nil
	}

	switch {
	case chains.IsEVMChain(chain.ChainId, additionalChains):
		for _, tss := range historicalTSSList {
			ethTssAddress, err := crypto.GetTssAddrEVM(tss.TssPubkey)
			if err != nil {
				continue
			}
			if ethTssAddress.Hex() == msg.Sender {
				return types.ErrMigrationFromOldTss
			}
		}
	case chains.IsBitcoinChain(chain.ChainId, additionalChains):
		bitcoinParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
		if err != nil {
			return err
		}
		for _, tss := range historicalTSSList {
			btcTssAddress, err := crypto.GetTssAddrBTC(tss.TssPubkey, bitcoinParams)
			if err != nil {
				continue
			}
			if btcTssAddress == msg.Sender {
				return types.ErrMigrationFromOldTss
			}
		}
	}

	return nil
}
