package indexdb

import (
	"database/sql"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/indexer/query"
)

type IndexDB struct {
	db      *sql.DB
	querier *query.ZetaQuerier
}

func NewIndexDB(sqldb *sql.DB, querier *query.ZetaQuerier) (*IndexDB, error) {

	return &IndexDB{
		querier: querier,
		db:      sqldb,
	}, nil
}

func (idb *IndexDB) Rebuild() error {
	// 1. create tables
	query := `
    CREATE TABLE IF NOT EXISTS finalized(
        sendHash TEXT PRIMARY KEY,
        inTxHash TEXT NOT NULL
    );
    `

	_, err := idb.db.Exec(query)
	if err != nil {
		return err
	}

	query = `
    CREATE TABLE IF NOT EXISTS mined(
        sendHash TEXT PRIMARY KEY,
        outTxHash TEXT NOT NULL
    );
    `

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	idb.querier.VisitAllTxEvents("SendFinalized", 0, func(res *sdk.TxResponse) error {
		//fmt.Println(res.Logs)
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				var sendHash, inTxHash string
				for _, attr := range vv.Attributes {
					//fmt.Println(attr.Key, attr.Value)
					if attr.Key == "Index" {
						fmt.Println(attr.Key, attr.Value)
						sendHash = attr.Value
					} else if attr.Key == "InTxHash" {
						fmt.Println(attr.Key, attr.Value)
						inTxHash = attr.Value
					}
				}
				_, err := idb.db.Exec("INSERT INTO finalized(sendHash, inTxHash) values(?,?)", sendHash, inTxHash)
				if err != nil {
					fmt.Println(err)
					return nil
				}
			}
			fmt.Println("####")
		}
		return nil
	})

	return nil
}
