//go:build test

// This file builds a full in-process app via zetasimulation.NewSimApp, which requires the `test`
// build tag; without it the EVM configurator panics. The tag keeps a bare `go test ./app/ante/...`
// (no -tags) from taking the rest of the package down with a confusing panic.

package ante_test

import (
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	serverconfig "github.com/zeta-chain/node/server/config"
	zetasimulation "github.com/zeta-chain/node/simulation"
	"github.com/zeta-chain/node/testutil/sample"
)

// TestProductionAnteHandler_RejectsVestingViaGroupProposal is an in-process integration check that
// wires the real app keepers into the full production ante chain (via app.DisabledAuthzMsgs) and
// runs a signed group.MsgSubmitProposal that wraps MsgCreateVestingAccount through it. It asserts
// the tx is rejected — the same outcome the localnet e2e test (disallow_vesting_via_group_proposal)
// verifies over the wire, but executable without Docker.
func TestProductionAnteHandler_RejectsVestingViaGroupProposal(t *testing.T) {
	// ARRANGE
	zetaApp, err := zetasimulation.NewSimApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		simtestutil.AppOptionsMap{},
	)
	require.NoError(t, err)

	encCfg := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID)

	// build the ante handler from the exact options app.New installs, so dropping any field there
	// (e.g. DisabledAuthzMsgs) fails this test rather than silently passing
	anteHandler, err := ante.NewAnteHandler(
		app.AnteHandlerOptions(zetaApp, encCfg.TxConfig.SignModeHandler()),
	)
	require.NoError(t, err)

	testPrivKey, testAddress := sample.PrivKeyAddressPair()
	_, testAddress2 := sample.PrivKeyAddressPair()

	vestingMsg := vesting.NewMsgCreateVestingAccount(
		testAddress, testAddress2,
		sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
		time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
	)
	proposalMsg, err := group.NewMsgSubmitProposal(
		sample.AccAddress(),
		[]string{testAddress.String()},
		[]sdk.Msg{vestingMsg},
		"",
		group.Exec_EXEC_UNSPECIFIED,
		"title",
		"summary",
	)
	require.NoError(t, err)

	tx, err := simtestutil.GenSignedMockTx(
		rand.New(rand.NewSource(1)),
		encCfg.TxConfig,
		[]sdk.Msg{proposalMsg},
		sdk.NewCoins(),
		simtestutil.DefaultGenTxGas,
		"testing-chain-id",
		[]uint64{0},
		[]uint64{0},
		testPrivKey,
	)
	require.NoError(t, err)

	ctx := zetaApp.NewContext(true)

	// ACT
	_, err = anteHandler(ctx, tx, false)

	// ASSERT: the ante decorator rejects the group-wrapped vesting-create before the group module runs
	require.Error(t, err)
	require.ErrorContains(t, err, "found disabled msg type")
}
