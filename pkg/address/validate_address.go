package address

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

// ValidateAddressForChain validates the address for the chain
// NOTE: since these checks are currently not used, we don't provide additional chains for simplicity
// https://github.com/zeta-chain/node/issues/2234
// https://github.com/zeta-chain/node/issues/2385
// NOTE: We should eventually not using these hard-coded checks at all for same reasons as above

// TODO : Use this function to validate Sender and Receiver address for CCTX
// https://github.com/zeta-chain/node/issues/2697
func ValidateAddressForChain(address string, chainID int64, additionalChains []chains.Chain) error {
	chain, found := chains.GetChainFromChainID(chainID, additionalChains)
	if !found {
		return &InvalidChainError{ChainID: chainID, Msg: "chain not supported"}
	}
	switch chain.Network {
	case chains.Network_eth:
		return ValidateEVMAddress(address)
	case chains.Network_zeta:
		return nil
	case chains.Network_btc:
		return ValidateBTCAddress(address, chainID)
	case chains.Network_polygon:
		return ValidateEVMAddress(address)
	case chains.Network_bsc:
		return ValidateEVMAddress(address)
	case chains.Network_optimism:
		return nil
	case chains.Network_base:
		return nil
	case chains.Network_solana:
		return nil
	default:
		return &InvalidNetworkError{Network: chain.Network.String(), Msg: "network not supported"}
	}
}

func ValidateEVMAddress(address string) error {
	if !ethcommon.IsHexAddress(address) {
		return &InvalidAddressError{Address: address, Msg: "not a valid hex address"}
	}
	if ethcommon.HexToAddress(address) == (ethcommon.Address{}) {
		return &InvalidAddressError{Address: address, Msg: "empty address"}
	}
	return nil
}

func ValidateBTCAddress(address string, chainID int64) error {
	addr, err := chains.DecodeBtcAddress(address, chainID)
	if err != nil {
		return &InvalidAddressError{Address: address, Msg: fmt.Sprintf("failed to decode address: %s chainId %d", err, chainID)}
	}
	if !chains.IsBtcAddressSupported(addr) {
		return &InvalidAddressError{Address: address, Msg: "address not supported"}
	}
	return nil
}
