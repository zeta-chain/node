package zetabridge

import (
	"context"
	"fmt"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// GetBlockHeight returns the current height for metachain blocks
// FIXME: deprecate this in favor of tendermint RPC?
func (b *ZetaCoreBridge) GetBlockHeight() (int64, error) {
	client := types.NewQueryClient(b.grpcConn)
	height, err := client.LastZetaHeight(
		context.Background(),
		&types.QueryLastZetaHeightRequest{},
	)
	if err != nil {
		return 0, err
	}

	fmt.Printf("block height: %d\n", height.Height)
	return height.Height, nil
}
