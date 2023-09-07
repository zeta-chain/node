package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/common"
)

func Chain(chainID int64) *common.Chain {
	r := newRandFromSeed(chainID)

	return &common.Chain{
		ChainName: common.ChainName(r.Intn(4)),
		ChainId:   chainID,
	}
}

func PubKeySet() *common.PubKeySet {
	pubKeySet := common.PubKeySet{
		Secp256k1: common.PubKey(secp256k1.GenPrivKey().PubKey().Bytes()),
		Ed25519:   common.PubKey(ed25519.GenPrivKey().PubKey().String()),
	}
	return &pubKeySet
}
