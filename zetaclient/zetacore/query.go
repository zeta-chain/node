// GetAdditionalChains returns the additional chains
func (c *Client) GetAdditionalChains() ([]chains.Chain, error) {
	// TODO AFTER MERGE
	client := authoritytypes.NewQueryClient(c.grpcConn)
	resp, err := client.ChainInfo(context.Background(), &authoritytypes.QueryGetChainInfoRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetChainInfo().Chains, nil
}
