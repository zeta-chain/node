package signer

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
	zctx "github.com/zeta-chain/node/zetaclient/context"
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

	gas    Gas
	nonce  uint64
	height uint64

	message []byte

	// cctxIndex field is the inbound message digest that is sent to the destination contract
	cctxIndex [32]byte

	// outboundParams field contains data detailing the receiver chain and outbound transaction
	outboundParams *types.OutboundParams

	// revertOptions field contains data detailing the revert options
	revertOptions types.RevertOptions
}

// NewOutboundData creates OutboundData from the given CCTX.
// returns `bool true` when transaction should be skipped.
func NewOutboundData(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	logger zerolog.Logger,
) (*OutboundData, bool, error) {
	if cctx == nil {
		return nil, false, errors.New("cctx is nil")
	}

	outboundParams := cctx.GetCurrentOutboundParam()
	if err := validateParams(outboundParams); err != nil {
		return nil, false, errors.Wrap(err, "invalid outboundParams")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, false, errors.Wrap(err, "unable to get app from context")
	}

	var (
		to        ethcommon.Address
		toChainID *big.Int
	)

	// in protocol contract v2, receiver is always set in the outbound
	if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V2 {
		to = ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
	} else {
		// recipient + destination chain
		var skip bool
		to, toChainID, skip = getDestination(cctx, logger)
		if skip {
			return nil, true, nil
		}
	}

	// ensure that chain exists in app's context
	if _, err := app.GetChain(toChainID.Int64()); err != nil {
		return nil, false, errors.Wrapf(err, "unable to get chain %d from app context", toChainID.Int64())
	}

	gas, err := gasFromCCTX(cctx, logger)
	if err != nil {
		return nil, false, errors.Wrap(err, "unable to make gas from CCTX")
	}

	cctxIndex, err := getCCTXIndex(cctx)
	if err != nil {
		return nil, false, errors.Wrap(err, "unable to get cctx index")
	}

	// Base64 decode message
	var message []byte
	if cctx.InboundParams.CoinType != coin.CoinType_Cmd {
		// protocol contract v2 uses hex encoding
		if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V2 {
			message, err = hex.DecodeString(cctx.RelayedMessage)
			if err != nil {
				logger.Err(err).Msgf("decode CCTX.Message %s error", cctx.RelayedMessage)
				message = []byte{}
			}
		} else {
			msg, errDecode := base64.StdEncoding.DecodeString(cctx.RelayedMessage)
			if errDecode != nil {
				logger.Err(err).Str("cctx.relayed_message", cctx.RelayedMessage).Msg("Unable to decode relayed message")
			} else {
				message = msg
			}
		}
	}

	return &OutboundData{
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

		outboundParams: outboundParams,

		revertOptions: cctx.RevertOptions,
	}, false, nil
}

func getCCTXIndex(cctx *types.CrossChainTx) ([32]byte, error) {
	// `0x` + `64 chars`. Two chars ranging `00...FF` represent one byte (64 chars = 32 bytes)
	if len(cctx.Index) != (2 + 64) {
		return [32]byte{}, fmt.Errorf("cctx index %q is invalid", cctx.Index)
	}

	// remove the leading `0x`
	cctxIndexSlice, err := hex.DecodeString(cctx.Index[2:])
	if err != nil || len(cctxIndexSlice) != 32 {
		return [32]byte{}, errors.Wrapf(err, "unable to decode cctx index %s", cctx.Index)
	}

	var cctxIndex [32]byte
	copy(cctxIndex[:32], cctxIndexSlice[:32])

	return cctxIndex, nil
}

// getDestination picks the destination address and Chain ID based on the status of the cross chain tx.
// returns true if transaction should be skipped.
func getDestination(cctx *types.CrossChainTx, logger zerolog.Logger) (ethcommon.Address, *big.Int, bool) {
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

func validateParams(params *types.OutboundParams) error {
	if params == nil || params.GasLimit == 0 {
		return errors.New("outboundParams is empty")
	}

	return nil
}
