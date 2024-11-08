package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// Simulation operation weights constants
// Operation weights are used by the simulation program to simulate the weight of different operations.
// This decides what percentage of a certain type of operation is part of a block.
// Based on the weights assigned in the cosmos sdk modules , 100 seems to the max weight used , and therefore guarantees that at least one operation of that type is present in a block.
// TODO Add more details to comment based on what the number represents in terms of percentage of operations in a block
// https://github.com/zeta-chain/node/issues/3100
const (
	// #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgEnableCCTX = "op_weight_msg_enable_crosschain_flags"
	// DefaultWeightMsgTypeMsgEnableCCTX We ues a high weight for this operation
	// to ensure that it is present in the block more number of times than any operation that changes the validator set

	// Arrived at this number based on the weights used in the cosmos sdk staking module and through some trial and error
	DefaultWeightMsgTypeMsgEnableCCTX = 3650
)

// WeightedOperations for observer module
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper,
) simulation.WeightedOperations {
	var weightMsgTypeMsgEnableCCTX int

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgEnableCCTX, &weightMsgTypeMsgEnableCCTX, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgEnableCCTX = DefaultWeightMsgTypeMsgEnableCCTX
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgTypeMsgEnableCCTX,
			SimulateMsgTypeMsgEnableCCTX(k),
		),
	}
}

// SimulateMsgTypeMsgEnableCCTX generates a MsgEnableCCTX and delivers it.
func SimulateMsgTypeMsgEnableCCTX(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policies, found := k.GetAuthorityKeeper().GetPolicies(ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableCCTX, "policies object not found"), nil, nil
		}
		if len(policies.Items) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableCCTX, "no policies found"), nil, nil
		}

		admin := policies.Items[0].Address
		address, err := types.GetOperatorAddressFromAccAddress(admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableCCTX, err.Error()), nil, err
		}
		simAccount, found := simtypes.FindAccount(accounts, address)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableCCTX, "admin account not found"), nil, nil
		}

		msg := types.MsgEnableCCTX{
			Creator:        simAccount.Address.String(),
			EnableInbound:  true,
			EnableOutbound: false,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
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
