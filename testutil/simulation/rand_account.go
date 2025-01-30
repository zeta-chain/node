package simulation

import (
	"fmt"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const DefaultRetryCount = 10

// GetSimAccount returns the simulation account associated with the observer address
func GetSimAccount(observerAddress string, accounts []simtypes.Account) (simtypes.Account, error) {
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

// RepeatCheck checks the function for a number of times and returns true if the function returns true
func RepeatCheck(fn func() bool) bool {
	for i := 0; i < DefaultRetryCount; i++ {
		if fn() {
			return true
		}
	}
	return false
}
