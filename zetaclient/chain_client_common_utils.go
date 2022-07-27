package zetaclient

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"os"
	"strconv"
	"time"
)

func (ob *ChainObserver) BuildBlockIndex(dbpath, chain string) error {
	path := fmt.Sprintf("%s/%s", dbpath, chain) // e.g. ~/.zetaclient/ETH
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}
	ob.db = db
	envvar := ob.chain.String() + "_SCAN_CURRENT"
	if os.Getenv(envvar) != "" {
		log.Info().Msgf("envvar %s is set; scan from current block", envvar)
		header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return err
		}
		ob.LastBlock = header.Number.Uint64()
	} else { // last observed block
		buf, err := db.Get([]byte(PosKey), nil)
		if err != nil {
			log.Info().Msg("db PosKey does not exist; read from ZetaCore")
			ob.LastBlock = ob.getLastHeight()
			// if ZetaCore does not have last heard block height, then use current
			if ob.LastBlock == 0 {
				header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
				if err != nil {
					return err
				}
				ob.LastBlock = header.Number.Uint64()
			}
			buf2 := make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(buf2, ob.LastBlock)
			err := db.Put([]byte(PosKey), buf2[:n], nil)
			if err != nil {
				log.Error().Err(err).Msg("error writing ob.LastBlock to db: ")
			}
		} else {
			ob.LastBlock, _ = binary.Uvarint(buf)
		}
	}
	return nil
}

func (ob *ChainObserver) BuildReceiptsMap() {
	iter := ob.db.NewIterator(util.BytesPrefix([]byte(NonceTxKeyPrefix)), nil)
	for iter.Next() {
		key := string(iter.Key())
		nonce, err := strconv.ParseInt(key[len(NonceTxKeyPrefix):], 10, 64)
		if err != nil {
			log.Error().Err(err).Msgf("error parsing nonce: %s", key)
			continue
		}
		var receipt ethtypes.Receipt
		err = receipt.UnmarshalJSON(iter.Value())
		if err != nil {
			log.Error().Err(err).Msgf("error unmarshalling receipt: %s", key)
			continue
		}
		ob.outTXConfirmedReceipts[int(nonce)] = &receipt
		log.Info().Msgf("chain %s reading nonce %d with receipt of tx %s", ob.chain, nonce, receipt.TxHash.Hex())
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Error().Err(err).Msg("error iterating over db")
	}
}

func (ob *ChainObserver) GetPriceQueriers(chain string, uniswapV3ABI, uniswapV2ABI abi.ABI) (*UniswapV3ZetaPriceQuerier, *UniswapV2ZetaPriceQuerier, *DummyZetaPriceQuerier) {
	uniswapv3querier := &UniswapV3ZetaPriceQuerier{
		UniswapV3Abi:        &uniswapV3ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain].PoolContractAddress),
		Chain:               ob.chain,
	}
	uniswapv2querier := &UniswapV2ZetaPriceQuerier{
		UniswapV2Abi:        &uniswapV2ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain].PoolContractAddress),
		Chain:               ob.chain,
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
	MIN_OB_INTERVAL := 24
	switch chain {
	case common.MumbaiChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.POLY_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.POLYGON_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv3querier
		ob.BlockTime = config.POLY_BLOCK_TIME

	case common.GoerliChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.ETH_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.ETH_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv3querier
		ob.BlockTime = config.ETH_BLOCK_TIME

	case common.BSCTestnetChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.BSC_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.BSC_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv2querier
		ob.BlockTime = config.BSC_BLOCK_TIME

	case common.RopstenChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.ROPSTEN_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.ROPSTEN_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv3querier
		ob.BlockTime = config.ROPSTEN_BLOCK_TIME
	}
}

func (ob *ChainObserver) SetMinAndMaxNonce(trackers []types.OutTxTracker) error {
	minNonce, maxNonce := -1, 0
	for _, tracker := range trackers {
		intNonce, err := strconv.Atoi(tracker.Nonce)
		if err != nil {
			return err
		}
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
		ob.MaxNonce = minNonce
	}
	if maxNonce > 0 {
		ob.MaxNonce = maxNonce
	}
	return nil
}
