package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zeta-chain/zetacore/zetaclient"
	"os"
	"strconv"
	"strings"
)

const (
	PosKey                 = zetaclient.PosKey
	NonceTxHashesKeyPrefix = zetaclient.NonceTxHashesKeyPrefix
	NonceTxKeyPrefix       = zetaclient.NonceTxKeyPrefix
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dbPath := flag.String("db-path", homeDir+"/.zetaclientd/chainobserver/GOERLI", "path to level db")
	flag.Parse()
	log.Info().Msgf("dbPath: %s", *dbPath)

	db, err := leveldb.OpenFile(*dbPath, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open db")
	}
	defer db.Close()

	// last block number
	buf, err := db.Get([]byte(PosKey), nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to get pos")
	} else {
		lastBlock, _ := binary.Uvarint(buf)
		log.Info().Msgf("lastBlock: %d", lastBlock)
	}

	// nonceTxHashesMap
	iter := db.NewIterator(util.BytesPrefix([]byte(NonceTxHashesKeyPrefix)), nil)
	for iter.Next() {
		key := string(iter.Key())
		nonce, err := strconv.ParseInt(key[len(NonceTxHashesKeyPrefix):], 10, 64)
		if err != nil {
			log.Error().Err(err).Msgf("error parsing nonce: %s", key)
			continue
		}
		txHashes := strings.Split(string(iter.Value()), ",")
		log.Info().Msgf("reading nonce %d with %d tx hashes", nonce, len(txHashes))
		for _, txHash := range txHashes {
			fmt.Printf("  %s\n", txHash)
		}
	}
	iter.Release()
	if err = iter.Error(); err != nil {
		log.Error().Err(err).Msg("error iterating over db")
	}

	// nonceTxMap
	{
		iter := db.NewIterator(util.BytesPrefix([]byte(NonceTxKeyPrefix)), nil)
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
			log.Info().Msgf("reading nonce %d with receipt of tx %s", nonce, receipt.TxHash.Hex())
		}
		iter.Release()
		if err = iter.Error(); err != nil {
			log.Error().Err(err).Msg("error iterating over db")
		}
	}

}
