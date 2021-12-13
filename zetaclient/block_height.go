package zetaclient

import (
	"context"
	"fmt"
	//"github.com/Meta-Protocol/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// GetBlockHeight returns the current height for metachain blocks
func (b *MetachainBridge) GetBlockHeight() (uint64, error) {
	client := types.NewQueryClient(b.grpcConn)
	height, err := client.LastMetaHeight(
		context.Background(),
		&types.QueryLastMetaHeightRequest{},
	)
	if err != nil {
		return 0, err
	}

	fmt.Printf("block height: %d\n", height.Height)
	return height.Height, nil
}

//func (b *MetachainBridge) GetLastBlockObserved(chain common.Chain) (uint64, error) {
//	client := types.NewQueryClient(b.grpcConn)
//	last_obs, err := client.LastBlockObserved(
//		context.Background(),
//		&types.QueryGetLastBlockObservedRequest{
//			Index: chain.String(),
//		},
//	)
//	if err != nil {
//		return 0, err
//	}
//
//	observed := last_obs.LastBlockObserved
//	fmt.Printf("last observed block height on chain %s: %d\n",
//		observed.Chain,
//		observed.Height)
//	return observed.Height, nil
//}
