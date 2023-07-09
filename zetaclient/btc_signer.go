package zetaclient

import (
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
func (signer *BTCSigner) SignWithdrawTx(to *btcutil.AddressWitnessPubKeyHash, amount float64, feeRate float64, utxos []btcjson.ListUnspentResult, db *gorm.DB, height uint64) (*wire.MsgTx, error) {
	var total float64
	var prevOuts []btcjson.ListUnspentResult
	// select N utxo sufficient to cover the amount
	//estimateFee := size (100 inputs + 2 output) * feeRate
	estimateFee := 0.00001 // FIXME: proper fee estimation
	for _, utxo := range utxos {
		// check for pending utxos
		if _, err := getPendingUTXO(db, utxoKey(utxo)); err != nil {
			if err == gorm.ErrRecordNotFound {
				total = total + utxo.Amount
				prevOuts = append(prevOuts, utxo)

				if total >= amount+estimateFee {
					break
				}
			} else {
				return nil, err
			}
		}
	}
	if total < amount {
		return nil, fmt.Errorf("not enough btc in reserve - available : %v , tx amount : %v", total, amount)
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
	btcFees := float64(tx.SerializeSize()) * feeRate / 1024 //FIXME: feeRate KB is 1000B or 1024B?
	fees, err := getSatoshis(btcFees)
	if err != nil {
		return nil, err
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
	txOut.Value = remainingSatoshis - fees
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

	// update pending utxos pendingUtxos
	err = signer.updatePendingUTXOs(db, prevOuts)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (signer *BTCSigner) Broadcast(signedTx *wire.MsgTx) error {
	fmt.Printf("BTCSigner: Broadcasting: %s\n", signedTx.TxHash().String())

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

	logger.Info().Msgf("BTC TryProcessOutTx: %s, value %d to %s", send.Index, send.GetCurrentOutTxParam().Amount.BigInt(), send.GetCurrentOutTxParam().Receiver)
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

	addr, err := btcutil.DecodeAddress(send.GetCurrentOutTxParam().Receiver, config.BitconNetParams)
	if err != nil {
		logger.Error().Err(err).Msgf("cannot decode address %s ", send.GetCurrentOutTxParam().Receiver)
		return
	}
	to, ok := addr.(*btcutil.AddressWitnessPubKeyHash)
	if err != nil || !ok {
		logger.Error().Err(err).Msgf("cannot decode address %s ", send.GetCurrentOutTxParam().Receiver)
		return
	}

	logger.Info().Msgf("SignWithdrawTx: to %s, value %d", addr.EncodeAddress(), send.GetCurrentOutTxParam().Amount.Uint64()/1e8)
	logger.Info().Msgf("using utxos: %v", btcClient.utxos)
	// FIXME: gas price?
	tx, err := signer.SignWithdrawTx(to, float64(send.GetCurrentOutTxParam().Amount.Uint64())/1e8, float64(gasprice.Int64())/1e8*1024, btcClient.utxos, btcClient.db, height)
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
	//}

}

func (signer *BTCSigner) updatePendingUTXOs(db *gorm.DB, utxos []btcjson.ListUnspentResult) error {
	for _, utxo := range utxos {
		// Try to find existing record in DB to populate primary key
		var pendingUTXO clienttypes.PendingUTXOSQLType
		db.Where("Key = ?", utxoKey(utxo)).First(&pendingUTXO)

		// If record doesn't exist, it will be created by the Save function
		pendingUTXO.UTXO = utxo
		pendingUTXO.Key = utxoKey(utxo)
		if err := db.Save(&pendingUTXO).Error; err != nil {
			return err
		}
	}
	return nil
}
