package address

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

func ValidateAddressForChain(address string, chainID int64, additionalChains []chains.Chain) error {
	chain, found := chains.GetChainFromChainID(chainID, additionalChains)
	if !found {
		return fmt.Errorf("chain id %d not supported", chainID)
	}
	switch chain.Network {
	case chains.Network_eth:
		return ValidateEthereumAddress(address)
	case chains.Network_zeta:
		return nil
	case chains.Network_btc:
		{
			addr, err := chains.DecodeBtcAddress(address, chainID)
			if err != nil {
				return fmt.Errorf("invalid address %s , chain %d: %s", address, chainID, err)
			}
			if !chains.IsBtcAddressSupported(addr) {
				return fmt.Errorf("unsupported address %s", address)
			}
			return nil
		}
	case chains.Network_polygon:
		return ValidateEthereumAddress(address)
	case chains.Network_bsc:
		return ValidateEthereumAddress(address)
	case chains.Network_optimism:
		return nil
	case chains.Network_base:
		return nil
	case chains.Network_solana:
		return nil
	default:
		return fmt.Errorf("invalid network %d", chain.Network)
	}
}

func ValidateEthereumAddress(address string) error {
	if !ethcommon.IsHexAddress(address) {
		return fmt.Errorf("invalid address %s ", address)
	}
	return nil
}
