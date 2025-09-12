package observer

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/logs"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// SaveBroadcastedTx saves successfully broadcasted transaction
func (ob *Observer) SaveBroadcastedTx(txHash string, nonce uint64) error {
	outboundID := ob.OutboundID(nonce)
	ob.Mu().Lock()
	ob.tssOutboundHashes[txHash] = true
	ob.broadcastedTx[outboundID] = txHash
	ob.Mu().Unlock()

	broadcastEntry := clienttypes.ToOutboundHashSQLType(txHash, outboundID)
	if err := ob.DB().Client().Save(&broadcastEntry).Error; err != nil {
		format := "failed to save broadcasted outbound hash %s for %s"
		return errors.Wrapf(err, format, txHash, outboundID)
	}
	ob.logger.Outbound.Info().
		Str(logs.FieldTx, txHash).
		Str(logs.FieldOutboundID, outboundID).
		Msg("saved broadcasted outbound hash to the database")

	return nil
}

// LoadLastBlockScanned loads the last scanned block from the database
func (ob *Observer) LoadLastBlockScanned(ctx context.Context) error {
	err := ob.Observer.LoadLastBlockScanned()
	if err != nil {
		return errors.Wrapf(err, "error LoadLastBlockScanned for chain %d", ob.Chain().ChainId)
	}

	// observer will scan from the last block when 'lastBlockScanned == 0', this happens when:
	// 1. environment variable is set explicitly to "latest"
	// 2. environment variable is empty and last scanned block is not found in DB
	if ob.LastBlockScanned() == 0 {
		blockNumber, err := ob.rpc.GetBlockCount(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to get block count")
		}
		// #nosec G115 always positive
		ob.WithLastBlockScanned(uint64(blockNumber))
	}

	// bitcoin regtest starts from hardcoded block 100
	if chains.IsBitcoinRegnet(ob.Chain().ChainId) {
		ob.WithLastBlockScanned(RegnetStartBlock)
	}
	ob.Logger().Chain.Info().Uint64("last_block_scanned", ob.LastBlockScanned()).Send()

	return nil
}

// loadBroadcastedTxMap loads broadcasted transactions from the database
func (ob *Observer) loadBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.OutboundHashSQLType

	tx := ob.DB().Client().Find(&broadcastedTransactions)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "unable to find broadcasted txs")
	}

	for _, entry := range broadcastedTransactions {
		ob.tssOutboundHashes[entry.Hash] = true
		ob.broadcastedTx[entry.Key] = entry.Hash
	}

	return nil
}
