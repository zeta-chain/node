package sample

import (
	"context"
	"math/big"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
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

func Proof() (tx_index int64, block *ethtypes.Block, header ethtypes.Header, headerRLP []byte, proof *common.Proof, tx *ethtypes.Transaction, err error) {
	tx_index = int64(9)
	RPC_URL := "https://rpc.ankr.com/eth_goerli"
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return
	}
	bn := int64(9889649)
	block, err = client.BlockByNumber(context.Background(), big.NewInt(bn))
	headerRLP, _ = rlp.EncodeToBytes(block.Header())
	err = rlp.DecodeBytes(headerRLP, &header)
	if err != nil {
		return
	}
	tr := ethereum.NewTrie(block.Transactions())
	var b []byte
	ib := rlp.AppendUint64(b, uint64(tx_index))
	p := ethereum.NewProof()
	err = tr.Prove(ib, 0, p)
	if err != nil {
		return
	}
	proof = common.NewEthereumProof(p)
	tx = block.Transactions()[tx_index]
	return
}
