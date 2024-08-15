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

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/logs"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
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

	// zetaConnectorABI is the ABI of the ZetaConnector contract
	zetaConnectorABI abi.ABI

	// erc20CustodyABI is the ABI of the ERC20Custody contract
	erc20CustodyABI abi.ABI

	// zetaConnectorAddress is the address of the ZetaConnector contract
	zetaConnectorAddress ethcommon.Address

	// er20CustodyAddress is the address of the ERC20Custody contract
	er20CustodyAddress ethcommon.Address
}

// NewSigner creates a new EVM signer
func NewSigner(
	ctx context.Context,
	chain chains.Chain,
	tss interfaces.TSSSigner,
	ts *metrics.TelemetryServer,
	logger base.Logger,
	endpoint string,
	zetaConnectorABI string,
	erc20CustodyABI string,
	zetaConnectorAddress ethcommon.Address,
	erc20CustodyAddress ethcommon.Address,
) (*Signer, error) {
	// create base signer
	baseSigner := base.NewSigner(chain, tss, ts, logger)

	// create EVM client
	client, ethSigner, err := getEVMRPC(ctx, endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create EVM client")
	}

	// prepare ABIs
	connectorABI, err := abi.JSON(strings.NewReader(zetaConnectorABI))
	if err != nil {
		return nil, errors.Wrap(err, "unable to build ZetaConnector ABI")
	}

	custodyABI, err := abi.JSON(strings.NewReader(erc20CustodyABI))
	if err != nil {
		return nil, errors.Wrap(err, "unable to build ERC20Custody ABI")
	}

	return &Signer{
		Signer:               baseSigner,
		client:               client,
		ethSigner:            ethSigner,
		zetaConnectorABI:     connectorABI,
		erc20CustodyABI:      custodyABI,
		zetaConnectorAddress: zetaConnectorAddress,
		er20CustodyAddress:   erc20CustodyAddress,
	}, nil
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
func (signer *Signer) SetGatewayAddress(_ string) {
	// Note: do nothing for now
	// gateway address will be needed in the future contract architecture
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
	// Note: return empty string for now
	// gateway address will be needed in the future contract architecture
	return ""
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

// SignOutbound
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
func (signer *Signer) SignOutbound(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	data, err = signer.zetaConnectorABI.Pack("onReceive",
		txData.sender.Bytes(),
		txData.srcChainID,
		txData.to,
		txData.amount,
		txData.message,
		txData.cctxIndex)
	if err != nil {
		return nil, fmt.Errorf("onReceive pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign onReceive error: %w", err)
	}

	return tx, nil
}

// SignRevertTx
// function onRevert(
// address originSenderAddress,
// uint256 originChainId,
// bytes calldata destinationAddress,
// uint256 destinationChainId,
// uint256 zetaAmount,
// bytes calldata message,
// bytes32 internalSendHash
// ) external override whenNotPaused onlyTssAddress
func (signer *Signer) SignRevertTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := signer.zetaConnectorABI.Pack("onRevert",
		txData.sender,
		txData.srcChainID,
		txData.to.Bytes(),
		txData.toChainID,
		txData.amount,
		txData.message,
		txData.cctxIndex,
	)
	if err != nil {
		return nil, fmt.Errorf("onRevert pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign onRevert error: %w", err)
	}

	return tx, nil
}

// SignCancelTx signs a transaction from TSS address to itself with a zero amount in order to increment the nonce
func (signer *Signer) SignCancelTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	txData.gas.Limit = evm.EthTransferGasLimit
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		signer.TSS().EVMAddress(),
		zeroValue, // zero out the amount to cancel the tx
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("SignCancelTx error: %w", err)
	}

	return tx, nil
}

// SignWithdrawTx signs a withdrawal transaction sent from the TSS address to the destination
func (signer *Signer) SignWithdrawTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	txData.gas.Limit = evm.EthTransferGasLimit
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		txData.to,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("SignWithdrawTx error: %w", err)
	}

	return tx, nil
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign an evm transaction using the TSS signer.
// It will then broadcast the signed transaction to the outbound chain.
// TODO(revamp): simplify function
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
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

	// https://github.com/zeta-chain/node/issues/2050
	var tx *ethtypes.Transaction
	// compliance check goes first
	if compliance.IsCctxRestricted(cctx) {
		compliance.PrintComplianceLog(
			logger,
			signer.Logger().Compliance,
			true,
			signer.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			txData.to.Hex(),
			cctx.GetCurrentOutboundParam().CoinType.String(),
		)

		tx, err = signer.SignCancelTx(ctx, txData) // cancel the tx
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.InboundParams.CoinType == coin.CoinType_Cmd { // admin command
		to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		if to == (ethcommon.Address{}) {
			logger.Error().Msgf("invalid receiver %s", cctx.GetCurrentOutboundParam().Receiver)
			return
		}
		msg := strings.Split(cctx.RelayedMessage, ":")
		if len(msg) != 2 {
			logger.Error().Msgf("invalid message %s", msg)
			return
		}
		// cmd field is used to determine whether to execute ERC20 whitelist or migrate TSS funds given that the coin type
		// from the cctx is coin.CoinType_Cmd
		cmd := msg[0]
		// params field is used to pass input parameters for command requests, currently it is used to pass the ERC20
		// contract address when a whitelist command is requested
		params := msg[1]
		tx, err = signer.SignAdminTx(ctx, txData, cmd, params)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if IsSenderZetaChain(cctx, zetacoreClient) {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignWithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignWithdrawTx(ctx, txData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf(
				"SignERC20WithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignERC20WithdrawTx(ctx, txData)
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignOutbound: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignOutbound(ctx, txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.OutboundParams[0].ReceiverChainId == zetacoreClient.Chain().ChainId {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignRevertTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			txData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
			txData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
			tx, err = signer.SignRevertTx(ctx, txData)
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignWithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignWithdrawTx(ctx, txData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignERC20WithdrawTx(ctx, txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		logger.Info().Msgf(
			"SignRevertTx: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			txData.gas.Price,
		)
		txData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)

		tx, err = signer.SignRevertTx(ctx, txData)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf(
			"SignOutbound: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			txData.gas.Price,
		)
		tx, err = signer.SignOutbound(ctx, txData)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
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

// BroadcastOutbound signed transaction through evm rpc client
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *ethtypes.Transaction,
	cctx *types.CrossChainTx,
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

// SignERC20WithdrawTx
// function withdraw(
// address recipient,
// address asset,
// uint256 amount,
// ) external onlyTssAddress
func (signer *Signer) SignERC20WithdrawTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := signer.erc20CustodyABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdraw error: %w", err)
	}

	return tx, nil
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
	cctx *types.CrossChainTx,
	zetacoreClient interfaces.ZetacoreClient,
) bool {
	return cctx.InboundParams.SenderChainId == zetacoreClient.Chain().ChainId &&
		cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound
}

// ErrorMsg returns a error message for SignOutbound failure with cctx data
func ErrorMsg(cctx *types.CrossChainTx) string {
	return fmt.Sprintf(
		"signer SignOutbound error: nonce %d chain %d",
		cctx.GetCurrentOutboundParam().TssNonce,
		cctx.GetCurrentOutboundParam().ReceiverChainId,
	)
}

// getEVMRPC is a helper function to set up the client and signer, also initializes a mock client for unit tests
func getEVMRPC(ctx context.Context, endpoint string) (interfaces.EVMRPCClient, ethtypes.Signer, error) {
	if endpoint == mocks.EVMRPCEnabled {
		chainID := big.NewInt(chains.BscMainnet.ChainId)
		ethSigner := ethtypes.NewLondonSigner(chainID)
		client := &mocks.MockEvmClient{}
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
