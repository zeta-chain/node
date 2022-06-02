package indexdb

import (
	"database/sql"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/cmd/indexer/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"time"
)

type IndexDB struct {
	db                 *sql.DB
	querier            *query.ZetaQuerier
	lastBlockProcessed int64
}

func NewIndexDB(sqldb *sql.DB, querier *query.ZetaQuerier) (*IndexDB, error) {

	return &IndexDB{
		querier: querier,
		db:      sqldb,
	}, nil
}

func (idb *IndexDB) Start() {
	err := idb.db.QueryRow("select max(blocknum) from block").Scan(&idb.lastBlockProcessed)
	if err != nil {
		log.Error().Err(err).Msg(" error querying max(blocknum) from block; please rebuild")
		return
	} else {
		log.Info().Msgf("latest indexed blocknum %d", idb.lastBlockProcessed)
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for range ticker.C {
			block, err := idb.querier.LatestBlock()
			if err != nil {
				log.Error().Err(err).Msg("LatestBlock error")
				continue
			}
			if block.Header.Height > idb.lastBlockProcessed {
				for i := idb.lastBlockProcessed + 1; i <= block.Header.Height; i++ {
					err = idb.processBlock(i)
					if err != nil {
						log.Error().Err(err).Msgf("processBlock on block %d error", i)
					}
					idb.lastBlockProcessed = i
					log.Info().Msgf("processed block %d; catching up to %d", i, block.Header.Height)
				}

			}
		}
	}()
}

func (idb *IndexDB) processBlock(bn int64) error {

	cnt, err := idb.querier.VisitAllTxEvents(types.InboundFinalized, bn, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				err := idb.processFinalized(res, kv)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Info().Msgf("block %d: %s events processed : %d", bn, types.InboundFinalized, cnt)

	cnt, err = idb.querier.VisitAllTxEvents(types.OutboundTxSuccessful, bn, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				err := idb.processOutboundSuccessful(res, kv)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Info().Msgf("block %d: %s events processed : %d", bn, types.OutboundTxSuccessful, cnt)

	cnt, err = idb.querier.VisitAllTxEvents(types.OutboundTxFailed, bn, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				err := idb.processOutboundFailed(res, kv)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Info().Msgf("block %d: %s events processed : %d", bn, types.OutboundTxFailed, cnt)

	err = idb.insertBlockTable(bn)
	log.Info().Msgf("block %d: logging block info", bn)
	if err != nil {
		return err
	}

	return nil
}

// #nosec -- suppress G201 warning: formating SQL query; arguments not from user inputs.
func (idb *IndexDB) Rebuild() error {
	// 0. clear existing tables
	drop := fmt.Sprintf("DROP TABLE IF EXISTS txs")
	_, err := idb.db.Exec(drop)
	if err != nil {
		return err
	}
	drop = fmt.Sprintf("DROP TABLE IF EXISTS block")
	_, err = idb.db.Exec(drop)
	if err != nil {
		return err
	}
	drop = fmt.Sprintf("DROP TABLE IF EXISTS %s", types.InboundFinalized)
	_, err = idb.db.Exec(drop)
	if err != nil {
		return err
	}
	drop = fmt.Sprintf("DROP TABLE IF EXISTS %s", types.OutboundTxFailed)
	_, err = idb.db.Exec(drop)
	if err != nil {
		return err
	}
	drop = fmt.Sprintf("DROP TABLE IF EXISTS %s", types.OutboundTxSuccessful)
	_, err = idb.db.Exec(drop)
	if err != nil {
		return err
	}

	// 1. create tables
	//#nosec G201
	query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS txs (
		%s TEXT PRIMARY KEY,
		%s TEXT,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		Status TEXT NOT NULL,
		lastupdate INTEGER
    );
    `, types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.ZetaBurnt, types.ZetaMint, types.Message)

	_, err = idb.db.Exec(query)
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
		timestamp TIMESTAMP NOT NULL,
		blocknumber INTEGER NOT NULL
    );
    `, types.InboundFinalized, types.SendHash, types.InTxHash, types.Sender, types.SenderChain,
		types.Receiver, types.ReceiverChain, types.NewStatus, types.ZetaBurnt, types.ZetaMint)

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	// #nosec G201
	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT Not NULL,
        %s TEXT PRIMARY KEY,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		blocknumber INTEGER NOT NULL
    );
    `, types.OutboundTxSuccessful, types.SendHash, types.OutTxHash, types.ZetaMint,
		types.Chain, types.OldStatus, types.NewStatus)

	_, err = idb.db.Exec(query)
	if err != nil {
		return err
	}

	// #nosec G201
	query = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s(
        %s TEXT Not NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
        %s TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
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
		blocktimestamp TIMESTAMP,
		querytimestamp TIMESTAMP,
		numtxs INTEGER,
		txhashes TEXT[]
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
	err = idb.insertBlockTable(block.Header.Height)
	if err != nil {
		fmt.Printf("cannot insert latest block from zetacore node: %s\n", err)
	}
	idb.lastBlockProcessed = block.Header.Height

	cnt, err := idb.querier.VisitAllTxEvents(types.InboundFinalized, -1, func(res *sdk.TxResponse) error {
		for _, v := range res.Logs {
			for _, vv := range v.Events {
				kv := AttributeToMap(vv.Attributes)
				err := idb.processFinalized(res, kv)
				if err != nil {
					return err
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
				err2 := idb.processOutboundSuccessful(res, kv)
				if err2 != nil {
					return err2
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
				err2 := idb.processOutboundFailed(res, kv)
				if err2 != nil {
					return err2
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

func (idb *IndexDB) insertBlockTable(bn int64) error {
	block, err := idb.querier.BlockByHeight(bn)
	if err != nil {
		fmt.Printf("cannot query TxResponsesByBlock from zetacore node: %s\n", err)
		return err
	}
	txResponses, err := idb.querier.TxResponsesByBlock(bn)
	if err != nil {
		fmt.Printf("cannot query TxResponsesByBlock from zetacore node: %s\n", err)
		return err
	}
	var txhashes []string
	for _, v := range txResponses {
		txhashes = append(txhashes, v.TxHash)
		j, _ := v.Tx.MarshalJSON()
		fmt.Printf("%s\n%s\n", v.Tx.GetTypeUrl(), j)
	}
	_, err = idb.db.Exec("INSERT INTO block(blocknum, blocktimestamp, querytimestamp, numtxs, txhashes) values($1,$2,$3,$4,$5)",
		block.Header.Height, block.Header.Time, time.Now().UTC(), len(txResponses), pq.Array(txhashes))
	if err != nil {
		fmt.Printf("cannot insert lastblock into database: %s\n", err)
		return err
	}
	return nil
}

func (idb *IndexDB) processOutboundFailed(res *sdk.TxResponse, kv map[string]string) error {
	fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.OutTxHash])
	_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, timestamp,blocknumber) values($1,$2,$3,$4,$5,$6,$7,$8)", types.OutboundTxFailed, types.SendHash, types.OutTxHash, types.ZetaMint, types.Chain, types.OldStatus, types.NewStatus),
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

	_, err = idb.db.Exec(fmt.Sprintf("UPDATE  txs set Status = $1, lastupdate = $2, %s = $4  where SendHash = $3", types.ZetaMint), kv[types.NewStatus], res.Height, kv[types.SendHash], kv[types.ZetaMint])
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (idb *IndexDB) processOutboundSuccessful(res *sdk.TxResponse, kv map[string]string) error {
	fmt.Printf("%s:%s\n", kv[types.SendHash], kv[types.OutTxHash])
	_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, timestamp,blocknumber) values($1,$2,$3,$4,$5,$6,$7,$8)", types.OutboundTxSuccessful, types.SendHash, types.OutTxHash, types.ZetaMint, types.Chain, types.OldStatus, types.NewStatus),
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
		return err
	}

	_, err = idb.db.Exec(fmt.Sprintf("UPDATE  txs set Status = $1, lastupdate=$2, %s = $4 where SendHash = $3", types.ZetaMint), kv[types.NewStatus], res.Height, kv[types.SendHash], kv[types.ZetaMint])
	if err != nil {
		return err
	}
	return nil
}

func (idb *IndexDB) processFinalized(res *sdk.TxResponse, kv map[string]string) error {
	_, err := idb.db.Exec(fmt.Sprintf("INSERT INTO  %s(%s, %s, %s, %s, %s, %s, %s, %s, %s, timestamp,blocknumber) values($1,$2,$3,$4,$5,$6,$7,$8, $9, $10, $11)",
		types.InboundFinalized, types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.NewStatus, types.ZetaBurnt, types.ZetaMint),
		kv[types.SendHash], kv[types.InTxHash], kv[types.Sender], kv[types.SenderChain], kv[types.Receiver], kv[types.ReceiverChain], kv[types.NewStatus], kv[types.ZetaBurnt], kv[types.ZetaMint], res.Timestamp, res.Height)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = idb.db.Exec(fmt.Sprintf("INSERT INTO  txs (%s, %s, %s, %s, %s, %s, %s, %s, %s, Status, lastupdate) values($1,$2,$3,$4,$5,$6,$7,$10, $11, $8, $9)",
		types.SendHash, types.InTxHash, types.Sender, types.SenderChain, types.Receiver, types.ReceiverChain, types.ZetaBurnt, types.ZetaMint, types.Message),
		kv[types.SendHash], kv[types.InTxHash], kv[types.Sender], kv[types.SenderChain], kv[types.Receiver], kv[types.ReceiverChain], kv[types.ZetaBurnt], kv[types.NewStatus],
		res.Height, kv[types.ZetaMint], kv[types.Message])
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func AttributeToMap(attr []sdk.Attribute) map[string]string {
	kv := make(map[string]string, len(attr))
	for _, v := range attr {
		kv[v.Key] = v.Value
	}
	return kv
}
