package ton

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
)

type GlobalConfigurationFile = config.GlobalConfigurationFile

// ConfigFromURL downloads & parses lite server config.
//
//nolint:gosec
func ConfigFromURL(ctx context.Context, url string) (*GlobalConfigurationFile, error) {
	const timeout = 3 * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download config file: %s", res.Status)
	}

	return config.ParseConfig(res.Body)
}

func ConfigFromPath(path string) (*GlobalConfigurationFile, error) {
	return config.ParseConfigFile(path)
}

// ConfigFromSource returns a parsed configuration file from a URL or a file path.
func ConfigFromSource(ctx context.Context, urlOrPath string) (*GlobalConfigurationFile, error) {
	if u, err := url.Parse(urlOrPath); err == nil {
		return ConfigFromURL(ctx, u.String())
	}

	return ConfigFromPath(urlOrPath)
}

// ConfigGetter represents LiteAPI config getter.
type ConfigGetter interface {
	GetConfigParams(ctx context.Context, mode liteapi.ConfigMode, params []uint32) (tlb.ConfigParams, error)
}

// FetchGasConfig fetches gas price from the config.
func FetchGasConfig(ctx context.Context, getter ConfigGetter) (tlb.GasLimitsPrices, error) {
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

// ParseGasPrice parses gas price from the config and returns price in tons per 1 gas unit.
func ParseGasPrice(cfg tlb.GasLimitsPrices) (uint64, error) {
	// from TON docs: gas_price: This parameter reflects
	// the price of gas in the network, in nano tons per 65536 gas units (2^16).
	switch cfg.SumType {
	case "GasPrices":
		return cfg.GasPrices.GasPrice >> 16, nil
	case "GasPricesExt":
		return cfg.GasPricesExt.GasPrice >> 16, nil
	case "GasFlatPfx":
		if cfg.GasFlatPfx.Other == nil {
			return 0, errors.New("GasFlatPfx.Other is nil")
		}
		return ParseGasPrice(*cfg.GasFlatPfx.Other)
	default:
		return 0, errors.Errorf("unknown SumType: %q", cfg.SumType)
	}
}
