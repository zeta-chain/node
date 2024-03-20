package types

import (
	"fmt"
	"regexp"

	"github.com/btcsuite/btcutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

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
