package signer

import (
	"fmt"
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

	// feeRateLatest is the latest median fee rate in satoshis/vByte
	// this value is fed by the zetacore when it bumps the gas price with gas stability pool
	feeRateLatest int64

	// feeRateBumpped is a flag to indicate if the fee rate in CCTX is bumped by zetacore
	feeRateBumped bool

	// minRelayFee is the minimum relay fee in unit of BTC
	minRelayFee float64

	// height is the ZetaChain block height
	height uint64

	// nonce is the nonce of the outbound
	nonce uint64

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

	// support coin type GAS and CMD only
	if cctx.InboundParams.CoinType != coin.CoinType_Gas && cctx.InboundParams.CoinType != coin.CoinType_Cmd {
		return nil, fmt.Errorf("invalid coin type %s", cctx.InboundParams.CoinType.String())
	}

	// parse fee rate
	feeRate, err := strconv.ParseInt(params.GasPrice, 10, 64)
	if err != nil || feeRate <= 0 {
		return nil, fmt.Errorf("invalid fee rate %s", params.GasPrice)
	}

	// check if zetacore has bumped the fee rate
	// 'GasPriorityFee' is always empty for Bitcoin unless zetacore bumps the fee rate
	var (
		feeRateBumped bool
		feeRateLatest int64
	)
	if params.GasPriorityFee != "" {
		gasPriorityFee, err := strconv.ParseInt(params.GasPriorityFee, 10, 64)
		if err != nil || gasPriorityFee <= 0 {
			return nil, fmt.Errorf("invalid gas priority fee %s", params.GasPriorityFee)
		}

		feeRateBumped = true
		feeRateLatest = gasPriorityFee
		logger.Info().Str("latest_fee_rate", params.GasPriorityFee).Msg("fee rate is bumped by zetacore")
	}

	// to avoid minRelayTxFee error, please do not use the minimum rate (1 sat/vB by default).
	// we simply add additional 1 sat/vB to 'minRate' to avoid tx rejection by Bitcoin core.
	// see: https://github.com/bitcoin/bitcoin/blob/master/src/policy/policy.h#L35
	minRate := common.FeeRateToSatPerByte(minRelayFee)
	if feeRate <= minRate {
		feeRate = minRate + 1
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
	// the float64 'amount' is used later to select UTXOs, precision does not matter
	amount := float64(params.Amount.Uint64()) / 1e8
	amountSats := params.Amount.BigInt().Int64()

	// compliance check
	restrictedCCTX := compliance.IsCctxRestricted(cctx)
	if restrictedCCTX {
		compliance.PrintComplianceLog(logger, loggerCompliance,
			true, params.ReceiverChainId, cctx.Index, cctx.InboundParams.Sender, params.Receiver, "BTC")
	}

	// check dust amount
	dustAmount := amountSats < constant.BTCWithdrawalDustAmount
	if dustAmount {
		logger.Warn().Int64("amount", amountSats).Msg("outbound will be cancelled due to dust amount")
	}

	// set the amount to 0 when the tx should be cancelled
	cancelTx := restrictedCCTX || dustAmount
	if cancelTx {
		amount = 0.0
		amountSats = 0
	}

	return &OutboundData{
		to:            to,
		amount:        amount,
		amountSats:    amountSats,
		feeRate:       feeRate,
		feeRateLatest: feeRateLatest,
		feeRateBumped: feeRateBumped,
		minRelayFee:   minRelayFee,
		height:        height,
		nonce:         params.TssNonce,
		cancelTx:      cancelTx,
	}, nil
}
