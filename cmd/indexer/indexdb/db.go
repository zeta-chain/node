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
        %s TEXT NOT NULL,
        %s TEXT PRIMARY KEY,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL
    );
    `, types.InboundFinalized, types.SendHash, types.InTxHash, types.Sender, types.SenderChain,
		types.Receiver, types.ReceiverChain, types.NewStatus, types.ZetaBurnt, types.ZetaMint)

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

	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT Not NULL,
        %s TEXT PRIMARY KEY,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL
    );
    `, types.OutboundTxFailed, types.SendHash, types.OutTxHash, types.ZetaMint,
		types.Chain, types.OldStatus, types.NewStatus)

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	cnt, err := idb.querier.VisitAllTxEvents(types.InboundFinalized, -1, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				//fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.InTxHash])
				_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, %s, %s, %s) values(?,?,?,?,?,?,?,?,?)",
					types.InboundFinalized, types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.NewStatus, types.ZetaBurnt, types.ZetaMint),
					kv[types.SendHash], kv[types.InTxHash], kv[types.Sender], kv[types.SenderChain], kv[types.Receiver], kv[types.ReceiverChain], kv[types.NewStatus], kv[types.ZetaBurnt], kv[types.ZetaMint])
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
	fmt.Printf("%s events processed : %d\n", types.InboundFinalized, cnt)

	cnt, err = idb.querier.VisitAllTxEvents(types.OutboundTxSuccessful, -1, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.InTxHash])
				_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s) values(?,?)", types.OutboundTxSuccessful, types.SendHash, types.OutTxHash),
					kv[types.SendHash], kv[types.OutTxHash])
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
	fmt.Printf("%s events processed : %d\n", types.OutboundTxSuccessful, cnt)

	cnt, err = idb.querier.VisitAllTxEvents(types.OutboundTxFailed, -1, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.OutTxHash])
				_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s) values(?,?,?,?,?,?)", types.OutboundTxFailed, types.SendHash, types.OutTxHash, types.ZetaMint, types.Chain, types.OldStatus, types.NewStatus),
					kv[types.SendHash],
					kv[types.OutTxHash],
					kv[types.ZetaMint],
					kv[types.Chain],
					kv[types.OldStatus],
					kv[types.NewStatus],
				)
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
	fmt.Printf("%s events processed : %d\n", types.OutboundTxFailed, cnt)

	return nil
}

func AttributeToMap(attr []sdk.Attribute) map[string]string {
	kv := make(map[string]string, len(attr))
	for _, v := range attr {
		kv[v.Key] = v.Value
	}
	return kv
}
