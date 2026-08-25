package e2etests

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
)

// TestDisallowVestingViaGroupProposal verifies that a vesting-account creation message cannot be
// smuggled onto the chain by wrapping it in a group proposal. The AuthzLimiterDecorator inspects
// the inner messages of group.MsgSubmitProposal at ante time, so the tx is rejected before the
// group module ever runs (no real group needs to exist).
func TestDisallowVestingViaGroupProposal(r *runner.E2ERunner, _ []string) {
	proposer := r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName)

	// a disabled message: create a vesting account
	vestingMsg := vesting.NewMsgCreateVestingAccount(
		sdk.MustAccAddressFromBech32(proposer),
		sdk.MustAccAddressFromBech32(sample.AccAddress()),
		sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
		time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
	)

	// wrap it in a group proposal; the group policy address is arbitrary since the tx is
	// rejected at ante before the group module validates it exists
	proposalMsg, err := group.NewMsgSubmitProposal(
		sample.AccAddress(),
		[]string{proposer},
		[]sdk.Msg{vestingMsg},
		"",
		group.Exec_EXEC_UNSPECIFIED,
		"create vesting account",
		"attempt to create a vesting account via a group proposal",
	)
	require.NoError(r, err)

	// broadcasting must fail: the ante decorator blocks the disabled inner message. Use the
	// no-retry variant since the rejection is deterministic — retrying only wastes ~25s and
	// floods the log with the same error six times.
	_, err = r.ZetaTxServer.BroadcastTxWithoutRetry(utils.OperationalPolicyName, proposalMsg)
	require.Error(r, err)
	require.ErrorContains(r, err, "found disabled msg type")
}
