package main

import (
	"context"
	"database/sql"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"math/big"
	"time"
)

type Indexer struct {
	chainName string
	ethClient *ethclient.Client
	db        *sql.DB
	lastBlock int64
	from      ethcommon.Address
}

var (
	NewIndexStartBlock = map[string]int64{
		"GOERLI":     7967968,
		"BSCTESTNET": 24664408,
		"MUMBAI":     29195198,
	}
)

const CreateTable string = `
  CREATE TABLE IF NOT EXISTS txs_from_address (
  nonce INTEGER NOT NULL,
  chain TEXT NOT NULL,
  from_address TEXT NOT NULL,
  to_address TEXT NOT NULL,
  block_number INTEGER NOT NULL,
  txhash TEXT NOT NULL,
  status INTEGER NOT NULL,
  time DATETIME NOT NULL,
  gas_price INTEGER NOT NULL,
  gas_fee INTEGER NOT NULL,
  PRIMARY KEY (nonce, chain, from_address)
  );`

func NewIndexer(chainName string, endpoint string, fromAddress string) (*Indexer, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "clean.sqlite3")
	if err != nil {
		return nil, err
	}
	if db.Ping() != nil {
		return nil, fmt.Errorf("failed to connect to database")
	}

	addr := ethcommon.HexToAddress(fromAddress)
	if addr == (ethcommon.Address{}) {
		return nil, fmt.Errorf("invalid address")
	}

	return &Indexer{
		chainName: chainName,
		ethClient: client,
		from:      addr,
		db:        db,
		lastBlock: 0,
	}, nil
}

func (i *Indexer) Start() {
	var startBlock int64
	if NewIndex {
		startBlock = NewIndexStartBlock[i.chainName]
		i.lastBlock = startBlock
		log.Info().Msgf("NewIndee mode: re-building db table; starting block %d", startBlock)
		if _, err := i.db.Exec(CreateTable); err != nil {
			log.Fatal().Err(err).Msgf("failed to create table")
			return
		}
	} else { // start from the latest index block in db
		row := i.db.QueryRow("SELECT MAX(block_number) FROM txs_from_address WHERE chain = ?", i.chainName)
		if err := row.Scan(&startBlock); err != nil {
			log.Fatal().Err(err).Msgf("failed to get last block number")
			return
		}
		log.Info().Msgf("start from block (max block indexed + 1) %d", startBlock+1)
		i.lastBlock = startBlock + 1
	}

	go i.WatchChain()
}

func (i *Indexer) Stop() {

}

func (i *Indexer) WatchChain() {
	bn, err := i.ethClient.BlockNumber(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to get block number")
		return
	}
	log.Info().Msgf("current block number: %d", bn)
	if i.lastBlock >= int64(bn) {
		log.Info().Msgf("no new block to index")
		return
	}
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		if i.lastBlock >= int64(bn) {
			log.Info().Msgf("no new block to index: current block %d; exit", bn)
		}
		block, err := i.ethClient.BlockByNumber(context.TODO(), big.NewInt(i.lastBlock))
		if err != nil {
			log.Error().Err(err).Msgf("failed to get block %d", i.lastBlock)
			continue
		}
		txs := block.Transactions()
		for _, tx := range txs {
			receipt, err := i.ethClient.TransactionReceipt(context.TODO(), tx.Hash())
			if err != nil {
				log.Error().Err(err).Msgf("failed to get receipt for tx %s", tx.Hash().Hex())
				continue
			}
			from, err := i.ethClient.TransactionSender(context.TODO(), tx, block.Hash(), receipt.TransactionIndex)
			if err != nil {
				log.Error().Err(err).Msgf("failed to get sender for tx %s", tx.Hash().Hex())
				continue
			}
			if from == i.from {
				log.Info().Msgf("  %s-%d, hash %s", i.chainName, tx.Nonce(), tx.Hash().Hex())
				if _, err := i.db.Exec("INSERT INTO txs_from_address (nonce, chain, from_address, to_address, block_number, txhash, status, time, gas_price, gas_fee) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
					tx.Nonce(), i.chainName, from.Hex(), tx.To().Hex(), block.Number().Int64(), tx.Hash().Hex(), receipt.Status, block.Time(), tx.GasPrice().Int64(), tx.GasPrice().Int64()*int64(receipt.GasUsed)); err != nil {
					log.Error().Err(err).Msgf("failed to insert tx %s", tx.Hash().Hex())
				}
			}
		}
		i.lastBlock++
	}
}
