package ton

import (
	"fmt"
	"net/http"

	"github.com/tonkeeper/tongo/config"
)

// ConfigFromURL downloads & parses config.
//
//nolint:gosec
func ConfigFromURL(url string) (*config.GlobalConfigurationFile, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download config file: %s", res.Status)
	}

	return config.ParseConfig(res.Body)
}
