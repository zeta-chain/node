package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverModuleTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type EVMSigner struct {
	client              *ethclient.Client
	chain               common.Chain
	chainID             *big.Int
	tssSigner           TSSSigner
	ethSigner           ethtypes.Signer
	abi                 abi.ABI
	metaContractAddress ethcommon.Address
	logger              zerolog.Logger
}

func NewEVMSigner(chain common.Chain, endpoint string, tssSigner TSSSigner, abiString string, metaContract ethcommon.Address) (*EVMSigner, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}
	ethSigner := ethtypes.LatestSignerForChainID(chainID)
	abi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}

	return &EVMSigner{
		client:              client,
		chain:               chain,
		tssSigner:           tssSigner,
		chainID:             chainID,
		ethSigner:           ethSigner,
		abi:                 abi,
		metaContractAddress: metaContract,
		logger:              log.With().Str("module", "EVMSigner").Logger(),
	}, nil
}

// given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *EVMSigner) Sign(data []byte, to ethcommon.Address, gasLimit uint64, gasPrice *big.Int, nonce uint64) (*ethtypes.Transaction, []byte, []byte, error) {
	tx := ethtypes.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes)
	if err != nil {
		return nil, nil, nil, err
	}
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
func (signer *EVMSigner) SignOutboundTx(sender ethcommon.Address, srcChainID *big.Int, to ethcommon.Address, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	if len(sendHash) < 32 {
		return nil, fmt.Errorf("sendHash len %d must be 32", len(sendHash))
	}
	var data []byte
	var err error

	data, err = signer.abi.Pack("onReceive", sender.Bytes(), srcChainID, to, amount, message, sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce)
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
func (signer *EVMSigner) SignRevertTx(sender ethcommon.Address, srcChainID *big.Int, to []byte, toChainID *big.Int, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	data, err = signer.abi.Pack("onRevert", sender, srcChainID, to, toChainID, amount, message, sendHash)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data, signer.metaContractAddress, gasLimit, gasPrice, nonce)
	if err != nil {
		return nil, fmt.Errorf("Sign error: %w", err)
	}

	return tx, nil
}

func (signer *EVMSigner) SignCancelTx(nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, signer.tssSigner.EVMAddress(), big.NewInt(0), 21000, gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes)
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

func (signer *EVMSigner) SignWithdrawTx(to ethcommon.Address, amount *big.Int, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes)
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

func (signer *EVMSigner) TryProcessOutTx(send *types.CrossChainTx, outTxMan *OutTxProcessorManager, evmClient ChainClient, zetaBridge *ZetaCoreBridge) {
	chain := GetTargetChain(send)
	outTxID := fmt.Sprintf("%s/%d", chain, send.OutBoundTxParams.OutBoundTxTSSNonce)

	logger := signer.logger.With().
		Str("sendHash", send.Index).
		Str("outTxID", outTxID).
		Logger()
	logger.Info().Msgf("start processing outTxID %s", outTxID)
	defer func() {
		outTxMan.EndTryProcess(outTxID)
	}()
	myid := zetaBridge.keys.GetSignerInfo().GetAddress().String()

	var to ethcommon.Address
	var err error
	var toChain common.Chain
	if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		to = ethcommon.HexToAddress(send.InBoundTxParams.Sender)
		toChain, err = common.ParseChain(send.InBoundTxParams.SenderChain)
		logger.Info().Msgf("Abort: reverting inbound")
	} else if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		to = ethcommon.HexToAddress(send.OutBoundTxParams.Receiver)
		toChain, err = common.ParseChain(send.OutBoundTxParams.ReceiverChain)
	}
	if err != nil {
		logger.Error().Err(err).Msg("ParseChain fail; skip")
		return
	}

	// Early return if the send is already processed
	included, confirmed, _ := evmClient.IsSendOutTxProcessed(send)
	if included || confirmed {
		logger.Info().Msgf("CCTX already processed; exit signer")
		return
	}

	message, err := base64.StdEncoding.DecodeString(send.RelayedMessage)
	if err != nil {
		logger.Err(err).Msgf("decode CCTX.Message %s error", send.RelayedMessage)
	}

	gasLimit := send.OutBoundTxParams.OutBoundTxGasLimit
	if gasLimit < 50_000 {
		gasLimit = 50_000
		logger.Warn().Msgf("gasLimit %d is too low; set to %d", send.OutBoundTxParams.OutBoundTxGasLimit, gasLimit)
	}
	if gasLimit > 1_000_000 {
		gasLimit = 1_000_000
		logger.Warn().Msgf("gasLimit %d is too high; set to %d", send.OutBoundTxParams.OutBoundTxGasLimit, gasLimit)
	}

	logger.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain, send.ZetaMint, to.Hex(), send.OutBoundTxParams.OutBoundTxTSSNonce, send.InBoundTxParams.InBoundTxFinalizedZetaHeight)
	sendHash, err := hex.DecodeString(send.Index[2:]) // remove the leading 0x
	if err != nil || len(sendHash) != 32 {
		logger.Error().Err(err).Msgf("decode CCTX %s error", send.Index)
		return
	}
	var sendhash [32]byte
	copy(sendhash[:32], sendHash[:32])
	gasprice, ok := new(big.Int).SetString(send.OutBoundTxParams.OutBoundTxGasPrice, 10)
	if !ok {
		logger.Error().Err(err).Msgf("cannot convert gas price  %s ", send.OutBoundTxParams.OutBoundTxGasPrice)
		return
	}
	// FIXME: remove this hack
	if toChain == common.GoerliChain {
		gasprice = gasprice.Mul(gasprice, big.NewInt(3))
		gasprice = gasprice.Div(gasprice, big.NewInt(2))
	}

	var tx *ethtypes.Transaction
	if send.InBoundTxParams.SenderChain == "ZETA" && send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf("SignWithdrawTx: %s => %s, nonce %d, gasprice %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
		tx, err = signer.SignWithdrawTx(to, send.ZetaMint.BigInt(), send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
	} else if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		srcChainID := config.Chains[send.InBoundTxParams.SenderChain].ChainID
		logger.Info().Msgf("SignRevertTx: %s => %s, nonce %d, gasprice %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
		toChainID := config.Chains[send.OutBoundTxParams.ReceiverChain].ChainID
		tx, err = signer.SignRevertTx(ethcommon.HexToAddress(send.InBoundTxParams.Sender), srcChainID, to.Bytes(), toChainID, send.ZetaMint.BigInt(), gasLimit, message, sendhash, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
	} else if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		srcChainID := config.Chains[send.InBoundTxParams.SenderChain].ChainID
		logger.Info().Msgf("SignOutboundTx: %s => %s, nonce %d, gasprice %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
		tx, err = signer.SignOutboundTx(ethcommon.HexToAddress(send.InBoundTxParams.Sender), srcChainID, to, send.ZetaMint.BigInt(), gasLimit, message, sendhash, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
	}

	if err != nil {
		logger.Warn().Err(err).Msgf("SignOutboundTx error: nonce %d chain %s", send.OutBoundTxParams.OutBoundTxTSSNonce, send.OutBoundTxParams.ReceiverChain)
		return
	}
	logger.Info().Msgf("Key-sign success: %s => %s, nonce %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce)
	//cnt, err := co.GetPromCounter(OutboundTxSignCount)
	//if err != nil {
	//	log.Error().Err(err).Msgf("GetPromCounter error")
	//} else {
	//	cnt.Inc()
	//}
	signers, err := zetaBridge.GetObserverList(toChain, zetaObserverModuleTypes.ObservationType_OutBoundTx.String())
	if err != nil {
		logger.Warn().Err(err).Msgf("unable to get observer list: chain %d observation %s", send.OutBoundTxParams.OutBoundTxTSSNonce, zetaObserverModuleTypes.ObservationType_OutBoundTx.String())

	}
	if tx != nil {
		outTxHash := tx.Hash().Hex()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", signer.chain, send.OutBoundTxParams.OutBoundTxTSSNonce, outTxHash, myid)
		if len(signers) == 0 || myid == signers[send.OutBoundTxParams.Broadcaster] || myid == signers[int(send.OutBoundTxParams.Broadcaster+1)%len(signers)] {
			backOff := 1000 * time.Millisecond
			// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
			for i := 0; i < 5; i++ {
				logger.Info().Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, i)
				// #nosec G404 randomness is not a security issue here
				time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) // FIXME: use backoff
				err := signer.Broadcast(tx)
				if err != nil {
					log.Warn().Err(err).Msgf("OutTx Broadcast error")
					retry, report := HandleBroadcastError(err, strconv.FormatUint(send.OutBoundTxParams.OutBoundTxTSSNonce, 10), toChain.String(), outTxHash)
					if report {
						zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(toChain.String(), tx.Nonce(), outTxHash)
						if err != nil {
							logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.OutBoundTxParams.OutBoundTxTSSNonce, toChain, outTxHash)
						}
						logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
					}
					if !retry {
						break
					}
					backOff *= 2
					continue
				}
				logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", send.OutBoundTxParams.OutBoundTxTSSNonce, toChain, outTxHash)
				zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(toChain.String(), tx.Nonce(), outTxHash)
				if err != nil {
					logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.OutBoundTxParams.OutBoundTxTSSNonce, toChain, outTxHash)
				}
				logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
				break // successful broadcast; no need to retry
			}

		}
	}

}
