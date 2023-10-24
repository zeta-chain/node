package sample

import (
	"context"
	"math/big"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
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

func EthHeader() (headerRLP []byte, err error) {
	url := "https://rpc.ankr.com/eth_goerli"
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	bn := int64(9889649)
	block, err := client.BlockByNumber(context.Background(), big.NewInt(bn))
	if err != nil {
		return
	}
	headerRLP, _ = rlp.EncodeToBytes(block.Header())
	return
}
