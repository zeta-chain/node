package config

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
)

type GlobalConfigurationFile = config.GlobalConfigurationFile

var ErrDownload = errors.New("failed to download config file")

// Getter represents LiteAPI config params getter.
// Don't be confused because config param in this case represent on-chain params,
// not lite-client's ADNL json config to connect to the network.
//
// Read more at https://docs.ton.org/develop/howto/blockchain-configs
type Getter interface {
	GetConfigParams(ctx context.Context, mode liteapi.ConfigMode, params []uint32) (tlb.ConfigParams, error)
}

// FromURL downloads & parses lite server config.
//
//nolint:gosec
func FromURL(ctx context.Context, url string) (*GlobalConfigurationFile, error) {
	const timeout = 3 * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(ErrDownload, err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Wrap(ErrDownload, res.Status)
	}

	return config.ParseConfig(res.Body)
}

// FromPath parses config file from path.
func FromPath(path string) (*GlobalConfigurationFile, error) {
	return config.ParseConfigFile(path)
}

// FromSource returns a parsed configuration file from a URL or a file path.
func FromSource(ctx context.Context, urlOrPath string) (*GlobalConfigurationFile, error) {
	if cfg, err := FromPath(urlOrPath); err == nil {
		return cfg, nil
	}

	u, err := url.Parse(urlOrPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	return FromURL(ctx, u.String())
}

// FetchGasConfig fetches gas price from the config.
func FetchGasConfig(ctx context.Context, getter Getter) (tlb.GasLimitsPrices, error) {
	// https://docs.ton.org/develop/howto/blockchain-configs
	// https://tonviewer.com/config#21
	const configKeyGas = 21

	response, err := getter.GetConfigParams(ctx, 0, []uint32{configKeyGas})
	if err != nil {
		return tlb.GasLimitsPrices{}, errors.Wrap(err, "failed to get config params")
	}

	ref, ok := response.Config.Get(configKeyGas)
	if !ok {
		return tlb.GasLimitsPrices{}, errors.Errorf("config key %d not found", configKeyGas)
	}

	var cfg tlb.ConfigParam21
	if err = tlb.Unmarshal(&ref.Value, &cfg); err != nil {
		return tlb.GasLimitsPrices{}, errors.Wrap(err, "failed to unmarshal config param")
	}

	return cfg.GasLimitsPrices, nil
}

type RPCGetter interface {
	GetConfigParam(ctx context.Context, index uint32) (*boc.Cell, error)
}

func FetchGasConfigRPC(ctx context.Context, rpc RPCGetter) (tlb.GasLimitsPrices, error) {
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
