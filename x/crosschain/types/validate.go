package types

import (
	"fmt"
	"regexp"

	"cosmossdk.io/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/zeta-chain/node/pkg/chains"
)

// ValidateCCTXIndex validates the CCTX index
func ValidateCCTXIndex(index string) error {
	if len(index) != CCTXIndexLength {
		return errors.Wrapf(ErrInvalidIndexValue, "invalid index length %d, expected: %d", len(index), CCTXIndexLength)
	}
	return nil
}

// ValidateHashForChain validates the hash for the chain
// NOTE: since these checks are currently not used, we don't provide additional chains for simplicity
// TODO: use authorityKeeper.GetChainInfo to provide additional chains
// https://github.com/zeta-chain/node/issues/2234
// https://github.com/zeta-chain/node/issues/2385
// NOTE: We should eventually not using these hard-coded checks at all since it might make the protocol too rigid
// Example: hash algorithm is changed for a chain: this required a upgrade on the protocol
func ValidateHashForChain(hash string, chainID int64) error {
	if chains.IsEthereumChain(chainID, []chains.Chain{}) || chains.IsZetaChain(chainID, []chains.Chain{}) {
		_, err := hexutil.Decode(hash)
		if err != nil {
			return fmt.Errorf("hash must be a valid ethereum hash %s", hash)
		}
		return nil
	}
	if chains.IsBitcoinChain(chainID, []chains.Chain{}) {
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
// NOTE: since these checks are currently not used, we don't provide additional chains for simplicity
// TODO: use authorityKeeper.GetChainInfo to provide additional chains
// https://github.com/zeta-chain/node/issues/2234
// https://github.com/zeta-chain/node/issues/2385
// NOTE: We should eventually not using these hard-coded checks at all for same reasons as above
func ValidateAddressForChain(address string, chainID int64) error {
	// we do not validate the address for zeta chain as the address field can be btc or eth address
	if chains.IsZetaChain(chainID, []chains.Chain{}) {
		return nil
	}
	if chains.IsEthereumChain(chainID, []chains.Chain{}) {
		if !ethcommon.IsHexAddress(address) {
			return fmt.Errorf("invalid address %s , chain %d", address, chainID)
		}
		return nil
	}
	if chains.IsBitcoinChain(chainID, []chains.Chain{}) {
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
