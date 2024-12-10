package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
)

// Simulation operation weights constants
// Operation weights are used by the `SimulateFromSeed`
// function to pick a random operation based on the weights.The functions with higher weights are more likely to be picked.

// Therefore, this decides the percentage of a certain operation that is part of a block.

// Based on the weights assigned in the cosmos sdk modules,
// 100 seems to the max weight used,and we should use relative weights
// to signify the number of each operation in a block.

// TODO Add more details to comment based on what the number represents in terms of percentage of operations in a block
// https://github.com/zeta-chain/node/issues/3100
const (
	// #nosec G101 not a hardcoded credential
	OpWeightMsgDeploySystemContracts      = "op_weight_msg_deploy_system_contracts"
	DefaultWeightMsgDeploySystemContracts = 5
)

// DeployedSystemContracts Use a flag to ensure that the system contracts are deployed only once
// https://github.com/zeta-chain/node/issues/3102
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgDeploySystemContracts int

	appParams.GetOrGenerate(cdc, OpWeightMsgDeploySystemContracts, &weightMsgDeploySystemContracts, nil,
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
		policies, found := k.GetAuthorityKeeper().GetPolicies(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgDeploySystemContracts,
				"policies object not found",
			), nil, nil
		}
		if len(policies.Items) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgDeploySystemContracts,
				"no policies found",
			), nil, nil
		}
		admin := policies.Items[0].Address

		address, err := observerTypes.GetOperatorAddressFromAccAddress(admin)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgDeploySystemContracts,
				"unable to get operator address",
			), nil, err
		}
		simAccount, found := simtypes.FindAccount(accounts, address)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgDeploySystemContracts,
				"sim account for admin address not found",
			), nil, nil
		}

		msg := types.MsgDeploySystemContracts{Creator: admin}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "failed to validate basic msg"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           &msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: k.GetAuthKeeper(),
			Bankkeeper:    k.GetBankKeeper(),
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
