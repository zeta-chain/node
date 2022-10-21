package eth

import (
	"context"
	"math/big"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/observer"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/store"
	dbInfra "github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/store/infra"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/config"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/logger"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/model"

	"github.com/ethereum/go-ethereum/ethclient"
)

var _ observer.ChainObserver = (*ETHChainObserver)(nil)

type ETHChainObserver struct {
	db     store.Repository
	chain  common.Chain
	cfg    *config.ChainConfig
	client *ethclient.Client
	log    logger.Logger
}

func NewETHChainObserver(ctx context.Context, chain common.Chain, log *logger.Logger) (*ETHChainObserver, error) {
	cfg := config.GetChainConfig(string(chain))
	db, err := dbInfra.NewLevelDBRepository(string(chain))
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(cfg.Endpoint)
	if err != nil {
		ob.logger.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.EvmClient = client
	return &ETHChainObserver{
		db:     db,
		chain:  chain,
		cfg:    cfg,
		client: client,
		log:    log,
	}, nil
}

func (obs *ETHChainObserver) GetBlockHeight(ctx context.Context) (uint64, error) {
	header, err := ob.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return header.Number.Uint64(), nil
}

func (obs *ETHChainObserver) GetZetaPrice(ctx context.Context) (*big.Int, uint64, error) {
	return nil, uint64(0), nil
}

func (obs *ETHChainObserver) GetGasPrice(ctx context.Context) (uint64, error) {
	return uint64(0), nil
}

func (obs *ETHChainObserver) GetConnectorEvents(ctx context.Context, start uint64, end *uint64) ([]*model.ConnectorEvent, error) {
	return nil, nil
}

func (obs *ETHChainObserver) GetTxByHash(ctx context.Context, hash string, nonce int64) (*model.Receipt, error) {
	return nil, nil
}

func (obs *ETHChainObserver) GetConnectorEvents(ctx context.Context, start uint64, end *uint64) ([]*model.ConnectorEvent, error) {
	return nil, nil
}

func (obs *ETHChainObserver) DB() store.Repository {
	return obs.db
}
