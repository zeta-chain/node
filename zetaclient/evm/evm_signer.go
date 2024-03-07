package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outtxprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
	zbridge "github.com/zeta-chain/zetacore/zetaclient/zetabridge"
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
	logger                      clientcommon.ClientLogger
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
	loggers clientcommon.ClientLogger,
	ts *metrics.TelemetryServer,
) (*Signer, error) {
	client, chainID, ethSigner, err := getEVMRPC(endpoint)
	if err != nil {
		return nil, err
	}
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
		logger: clientcommon.ClientLogger{
			Std:        loggers.Std.With().Str("chain", chain.ChainName.String()).Str("module", "EVMSigner").Logger(),
			Compliance: loggers.Compliance,
		},
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
		signer.logger.Std.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Std.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
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
func (signer *Signer) SignOutboundTx(txData *OutBoundTransactionData) (*ethtypes.Transaction, error) {
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
		return nil, fmt.Errorf("onReceive pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(data,
		signer.metaContractAddress,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height)
	if err != nil {
		return nil, fmt.Errorf("onReceive sign error: %w", err)
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
func (signer *Signer) SignRevertTx(txData *OutBoundTransactionData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	data, err = signer.abi.Pack("onRevert",
		txData.sender,
		txData.srcChainID,
		txData.to.Bytes(),
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

// SignCancelTx signs a transaction from TSS address to itself with a zero amount in order to increment the nonce
func (signer *Signer) SignCancelTx(nonce uint64, gasPrice *big.Int, height uint64) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(nonce, signer.tssSigner.EVMAddress(), big.NewInt(0), 21000, gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes, height, nonce, signer.chain, "")
	if err != nil {
		return nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Std.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Std.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, err
	}

	return signedTX, nil
}

// SignWithdrawTx signs a withdrawal transaction sent from the TSS address to the destination
func (signer *Signer) SignWithdrawTx(txData *OutBoundTransactionData) (*ethtypes.Transaction, error) {
	tx := ethtypes.NewTransaction(txData.nonce, txData.to, txData.amount, 21000, txData.gasPrice, nil)
	hashBytes := signer.ethSigner.Hash(tx).Bytes()
	sig, err := signer.tssSigner.Sign(hashBytes, txData.height, txData.nonce, signer.chain, "")
	if err != nil {
		return nil, err
	}
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.logger.Std.Error().Err(err).Msgf("SigToPub error")
	}
	addr := crypto.PubkeyToAddress(*pubk)
	signer.logger.Std.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, err
	}

	return signedTX, nil
}

// SignCommandTx signs a transaction based on the given command includes:
//
//	cmd_whitelist_erc20
//	cmd_migrate_tss_funds
func (signer *Signer) SignCommandTx(txData *OutBoundTransactionData, cmd string, params string) (*ethtypes.Transaction, error) {
	switch cmd {
	case common.CmdWhitelistERC20:
		return signer.SignWhitelistERC20Cmd(txData, params)
	case common.CmdMigrateTssFunds:
		return signer.SignMigrateTssFundsCmd(txData)
	}
	return nil, fmt.Errorf("SignCommandTx: unknown command %s", cmd)
}

// TryProcessOutTx - signer interface implementation
// This function will attempt to build and sign an evm transaction using the TSS signer.
// It will then broadcast the signed transaction to the outbound chain.
func (signer *Signer) TryProcessOutTx(
	cctx *types.CrossChainTx,
	outTxMan *outtxprocessor.Processor,
	outTxID string,
	chainclient interfaces.ChainClient,
	zetaBridge interfaces.ZetaCoreBridger,
	height uint64,
) {
	logger := signer.logger.Std.With().
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
	txData, skipTx, err := NewOutBoundTransactionData(cctx, evmClient, signer.client, logger, height)
	if err != nil {
		logger.Err(err).Msg("error setting up transaction input fields")
		return
	}
	if skipTx {
		return
	}

	// Get destination chain for logging
	toChain := common.GetChainFromChainID(txData.toChainID.Int64())

	// Get cross-chain flags
	crossChainflags, err := zetaBridge.GetCrosschainFlags()
	if err != nil {
		logger.Err(err).Msg("couldn't retrieve crosschain flags from core")
		return
	}

	var tx *ethtypes.Transaction
	// compliance check goes first
	if clientcommon.IsCctxRestricted(cctx) {
		clientcommon.PrintComplianceLog(logger, signer.logger.Compliance,
			true, evmClient.chain.ChainId, cctx.Index, cctx.InboundTxParams.Sender, txData.to.Hex(), cctx.GetCurrentOutTxParam().CoinType.String())
		tx, err = signer.SignCancelTx(txData.nonce, txData.gasPrice, height) // cancel the tx
		if err != nil {
			logger.Warn().Err(err).Msg(SignerErrorMsg(cctx))
			return
		}
	} else if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Cmd { // admin command
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
		// cmd field is used to determine whether to execute ERC20 whitelist or migrate TSS funds given that the coin type
		// from the cctx is common.CoinType_Cmd
		cmd := msg[0]
		// params field is used to pass input parameters for command requests, currently it is used to pass the ERC20
		// contract address when a whitelist command is requested
		params := msg[1]
		tx, err = signer.SignCommandTx(txData, cmd, params)
		if err != nil {
			logger.Warn().Err(err).Msg(SignerErrorMsg(cctx))
			return
		}
	} else if IsSenderZetaChain(cctx, zetaBridge, &crossChainflags) {
		switch cctx.GetCurrentOutTxParam().CoinType {
		case common.CoinType_Gas:
			logger.Info().Msgf("SignWithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignWithdrawTx(txData)
		case common.CoinType_ERC20:
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignERC20WithdrawTx(txData)
		case common.CoinType_Zeta:
			logger.Info().Msgf("SignOutboundTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignOutboundTx(txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(SignerErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.OutboundTxParams[0].ReceiverChainId == zetaBridge.ZetaChain().ChainId {
		switch cctx.GetCurrentOutTxParam().CoinType {
		case common.CoinType_Gas:
			logger.Info().Msgf("SignWithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignWithdrawTx(txData)
		case common.CoinType_ERC20:
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
			tx, err = signer.SignERC20WithdrawTx(txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(SignerErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		logger.Info().Msgf("SignRevertTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
		txData.srcChainID = big.NewInt(cctx.OutboundTxParams[0].ReceiverChainId)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutTxParam().ReceiverChainId)

		tx, err = signer.SignRevertTx(txData)
		if err != nil {
			logger.Warn().Err(err).Msg(SignerErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf("SignOutboundTx: %d => %s, nonce %d, gasPrice %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, txData.gasPrice)
		tx, err = signer.SignOutboundTx(txData)
		if err != nil {
			logger.Warn().Err(err).Msg(SignerErrorMsg(cctx))
			return
		}
	}

	logger.Info().Msgf("Key-sign success: %d => %s, nonce %d", cctx.InboundTxParams.SenderChainId, toChain, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)

	// Broadcast Signed Tx
	signer.BroadcastOutTx(tx, cctx, logger, myID, zetaBridge, txData)
}

// BroadcastOutTx signed transaction through evm rpc client
func (signer *Signer) BroadcastOutTx(
	tx *ethtypes.Transaction,
	cctx *types.CrossChainTx,
	logger zerolog.Logger,
	myID sdk.AccAddress,
	zetaBridge interfaces.ZetaCoreBridger,
	txData *OutBoundTransactionData) {
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
func (signer *Signer) SignERC20WithdrawTx(txData *OutBoundTransactionData) (*ethtypes.Transaction, error) {
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

// Exported for unit tests

func (signer *Signer) GetReportedTxList() *map[string]bool {
	return &signer.outTxHashBeingReported
}
func (signer *Signer) EvmClient() interfaces.EVMRPCClient {
	return signer.client
}
func (signer *Signer) EvmSigner() ethtypes.Signer {
	return signer.ethSigner
}

// ________________________

// getEVMRPC is a helper function to set up the client and signer, also initializes a mock client for unit tests
func getEVMRPC(endpoint string) (interfaces.EVMRPCClient, *big.Int, ethtypes.Signer, error) {
	if endpoint == stub.EVMRPCEnabled {
		chainID := big.NewInt(common.BscMainnetChain().ChainId)
		ethSigner := ethtypes.NewEIP155Signer(chainID)
		client := stub.EvmClient{}
		return client, chainID, ethSigner, nil
	}

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, nil, nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, nil, nil, err
	}
	ethSigner := ethtypes.LatestSignerForChainID(chainID)
	return client, chainID, ethSigner, nil
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

func IsSenderZetaChain(cctx *types.CrossChainTx, zetaBridge interfaces.ZetaCoreBridger, flags *observertypes.CrosschainFlags) bool {
	return cctx.InboundTxParams.SenderChainId == zetaBridge.ZetaChain().ChainId && cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound && flags.IsOutboundEnabled
}

func SignerErrorMsg(cctx *types.CrossChainTx) string {
	return fmt.Sprintf("signer SignOutbound error: nonce %d chain %d", cctx.GetCurrentOutTxParam().OutboundTxTssNonce, cctx.GetCurrentOutTxParam().ReceiverChainId)
}

func (signer *Signer) SignWhitelistERC20Cmd(txData *OutBoundTransactionData, params string) (*ethtypes.Transaction, error) {
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
		return nil, err
	}
	tx, _, _, err := signer.Sign(data, txData.to, txData.gasLimit, txData.gasPrice, outboundParams.OutboundTxTssNonce, txData.height)
	if err != nil {
		return nil, fmt.Errorf("sign error: %w", err)
	}
	return tx, nil
}

func (signer *Signer) SignMigrateTssFundsCmd(txData *OutBoundTransactionData) (*ethtypes.Transaction, error) {
	outboundParams := txData.outboundParams
	tx, _, _, err := signer.Sign(nil, txData.to, txData.gasLimit, txData.gasPrice, outboundParams.OutboundTxTssNonce, txData.height)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
