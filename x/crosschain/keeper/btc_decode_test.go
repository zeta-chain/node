package keeper_test

import (
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
)

func TestIsDecode(t *testing.T) {
	params, err := common.GetBTCChainParams(18332)
	assert.NoError(t, err)
	_, err = btcutil.DecodeAddress("14CEjTd5ci3228J45GdnGeUKLSSeCWUQxK", params)
	assert.NoError(t, err)
}
