package bitcoin

import (
	"context"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/observer"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/store"
	dbInfra "github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/store/infra"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/config"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/logger"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/model"
)

var _ observer.ChainObserver = (*BTCChainObserver)(nil)

type BTCChainObserver struct {
	db     store.Repository
	chain  common.Chain
	cfg    *config.ChainConfig
	client *rpcclient.Client
	log    logger.Logger
}

func NewBTCChainObserver(ctx context.Context, mainCfg *config.Configuration, chain common.Chain, log logger.Logger) (*BTCChainObserver, error) {
	cfg := config.GetChainConfig(mainCfg, string(chain))
	db, err := dbInfra.NewLevelDBRepository(string(chain))
	if err != nil {
		log.Errorw("creating level db", "error", err.Error())
		return nil, err
	}
	connConfig := &rpcclient.ConnConfig{
		Host:              cfg.Endpoint,
		Endpoint:          "",
		HTTPPostMode:      true,
		EnableBCInfoHacks: true,
	}
	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		log.Errorw("dialing btc rpcclient", "error", err.Error())
		return nil, err
	}
	return &BTCChainObserver{
		db:     db,
		chain:  chain,
		cfg:    cfg,
		client: client,
		log:    log,
	}, nil
}

func (obs *BTCChainObserver) GetBlockHeight(ctx context.Context) (uint64, error) {
	blockNum, err := obs.client.GetBlockCount()
	return uint64(blockNum), err
}

func (obs *BTCChainObserver) GetZetaPrice(ctx context.Context) (*big.Int, uint64, error) {
	return nil, uint64(0), nil
}

func (obs *BTCChainObserver) GetGasPrice(ctx context.Context) (uint64, error) {
	return uint64(0), nil
}

func (obs *BTCChainObserver) GetConnectorEvents(ctx context.Context, start uint64, end *uint64, filter model.EventFilter) ([]*model.ConnectorEvent, error) {
	var events []*model.ConnectorEvent
	for i := start; i <= end; i++ {
		blockHash, err := obs.client.GetBlockHash(i)
		if err != nil {
			return nil, err
		}
		txResult, err := obs.client.GetBlockVerboseTx(blockHash)
		if err != nil {
			return nil, err
		}
		//txResult, err := obs.client.GetRawTransactionVerbose(txHashh)
		for _, rawTx := range txResult.RawTx {
			if event, ok := GetEventFromRawTx(rawTx); ok {
				events = append(events, GetEventFromRawTx(rawTx))
			}
		}
	}
	return events, nil
}

func (obs *BTCChainObserver) GetTxByHash(ctx context.Context, hash string, nonce int64) (*model.Receipt, error) {
	txHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
	result, err := obs.client.GetRawTransactionVerbose(txHash)
	if err != nil {
		return nil, err
	}
	return GetReceipt(result), nil
}

func (obs *BTCChainObserver) PrepareTx(ctx context.Context, outTx *model.OutTx) (*model.OutTx, error) {
	//var inputs []btcjson.TransactionInput
	//var amounts map[btcutil.Address]btcutil.Amount
	var lockTime int64
	inputs, amounts := GetInputsAmounts(outTx)
	wireMsgTx, err := obs.client.CreateRawTransaction(inputs, amounts, lockTime)
	return GetOutTx(wireMsgTx), nil
}

func (obs *BTCChainObserver) SendTx(ctx context.Context, outTx *model.OutTx) (*model.OutTxReceipt, error) {
	//var msgTx *wire.MsgTx
	var allowHighFees bool
	msgTx = GetMsgTx(outTx)
	hash, err := obs.client.SendRawTransaction(msgTx, allowHighFees)
	if err != nil {
		return nil, err
	}
	return GetTxReceipt(hash), nil
}

func (obs *BTCChainObserver) DB() store.Repository {
	return obs.db
}
