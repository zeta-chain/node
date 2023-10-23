package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// GetCurrentOutTxParam returns the current outbound tx params.
// There can only be one active outtx.
// OutboundTxParams[0] is the original outtx, if it reverts, then
// OutboundTxParams[1] is the new outtx.
func (m CrossChainTx) GetCurrentOutTxParam() *OutboundTxParams {
	if len(m.OutboundTxParams) == 0 {
		return &OutboundTxParams{}
	}
	return m.OutboundTxParams[len(m.OutboundTxParams)-1]
}

// IsCurrentOutTxRevert returns true if the current outbound tx is the revert tx.
func (m CrossChainTx) IsCurrentOutTxRevert() bool {
	return len(m.OutboundTxParams) == 2
}

// OriginalDestinationChainID returns the original destination of the outbound tx, reverted or not
// If there is no outbound tx, return -1
func (m CrossChainTx) OriginalDestinationChainID() int64 {
	if len(m.OutboundTxParams) == 0 {
		return -1
	}
	return m.OutboundTxParams[0].ReceiverChainId
}

// GetAllAuthzZetaclientTxTypes returns all the authz types for zetaclient
func GetAllAuthzZetaclientTxTypes() []string {
	return []string{
		sdk.MsgTypeURL(&MsgNonceVoter{}),
		sdk.MsgTypeURL(&MsgGasPriceVoter{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedInboundTx{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedOutboundTx{}),
		sdk.MsgTypeURL(&MsgSetNodeKeys{}),
		sdk.MsgTypeURL(&MsgCreateTSSVoter{}),
		sdk.MsgTypeURL(&MsgAddToOutTxTracker{}),
		sdk.MsgTypeURL(&MsgSetNodeKeys{}),
		sdk.MsgTypeURL(&observertypes.MsgAddBlameVote{}),
		sdk.MsgTypeURL(&observertypes.MsgAddBlockHeader{}),
	}
}

// GetGasPrice returns the gas price of the outbound tx
func (m OutboundTxParams) GetGasPrice() (uint64, error) {
	gasPrice, err := strconv.ParseUint(m.OutboundTxGasPrice, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse cctx gas price %s: %s", m.OutboundTxGasPrice, err.Error())
	}

	return gasPrice, nil
}
