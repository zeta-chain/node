package simulation

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
)

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

// GetRandomChainID returns a random chainID from the list of chains
func GetRandomChainID(r *rand.Rand, chains []chains.Chain) int64 {
	idx := r.Intn(len(chains))
	return chains[idx].ChainId
}
