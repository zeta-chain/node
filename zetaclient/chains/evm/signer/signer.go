// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/compliance"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = 1000 * time.Millisecond

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 5
)

var (
	_ interfaces.ChainSigner = (*Signer)(nil)

	// zeroValue is for outbounds that carry no ETH (gas token) value
	zeroValue = big.NewInt(0)
)

// Signer deals with the signing EVM transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the EVM RPC client to interact with the EVM chain
	client interfaces.EVMRPCClient

	// ethSigner encapsulates EVM transaction signature handling
	ethSigner ethtypes.Signer

	// zetaConnectorAddress is the address of the ZetaConnector contract
	zetaConnectorAddress ethcommon.Address

	// er20CustodyAddress is the address of the ERC20Custody contract
	er20CustodyAddress ethcommon.Address

	// gatewayAddress is the address of the Gateway contract
	gatewayAddress ethcommon.Address
}

// NewSigner creates a new EVM signer
func NewSigner(
	ctx context.Context,
	chain chains.Chain,
	tss interfaces.TSSSigner,
	ts *metrics.TelemetryServer,
	logger base.Logger,
	endpoint string,
	zetaConnectorAddress ethcommon.Address,
	erc20CustodyAddress ethcommon.Address,
	gatewayAddress ethcommon.Address,
) (*Signer, error) {
	// create base signer
	baseSigner := base.NewSigner(chain, tss, ts, logger)

	// create EVM client
	client, ethSigner, err := getEVMRPC(ctx, endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create EVM client")
	}

	return &Signer{
		Signer:               baseSigner,
		client:               client,
		ethSigner:            ethSigner,
		zetaConnectorAddress: zetaConnectorAddress,
		er20CustodyAddress:   erc20CustodyAddress,
		gatewayAddress:       gatewayAddress,
	}, nil
}

// WithEvmClient attaches a new client to the signer
func (signer *Signer) WithEvmClient(client interfaces.EVMRPCClient) {
	signer.client = client
}

// SetZetaConnectorAddress sets the zeta connector address
func (signer *Signer) SetZetaConnectorAddress(addr ethcommon.Address) {
	signer.Lock()
	defer signer.Unlock()
	signer.zetaConnectorAddress = addr
}

// SetERC20CustodyAddress sets the erc20 custody address
func (signer *Signer) SetERC20CustodyAddress(addr ethcommon.Address) {
	signer.Lock()
	defer signer.Unlock()
	signer.er20CustodyAddress = addr
}

// SetGatewayAddress sets the gateway address
func (signer *Signer) SetGatewayAddress(addr string) {
	signer.Lock()
	defer signer.Unlock()
	signer.gatewayAddress = ethcommon.HexToAddress(addr)
}

// GetZetaConnectorAddress returns the zeta connector address
func (signer *Signer) GetZetaConnectorAddress() ethcommon.Address {
	signer.Lock()
	defer signer.Unlock()
	return signer.zetaConnectorAddress
}

// GetERC20CustodyAddress returns the erc20 custody address
func (signer *Signer) GetERC20CustodyAddress() ethcommon.Address {
	signer.Lock()
	defer signer.Unlock()
	return signer.er20CustodyAddress
}

// GetGatewayAddress returns the gateway address
func (signer *Signer) GetGatewayAddress() string {
	signer.Lock()
	defer signer.Unlock()
	return signer.gatewayAddress.String()
}

// Sign given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *Signer) Sign(
	ctx context.Context,
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
	height uint64,
) (*ethtypes.Transaction, []byte, []byte, error) {
	signer.Logger().Std.Debug().
		Str("tss_pub_key", signer.TSS().EVMAddress().String()).
		Msg("Signing evm transaction")

	chainID := big.NewInt(signer.Chain().ChainId)
	tx, err := newTx(chainID, data, to, amount, gas, nonce)
	if err != nil {
		return nil, nil, nil, err
	}

	hashBytes := signer.ethSigner.Hash(tx).Bytes()

	sig, err := signer.TSS().Sign(ctx, hashBytes, height, nonce, signer.Chain().ChainId, "")
	if err != nil {
		return nil, nil, nil, err
	}

	log.Debug().Msgf("Sign: Signature: %s", hex.EncodeToString(sig[:]))
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.Logger().Std.Error().Err(err).Msgf("SigToPub error")
	}

	addr := crypto.PubkeyToAddress(*pubk)
	signer.Logger().Std.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, nil, nil, err
	}

	return signedTX, sig[:], hashBytes[:], nil
}

func newTx(
	chainID *big.Int,
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
) (*ethtypes.Transaction, error) {
	if err := gas.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid gas parameters")
	}

	if gas.isLegacy() {
		return ethtypes.NewTx(&ethtypes.LegacyTx{
			To:       &to,
			Value:    amount,
			Data:     data,
			GasPrice: gas.Price,
			Gas:      gas.Limit,
			Nonce:    nonce,
		}), nil
	}

	return ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   chainID,
		To:        &to,
		Value:     amount,
		Data:      data,
		GasFeeCap: gas.Price,
		GasTipCap: gas.PriorityFee,
		Gas:       gas.Limit,
		Nonce:     nonce,
	}), nil
}

func (signer *Signer) broadcast(ctx context.Context, tx *ethtypes.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	return signer.client.SendTransaction(ctx, tx)
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign an evm transaction using the TSS signer.
// It will then broadcast the signed transaction to the outbound chain.
// TODO(revamp): simplify function
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	outboundProc *outboundprocessor.Processor,
	outboundID string,
	_ interfaces.ChainObserver,
	zetacoreClient interfaces.ZetacoreClient,
	height uint64,
) {
	// end outbound process on panic
	defer func() {
		outboundProc.EndTryProcess(outboundID)
		if r := recover(); r != nil {
			signer.Logger().Std.Error().Msgf("TryProcessOutbound: %s, caught panic error: %v", cctx.Index, r)
		}
	}()

	// prepare logger and a few local variables
	var (
		params = cctx.GetCurrentOutboundParam()
		myID   = zetacoreClient.GetKeys().GetOperatorAddress()
		logger = signer.Logger().Std.With().
			Str(logs.FieldMethod, "TryProcessOutbound").
			Int64(logs.FieldChain, signer.Chain().ChainId).
			Uint64(logs.FieldNonce, params.TssNonce).
			Str(logs.FieldCctx, cctx.Index).
			Str("cctx.receiver", params.Receiver).
			Str("cctx.amount", params.Amount.String()).
			Logger()
	)
	logger.Info().Msgf("TryProcessOutbound")

	// retrieve app context
	app, err := zctx.FromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error getting app context")
		return
	}

	// Setup Transaction input
	txData, skipTx, err := NewOutboundData(ctx, cctx, height, logger)
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
		logger.Error().Err(err).Msgf("error getting toChain %d", txData.toChainID.Int64())
		return
	case toChain.IsZeta():
		// should not happen
		logger.Error().Msgf("unable to TryProcessOutbound when toChain is zetaChain (%d)", toChain.ID())
		return
	}

	// sign outbound
	tx, err := signer.SignOutboundFromCCTX(
		ctx,
		logger,
		cctx,
		txData,
		zetacoreClient,
		toChain,
	)
	if err != nil {
		logger.Err(err).Msg("error signing outbound")
		return
	}

	logger.Info().Msgf(
		"Key-sign success: %d => %d, nonce %d",
		cctx.InboundParams.SenderChainId,
		toChain.ID(),
		cctx.GetCurrentOutboundParam().TssNonce,
	)

	// Broadcast Signed Tx
	signer.BroadcastOutbound(ctx, tx, cctx, logger, myID, zetacoreClient, txData)
}

// SignOutboundFromCCTX signs an outbound transaction from a given cctx
// TODO: simplify logic with all if else
// https://github.com/zeta-chain/node/issues/2050
func (signer *Signer) SignOutboundFromCCTX(
	ctx context.Context,
	logger zerolog.Logger,
	cctx *crosschaintypes.CrossChainTx,
	outboundData *OutboundData,
	zetacoreClient interfaces.ZetacoreClient,
	toChain zctx.Chain,
) (*ethtypes.Transaction, error) {
	if compliance.IsCctxRestricted(cctx) {
		// restricted cctx
		compliance.PrintComplianceLog(
			logger,
			signer.Logger().Compliance,
			true,
			signer.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			outboundData.to.Hex(),
			cctx.GetCurrentOutboundParam().CoinType.String(),
		)

		return signer.SignCancel(ctx, outboundData)
	} else if cctx.InboundParams.CoinType == coin.CoinType_Cmd {
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
		return signer.SignAdminTx(ctx, outboundData, cmd, params)
	} else if cctx.ProtocolContractVersion == crosschaintypes.ProtocolContractVersion_V2 {
		// call sign outbound from cctx for v2 protocol contracts
		return signer.SignOutboundFromCCTXV2(ctx, cctx, outboundData)
	} else if IsSenderZetaChain(cctx, zetacoreClient) {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignGasWithdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignGasWithdraw(ctx, outboundData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf(
				"SignERC20Withdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignERC20Withdraw(ctx, outboundData)
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignConnectorOnReceive: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignConnectorOnReceive(ctx, outboundData)
		}
	} else if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingRevert && cctx.OutboundParams[0].ReceiverChainId == zetacoreClient.Chain().ChainId {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignConnectorOnRevert: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			outboundData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
			outboundData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
			return signer.SignConnectorOnRevert(ctx, outboundData)
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignGasWithdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignGasWithdraw(ctx, outboundData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf("SignERC20Withdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignERC20Withdraw(ctx, outboundData)
		}
	} else if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingRevert {
		logger.Info().Msgf(
			"SignConnectorOnRevert: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			outboundData.gas.Price,
		)
		outboundData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
		outboundData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
		return signer.SignConnectorOnRevert(ctx, outboundData)
	} else if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingOutbound {
		logger.Info().Msgf(
			"SignConnectorOnReceive: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			outboundData.gas.Price,
		)
		return signer.SignConnectorOnReceive(ctx, outboundData)
	}

	return nil, fmt.Errorf("SignOutboundFromCCTX: can't determine how to sign outbound from cctx %s", cctx.String())
}

// BroadcastOutbound signed transaction through evm rpc client
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *ethtypes.Transaction,
	cctx *crosschaintypes.CrossChainTx,
	logger zerolog.Logger,
	myID sdk.AccAddress,
	zetacoreClient interfaces.ZetacoreClient,
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
		logger.Error().Err(err).Msgf("error getting toChain %d", txData.toChainID.Int64())
		return
	case toChain.IsZeta():
		// should not happen
		logger.Error().Msgf("unable to broadcast when toChain is zetaChain (%d)", toChain.ID())
		return
	case tx == nil:
		logger.Warn().Msgf("BroadcastOutbound: no tx to broadcast %s", cctx.Index)
		return
	}

	// broadcast transaction
	outboundHash := tx.Hash().Hex()

	// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s) in case of RPC error
	backOff := broadcastBackoff
	for i := 0; i < broadcastRetries; i++ {
		time.Sleep(backOff)
		err := signer.broadcast(ctx, tx)
		if err != nil {
			log.Warn().
				Err(err).
				Msgf("BroadcastOutbound: error broadcasting tx %s on chain %d nonce %d retry %d signer %s",
					outboundHash, toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce, i, myID)
			retry, report := zetacore.HandleBroadcastError(
				err,
				strconv.FormatUint(cctx.GetCurrentOutboundParam().TssNonce, 10),
				fmt.Sprintf("%d", toChain.ID()),
				outboundHash,
			)
			if report {
				signer.reportToOutboundTracker(ctx, zetacoreClient, toChain.ID(), tx.Nonce(), outboundHash, logger)
			}
			if !retry {
				break
			}
			backOff *= 2
			continue
		}
		logger.Info().Msgf("BroadcastOutbound: broadcasted tx %s on chain %d nonce %d signer %s",
			outboundHash, toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce, myID)
		signer.reportToOutboundTracker(ctx, zetacoreClient, toChain.ID(), tx.Nonce(), outboundHash, logger)
		break // successful broadcast; no need to retry
	}
}

// EvmClient returns the EVM RPC client
func (signer *Signer) EvmClient() interfaces.EVMRPCClient {
	return signer.client
}

// EvmSigner returns the EVM signer object for the signer
func (signer *Signer) EvmSigner() ethtypes.Signer {
	// TODO(revamp): rename field into evmSigner
	return signer.ethSigner
}

// IsSenderZetaChain checks if the sender chain is ZetaChain
// TODO(revamp): move to another package more general for cctx functions
func IsSenderZetaChain(
	cctx *crosschaintypes.CrossChainTx,
	zetacoreClient interfaces.ZetacoreClient,
) bool {
	return cctx.InboundParams.SenderChainId == zetacoreClient.Chain().ChainId &&
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

// getEVMRPC is a helper function to set up the client and signer, also initializes a mock client for unit tests
func getEVMRPC(ctx context.Context, endpoint string) (interfaces.EVMRPCClient, ethtypes.Signer, error) {
	if endpoint == testutils.MockEVMRPCEndpoint {
		chainID := big.NewInt(chains.BscMainnet.ChainId)
		ethSigner := ethtypes.NewLondonSigner(chainID)
		client := &mocks.EVMRPCClient{}
		return client, ethSigner, nil
	}
	httpClient, err := metrics.GetInstrumentedHTTPClient(endpoint)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get instrumented HTTP client")
	}

	rpcClient, err := ethrpc.DialHTTPWithClient(endpoint, httpClient)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to dial EVM client (endpoint %q)", endpoint)
	}
	client := ethclient.NewClient(rpcClient)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get chain ID")
	}

	ethSigner := ethtypes.LatestSignerForChainID(chainID)

	return client, ethSigner, nil
}
