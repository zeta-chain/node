package signer

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
)

const (
	MinGasLimit = 100_000
	MaxGasLimit = 1_000_000
)

// OutboundData is a data structure containing input fields used to construct each type of transaction.
// This is populated using cctx and other input parameters passed to TryProcessOutbound
type OutboundData struct {
	srcChainID *big.Int
	sender     ethcommon.Address

	toChainID *big.Int
	to        ethcommon.Address

	asset  ethcommon.Address
	amount *big.Int

	gas    gas
	nonce  uint64
	height uint64

	message []byte

	// cctxIndex field is the inbound message digest that is sent to the destination contract
	cctxIndex [32]byte

	// outboundParams field contains data detailing the receiver chain and outbound transaction
	outboundParams *types.OutboundParams
}

// NewOutboundData populates transaction input fields parsed from the cctx and other parameters
// returns
//  1. New NewOutboundData Data struct or nil if an error occurred.
//  2. bool (skipTx) - if the transaction doesn't qualify to be processed the function will return true, meaning that this
//     cctx will be skipped and false otherwise.
//  3. error
func NewOutboundData(
	cctx *types.CrossChainTx,
	evmObserver *observer.Observer,
	height uint64,
	logger zerolog.Logger,
) (*OutboundData, bool, error) {
	outboundParams := cctx.GetCurrentOutboundParam()
	if outboundParams == nil {
		return nil, false, errors.New("outboundParams is nil")
	}

	// Check if the CCTX has already been processed
	alreadyIncluded, alreadyConfirmed, err := evmObserver.IsOutboundProcessed(cctx, logger)
	switch {
	case err != nil:
		return nil, true, errors.Wrap(err, "failed to check if outbound is processed")
	case alreadyIncluded || alreadyConfirmed:
		logger.Info().Msgf("CCTX already processed; skipping")
		return nil, true, nil
	}

	// Determine destination chain and address
	to, toChainID, skip := determineDestination(cctx, logger)
	if skip {
		return nil, true, nil
	}

	toChain := chains.GetChainFromChainID(toChainID.Int64())
	if toChain == nil {
		return nil, false, fmt.Errorf("unknown chain: %d", toChainID.Int64())
	}

	nonce := outboundParams.TssNonce

	// Get sendHash
	cctxIndex, err := getCCTXIndex(cctx)
	if err != nil {
		return nil, false, err
	}

	// Determine gas fees
	gas, err := determineGas(cctx, logger)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to determine gas fees")
	}

	// In case there is a pending transaction, make sure this keysign is a transaction replacement
	if tx := evmObserver.GetPendingTx(nonce); tx != nil {
		newFeeIsLessThanPending := gas.MaxFeePerUnit.Cmp(tx.GasPrice()) <= 0

		evt := logger.Info().
			Str("cctx.pending_tx_hash", tx.Hash().Hex()).
			Uint64("cctx.pending_tx_nonce", nonce)

		if newFeeIsLessThanPending {
			evt.Msg("Please wait for pending outbound to be included in the block")
			return nil, true, nil
		}

		evt.
			Uint64("cctx.max_gas_fee", gas.MaxFeePerUnit.Uint64()).
			Msg("Replacing pending outbound transaction with higher gas fees")
	}

	// Base64 decode message
	var message []byte
	if cctx.InboundParams.CoinType != coin.CoinType_Cmd {
		msg, errDecode := base64.StdEncoding.DecodeString(cctx.RelayedMessage)
		if errDecode != nil {
			logger.Err(err).Str("cctx.relayed_message", cctx.RelayedMessage).Msg("Unable to decode relayed message")
		} else {
			message = msg
		}
	}

	logger.Info().
		Str("cctx.to_chain_name", toChain.GetChainName().String()).
		Int64("cctx.to_chain_id", toChain.ChainId).
		Str("cctx.to_recipient", to.Hex()).
		Uint64("cctx.nonce", nonce).
		Uint64("cctx.finalized_zeta_height", cctx.InboundParams.FinalizedZetaHeight).
		Msg("Constructed OutboundData")

	return &OutboundData{
		outboundParams: outboundParams,

		srcChainID: big.NewInt(cctx.InboundParams.SenderChainId),
		sender:     ethcommon.HexToAddress(cctx.InboundParams.Sender),

		toChainID: toChainID,
		to:        to,

		asset:  ethcommon.HexToAddress(cctx.InboundParams.Asset),
		amount: outboundParams.Amount.BigInt(),

		gas:    gas,
		nonce:  outboundParams.TssNonce,
		height: height,

		message: message,

		cctxIndex: cctxIndex,
	}, false, nil
}

func getCCTXIndex(cctx *types.CrossChainTx) ([32]byte, error) {
	cctxIndexSlice, err := hex.DecodeString(cctx.Index[2:]) // remove the leading 0x
	if err != nil || len(cctxIndexSlice) != 32 {
		return [32]byte{}, errors.Wrapf(err, "unable to decode cctx index %s", cctx.Index)
	}

	var cctxIndex [32]byte
	copy(cctxIndex[:32], cctxIndexSlice[:32])

	return cctxIndex, nil
}

// determineDestination picks the destination address and Chain ID based on the status of the cross chain tx.
// returns true if transaction should be skipped.
func determineDestination(cctx *types.CrossChainTx, logger zerolog.Logger) (ethcommon.Address, *big.Int, bool) {
	switch cctx.CctxStatus.Status {
	case types.CctxStatus_PendingRevert:
		to := ethcommon.HexToAddress(cctx.InboundParams.Sender)
		chainID := big.NewInt(cctx.InboundParams.SenderChainId)

		logger.Info().
			Str("cctx.index", cctx.Index).
			Int64("cctx.chain_id", chainID.Int64()).
			Msgf("Abort: reverting inbound")

		return to, chainID, false
	case types.CctxStatus_PendingOutbound:
		to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		chainID := big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)

		return to, chainID, false
	}

	logger.Info().
		Str("cctx.index", cctx.Index).
		Str("cctx.status", cctx.CctxStatus.String()).
		Msgf("CCTX doesn't need to be processed")

	return ethcommon.Address{}, nil, true
}

func determineGas(cctx *types.CrossChainTx, logger zerolog.Logger) (gas, error) {
	var (
		outboundParams = cctx.GetCurrentOutboundParam()
		limit          = outboundParams.GasLimit
	)

	switch {
	case limit < MinGasLimit:
		limit = MinGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", outboundParams.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too low. Setting to the minimum (%d)", MinGasLimit)
	case limit > MaxGasLimit:
		limit = MaxGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", outboundParams.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too high; Setting to the maximum (%d)", MaxGasLimit)
	}

	maxFee, ok := new(big.Int).SetString(outboundParams.GasPrice, 10)
	if !ok {
		return gas{}, errors.New("unable to parse gasPrice from " + outboundParams.GasPrice)
	}

	// TODO RELY ONLY ON gas{} data.
	// use dynamic gas price for ethereum chains.
	// The code below is a fix for https://github.com/zeta-chain/node/issues/1085
	// doesn't close directly the issue because we should determine if we want to keep using SuggestGasPrice if no GasPrice
	// we should possibly remove it completely and return an error if no GasPrice is provided because it means no fee is processed on ZetaChain
	//specified, ok := new(big.Int).SetString(cctx.GetCurrentOutboundParam().GasPrice, 10)
	//if !ok {
	//	if chains.IsEthereumChain(chain.ChainId) {
	//		suggested, err := client.SuggestGasPrice(context.Background())
	//		if err != nil {
	//			return errors.Join(err, fmt.Errorf("cannot get gas price from chain %s ", chain))
	//		}
	//		txData.gasPrice = roundUpToNearestGwei(suggested)
	//	} else {
	//		return fmt.Errorf("cannot convert gas price  %s ", cctx.GetCurrentOutboundParam().GasPrice)
	//	}
	//} else {
	//	txData.gasPrice = specified
	//}

	return gas{
		Limit:              limit,
		MaxFeePerUnit:      maxFee,
		PriorityFeePerUnit: big.NewInt(0), // todo!
	}, nil
}
