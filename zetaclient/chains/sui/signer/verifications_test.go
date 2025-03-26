package signer

import (
	"context"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestLiveCheckObjectsAreShared tests checkObjectsAreShared with testnet data
func TestLiveCheckObjectsAreShared(t *testing.T) {
	// uncomment to enable test
	//t.Skip("skipping live test")
	ctx := context.Background()
	client := sui.NewSuiClient("https://fullnode.testnet.sui.io")

	// all these objects are shared
	objectIds := []string{
		"0x6fc08f682551e52c2cc34362a20f744ba6a3d8d17f6583fa2f774887c4079700", // ZetaChain gateway - shared
		"0x0000000000000000000000000000000000000000000000000000000000000006", // Sui universal clock object - shared
		"0xf5ff7d5ba73b581bca6b4b9fa0049cd320360abd154b809f8700a8fd3cfaf7ca", // Cetus config - immutable
	}

	err := checkObjectsAreShared(ctx, client, objectIds)
	require.NoError(t, err)

	// contains a owned object
	objectIds = []string{
		"0x6fc08f682551e52c2cc34362a20f744ba6a3d8d17f6583fa2f774887c4079700",
		"0x01f724edef5461e280e533649b08c191a39af5f677025a226664a6626373b393", // TSS withdraw cap - owned
		"0x0000000000000000000000000000000000000000000000000000000000000006",
		"0xf5ff7d5ba73b581bca6b4b9fa0049cd320360abd154b809f8700a8fd3cfaf7ca",
	}

	err = checkObjectsAreShared(ctx, client, objectIds)
	require.Error(t, err)

	// contains a non existing object
	objectIds = []string{
		"0x6fc08f682551e52c2cc34362a20f744ba6a3d8d17f6583fa2f774887c4079700",
		"0x000000000000000000000000000000000000000000000000000000000000aaaa", // doesn't exist
		"0x0000000000000000000000000000000000000000000000000000000000000006",
		"0xf5ff7d5ba73b581bca6b4b9fa0049cd320360abd154b809f8700a8fd3cfaf7ca",
	}

	err = checkObjectsAreShared(ctx, client, objectIds)
	require.Error(t, err)
}
