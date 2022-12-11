package types

import (
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var (
	RPCClient rpcclient.Client
)

func RegisterRPCClient(rpcclient rpcclient.Client) error {
	RPCClient = rpcclient
	return nil
}
