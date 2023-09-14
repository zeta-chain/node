package sample

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEthAddress(t *testing.T) {
	ethAddress := EthAddress()
	require.NotEqual(t, ethcommon.Address{}, ethAddress)

	// don't generate the same address
	ethAddress2 := EthAddress()
	require.NotEqual(t, ethAddress, ethAddress2)
}
