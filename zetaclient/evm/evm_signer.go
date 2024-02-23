package evm

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outtxprocessor"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/zetacore/common"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	zbridge "github.com/zeta-chain/zetacore/zetaclient/zetabridge"
)

type TransactionData struct {
	srcChainID *big.Int
	toChainID  *big.Int
	sender     ethcommon.Address
	to         ethcommon.Address
	asset      ethcommon.Address
	amount     *big.Int
	gasPrice   *big.Int
	gasLimit   uint64
	message    []byte
	sendHash   [32]byte
	nonce      uint64
	height     uint64

	cmd            string
	params         string
	outboundParams *types.OutboundTxParams
	flags          observertypes.CrosschainFlags
}

const (
	OutTxInclusionTimeout     = 20 * time.Minute
	OutTxTrackerReportTimeout = 10 * time.Minute
	ZetaBlockTime             = 6500 * time.Millisecond
)

type Signer struct {
	client                      interfaces.EVMRPCClient
	chain                       *common.Chain
	chainID                     *big.Int
	tssSigner                   interfaces.TSSSigner
	ethSigner                   ethtypes.Signer
	abi                         abi.ABI
	erc20CustodyABI             abi.ABI
	metaContractAddress         ethcommon.Address
	erc20CustodyContractAddress ethcommon.Address
	logger                      zerolog.Logger
	ts                          *metrics.TelemetryServer

	// for outTx tracker report only
	mu                     *sync.Mutex
	outTxHashBeingReported map[string]bool
}

var _ interfaces.ChainSigner = &Signer{}

func NewEVMSigner(
	chain common.Chain,
	endpoint string,
	tssSigner interfaces.TSSSigner,
	abiString string,
	erc20CustodyABIString string,
	metaContract ethcommon.Address,
	erc20CustodyContract ethcommon.Address,
	logger zerolog.Logger,
	ts *metrics.TelemetryServer,
) (*Signer, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}
	ethSigner := ethtypes.LatestSignerForChainID(chainID)
	connectorABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	erc20CustodyABI, err := abi.JSON(strings.NewReader(erc20CustodyABIString))
	if err != nil {
		return nil, err
	}

	return &Signer{
		client:                      client,
		chain:                       &chain,
		tssSigner:                   tssSigner,
		chainID:                     chainID,
		ethSigner:                   ethSigner,
		abi:                         connectorABI,
		erc20CustodyABI:             erc20CustodyABI,
		metaContractAddress:         metaContract,
		erc20CustodyContractAddress: erc20CustodyContract,
		logger: logger.With().
			Str("chain", chain.ChainName.String()).
			Str("module", "EVMSigner").Logger(),
		ts:                     ts,
		mu:                     &sync.Mutex{},
		outTxHashBeingReported: make(map[string]bool),
	}, nil
}

// Sign given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *Signer) Sign(
	data []byte,
	to ethcommon.Address,
	gasLimit uint64,
	gasPrice *big.Int,
	nonce uint64,
	height uint64,
) (*ethtypes.Transaction, []byte, []byte, error) {
	log.Debug().Msgf("TSS SIGNER: %s", signer.tssSigner.Pubkey())
	tx := ethtypes.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()

	sig, err := signer.tssSigner.Sign(hashBytes, height, nonce, signer.chain, "")
	if err != nil {
		return nil, nil, nil, err
	}
	log.Debug().Msgf("Sign: Signature: %s", hex.EncodeToString(sig[:]))
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, nil, nil, err
	}
	return signedTX, sig[:], hashBytes[:], nil
}

// Broadcast takes in signed tx, broadcast to external chain node
func (signer *Signer) Broadcast(tx *ethtypes.Transaction) error {
	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return signer.client.SendTransaction(ctxt, tx)
}

// SignOutboundTx
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
func (signer *Signer) SignOutboundTx(txData *TransactionData) (*ethtypes.Transaction, error) {

	if len(txData.sendHash) < 32 {
		return nil, fmt.Errorf("sendHash len %d must be 32", len(txData.sendHash))
	}
	var data []byte
	var err error

	data, err = signer.abi.Pack("onReceive",
		txData.sender.Bytes(),
		txData.srcChainID,
		txData.to,
		txData.amount,
		txData.message,
		txData.sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data,
		signer.metaContractAddress,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
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
func (signer *Signer) SignRevertTx(txData *TransactionData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	data, err = signer.abi.Pack("onRevert",
		txData.sender,
		txData.srcChainID,
		txData.to,
		txData.toChainID,
		txData.amount,
		txData.message,
		txData.sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data,
		signer.metaContractAddress,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
	}

	return tx, nil
}

func (signer *Signer) SignCancelTx(nonce uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, signer.tssSigner.EVMAddress(), big.NewInt(0), 21000, gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes, height, nonce, signer.chain, "")
	if err != nil {
		return nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, err
	}

	return signedTX, nil
}

func (signer *Signer) SignWithdrawTx(txData *TransactionData) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(txData.nonce, txData.to, txData.amount, 21000, txData.gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes, txData.height, txData.nonce, signer.chain, "")
	if err != nil {
		return nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, err
	}

	return signedTX, nil
}

func (signer *Signer) SignCommandTx(txData *TransactionData) (*ethtypes.Transaction, error) {
	outboundParams := txData.outboundParams
	if txData.cmd == common.CmdWhitelistERC20 {
		erc20 := ethcommon.HexToAddress(txData.params)
		if erc20 == (ethcommon.Address{}) {
			return nil, fmt.Errorf("SignCommandTx: invalid erc20 address %s", txData.params)
		}
		custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
		if err != nil {
			return nil, err
		}
		data, err := custodyAbi.Pack("whitelist", erc20)
		if err != nil {
			return nil, err
		}
		tx, _, _, err := signer.Sign(data, txData.to, txData.gasLimit, txData.gasPrice, outboundParams.OutboundTxTssNonce, txData.height)
		if err != nil {
			return nil, fmt.Errorf("sign error: %w", err)
		}
		return tx, nil
	}
	if txData.cmd == common.CmdMigrateTssFunds {
		tx := ethtypes.NewTransaction(outboundParams.OutboundTxTssNonce, txData.to, outboundParams.Amount.BigInt(), outboundParams.OutboundTxGasLimit, txData.gasPrice, nil)
		hashBytes := signer.ethSigner.Hash(tx).Bytes()
		sig, err := signer.tssSigner.Sign(hashBytes, txData.height, outboundParams.OutboundTxTssNonce, signer.chain, "")
		if err != nil {
			return nil, err
		}
		pubk, err := crypto.SigToPub(hashBytes, sig[:])
		if err != nil {
			signer.logger.Error().Err(err).Msgf("SigToPub error")
		}
		addr := crypto.PubkeyToAddress(*pubk)
		signer.logger.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
		signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
		if err != nil {
			return nil, err
		}

		return signedTX, nil
	}

	return nil, fmt.Errorf("SignCommandTx: unknown command %s", txData.cmd)
}

func setChainAndSender(cctx *types.CrossChainTx, logger zerolog.Logger, txData *TransactionData) bool {
	if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		txData.to = ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
		txData.toChainID = big.NewInt(cctx.InboundTxParams.SenderChainId) //common.GetChainFromChainID(cctx.InboundTxParams.SenderChainId)
		logger.Info().Msgf("Abort: reverting inbound")
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		txData.to = ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutTxParam().ReceiverChainId) //common.GetChainFromChainID(cctx.GetCurrentOutTxParam().ReceiverChainId)
	} else {
		logger.Info().Msgf("Transaction doesn't need to be processed status: %d", cctx.CctxStatus.Status)
		return true
	}
	return false
}

func setupGas(cctx *types.CrossChainTx,
	logger zerolog.Logger,
	client interfaces.EVMRPCClient,
	chain *common.Chain,
	txData *TransactionData) error {

	txData.gasLimit = cctx.GetCurrentOutTxParam().OutboundTxGasLimit
	if txData.gasLimit < 100_000 {
		txData.gasLimit = 100_000
		logger.Warn().Msgf("gasLimit %d is too low; set to %d", cctx.GetCurrentOutTxParam().OutboundTxGasLimit, txData.gasLimit)
	}
	if txData.gasLimit > 1_000_000 {
		txData.gasLimit = 1_000_000
		logger.Warn().Msgf("gasLimit %d is too high; set to %d", cctx.GetCurrentOutTxParam().OutboundTxGasLimit, txData.gasLimit)
	}

	// use dynamic gas price for ethereum chains.
	// The code below is a fix for https://github.com/zeta-chain/node/issues/1085
	// doesn't close directly the issue because we should determine if we want to keep using SuggestGasPrice if no OutboundTxGasPrice
	// we should possibly remove it completely and return an error if no OutboundTxGasPrice is provided because it means no fee is processed on ZetaChain
	specified, ok := new(big.Int).SetString(cctx.GetCurrentOutTxParam().OutboundTxGasPrice, 10)
	if !ok {
		if common.IsEthereumChain(chain.ChainId) {
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

func setTransactionData(
	cctx *types.CrossChainTx,
	evmClient *ChainClient,
	evmRPC interfaces.EVMRPCClient,
	logger zerolog.Logger,
	zetaBridge interfaces.ZetaCoreBridger,
	txData *TransactionData) (bool, error) {

	txData.outboundParams = cctx.GetCurrentOutTxParam()
	txData.amount = cctx.GetCurrentOutTxParam().Amount.BigInt()
	txData.nonce = cctx.GetCurrentOutTxParam().OutboundTxTssNonce
	txData.sender = ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
	txData.srcChainID = big.NewInt(cctx.InboundTxParams.SenderChainId)
	txData.asset = ethcommon.HexToAddress(cctx.InboundTxParams.Asset)

	skipTx := setChainAndSender(cctx, logger, txData)
	if skipTx {
		return true, nil
	}

	toChain := common.GetChainFromChainID(txData.toChainID.Int64())
	if toChain == nil {
		return true, fmt.Errorf("unknown chain: %d", txData.toChainID.Int64())
	}

	// Get nonce, Early return if the cctx is already processed
	nonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce
	included, confirmed, err := evmClient.IsSendOutTxProcessed(cctx.Index, nonce, cctx.GetCurrentOutTxParam().CoinType, logger)
	if err != nil {
		return true, errors.New("IsSendOutTxProcessed failed")
	}
	if included || confirmed {
		logger.Info().Msgf("CCTX already processed; exit signer")
		return true, nil
	}

	// Set up gas limit and gas price
	err = setupGas(cctx, logger, evmRPC, toChain, txData)
	if err != nil {
		return true, err
	}

	// Get sendHash
	logger.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain, cctx.InboundTxParams.Amount, txData.to.Hex(), nonce, cctx.InboundTxParams.InboundTxFinalizedZetaHeight)
	sendHash, err := hex.DecodeString(cctx.Index[2:]) // remove the leading 0x
	if err != nil || len(sendHash) != 32 {
		return true, fmt.Errorf("decode CCTX %s error", cctx.Index)
	}
	copy(txData.sendHash[:32], sendHash[:32])

	// In case there is a pending transaction, make sure this keysign is a transaction replacement
	pendingTx := evmClient.GetPendingTx(nonce)
	if pendingTx != nil {
		if txData.gasPrice.Cmp(pendingTx.GasPrice()) > 0 {
			logger.Info().Msgf("replace pending outTx %s nonce %d using gas price %d", pendingTx.Hash().Hex(), nonce, txData.gasPrice)
		} else {
			logger.Info().Msgf("please wait for pending outTx %s nonce %d to be included", pendingTx.Hash().Hex(), nonce)
			return true, nil
		}
	}

	// Base64 decode message
	if cctx.GetCurrentOutTxParam().CoinType != common.CoinType_Cmd {
		txData.message, err = base64.StdEncoding.DecodeString(cctx.RelayedMessage)
		if err != nil {
			logger.Err(err).Msgf("decode CCTX.Message %s error", cctx.RelayedMessage)
		}
	}

	// Get cross-chain flags
	txData.flags, err = zetaBridge.GetCrosschainFlags()
	if err != nil {
		return true, errors.New("cannot get crosschain flags")
	}

	return false, nil
}

func (signer *Signer) TryProcessOutTx(
	cctx *types.CrossChainTx,
	outTxMan *outtxprocessor.Processor,
	outTxID string,
	chainclient interfaces.ChainClient,
	zetaBridge interfaces.ZetaCoreBridger,
	height uint64,
) {
	logger := signer.logger.With().
		Str("outTxID", outTxID).
		Str("SendHash", cctx.Index).
		Logger()
	logger.Info().Msgf("start processing outTxID %s", outTxID)
	logger.Info().Msgf("EVM Chain TryProcessOutTx: %s, value %d to %s", cctx.Index, cctx.GetCurrentOutTxParam().Amount.BigInt(), cctx.GetCurrentOutTxParam().Receiver)

	defer func() {
		outTxMan.EndTryProcess(outTxID)
	}()
	myID := zetaBridge.GetKeys().GetOperatorAddress()

	evmClient, ok := chainclient.(*ChainClient)
	if !ok {
		logger.Error().Msg("chain client is not an EVMChainClient")
		return
	}

	// Setup Transaction input
	txData := TransactionData{}
	txData.height = height
	skipTx, err := setTransactionData(cctx, evmClient, signer.client, logger, zetaBridge, &txData)
	if err != nil {
		logger.Error().Msg(err.Error())
		return
	}
	if skipTx {
		return
	}

	// Get destination chain for logging
	toChain := common.GetChainFromChainID(txData.toChainID.Int64())

	var tx *ethtypes.Transaction
	// Sign transaction
	if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Cmd { // admin command
		to := ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver)
		if to == (ethcommon.Address{}) {
			logger.Error().Msgf("invalid receiver %s", cctx.GetCurrentOutTxParam().Receiver)
			return
		}
		msg := strings.Split(cctx.RelayedMessage, ":")
		if len(msg) != 2 {
			logger.Error().Msgf("invalid message %s", msg)
			return
		}
		txData.cmd = msg[0]
		txData.params = msg[1]
		tx, err = signer.SignCommandTx(&txData)
	} else if cctx.InboundTxParams.SenderChainId == zetaBridge.ZetaChain().ChainId && cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound && txData.flags.IsOutboundEnabled {
		if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Gas {
			logger.Info().Msgf("SignWithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignWithdrawTx(&txData)
		}
		if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_ERC20 {

			logger.Info().Msgf("SignERC20WithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignERC20WithdrawTx(&txData)
		}
		if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Zeta {
			logger.Info().Msgf("SignOutboundTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignOutboundTx(&txData)
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.OutboundTxParams[0].ReceiverChainId == zetaBridge.ZetaChain().ChainId {
		if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Gas {
			logger.Info().Msgf("SignWithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignWithdrawTx(&txData)
		}
		if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_ERC20 {
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignERC20WithdrawTx(&txData)
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		logger.Info().Msgf("SignRevertTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
		txData.srcChainID = big.NewInt(cctx.OutboundTxParams[0].ReceiverChainId)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutTxParam().ReceiverChainId)

		tx, err = signer.SignRevertTx(&txData)
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf("SignOutboundTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
		tx, err = signer.SignOutboundTx(&txData)
	}

	if err != nil {
		logger.Warn().Err(err).Msgf("signer SignOutbound error: nonce %d chain %d", cctx.GetCurrentOutTxParam().OutboundTxTssNonce, cctx.GetCurrentOutTxParam().ReceiverChainId)
		return
	}
	logger.Info().Msgf("Key-sign success: %d => %s, nonce %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)

	_, err = zetaBridge.GetObserverList()
	if err != nil {
		logger.Warn().Err(err).Msgf("unable to get observer list: chain %d observation %s", cctx.GetCurrentOutTxParam().OutboundTxTssNonce, observertypes.ObservationType_OutBoundTx.String())

	}

	// Broadcast Signed Tx
	signer.broadcastOutTx(tx, cctx, logger, myID, zetaBridge, &txData)
}

func (signer *Signer) broadcastOutTx(
	tx *ethtypes.Transaction,
	cctx *types.CrossChainTx,
	logger zerolog.Logger,
	myID sdk.AccAddress,
	zetaBridge interfaces.ZetaCoreBridger,
	txData *TransactionData) {
	// Get destination chain for logging
	toChain := common.GetChainFromChainID(txData.toChainID.Int64())

	// Try to broadcast transaction
	if tx != nil {
		outTxHash := tx.Hash().Hex()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", signer.chain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, outTxHash, myID)
		//if len(signers) == 0 || myid == signers[send.OutboundTxParams.Broadcaster] || myid == signers[int(send.OutboundTxParams.Broadcaster+1)%len(signers)] {
		backOff := 1000 * time.Millisecond
		// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
		for i := 0; i < 5; i++ {
			logger.Info().Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, i)
			// #nosec G404 randomness is not a security issue here
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) // FIXME: use backoff
			err := signer.Broadcast(tx)
			if err != nil {
				log.Warn().Err(err).Msgf("OutTx Broadcast error")
				retry, report := zbridge.HandleBroadcastError(err, strconv.FormatUint(cctx.GetCurrentOutTxParam().OutboundTxTssNonce, 10), toChain.String(), outTxHash)
				if report {
					signer.reportToOutTxTracker(zetaBridge, toChain.ChainId, tx.Nonce(), outTxHash, logger)
				}
				if !retry {
					break
				}
				backOff *= 2
				continue
			}
			logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", cctx.GetCurrentOutTxParam().OutboundTxTssNonce, toChain, outTxHash)
			signer.reportToOutTxTracker(zetaBridge, toChain.ChainId, tx.Nonce(), outTxHash, logger)
			break // successful broadcast; no need to retry
		}
	}
}

// reportToOutTxTracker reports outTxHash to tracker only when tx receipt is available
func (signer *Signer) reportToOutTxTracker(zetaBridge interfaces.ZetaCoreBridger, chainID int64, nonce uint64, outTxHash string, logger zerolog.Logger) {
	// skip if already being reported
	signer.mu.Lock()
	defer signer.mu.Unlock()
	if _, found := signer.outTxHashBeingReported[outTxHash]; found {
		logger.Info().Msgf("reportToOutTxTracker: outTxHash %s for chain %d nonce %d is being reported", outTxHash, chainID, nonce)
		return
	}
	signer.outTxHashBeingReported[outTxHash] = true // mark as being reported

	// report to outTx tracker with goroutine
	go func() {
		defer func() {
			signer.mu.Lock()
			delete(signer.outTxHashBeingReported, outTxHash)
			signer.mu.Unlock()
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
			if time.Since(tStart) > OutTxInclusionTimeout {
				// if tx is still pending after timeout, report to outTxTracker anyway as we cannot monitor forever
				if isPending {
					report = true // probably will be included later
				}
				logger.Info().Msgf("reportToOutTxTracker: timeout waiting tx inclusion for chain %d nonce %d outTxHash %s report %v", chainID, nonce, outTxHash, report)
				break
			}
			// try getting the tx
			_, isPending, err = signer.client.TransactionByHash(context.TODO(), ethcommon.HexToHash(outTxHash))
			if err != nil {
				logger.Info().Err(err).Msgf("reportToOutTxTracker: error getting tx for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
				continue
			}
			// if tx is include in a block, try getting receipt
			if !isPending {
				report = true // included
				receipt, err := signer.client.TransactionReceipt(context.TODO(), ethcommon.HexToHash(outTxHash))
				if err != nil {
					logger.Info().Err(err).Msgf("reportToOutTxTracker: error getting receipt for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
				}
				if receipt != nil {
					blockNumber = receipt.BlockNumber.Uint64()
				}
				break
			}
			// keep monitoring pending tx
			logger.Info().Msgf("reportToOutTxTracker: tx has not been included yet for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
		}

		// try adding to outTx tracker for 10 minutes
		if report {
			tStart := time.Now()
			for {
				// give up after 10 minutes of retrying
				if time.Since(tStart) > OutTxTrackerReportTimeout {
					logger.Info().Msgf("reportToOutTxTracker: timeout adding outtx tracker for chain %d nonce %d outTxHash %s, please add manually", chainID, nonce, outTxHash)
					break
				}
				// stop if the cctx is already finalized
				cctx, err := zetaBridge.GetCctxByNonce(chainID, nonce)
				if err != nil {
					logger.Err(err).Msgf("reportToOutTxTracker: error getting cctx for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
				} else if !crosschainkeeper.IsPending(*cctx) {
					logger.Info().Msgf("reportToOutTxTracker: cctx already finalized for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
					break
				}
				// report to outTx tracker
				zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(chainID, nonce, outTxHash, nil, "", -1)
				if err != nil {
					logger.Err(err).Msgf("reportToOutTxTracker: error adding to outtx tracker for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
				} else if zetaHash != "" {
					logger.Info().Msgf("reportToOutTxTracker: added outTxHash to core successful %s, chain %d nonce %d outTxHash %s block %d",
						zetaHash, chainID, nonce, outTxHash, blockNumber)
				} else {
					// stop if the tracker contains the outTxHash
					logger.Info().Msgf("reportToOutTxTracker: outtx tracker contains outTxHash %s for chain %d nonce %d", outTxHash, chainID, nonce)
					break
				}
				// retry otherwise
				time.Sleep(ZetaBlockTime * 3)
			}
		}
	}()
}

// SignERC20WithdrawTx
// function withdraw(
// address recipient,
// address asset,
// uint256 amount,
// ) external onlyTssAddress
func (signer *Signer) SignERC20WithdrawTx(txData *TransactionData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error
	data, err = signer.erc20CustodyABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.erc20CustodyContractAddress, txData.gasLimit, txData.gasPrice, txData.nonce, txData.height)
	if err != nil {
		return nil, fmt.Errorf("sign error: %w", err)
	}

	return tx, nil
}

// SignWhitelistTx
// function whitelist(
// address asset,
// ) external onlyTssAddress
// function unwhitelist(
// address asset,
// ) external onlyTssAddress
func (signer *Signer) SignWhitelistTx(
	action string,
	_ ethcommon.Address,
	asset ethcommon.Address,
	gasLimit uint64,
	nonce uint64,
	gasPrice *big.Int,
	height uint64,
) (*ethtypes.Transaction, error) {
	var data []byte

	var err error

	data, err = signer.erc20CustodyABI.Pack(action, asset)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.erc20CustodyContractAddress, gasLimit, gasPrice, nonce, height)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
	}

	return tx, nil
}

func roundUpToNearestGwei(gasPrice *big.Int) *big.Int {
	oneGwei := big.NewInt(1_000_000_000) // 1 Gwei
	mod := new(big.Int)
	mod.Mod(gasPrice, oneGwei)
	if mod.Cmp(big.NewInt(0)) == 0 { // gasprice is already a multiple of 1 Gwei
		return gasPrice
	}
	return new(big.Int).Add(gasPrice, new(big.Int).Sub(oneGwei, mod))
}
