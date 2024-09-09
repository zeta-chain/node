package rpc

import (
	"context"

	"github.com/zeta-chain/node/pkg/chains"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
)

// GetAdditionalChains returns the additional chains
func (c *Clients) GetAdditionalChains(ctx context.Context) ([]chains.Chain, error) {
	resp, err := c.Authority.ChainInfo(ctx, &authoritytypes.QueryGetChainInfoRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetChainInfo().Chains, nil
}
