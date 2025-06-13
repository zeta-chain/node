package simulation

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// GetRandomAccountAndObserver returns a random account and the associated observer address
func GetRandomAccountAndObserver(
	r *rand.Rand,
	ctx sdk.Context,
	k ObserverKeeper,
	accounts []simtypes.Account,
) (simtypes.Account, string, []string, error) {
	observerList := []string{}
	observers, found := k.GetObserverSet(ctx)
	if !found {
		return simtypes.Account{}, "", observerList, fmt.Errorf("observer set not found")
	}

	observerList = observers.ObserverList

	if len(observers.ObserverList) == 0 {
		return simtypes.Account{}, "", observerList, fmt.Errorf("no observers present in observer set found")
	}

	randomObserver := ""
	foundObserver := RepeatCheck(func() bool {
		randomObserver = GetRandomObserver(r, observerList)
		_, foundNodeAccount := k.GetNodeAccount(ctx, randomObserver)
		if !foundNodeAccount {
			return false
		}
		return k.CheckObserverCanVote(ctx, randomObserver) == nil
	})

	if !foundObserver {
		return simtypes.Account{}, "", nil, fmt.Errorf("no observer found")
	}

	simAccount, err := GetSimAccount(randomObserver, accounts)
	if err != nil {
		return simtypes.Account{}, "", observerList, err
	}
	return simAccount, randomObserver, observerList, nil
}

// GetRandomNodeAccount returns a random node account and the associated simulation account
func GetRandomNodeAccount(
	r *rand.Rand,
	ctx sdk.Context,
	k ObserverKeeper,
	accounts []simtypes.Account,
) (simtypes.Account, string, error) {
	nodeAccounts := k.GetAllNodeAccount(ctx)

	if len(nodeAccounts) == 0 {
		return simtypes.Account{}, "", fmt.Errorf("no node accounts present")
	}

	randomNodeAccount := nodeAccounts[r.Intn(len(nodeAccounts))].Operator

	simAccount, err := GetSimAccount(randomNodeAccount, accounts)
	if err != nil {
		return simtypes.Account{}, "", err
	}
	return simAccount, randomNodeAccount, nil
}

// GetRandomObserver returns a random observer address from the list of observers
func GetRandomObserver(r *rand.Rand, observerList []string) string {
	idx := r.Intn(len(observerList))
	return observerList[idx]
}

// GetObserverAccount returns the simulation account associated with the observer address
func GetObserverAccount(observerAddress string, accounts []simtypes.Account) (simtypes.Account, error) {
	operatorAddress, err := observertypes.GetOperatorAddressFromAccAddress(observerAddress)
	if err != nil {
		return simtypes.Account{}, fmt.Errorf("validator not found for observer ")
	}

	simAccount, found := simtypes.FindAccount(accounts, operatorAddress)
	if !found {
		return simtypes.Account{}, fmt.Errorf("operator account not found")
	}
	return simAccount, nil
}
