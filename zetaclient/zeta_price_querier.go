package zetaclient

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"math"
	"math/big"
)

type ZetaPriceQuerier interface {
	// returns price (gasasset/zeta), blockNum, error
	GetZetaPrice() (*big.Int, uint64, error)
}

type UniswapV3ZetaPriceQuerier struct {
	UniswapV3Abi        *abi.ABI
	Client              *ethclient.Client
	PoolContractAddress ethcommon.Address
	Chain               common.Chain
}

var _ ZetaPriceQuerier = &UniswapV3ZetaPriceQuerier{}

type DummyZetaPriceQuerier struct {
	Chain  common.Chain
	Client *ethclient.Client
}

var _ ZetaPriceQuerier = &DummyZetaPriceQuerier{}

type UniswapV2ZetaPriceQuerier struct {
	UniswapV2Abi        *abi.ABI
	Client              *ethclient.Client
	PoolContractAddress ethcommon.Address
	Chain               common.Chain
}

var _ ZetaPriceQuerier = &UniswapV2ZetaPriceQuerier{}

// return the ratio GAS(ETH, BNB, MATIC, etc)/ZETA from Uniswap v3
// return price (gasasset/zeta), blockNum, error
func (q *UniswapV3ZetaPriceQuerier) GetZetaPrice() (*big.Int, uint64, error) {
	TIME_WINDOW := 600 // time weighted average price over last 10min (600s) period
	input, err := q.UniswapV3Abi.Pack("observe", []uint32{0, uint32(TIME_WINDOW)})
	if err != nil {
		return nil, 0, fmt.Errorf("fail to pack observe")
	}

	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	toAddr := q.PoolContractAddress
	res, err := q.Client.CallContract(context.TODO(), ethereum.CallMsg{
		From: fromAddr,
		To:   &toAddr,
		Data: input,
	}, nil)
	if err != nil {
		log.Err(err).Msgf("%s CallContract error", q.Chain)
		return nil, 0, err
	}
	bn, err := q.Client.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msgf("%s BlockNumber error", q.Chain)
		return nil, 0, err
	}
	output, err := q.UniswapV3Abi.Unpack("observe", res)
	if err != nil || len(output) != 2 {
		log.Err(err).Msgf("%s Unpack error or len(output) (%d) != 2", q.Chain, len(output))
		return nil, 0, err
	}
	cumTicks := *abi.ConvertType(output[0], new([2]*big.Int)).(*[2]*big.Int)
	tickDiff := big.NewInt(0).Div(big.NewInt(0).Sub(cumTicks[0], cumTicks[1]), big.NewInt(int64(TIME_WINDOW)))
	price := math.Pow(1.0001, float64(tickDiff.Int64())) * 1e18 // price is fixed point with decimal 18
	v, _ := big.NewFloat(price).Int(nil)
	return v, bn, nil
}

// return the ratio GAS(ETH, BNB, MATIC, etc)/ZETA from Uniswap v2 and its clone
// return price (gasasset/zeta), blockNum, error
func (q *UniswapV2ZetaPriceQuerier) GetZetaPrice() (*big.Int, uint64, error) {
	input, err := q.UniswapV2Abi.Pack("getReserves")
	if err != nil {
		return nil, 0, fmt.Errorf("fail to pack getReserves")
	}

	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	toAddr := q.PoolContractAddress
	res, err := q.Client.CallContract(context.TODO(), ethereum.CallMsg{
		From: fromAddr,
		To:   &toAddr,
		Data: input,
	}, nil)
	if err != nil {
		log.Err(err).Msgf("%s CallContract error", q.Chain)
		return nil, 0, err
	}
	bn, err := q.Client.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msgf("%s BlockNumber error", q.Chain)
		return nil, 0, err
	}
	output, err := q.UniswapV2Abi.Unpack("getReserves", res)
	if err != nil || len(output) != 3 {
		log.Err(err).Msgf("%s Unpack error or len(output) (%d) != 3", q.Chain, len(output))
		return nil, 0, err
	}
	reserve0 := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	reserve1 := *abi.ConvertType(output[1], new(*big.Int)).(**big.Int)
	r0, acc0 := big.NewFloat(0).SetInt(reserve0).Float64()
	r1, acc1 := big.NewFloat(0).SetInt(reserve1).Float64()

	if r0 <= 0 || r1 <= 0 || acc0 != big.Exact || acc1 != big.Exact {
		log.Err(err).Msgf("%s inexact conversion acc0=%s acc1=%s r0=%d r1=%d", q.Chain, acc0, acc1, reserve0, reserve1)
		return nil, 0, err
	}
	v, _ := big.NewFloat(r0 / r1 * 1.0e18).Int(nil)
	return v, bn, nil
}

// dummy price: always 1; returns 1e18, bn, and error
func (q *DummyZetaPriceQuerier) GetZetaPrice() (*big.Int, uint64, error) {
	bn, err := q.Client.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msgf("%s BlockNumber error", q.Chain)
		return nil, 0, err
	}
	v, _ := big.NewFloat(1.0e18).Int(nil)
	return v, bn, nil
}
