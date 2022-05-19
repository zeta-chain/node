package indexdb

import (
	"database/sql"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/indexer/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
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
	query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT PRIMARY KEY,
        %s TEXT NOT NULL
    );
    `, types.InboundFinalized, types.SendHash, types.InTxHash)

	_, err := idb.db.Exec(query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT PRIMARY KEY,
        %s TEXT NOT NULL
    );
    `, types.OutboundTxSuccessful, types.SendHash, types.OutTxHash)

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	_, err = idb.querier.VisitAllTxEvents(types.InboundFinalized, 0, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				_, err := idb.db.Exec("INSERT INTO ?(?, ?) values(?,?)",
					types.InboundFinalized, types.SendHash, types.InTxHash,
					kv[types.SendHash], kv[types.InTxHash])
				if err != nil {
					fmt.Println(err)
					return nil
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func AttributeToMap(attr []sdk.Attribute) map[string]string {
	kv := make(map[string]string, len(attr))
	for _, v := range attr {
		kv[v.Key] = kv[v.Value]
	}
	return kv
}
