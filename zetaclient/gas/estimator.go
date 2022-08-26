package gas

import (
	"container/ring"
	"context"
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

const (
	GAS_ESTIMATOR_BLOCKS = 100
	DEFAULT_MIN_PRICE    = 100000000
	MAX_GAS_PRICE        = 1000000000000
	// thresholds
	SAFE     = .35
	STANDARD = .60
	FAST     = .90
	FASTEST  = 1.0
)

type pricesByGas []*big.Int

func (pbg pricesByGas) Len() int           { return len(pbg) }
func (pbg pricesByGas) Swap(i, j int)      { pbg[i], pbg[j] = pbg[j], pbg[i] }
func (pbg pricesByGas) Less(i, j int) bool { return pbg[i].Cmp(pbg[j]) < 0 }

type GasEstimator struct {
	client       *ethclient.Client
	logger       zerolog.Logger
	ring         *ring.Ring
	lastBlock    *big.Int
	lastPrice    *big.Int
	count        int
	defaultPrice *big.Int
}

func NewGasEstimator(client *ethclient.Client, logger zerolog.Logger) *GasEstimator {
	return &GasEstimator{
		client: client,
		logger: logger,
		ring:   ring.New(GAS_ESTIMATOR_BLOCKS),
	}
}

func (ge *GasEstimator) GetPrice(ctx context.Context, blocknum *big.Int) (*big.Int, error) {
	if blocknum == nil {
		return nil, errors.New("invalid blocknum")
	}
	// blocknum was already processed return previous price
	if ge.lastBlock != nil && blocknum.Cmp(ge.lastBlock) <= 0 {
		return ge.lastPrice, nil
	}
	// get default price
	defaultPrice, err := ge.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	ge.defaultPrice = defaultPrice

	err = ge.fetchBlocks(ctx, blocknum)
	if err != nil {
		ge.logger.Warn().Err(err).Msgf("gas.GasEstimator.fetchBlocks")
		return defaultPrice, nil
	}

	price := ge.getPrice()
	ge.lastPrice = price
	ge.lastBlock = blocknum
	return price, nil
}

func (ge *GasEstimator) fetchBlocks(ctx context.Context, blocknum *big.Int) error {
	dif := big.NewInt(GAS_ESTIMATOR_BLOCKS)
	if ge.lastBlock != nil {
		dif = dif.Sub(blocknum, ge.lastBlock)
	}
	one := big.NewInt(1)
	first := big.NewInt(0).Sub(blocknum, dif)
	for i := first.Add(first, one); i.Cmp(blocknum) <= 0; i = i.Add(i, one) {
		block, err := ge.client.BlockByNumber(ctx, i)
		if err != nil {
			return err
		}
		minGas := ge.fetchMinBlockGas(block.Transactions())
		ge.ring.Value = minGas
		ge.ring = ge.ring.Next()
	}
	return nil
}

func (ge *GasEstimator) getPrice() *big.Int {
	prices := make([]*big.Int, GAS_ESTIMATOR_BLOCKS)
	var ix int
	threshold := FAST
	ge.ring.Do(func(elem any) {
		el, ok := elem.(*big.Int)
		if ok {
			prices[ix] = el
		} else {
			prices[ix] = big.NewInt(DEFAULT_MIN_PRICE)
		}
		ix++
	})
	// sort prices by min gas
	sort.Sort(pricesByGas(prices))
	// get price
	index := roundDown(float64(GAS_ESTIMATOR_BLOCKS) * threshold)
	if index >= len(prices) {
		index = len(prices) - 1
	}
	gasPrice := prices[index]
	minGasPrice := big.NewInt(DEFAULT_MIN_PRICE)
	maxGasPrice := big.NewInt(MAX_GAS_PRICE)
	if gasPrice.Cmp(maxGasPrice) > 0 {
		return ge.defaultPrice
	}
	if gasPrice.Cmp(minGasPrice) < 0 {
		return ge.defaultPrice
	}
	return gasPrice
}

func (ge *GasEstimator) fetchMinBlockGas(txs types.Transactions) *big.Int {
	zero := big.NewInt(0)
	minGasPrice := big.NewInt(0)
	for _, tx := range txs {
		txPrice := tx.GasPrice()
		if txPrice.Cmp(minGasPrice) < 0 || minGasPrice.Cmp(zero) == 0 {
			minGasPrice = tx.GasPrice()
		}
	}
	if minGasPrice.Cmp(zero) == 0 {
		minGasPrice.SetInt64(DEFAULT_MIN_PRICE)
	}
	return minGasPrice
}

func roundDown(val float64) int {
	if val < 0 {
		return int(val - 1.0)
	}
	return int(val)
}
