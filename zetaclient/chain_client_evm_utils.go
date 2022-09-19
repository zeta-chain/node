package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	"math/big"
	"strings"
)

func (ob *ChainObserver) ExternalChainWatcher() {
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

func (ob *ChainObserver) observeInTX() error {
	header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	counter, err := ob.GetPromCounter("rpc_getBlockByNumber_count")
	if err != nil {
		ob.logger.Error().Err(err).Msg("GetPromCounter:")
	}
	counter.Inc()

	// "confirmed" current block number
	confirmedBlockNum := header.Number.Uint64() - ob.confCount
	// skip if no new block is produced.
	if confirmedBlockNum <= ob.GetLastBlock() {
		return nil
	}
	toBlock := ob.GetLastBlock() + config.MaxBlocksPerPeriod // read at most 10 blocks in one go
	if toBlock >= confirmedBlockNum {
		toBlock = confirmedBlockNum
	}
	ob.sampleLogger.Info().Msgf("%s current block %d, querying from %d to %d, %d blocks left to catch up, watching MPI address %s", ob.chain, header.Number.Uint64(), ob.GetLastBlock()+1, toBlock, int(toBlock)-int(confirmedBlockNum), ob.ConnectorAddress.Hex())
	startBlock := ob.GetLastBlock() + 1
	// ============= query the Connector contract =============
	// Finally query the for the logs
	logs, err := ob.Connector.FilterZetaSent(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
		Context: context.TODO(),
	}, []ethcommon.Address{}, []*big.Int{})

	if err != nil {
		return err
	}
	cnt, err := ob.GetPromCounter("rpc_getLogs_count")
	if err != nil {
		return err
	}
	cnt.Inc()
	// Pull out arguments from logs
	for logs.Next() {
		event := logs.Event
		ob.logger.Info().Msgf("TxBlockNumber %d Transaction Hash: %s", event.Raw.BlockNumber, event.Raw.TxHash)

		destChain := config.FindChainByID(event.DestinationChainId)
		destAddr := types.BytesToEthHex(event.DestinationAddress)
		if strings.EqualFold(destAddr, config.Chains[destChain].ZETATokenContractAddress) {
			ob.logger.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
		}
		zetaHash, err := ob.zetaClient.PostSend(
			event.ZetaTxSenderAddress.Hex(),
			ob.chain.String(),
			types.BytesToEthHex(event.DestinationAddress),
			config.FindChainByID(event.DestinationChainId),
			event.ZetaValueAndGas.String(),
			event.ZetaValueAndGas.String(),
			base64.StdEncoding.EncodeToString(event.Message),
			event.Raw.TxHash.Hex(),
			event.Raw.BlockNumber,
			event.DestinationGasLimit.Uint64(),
			common.CoinType_Zeta,
		)
		if err != nil {
			ob.logger.Error().Err(err).Msg("error posting to zeta core")
			continue
		}
		ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
	}
	// ============= end of query the Connector contract =============

	// ============= query the incoming tx to TSS address ==============
	tssAddress := ob.Tss.Address()
	for bn := startBlock; bn <= toBlock; bn++ {
		block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
		if err != nil {
			ob.logger.Error().Err(err).Msg("error getting block")
			continue
		}
		for _, tx := range block.Transactions() {
			if tx.To() == nil {
				continue
			}
			if *tx.To() == tssAddress {
				receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), tx.Hash())
				if receipt.Status != 1 { // 1: successful, 0: failed
					ob.logger.Info().Msgf("tx %s failed; don't act", tx.Hash().Hex())
					continue
				}
				if err != nil {
					ob.logger.Err(err).Msg("TransactionReceipt")
					continue
				}
				from, err := ob.EvmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					ob.logger.Err(err).Msg("TransactionSender")
					continue
				}
				ob.logger.Info().Msgf("TSS inTx detected: %s, blocknum %d", tx.Hash().Hex(), receipt.BlockNumber)
				ob.logger.Info().Msgf("TSS inTx value: %s", tx.Value().String())
				ob.logger.Info().Msgf("TSS inTx from: %s", from.Hex())
				message := ""
				if len(tx.Data()) != 0 {
					message = hex.EncodeToString(tx.Data())
				}
				zetaHash, err := ob.zetaClient.PostSend(
					from.Hex(),
					ob.chain.String(),
					from.Hex(),
					"ZETA",
					tx.Value().String(),
					tx.Value().String(),
					message,
					tx.Hash().Hex(),
					receipt.BlockNumber.Uint64(),
					90_000,
					common.CoinType_Gas,
				)
				if err != nil {
					ob.logger.Error().Err(err).Msg("error posting to zeta core")
					continue
				}
				ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
			}
		}
	}
	// ============= end of query the incoming tx to TSS address ==============

	//ob.LastBlock = toBlock
	ob.setLastBlock(toBlock)
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, toBlock)
	err = ob.db.Put([]byte(PosKey), buf[:n], nil)
	if err != nil {
		ob.logger.Error().Err(err).Msg("error writing toBlock to db")
	}
	return nil
}

// query the base gas price for the block number bn.
func (ob *ChainObserver) GetBaseGasPrice() *big.Int {
	gasPrice, err := ob.EvmClient.SuggestGasPrice(context.TODO())
	if err != nil {
		ob.logger.Err(err).Msg("GetBaseGasPrice")
		return nil
	}
	return gasPrice
}
