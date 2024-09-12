package ton

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tonkeeper/tongo/config"
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
