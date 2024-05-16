package signer

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

const (
	MinGasLimit = 100_000
	MaxGasLimit = 1_000_000
)

// OutboundTransactionData is a data structure containing input fields used to construct each type of transaction.
// This is populated using cctx and other input parameters passed to TryProcessOutTx
type OutboundTransactionData struct {
	srcChainID *big.Int
	toChainID  *big.Int
	sender     ethcommon.Address
	to         ethcommon.Address
	asset      ethcommon.Address
	amount     *big.Int
	gasPrice   *big.Int
	gasLimit   uint64
	message    []byte
	nonce      uint64
	height     uint64

	// sendHash field is the inbound message digest that is sent to the destination contract
	sendHash [32]byte

	// outboundParams field contains data detailing the receiver chain and outbound transaction
	outboundParams *types.OutboundTxParams
}

// SetChainAndSender populates the destination address and Chain ID based on the status of the cross chain tx
// returns true if transaction should be skipped
// returns false otherwise
func (txData *OutboundTransactionData) SetChainAndSender(cctx *types.CrossChainTx, logger zerolog.Logger) bool {
	switch cctx.CctxStatus.Status {
	case types.CctxStatus_PendingRevert:
		txData.to = ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
		txData.toChainID = big.NewInt(cctx.InboundTxParams.SenderChainId)
		logger.Info().Msgf("Abort: reverting inbound")
	case types.CctxStatus_PendingOutbound:
		txData.to = ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutTxParam().ReceiverChainId)
	default:
		logger.Info().Msgf("Transaction doesn't need to be processed status: %d", cctx.CctxStatus.Status)
		return true
	}
	return false
}

// SetupGas sets the gas limit and price
func (txData *OutboundTransactionData) SetupGas(
	cctx *types.CrossChainTx,
	logger zerolog.Logger,
	client interfaces.EVMRPCClient,
	chain *chains.Chain,
) error {

	txData.gasLimit = cctx.GetCurrentOutTxParam().OutboundTxGasLimit
	if txData.gasLimit < MinGasLimit {
		txData.gasLimit = MinGasLimit
		logger.Warn().Msgf("gasLimit %d is too low; set to %d", cctx.GetCurrentOutTxParam().OutboundTxGasLimit, txData.gasLimit)
	}
	if txData.gasLimit > MaxGasLimit {
		txData.gasLimit = MaxGasLimit
		logger.Warn().Msgf("gasLimit %d is too high; set to %d", cctx.GetCurrentOutTxParam().OutboundTxGasLimit, txData.gasLimit)
	}

	// use dynamic gas price for ethereum chains.
	// The code below is a fix for https://github.com/zeta-chain/node/issues/1085
	// doesn't close directly the issue because we should determine if we want to keep using SuggestGasPrice if no OutboundTxGasPrice
	// we should possibly remove it completely and return an error if no OutboundTxGasPrice is provided because it means no fee is processed on ZetaChain
	specified, ok := new(big.Int).SetString(cctx.GetCurrentOutTxParam().OutboundTxGasPrice, 10)
	if !ok {
		if chains.IsEthereumChain(chain.ChainId) {
			suggested, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				return errors.Join(err, fmt.Errorf("cannot get gas price from chain %s ", chain))
			}
			txData.gasPrice = roundUpToNearestGwei(suggested)
		} else {
			return fmt.Errorf("cannot convert gas price  %s ", cctx.GetCurrentOutTxParam().OutboundTxGasPrice)
		}
	} else {
		txData.gasPrice = specified
	}
	return nil
}

// NewOutBoundTransactionData populates transaction input fields parsed from the cctx and other parameters
// returns
//  1. New OutBoundTransaction Data struct or nil if an error occurred.
//  2. bool (skipTx) - if the transaction doesn't qualify to be processed the function will return true, meaning that this
//     cctx will be skipped and false otherwise.
//  3. error
func NewOutBoundTransactionData(
	cctx *types.CrossChainTx,
	evmObserver *observer.Observer,
	evmRPC interfaces.EVMRPCClient,
	logger zerolog.Logger,
	height uint64,
) (*OutboundTransactionData, bool, error) {
	txData := OutboundTransactionData{}
	txData.outboundParams = cctx.GetCurrentOutTxParam()
	txData.amount = cctx.GetCurrentOutTxParam().Amount.BigInt()
	txData.nonce = cctx.GetCurrentOutTxParam().OutboundTxTssNonce
	txData.sender = ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
	txData.srcChainID = big.NewInt(cctx.InboundTxParams.SenderChainId)
	txData.asset = ethcommon.HexToAddress(cctx.InboundTxParams.Asset)
	txData.height = height

	skipTx := txData.SetChainAndSender(cctx, logger)
	if skipTx {
		return nil, true, nil
	}

	toChain := chains.GetChainFromChainID(txData.toChainID.Int64())
	if toChain == nil {
		return nil, true, fmt.Errorf("unknown chain: %d", txData.toChainID.Int64())
	}

	// Get nonce, Early return if the cctx is already processed
	nonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce
	included, confirmed, err := evmObserver.IsOutboundProcessed(cctx, logger)
	if err != nil {
		return nil, true, errors.New("IsOutboundProcessed failed")
	}
	if included || confirmed {
		logger.Info().Msgf("CCTX already processed; exit signer")
		return nil, true, nil
	}

	// Set up gas limit and gas price
	err = txData.SetupGas(cctx, logger, evmRPC, toChain)
	if err != nil {
		return nil, true, err
	}

	// Get sendHash
	logger.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain, cctx.InboundTxParams.Amount, txData.to.Hex(), nonce, cctx.InboundTxParams.InboundTxFinalizedZetaHeight)
	sendHash, err := hex.DecodeString(cctx.Index[2:]) // remove the leading 0x
	if err != nil || len(sendHash) != 32 {
		return nil, true, fmt.Errorf("decode CCTX %s error", cctx.Index)
	}
	copy(txData.sendHash[:32], sendHash[:32])

	// In case there is a pending transaction, make sure this keysign is a transaction replacement
	pendingTx := evmObserver.GetPendingTx(nonce)
	if pendingTx != nil {
		if txData.gasPrice.Cmp(pendingTx.GasPrice()) > 0 {
			logger.Info().Msgf("replace pending outTx %s nonce %d using gas price %d", pendingTx.Hash().Hex(), nonce, txData.gasPrice)
		} else {
			logger.Info().Msgf("please wait for pending outTx %s nonce %d to be included", pendingTx.Hash().Hex(), nonce)
			return nil, true, nil
		}
	}

	// Base64 decode message
	if cctx.InboundTxParams.CoinType != coin.CoinType_Cmd {
		txData.message, err = base64.StdEncoding.DecodeString(cctx.RelayedMessage)
		if err != nil {
			logger.Err(err).Msgf("decode CCTX.Message %s error", cctx.RelayedMessage)
		}
	}

	return &txData, false, nil
}
