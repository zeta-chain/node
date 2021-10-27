package metaclientd

import (
	"context"
	"github.com/Meta-Protocol/metacore/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type TSSSigner interface {
	Pubkey() []byte
	Sign(data []byte) [65]byte
	Address() ethcommon.Address
}

type Signer struct {
	client    *ethclient.Client
	nonce     uint64
	chain     common.Chain
	chainID   *big.Int
	tssSigner TSSSigner
	ethSigner ethtypes.Signer
}

func NewSigner(chain common.Chain, endpoint string, tssAddress ethcommon.Address, tssSigner TSSSigner) (*Signer, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	nonce, err := client.NonceAt(context.TODO(), tssAddress, nil)
	if err != nil {
		return nil, err
	}
	chainID, err := client.ChainID(context.TODO())
	ethSigner := ethtypes.LatestSignerForChainID(chainID)

	return &Signer{
		client:    client,
		nonce:     nonce,
		chain:     chain,
		tssSigner: tssSigner,
		chainID:   chainID,
		ethSigner: ethSigner,
	}, nil
}

// given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *Signer) Sign(data []byte, to ethcommon.Address, gasLimit uint64, gasPrice *big.Int) (*ethtypes.Transaction, []byte,[]byte,  error) {
	tx := ethtypes.NewTransaction(signer.nonce, to, big.NewInt(0), gasLimit, gasPrice, data)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig := signer.tssSigner.Sign(hashBytes)
	signer.nonce++
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, nil, nil, err
	}
	return signedTX, sig[:], hashBytes[:], nil
}

// takes in signed tx, broadcast to external chain node
func (signer *Signer) Broadcast(tx *ethtypes.Transaction) error {
	return signer.client.SendTransaction(context.TODO(), tx)
}

