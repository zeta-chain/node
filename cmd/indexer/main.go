package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // this registers a sql driver
)

func main() {
	node := flag.String("node-ip", "127.0.0.1", "The IP address of the ZetaCore node")
	dbpath := flag.String("dbpath", "db.sqlite", "File path to the database")
	rebuild := flag.Bool("rebuild", false, "Rebuild the database from scratch (will erase and rebuild dbfile)")
	flag.Parse()

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

	//indexdb, err := indexdb2.NewIndexDB(db)
	_ = db
	_ = node
	_ = rebuild
}
