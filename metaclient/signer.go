package metaclient

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
func (signer *Signer) Sign(data []byte, to ethcommon.Address, gasLimit uint64, gasPrice *big.Int, nonce uint64) (*ethtypes.Transaction, []byte, []byte, error) {
	tx := ethtypes.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig := signer.tssSigner.Sign(hashBytes)
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
func (signer *Signer) MMint(amount *big.Int, to ethcommon.Address, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (string, error) {
	if len(sendHash) < 32 {
		return "", fmt.Errorf("sendHash len %d must be 32", len(sendHash))
	}
	data, err := signer.abi.Pack("mint", to, amount, sendHash)
	if err != nil {
		return "", fmt.Errorf("pack error: %w", err)
	}
	//gasPriceSuggested, err := signer.client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return "", fmt.Errorf("SuggestGasPrice error: %w", err)
	//}
	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce)
	if err != nil {
		return "", fmt.Errorf("sign error: %w", err)
	}
	err = signer.Broadcast(tx)
	if err != nil {
		return "", fmt.Errorf("Broadcast error: %w", err)
	}
	return tx.Hash().Hex(), nil
}

type TestSigner struct {
	PrivKey *ecdsa.PrivateKey
}

func (s TestSigner) Sign(digest []byte) [65]byte {
	sig, _ := crypto.Sign(digest, s.PrivKey)
	var sigbyte [65]byte
	copy(sigbyte[:], sig[:65])
	return sigbyte
}

func (s TestSigner) Pubkey() []byte {
	publicKeyBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	return publicKeyBytes
}

func (s TestSigner) Address() ethcommon.Address {
	return crypto.PubkeyToAddress(s.PrivKey.PublicKey)
}
