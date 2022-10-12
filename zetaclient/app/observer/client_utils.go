package observer

import (
	"encoding/base64"
	"encoding/binary"
	"strings"

	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/model"
	"github.com/zeta-chain/zetacore/zetaclient/types"
)

func (ob *Observer) ExternalChainWatcher() {
	// At each tick, query the Connector contract
	for {
		select {
		case <-ob.ticker.C:
			err := ob.observeInTX()
			if err != nil {
				ob.logger.Err(err).Msg("observeInTX error")
				continue
			}
		case <-ob.stop:
			ob.logger.Info().Msg("ExternalChainWatcher stopped")
			return
		}
	}
}

func (ob *Observer) observeInTX() error {
	blockNumber, err := ob.ChainObserver.GetBlockHeight(ob.ctx)
	if err != nil {
		return err
	}
	counter, err := ob.GetPromCounter("rpc_getBlockByNumber_count")
	if err != nil {
		ob.logger.Error().Err(err).Msg("GetPromCounter:")
	}
	counter.Inc()

	// "confirmed" current block number
	confirmedBlockNum := blockNumber - ob.confCount
	// skip if no new block is produced.
	if confirmedBlockNum <= ob.GetLastBlock() {
		return nil
	}
	toBlock := ob.GetLastBlock() + config.MaxBlocksPerPeriod // read at most 10 blocks in one go
	if toBlock >= confirmedBlockNum {
		toBlock = confirmedBlockNum
	}
	ob.sampleLogger.Info().Msgf("%s current block %d, querying from %d to %d, %d blocks left to catch up", ob.chain, header.Number.Uint64(), ob.GetLastBlock()+1, toBlock, int(toBlock)-int(confirmedBlockNum))

	// Finally query connector for the logs
	events, err := ob.ChainObserver.GetConnectorEvents(ob.ctx, start, end)
	if err != nil {
		return err
	}
	cnt, err := ob.GetPromCounter("rpc_getLogs_count")
	if err != nil {
		return err
	}
	cnt.Inc()

	// Pull out arguments from logs
	for _, event := range events {
		ob.logger.Info().Msgf("TxBlockNumber %d Transaction Hash: %s", event.BlockNumber, event.TxHash)

		destChain := config.FindChainByID(event.DestinationChainId)
		destAddr := types.BytesToEthHex(event.DestinationAddress)
		if strings.EqualFold(destAddr, config.Chains[destChain].ZETATokenContractAddress) {
			ob.logger.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
		}
		zetaHash, err := ob.zetaClient.PostSend(
			event.ZetaTxSenderAddress,
			ob.chain.String(),
			event.DestinationAddress,
			config.FindChainByID(event.DestinationChainId),
			event.ZetaValueAndGas.String(),
			event.ZetaValueAndGas.String(),
			base64.StdEncoding.EncodeToString(event.Message),
			event.TxHash,
			event.BlockNumber,
			event.DestinationGasLimit.Uint64(),
		)
		if err != nil {
			ob.logger.Error().Err(err).Msg("error posting to zeta core")
			continue
		}
		ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
	}

	ob.LastBlock = toBlock
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, toBlock)
	err = ob.db.Put([]byte(model.PosKey), buf[:n], nil)
	if err != nil {
		ob.logger.Error().Err(err).Msg("error writing toBlock to db")
	}
	return nil
}
