package rpc

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type ConfigGetter interface {
	GetConfigParam(ctx context.Context, index uint32) (*boc.Cell, error)
}

// FetchGasConfigRPC query chain's config params and extract gas prices.
// Gas prices can be changes only by gov proposal (ie infrequently).
func FetchGasConfigRPC(ctx context.Context, rpc ConfigGetter) (tlb.GasLimitsPrices, error) {
	// https://docs.ton.org/develop/howto/blockchain-configs
	// https://tonviewer.com/config#21
	const configKeyGas = 21

	cell, err := rpc.GetConfigParam(ctx, configKeyGas)
	if err != nil {
		return tlb.GasLimitsPrices{}, errors.Wrap(err, "failed to get config param")
	}

	var cfg tlb.ConfigParam21
	if err = tlb.Unmarshal(cell, &cfg); err != nil {
		return tlb.GasLimitsPrices{}, errors.Wrap(err, "failed to unmarshal config param")
	}

	return cfg.GasLimitsPrices, nil
}

// ParseGasPrice parses gas price from the config and returns price in tons per 1 gas unit.
// You can take a look at definitions here:
// https://github.com/ton-blockchain/ton/blob/master/crypto/block/block.tlb
// https://docs.ton.org/develop/howto/blockchain-configs#param-20-and-21
//
// gas_prices#dd gas_price:uint64 gas_limit:uint64 gas_credit:uint64
// block_gas_limit:uint64 freeze_due_limit:uint64 delete_due_limit:uint64 = GasLimitsPrices;
//
// gas_prices_ext#de gas_price:uint64 gas_limit:uint64 special_gas_limit:uint64 gas_credit:uint64
// block_gas_limit:uint64 freeze_due_limit:uint64 delete_due_limit:uint64 = GasLimitsPrices;
//
// gas_flat_pfx#d1 flat_gas_limit:uint64 flat_gas_price:uint64 other:GasLimitsPrices = GasLimitsPrices;
func ParseGasPrice(cfg tlb.GasLimitsPrices) (uint64, error) {
	// tongo lib uses a concept of "sum types"
	// to decode different (sub)type of entities.
	// Basically, sumType is a struct property that is not empty (i.e. decoded).
	const (
		sumTypeGasPrices    = "GasPrices"
		sumTypeGasPricesExt = "GasPricesExt"
		sumTypeGasFlatPfx   = "GasFlatPfx"
	)

	// from TON docs: gas_price: This parameter reflects
	// the price of gas in the network, in nano tons per 65536 gas units (2^16).
	// We have 3 cases because TON node might return on of these 3 structs.
	switch cfg.SumType {
	case sumTypeGasPrices:
		return cfg.GasPrices.GasPrice >> 16, nil
	case sumTypeGasPricesExt:
		return cfg.GasPricesExt.GasPrice >> 16, nil
	case sumTypeGasFlatPfx:
		if cfg.GasFlatPfx.Other == nil {
			return 0, errors.New("GasFlatPfx.Other is nil")
		}
		return ParseGasPrice(*cfg.GasFlatPfx.Other)
	default:
		return 0, errors.Errorf("unknown SumType: %q", cfg.SumType)
	}
}
