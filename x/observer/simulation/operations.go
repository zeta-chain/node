package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/x/observer/keeper"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

var (
	TypeMsgEnableCCTX                  = sdk.MsgTypeURL(&observertypes.MsgEnableCCTX{})
	TypeMsgDisableCCTX                 = sdk.MsgTypeURL(&observertypes.MsgDisableCCTX{})
	TypeMsgVoteTSS                     = sdk.MsgTypeURL(&observertypes.MsgVoteTSS{})
	TypeMsgUpdateKeygen                = sdk.MsgTypeURL(&observertypes.MsgUpdateKeygen{})
	TypeMsgUpdateObserver              = sdk.MsgTypeURL(&observertypes.MsgUpdateObserver{})
	TypeMsgUpdateChainParams           = sdk.MsgTypeURL(&observertypes.MsgUpdateChainParams{})
	TypeMsgRemoveChainParams           = sdk.MsgTypeURL(&observertypes.MsgRemoveChainParams{})
	TypeMsgResetChainNonces            = sdk.MsgTypeURL(&observertypes.MsgResetChainNonces{})
	TypeMsgUpdateGasPriceIncreaseFlags = sdk.MsgTypeURL(&observertypes.MsgUpdateGasPriceIncreaseFlags{})
	TypeMsgAddObserver                 = sdk.MsgTypeURL(&observertypes.MsgAddObserver{})
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
	OpWeightMsgTypeMsgEnableCCTX                  = "op_weight_msg_enable_crosschain_flags"         // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgDisableCCTX                 = "op_weight_msg_disable_crosschain_flags"        // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgVoteTSS                     = "op_weight_msg_vote_tss"                        // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateKeygen                = "op_weight_msg_update_keygen"                   // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateObserver              = "op_weight_msg_update_observer"                 // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateChainParams           = "op_weight_msg_update_chain_params"             // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgRemoveChainParams           = "op_weight_msg_remove_chain_params"             // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgResetChainNonces            = "op_weight_msg_reset_chain_nonces"              // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateGasPriceIncreaseFlags = "op_weight_msg_update_gas_price_increase_flags" // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgAddObserver                 = "op_weight_msg_add_observer"                    // #nosec G101 not a hardcoded credential

	// DefaultWeightMsgTypeMsgEnableCCTX We use a high weight for this operation
	// to ensure that it is present in the block more number of times than any operation that changes the validator set
	// Arrived at this number based on the weights used in the cosmos sdk staking module and through some trial and error
	DefaultWeightMsgTypeMsgEnableCCTX                  = 100
	DefaultWeightMsgTypeMsgDisableCCTX                 = 10
	DefaultWeightMsgTypeMsgVoteTSS                     = 10
	DefaultWeightMsgTypeMsgUpdateKeygen                = 10
	DefaultWeightMsgTypeMsgUpdateObserver              = 10
	DefaultWeightMsgTypeMsgUpdateChainParams           = 10
	DefaultWeightMsgTypeMsgRemoveChainParams           = 10
	DefaultWeightMsgTypeMsgResetChainNonces            = 5
	DefaultWeightMsgTypeMsgUpdateGasPriceIncreaseFlags = 10
	DefaultWeightMsgTypeMsgAddObserver                 = 5
	moduleName                                         = observertypes.ModuleName
)

// WeightedOperations for observer module
func WeightedOperations(
	appParams simtypes.AppParams, k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgTypeMsgEnableCCTX                  int
		weightMsgTypeMsgDisableCCTX                 int
		weightMsgTypeMsgVoteTSS                     int
		weightMsgTypeMsgUpdateKeygen                int
		weightMsgTypeMsgUpdateObserver              int
		weightMsgTypeMsgUpdateChainParams           int
		weightMsgTypeMsgRemoveChainParams           int
		weightMsgTypeMsgResetChainNonces            int
		weightMsgTypeMsgUpdateGasPriceIncreaseFlags int
		weightMsgTypeMsgAddObserver                 int
	)

	appParams.GetOrGenerate(OpWeightMsgTypeMsgEnableCCTX, &weightMsgTypeMsgEnableCCTX, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgEnableCCTX = DefaultWeightMsgTypeMsgEnableCCTX
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgDisableCCTX, &weightMsgTypeMsgDisableCCTX, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgDisableCCTX = DefaultWeightMsgTypeMsgDisableCCTX
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgVoteTSS, &weightMsgTypeMsgVoteTSS, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgVoteTSS = DefaultWeightMsgTypeMsgVoteTSS
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgUpdateKeygen, &weightMsgTypeMsgUpdateKeygen, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateKeygen = DefaultWeightMsgTypeMsgUpdateKeygen
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgUpdateObserver, &weightMsgTypeMsgUpdateObserver, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateObserver = DefaultWeightMsgTypeMsgUpdateObserver
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgUpdateChainParams, &weightMsgTypeMsgUpdateChainParams, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateChainParams = DefaultWeightMsgTypeMsgUpdateChainParams
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgRemoveChainParams, &weightMsgTypeMsgRemoveChainParams, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgRemoveChainParams = DefaultWeightMsgTypeMsgRemoveChainParams
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgResetChainNonces, &weightMsgTypeMsgResetChainNonces, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgResetChainNonces = DefaultWeightMsgTypeMsgResetChainNonces
		})

	appParams.GetOrGenerate(
		OpWeightMsgTypeMsgUpdateGasPriceIncreaseFlags,
		&weightMsgTypeMsgUpdateGasPriceIncreaseFlags,
		nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateGasPriceIncreaseFlags = DefaultWeightMsgTypeMsgUpdateGasPriceIncreaseFlags
		},
	)

	appParams.GetOrGenerate(OpWeightMsgTypeMsgAddObserver, &weightMsgTypeMsgAddObserver, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgAddObserver = DefaultWeightMsgTypeMsgAddObserver
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgTypeMsgEnableCCTX,
			SimulateEnableCCTX(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgDisableCCTX,
			SimulateDisableCCTX(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateKeygen,
			SimulateUpdateKeygen(k),
		),
		//
		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateChainParams,
			SimulateUpdateChainParams(k),
		),
		//
		//simulation.NewWeightedOperation(
		//	weightMsgTypeMsgRemoveChainParams,
		//	SimulateMsgRemoveChainParams(k),
		//),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgResetChainNonces,
			SimulateResetChainNonces(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateGasPriceIncreaseFlags,
			SimulateUpdateGasPriceIncreaseFlags(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateUpdateObserver(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateAddObserverNodeAccount(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateAddObserver(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgVoteTSS,
			SimulateMsgVoteTSS(k),
		),
	}
}
