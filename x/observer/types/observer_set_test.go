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
	err := types.CheckReceiveStatus(chains.ReceiveStatus_Success)
	require.NoError(t, err)
	err = types.CheckReceiveStatus(chains.ReceiveStatus_Failed)
	require.NoError(t, err)
	err = types.CheckReceiveStatus(chains.ReceiveStatus_Created)
	require.Error(t, err)
}
