package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

const (
	OutTxTrackerReportTimeout = 10 * time.Minute
)

type EVMSigner struct {
	client                      EVMRPCClient
	chain                       *common.Chain
	chainID                     *big.Int
	tssSigner                   TSSSigner
	ethSigner                   ethtypes.Signer
	abi                         abi.ABI
	erc20CustodyABI             abi.ABI
	metaContractAddress         ethcommon.Address
	erc20CustodyContractAddress ethcommon.Address
	logger                      zerolog.Logger
	ts                          *TelemetryServer

	// for outTx tracker report only
	mu                     *sync.Mutex
	outTxHashBeingReported map[string]bool
}

var _ ChainSigner = &EVMSigner{}

func NewEVMSigner(
	chain common.Chain,
	endpoint string,
	tssSigner TSSSigner,
	abiString string,
	erc20CustodyABIString string,
	metaContract ethcommon.Address,
	erc20CustodyContract ethcommon.Address,
	logger zerolog.Logger,
	ts *TelemetryServer,
) (*EVMSigner, error) {
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

	return &EVMSigner{
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
func (signer *EVMSigner) Sign(
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
func (signer *EVMSigner) Broadcast(tx *ethtypes.Transaction) error {
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
func (signer *EVMSigner) SignOutboundTx(sender ethcommon.Address,
	srcChainID *big.Int,
	to ethcommon.Address,
	amount *big.Int,
	gasLimit uint64,
	message []byte,
	sendHash [32]byte,
	nonce uint64,
	gasPrice *big.Int,
	height uint64) (*ethtypes.Transaction, error) {

	if len(sendHash) < 32 {
		return nil, fmt.Errorf("sendHash len %d must be 32", len(sendHash))
	}
	var data []byte
	var err error

	data, err = signer.abi.Pack("onReceive", sender.Bytes(), srcChainID, to, amount, message, sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce, height)
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
func (signer *EVMSigner) SignRevertTx(
	sender ethcommon.Address,
	srcChainID *big.Int,
	to []byte,
	toChainID *big.Int,
	amount *big.Int,
	gasLimit uint64,
	message []byte,
	sendHash [32]byte,
	nonce uint64,
	gasPrice *big.Int,
	height uint64,
) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	data, err = signer.abi.Pack("onRevert", sender, srcChainID, to, toChainID, amount, message, sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce, height)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
	}

	return tx, nil
}

func (signer *EVMSigner) SignCancelTx(nonce uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
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

func (signer *EVMSigner) SignWithdrawTx(
	to ethcommon.Address,
	amount *big.Int,
	nonce uint64,
	gasPrice *big.Int,
	height uint64,
) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)
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

func (signer *EVMSigner) SignCommandTx(
	cmd string,
	params string,
	to ethcommon.Address,
	outboundParams *types.OutboundTxParams,
	gasLimit uint64,
	gasPrice *big.Int,
	height uint64,
) (*ethtypes.Transaction, error) {
	if cmd == common.CmdWhitelistERC20 {
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
			return nil, err
		}
		tx, _, _, err := signer.Sign(data, to, gasLimit, gasPrice, outboundParams.OutboundTxTssNonce, height)
		if err != nil {
			return nil, fmt.Errorf("sign error: %w", err)
		}
		return tx, nil
	}
	if cmd == common.CmdMigrateTssFunds {
		tx := ethtypes.NewTransaction(outboundParams.OutboundTxTssNonce, to, outboundParams.Amount.BigInt(), 21000, gasPrice, nil)
		hashBytes := signer.ethSigner.Hash(tx).Bytes()
		sig, err := signer.tssSigner.Sign(hashBytes, height, outboundParams.OutboundTxTssNonce, signer.chain, "")
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

	return nil, fmt.Errorf("SignCommandTx: unknown command %s", cmd)
}

func (signer *EVMSigner) TryProcessOutTx(
	send *types.CrossChainTx,
	outTxMan *OutTxProcessorManager,
	outTxID string,
	chainclient ChainClient,
	zetaBridge ZetaCoreBridger,
	height uint64,
) {
	logger := signer.logger.With().
		Str("outTxID", outTxID).
		Str("SendHash", send.Index).
		Logger()
	logger.Info().Msgf("start processing outTxID %s", outTxID)
	logger.Info().Msgf("EVM Chain TryProcessOutTx: %s, value %d to %s", send.Index, send.GetCurrentOutTxParam().Amount.BigInt(), send.GetCurrentOutTxParam().Receiver)

	defer func() {
		outTxMan.EndTryProcess(outTxID)
	}()
	myID := zetaBridge.GetKeys().GetOperatorAddress()

	var to ethcommon.Address
	var toChain *common.Chain
	if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		to = ethcommon.HexToAddress(send.InboundTxParams.Sender)
		toChain = common.GetChainFromChainID(send.InboundTxParams.SenderChainId)
		if toChain == nil {
			logger.Error().Msgf("Unknown chain: %d", send.InboundTxParams.SenderChainId)
			return
		}
		logger.Info().Msgf("Abort: reverting inbound")
	} else if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		to = ethcommon.HexToAddress(send.GetCurrentOutTxParam().Receiver)
		toChain = common.GetChainFromChainID(send.GetCurrentOutTxParam().ReceiverChainId)
		if toChain == nil {
			logger.Error().Msgf("Unknown chain: %d", send.GetCurrentOutTxParam().ReceiverChainId)
			return
		}
	} else {
		logger.Info().Msgf("Transaction doesn't need to be processed status: %d", send.CctxStatus.Status)
		return
	}
	evmClient, ok := chainclient.(*EVMChainClient)
	if !ok {
		logger.Error().Msgf("chain client is not an EVMChainClient")
		return
	}

	// Early return if the cctx is already processed
	nonce := send.GetCurrentOutTxParam().OutboundTxTssNonce
	included, confirmed, err := evmClient.IsSendOutTxProcessed(send.Index, nonce, send.GetCurrentOutTxParam().CoinType, logger)
	if err != nil {
		logger.Error().Err(err).Msg("IsSendOutTxProcessed failed")
	}
	if included || confirmed {
		logger.Info().Msgf("CCTX already processed; exit signer")
		return
	}

	var message []byte
	if send.GetCurrentOutTxParam().CoinType != common.CoinType_Cmd {
		message, err = base64.StdEncoding.DecodeString(send.RelayedMessage)
		if err != nil {
			logger.Err(err).Msgf("decode CCTX.Message %s error", send.RelayedMessage)
		}
	}

	gasLimit := send.GetCurrentOutTxParam().OutboundTxGasLimit
	if gasLimit < 100_000 {
		gasLimit = 100_000
		logger.Warn().Msgf("gasLimit %d is too low; set to %d", send.GetCurrentOutTxParam().OutboundTxGasLimit, gasLimit)
	}
	if gasLimit > 1_000_000 {
		gasLimit = 1_000_000
		logger.Warn().Msgf("gasLimit %d is too high; set to %d", send.GetCurrentOutTxParam().OutboundTxGasLimit, gasLimit)
	}

	logger.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain, send.InboundTxParams.Amount, to.Hex(), nonce, send.InboundTxParams.InboundTxFinalizedZetaHeight)
	sendHash, err := hex.DecodeString(send.Index[2:]) // remove the leading 0x
	if err != nil || len(sendHash) != 32 {
		logger.Error().Err(err).Msgf("decode CCTX %s error", send.Index)
		return
	}
	var sendhash [32]byte
	copy(sendhash[:32], sendHash[:32])

	// use dynamic gas price for ethereum chains
	var gasprice *big.Int

	// The code below is a fix for https://github.com/zeta-chain/node/issues/1085
	// doesn't close directly the issue because we should determine if we want to keep using SuggestGasPrice if no OutboundTxGasPrice
	// we should possibly remove it completely and return an error if no OutboundTxGasPrice is provided because it means no fee is processed on ZetaChain
	specified, ok := new(big.Int).SetString(send.GetCurrentOutTxParam().OutboundTxGasPrice, 10)
	if !ok {
		if common.IsEthereumChain(toChain.ChainId) {
			suggested, err := signer.client.SuggestGasPrice(context.Background())
			if err != nil {
				logger.Error().Err(err).Msgf("cannot get gas price from chain %s ", toChain)
				return
			}
			gasprice = roundUpToNearestGwei(suggested)
		} else {
			logger.Error().Err(err).Msgf("cannot convert gas price  %s ", send.GetCurrentOutTxParam().OutboundTxGasPrice)
			return
		}
	} else {
		gasprice = specified
	}

	// In case there is a pending transaction, make sure this keysign is a transaction replacement
	pendingTx := evmClient.GetPendingTx(nonce)
	if pendingTx != nil {
		if gasprice.Cmp(pendingTx.GasPrice()) > 0 {
			logger.Info().Msgf("replace pending outTx %s nonce %d using gas price %d", pendingTx.Hash().Hex(), nonce, gasprice)
		} else {
			logger.Info().Msgf("please wait for pending outTx %s nonce %d to be included", pendingTx.Hash().Hex(), nonce)
			return
		}
	}

	flags, err := zetaBridge.GetCrosschainFlags()
	if err != nil {
		logger.Error().Err(err).Msgf("cannot get crosschain flags")
		return
	}

	var tx *ethtypes.Transaction

	if send.GetCurrentOutTxParam().CoinType == common.CoinType_Cmd { // admin command
		to := ethcommon.HexToAddress(send.GetCurrentOutTxParam().Receiver)
		if to == (ethcommon.Address{}) {
			logger.Error().Msgf("invalid receiver %s", send.GetCurrentOutTxParam().Receiver)
			return
		}

		msg := strings.Split(send.RelayedMessage, ":")
		if len(msg) != 2 {
			logger.Error().Msgf("invalid message %s", msg)
			return
		}
		tx, err = signer.SignCommandTx(msg[0], msg[1], to, send.GetCurrentOutTxParam(), gasLimit, gasprice, height)
	} else if send.InboundTxParams.SenderChainId == zetaBridge.ZetaChain().ChainId && send.CctxStatus.Status == types.CctxStatus_PendingOutbound && flags.IsOutboundEnabled {
		if send.GetCurrentOutTxParam().CoinType == common.CoinType_Gas {
			logger.Info().Msgf("SignWithdrawTx: %d => %s, nonce %d, gasprice %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, gasprice)
			tx, err = signer.SignWithdrawTx(
				to,
				send.GetCurrentOutTxParam().Amount.BigInt(),
				send.GetCurrentOutTxParam().OutboundTxTssNonce,
				gasprice,
				height,
			)
		}
		if send.GetCurrentOutTxParam().CoinType == common.CoinType_ERC20 {
			asset := ethcommon.HexToAddress(send.InboundTxParams.Asset)
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %s, nonce %d, gasprice %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, gasprice)
			tx, err = signer.SignERC20WithdrawTx(
				to,
				asset,
				send.GetCurrentOutTxParam().Amount.BigInt(),
				gasLimit,
				send.GetCurrentOutTxParam().OutboundTxTssNonce,
				gasprice,
				height,
			)
		}
		if send.GetCurrentOutTxParam().CoinType == common.CoinType_Zeta {
			logger.Info().Msgf("SignOutboundTx: %d => %s, nonce %d, gasprice %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, gasprice)
			tx, err = signer.SignOutboundTx(
				ethcommon.HexToAddress(send.InboundTxParams.Sender),
				big.NewInt(send.InboundTxParams.SenderChainId),
				to,
				send.GetCurrentOutTxParam().Amount.BigInt(),
				gasLimit,
				message,
				sendhash,
				send.GetCurrentOutTxParam().OutboundTxTssNonce,
				gasprice,
				height,
			)
		}
	} else if send.CctxStatus.Status == types.CctxStatus_PendingRevert && send.OutboundTxParams[0].ReceiverChainId == zetaBridge.ZetaChain().ChainId {
		if send.GetCurrentOutTxParam().CoinType == common.CoinType_Gas {
			logger.Info().Msgf("SignWithdrawTx: %d => %s, nonce %d, gasprice %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, gasprice)
			tx, err = signer.SignWithdrawTx(
				to,
				send.GetCurrentOutTxParam().Amount.BigInt(),
				send.GetCurrentOutTxParam().OutboundTxTssNonce,
				gasprice,
				height,
			)
		}
		if send.GetCurrentOutTxParam().CoinType == common.CoinType_ERC20 {
			asset := ethcommon.HexToAddress(send.InboundTxParams.Asset)
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %s, nonce %d, gasprice %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, gasprice)
			tx, err = signer.SignERC20WithdrawTx(
				to,
				asset,
				send.GetCurrentOutTxParam().Amount.BigInt(),
				gasLimit,
				send.GetCurrentOutTxParam().OutboundTxTssNonce,
				gasprice,
				height,
			)
		}
	} else if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		logger.Info().Msgf("SignRevertTx: %d => %s, nonce %d, gasprice %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, gasprice)
		tx, err = signer.SignRevertTx(
			ethcommon.HexToAddress(send.InboundTxParams.Sender),
			big.NewInt(send.OutboundTxParams[0].ReceiverChainId),
			to.Bytes(),
			big.NewInt(send.GetCurrentOutTxParam().ReceiverChainId),
			send.GetCurrentOutTxParam().Amount.BigInt(),
			gasLimit,
			message,
			sendhash,
			send.GetCurrentOutTxParam().OutboundTxTssNonce,
			gasprice,
			height,
		)
	} else if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf("SignOutboundTx: %d => %s, nonce %d, gasprice %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, gasprice)
		tx, err = signer.SignOutboundTx(
			ethcommon.HexToAddress(send.InboundTxParams.Sender),
			big.NewInt(send.InboundTxParams.SenderChainId),
			to,
			send.GetCurrentOutTxParam().Amount.BigInt(),
			gasLimit,
			message,
			sendhash,
			send.GetCurrentOutTxParam().OutboundTxTssNonce,
			gasprice,
			height,
		)
	}

	if err != nil {
		logger.Warn().Err(err).Msgf("signer SignOutbound error: nonce %d chain %d", send.GetCurrentOutTxParam().OutboundTxTssNonce, send.GetCurrentOutTxParam().ReceiverChainId)
		return
	}
	logger.Info().Msgf("Key-sign success: %d => %s, nonce %d", send.InboundTxParams.SenderChainId, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce)

	_, err = zetaBridge.GetObserverList()
	if err != nil {
		logger.Warn().Err(err).Msgf("unable to get observer list: chain %d observation %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, observertypes.ObservationType_OutBoundTx.String())

	}
	if tx != nil {
		outTxHash := tx.Hash().Hex()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", signer.chain, send.GetCurrentOutTxParam().OutboundTxTssNonce, outTxHash, myID)
		//if len(signers) == 0 || myid == signers[send.OutboundTxParams.Broadcaster] || myid == signers[int(send.OutboundTxParams.Broadcaster+1)%len(signers)] {
		backOff := 1000 * time.Millisecond
		// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
		for i := 0; i < 5; i++ {
			logger.Info().Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, toChain, send.GetCurrentOutTxParam().OutboundTxTssNonce, i)
			// #nosec G404 randomness is not a security issue here
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) // FIXME: use backoff
			err := signer.Broadcast(tx)
			if err != nil {
				log.Warn().Err(err).Msgf("OutTx Broadcast error")
				retry, report := HandleBroadcastError(err, strconv.FormatUint(send.GetCurrentOutTxParam().OutboundTxTssNonce, 10), toChain.String(), outTxHash)
				if report {
					signer.reportToOutTxTracker(zetaBridge, toChain.ChainId, tx.Nonce(), outTxHash, logger)
				}
				if !retry {
					break
				}
				backOff *= 2
				continue
			}
			logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, toChain, outTxHash)
			signer.reportToOutTxTracker(zetaBridge, toChain.ChainId, tx.Nonce(), outTxHash, logger)
			break // successful broadcast; no need to retry
		}
	}
}

// reportToOutTxTracker reports outTxHash to tracker only when tx receipt is available
func (signer *EVMSigner) reportToOutTxTracker(zetaBridge ZetaCoreBridger, chainID int64, nonce uint64, outTxHash string, logger zerolog.Logger) {
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

		// try fetching tx receipt for 10 minutes
		tStart := time.Now()
		for {
			if time.Since(tStart) > OutTxTrackerReportTimeout { // give up after 10 minutes
				logger.Info().Msgf("reportToOutTxTracker: outTxHash report timeout for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
				return
			}
			receipt, err := signer.client.TransactionReceipt(context.TODO(), ethcommon.HexToHash(outTxHash))
			if err != nil {
				logger.Info().Err(err).Msgf("reportToOutTxTracker: receipt not available for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
				time.Sleep(10 * time.Second)
				continue
			}
			if receipt != nil {
				_, isPending, err := signer.client.TransactionByHash(context.TODO(), ethcommon.HexToHash(outTxHash))
				if err != nil || isPending {
					logger.Info().Err(err).Msgf("reportToOutTxTracker: error getting tx or tx is pending for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
					time.Sleep(10 * time.Second)
					continue
				}
				break
			}
		}

		// report to outTxTracker
		zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(chainID, nonce, outTxHash, nil, "", -1)
		if err != nil {
			logger.Err(err).Msgf("reportToOutTxTracker: unable to add to tracker on ZetaCore for chain %d nonce %d outTxHash %s", chainID, nonce, outTxHash)
		}
		logger.Info().Msgf("reportToOutTxTracker: reported outTxHash to core successful %s, chain %d nonce %d outTxHash %s", zetaHash, chainID, nonce, outTxHash)
	}()
}

// SignERC20WithdrawTx
// function withdraw(
// address recipient,
// address asset,
// uint256 amount,
// ) external onlyTssAddress
func (signer *EVMSigner) SignERC20WithdrawTx(
	recipient ethcommon.Address,
	asset ethcommon.Address,
	amount *big.Int,
	gasLimit uint64,
	nonce uint64,
	gasPrice *big.Int,
	height uint64,
) (*ethtypes.Transaction, error) {
	var data []byte
	var err error
	data, err = signer.erc20CustodyABI.Pack("withdraw", recipient, asset, amount)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.erc20CustodyContractAddress, gasLimit, gasPrice, nonce, height)
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
func (signer *EVMSigner) SignWhitelistTx(
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
