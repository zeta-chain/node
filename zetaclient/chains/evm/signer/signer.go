// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"fmt"
	"math/big"
	"runtime/debug"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	ethcommon "github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/retry"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = time.Second * 6

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 5

	// broadcastTimeout is the timeout for broadcasting a transaction
	// we should allow enough time for the tx submission and avoid fast timeout
	broadcastTimeout = time.Second * 15
)

var (
	// zeroValue is for outbounds that carry no ETH (gas token) value
	zeroValue = big.NewInt(0)

	// ErrWaitForSignature is the error returned when waiting for transaction signature
	ErrWaitForSignature = errors.New("waiting for transaction signature")
)

type EVMClient interface {
	NonceAt(_ context.Context, account ethcommon.Address, blockNumber *big.Int) (uint64, error)

	IsTxConfirmed(_ context.Context, txHash string, confirmations uint64) (bool, error)

	Signer() eth.Signer

	// This is a mutating function that does not get called when zetaclient is in dry-mode.
	SendTransaction(context.Context, *eth.Transaction) error
}

// Signer deals with the signing EVM transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// evmClient is the EVM RPC client used to interact with the EVM chain
	evmClient EVMClient

	// zetaConnectorAddress is the address of the ZetaConnector contract
	zetaConnectorAddress ethcommon.Address

	// er20CustodyAddress is the address of the ERC20Custody contract
	er20CustodyAddress ethcommon.Address

	// gatewayAddress is the address of the Gateway contract
	gatewayAddress ethcommon.Address
}

// New Signer constructor
func New(
	baseSigner *base.Signer,
	evmClient EVMClient,
	zetaConnectorAddress ethcommon.Address,
	erc20CustodyAddress ethcommon.Address,
	gatewayAddress ethcommon.Address,
) (*Signer, error) {
	return &Signer{
		Signer:               baseSigner,
		evmClient:            evmClient,
		zetaConnectorAddress: zetaConnectorAddress,
		er20CustodyAddress:   erc20CustodyAddress,
		gatewayAddress:       gatewayAddress,
	}, nil
}

// SetZetaConnectorAddress sets the zeta connector address
func (signer *Signer) SetZetaConnectorAddress(addr ethcommon.Address) {
	// noop
	if (addr == ethcommon.Address{}) || signer.zetaConnectorAddress == addr {
		return
	}

	signer.Logger().Std.Info().
		Stringer("signer_old_zeta_connector_address", signer.zetaConnectorAddress).
		Stringer("signer_new_zeta_connector_address", addr).
		Msg("updated zeta connector address")

	signer.Lock()
	signer.zetaConnectorAddress = addr
	signer.Unlock()
}

// SetERC20CustodyAddress sets the erc20 custody address
func (signer *Signer) SetERC20CustodyAddress(addr ethcommon.Address) {
	// noop
	if (addr == ethcommon.Address{}) || signer.er20CustodyAddress == addr {
		return
	}

	signer.Logger().Std.Info().
		Stringer("signer_old_erc20_custody_address", signer.er20CustodyAddress).
		Stringer("signer_new_erc20_custody_address", addr).
		Msg("updated ERC-20 custody address")

	signer.Lock()
	signer.er20CustodyAddress = addr
	signer.Unlock()
}

// SetGatewayAddress sets the gateway address
func (signer *Signer) SetGatewayAddress(addrRaw string) {
	addr := ethcommon.HexToAddress(addrRaw)

	// noop
	if (addr == ethcommon.Address{}) || signer.gatewayAddress == addr {
		return
	}

	signer.Logger().Std.Info().
		Stringer("signer_old_gateway_address", signer.gatewayAddress).
		Stringer("signer_new_gateway_address", addr).
		Msg("Updated gateway address")

	signer.Lock()
	signer.gatewayAddress = addr
	signer.Unlock()
}

// GetZetaConnectorAddress returns the zeta connector address
func (signer *Signer) GetZetaConnectorAddress() ethcommon.Address {
	return signer.zetaConnectorAddress
}

// GetERC20CustodyAddress returns the erc20 custody address
func (signer *Signer) GetERC20CustodyAddress() ethcommon.Address {
	return signer.er20CustodyAddress
}

// GetGatewayAddress returns the gateway address
func (signer *Signer) GetGatewayAddress() string {
	return signer.gatewayAddress.String()
}

// NextTSSNonce returns the next nonce of the TSS account
func (signer *Signer) NextTSSNonce(ctx context.Context) (uint64, error) {
	nextNonce, err := signer.evmClient.NonceAt(ctx, signer.TSS().PubKey().AddressEVM(), nil)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get TSS account nonce")
	}

	// update next TSS nonce metrics
	signer.SetNextTSSNonce(nextNonce)

	return nextNonce, nil
}

// Sign given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *Signer) Sign(
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
) (*eth.Transaction, []byte, []byte, error) {
	chainID := big.NewInt(signer.Chain().ChainId)
	tx, err := newTx(chainID, data, to, amount, gas, nonce)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to create new tx")
	}

	hashBytes := signer.evmClient.Signer().Hash(tx).Bytes()

	// get cached signature if available, otherwise add digest and wait for keysign
	sig, found := signer.GetSignatureOrAddDigest(nonce, hashBytes)
	if !found {
		return nil, nil, nil, ErrWaitForSignature
	}

	_, err = crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to derive pub key from signature")
	}

	signedTX, err := tx.WithSignature(signer.evmClient.Signer(), sig[:])
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to set tx signature")
	}

	return signedTX, sig[:], hashBytes[:], nil
}

func newTx(
	_ *big.Int,
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
) (*eth.Transaction, error) {
	if err := gas.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid gas parameters")
	}

	// https://github.com/zeta-chain/node/issues/3221
	//if gas.isLegacy() {
	return eth.NewTx(&eth.LegacyTx{
		To:       &to,
		Value:    amount,
		Data:     data,
		GasPrice: gas.Price,
		Gas:      gas.Limit,
		Nonce:    nonce,
	}), nil
	//}
	//
	//return ethtypes.NewTx(&ethtypes.DynamicFeeTx{
	//	ChainID:   chainID,
	//	To:        &to,
	//	Value:     amount,
	//	Data:      data,
	//	GasFeeCap: gas.Price,
	//	GasTipCap: gas.PriorityFee,
	//	Gas:       gas.Limit,
	//	Nonce:     nonce,
	//}), nil
}

func (signer *Signer) broadcast(ctx context.Context, tx *eth.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, broadcastTimeout)
	defer cancel()

	return signer.evmClient.SendTransaction(ctx, tx)
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign an evm transaction using the TSS signer.
// It will then broadcast the signed transaction to the outbound chain.
// TODO(revamp): simplify function
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	zetaRepo *zrepo.ZetaRepo,
	_ uint64,
) {
	outboundID := base.OutboundIDFromCCTX(cctx)
	signer.MarkOutbound(outboundID, true)

	// end outbound process on panic
	defer func() {
		signer.MarkOutbound(outboundID, false)
		if r := recover(); r != nil {
			signer.Logger().
				Std.Error().
				Str(logs.FieldCctxIndex, cctx.Index).
				Interface("panic", r).
				Str("stack_trace", string(debug.Stack())).
				Msg("caught panic error")
		}
	}()

	// prepare logger and a few local variables
	var (
		params = cctx.GetCurrentOutboundParam()
		myID   = zetaRepo.GetOperatorAddress()
		logger = signer.Logger().Std.With().
			Int64(logs.FieldChain, signer.Chain().ChainId).
			Uint64(logs.FieldNonce, params.TssNonce).
			Str(logs.FieldCctxIndex, cctx.Index).
			Str("cctx_receiver", params.Receiver).
			Stringer("cctx_amount", params.Amount).
			Str("signer", myID).
			Logger()
	)
	logger.Debug().Msg("TryProcessOutbound")

	// retrieve app context
	app, err := zctx.FromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error getting app context")
		return
	}

	// Setup Transaction input
	txData, skipTx, err := NewOutboundData(ctx, cctx, logger)
	if err != nil {
		logger.Err(err).Msg("error setting up transaction input fields")
		return
	}

	if skipTx {
		return
	}

	toChain, err := app.GetChain(txData.toChainID.Int64())
	switch {
	case err != nil:
		logger.Error().
			Err(err).
			Int64("to_chain_id", txData.toChainID.Int64()).
			Msg("error getting toChain")
		return
	case toChain.IsZeta():
		// should not happen
		logger.Error().
			Int64("to_chain_id", toChain.ID()).
			Msg("unable to TryProcessOutbound when toChain is zetaChain")
		return
	}

	logger = logger.With().Uint64("gas_price", txData.gas.Price.Uint64()).Logger()

	if signer.ClientMode.IsDryMode() {
		logger.Info().Stringer(logs.FieldMode, mode.DryMode).Msg("skipping outbound processing")
		return
	}

	// sign outbound
	tx, err := signer.SignOutboundFromCCTX(
		logger,
		cctx,
		txData,
		zetaRepo,
		toChain,
	)
	if errors.Is(err, ErrWaitForSignature) {
		return
	} else if err != nil {
		logger.Err(err).Msg("error signing outbound")
		return
	}

	// attach tx hash to logger and print log
	logger = logger.With().Str(logs.FieldTx, tx.Hash().Hex()).Logger()

	// Broadcast Signed Tx
	signer.BroadcastOutbound(ctx, tx, cctx, logger, zetaRepo, txData)
}

// SignOutboundFromCCTX signs an outbound transaction from a given cctx
// TODO: simplify logic with all if else
// https://github.com/zeta-chain/node/issues/2050
func (signer *Signer) SignOutboundFromCCTX(
	logger zerolog.Logger,
	cctx *crosschaintypes.CrossChainTx,
	outboundData *OutboundData,
	zetaRepo *zrepo.ZetaRepo,
	_ zctx.Chain,
) (*eth.Transaction, error) {
	switch {
	case !signer.PassesCompliance(cctx):
		// restricted cctx
		return signer.SignCancel(outboundData)
	case cctx.InboundParams.CoinType == coin.CoinType_Cmd:
		// admin command
		to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		if to == (ethcommon.Address{}) {
			return nil, fmt.Errorf("invalid receiver %s", cctx.GetCurrentOutboundParam().Receiver)
		}
		msg := strings.Split(cctx.RelayedMessage, ":")
		if len(msg) != 2 {
			return nil, fmt.Errorf("invalid message %s", msg)
		}
		// cmd field is used to determine whether to execute ERC20 whitelist or migrate TSS funds given that the coin type
		// from the cctx is coin.CoinType_Cmd
		cmd := msg[0]
		// params field is used to pass input parameters for command requests, currently it is used to pass the ERC20
		// contract address when a whitelist command is requested
		params := msg[1]
		return signer.SignAdminTx(outboundData, cmd, params)
	case cctx.ProtocolContractVersion == crosschaintypes.ProtocolContractVersion_V2:
		// call sign outbound from cctx for v2 protocol contracts
		return signer.SignOutboundFromCCTXV2(cctx, outboundData)
	case IsPendingOutboundFromZetaChain(cctx, zetaRepo):
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Gas:
			logger.Info().Msg("calling SignGasWithdraw")
			return signer.SignGasWithdraw(outboundData)
		case coin.CoinType_ERC20:
			logger.Info().Msg("calling SignERC20Withdraw")
			return signer.SignERC20Withdraw(outboundData)
		case coin.CoinType_Zeta:
			logger.Info().Msg("calling SignConnectorOnReceive")
			return signer.SignConnectorOnReceive(outboundData)
		}
	case cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingRevert && cctx.OutboundParams[0].ReceiverChainId == zetaRepo.ZetaChain().ChainId:
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Zeta:
			logger.Info().Msg("calling SignConnectorOnRevert")
			outboundData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
			outboundData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
			return signer.SignConnectorOnRevert(outboundData)
		case coin.CoinType_Gas:
			logger.Info().Msg("calling SignGasWithdraw")
			return signer.SignGasWithdraw(outboundData)
		case coin.CoinType_ERC20:
			logger.Info().Msg("calling SignERC20Withdraw")
			return signer.SignERC20Withdraw(outboundData)
		}
	case cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingRevert:
		logger.Info().Msg("calling SignConnectorOnRevert")
		outboundData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
		outboundData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
		return signer.SignConnectorOnRevert(outboundData)
	case cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingOutbound:
		logger.Info().Msg("calling SignConnectorOnReceive")
		return signer.SignConnectorOnReceive(outboundData)
	}

	return nil, fmt.Errorf("unknown signing method for cctx %s", cctx.String())
}

// BroadcastOutbound signed transaction through evm rpc client
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *eth.Transaction,
	cctx *crosschaintypes.CrossChainTx,
	logger zerolog.Logger,
	zetaRepo *zrepo.ZetaRepo,
	txData *OutboundData,
) {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		logger.Err(err).Msg("error getting app context")
		return
	}

	toChain, err := app.GetChain(txData.toChainID.Int64())
	switch {
	case err != nil:
		logger.Error().
			Err(err).
			Int64("to_chain_id", txData.toChainID.Int64()).
			Msg("error getting toChain")
		return
	case toChain.IsZeta():
		// should not happen
		logger.Error().
			Int64("to_chain_id", toChain.ID()).
			Msg("unable to broadcast when toChain is zetaChain")
		return
	case tx == nil:
		logger.Warn().Str(logs.FieldCctxIndex, cctx.Index).Msg("no outbound tx to broadcast")
		return
	}

	var (
		outboundHash = tx.Hash().Hex()
		nonce        = cctx.GetCurrentOutboundParam().TssNonce
	)

	// define broadcast function
	broadcast := func() error {
		// get latest TSS account nonce
		latestNonce, err := signer.evmClient.NonceAt(ctx, signer.TSS().PubKey().AddressEVM(), nil)
		if err != nil {
			return errors.Wrap(err, "unable to get latest TSS account nonce")
		}

		// if TSS nonce is higher than CCTX nonce, there is no need to broadcast
		// this avoids foreseeable "nonce too low" error and unnecessary tracker report
		// Note: the latest finalized nonce is used here, not the pending nonce, making it possible
		//       to replace pending txs
		if latestNonce > nonce {
			logger.Info().
				Uint64("latest_nonce", latestNonce).
				Msg("skipped broadcasting tx because CCTX nonce is too low")
			return nil
		}

		// broadcast success, report to tracker
		if err = signer.broadcast(ctx, tx); err == nil {
			signer.reportToOutboundTracker(ctx, zetaRepo, toChain.ID(), nonce, outboundHash, logger)
			return nil
		}

		// handle different broadcast errors
		retry, report := zetacore.HandleBroadcastError(err, nonce, toChain.ID(), outboundHash)
		if report {
			signer.reportToOutboundTracker(ctx, zetaRepo, toChain.ID(), nonce, outboundHash, logger)
			return nil
		}
		if retry {
			return errors.Wrap(err, "unable to broadcast tx, retrying")
		}

		// no re-broadcast, no report, stop retry
		// e.g. "replacement transaction underpriced"
		return nil
	}

	// broadcast transaction with backoff to tolerate RPC error
	bo := backoff.NewConstantBackOff(broadcastBackoff)
	boWithMaxRetries := backoff.WithMaxRetries(bo, broadcastRetries)
	if err := retry.DoWithBackoff(broadcast, boWithMaxRetries); err != nil {
		logger.Error().Err(err).Msg("unable to broadcast EVM outbound")
	}

	logger.Info().Msg("broadcasted EVM outbound")
}

// IsPendingOutboundFromZetaChain checks if the sender chain is ZetaChain and if status is PendingOutbound
// TODO(revamp): move to another package more general for cctx functions
func IsPendingOutboundFromZetaChain(
	cctx *crosschaintypes.CrossChainTx,
	zetaRepo *zrepo.ZetaRepo,
) bool {
	return cctx.InboundParams.SenderChainId == zetaRepo.ZetaChain().ChainId &&
		cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingOutbound
}

// ErrorMsg returns a error message for SignConnectorOnReceive failure with cctx data
func ErrorMsg(cctx *crosschaintypes.CrossChainTx) string {
	return fmt.Sprintf(
		"signer SignConnectorOnReceive error: nonce %d chain %d",
		cctx.GetCurrentOutboundParam().TssNonce,
		cctx.GetCurrentOutboundParam().ReceiverChainId,
	)
}
