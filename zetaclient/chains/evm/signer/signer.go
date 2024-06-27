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
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/constant"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	clientcontext "github.com/zeta-chain/zetacore/zetaclient/context"
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
	_ interfaces.ChainSigner = &Signer{}

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

	// outboundHashBeingReported is a map of outboundHash being reported
	outboundHashBeingReported map[string]bool
}

// NewSigner creates a new EVM signer
func NewSigner(
	chain chains.Chain,
	zetacoreContext *clientcontext.ZetacoreContext,
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
	baseSigner := base.NewSigner(chain, zetacoreContext, tss, ts, logger)

	// create EVM client
	client, ethSigner, err := getEVMRPC(endpoint)
	if err != nil {
		return nil, err
	}

	// prepare ABIs
	connectorABI, err := abi.JSON(strings.NewReader(zetaConnectorABI))
	if err != nil {
		return nil, err
	}
	custodyABI, err := abi.JSON(strings.NewReader(erc20CustodyABI))
	if err != nil {
		return nil, err
	}

	return &Signer{
		Signer:                    baseSigner,
		client:                    client,
		ethSigner:                 ethSigner,
		zetaConnectorABI:          connectorABI,
		erc20CustodyABI:           custodyABI,
		zetaConnectorAddress:      zetaConnectorAddress,
		er20CustodyAddress:        erc20CustodyAddress,
		outboundHashBeingReported: make(map[string]bool),
	}, nil
}

// SetZetaConnectorAddress sets the zeta connector address
func (s *Signer) SetZetaConnectorAddress(addr ethcommon.Address) {
	s.Lock()
	defer s.Unlock()
	s.zetaConnectorAddress = addr
}

// SetERC20CustodyAddress sets the erc20 custody address
func (s *Signer) SetERC20CustodyAddress(addr ethcommon.Address) {
	s.Lock()
	defer s.Unlock()
	s.er20CustodyAddress = addr
}

// GetZetaConnectorAddress returns the zeta connector address
func (s *Signer) GetZetaConnectorAddress() ethcommon.Address {
	s.Lock()
	defer s.Unlock()
	return s.zetaConnectorAddress
}

// GetERC20CustodyAddress returns the erc20 custody address
func (s *Signer) GetERC20CustodyAddress() ethcommon.Address {
	s.Lock()
	defer s.Unlock()
	return s.er20CustodyAddress
}

// Sign given data, and metadata (gas, nonce, etc.)
// returns a signed transaction, sig bytes, hash bytes, and error
func (s *Signer) Sign(
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
	height uint64,
) (*ethtypes.Transaction, []byte, []byte, error) {
	s.Logger().Debug().
		Str("tss_pub_key", string(s.TSS().Pubkey())).
		Msg("Signing evm transaction")

	chainID := big.NewInt(s.Chain().ChainId)
	tx, err := newTx(chainID, data, to, amount, gas, nonce)
	if err != nil {
		return nil, nil, nil, err
	}

	hashBytes := s.ethSigner.Hash(tx).Bytes()

	sig, err := s.TSS().Sign(hashBytes, height, nonce, s.Chain().ChainId, "")
	if err != nil {
		return nil, nil, nil, err
	}

	log.Debug().Msgf("Sign: Signature: %s", hex.EncodeToString(sig[:]))
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		s.Logger().Error().Err(err).Msgf("SigToPub error")
	}

	addr := crypto.PubkeyToAddress(*pubk)
	s.Logger().Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(s.ethSigner, sig[:])
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
			GasPrice: gas.MaxFeePerUnit,
			Gas:      gas.Limit,
			Nonce:    nonce,
		}), nil
	}

	return ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   chainID,
		To:        &to,
		Value:     amount,
		Data:      data,
		GasFeeCap: gas.MaxFeePerUnit,
		GasTipCap: gas.PriorityFeePerUnit,
		Gas:       gas.Limit,
		Nonce:     nonce,
	}), nil
}

// Broadcast takes in signed tx, broadcast to external chain node
func (s *Signer) Broadcast(tx *ethtypes.Transaction) error {
	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.client.SendTransaction(ctxt, tx)
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
func (s *Signer) SignOutbound(txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := s.zetaConnectorABI.Pack("onReceive",
		txData.sender.Bytes(),
		txData.srcChainID,
		txData.to,
		txData.amount,
		txData.message,
		txData.cctxIndex)

	if err != nil {
		return nil, fmt.Errorf("onReceive pack error: %w", err)
	}

	tx, _, _, err := s.Sign(
		data,
		s.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign onReceive")
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
func (s *Signer) SignRevertTx(txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := s.zetaConnectorABI.Pack("onRevert",
		txData.sender,
		txData.srcChainID,
		txData.to.Bytes(),
		txData.toChainID,
		txData.amount,
		txData.message,
		txData.cctxIndex)
	if err != nil {
		return nil, fmt.Errorf("onRevert pack error: %w", err)
	}

	tx, _, _, err := s.Sign(
		data,
		s.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign onRevert")
	}

	return tx, nil
}

// SignCancelTx signs a transaction from TSS address to itself with a zero amount in order to increment the nonce
func (s *Signer) SignCancelTx(txData *OutboundData) (*ethtypes.Transaction, error) {
	// todo LIMIT=evm.EthTransferGasLimit

	tx, _, _, err := s.Sign(
		nil,
		s.TSS().EVMAddress(),
		zeroValue, // zero out the amount to cancel the tx
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign cancellation tx")
	}

	return tx, nil
}

// SignWithdrawTx signs a withdrawal transaction sent from the TSS address to the destination
func (s *Signer) SignWithdrawTx(txData *OutboundData) (*ethtypes.Transaction, error) {
	// todo LIMIT=evm.EthTransferGasLimit

	tx, _, _, err := s.Sign(
		nil,
		txData.to,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign withdraw tx")
	}

	return tx, nil
}

// SignCommandTx signs a transaction based on the given command includes:
//
//	cmd_whitelist_erc20
//	cmd_migrate_tss_funds
func (s *Signer) SignCommandTx(txData *OutboundData, cmd string, params string) (*ethtypes.Transaction, error) {
	switch cmd {
	case constant.CmdWhitelistERC20:
		return s.SignWhitelistERC20Cmd(txData, params)
	case constant.CmdMigrateTssFunds:
		return s.SignMigrateTssFundsCmd(txData)
	}
	return nil, fmt.Errorf("SignCommandTx: unknown command %s", cmd)
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign an evm transaction using the TSS signer.
// It will then broadcast the signed transaction to the outbound chain.
func (s *Signer) TryProcessOutbound(
	cctx *types.CrossChainTx,
	outboundProc *outboundprocessor.Processor,
	outboundID string,
	chainObserver interfaces.ChainObserver,
	zetacoreClient interfaces.ZetacoreClient,
	height uint64,
) {
	logger := s.Logger().With().
		Str("outbound_id", outboundID).
		Str("cctx.index", cctx.Index).
		Logger()

	params := cctx.GetCurrentOutboundParam()

	logger.Info().
		Str("cctx.amount", params.Amount.String()).
		Str("cctx.receiver", params.Receiver).
		Msgf("start processing outbound")

	defer outboundProc.EndTryProcess(outboundID)

	myID := zetacoreClient.GetKeys().GetOperatorAddress()

	evmObserver, ok := chainObserver.(*observer.Observer)
	if !ok {
		logger.Error().Msg("chain observer is not an EVM observer")
		return
	}

	// Setup Transaction input
	txData, skipTx, err := NewOutboundData(cctx, evmObserver, height, logger)
	if err != nil {
		logger.Err(err).Msg("error setting up transaction input fields")
		return
	}
	if skipTx {
		return
	}

	// Get destination chain for logging
	toChain := chains.GetChainFromChainID(txData.toChainID.Int64())

	// Get cross-chain flags
	crossChainflags := s.ZetacoreContext().GetCrossChainFlags()
	// https://github.com/zeta-chain/node/issues/2050
	var tx *ethtypes.Transaction
	// compliance check goes first
	if compliance.IsCctxRestricted(cctx) {
		compliance.PrintComplianceLog(
			logger,
			s.Logger().Compliance,
			true,
			s.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			txData.to.Hex(),
			cctx.GetCurrentOutboundParam().CoinType.String(),
		)

		tx, err = s.SignCancelTx(txData) // cancel the tx
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
		tx, err = s.SignCommandTx(txData, cmd, params)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if IsSenderZetaChain(cctx, zetacoreClient, &crossChainflags) {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignWithdrawTx: %d => %s, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain,
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.MaxFeePerUnit,
			)
			tx, err = s.SignWithdrawTx(txData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf(
				"SignERC20WithdrawTx: %d => %s, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain,
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.MaxFeePerUnit,
			)
			tx, err = s.SignERC20WithdrawTx(txData)
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignOutbound: %d => %s, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain,
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.MaxFeePerUnit,
			)
			tx, err = s.SignOutbound(txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.OutboundParams[0].ReceiverChainId == zetacoreClient.Chain().ChainId {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignRevertTx: %d => %s, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain, cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.MaxFeePerUnit,
			)
			txData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
			txData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
			tx, err = s.SignRevertTx(txData)
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignWithdrawTx: %d => %s, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain,
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.MaxFeePerUnit,
			)
			tx, err = s.SignWithdrawTx(txData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %s, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain,
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.MaxFeePerUnit,
			)
			tx, err = s.SignERC20WithdrawTx(txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		logger.Info().Msgf(
			"SignRevertTx: %d => %s, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain,
			cctx.GetCurrentOutboundParam().TssNonce,
			txData.gas.MaxFeePerUnit,
		)
		txData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)

		tx, err = s.SignRevertTx(txData)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf(
			"SignOutbound: %d => %s, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain,
			cctx.GetCurrentOutboundParam().TssNonce,
			txData.gas.MaxFeePerUnit,
		)
		tx, err = s.SignOutbound(txData)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	}

	logger.Info().Msgf(
		"Key-sign success: %d => %s, nonce %d",
		cctx.InboundParams.SenderChainId,
		toChain,
		cctx.GetCurrentOutboundParam().TssNonce,
	)

	// Broadcast Signed Tx
	s.BroadcastOutbound(tx, cctx, logger, myID, zetacoreClient, txData)
}

// BroadcastOutbound signed transaction through evm rpc client
func (s *Signer) BroadcastOutbound(
	tx *ethtypes.Transaction,
	cctx *types.CrossChainTx,
	logger zerolog.Logger,
	myID sdk.AccAddress,
	zetacoreClient interfaces.ZetacoreClient,
	txData *OutboundData,
) {
	// Get destination chain for logging
	toChain := chains.GetChainFromChainID(txData.toChainID.Int64())
	if tx == nil {
		logger.Warn().Msgf("BroadcastOutbound: no tx to broadcast %s", cctx.Index)
	}

	// broadcast transaction
	if tx != nil {
		outboundHash := tx.Hash().Hex()

		// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s) in case of RPC error
		backOff := broadcastBackoff
		for i := 0; i < broadcastRetries; i++ {
			time.Sleep(backOff)
			err := s.Broadcast(tx)
			if err != nil {
				log.Warn().
					Err(err).
					Msgf("BroadcastOutbound: error broadcasting tx %s on chain %d nonce %d retry %d signer %s",
						outboundHash, toChain.ChainId, cctx.GetCurrentOutboundParam().TssNonce, i, myID)
				retry, report := zetacore.HandleBroadcastError(
					err,
					strconv.FormatUint(cctx.GetCurrentOutboundParam().TssNonce, 10),
					toChain.String(),
					outboundHash,
				)
				if report {
					s.reportToOutboundTracker(zetacoreClient, toChain.ChainId, tx.Nonce(), outboundHash, logger)
				}
				if !retry {
					break
				}
				backOff *= 2
				continue
			}
			logger.Info().Msgf("BroadcastOutbound: broadcasted tx %s on chain %d nonce %d signer %s",
				outboundHash, toChain.ChainId, cctx.GetCurrentOutboundParam().TssNonce, myID)
			s.reportToOutboundTracker(zetacoreClient, toChain.ChainId, tx.Nonce(), outboundHash, logger)
			break // successful broadcast; no need to retry
		}
	}
}

// SignERC20WithdrawTx
// function withdraw(
// address recipient,
// address asset,
// uint256 amount,
// ) external onlyTssAddress
func (s *Signer) SignERC20WithdrawTx(txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := s.erc20CustodyABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := s.Sign(
		data,
		s.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign withdraw")
	}

	return tx, nil
}

// Exported for unit tests

// GetReportedTxList returns a list of outboundHash being reported
// TODO: investigate pointer usage
// https://github.com/zeta-chain/node/issues/2084
func (s *Signer) GetReportedTxList() *map[string]bool {
	return &s.outboundHashBeingReported
}

func (s *Signer) EvmClient() interfaces.EVMRPCClient {
	return s.client
}

func (s *Signer) EvmSigner() ethtypes.Signer {
	return s.ethSigner
}

func IsSenderZetaChain(
	cctx *types.CrossChainTx,
	zetacoreClient interfaces.ZetacoreClient,
	flags *observertypes.CrosschainFlags,
) bool {
	return cctx.InboundParams.SenderChainId == zetacoreClient.Chain().ChainId &&
		cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound && flags.IsOutboundEnabled
}

func ErrorMsg(cctx *types.CrossChainTx) string {
	return fmt.Sprintf(
		"signer SignOutbound error: nonce %d chain %d",
		cctx.GetCurrentOutboundParam().TssNonce,
		cctx.GetCurrentOutboundParam().ReceiverChainId,
	)
}

func (s *Signer) SignWhitelistERC20Cmd(txData *OutboundData, params string) (*ethtypes.Transaction, error) {
	outboundParams := txData.outboundParams

	erc20 := ethcommon.HexToAddress(params)
	if erc20 == (ethcommon.Address{}) {
		return nil, fmt.Errorf("SignCommandTx: invalid erc20 address %s", params)
	}

	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := custodyAbi.Pack("whitelist", erc20)
	if err != nil {
		return nil, fmt.Errorf("whitelist pack error: %w", err)
	}

	tx, _, _, err := s.Sign(
		data,
		txData.to,
		zeroValue,
		txData.gas,
		outboundParams.TssNonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign whitelist error: %w", err)
	}
	return tx, nil
}

func (s *Signer) SignMigrateTssFundsCmd(txData *OutboundData) (*ethtypes.Transaction, error) {
	tx, _, _, err := s.Sign(
		nil,
		txData.to,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign migrate tss funds")
	}

	return tx, nil
}

// reportToOutboundTracker reports outboundHash to tracker only when tx receipt is available
func (s *Signer) reportToOutboundTracker(
	zetacoreClient interfaces.ZetacoreClient,
	chainID int64,
	nonce uint64,
	outboundHash string,
	logger zerolog.Logger,
) {
	// skip if already being reported
	s.Lock()
	defer s.Unlock()
	if _, found := s.outboundHashBeingReported[outboundHash]; found {
		logger.Info().
			Msgf("reportToOutboundTracker: outboundHash %s for chain %d nonce %d is being reported", outboundHash, chainID, nonce)
		return
	}
	s.outboundHashBeingReported[outboundHash] = true // mark as being reported

	// report to outbound tracker with goroutine
	go func() {
		defer func() {
			s.Lock()
			delete(s.outboundHashBeingReported, outboundHash)
			s.Unlock()
		}()

		// try monitoring tx inclusion status for 10 minutes
		var err error
		report := false
		isPending := false
		blockNumber := uint64(0)
		tStart := time.Now()
		for {
			// give up after 10 minutes of monitoring
			time.Sleep(10 * time.Second)

			if time.Since(tStart) > evm.OutboundInclusionTimeout {
				// if tx is still pending after timeout, report to outboundTracker anyway as we cannot monitor forever
				if isPending {
					report = true // probably will be included later
				}
				logger.Info().
					Msgf("reportToOutboundTracker: timeout waiting tx inclusion for chain %d nonce %d outboundHash %s report %v", chainID, nonce, outboundHash, report)
				break
			}
			// try getting the tx
			_, isPending, err = s.client.TransactionByHash(context.TODO(), ethcommon.HexToHash(outboundHash))
			if err != nil {
				logger.Info().
					Err(err).
					Msgf("reportToOutboundTracker: error getting tx for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				continue
			}
			// if tx is include in a block, try getting receipt
			if !isPending {
				report = true // included
				receipt, err := s.client.TransactionReceipt(context.TODO(), ethcommon.HexToHash(outboundHash))
				if err != nil {
					logger.Info().
						Err(err).
						Msgf("reportToOutboundTracker: error getting receipt for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				}
				if receipt != nil {
					blockNumber = receipt.BlockNumber.Uint64()
				}
				break
			}
			// keep monitoring pending tx
			logger.Info().
				Msgf("reportToOutboundTracker: tx has not been included yet for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
		}

		// try adding to outbound tracker for 10 minutes
		if report {
			tStart := time.Now()
			for {
				// give up after 10 minutes of retrying
				if time.Since(tStart) > evm.OutboundTrackerReportTimeout {
					logger.Info().
						Msgf("reportToOutboundTracker: timeout adding outbound tracker for chain %d nonce %d outboundHash %s, please add manually", chainID, nonce, outboundHash)
					break
				}
				// stop if the cctx is already finalized
				cctx, err := zetacoreClient.GetCctxByNonce(chainID, nonce)
				if err != nil {
					logger.Err(err).
						Msgf("reportToOutboundTracker: error getting cctx for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				} else if !crosschainkeeper.IsPending(cctx) {
					logger.Info().Msgf("reportToOutboundTracker: cctx already finalized for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
					break
				}
				// report to outbound tracker
				zetaHash, err := zetacoreClient.AddOutboundTracker(chainID, nonce, outboundHash, nil, "", -1)
				if err != nil {
					logger.Err(err).
						Msgf("reportToOutboundTracker: error adding to outbound tracker for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				} else if zetaHash != "" {
					logger.Info().Msgf("reportToOutboundTracker: added outboundHash to core successful %s, chain %d nonce %d outboundHash %s block %d",
						zetaHash, chainID, nonce, outboundHash, blockNumber)
				} else {
					// stop if the tracker contains the outboundHash
					logger.Info().Msgf("reportToOutboundTracker: outbound tracker contains outboundHash %s for chain %d nonce %d", outboundHash, chainID, nonce)
					break
				}
				// retry otherwise
				time.Sleep(evm.ZetaBlockTime * 3)
			}
		}
	}()
}

// getEVMRPC is a helper function to set up the client and signer, also initializes a mock client for unit tests
func getEVMRPC(endpoint string) (interfaces.EVMRPCClient, ethtypes.Signer, error) {
	if endpoint == mocks.EVMRPCEnabled {
		chainID := big.NewInt(chains.BscMainnet.ChainId)
		ethSigner := ethtypes.NewLondonSigner(chainID)
		client := &mocks.MockEvmClient{}
		return client, ethSigner, nil
	}

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, nil, err
	}
	ethSigner := ethtypes.LatestSignerForChainID(chainID)
	return client, ethSigner, nil
}
