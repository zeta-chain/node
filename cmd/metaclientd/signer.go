package metaclientd

import (
	"context"
	"github.com/Meta-Protocol/metacore/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

type TSSSigner interface {
	Pubkey() []byte
	Sign(data []byte) [65]byte
	Address() ethcommon.Address
}

type Signer struct {
	client              *ethclient.Client
	nonce               uint64
	chain               common.Chain
	chainID             *big.Int
	tssSigner           TSSSigner
	ethSigner           ethtypes.Signer
	abi                 abi.ABI
	metaContractAddress ethcommon.Address
}

func NewSigner(chain common.Chain, endpoint string, tssAddress ethcommon.Address, tssSigner TSSSigner, abiString string, metaContract ethcommon.Address) (*Signer, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	nonce, err := client.NonceAt(context.TODO(), tssAddress, nil)
	if err != nil {
		return nil, err
	}
	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}
	ethSigner := ethtypes.LatestSignerForChainID(chainID)
	abi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}

	return &Signer{
		client:              client,
		nonce:               nonce,
		chain:               chain,
		tssSigner:           tssSigner,
		chainID:             chainID,
		ethSigner:           ethSigner,
		abi:                 abi,
		metaContractAddress: metaContract,
	}, nil
}

// given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *Signer) Sign(data []byte, to ethcommon.Address, gasLimit uint64, gasPrice *big.Int) (*ethtypes.Transaction, []byte, []byte, error) {
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

// send outbound tx to smart contract
func (signer *Signer) MMint(amount *big.Int, to ethcommon.Address, gasLimit uint64, message []byte) (string, error) {
	data, err := signer.abi.Pack("mint", to, amount)
	if err != nil {
		return "", err
	}
	gasPrice, err := signer.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice)
	if err != nil {
		return "", err
	}
	err = signer.Broadcast(tx)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}
