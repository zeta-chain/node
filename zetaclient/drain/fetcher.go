//go:build drain

package drain

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/draintx"
)

// httpFetchTimeout bounds a single payload fetch.
const httpFetchTimeout = 15 * time.Second

// HTTPFetcher fetches the drain payload over HTTP (e.g. an object-store URL).
type HTTPFetcher struct {
	URL    string
	client *http.Client
}

// NewHTTPFetcher creates an HTTPFetcher for the given URL.
func NewHTTPFetcher(url string) *HTTPFetcher {
	return &HTTPFetcher{URL: url, client: &http.Client{Timeout: httpFetchTimeout}}
}

// Fetch retrieves and decodes the current payload.
func (f *HTTPFetcher) Fetch(ctx context.Context) (draintx.Payload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.URL, nil)
	if err != nil {
		return draintx.Payload{}, errors.Wrap(err, "unable to build request")
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return draintx.Payload{}, errors.Wrap(err, "unable to fetch payload")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return draintx.Payload{}, errors.Errorf("unexpected status %d", resp.StatusCode)
	}

	var payload draintx.Payload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return draintx.Payload{}, errors.Wrap(err, "unable to decode payload")
	}
	return payload, nil
}
