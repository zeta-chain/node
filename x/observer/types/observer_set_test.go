package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestObserverSet_Validate(t *testing.T) {
	observer1Address := sample.AccAddress()
	tt := []struct {
		name     string
		observer types.ObserverSet
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name:     "observer set with duplicate observer",
			observer: types.ObserverSet{ObserverList: []string{observer1Address, observer1Address}},
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, types.ErrDuplicateObserver)
			},
		},
		{
			name:     "observer set with invalid observer",
			observer: types.ObserverSet{ObserverList: []string{"invalid"}},
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "decoding bech32 failed")
			},
		},
		{
			name:     "observer set with valid observer",
			observer: types.ObserverSet{ObserverList: []string{observer1Address}},
			wantErr:  require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.wantErr(t, tc.observer.Validate())
		})

	}
}

func TestCheckReceiveStatus(t *testing.T) {
	err := types.CheckReceiveStatus(chains.ReceiveStatus_success)
	require.NoError(t, err)
	err = types.CheckReceiveStatus(chains.ReceiveStatus_failed)
	require.NoError(t, err)
	err = types.CheckReceiveStatus(chains.ReceiveStatus_created)
	require.Error(t, err)
}
