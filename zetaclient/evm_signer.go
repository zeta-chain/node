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

type EVMSigner struct {
	client                      *ethclient.Client
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
}

var _ ChainSigner = &EVMSigner{}

func NewEVMSigner(chain common.Chain, endpoint string, tssSigner TSSSigner, abiString string, erc20CustodyABIString string, metaContract ethcommon.Address, erc20CustodyContract ethcommon.Address, logger zerolog.Logger, ts *TelemetryServer) (*EVMSigner, error) {
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
		ts: ts,
	}, nil
}

// given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *EVMSigner) Sign(data []byte, to ethcommon.Address, gasLimit uint64, gasPrice *big.Int, nonce uint64, height uint64) (*ethtypes.Transaction, []byte, []byte, error) {
	log.Debug().Msgf("TSS SIGNER: %s", signer.tssSigner.Pubkey())
	tx := ethtypes.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()

	sig, err := signer.tssSigner.Sign(hashBytes, height, signer.chain, "")
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

// takes in signed tx, broadcast to external chain node
func (signer *EVMSigner) Broadcast(tx *ethtypes.Transaction) error {
	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return signer.client.SendTransaction(ctxt, tx)
}

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

// function onRevert(
// address originSenderAddress,
// uint256 originChainId,
// bytes calldata destinationAddress,
// uint256 destinationChainId,
// uint256 zetaAmount,
// bytes calldata message,
// bytes32 internalSendHash
// ) external override whenNotPaused onlyTssAddress
func (signer *EVMSigner) SignRevertTx(sender ethcommon.Address, srcChainID *big.Int, to []byte, toChainID *big.Int, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
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
	sig, err := signer.tssSigner.Sign(hashBytes, height, signer.chain, "")
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

func (signer *EVMSigner) SignWithdrawTx(to ethcommon.Address, amount *big.Int, nonce uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes, height, signer.chain, "")
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

func (signer *EVMSigner) SignCommandTx(cmd string, params string, to ethcommon.Address, outboundParams *types.OutboundTxParams, gasLimit uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
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
		sig, err := signer.tssSigner.Sign(hashBytes, height, signer.chain, "")
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

func (signer *EVMSigner) TryProcessOutTx(send *types.CrossChainTx, outTxMan *OutTxProcessorManager, outTxID string, evmClient ChainClient, zetaBridge *ZetaCoreBridge, height uint64) {
	logger := signer.logger.With().
		Str("outTxID", outTxID).
		Str("SendHash", send.Index).
		Logger()
	logger.Info().Msgf("start processing outTxID %s", outTxID)
	logger.Info().Msgf("EVM Chain TryProcessOutTx: %s, value %d to %s", send.Index, send.GetCurrentOutTxParam().Amount.BigInt(), send.GetCurrentOutTxParam().Receiver)

	defer func() {
		outTxMan.EndTryProcess(outTxID)
	}()
	myid := zetaBridge.keys.GetOperatorAddress()

	var to ethcommon.Address
	var err error
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
	if err != nil {
		logger.Error().Err(err).Msg("ParseChain fail; skip")
		return
	}
	fmt.Println("Trying to process out tx")
	// Early return if the cctx is already processed
	included, confirmed, err := evmClient.IsSendOutTxProcessed(send.Index, send.GetCurrentOutTxParam().OutboundTxTssNonce, send.GetCurrentOutTxParam().CoinType, logger)
	if err != nil {
		logger.Error().Err(err).Msg("IsSendOutTxProcessed failed")
	}
	if included || confirmed {
		logger.Info().Msgf("CCTX already processed; exit signer")
		return
	}

	//message, err := base64.StdEncoding.DecodeString(send.RelayedMessage)
	//if err != nil {
	//	logger.Err(err).Msgf("decode CCTX.Message %s error", send.RelayedMessage)
	//}
	message, err := base64.StdEncoding.DecodeString("")
	if err != nil {
		logger.Err(err).Msgf("decode CCTX.Message %s error", send.RelayedMessage)
	}

	gasLimit := send.GetCurrentOutTxParam().OutboundTxGasLimit
	if gasLimit < 50_000 {
		gasLimit = 50_000
		logger.Warn().Msgf("gasLimit %d is too low; set to %d", send.GetCurrentOutTxParam().OutboundTxGasLimit, gasLimit)
	}
	if gasLimit > 1_000_000 {
		gasLimit = 1_000_000
		logger.Warn().Msgf("gasLimit %d is too high; set to %d", send.GetCurrentOutTxParam().OutboundTxGasLimit, gasLimit)
	}

	logger.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain, send.InboundTxParams.Amount, to.Hex(), send.GetCurrentOutTxParam().OutboundTxTssNonce, send.InboundTxParams.InboundTxFinalizedZetaHeight)
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
	//if common.IsEthereumChain(toChain.ChainId) {
	//	suggested, err := signer.client.SuggestGasPrice(context.Background())
	//	if err != nil {
	//		logger.Error().Err(err).Msgf("cannot get gas price from chain %s ", toChain)
	//		return
	//	}
	//	gasprice = roundUpToNearestGwei(suggested)
	//} else {
	//	specified, ok := new(big.Int).SetString(send.GetCurrentOutTxParam().OutboundTxGasPrice, 10)
	//	if !ok {
	//		logger.Error().Err(err).Msgf("cannot convert gas price  %s ", send.GetCurrentOutTxParam().OutboundTxGasPrice)
	//		return
	//	}
	//	gasprice = specified
	//}

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
	} else if send.InboundTxParams.SenderChainId == common.ZetaChain().ChainId && send.CctxStatus.Status == types.CctxStatus_PendingOutbound && flags.IsOutboundEnabled {
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
	} else if send.CctxStatus.Status == types.CctxStatus_PendingRevert && send.OutboundTxParams[0].ReceiverChainId == common.ZetaChain().ChainId {
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

	_, err = zetaBridge.GetObserverList(*toChain)
	if err != nil {
		logger.Warn().Err(err).Msgf("unable to get observer list: chain %d observation %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, observertypes.ObservationType_OutBoundTx.String())

	}
	if tx != nil {
		outTxHash := tx.Hash().Hex()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", signer.chain, send.GetCurrentOutTxParam().OutboundTxTssNonce, outTxHash, myid)
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
					zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(toChain.ChainId, tx.Nonce(), outTxHash, nil, "", -1)
					if err != nil {
						logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, toChain, outTxHash)
					}
					logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
				}
				if !retry {
					break
				}
				backOff *= 2
				continue
			}
			logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, toChain, outTxHash)
			zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(toChain.ChainId, tx.Nonce(), outTxHash, nil, "", -1)
			if err != nil {
				logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, toChain, outTxHash)
			}
			logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
			break // successful broadcast; no need to retry
		}

	}
	//}

}

// function withdraw(
// address recipient,
// address asset,
// uint256 amount,
// ) external onlyTssAddress
func (signer *EVMSigner) SignERC20WithdrawTx(recipient ethcommon.Address, asset ethcommon.Address, amount *big.Int, gasLimit uint64, nonce uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
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

// function whitelist(
// address asset,
// ) external onlyTssAddress
// function unwhitelist(
// address asset,
// ) external onlyTssAddress
func (signer *EVMSigner) SignWhitelistTx(action string, _ ethcommon.Address, asset ethcommon.Address, gasLimit uint64, nonce uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
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
