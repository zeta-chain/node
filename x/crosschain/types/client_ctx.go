package types

import (
	"github.com/cosmos/cosmos-sdk/client"
)

var (
	ClientCtx client.Context
)

func RegisterClientCtx(clientCtx client.Context) error {
	ClientCtx = clientCtx
	return nil
}
