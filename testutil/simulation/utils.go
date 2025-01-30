package simulation

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const DefaultRetryCount = 10

// GetPolicyAccount returns the simulation account associated with the policy address
func GetPolicyAccount(ctx sdk.Context, k AuthorityKeeper, accounts []simtypes.Account) (simtypes.Account, error) {
	policies, found := k.GetPolicies(ctx)
	if !found {
		return simtypes.Account{}, fmt.Errorf("policies object not found")
	}
	if len(policies.Items) == 0 {
		return simtypes.Account{}, fmt.Errorf("no policies found")
	}

	admin := policies.Items[0].Address
	address, err := observertypes.GetOperatorAddressFromAccAddress(admin)
	if err != nil {
		return simtypes.Account{}, err
	}
	simAccount, found := simtypes.FindAccount(accounts, address)
	if !found {
		return simtypes.Account{}, fmt.Errorf("admin account not found in list of simulation accounts")
	}
	return simAccount, nil
}

// GetAsset returns the asset associated with the chainID
func GetAsset(ctx sdk.Context, k FungibleKeeper, chainID int64) (string, error) {
	foreignCoins := k.GetAllForeignCoins(ctx)
	asset := ""

	for _, coin := range foreignCoins {
		if coin.ForeignChainId == chainID {
			return coin.Asset, nil
		}
	}

	return asset, fmt.Errorf("asset not found for chain %d", chainID)
}

// GetExternalChain returns a random external chain from the list of supported chains
func GetExternalChain(ctx sdk.Context, k ObserverKeeper, r *rand.Rand) (chains.Chain, error) {
	supportedChains := k.GetSupportedChains(ctx)
	if len(supportedChains) == 0 {
		return chains.Chain{}, fmt.Errorf("no supported chains found")
	}
	externalChain := chains.Chain{}
	foundExternalChain := RepeatCheck(func() bool {
		c := supportedChains[r.Intn(len(supportedChains))]
		if !c.IsZetaChain() {
			externalChain = c
			return true
		}
		return false
	})

	if !foundExternalChain {
		return chains.Chain{}, fmt.Errorf("no external chain found")
	}
	return externalChain, nil
}

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
		ok := k.IsNonTombstonedObserver(ctx, randomObserver)
		if ok {
			return true
		}
		return false
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

// GetRandomChainID returns a random chainID from the list of chains
func GetRandomChainID(r *rand.Rand, chains []chains.Chain) int64 {
	idx := r.Intn(len(chains))
	return chains[idx].ChainId
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

// RepeatCheck checks the function for a number of times and returns true if the function returns true
func RepeatCheck(fn func() bool) bool {
	for i := 0; i < DefaultRetryCount; i++ {
		if fn() {
			return true
		}
	}
	return false
}
