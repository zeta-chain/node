package sample

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestEthAddress(t *testing.T) {
	ethAddress := EthAddress()
	assert.NotEqual(t, ethcommon.Address{}, ethAddress)

	// don't generate the same address
	ethAddress2 := EthAddress()
	assert.NotEqual(t, ethAddress, ethAddress2)
}
