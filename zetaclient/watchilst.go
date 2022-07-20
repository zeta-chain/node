package zetaclient

import (
	"fmt"
	"github.com/rs/zerolog/log"
	zetatypes "github.com/zeta-chain/zetacore/x/zetacore/types"
	"strings"
	"time"
)

// AddTxHashToWatchList adds an outbound TX hash to the watchlist
func (ob *ChainObserver) AddTxHashToWatchList(txHash string, nonce int, sendHash string) {
	outTx := OutTx{
		TxHash:   txHash,
		Nonce:    nonce,
		SendHash: sendHash,
	}

	if outTx.TxHash != "" { // TODO: this seems unnecessary
		ob.mu.Lock()
		ob.outTXPending[outTx.Nonce] = append(ob.outTXPending[outTx.Nonce], outTx.TxHash)
		ob.mu.Unlock()
		key := []byte(NonceTxHashesKeyPrefix + fmt.Sprintf("%d", outTx.Nonce))
		value := []byte(strings.Join(ob.outTXPending[outTx.Nonce], ","))
		if err := ob.db.Put(key, value, nil); err != nil {
			log.Error().Err(err).Msgf("AddTxHashToWatchList: error adding nonce %d tx hashes to db", outTx.Nonce)
		}

		log.Info().Msgf("add %s nonce %d TxHash watch list length: %d", ob.chain, outTx.Nonce, len(ob.outTXPending[outTx.Nonce]))
		ob.fileLogger.Info().Msgf("add %s nonce %d TxHash watch list length: %d", ob.chain, outTx.Nonce, len(ob.outTXPending[outTx.Nonce]))
	}
}

// PurgeTxHashWatchList  txhash from watch list which have no corresponding sendPending in zetacore.
// Returns the min/max nonce after purge
func (ob *ChainObserver) PurgeTxHashWatchList() (int, int, error) {
	purgedTxHashCount := 0
	sends, err := ob.zetaClient.GetAllPendingSend()
	if err != nil {
		return 0, 0, err
	}
	pendingNonces := make(map[int]bool)
	for _, send := range sends {
		if send.Status == zetatypes.SendStatus_PendingRevert && send.SenderChain == ob.chain.String() {
			pendingNonces[int(send.Nonce)] = true
		} else if send.Status == zetatypes.SendStatus_PendingOutbound && send.ReceiverChain == ob.chain.String() {
			pendingNonces[int(send.Nonce)] = true
		}
	}
	tNow := time.Now()
	ob.mu.Lock()
	for nonce, _ := range ob.outTXPending {
		if _, found := pendingNonces[nonce]; !found {
			txHashes := ob.outTXPending[nonce]
			delete(ob.outTXPending, nonce)
			if err = ob.db.Delete([]byte(NonceTxHashesKeyPrefix+fmt.Sprintf("%d", nonce)), nil); err != nil {
				log.Error().Err(err).Msgf("PurgeTxHashWatchList: error deleting nonce %d tx hashes from db", nonce)
			}
			purgedTxHashCount++
			log.Info().Msgf("PurgeTxHashWatchList: chain %s nonce %d removed", ob.chain, nonce)
			ob.fileLogger.Info().Msgf("PurgeTxHashWatchList: chain %s nonce %d removed txhashes %v", ob.chain, nonce, txHashes)
		}
	}
	ob.mu.Unlock()
	if purgedTxHashCount > 0 {
		log.Info().Msgf("PurgeTxHashWatchList: chain %s purged %d txhashes in %v", ob.chain, purgedTxHashCount, time.Since(tNow))
	}
	minNonce, maxNonce := -1, 0
	if len(pendingNonces) > 0 {
		for nonce, _ := range pendingNonces {
			if minNonce == -1 {
				minNonce = nonce
			}
			if nonce < minNonce {
				minNonce = nonce
			}
			if nonce > maxNonce {
				maxNonce = nonce
			}
		}
	}
	return minNonce, maxNonce, nil
}
