package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // this registers a sql driver
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/cmd/indexer/indexdb"
	"github.com/zeta-chain/zetacore/cmd/indexer/query"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	node := flag.String("node-ip", "127.0.0.1", "The IP address of the ZetaCore node")
	dbpath := flag.String("dbpath", "db.sqlite", "File path to the database")
	rebuild := flag.Bool("rebuild", false, "Rebuild the database from scratch (will erase and rebuild dbfile)")
	flag.Parse()

	_ = rebuild
	_ = node

	db, err := sql.Open("sqlite3", *dbpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	querier, err := query.NewZetaQuerier("3.20.194.40")
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Remove("test.db")
	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	idb, err := indexdb.NewIndexDB(db, querier)

	log.Info().Msgf("Rebuilding database...")
	start := time.Now()
	idb.Rebuild()
	duration := time.Since(start)
	log.Info().Msgf("Rebuilding database takes %s", duration)

	log.Info().Msgf("Start watching events...")
	idb.Start()

	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received; exit")
}
