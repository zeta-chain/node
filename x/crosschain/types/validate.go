package types

import (
	"fmt"
	"regexp"

	"cosmossdk.io/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

// ValidateZetaIndex validates the zeta index
func ValidateZetaIndex(index string) error {
	if len(index) != ZetaIndexLength {
		return errors.Wrap(ErrInvalidIndexValue, fmt.Sprintf("invalid index length %d", len(index)))
	}
	return nil
}

// ValidateHashForChain validates the hash for the chain
func ValidateHashForChain(hash string, chainID int64) error {
	if chains.IsEthereumChain(chainID) || chains.IsZetaChain(chainID) {
		_, err := hexutil.Decode(hash)
		if err != nil {
			return fmt.Errorf("hash must be a valid ethereum hash %s", hash)
		}
		return nil
	}
	if chains.IsBitcoinChain(chainID) {
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

// ValidateAddressForChain validates the address for the chain
func ValidateAddressForChain(address string, chainID int64) error {
	// we do not validate the address for zeta chain as the address field can be btc or eth address
	if chains.IsZetaChain(chainID) {
		return nil
	}
	if chains.IsEthereumChain(chainID) {
		if !ethcommon.IsHexAddress(address) {
			return fmt.Errorf("invalid address %s , chain %d", address, chainID)
		}
		return nil
	}
	if chains.IsBitcoinChain(chainID) {
		addr, err := chains.DecodeBtcAddress(address, chainID)
		if err != nil {
			return fmt.Errorf("invalid address %s , chain %d: %s", address, chainID, err)
		}
		if !chains.IsBtcAddressSupported(addr) {
			return fmt.Errorf("unsupported address %s", address)
		}
		return nil
	}
	return fmt.Errorf("invalid chain id %d", chainID)
}
