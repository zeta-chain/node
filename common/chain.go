package common

import (
	"errors"
	"fmt"
	"strings"
)

var (
	EmptyChain   = Chain("")
	BSCChain     = Chain("BSC")
	ETHChain     = Chain("ETH")
	POLYGONChain = Chain("POLYGON")
	ZETAChain    = Chain("ZETA")

	SigningAlgoSecp256k1 = SigninAlgo("secp256k1")
	SigningAlgoEd25519   = SigninAlgo("ed25519")
)

type SigninAlgo string

// Chain is an alias of string , represent a block chain
type Chain string

// Chains represent a slice of Chain
type Chains []Chain

// Validate validates chain format, should consist only of uppercase letters
func (c Chain) Validate() error {
	if len(c) < 3 {
		return errors.New("chain id len is less than 3")
	}
	if len(c) > 10 {
		return errors.New("chain id len is more than 10")
	}
	for _, ch := range string(c) {
		if ch < 'A' || ch > 'Z' {
			return errors.New("chain id can consist only of uppercase letters")
		}
	}
	return nil
}

// NewChain create a new Chain and default the siging_algo to Secp256k1
func NewChain(chainID string) (Chain, error) {
	chain := Chain(strings.ToUpper(chainID))
	if err := chain.Validate(); err != nil {
		return chain, err
	}
	return chain, nil
}

func ParseChain(chainID string) (Chain, error) {
	switch chainID {
	case "ETH":
		return ETHChain, nil
	case "BSC":
		return BSCChain, nil
	case "POLYGON":
		return POLYGONChain, nil
	default:
		return EmptyChain, fmt.Errorf("Unsupported chain %s", chainID)
	}
}

func (chain Chain) GetNativeTokenSymbol() string {
	switch chain {
	case ETHChain:
		return "ETH"
	case BSCChain:
		return "BNB"
	case POLYGONChain:
		return "MATIC"
	default:
		return "" // should not happen
	}
}

// Equals compare two chain to see whether they represent the same chain
func (c Chain) Equals(c2 Chain) bool {
	return strings.EqualFold(c.String(), c2.String())
}

func (c Chain) IsZetaChain() bool {
	return c.Equals(ZETAChain)
}

// IsEmpty is to determinate whether the chain is empty
func (c Chain) IsEmpty() bool {
	return strings.TrimSpace(c.String()) == ""
}

// String implement fmt.Stringer
func (c Chain) String() string {
	// convert it to upper case again just in case someone created a ticker via Chain("rune")
	return strings.ToUpper(string(c))
}

// GetSigningAlgo get the signing algorithm for the given chain
func (c Chain) GetSigningAlgo() SigninAlgo {
	switch c {
	case ETHChain, POLYGONChain, BSCChain:
		return SigningAlgoSecp256k1
	default:
		return SigningAlgoSecp256k1
	}
}

// Has check whether chain c is in the list
func (chains Chains) Has(c Chain) bool {
	for _, ch := range chains {
		if ch.Equals(c) {
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
