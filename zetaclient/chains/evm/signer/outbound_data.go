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
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
)

const (
	MinGasLimit = 100_000
	MaxGasLimit = 1_000_000
)

// OutboundData is a data structure containing input fields used to construct each type of transaction.
// This is populated using cctx and other input parameters passed to TryProcessOutbound
type OutboundData struct {
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

	// cctxIndex field is the inbound message digest that is sent to the destination contract
	cctxIndex [32]byte

	// outboundParams field contains data detailing the receiver chain and outbound transaction
	outboundParams *types.OutboundParams
}

// SetChainAndSender populates the destination address and Chain ID based on the status of the cross chain tx
// returns true if transaction should be skipped
// returns false otherwise
func (txData *OutboundData) SetChainAndSender(cctx *types.CrossChainTx, logger zerolog.Logger) bool {
	switch cctx.CctxStatus.Status {
	case types.CctxStatus_PendingRevert:
		txData.to = ethcommon.HexToAddress(cctx.InboundParams.Sender)
		txData.toChainID = big.NewInt(cctx.InboundParams.SenderChainId)
		logger.Info().Msgf("Abort: reverting inbound")
	case types.CctxStatus_PendingOutbound:
		txData.to = ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
	default:
		logger.Info().Msgf("Transaction doesn't need to be processed status: %d", cctx.CctxStatus.Status)
		return true
	}
	return false
}

// SetupGas sets the gas limit and price
func (txData *OutboundData) SetupGas(
	cctx *types.CrossChainTx,
	logger zerolog.Logger,
	client interfaces.EVMRPCClient,
	chain chains.Chain,
) error {
	txData.gasLimit = cctx.GetCurrentOutboundParam().GasLimit
	if txData.gasLimit < MinGasLimit {
		txData.gasLimit = MinGasLimit
		logger.Warn().
			Msgf("gasLimit %d is too low; set to %d", cctx.GetCurrentOutboundParam().GasLimit, txData.gasLimit)
	}
	if txData.gasLimit > MaxGasLimit {
		txData.gasLimit = MaxGasLimit
		logger.Warn().
			Msgf("gasLimit %d is too high; set to %d", cctx.GetCurrentOutboundParam().GasLimit, txData.gasLimit)
	}

	// use dynamic gas price for ethereum chains.
	// The code below is a fix for https://github.com/zeta-chain/node/issues/1085
	// doesn't close directly the issue because we should determine if we want to keep using SuggestGasPrice if no GasPrice
	// we should possibly remove it completely and return an error if no GasPrice is provided because it means no fee is processed on ZetaChain
	specified, ok := new(big.Int).SetString(cctx.GetCurrentOutboundParam().GasPrice, 10)
	if !ok {
		if chain.Network == chains.Network_eth {
			suggested, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				return errors.Join(err, fmt.Errorf("cannot get gas price from chain %s ", chain.String()))
			}
			txData.gasPrice = roundUpToNearestGwei(suggested)
		} else {
			return fmt.Errorf("cannot convert gas price  %s ", cctx.GetCurrentOutboundParam().GasPrice)
		}
	} else {
		txData.gasPrice = specified
	}
	return nil
}

// NewOutboundData populates transaction input fields parsed from the cctx and other parameters
// returns
//  1. New NewOutboundData Data struct or nil if an error occurred.
//  2. bool (skipTx) - if the transaction doesn't qualify to be processed the function will return true, meaning that this
//     cctx will be skipped and false otherwise.
//  3. error
func NewOutboundData(
	ctx context.Context,
	cctx *types.CrossChainTx,
	evmObserver *observer.Observer,
	evmRPC interfaces.EVMRPCClient,
	logger zerolog.Logger,
	height uint64,
) (*OutboundData, bool, error) {
	txData := OutboundData{}
	txData.outboundParams = cctx.GetCurrentOutboundParam()
	txData.amount = cctx.GetCurrentOutboundParam().Amount.BigInt()
	txData.nonce = cctx.GetCurrentOutboundParam().TssNonce
	txData.sender = ethcommon.HexToAddress(cctx.InboundParams.Sender)
	txData.srcChainID = big.NewInt(cctx.InboundParams.SenderChainId)
	txData.asset = ethcommon.HexToAddress(cctx.InboundParams.Asset)

	txData.height = height

	skipTx := txData.SetChainAndSender(cctx, logger)
	if skipTx {
		return nil, true, nil
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, false, err
	}

	toChain, found := chains.GetChainFromChainID(txData.toChainID.Int64(), app.GetAdditionalChains())
	if !found {
		return nil, true, fmt.Errorf("unknown chain: %d", txData.toChainID.Int64())
	}

	// Get nonce, Early return if the cctx is already processed
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	included, confirmed, err := evmObserver.IsOutboundProcessed(ctx, cctx, logger)
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
	logger.Info().
		Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain.String(), cctx.InboundParams.Amount, txData.to.Hex(), nonce, cctx.InboundParams.FinalizedZetaHeight)
	cctxIndex, err := hex.DecodeString(cctx.Index[2:]) // remove the leading 0x
	if err != nil || len(cctxIndex) != 32 {
		return nil, true, fmt.Errorf("decode CCTX %s error", cctx.Index)
	}
	copy(txData.cctxIndex[:32], cctxIndex[:32])

	// In case there is a pending transaction, make sure this keysign is a transaction replacement
	pendingTx := evmObserver.GetPendingTx(nonce)
	if pendingTx != nil {
		if txData.gasPrice.Cmp(pendingTx.GasPrice()) > 0 {
			logger.Info().
				Msgf("replace pending outbound %s nonce %d using gas price %d", pendingTx.Hash().Hex(), nonce, txData.gasPrice)
		} else {
			logger.Info().Msgf("please wait for pending outbound %s nonce %d to be included", pendingTx.Hash().Hex(), nonce)
			return nil, true, nil
		}
	}

	// Base64 decode message
	if cctx.InboundParams.CoinType != coin.CoinType_Cmd {
		txData.message, err = base64.StdEncoding.DecodeString(cctx.RelayedMessage)
		if err != nil {
			logger.Err(err).Msgf("decode CCTX.Message %s error", cctx.RelayedMessage)
		}
	}

	return &txData, false, nil
}
