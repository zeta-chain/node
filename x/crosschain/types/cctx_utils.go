package types

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/btcsuite/btcutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	if m.Index != "" {
		err := ValidateZetaIndex(m.Index)
		if err != nil {
			return err
		}
	}
	err := m.InboundTxParams.Validate()
	if err != nil {
		return err
	}
	for _, outboundTxParam := range m.OutboundTxParams {
		err = outboundTxParam.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m InboundTxParams) Validate() error {
	if m.Sender == "" {
		return fmt.Errorf("sender cannot be empty")
	}
	if common.GetChainFromChainID(m.SenderChainId) == nil {
		return fmt.Errorf("invalid sender chain id %d", m.SenderChainId)
	}
	err := ValidateAddressForChain(m.Sender, m.SenderChainId)
	if err != nil {
		return err
	}

	if m.TxOrigin != "" {
		errTxOrigin := ValidateAddressForChain(m.TxOrigin, m.SenderChainId)
		if errTxOrigin != nil {
			return errTxOrigin
		}
	}
	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	err = ValidateHashForChain(m.InboundTxObservedHash, m.SenderChainId)
	if err != nil {
		return err
	}
	if m.InboundTxBallotIndex != "" {
		err = ValidateZetaIndex(m.InboundTxBallotIndex)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m OutboundTxParams) Validate() error {
	if m.Receiver == "" {
		return fmt.Errorf("receiver cannot be empty")
	}
	if common.GetChainFromChainID(m.ReceiverChainId) == nil {
		return fmt.Errorf("invalid receiver chain id %d", m.ReceiverChainId)
	}
	err := ValidateAddressForChain(m.Receiver, m.ReceiverChainId)
	if err != nil {
		return err
	}
	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	if m.OutboundTxBallotIndex != "" {
		err = ValidateZetaIndex(m.OutboundTxBallotIndex)
		if err != nil {
			return err
		}
	}
	if m.OutboundTxHash != "" {
		err = ValidateHashForChain(m.OutboundTxHash, m.ReceiverChainId)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateZetaIndex(index string) error {
	if len(index) != 66 {
		return fmt.Errorf("invalid index hash %s", index)
	}
	return nil
}
func ValidateHashForChain(hash string, chainID int64) error {
	if common.IsEthereumChain(chainID) || common.IsZetaChain(chainID) {
		_, err := hexutil.Decode(hash)
		if err != nil {
			return fmt.Errorf("hash must be a valid ethereum hash %s", hash)
		}
		return nil
	}
	if common.IsBitcoinChain(chainID) {
		r, err := regexp.Compile("^[a-fA-F0-9]{64}$")
		if err != nil {
			return fmt.Errorf("error compiling regex")
		}
		if !r.MatchString(hash) {
			return fmt.Errorf("hash must be a valid bitcoin hash %s", hash)
		}
		return nil
	}
	return fmt.Errorf("invalid chain id %d", chainID)
}

func ValidateAddressForChain(address string, chainID int64) error {
	// we do not validate the address for zeta chain as the address field can be btc or eth address
	if common.IsZetaChain(chainID) {
		return nil
	}
	if common.IsEthereumChain(chainID) {
		if !ethcommon.IsHexAddress(address) {
			return fmt.Errorf("invalid address %s , chain %d", address, chainID)
		}
		return nil
	}
	if common.IsBitcoinChain(chainID) {
		addr, err := common.DecodeBtcAddress(address, chainID)
		if err != nil {
			return fmt.Errorf("invalid address %s , chain %d: %s", address, chainID, err)
		}
		_, ok := addr.(*btcutil.AddressWitnessPubKeyHash)
		if !ok {
			return fmt.Errorf(" invalid address %s (not P2WPKH address)", address)
		}
		return nil
	}
	return fmt.Errorf("invalid chain id %d", chainID)
}

// GetGasPrice returns the gas price of the outbound tx
func (m OutboundTxParams) GetGasPrice() (uint64, error) {
	gasPrice, err := strconv.ParseUint(m.OutboundTxGasPrice, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse cctx gas price %s: %s", m.OutboundTxGasPrice, err.Error())
	}

	return gasPrice, nil
}
