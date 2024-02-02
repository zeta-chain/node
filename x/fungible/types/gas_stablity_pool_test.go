package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestGetGasStabilityPoolAddress(t *testing.T) {
	address := types.GasStabilityPoolAddress()
	assert.False(t, address.Empty())
}

func TestGetGasStabilityPoolAddressEVM(t *testing.T) {
	address := types.GasStabilityPoolAddressEVM()
	assert.NotEmpty(t, address)
}
