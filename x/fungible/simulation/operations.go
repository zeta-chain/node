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
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	OpWeightMsgDeploySystemContracts      = "op_weight_msg_deploy_system_contracts"
	DefaultWeightMsgDeploySystemContracts = 100
)

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

func SimulateMsgDeploySystemContracts(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, chainID string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		if DeployedSystemContracts {
			return simtypes.OperationMsg{}, nil, nil
		}

		policies, found := k.GetAuthorityKeeper().GetPolicies(ctx)
		if !found {
			fmt.Println("Policies not found")
		}

		admin := policies.Items[0].Address
		address, err := observerTypes.GetOperatorAddressFromAccAddress(admin)
		if err != nil {
			panic(err)
		}
		simAccount, found := simtypes.FindAccount(accounts, address)
		if !found {
			// TODO : remove panic
			panic("admin account not found")
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
