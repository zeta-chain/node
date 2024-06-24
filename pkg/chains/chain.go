package chains

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Validate checks whether the chain is valid
// The function check the chain ID is positive and all enum fields have a defined value
func (chain Chain) Validate() error {
	if chain.ChainId <= 0 {
		return fmt.Errorf("chain ID must be positive")
	}

	if _, ok := ChainName_name[int32(chain.ChainName)]; !ok {
		return fmt.Errorf("invalid chain name %d", int32(chain.ChainName))
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
	if IsEVMChain(chain.ChainId) {
		addr := ethcommon.BytesToAddress(b)
		if addr == (ethcommon.Address{}) {
			return "", fmt.Errorf("invalid EVM address")
		}
		return addr.Hex(), nil
	} else if IsBitcoinChain(chain.ChainId) {
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
	}
	return "", fmt.Errorf("chain (%d) not supported", chain.ChainId)
}

// DecodeAddressFromChainID decode the address string to bytes
func DecodeAddressFromChainID(chainID int64, addr string) ([]byte, error) {
	if IsEVMChain(chainID) {
		return ethcommon.HexToAddress(addr).Bytes(), nil
	} else if IsBitcoinChain(chainID) {
		return []byte(addr), nil
	}
	return nil, fmt.Errorf("chain (%d) not supported", chainID)
}

// IsEVMChain returns true if the chain is an EVM chain or uses the ethereum consensus mechanism for block finality
func IsEVMChain(chainID int64) bool {
	return ChainIDInChainList(chainID, ChainListByConsensus(Consensus_ethereum))
}

// IsBitcoinChain returns true if the chain is a Bitcoin-based chain or uses the bitcoin consensus mechanism for block finality
func IsBitcoinChain(chainID int64) bool {
	return ChainIDInChainList(chainID, ChainListByConsensus(Consensus_bitcoin))
}

// IsEthereumChain returns true if the chain is an Ethereum chain
func IsEthereumChain(chainID int64) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_eth))
}

// IsZetaChain returns true if the chain is a Zeta chain
func IsZetaChain(chainID int64) bool {
	return ChainIDInChainList(chainID, ChainListByNetwork(Network_zeta))
}

// IsHeaderSupportedChain returns true if the chain's consensus supports block header-based verification
func IsHeaderSupportedChain(chainID int64) bool {
	return ChainIDInChainList(chainID, ChainListForHeaderSupport())
}

// SupportMerkleProof returns true if the chain supports block header-based verification
func (chain Chain) SupportMerkleProof() bool {
	return IsEVMChain(chain.ChainId) || IsBitcoinChain(chain.ChainId)
}

// IsEmpty is to determinate whether the chain is empty
func (chain Chain) IsEmpty() bool {
	return strings.TrimSpace(chain.String()) == ""
}

// GetChainFromChainID returns the chain from the chain ID
func GetChainFromChainID(chainID int64) *Chain {
	chains := DefaultChainsList()
	for _, chain := range chains {
		if chainID == chain.ChainId {
			return chain
		}
	}
	return nil
}

// GetBTCChainParams returns the bitcoin chain config params from the chain ID
func GetBTCChainParams(chainID int64) (*chaincfg.Params, error) {
	switch chainID {
	case 18444:
		return &chaincfg.RegressionNetParams, nil
	case 18332:
		return &chaincfg.TestNet3Params, nil
	case 8332:
		return &chaincfg.MainNetParams, nil
	default:
		return nil, fmt.Errorf("error chainID %d is not a bitcoin chain", chainID)
	}
}

// GetBTCChainIDFromChainParams returns the bitcoin chain ID from the chain config params
func GetBTCChainIDFromChainParams(params *chaincfg.Params) (int64, error) {
	switch params.Name {
	case chaincfg.RegressionNetParams.Name:
		return 18444, nil
	case chaincfg.TestNet3Params.Name:
		return 18332, nil
	case chaincfg.MainNetParams.Name:
		return 8332, nil
	default:
		return 0, fmt.Errorf("error chain %s is not a bitcoin chain", params.Name)
	}
}

// ChainIDInChainList checks whether the chainID is in the chain list
func ChainIDInChainList(chainID int64, chainList []*Chain) bool {
	for _, c := range chainList {
		if chainID == c.ChainId {
			return true
		}
	}
	return false
}
