package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestObserverSet(t *testing.T) {
	observerSet := sample.ObserverSet(4)

	require.Equal(t, int(4), observerSet.Len())
	require.Equal(t, uint64(4), observerSet.LenUint())
	err := observerSet.Validate()
	require.NoError(t, err)

	observerSet.ObserverList[0] = "invalid"
	err = observerSet.Validate()
	require.Error(t, err)
}

func TestCheckReceiveStatus(t *testing.T) {
	err := types.CheckReceiveStatus(chains.ReceiveStatus_success)
	require.NoError(t, err)
	err = types.CheckReceiveStatus(chains.ReceiveStatus_failed)
	require.NoError(t, err)
	err = types.CheckReceiveStatus(chains.ReceiveStatus_created)
	require.Error(t, err)
}
