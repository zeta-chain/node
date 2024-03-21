package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/pkg"
)

func Chain(chainID int64) *pkg.Chain {
	r := newRandFromSeed(chainID)

	return &pkg.Chain{
		ChainName: pkg.ChainName(r.Intn(4)),
		ChainId:   chainID,
	}
}

func PubKeySet() *pkg.PubKeySet {
	pubKeySet := pkg.PubKeySet{
		Secp256k1: pkg.PubKey(secp256k1.GenPrivKey().PubKey().Bytes()),
		Ed25519:   pkg.PubKey(ed25519.GenPrivKey().PubKey().String()),
	}
	return &pubKeySet
}

func EventIndex() uint64 {
	r := newRandFromSeed(1)
	return r.Uint64()
}
