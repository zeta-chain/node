package sample

import (
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/crypto"
)

func Chain(chainID int64) *chains.Chain {
	r := newRandFromSeed(chainID)

	return &chains.Chain{
		ChainName: chains.ChainName(r.Intn(4)),
		ChainId:   chainID,
	}
}

func PubKeySet() *crypto.PubKeySet {
	pubKeySet := crypto.PubKeySet{
		Secp256k1: crypto.PubKey(secp256k1.GenPrivKey().PubKey().Bytes()),
		Ed25519:   crypto.PubKey(ed25519.GenPrivKey().PubKey().String()),
	}
	return &pubKeySet
}

func EventIndex() uint64 {
	r := newRandFromSeed(1)
	return r.Uint64()
}
