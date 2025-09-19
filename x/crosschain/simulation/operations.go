package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

var (
	TypeMsgAddOutboundTracker     = sdk.MsgTypeURL(&types.MsgAddOutboundTracker{})
	TypeMsgAddInboundTracker      = sdk.MsgTypeURL(&types.MsgAddInboundTracker{})
	TypeMsgRemoveOutboundTracker  = sdk.MsgTypeURL(&types.MsgRemoveOutboundTracker{})
	TypeMsgVoteGasPrice           = sdk.MsgTypeURL(&types.MsgVoteGasPrice{})
	TypeMsgVoteOutbound           = sdk.MsgTypeURL(&types.MsgVoteOutbound{})
	TypeMsgVoteInbound            = sdk.MsgTypeURL(&types.MsgVoteInbound{})
	TypeMsgWhitelistAsset         = sdk.MsgTypeURL(&types.MsgWhitelistAsset{})
	TypeMsgMigrateTssFunds        = sdk.MsgTypeURL(&types.MsgMigrateTssFunds{})
	TypeMsgUpdateTssAddress       = sdk.MsgTypeURL(&types.MsgUpdateTssAddress{})
	TypeMsgAbortStuckCCTX         = sdk.MsgTypeURL(&types.MsgAbortStuckCCTX{})
	TypeMsgUpdateRateLimiterFlags = sdk.MsgTypeURL(&types.MsgUpdateRateLimiterFlags{})
	TypeMsgRefundAbortedCCTX      = sdk.MsgTypeURL(&types.MsgRefundAbortedCCTX{})
)

// Simulation operation weights constants
// Operation weights are used by the simulation program to simulate the weight of different operations.
// This decides what percentage of a certain type of operation is part of a block.
// Based on the weights assigned in the cosmos sdk modules , 100 seems to the max weight used , and therefore guarantees that at least one operation of that type is present in a block.
// Operation weights are used by the `SimulateFromSeed`
// function to pick a random operation based on the weights.The functions with higher weights are more likely to be picked.

// Therefore, this decides the percentage of a certain operation that is part of a block.

// Based on the weights assigned in the cosmos sdk modules,
// 100 seems to the max weight used,and we should use relative weights
// to signify the number of each operation in a block.
const (
	DefaultWeightAddOutboundTracker     = 10
	DefaultWeightAddInboundTracker      = 10
	DefaultWeightRemoveOutboundTracker  = 10
	DefaultWeightVoteGasPrice           = 50
	DefaultWeightVoteOutbound           = 10
	DefaultWeightVoteInbound            = 10
	DefaultWeightWhitelistAsset         = 10
	DefaultWeightMigrateTssFunds        = 1
	DefaultWeightUpdateTssAddress       = 10
	DefaultWeightAbortStuckCCTX         = 5
	DefaultWeightUpdateRateLimiterFlags = 10
	DefaultWeightRefundAbortedCCTX      = 10

	OpWeightMsgAddOutboundTracker  = "op_weight_msg_add_outbound_tracker"      // #nosec G101 not a hardcoded credential
	OpWeightAddInboundTracker      = "op_weight_msg_add_inbound_tracker"       // #nosec G101 not a hardcoded credential
	OpWeightRemoveOutboundTracker  = "op_weight_msg_remove_outbound_tracker"   // #nosec G101 not a hardcoded credential
	OpWeightVoteGasPrice           = "op_weight_msg_vote_gas_price"            // #nosec G101 not a hardcoded credential
	OpWeightVoteOutbound           = "op_weight_msg_vote_outbound"             // #nosec G101 not a hardcoded credential
	OpWeightVoteInbound            = "op_weight_msg_vote_inbound"              // #nosec G101 not a hardcoded credential
	OpWeightWhitelistAsset         = "op_weight_msg_whitelist_asset"           // #nosec G101 not a hardcoded credential
	OpWeightMigrateTssFunds        = "op_weight_msg_migrate_tss_funds"         // #nosec G101 not a hardcoded credential
	OpWeightUpdateTssAddress       = "op_weight_msg_update_tss_address"        // #nosec G101 not a hardcoded credential
	OpWeightAbortStuckCCTX         = "op_weight_msg_abort_stuck_cctx"          // #nosec G101 not a hardcoded credential
	OpWeightUpdateRateLimiterFlags = "op_weight_msg_update_rate_limiter_flags" // #nosec G101 not a hardcoded credential
	OpWeightRefundAbortedCCTX      = "op_weight_msg_refund_aborted_cctx"       // #nosec G101 not a hardcoded credential
)

func WeightedOperations(
	appParams simtypes.AppParams, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightAddOutboundTracker     int
		weightAddInboundTracker      int
		weightRemoveOutboundTracker  int
		weightVoteGasPrice           int
		weightVoteOutbound           int
		weightVoteInbound            int
		weightWhitelistAsset         int
		weightMigrateTssFunds        int
		weightUpdateTssAddress       int
		weightAbortStuckCCTX         int
		weightUpdateRateLimiterFlags int
		weightRefundAbortedCCTX      int
	)

	appParams.GetOrGenerate(OpWeightMsgAddOutboundTracker, &weightAddOutboundTracker, nil,
		func(_ *rand.Rand) {
			weightAddOutboundTracker = DefaultWeightAddOutboundTracker
		},
	)

	appParams.GetOrGenerate(OpWeightAddInboundTracker, &weightAddInboundTracker, nil,
		func(_ *rand.Rand) {
			weightAddInboundTracker = DefaultWeightAddInboundTracker
		},
	)

	appParams.GetOrGenerate(OpWeightRemoveOutboundTracker, &weightRemoveOutboundTracker, nil,
		func(_ *rand.Rand) {
			weightRemoveOutboundTracker = DefaultWeightRemoveOutboundTracker
		},
	)

	appParams.GetOrGenerate(OpWeightVoteGasPrice, &weightVoteGasPrice, nil,
		func(_ *rand.Rand) {
			weightVoteGasPrice = DefaultWeightVoteGasPrice
		},
	)

	appParams.GetOrGenerate(OpWeightVoteOutbound, &weightVoteOutbound, nil,
		func(_ *rand.Rand) {
			weightVoteOutbound = DefaultWeightVoteOutbound
		},
	)

	appParams.GetOrGenerate(OpWeightVoteInbound, &weightVoteInbound, nil,
		func(_ *rand.Rand) {
			weightVoteInbound = DefaultWeightVoteInbound
		},
	)

	appParams.GetOrGenerate(OpWeightWhitelistAsset, &weightWhitelistAsset, nil,
		func(_ *rand.Rand) {
			weightWhitelistAsset = DefaultWeightWhitelistAsset
		},
	)

	appParams.GetOrGenerate(OpWeightMigrateTssFunds, &weightMigrateTssFunds, nil,
		func(_ *rand.Rand) {
			weightMigrateTssFunds = DefaultWeightMigrateTssFunds
		},
	)

	appParams.GetOrGenerate(OpWeightUpdateTssAddress, &weightUpdateTssAddress, nil,
		func(_ *rand.Rand) {
			weightUpdateTssAddress = DefaultWeightUpdateTssAddress
		},
	)

	appParams.GetOrGenerate(OpWeightAbortStuckCCTX, &weightAbortStuckCCTX, nil,
		func(_ *rand.Rand) {
			weightAbortStuckCCTX = DefaultWeightAbortStuckCCTX
		},
	)

	appParams.GetOrGenerate(OpWeightUpdateRateLimiterFlags, &weightUpdateRateLimiterFlags, nil,
		func(_ *rand.Rand) {
			weightUpdateRateLimiterFlags = DefaultWeightUpdateRateLimiterFlags
		},
	)

	appParams.GetOrGenerate(OpWeightRefundAbortedCCTX, &weightRefundAbortedCCTX, nil,
		func(_ *rand.Rand) {
			weightRefundAbortedCCTX = DefaultWeightRefundAbortedCCTX
		},
	)

	// TODO : Add the new test for MsgRemoveInboundTracker
	// https://github.com/zeta-chain/node/issues/3479

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightVoteGasPrice,
			SimulateMsgVoteGasPrice(k),
		),
		simulation.NewWeightedOperation(
			weightVoteInbound,
			SimulateVoteInbound(k),
		),
		simulation.NewWeightedOperation(
			weightVoteOutbound,
			SimulateVoteOutbound(k),
		),
		simulation.NewWeightedOperation(
			weightAddInboundTracker,
			SimulateMsgAddInboundTracker(k),
		),
		simulation.NewWeightedOperation(
			weightAddOutboundTracker,
			SimulateMsgAddOutboundTracker(k),
		),
		simulation.NewWeightedOperation(
			weightRemoveOutboundTracker,
			SimulateMsgRemoveOutboundTracker(k),
		),
		simulation.NewWeightedOperation(
			weightWhitelistAsset,
			SimulateMsgWhitelistAsset(k),
		),
		simulation.NewWeightedOperation(
			weightAbortStuckCCTX,
			SimulateMsgAbortStuckCCTX(k),
		),
		simulation.NewWeightedOperation(
			weightRefundAbortedCCTX,
			SimulateMsgRefundAbortedCCTX(k),
		),
		simulation.NewWeightedOperation(
			weightUpdateRateLimiterFlags,
			SimulateMsgUpdateRateLimiterFlags(k),
		),
		simulation.NewWeightedOperation(
			weightUpdateTssAddress,
			SimulateMsgUpdateTssAddress(k),
		),
	}
}
