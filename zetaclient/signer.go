package zetaclient

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"math/big"
	"strings"
	"time"
)

type TSSSigner interface {
	Pubkey() []byte
	Sign(data []byte) ([65]byte, error)
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

func NewSigner(chain common.Chain, endpoint string, tssSigner TSSSigner, abiString string, metaContract ethcommon.Address) (*Signer, error) {
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
	sig, err := signer.tssSigner.Sign(hashBytes)
	if err != nil {
		return nil, nil, nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		log.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	log.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, nil, nil, err
	}
	return signedTX, sig[:], hashBytes[:], nil
}

// takes in signed tx, broadcast to external chain node
func (signer *Signer) Broadcast(tx *ethtypes.Transaction) error {
	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return signer.client.SendTransaction(ctxt, tx)
}

//    function onReceive(
//        bytes calldata originSenderAddress,
//        uint256 originChainId,
//        address destinationAddress,
//        uint zetaAmount,
//        bytes calldata message,
//        bytes32 internalSendHash
//    ) external virtual {}
func (signer *Signer) SignOutboundTx(sender ethcommon.Address, srcChainID *big.Int, to ethcommon.Address, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	if len(sendHash) < 32 {
		return nil, fmt.Errorf("sendHash len %d must be 32", len(sendHash))
	}
	var data []byte
	var err error

	data, err = signer.abi.Pack("onReceive", sender.Bytes(), srcChainID, to, amount, message, sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
	}

	return tx, nil
}

//function onRevert(
//address originSenderAddress,
//uint256 originChainId,
//bytes calldata destinationAddress,
//uint256 destinationChainId,
//uint256 zetaAmount,
//bytes calldata message,
//bytes32 internalSendHash
//) external override whenNotPaused onlyTssAddress
func (signer *Signer) SignRevertTx(sender ethcommon.Address, srcChainID *big.Int, to []byte, toChainId *big.Int, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	if len(sendHash) < 32 {
		return nil, fmt.Errorf("sendHash len %d must be 32", len(sendHash))
	}
	var data []byte
	var err error

	data, err = signer.abi.Pack("onRevert", sender, srcChainID, to, toChainId, amount, message, sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
	}

	return tx, nil
}

type TestSigner struct {
	PrivKey *ecdsa.PrivateKey
}

func (s TestSigner) Sign(digest []byte) ([65]byte, error) {
	sig, err := crypto.Sign(digest, s.PrivKey)
	if err != nil {
		return [65]byte{}, err
	}
	var sigbyte [65]byte
	copy(sigbyte[:], sig[:65])
	return sigbyte, nil
}

func (s TestSigner) Pubkey() []byte {
	publicKeyBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	return publicKeyBytes
}

func (s TestSigner) Address() ethcommon.Address {
	return crypto.PubkeyToAddress(s.PrivKey.PublicKey)
}
