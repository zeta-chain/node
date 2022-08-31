package zetaclient

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func (ob *ChainObserver) BuildBlockIndex(dbpath, chain string) error {
	logger := ob.logger
	path := fmt.Sprintf("%s/%s", dbpath, chain) // e.g. ~/.zetaclient/ETH
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}
	ob.db = db
	envvar := ob.chain.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		logger.Info().Msgf("envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == "latest" {
			header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.setLastBlock(header.Number.Uint64())
		} else {
			scanFromBlockInt, err := strconv.ParseInt(scanFromBlock, 10, 64)
			if err != nil {
				return err
			}
			ob.setLastBlock(uint64(scanFromBlockInt))
		}
	} else { // last observed block
		buf, err := db.Get([]byte(PosKey), nil)
		if err != nil {
			logger.Info().Msg("db PosKey does not exist; read from ZetaCore")
			ob.setLastBlock(ob.getLastHeight())
			// if ZetaCore does not have last heard block height, then use current
			if ob.GetLastBlock() == 0 {
				header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
				if err != nil {
					return err
				}
				ob.setLastBlock(header.Number.Uint64())
			}
			buf2 := make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(buf2, ob.GetLastBlock())
			err := db.Put([]byte(PosKey), buf2[:n], nil)
			if err != nil {
				logger.Error().Err(err).Msg("error writing ob.LastBlock to db: ")
			}
		} else {
			lastBlock, _ := binary.Uvarint(buf)
			ob.setLastBlock(lastBlock)
		}
	}
	return nil
}

func (ob *ChainObserver) BuildReceiptsMap() {
	logger := ob.logger
	iter := ob.db.NewIterator(util.BytesPrefix([]byte(NonceTxKeyPrefix)), nil)
	for iter.Next() {
		key := string(iter.Key())
		nonce, err := strconv.ParseInt(key[len(NonceTxKeyPrefix):], 10, 64)
		if err != nil {
			logger.Error().Err(err).Msgf("error parsing nonce: %s", key)
			continue
		}
		var receipt ethtypes.Receipt
		err = receipt.UnmarshalJSON(iter.Value())
		if err != nil {
			logger.Error().Err(err).Msgf("error unmarshalling receipt: %s", key)
			continue
		}
		ob.outTXConfirmedReceipts[int(nonce)] = &receipt
		//log.Info().Msgf("chain %s reading nonce %d with receipt of tx %s", ob.chain, nonce, receipt.TxHash.Hex())
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		logger.Error().Err(err).Msg("error iterating over db")
	}
}

func (ob *ChainObserver) GetPriceQueriers(chain string, uniswapV3ABI, uniswapV2ABI abi.ABI) (*UniswapV3ZetaPriceQuerier, *UniswapV2ZetaPriceQuerier, *DummyZetaPriceQuerier) {
	uniswapv3querier := &UniswapV3ZetaPriceQuerier{
		UniswapV3Abi:        &uniswapV3ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain].PoolContractAddress),
		Chain:               ob.chain,
		TokenOrder:          config.Chains[chain].PoolTokenOrder,
	}
	uniswapv2querier := &UniswapV2ZetaPriceQuerier{
		UniswapV2Abi:        &uniswapV2ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain].PoolContractAddress),
		Chain:               ob.chain,
		TokenOrder:          config.Chains[chain].PoolTokenOrder,
	}
	dummyQuerier := &DummyZetaPriceQuerier{
		Chain:  ob.chain,
		Client: ob.EvmClient,
	}
	return uniswapv3querier, uniswapv2querier, dummyQuerier
}

func (ob *ChainObserver) SetChainDetails(chain common.Chain,
	uniswapv3querier *UniswapV3ZetaPriceQuerier,
	uniswapv2querier *UniswapV2ZetaPriceQuerier) {
	MinObInterval := 24
	switch chain {
	case common.MumbaiChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.PolygonBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.PolygonConfirmationCount
		ob.BlockTime = config.PolygonBlockTime

	case common.GoerliChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.EthBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.EthConfirmationCount
		ob.BlockTime = config.EthBlockTime

	case common.BSCTestnetChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.BscBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.BscConfirmationCount
		ob.BlockTime = config.BscBlockTime

	case common.RopstenChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.RopstenBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.RopstenConfirmationCount
		ob.BlockTime = config.RopstenBlockTime
	}
	switch config.Chains[chain.String()].PoolContract {
	case clienttypes.UniswapV2:
		ob.ZetaPriceQuerier = uniswapv2querier
	case clienttypes.UniswapV3:
		ob.ZetaPriceQuerier = uniswapv3querier
	default:
		ob.logger.Error().Msgf("unknown pool contract type: %d", config.Chains[chain.String()].PoolContract)
	}
}

func (ob *ChainObserver) SetMinAndMaxNonce(trackers []types.OutTxTracker) error {
	minNonce, maxNonce := int64(-1), int64(0)
	for _, tracker := range trackers {
		conv, err := strconv.Atoi(tracker.Nonce)
		if err != nil {
			return err
		}
		intNonce := int64(conv)
		if minNonce == -1 {
			minNonce = intNonce
		}
		if intNonce < minNonce {
			minNonce = intNonce
		}
		if intNonce > maxNonce {
			maxNonce = intNonce
		}
	}
	if minNonce != -1 {
		atomic.StoreInt64(&ob.MinNonce, minNonce)
	}
	if maxNonce > 0 {
		atomic.StoreInt64(&ob.MaxNonce, maxNonce)
	}
	return nil
}
