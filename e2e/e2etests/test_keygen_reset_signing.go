package e2etests

import (
	"math"
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestKeygenResetSigning asserts that resetting the keygen record does not stop the network
// from signing with the TSS key it already has.
//
// # HOW THIS REPRODUCES THE MAINNET INCIDENT — READ BEFORE CHANGING
//
// On mainnet nobody sent a transaction. An observer's validator was jailed and began
// unbonding, the staking hook dropped it from the observer set, and the next BeginBlocker
// (x/observer/abci.go) saw the observer count change and ran:
//
//	k.DisableInboundOnly(ctx)
//	k.SetKeygen(ctx, types.Keygen{BlockNumber: math.MaxInt64})
//
// That struct literal zeroes every field it does not name, so the grantee list was erased and
// the status reset to PendingKeygen. Every zetaclient then restarted and none could start
// again: an empty grantee list meant an empty p2p whitelist.
//
// Driving a validator through jailing and unbonding inside an e2e test is slow and fragile.
// MsgAddObserver with AddNodeAccountOnly is used INSTEAD because its handler runs that exact
// same pair of calls (msg_server_add_observer.go), so it writes byte-for-byte the state the
// BeginBlocker wrote, in one transaction. It is a stand-in for convenience.
//
// It is NOT what happened on mainnet, and no add-observer transaction was involved there. Do
// not read this test as evidence of the cause.
//
// The message targets an observer that already exists, so SetNodeAccount rewrites an identical
// record and the observer set is untouched. The only effects are the keygen reset and the
// inbound pause, which is what we want to isolate.
func TestKeygenResetSigning(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	// Fund the deployer while the network is still healthy, so the withdraw later is the only
	// thing under test. This doubles as the "before" half of the comparison: an inbound that
	// works now, and an outbound that must still work once the record is erased.
	depositHash := r.DepositEtherToDeployer()
	depositCCTX := utils.WaitCctxMinedByInboundHash(r.Ctx, depositHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, depositCCTX, crosschaintypes.CctxStatus_OutboundMined)

	tssBefore, err := r.ObserverClient.TSS(r.Ctx, &observertypes.QueryGetTSSRequest{})
	require.NoError(r, err)
	require.NotEmpty(r, tssBefore.TSS.TssPubkey, "test needs a TSS that has already been generated")

	keygenBefore, err := r.ObserverClient.Keygen(r.Ctx, &observertypes.QueryGetKeygenRequest{})
	require.NoError(r, err)
	require.NotNil(r, keygenBefore.Keygen)
	require.NotEmpty(
		r,
		keygenBefore.Keygen.GranteePubkeys,
		"keygen record must start populated, otherwise this test proves nothing",
	)

	// Reuse an existing node account so the message adds nobody and changes no observer.
	nodeAccounts, err := r.ObserverClient.NodeAccountAll(r.Ctx, &observertypes.QueryAllNodeAccountRequest{})
	require.NoError(r, err)
	require.NotEmpty(r, nodeAccounts.NodeAccount)

	existing := nodeAccounts.NodeAccount[0]
	require.NotNil(r, existing.GranteePubkey)

	msgReset := observertypes.NewMsgAddObserver(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		existing.Operator,
		existing.GranteePubkey.Secp256k1.String(),
		true,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgReset)
	require.NoError(r, err)

	// The record is now in the state that took mainnet down.
	keygenAfter, err := r.ObserverClient.Keygen(r.Ctx, &observertypes.QueryGetKeygenRequest{})
	require.NoError(r, err)
	require.NotNil(r, keygenAfter.Keygen)
	require.Empty(r, keygenAfter.Keygen.GranteePubkeys, "grantee list should have been erased")
	require.Equal(r, observertypes.KeygenStatus_PendingKeygen, keygenAfter.Keygen.Status)
	require.Equal(r, int64(math.MaxInt64), keygenAfter.Keygen.BlockNumber)

	// The same handler pauses inbound, so restore it before expecting a deposit to flow.
	reEnableInbound(r)

	// The real assertion: a withdraw needs a TSS keysign, so completing one proves the
	// signers stayed up and kept using the key they already held.
	r.ApproveETHZRC20(r.GatewayZEVMAddr)
	tx := r.ETHWithdraw(r.EVMAddress(), amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw after keygen reset")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// Signed by the original key, not a replacement.
	tssAfter, err := r.ObserverClient.TSS(r.Ctx, &observertypes.QueryGetTSSRequest{})
	require.NoError(r, err)
	require.Equal(r, tssBefore.TSS.TssPubkey, tssAfter.TSS.TssPubkey, "TSS key must not have rotated")

	// Still reset: nothing repaired the record, so signing survived it rather than dodged it.
	keygenEnd, err := r.ObserverClient.Keygen(r.Ctx, &observertypes.QueryGetKeygenRequest{})
	require.NoError(r, err)
	require.Empty(r, keygenEnd.Keygen.GranteePubkeys)
}
