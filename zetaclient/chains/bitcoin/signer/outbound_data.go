package signer

import (
	"fmt"
	"math"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/compliance"
)

// OutboundData is a data structure containing necessary data to construct a BTC outbound transaction
type OutboundData struct {
	// to is the recipient address
	to btcutil.Address

	// amount is the amount in BTC
	amount float64

	// amountSats is the amount in satoshis
	amountSats int64

	// feeRate is the fee rate in satoshis/vByte
	feeRate int64

	// txSize is the average size of a BTC outbound transaction
	// user is charged (in ZRC20 contract) at a static txSize on each withdrawal
	txSize int64

	// nonce is the nonce of the outbound
	nonce uint64

	// height is the ZetaChain block height
	height uint64

	// cancelTx is a flag to indicate if this outbound should be cancelled
	cancelTx bool
}

// NewOutboundData creates OutboundData from the given CCTX.
func NewOutboundData(
	cctx *types.CrossChainTx,
	height uint64,
	minRelayFee float64,
	logger, loggerCompliance zerolog.Logger,
) (*OutboundData, error) {
	if cctx == nil {
		return nil, errors.New("cctx is nil")
	}
	params := cctx.GetCurrentOutboundParam()

	// support gas token only for Bitcoin outbound
	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return nil, fmt.Errorf("invalid coin type %s", cctx.InboundParams.CoinType.String())
	}

	// initial fee rate
	feeRate, err := strconv.ParseInt(params.GasPrice, 10, 64)
	if err != nil || feeRate <= 0 {
		return nil, fmt.Errorf("invalid fee rate %s", params.GasPrice)
	}

	// use current gas rate if fed by zetacore
	newRate, err := strconv.ParseInt(params.GasPriorityFee, 10, 64)
	if err == nil && newRate > 0 && newRate != feeRate {
		logger.Info().Msgf("use new fee rate %d sat/vB instead of %d sat/vB", newRate, feeRate)
		feeRate = newRate
	}

	// check receiver address
	to, err := chains.DecodeBtcAddress(params.Receiver, params.ReceiverChainId)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}
	if !chains.IsBtcAddressSupported(to) {
		return nil, fmt.Errorf("unsupported receiver address %s", to.EncodeAddress())
	}

	// amount in BTC and satoshis
	amount := float64(params.Amount.Uint64()) / 1e8
	amountSats := params.Amount.BigInt().Int64()

	// check gas limit
	if params.CallOptions == nil {
		// never happens, 'GetCurrentOutboundParam' will create it
		return nil, errors.New("call options is nil")
	}
	if params.CallOptions.GasLimit > math.MaxInt64 {
		return nil, fmt.Errorf("invalid gas limit %d", params.CallOptions.GasLimit)
	}

	// add minimum relay fee (1000 satoshis/KB by default) to gasPrice to avoid minRelayTxFee error
	// see: https://github.com/bitcoin/bitcoin/blob/master/src/policy/policy.h#L35
	satPerByte := common.FeeRateToSatPerByte(minRelayFee)
	feeRate += satPerByte

	// compliance check
	restrictedCCTX := compliance.IsCctxRestricted(cctx)
	if restrictedCCTX {
		compliance.PrintComplianceLog(logger, loggerCompliance,
			true, params.ReceiverChainId, cctx.Index, cctx.InboundParams.Sender, params.Receiver, "BTC")
	}

	// check dust amount
	dustAmount := params.Amount.Uint64() < constant.BTCWithdrawalDustAmount
	if dustAmount {
		logger.Warn().Msgf("dust amount %d sats, canceling tx", params.Amount.Uint64())
	}

	// set the amount to 0 when the tx should be cancelled
	cancelTx := restrictedCCTX || dustAmount
	if cancelTx {
		amount = 0.0
		amountSats = 0
	}

	return &OutboundData{
		to:         to,
		amount:     amount,
		amountSats: amountSats,
		feeRate:    feeRate,
		// #nosec G115 checked in range
		txSize:   int64(params.CallOptions.GasLimit),
		nonce:    params.TssNonce,
		height:   height,
		cancelTx: cancelTx,
	}, nil
}
