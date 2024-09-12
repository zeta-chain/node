package ton

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	ton "github.com/tonkeeper/tongo/liteapi"
)

type Client struct {
	*ton.Client
	*SidecarClient
}

func (c *Client) WaitForBlocks(ctx context.Context) error {
	const (
		blocksToWait = 3
		interval     = 3 * time.Second
	)

	block, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return err
	}

	waitFor := block.Last.Seqno + blocksToWait

	for {
		freshBlock, err := c.GetMasterchainInfo(ctx)
		if err != nil {
			return err
		}

		if waitFor < freshBlock.Last.Seqno {
			return nil
		}

		time.Sleep(interval)
	}
}

type SidecarClient struct {
	baseURL string
	c       *http.Client
}

var ErrNotHealthy = fmt.Errorf("TON node is not healthy yet")

func NewSidecarClient(baseURL string) *SidecarClient {
	c := &http.Client{Timeout: 3 * time.Second}
	return &SidecarClient{baseURL, c}
}

// Faucet represents the faucet information.
//
//nolint:revive,stylecheck // comes from my-local-ton
type Faucet struct {
	InitialBalance   int64  `json:"initialBalance"`
	PrivateKey       string `json:"privateKey"`
	PublicKey        string `json:"publicKey"`
	WalletRawAddress string `json:"walletRawAddress"`
	Mnemonic         string `json:"mnemonic"`
	WalletVersion    string `json:"walletVersion"`
	WorkChain        int32  `json:"workChain"`
	SubWalletId      int    `json:"subWalletId"`
	Created          bool   `json:"created"`
}

// LiteServerURL returns the URL to the lite server config
func (c *SidecarClient) LiteServerURL() string {
	return fmt.Sprintf("%s/lite-client.json", c.baseURL)
}

// GetFaucet returns the faucet information.
func (c *SidecarClient) GetFaucet(ctx context.Context) (Faucet, error) {
	resp, err := c.get(ctx, "faucet.json")
	if err != nil {
		return Faucet{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Faucet{}, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	var faucet Faucet
	if err := json.NewDecoder(resp.Body).Decode(&faucet); err != nil {
		return Faucet{}, err
	}

	return faucet, nil
}

// Status checks the health of the TON node. Returns ErrNotHealthy or nil.
func (c *SidecarClient) Status(ctx context.Context) error {
	resp, err := c.get(ctx, "status")
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(ErrNotHealthy, "status %d. %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *SidecarClient) get(ctx context.Context, path string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.c.Do(req)
}
