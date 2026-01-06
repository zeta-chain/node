package common

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zeta-chain/node/pkg/chains"
)

// ResolveChain resolves a chain from either a chain ID (numeric string) or chain name
// Examples:
//   - "7000" -> ZetaChainMainnet (by chain ID)
//   - "zeta_mainnet" -> ZetaChainMainnet (by name)
//   - "eth_mainnet" -> Ethereum (by name)
func ResolveChain(chainIdentifier string) (chains.Chain, error) {
	if chainID, err := strconv.ParseInt(chainIdentifier, 10, 64); err == nil {
		chain, found := chains.GetChainFromChainID(chainID, []chains.Chain{})
		if !found {
			return chains.Chain{}, fmt.Errorf("chain with ID %d not found", chainID)
		}
		return chain, nil
	}
	return GetChainByName(chainIdentifier)
}

// GetChainByName finds a chain by its name (case-insensitive)
func GetChainByName(name string) (chains.Chain, error) {
	allChains := chains.CombineDefaultChainsList([]chains.Chain{})
	nameLower := strings.ToLower(name)

	for _, chain := range allChains {
		if strings.ToLower(chain.Name) == nameLower {
			return chain, nil
		}
	}

	return chains.Chain{}, fmt.Errorf("chain with name %q not found", name)
}

// NetworkTypeFromChain returns the network type string from a chain
func NetworkTypeFromChain(chain chains.Chain) string {
	switch chain.NetworkType {
	case chains.NetworkType_mainnet:
		return "mainnet"
	case chains.NetworkType_testnet:
		return "testnet"
	case chains.NetworkType_privnet:
		return "localnet"
	case chains.NetworkType_devnet:
		return "devnet"
	default:
		return "mainnet"
	}
}

// ListAvailableChains returns a formatted string of all available chains
func ListAvailableChains() string {
	allChains := chains.CombineDefaultChainsList([]chains.Chain{})
	var sb strings.Builder

	sb.WriteString("Available chains:\n")
	for _, chain := range allChains {
		sb.WriteString(fmt.Sprintf("  %s (ID: %d)\n", chain.Name, chain.ChainId))
	}

	return sb.String()
}
