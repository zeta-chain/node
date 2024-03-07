package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
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
		sdk.MsgTypeURL(&MsgGasPriceVoter{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedInboundTx{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedOutboundTx{}),
		sdk.MsgTypeURL(&MsgCreateTSSVoter{}),
		sdk.MsgTypeURL(&MsgAddToOutTxTracker{}),
		sdk.MsgTypeURL(&observertypes.MsgAddBlameVote{}),
		sdk.MsgTypeURL(&observertypes.MsgAddBlockHeader{}),
	}
}

func (m CrossChainTx) Validate() error {
	if m.InboundTxParams == nil {
		return fmt.Errorf("inbound tx params cannot be nil")
	}
	if m.OutboundTxParams == nil {
		return fmt.Errorf("outbound tx params cannot be nil")
	}
	if m.CctxStatus == nil {
		return fmt.Errorf("cctx status cannot be nil")
	}
	if len(m.OutboundTxParams) > 2 {
		return fmt.Errorf("outbound tx params cannot be more than 2")
	}
	if len(m.Index) != 66 {
		return ErrInvalidCCTXIndex
	}

	return nil
}

func (m InboundTxParams) Validate() error {
	if m.Sender == "" {
		return fmt.Errorf("sender cannot be empty")
	}
	if m.InboundTxObservedHash == "" {
		return fmt.Errorf("inbound tx observed hash cannot be empty")
	}
	if len(m.InboundTxBallotIndex) != 66 {
		return fmt.Errorf("inbound tx ballot index must be 66 characters")
	}
	if common.IsEthereumChain(m.SenderChainId) {
		if !ethcommon.IsHexAddress(m.Sender) {
			return fmt.Errorf("sender a valid ethereum address")
		}
	}
	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	if common.IsBitcoinChain(m.SenderChainId) {
		//if _, err := common.BitcoinAddressToPubKeyHash(m.Sender); err != nil {
		//	return fmt.Errorf("sender must be a valid bitcoin address")
		//}
	}
	return nil
}

func (m OutboundTxParams) Validate() error {
	if m.Receiver == "" {
		return fmt.Errorf("receiver cannot be empty")
	}
	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	if m.OutboundTxGasPrice == "" {
		return fmt.Errorf("outbound tx gas price cannot be empty")
	}
	if m.GasLimit == 0 {
		return fmt.Errorf("gas limit cannot be 0")
	}
	if m.ReceiverChainId == 0 {
		return fmt.Errorf("receiver chain id cannot be 0")
	}
	if common.IsEthereumChain(m.ReceiverChainId) {
		if !ethcommon.IsHexAddress(m.Receiver) {
			return fmt.Errorf("receiver must be a valid ethereum address")
		}
	}
	if common.IsBitcoinChain(m.ReceiverChainId) {
		//if _, err := common.BitcoinAddressToPubKeyHash(m.Receiver); err != nil {
		//	return fmt.Errorf("receiver must be a valid bitcoin address")
		//}
	}
	return nil
}

// GetGasPrice returns the gas price of the outbound tx
func (m OutboundTxParams) GetGasPrice() (uint64, error) {
	gasPrice, err := strconv.ParseUint(m.OutboundTxGasPrice, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse cctx gas price %s: %s", m.OutboundTxGasPrice, err.Error())
	}

	return gasPrice, nil
}
