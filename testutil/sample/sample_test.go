package sample

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEthAddress(t *testing.T) {
	ethAddress := EthAddress()
	require.NotEqual(t, ethcommon.Address{}, ethAddress)
}
