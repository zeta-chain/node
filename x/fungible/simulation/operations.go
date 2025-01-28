package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	zetasimulation "github.com/zeta-chain/node/testutil/simulation"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

var TypeMsgDeploySystemContracts = sdk.MsgTypeURL(&types.MsgDeploySystemContracts{})

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
	// #nosec G101 not a hardcoded credential
	OpWeightMsgDeploySystemContracts      = "op_weight_msg_deploy_system_contracts"
	DefaultWeightMsgDeploySystemContracts = 5
)

// DeployedSystemContracts Use a flag to ensure that the system contracts are deployed only once
// https://github.com/zeta-chain/node/issues/3102
func WeightedOperations(
	appParams simtypes.AppParams, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgDeploySystemContracts int

	appParams.GetOrGenerate(OpWeightMsgDeploySystemContracts, &weightMsgDeploySystemContracts, nil,
		func(_ *rand.Rand) {
			weightMsgDeploySystemContracts = DefaultWeightMsgDeploySystemContracts
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgDeploySystemContracts,
			SimulateMsgDeploySystemContracts(k),
		),
	}
}

// SimulateMsgDeploySystemContracts deploy system contracts.It is run only once in first block.
func SimulateMsgDeploySystemContracts(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := zetasimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgDeploySystemContracts, err.Error()), nil, nil
		}

		msg := types.MsgDeploySystemContracts{Creator: policyAccount.Address.String()}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgDeploySystemContracts,
				"failed to validate basic msg",
			), nil, err
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           &msg,
			Context:       ctx,
			SimAccount:    policyAccount,
			AccountKeeper: k.GetAuthKeeper(),
			Bankkeeper:    k.GetBankKeeper(),
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
