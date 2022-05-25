package indexdb

import (
	"database/sql"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/indexer/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"time"
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

func (idb *IndexDB) processBlock(bn int64) error {

	cnt, err := idb.querier.VisitAllTxEvents(types.InboundFinalized, bn, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				//fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.InTxHash])
				_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, %s, %s, %s, timestamp,blocknumber) values(?,?,?,?,?,?,?,?,?, ?,?)",
					types.InboundFinalized, types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.NewStatus, types.ZetaBurnt, types.ZetaMint),
					kv[types.SendHash], kv[types.InTxHash], kv[types.Sender], kv[types.SenderChain], kv[types.Receiver], kv[types.ReceiverChain], kv[types.NewStatus], kv[types.ZetaBurnt], kv[types.ZetaMint], res.Timestamp, res.Height)
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

	block, err := idb.querier.BlockByHeight(bn)
	if err != nil {
		fmt.Printf("cannot query latest block from zetacore node: %s\n", err)
		return err
	}
	_, err = idb.db.Exec("INSERT INTO block(blocknum, blocktimestamp, querytimestamp, numtxs) values(?,?,?,?)",
		block.Header.Height, block.Header.Time, time.Now().UTC(), len(block.Data.Txs))
	if err != nil {
		fmt.Printf("cannot insert lastblock into database: %s\n", err)
		return err
	}
	return nil
}

func (idb *IndexDB) Rebuild() error {
	// 1. create tables
	query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS txs (
		%s TEXT PRIMARY KEY,
		%s TEXT,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		Status TEXT NOT NULL,
		lastupdate INTEGER
    );
    `, types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.ZetaBurnt)

	_, err := idb.db.Exec(query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT NOT NULL,
        %s TEXT PRIMARY KEY,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		blocknumber INTEGER NOT NULL
    );
    `, types.InboundFinalized, types.SendHash, types.InTxHash, types.Sender, types.SenderChain,
		types.Receiver, types.ReceiverChain, types.NewStatus, types.ZetaBurnt, types.ZetaMint)

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT Not NULL,
        %s TEXT PRIMARY KEY,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		blocknumber INTEGER NOT NULL
    );
    `, types.OutboundTxSuccessful, types.SendHash, types.OutTxHash, types.ZetaMint,
		types.Chain, types.OldStatus, types.NewStatus)

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT Not NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		blocknumber INTEGER NOT NULL,
		PRIMARY KEY ( %s, %s)
    );
    `, types.OutboundTxFailed, types.SendHash, types.OutTxHash, types.ZetaMint,
		types.Chain, types.OldStatus, types.NewStatus, types.SendHash, types.OutTxHash)

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS block (
		blocknum INTEGER PRIMARY KEY,
		blocktimestamp DATETIME,
		querytimestamp DATETIME,
		numtxs INTEGER
    );
    `)
	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	block, err := idb.querier.LatestBlock()
	if err != nil {
		fmt.Printf("cannot query latest block from zetacore node: %s\n", err)
	}
	_, err = idb.db.Exec("INSERT INTO block(blocknum, blocktimestamp, querytimestamp, numtxs) values(?,?,?,?)",
		block.Header.Height, block.Header.Time, time.Now().UTC(), len(block.Data.Txs))
	if err != nil {
		fmt.Printf("cannot insert lastblock into database: %s\n", err)
	}

	cnt, err := idb.querier.VisitAllTxEvents(types.InboundFinalized, -1, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				//fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.InTxHash])
				_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, %s, %s, %s, timestamp,blocknumber) values(?,?,?,?,?,?,?,?,?, ?,?)",
					types.InboundFinalized, types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.NewStatus, types.ZetaBurnt, types.ZetaMint),
					kv[types.SendHash], kv[types.InTxHash], kv[types.Sender], kv[types.SenderChain], kv[types.Receiver], kv[types.ReceiverChain], kv[types.NewStatus], kv[types.ZetaBurnt], kv[types.ZetaMint], res.Timestamp, res.Height)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				_, err = idb.db.Exec(fmt.Sprintf("INSERT INTO  txs (%s, %s, %s, %s, %s, %s, %s, Status, lastupdate) values(?,?,?,?,?,?,?,?,?)",
					types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.ZetaBurnt),
					kv[types.SendHash], kv[types.InTxHash], kv[types.Sender], kv[types.SenderChain], kv[types.Receiver], kv[types.ReceiverChain], kv[types.ZetaBurnt], kv[types.NewStatus],
					res.Height)
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
				fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.OutTxHash])
				_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, timestamp,blocknumber) values(?,?,?,?,?,?,?,?)", types.OutboundTxSuccessful, types.SendHash, types.OutTxHash, types.ZetaMint, types.Chain, types.OldStatus, types.NewStatus),
					kv[types.SendHash],
					kv[types.OutTxHash],
					kv[types.ZetaMint],
					kv[types.Chain],
					kv[types.OldStatus],
					kv[types.NewStatus],
					res.Timestamp,
					res.Height,
				)
				if err != nil {
					fmt.Println(err)
					return err
				}

				_, err = idb.db.Exec(fmt.Sprintf("UPDATE  txs set Status = ?, lastupdate=?  where SendHash = ?"), kv[types.NewStatus], res.Height, kv[types.SendHash])
				if err != nil {
					fmt.Println(err)
					return err
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
				_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, timestamp,blocknumber) values(?,?,?,?,?,?, ?,?)", types.OutboundTxFailed, types.SendHash, types.OutTxHash, types.ZetaMint, types.Chain, types.OldStatus, types.NewStatus),
					kv[types.SendHash],
					kv[types.OutTxHash],
					kv[types.ZetaMint],
					kv[types.Chain],
					kv[types.OldStatus],
					kv[types.NewStatus],
					res.Timestamp,
					res.Height,
				)
				if err != nil {
					fmt.Println(err)
					return err
				}

				_, err = idb.db.Exec(fmt.Sprintf("UPDATE  txs set Status = ?, lastupdate = ?  where SendHash = ?"), kv[types.NewStatus], res.Height, kv[types.SendHash])
				if err != nil {
					fmt.Println(err)
					return err
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
