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
// Operation weights are used by the simulation program to simulate the weight of different operations.
// This decides what percentage of a certain type of operation is part of a block.
// Based on the weights assigned in the cosmos sdk modules , 100 seems to the max weight used , and therefore guarantees that at least one operation of that type is present in a block.
// TODO Add more details to comment based on what the number represents in terms of percentage of operations in a block
// https://github.com/zeta-chain/node/issues/3100
const (
	OpWeightMsgDeploySystemContracts      = "op_weight_msg_deploy_system_contracts"
	DefaultWeightMsgDeploySystemContracts = 100
)

// DeployedSystemContracts Use a flag to ensure that the system contracts are deployed only once
// https://github.com/zeta-chain/node/issues/3102
var (
	DeployedSystemContracts = false
)

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgDeploySystemContracts int

	appParams.GetOrGenerate(cdc, OpWeightMsgDeploySystemContracts, &weightMsgDeploySystemContracts, nil,
		func(r *rand.Rand) {
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
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, chainID string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		if DeployedSystemContracts {
			return simtypes.OperationMsg{}, nil, nil
		}

		policies, found := k.GetAuthorityKeeper().GetPolicies(ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDeploySystemContracts, "policies not found"), nil, nil
		}

		admin := policies.Items[0].Address

		address, err := observerTypes.GetOperatorAddressFromAccAddress(admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDeploySystemContracts, "unable to get operator address"), nil, err
		}
		simAccount, found := simtypes.FindAccount(accounts, address)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDeploySystemContracts, "sim account for admin address not found"), nil, nil
		}

		msg := types.MsgDeploySystemContracts{Creator: admin}

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

		DeployedSystemContracts = true
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
