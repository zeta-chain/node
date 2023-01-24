package zetaclient

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"math/big"
	"strings"
	"time"
)

type EVMSigner struct {
	client              *ethclient.Client
	chain               *common.Chain
	chainID             *big.Int
	tssSigner           TSSSigner
	ethSigner           ethtypes.Signer
	abi                 abi.ABI
	metaContractAddress ethcommon.Address
	logger              zerolog.Logger
}

func NewEVMSigner(chain common.Chain, endpoint string, tssSigner TSSSigner, abiString string, metaContract ethcommon.Address) (*EVMSigner, error) {
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

	return &EVMSigner{
		client:              client,
		chain:               &chain,
		tssSigner:           tssSigner,
		chainID:             chainID,
		ethSigner:           ethSigner,
		abi:                 abi,
		metaContractAddress: metaContract,
		logger:              log.With().Str("module", "EVMSigner").Logger(),
	}, nil
}

// given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *EVMSigner) Sign(data []byte, to ethcommon.Address, gasLimit uint64, gasPrice *big.Int, nonce uint64) (*ethtypes.Transaction, []byte, []byte, error) {
	tx := ethtypes.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes)
	if err != nil {
		return nil, nil, nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, nil, nil, err
	}
	return signedTX, sig[:], hashBytes[:], nil
}

// takes in signed tx, broadcast to external chain node
func (signer *EVMSigner) Broadcast(tx *ethtypes.Transaction) error {
	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return signer.client.SendTransaction(ctxt, tx)
}

// function onReceive(
//
//	bytes calldata originSenderAddress,
//	uint256 originChainId,
//	address destinationAddress,
//	uint zetaAmount,
//	bytes calldata message,
//	bytes32 internalSendHash
//
// ) external virtual {}
func (signer *EVMSigner) SignOutboundTx(sender ethcommon.Address, srcChainID *big.Int, to ethcommon.Address, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
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

// function onRevert(
// address originSenderAddress,
// uint256 originChainId,
// bytes calldata destinationAddress,
// uint256 destinationChainId,
// uint256 zetaAmount,
// bytes calldata message,
// bytes32 internalSendHash
// ) external override whenNotPaused onlyTssAddress
func (signer *EVMSigner) SignRevertTx(sender ethcommon.Address, srcChainID *big.Int, to []byte, toChainID *big.Int, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	data, err = signer.abi.Pack("onRevert", sender, srcChainID, to, toChainID, amount, message, sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
	}

	return tx, nil
}

func (signer *EVMSigner) SignCancelTx(nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, signer.tssSigner.EVMAddress(), big.NewInt(0), 21000, gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes)
	if err != nil {
		return nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, err
	}

	return signedTX, nil
}

func (signer *EVMSigner) SignWithdrawTx(to ethcommon.Address, amount *big.Int, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes)
	if err != nil {
		return nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, err
	}

	return signedTX, nil
}
