package ante_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	serverconfig "github.com/zeta-chain/node/server/config"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestSystemTxPriorityDecorator_AnteHandle(t *testing.T) {
	txConfig := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID).TxConfig

	testPrivKey, _ := sample.PrivKeyAddressPair()

	decorator := ante.NewSystemPriorityDecorator()
	mmd := MockAnteHandler{}
	// set priority to 10 before ante handler
	ctx := sdk.Context{}.WithIsCheckTx(true).WithPriority(10)

	tx, err := simtestutil.GenSignedMockTx(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		txConfig,
		[]sdk.Msg{},
		sdk.NewCoins(),
		simtestutil.DefaultGenTxGas,
		"testing-chain-id",
		[]uint64{0},
		[]uint64{0},
		testPrivKey,
	)
	require.NoError(t, err)
	ctx, err = decorator.AnteHandle(ctx, tx, false, mmd.AnteHandle)
	require.NoError(t, err)

	// check that priority is set to max int64
	priorityAfter := ctx.Priority()
	require.Equal(t, math.MaxInt64, int(priorityAfter))
}
