package chains

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// Validate checks whether the chain is valid
// The function check the chain ID is positive and all enum fields have a defined value
func (chain Chain) Validate() error {
	if chain.ChainId <= 0 {
		return fmt.Errorf("chain ID must be positive")
	}

	if _, ok := Network_name[int32(chain.Network)]; !ok {
		return fmt.Errorf("invalid network %d", int32(chain.Network))
	}

	if _, ok := NetworkType_name[int32(chain.NetworkType)]; !ok {
		return fmt.Errorf("invalid network type %d", int32(chain.NetworkType))
	}

	if _, ok := Vm_name[int32(chain.Vm)]; !ok {
		return fmt.Errorf("invalid vm %d", int32(chain.Vm))
	}

	if _, ok := Consensus_name[int32(chain.Consensus)]; !ok {
		return fmt.Errorf("invalid consensus %d", int32(chain.Consensus))
	}

	if chain.Name == "" {
		return errors.New("chain name cannot be empty")
	}

	return nil
}

// IsZetaChain returns true if the chain is a ZetaChain chain
func (chain Chain) IsZetaChain() bool {
	return chain.Network == Network_zeta
}

// IsExternalChain returns true if the chain is an ExternalChain chain, not ZetaChain
func (chain Chain) IsExternalChain() bool {
	return chain.IsExternal
}

// EncodeAddress bytes representations of address
// on EVM chain, it is 20Bytes
// on Bitcoin chain, it is P2WPKH address, []byte(bech32 encoded string)
func (chain Chain) EncodeAddress(b []byte) (string, error) {
	if chain.Vm == Vm_evm {
		addr := ethcommon.BytesToAddress(b)
		if addr == (ethcommon.Address{}) {
			return "", fmt.Errorf("invalid EVM address")
		}
		return addr.Hex(), nil
	}
	switch chain.Consensus {
	case Consensus_bitcoin:
		addrStr := string(b)
		chainParams, err := GetBTCChainParams(chain.ChainId)
		if err != nil {
			return "", err
		}
		addr, err := DecodeBtcAddress(addrStr, chain.ChainId)
		if err != nil {
			return "", err
		}
		if !addr.IsForNet(chainParams) {
			return "", fmt.Errorf("address is not for network %s", chainParams.Name)
		}
		return addrStr, nil
	case Consensus_solana_consensus:
		pk, err := DecodeSolanaWalletAddress(string(b))
		if err != nil {
			return "", err
		}
		return pk.String(), nil
	case Consensus_catchain_consensus:
		acc, err := ton.ParseAccountID(string(b))
		if err != nil {
			return "", err
		}
		return acc.ToRaw(), nil
	case Consensus_sui_consensus:
		return "0x" + hex.EncodeToString(b), nil
	default:
		return "", fmt.Errorf("chain id %d not supported", chain.ChainId)
	}
}

func (chain Chain) IsEVMChain() bool {
	return chain.Vm == Vm_evm
}

func (chain Chain) IsSolanaChain() bool {
	return chain.Consensus == Consensus_solana_consensus
}

func (chain Chain) IsBitcoinChain() bool {
	return chain.Consensus == Consensus_bitcoin
}

func (chain Chain) IsTONChain() bool {
	return chain.Consensus == Consensus_catchain_consensus
}

func (chain Chain) LogFields() map[string]any {
	return map[string]any{
		logs.FieldChain:   chain.ChainId,
		logs.FieldNetwork: chain.Network.String(),
	}
}

// DecodeAddressFromChainID decode the address string to bytes
// additionalChains is a list of additional chains to search from
// in practice, it is used in the protocol to dynamically support new chains without doing an upgrade
func DecodeAddressFromChainID(chainID int64, addr string, additionalChains []Chain) ([]byte, error) {
	switch {
	case IsEVMChain(chainID, additionalChains):
		return ethcommon.HexToAddress(addr).Bytes(), nil
	case IsBitcoinChain(chainID, additionalChains):
		return []byte(addr), nil
	case IsSolanaChain(chainID, additionalChains):
		return []byte(addr), nil
	case IsTONChain(chainID, additionalChains):
		// e.g. `0:55798cb7b87168251a7c39f6806b8c202f6caa0f617a76f4070b3fdacfd056a1`
		acc, err := ton.ParseAccountID(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid TON address %q: %w", addr, err)
		}
		return []byte(acc.ToRaw()), nil
	case IsSuiChain(chainID, additionalChains):
		addrBytes, err := sui.EncodeAddress(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid Sui address %q: %w", addr, err)
		}
		return addrBytes, nil
	default:
		return nil, fmt.Errorf("chain (%d) not supported", chainID)
	}
}

// IsEVMChain returns true if the chain is an EVM chain
// additionalChains is a list of additional chains to search from
// in practice, it is used in the protocol to dynamically support new chains without doing an upgrade
func IsEVMChain(chainID int64, additionalChains []Chain) bool {
	chain, found := GetChainFromChainID(chainID, additionalChains)
	if !found {
		return false
	}
	return chain.IsEVMChain()
}

// IsBitcoinChain returns true if the chain is a Bitcoin-based chain or uses the bitcoin consensus mechanism for block finality
// additionalChains is a list of additional chains to search from
// in practice, it is used in the protocol to dynamically support new chains without doing an upgrade
func IsBitcoinChain(chainID int64, additionalChains []Chain) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_btc, additionalChains))
}

// IsSolanaChain returns true if the chain is a Solana chain
func IsSolanaChain(chainID int64, additionalChains []Chain) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_solana, additionalChains))
}

// IsSuiChain returns true if the chain is Sui chain.
func IsSuiChain(chainID int64, additionalChains []Chain) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_sui, additionalChains))
}

// IsTONChain returns true is the chain is TON chain
func IsTONChain(chainID int64, additionalChains []Chain) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_ton, additionalChains))
}

// IsEthereumChain returns true if the chain is an Ethereum chain
// additionalChains is a list of additional chains to search from
// in practice, it is used in the protocol to dynamically support new chains without doing an upgrade
func IsEthereumChain(chainID int64, additionalChains []Chain) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_eth, additionalChains))
}

// IsZetaChain returns true if the chain is a Zeta chain
// additionalChains is a list of additional chains to search from
// in practice, it is used in the protocol to dynamically support new chains without doing an upgrade
func IsZetaChain(chainID int64, additionalChains []Chain) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_zeta, additionalChains))
}

// IsEmpty is to determinate whether the chain is empty
func (chain Chain) IsEmpty() bool {
	return strings.TrimSpace(chain.String()) == ""
}

// GetChainFromChainID returns the chain from the chain ID
// additionalChains is a list of additional chains to search from
// in practice, it is used in the protocol to dynamically support new chains without doing an upgrade
func GetChainFromChainID(chainID int64, additionalChains []Chain) (Chain, bool) {
	chains := CombineDefaultChainsList(additionalChains)
	for _, chain := range chains {
		if chainID == chain.ChainId {
			return chain, true
		}
	}
	return Chain{}, false
}

// GetBTCChainParams returns the bitcoin chain config params from the chain ID
func GetBTCChainParams(chainID int64) (*chaincfg.Params, error) {
	switch chainID {
	case BitcoinRegtest.ChainId:
		return &chaincfg.RegressionNetParams, nil
	case BitcoinTestnet.ChainId:
		return &chaincfg.TestNet3Params, nil
	case BitcoinMainnet.ChainId:
		return &chaincfg.MainNetParams, nil
	case BitcoinSignetTestnet.ChainId:
		return &chaincfg.SigNetParams, nil
	case BitcoinTestnet4.ChainId:
		return &TestNet4Params, nil
	default:
		return nil, fmt.Errorf("error chainID %d is not a bitcoin chain", chainID)
	}
}

// GetBTCChainIDFromChainParams returns the bitcoin chain ID from the chain config params
func GetBTCChainIDFromChainParams(params *chaincfg.Params) (int64, error) {
	switch params.Name {
	case chaincfg.RegressionNetParams.Name:
		return BitcoinRegtest.ChainId, nil
	case chaincfg.TestNet3Params.Name:
		return BitcoinTestnet.ChainId, nil
	case chaincfg.MainNetParams.Name:
		return BitcoinMainnet.ChainId, nil
	case chaincfg.SigNetParams.Name:
		return BitcoinSignetTestnet.ChainId, nil
	case TestNet4Params.Name:
		return BitcoinTestnet4.ChainId, nil
	default:
		return 0, fmt.Errorf("error chain %s is not a bitcoin chain", params.Name)
	}
}

// ChainIDInChainList checks whether the chainID is in the chain list
func ChainIDInChainList(chainID int64, chainList []Chain) bool {
	for _, c := range chainList {
		if chainID == c.ChainId {
			return true
		}
	}
	return false
}
