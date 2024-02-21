package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outtxprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/tss"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const (
	maxNoOfInputsPerTx = 20
	consolidationRank  = 10           // the rank below (or equal to) which we consolidate UTXOs
	outTxBytesMin      = uint64(239)  // 239vB == EstimateSegWitTxSize(2, 3)
	outTxBytesMax      = uint64(1531) // 1531v == EstimateSegWitTxSize(21, 3)
)

type BTCSigner struct {
	tssSigner        interfaces.TSSSigner
	rpcClient        interfaces.BTCRPCClient
	logger           zerolog.Logger
	loggerCompliance zerolog.Logger
	ts               *metrics.TelemetryServer
}

var _ interfaces.ChainSigner = &BTCSigner{}

func NewBTCSigner(
	cfg config.BTCConfig,
	tssSigner interfaces.TSSSigner,
	loggers clientcommon.ClientLogger,
	ts *metrics.TelemetryServer) (*BTCSigner, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         cfg.RPCHost,
		User:         cfg.RPCUsername,
		Pass:         cfg.RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       cfg.RPCParams,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating bitcoin rpc client: %s", err)
	}

	return &BTCSigner{
		tssSigner:        tssSigner,
		rpcClient:        client,
		logger:           loggers.Std.With().Str("chain", "BTC").Str("module", "BTCSigner").Logger(),
		loggerCompliance: loggers.Compliance,
		ts:               ts,
	}, nil
}

// SignWithdrawTx receives utxos sorted by value, amount in BTC, feeRate in BTC per Kb
func (signer *BTCSigner) SignWithdrawTx(
	to *btcutil.AddressWitnessPubKeyHash,
	amount float64,
	gasPrice *big.Int,
	sizeLimit uint64,
	btcClient *BTCChainClient,
	height uint64,
	nonce uint64,
	chain *common.Chain,
) (*wire.MsgTx, error) {
	estimateFee := float64(gasPrice.Uint64()*outTxBytesMax) / 1e8
	nonceMark := common.NonceMarkAmount(nonce)

	// refresh unspent UTXOs and continue with keysign regardless of error
	err := btcClient.FetchUTXOS()
	if err != nil {
		signer.logger.Error().Err(err).Msgf("SignWithdrawTx: FetchUTXOS error: nonce %d chain %d", nonce, chain.ChainId)
	}

	// select N UTXOs to cover the total expense
	prevOuts, total, consolidatedUtxo, consolidatedValue, err := btcClient.SelectUTXOs(amount+estimateFee+float64(nonceMark)*1e-8, maxNoOfInputsPerTx, nonce, consolidationRank, false)
	if err != nil {
		return nil, err
	}

	// build tx with selected unspents
	tx := wire.NewMsgTx(wire.TxVersion)
	for _, prevOut := range prevOuts {
		hash, err := chainhash.NewHashFromStr(prevOut.TxID)
		if err != nil {
			return nil, err
		}
		outpoint := wire.NewOutPoint(hash, prevOut.Vout)
		txIn := wire.NewTxIn(outpoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	amountSatoshis, err := GetSatoshis(amount)
	if err != nil {
		return nil, err
	}

	// size checking
	// #nosec G701 always positive
	txSize := EstimateSegWitTxSize(uint64(len(prevOuts)), 3)
	if sizeLimit < BtcOutTxBytesWithdrawer { // ZRC20 'withdraw' charged less fee from end user
		signer.logger.Info().Msgf("sizeLimit %d is less than BtcOutTxBytesWithdrawer %d for nonce %d", sizeLimit, txSize, nonce)
	}
	if txSize < outTxBytesMin { // outbound shouldn't be blocked a low sizeLimit
		signer.logger.Warn().Msgf("txSize %d is less than outTxBytesMin %d; use outTxBytesMin", txSize, outTxBytesMin)
		txSize = outTxBytesMin
	}
	if txSize > outTxBytesMax { // in case of accident
		signer.logger.Warn().Msgf("txSize %d is greater than outTxBytesMax %d; use outTxBytesMax", txSize, outTxBytesMax)
		txSize = outTxBytesMax
	}

	// fee calculation
	// #nosec G701 always in range (checked above)
	fees := new(big.Int).Mul(big.NewInt(int64(txSize)), gasPrice)
	signer.logger.Info().Msgf("bitcoin outTx nonce %d gasPrice %s size %d fees %s consolidated %d utxos of value %v",
		nonce, gasPrice.String(), txSize, fees.String(), consolidatedUtxo, consolidatedValue)

	// calculate remaining btc to TSS self
	tssAddrWPKH := signer.tssSigner.BTCAddressWitnessPubkeyHash()
	payToSelf, err := PayToWitnessPubKeyHashScript(tssAddrWPKH.WitnessProgram())
	if err != nil {
		return nil, err
	}
	remaining := total - amount
	remainingSats, err := GetSatoshis(remaining)
	if err != nil {
		return nil, err
	}
	remainingSats -= fees.Int64()
	remainingSats -= nonceMark
	if remainingSats < 0 {
		return nil, fmt.Errorf("remainder value is negative: %d", remainingSats)
	} else if remainingSats == nonceMark {
		signer.logger.Info().Msgf("SignWithdrawTx: adjust remainder value to avoid duplicate nonce-mark: %d", remainingSats)
		remainingSats--
	}

	// 1st output: the nonce-mark btc to TSS self
	txOut1 := wire.NewTxOut(nonceMark, payToSelf)
	tx.AddTxOut(txOut1)

	// 2nd output: the payment to the recipient
	pkScript, err := PayToWitnessPubKeyHashScript(to.WitnessProgram())
	if err != nil {
		return nil, err
	}
	txOut2 := wire.NewTxOut(amountSatoshis, pkScript)
	tx.AddTxOut(txOut2)

	// 3rd output: the remaining btc to TSS self
	if remainingSats > 0 {
		txOut3 := wire.NewTxOut(remainingSats, payToSelf)
		tx.AddTxOut(txOut3)
	}

	// sign the tx
	sigHashes := txscript.NewTxSigHashes(tx)
	witnessHashes := make([][]byte, len(tx.TxIn))
	for ix := range tx.TxIn {
		amt, err := GetSatoshis(prevOuts[ix].Amount)
		if err != nil {
			return nil, err
		}
		pkScript, err := hex.DecodeString(prevOuts[ix].ScriptPubKey)
		if err != nil {
			return nil, err
		}
		witnessHashes[ix], err = txscript.CalcWitnessSigHash(pkScript, sigHashes, txscript.SigHashAll, tx, ix, amt)
		if err != nil {
			return nil, err
		}
	}
	tss, ok := signer.tssSigner.(*tss.TSS)
	if !ok {
		return nil, fmt.Errorf("tssSigner is not a TSS")
	}
	sig65Bs, err := tss.SignBatch(witnessHashes, height, nonce, chain)
	if err != nil {
		return nil, fmt.Errorf("SignBatch error: %v", err)
	}

	for ix := range tx.TxIn {
		sig65B := sig65Bs[ix]
		R := big.NewInt(0).SetBytes(sig65B[:32])
		S := big.NewInt(0).SetBytes(sig65B[32:64])
		sig := btcec.Signature{
			R: R,
			S: S,
		}

		pkCompressed := signer.tssSigner.PubKeyCompressedBytes()
		hashType := txscript.SigHashAll
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(hashType)), pkCompressed}
		tx.TxIn[ix].Witness = txWitness
	}
	return tx, nil
}

func (signer *BTCSigner) Broadcast(signedTx *wire.MsgTx) error {
	fmt.Printf("BTCSigner: Broadcasting: %s\n", signedTx.TxHash().String())

	var outBuff bytes.Buffer
	err := signedTx.Serialize(&outBuff)
	if err != nil {
		return err
	}
	str := hex.EncodeToString(outBuff.Bytes())
	fmt.Printf("BTCSigner: Transaction Data: %s\n", str)

	hash, err := signer.rpcClient.SendRawTransaction(signedTx, true)
	if err != nil {
		return err
	}
	signer.logger.Info().Msgf("Broadcasting BTC tx , hash %s ", hash)
	return nil
}

func (signer *BTCSigner) TryProcessOutTx(
	cctx *types.CrossChainTx,
	outTxMan *outtxprocessor.Processor,
	outTxID string,
	chainclient interfaces.ChainClient,
	zetaBridge interfaces.ZetaCoreBridger,
	height uint64,
) {
	defer func() {
		outTxMan.EndTryProcess(outTxID)
		if err := recover(); err != nil {
			signer.logger.Error().Msgf("BTC TryProcessOutTx: %s, caught panic error: %v", cctx.Index, err)
		}
	}()

	logger := signer.logger.With().
		Str("OutTxID", outTxID).
		Str("SendHash", cctx.Index).
		Logger()

	params := cctx.GetCurrentOutTxParam()
	if params.CoinType == common.CoinType_Zeta || params.CoinType == common.CoinType_ERC20 {
		logger.Error().Msgf("BTC TryProcessOutTx: can only send BTC to a BTC network")
		return
	}

	logger.Info().Msgf("BTC TryProcessOutTx: %s, value %d to %s", cctx.Index, params.Amount.BigInt(), params.Receiver)
	btcClient, ok := chainclient.(*BTCChainClient)
	if !ok {
		logger.Error().Msgf("chain client is not a bitcoin client")
		return
	}
	flags, err := zetaBridge.GetCrosschainFlags()
	if err != nil {
		logger.Error().Err(err).Msgf("cannot get crosschain flags")
		return
	}
	if !flags.IsOutboundEnabled {
		logger.Info().Msgf("outbound is disabled")
		return
	}
	myid := zetaBridge.GetKeys().GetAddress()
	outboundTxTssNonce := params.OutboundTxTssNonce

	sizelimit := params.OutboundTxGasLimit
	gasprice, ok := new(big.Int).SetString(params.OutboundTxGasPrice, 10)
	if !ok || gasprice.Cmp(big.NewInt(0)) < 0 {
		logger.Error().Msgf("cannot convert gas price  %s ", params.OutboundTxGasPrice)
		return
	}

	// Check receiver P2WPKH address
	bitcoinNetParams, err := common.BitcoinNetParamsFromChainID(params.ReceiverChainId)
	if err != nil {
		logger.Error().Err(err).Msgf("cannot get bitcoin net params%v", err)
		return
	}
	addr, err := common.DecodeBtcAddress(params.Receiver, params.ReceiverChainId)
	if err != nil {
		logger.Error().Err(err).Msgf("cannot decode address %s ", params.Receiver)
		return
	}
	if !addr.IsForNet(bitcoinNetParams) {
		logger.Error().Msgf(
			"address %s is not for network %s",
			params.Receiver,
			bitcoinNetParams.Name,
		)
		return
	}
	to, ok := addr.(*btcutil.AddressWitnessPubKeyHash)
	if err != nil || !ok {
		logger.Error().Err(err).Msgf("cannot convert address %s to P2WPKH address", params.Receiver)
		return
	}
	amount := float64(params.Amount.Uint64()) / 1e8

	// Add 1 satoshi/byte to gasPrice to avoid minRelayTxFee issue
	networkInfo, err := signer.rpcClient.GetNetworkInfo()
	if err != nil {
		logger.Error().Err(err).Msgf("cannot get bitcoin network info")
		return
	}
	satPerByte := FeeRateToSatPerByte(networkInfo.RelayFee)
	gasprice.Add(gasprice, satPerByte)

	// compliance check
	if clientcommon.IsCctxBanned(cctx) {
		logMsg := fmt.Sprintf("Banned address detected in cctx: sender %s receiver %s chain %d nonce %d",
			cctx.InboundTxParams.Sender, to, params.ReceiverChainId, outboundTxTssNonce)
		logger.Warn().Msg(logMsg)
		signer.loggerCompliance.Warn().Msg(logMsg)
		amount = 0 // zero out the amount to cancel the tx
	}

	logger.Info().Msgf("SignWithdrawTx: to %s, value %d sats", addr.EncodeAddress(), params.Amount.Uint64())
	logger.Info().Msgf("using utxos: %v", btcClient.utxos)

	tx, err := signer.SignWithdrawTx(
		to,
		amount,
		gasprice,
		sizelimit,
		btcClient,
		height,
		outboundTxTssNonce,
		&btcClient.chain,
	)
	if err != nil {
		logger.Warn().Err(err).Msgf("SignOutboundTx error: nonce %d chain %d", outboundTxTssNonce, params.ReceiverChainId)
		return
	}
	logger.Info().Msgf("Key-sign success: %d => %s, nonce %d", cctx.InboundTxParams.SenderChainId, btcClient.chain.ChainName, outboundTxTssNonce)

	// FIXME: add prometheus metrics
	_, err = zetaBridge.GetObserverList()
	if err != nil {
		logger.Warn().Err(err).Msgf("unable to get observer list: chain %d observation %s", outboundTxTssNonce, observertypes.ObservationType_OutBoundTx.String())
	}
	if tx != nil {
		outTxHash := tx.TxHash().String()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", btcClient.chain.ChainName, outboundTxTssNonce, outTxHash, myid)
		// TODO: pick a few broadcasters.
		//if len(signers) == 0 || myid == signers[send.OutboundTxParams.Broadcaster] || myid == signers[int(send.OutboundTxParams.Broadcaster+1)%len(signers)] {
		// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
		for i := 0; i < 5; i++ {
			// #nosec G404 randomness is not a security issue here
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) //random delay to avoid sychronized broadcast
			err := signer.Broadcast(tx)
			if err != nil {
				logger.Warn().Err(err).Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, btcClient.chain.ChainName, outboundTxTssNonce, i)
				continue
			}
			logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", outboundTxTssNonce, btcClient.chain.String(), outTxHash)
			zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(btcClient.chain.ChainId, outboundTxTssNonce, outTxHash, nil, "", -1)
			if err != nil {
				logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", outboundTxTssNonce, btcClient.chain.ChainName, outTxHash)
			}
			logger.Info().Msgf("Broadcast to core successful %s", zetaHash)

			// Save successfully broadcasted transaction to btc chain client
			btcClient.SaveBroadcastedTx(outTxHash, outboundTxTssNonce)

			break // successful broadcast; no need to retry
		}
	}
}
