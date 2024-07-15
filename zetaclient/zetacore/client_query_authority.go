package zetacore

import (
	"context"

	"github.com/zeta-chain/zetacore/pkg/chains"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
)

// GetAdditionalChains returns the additional chains
func (c *Client) GetAdditionalChains(ctx context.Context) ([]chains.Chain, error) {
	resp, err := c.client.authority.ChainInfo(ctx, &authoritytypes.QueryGetChainInfoRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetChainInfo().Chains, nil
}
