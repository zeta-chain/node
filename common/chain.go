package common

import (
	"strings"
)

var (
	SigningAlgoSecp256k1 = SigninAlgo("secp256k1")
	SigningAlgoEd25519   = SigninAlgo("ed25519")
)

// FIXME: should it return err?
func ParseStringToObserverChain(chain string) ChainName {
	c := ChainName_value[chain]
	return ChainName(c)
}

type SigninAlgo string

// Chain is an alias of string , represent a block chain
//type Chain string

// Chains represent a slice of Chain
type Chains []Chain

//// Validate validates chain format, should consist only of uppercase letters
//func (chain Chain) Validate() error {
//	if len(chain) < 3 {
//		return errors.New("chain id len is less than 3")
//	}
//	if len(chain) > 10 {
//		return errors.New("chain id len is more than 10")
//	}
//	for _, ch := range string(chain) {
//		if ch < 'A' || ch > 'Z' {
//			return errors.New("chain id can consist only of uppercase letters")
//		}
//	}
//	return nil
//}
//
//// NewChain create a new Chain and default the siging_algo to Secp256k1
//func NewChain(chainID string) (Chain, error) {
//	chain := Chain(strings.ToUpper(chainID))
//	if err := chain.Validate(); err != nil {
//		return chain, err
//	}
//	return chain, nil
//}

// Equals compare two chain to see whether they represent the same chain
func (chain Chain) IsEqual(c Chain) bool {
	if chain.ChainName == c.ChainName && chain.ChainId == c.ChainId {
		return true
	}
	return false
}

func (chain Chain) IsZetaChain() bool {
	return chain.IsEqual(ZetaChain()) || chain.IsEqual(ZetaLocalNetChain()) || chain.IsEqual(ZetaAthensNetChain())
}

func (chain Chain) IsEVMChain() bool {
	return chain.IsEqual(EthChain()) || chain.IsEqual(BscMainnetChain()) || chain.IsEqual(PolygonChain()) || chain.IsEqual(GoerliChain()) ||
		chain.IsEqual(MumbaiChain()) || chain.IsEqual(BscTestnetChain()) || chain.IsEqual(BaobabChain()) || chain.IsEqual(GoeriliLocalNetChain())
}

func (chain Chain) IsKlaytnChain() bool {
	return chain.IsEqual(BaobabChain())
}

func (chain Chain) IsBitcoinChain() bool {
	return chain.IsEqual(BtcMainnetChain()) || chain.IsEqual(BtcTestNetChain())
}

// IsEmpty is to determinate whether the chain is empty
func (chain Chain) IsEmpty() bool {
	return strings.TrimSpace(chain.String()) == ""
}

// GetSigningAlgo get the signing algorithm for the given chain
func (chain Chain) GetSigningAlgo() SigninAlgo {
	switch chain {
	case EthChain(), PolygonChain(), BscMainnetChain(), BtcMainnetChain(), BtcTestNetChain():
		return SigningAlgoSecp256k1
	default:
		return SigningAlgoSecp256k1
	}
}

// Has check whether chain c is in the list
func (chains Chains) Has(c Chain) bool {
	for _, ch := range chains {
		if ch.IsEqual(c) {
			return true
		}
	}
	return false
}

// Distinct return a distinct set of chains, no duplicates
func (chains Chains) Distinct() Chains {
	var newChains Chains
	for _, chain := range chains {
		if !newChains.Has(chain) {
			newChains = append(newChains, chain)
		}
	}
	return newChains
}

func (chains Chains) Strings() []string {
	strings := make([]string, len(chains))
	for i, c := range chains {
		strings[i] = c.String()
	}
	return strings
}
