package simulation

import (
	"fmt"
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

const (
	OpWeightMsgTypeMsgEnableCCTX      = "op_weight_msg_enable_crosschain_flags"
	DefaultWeightMsgTypeMsgEnableCCTX = 650
)

// WeightedOperations for observer module
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper,
) simulation.WeightedOperations {
	var weightMsgTypeMsgEnableCCTX int

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgEnableCCTX, &weightMsgTypeMsgEnableCCTX, nil,
		func(r *rand.Rand) {
			weightMsgTypeMsgEnableCCTX = DefaultWeightMsgTypeMsgEnableCCTX
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgTypeMsgEnableCCTX,
			SimulateMsgTypeMsgEnableCCTX(k),
		),
	}
}

func SimulateMsgTypeMsgEnableCCTX(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, chainID string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {

		policies, found := k.GetAuthorityKeeper().GetPolicies(ctx)
		if !found {
			fmt.Println("Policies not found")
		}

		// TODO setup simutils package
		admin := policies.Items[0].Address
		address, err := types.GetOperatorAddressFromAccAddress(admin)
		if err != nil {
			panic(err)
		}
		simAccount, found := simtypes.FindAccount(accounts, address)
		if !found {
			// TODO : remove panic
			panic("admin account not found")
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
