package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestParamKeyTable(t *testing.T) {
	kt := ParamKeyTable()

	ps := Params{}
	for _, psp := range ps.ParamSetPairs() {
		require.PanicsWithValue(t, "duplicate parameter key", func() {
			kt.RegisterType(psp)
		})
	}
}

func TestParamSetPairs(t *testing.T) {
	params := DefaultParams()
	pairs := params.ParamSetPairs()

	require.Equal(t, 0, len(pairs), "The number of param set pairs should match the expected count")
}

func TestParamsString(t *testing.T) {
	params := DefaultParams()
	out, err := yaml.Marshal(params)
	require.NoError(t, err)
	require.Equal(t, string(out), params.String())
}

func TestNewParams(t *testing.T) {
	params := NewParams()
	assert.True(t, params.Enabled)
	assert.Nil(t, params.Validate())
}
