package zetaclient

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"gorm.io/gorm"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverModuleTypes "github.com/zeta-chain/zetacore/x/observer/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	maxNoOfInputsPerTx = 20
)

type BTCSigner struct {
	tssSigner TSSSigner
	rpcClient *rpcclient.Client
	logger    zerolog.Logger
	ts        *TelemetryServer
}

var _ ChainSigner = &BTCSigner{}

func NewBTCSigner(tssSigner TSSSigner, rpcClient *rpcclient.Client, logger zerolog.Logger, ts *TelemetryServer) (*BTCSigner, error) {
	return &BTCSigner{
		tssSigner: tssSigner,
		rpcClient: rpcClient,
		logger: logger.With().
			Str("chain", "BTC").
			Str("module", "BTCSigner").Logger(),
		ts: ts,
	}, nil
}

// SignWithdrawTx receives utxos sorted by value, amount in BTC, feeRate in BTC per Kb
func (signer *BTCSigner) SignWithdrawTx(to *btcutil.AddressWitnessPubKeyHash, amount float64, gasPrice *big.Int, utxos []btcjson.ListUnspentResult, db *gorm.DB, height uint64) (*wire.MsgTx, error) {
	// select N UTXOs to cover the amount
	estimateFee := 0.0001 // 10,000 sats, should be good for testnet
	minFee := 0.00005
	prevOuts, total, err := selectUTXOs(utxos, amount+estimateFee, maxNoOfInputsPerTx)
	if err != nil {
		return nil, err
	}
	remaining := total - amount

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

	amountSatoshis, err := getSatoshis(amount)
	if err != nil {
		return nil, err
	}
	// add txout with remaining btc
	fees := new(big.Int).Mul(big.NewInt(int64(tx.SerializeSize())), gasPrice)
	fees.Div(fees, big.NewInt(1000)) //FIXME: feeRate KB is 1000B or 1024B?
	if fees.Int64() < int64(minFee*1e8) {
		fmt.Printf("fees %d is less than minFee %f; use minFee", fees, minFee*1e8)
		fees = big.NewInt(int64(minFee * 1e8))
	}

	tssAddrWPKH := signer.tssSigner.BTCAddressWitnessPubkeyHash()
	pkScript2, err := payToWitnessPubKeyHashScript(tssAddrWPKH.WitnessProgram())
	if err != nil {
		return nil, err
	}
	remainingSatoshis, err := getSatoshis(remaining)
	if err != nil {
		return nil, err
	}
	txOut := wire.NewTxOut(remainingSatoshis, pkScript2)

	remainderValue := remainingSatoshis - fees.Int64()
	if remainderValue < 0 {
		fmt.Printf("BTCSigner: SignWithdrawTx: Remainder Value is negative! : %d\n", remainderValue)
		fmt.Printf("BTCSigner: SignWithdrawTx: Number of inputs : %d\n", len(tx.TxIn))
		return nil, fmt.Errorf("remainder value is negative")
	}

	txOut.Value = remainderValue
	tx.AddTxOut(txOut)

	// add txout
	pkScript, err := payToWitnessPubKeyHashScript(to.WitnessProgram())
	if err != nil {
		return nil, err
	}
	txOut2 := wire.NewTxOut(amountSatoshis, pkScript)
	tx.AddTxOut(txOut2)

	// sign the tx
	sigHashes := txscript.NewTxSigHashes(tx)
	witnessHashes := make([][]byte, len(tx.TxIn))
	for ix := range tx.TxIn {
		amt, err := getSatoshis(prevOuts[ix].Amount)
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
	tss, ok := signer.tssSigner.(*TSS)
	if !ok {
		return nil, fmt.Errorf("tssSigner is not a TSS")
	}
	sig65Bs, err := tss.SignBatch(witnessHashes, height)
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
		if err != nil {
			return nil, err
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
	_ = signedTx.Serialize(&outBuff)
	str := hex.EncodeToString(outBuff.Bytes())
	fmt.Printf("BTCSigner: Transaction Data: %s\n", str)

	hash, err := signer.rpcClient.SendRawTransaction(signedTx, true)
	if err != nil {
		return err
	}
	signer.logger.Info().Msgf("Broadcasting BTC tx , hash %s ", hash)
	return err
}

func (signer *BTCSigner) TryProcessOutTx(send *types.CrossChainTx, outTxMan *OutTxProcessorManager, outTxID string, chainclient ChainClient, zetaBridge *ZetaCoreBridge, height uint64) {
	defer func() {
		if err := recover(); err != nil {
			signer.logger.Error().Msgf("BTC TryProcessOutTx: %s, caught panic error: %v", send.Index, err)
		}
	}()

	logger := signer.logger.With().
		Str("OutTxID", outTxID).
		Str("SendHash", send.Index).
		Logger()
	if send.GetCurrentOutTxParam().CoinType != common.CoinType_Gas {
		logger.Error().Msgf("BTC TryProcessOutTx: can only send BTC to a BTC network")
		return
	}
	toAddr := send.GetCurrentOutTxParam().Receiver

	logger.Info().Msgf("BTC TryProcessOutTx: %s, value %d to %s", send.Index, send.GetCurrentOutTxParam().Amount.BigInt(), toAddr)
	defer func() {
		outTxMan.EndTryProcess(outTxID)
	}()
	btcClient, ok := chainclient.(*BitcoinChainClient)
	if !ok {
		logger.Error().Msgf("chain client is not a bitcoin client")
		return
	}

	myid := zetaBridge.keys.GetAddress()
	// Early return if the send is already processed
	// FIXME: handle revert case
	included, confirmed, _ := btcClient.IsSendOutTxProcessed(send.Index, int(send.GetCurrentOutTxParam().OutboundTxTssNonce), common.CoinType_Gas, logger)
	if included || confirmed {
		logger.Info().Msgf("CCTX already processed; exit signer")
		return
	}

	gasprice, ok := new(big.Int).SetString(send.GetCurrentOutTxParam().OutboundTxGasPrice, 10)
	if !ok {
		logger.Error().Msgf("cannot convert gas price  %s ", send.GetCurrentOutTxParam().OutboundTxGasPrice)
		return
	}

	// FIXME: config chain params
	addr, err := btcutil.DecodeAddress(toAddr, config.BitconNetParams)
	if err != nil {
		logger.Error().Err(err).Msgf("cannot decode address %s ", send.GetCurrentOutTxParam().Receiver)
		return
	}
	to, ok := addr.(*btcutil.AddressWitnessPubKeyHash)
	if err != nil || !ok {
		logger.Error().Err(err).Msgf("cannot decode address %s ", send.GetCurrentOutTxParam().Receiver)
		return
	}

	logger.Info().Msgf("SignWithdrawTx: to %s, value %d sats", addr.EncodeAddress(), send.GetCurrentOutTxParam().Amount.Uint64())
	logger.Info().Msgf("using utxos: %v", btcClient.utxos)
	tx, err := signer.SignWithdrawTx(to, float64(send.GetCurrentOutTxParam().Amount.Uint64())/1e8, gasprice, btcClient.utxos, btcClient.db, height)
	if err != nil {
		logger.Warn().Err(err).Msgf("SignOutboundTx error: nonce %d chain %d", send.GetCurrentOutTxParam().OutboundTxTssNonce, send.GetCurrentOutTxParam().ReceiverChainId)
		return
	}
	logger.Info().Msgf("Key-sign success: %d => %s, nonce %d", send.InboundTxParams.SenderChainId, btcClient.chain.ChainName, send.GetCurrentOutTxParam().OutboundTxTssNonce)
	// FIXME: add prometheus metrics
	_, err = zetaBridge.GetObserverList(btcClient.chain)
	if err != nil {
		logger.Warn().Err(err).Msgf("unable to get observer list: chain %d observation %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, zetaObserverModuleTypes.ObservationType_OutBoundTx.String())
	}
	if tx != nil {
		outTxHash := tx.TxHash().String()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", btcClient.chain.ChainName, send.GetCurrentOutTxParam().OutboundTxTssNonce, outTxHash, myid)
		// TODO: pick a few broadcasters.
		//if len(signers) == 0 || myid == signers[send.OutboundTxParams.Broadcaster] || myid == signers[int(send.OutboundTxParams.Broadcaster+1)%len(signers)] {
		// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
		for i := 0; i < 5; i++ {
			// #nosec G404 randomness is not a security issue here
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) //random delay to avoid sychronized broadcast
			err := signer.Broadcast(tx)
			if err != nil {
				logger.Warn().Err(err).Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, btcClient.chain.ChainName, send.GetCurrentOutTxParam().OutboundTxTssNonce, i)
				continue
			}
			logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, btcClient.chain.String(), outTxHash)
			zetaHash, err := zetaBridge.AddTxHashToOutTxTracker(btcClient.chain.ChainId, send.GetCurrentOutTxParam().OutboundTxTssNonce, outTxHash)
			if err != nil {
				logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.GetCurrentOutTxParam().OutboundTxTssNonce, btcClient.chain.ChainName, outTxHash)
			}
			logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
			outTxID := btcClient.GetTxID(send.GetCurrentOutTxParam().OutboundTxTssNonce)

			// Save successfully broadcasted transaction to btc chain client
			btcClient.mu.Lock()
			btcClient.broadcastedTx[outTxID] = tx.TxHash()
			btcClient.mu.Unlock()
			broadcastEntry := clienttypes.ToTransactionHashSQLType(tx.TxHash(), outTxID)
			if err := btcClient.db.Create(&broadcastEntry).Error; err != nil {
				btcClient.logger.ObserveOutTx.Error().Err(err).Msg("observeOutTx: error saving broadcasted tx")
			}

			break // successful broadcast; no need to retry
		}
	}
}

// Selects a sublist of utxos to be used as inputs.
//
// Parameters:
//   - utxos: A list of ordered (by amount in ascending order) UTXO
//   - amount: The desired minimum total value of the selected UTXOs.
//   - utxoCap: The maximum number of UTXOs to be selected.
//
// Returns: a sublist of selected UTXOs or an error if the qulifying sublist cannot be found.
func selectUTXOs(utxos []btcjson.ListUnspentResult, amount float64, utxoCap uint8) ([]btcjson.ListUnspentResult, float64, error) {
	total := 0.0
	left, right := 0, 0

	for total < amount && right < len(utxos) {
		if utxoCap > 0 { // expand sublist
			total += utxos[right].Amount
			right++
			utxoCap--
		} else { // pop the smallest utxo and append the current one
			total -= utxos[left].Amount
			total += utxos[right].Amount
			left++
			right++
		}
	}
	if total < amount {
		return nil, 0, fmt.Errorf("not enough btc in reserve - available : %v , tx amount : %v", total, amount)
	}
	return utxos[left:right], total, nil
}
