package net

import (
	"context"
	"fmt"

	rpcclient "github.com/cometbft/cometbft/rpc/client"

	"github.com/cosmos/evm/server/config"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
)

// PublicAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicAPI struct {
	networkVersion uint64
	tmClient       rpcclient.Client
}

// NewPublicAPI creates an instance of the public Net Web3 API.
func NewPublicAPI(ctx *server.Context, clientCtx client.Context) *PublicAPI {
	cfg, err := config.GetConfig(ctx.Viper)
	if err != nil {
		panic(err)
	}
	return &PublicAPI{
		networkVersion: cfg.EVM.EVMChainID,
		tmClient:       clientCtx.Client.(rpcclient.Client),
	}
}

// Version returns the current ethereum protocol version.
func (s *PublicAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}

// Listening returns if client is actively listening for network connections.
func (s *PublicAPI) Listening() bool {
	ctx := context.Background()
	netInfo, err := s.tmClient.NetInfo(ctx)
	if err != nil {
		return false
	}
	return netInfo.Listening
}

// PeerCount returns the number of peers currently connected to the client.
func (s *PublicAPI) PeerCount() int {
	ctx := context.Background()
	netInfo, err := s.tmClient.NetInfo(ctx)
	if err != nil {
		return 0
	}
	return len(netInfo.Peers)
}
