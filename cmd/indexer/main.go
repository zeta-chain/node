package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq" // this registers a sql driver
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/cmd/indexer/indexdb"
	"github.com/zeta-chain/zetacore/cmd/indexer/query"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	EXTERNAL_CHAINS = []string{"GOERLI", "BSCTESTNET", "MUMBAI", "ROPSTEN"}
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get current username")
		return
	}

	node := flag.String("node-ip", "127.0.0.1", "The IP address of the ZetaCore node")
	rebuild := flag.Bool("rebuild", false, "Rebuild the database from scratch (will erase and rebuild dbfile)")
	dbhost := flag.String("dbhost", "localhost", "host URL of the PostgreSQL database")
	dbport := flag.Int64("dbport", 5432, "port of the PostgresSQL database")
	dbuser := flag.String("dbuser", user.Username, "username of PostgresSQL database")
	dbpasswd := flag.String("dbpasswd", "", "password of PostgresSQL database")
	dbname := flag.String("dbname", "testdb", "database name of PostgresSQL database")
	scanRange := flag.String("scan-range", "0:9223372036854775807", "rescan from this block")
	flag.Parse()

	var startBlock, endBlock int64
	var err1, err2 error
	if *scanRange != "" {
		parts := strings.Split(*scanRange, ":")
		if len(parts) != 2 {
			fmt.Println("scan-range must be of the form <start>:<end> both inclusive")
			return
		}
		startBlock, err1 = strconv.ParseInt(parts[0], 10, 64)
		endBlock, err2 = strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil || startBlock > endBlock {
			fmt.Println("scan-range must be of the form <start>:<end> both inclusive")
			return
		}
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", *dbhost, *dbport, *dbuser, *dbpasswd, *dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Info().Msgf("connected to psql server %s", psqlInfo)

	querier, err := query.NewZetaQuerier(*node)
	if err != nil {
		log.Error().Err(err).Msg("NewZetaQuerier error")
		return
	}

	clientMap := make(map[string]*ethclient.Client)
	for _, chain := range EXTERNAL_CHAINS {
		envvar := chain + "_ENDPOINT"
		endpoint := os.Getenv(envvar)
		log.Info().Msgf("%s=%s, connecting...", envvar, endpoint)
		if len(endpoint) != 0 {
			client, err := ethclient.Dial(endpoint)
			if err != nil {
				log.Error().Err(err)
				continue
			}
			clientMap[chain] = client
		}
	}

	idb, err := indexdb.NewIndexDB(db, querier, clientMap)
	if err != nil {
		log.Error().Err(err).Msg("NewIndexDB error")
		return
	}

	if *rebuild {
		log.Info().Msgf("Rebuilding database...")
		start := time.Now()
		err = idb.Rebuild()
		duration := time.Since(start)
		log.Info().Err(err).Msgf("Rebuilding database takes %s", duration)
	}

	idb.LastBlockProcessed = startBlock
	idb.EndBlock = endBlock
	log.Info().Msgf("Start watching events from block %d to block %d...", startBlock, endBlock)
	done := make(chan bool)
	idb.Start(done)

	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Info().Msgf("awaiting signal...")
	select {
	case <-ch:
	case <-done:
	}
	log.Info().Msg("stop signal received; exit")
}
