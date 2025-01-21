package observer

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// SaveBroadcastedTx saves successfully broadcasted transaction
func (ob *Observer) SaveBroadcastedTx(txHash string, nonce uint64) {
	outboundID := ob.OutboundID(nonce)
	ob.Mu().Lock()
	ob.tssOutboundHashes[txHash] = true
	ob.broadcastedTx[outboundID] = txHash
	ob.Mu().Unlock()

	broadcastEntry := clienttypes.ToOutboundHashSQLType(txHash, outboundID)
	if err := ob.DB().Client().Save(&broadcastEntry).Error; err != nil {
		ob.logger.Outbound.Error().
			Err(err).
			Msgf("SaveBroadcastedTx: error saving broadcasted txHash %s for outbound %s", txHash, outboundID)
	}
	ob.logger.Outbound.Info().Msgf("SaveBroadcastedTx: saved broadcasted txHash %s for outbound %s", txHash, outboundID)
}

// LoadLastBlockScanned loads the last scanned block from the database
func (ob *Observer) LoadLastBlockScanned(ctx context.Context) error {
	err := ob.Observer.LoadLastBlockScanned(ob.Logger().Chain)
	if err != nil {
		return errors.Wrapf(err, "error LoadLastBlockScanned for chain %d", ob.Chain().ChainId)
	}

	// observer will scan from the last block when 'lastBlockScanned == 0', this happens when:
	// 1. environment variable is set explicitly to "latest"
	// 2. environment variable is empty and last scanned block is not found in DB
	if ob.LastBlockScanned() == 0 {
		blockNumber, err := ob.rpc.GetBlockCount(ctx)
		if err != nil {
			return errors.Wrapf(err, "error GetBlockCount for chain %d", ob.Chain().ChainId)
		}
		// #nosec G115 always positive
		ob.WithLastBlockScanned(uint64(blockNumber))
	}

	// bitcoin regtest starts from hardcoded block 100
	if chains.IsBitcoinRegnet(ob.Chain().ChainId) {
		ob.WithLastBlockScanned(RegnetStartBlock)
	}
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from block %d", ob.Chain().ChainId, ob.LastBlockScanned())

	return nil
}

// LoadBroadcastedTxMap loads broadcasted transactions from the database
func (ob *Observer) LoadBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.OutboundHashSQLType
	if err := ob.DB().Client().Find(&broadcastedTransactions).Error; err != nil {
		ob.logger.Chain.Error().Err(err).Msgf("error iterating over db for chain %d", ob.Chain().ChainId)
		return err
	}
	for _, entry := range broadcastedTransactions {
		ob.tssOutboundHashes[entry.Hash] = true
		ob.broadcastedTx[entry.Key] = entry.Hash
	}
	return nil
}
