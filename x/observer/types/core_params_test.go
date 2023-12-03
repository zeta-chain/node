package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestCoreParamsList_Validate(t *testing.T) {
	t.Run("should return no error for default list", func(t *testing.T) {
		list := types.GetCoreParams()
		err := list.Validate()
		require.NoError(t, err)
	})

	t.Run("should return error for invalid chain id", func(t *testing.T) {
		list := types.GetCoreParams()
		list.CoreParams[0].ChainId = 999
		err := list.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found in chain list")
	})

	t.Run("should return error for duplicated chain ID", func(t *testing.T) {
		list := types.GetCoreParams()
		list.CoreParams = append(list.CoreParams, list.CoreParams[0])
		err := list.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicated chain id")
	})
}
