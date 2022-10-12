package observer

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/adapters/signer"
)

type ChainObserver interface {
	// control methods
	Start()
	Stop()
	//ExternalChainWatcher()
	//GasPriceWatcher()
	//ExchangeRateWatcher()
	//OutTxWatcher()
	// getters
	Chain() *common.Chain
	Endpoint() string
	Ticker() *time.Ticker
	TSSSigner() signer.TSSSigner
	LastBlock() uint64
	ConfirmationsCount() uint64
	//BlockTimeSeconds() uint64
	//TxWatchList() map[string]string
	//OutTxChan() chan model.OutTx // TODO: make OutTx generic
	IsSendOutTxProcessed(string, int) (bool, bool, error)
	//StopChan() chan struct{}
	//CriticalLog() *zerolog.Logger
	//Log() zerolog.Logger
	PostNonceIfNotRecorded() error
	// metrics
	GetPromGauge(string) (prometheus.Gauge, error)
}
