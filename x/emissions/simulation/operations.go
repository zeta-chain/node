package simulation

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/zeta-chain/node/x/emissions/keeper"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
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
	DefaultWeightMsgWithdrawEmissionType = 100
	DefaultWeightMsgUpdateParams         = 100

	OpWeightMsgWithdrawEmissionType = "op_weight_msg_withdraw_emission_type"
	OpWeightMsgUpdateParams         = "op_weight_msg_update_params"
)

func WeightedOperations(
	appParams simtypes.AppParams, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgWithdrawEmissionType int
		weightMsgUpdateParams         int
	)

	appParams.GetOrGenerate(OpWeightMsgWithdrawEmissionType, &weightMsgWithdrawEmissionType, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawEmissionType = DefaultWeightMsgWithdrawEmissionType
		})

	appParams.GetOrGenerate(OpWeightMsgUpdateParams, &weightMsgUpdateParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateParams = DefaultWeightMsgUpdateParams
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgWithdrawEmissionType,
			SimulateMsgWithdrawEmissions(k),
		),
	}
}

// GetRandomAccountAndObserver returns a random account and the associated observer address
func GetRandomAccountAndObserver(
	r *rand.Rand,
	ctx sdk.Context,
	k keeper.Keeper,
	accounts []simtypes.Account,
) (simtypes.Account, string, error) {
	observers, found := k.GetObserverKeeper().GetObserverSet(ctx)
	if !found {
		return simtypes.Account{}, "", fmt.Errorf("observer set not found")
	}

	if len(observers.ObserverList) == 0 {
		return simtypes.Account{}, "", fmt.Errorf("no observers present in observer set found")
	}

	randomObserver := ""
	foundObserver := false
	for i := 0; i < 10; i++ {
		randomObserver = GetRandomObserver(r, observers.ObserverList)
		ok := k.GetObserverKeeper().IsNonTombstonedObserver(ctx, randomObserver)
		if ok {
			foundObserver = true
			break
		}
	}

	if !foundObserver {
		return simtypes.Account{}, "", fmt.Errorf("no observer found")
	}

	simAccount, err := GetObserverAccount(randomObserver, accounts)
	if err != nil {
		return simtypes.Account{}, "", err
	}
	return simAccount, randomObserver, nil
}

func GetRandomObserver(r *rand.Rand, observerList []string) string {
	idx := r.Intn(len(observerList))
	return observerList[idx]
}

// GetObserverAccount returns the account associated with the observer address from the list of accounts provided
// GetObserverAccount can fail if all the observers are removed from the observer set ,this can happen
//if the other modules create transactions which affect the validator
//and triggers any of the staking hooks defined in the observer modules

func GetObserverAccount(observerAddress string, accounts []simtypes.Account) (simtypes.Account, error) {
	operatorAddress, err := observerTypes.GetOperatorAddressFromAccAddress(observerAddress)
	if err != nil {
		return simtypes.Account{}, fmt.Errorf("validator not found for observer ")
	}

	simAccount, found := simtypes.FindAccount(accounts, operatorAddress)
	if !found {
		return simtypes.Account{}, fmt.Errorf("operator account not found")
	}
	return simAccount, nil
}
