package main

import (
	"context"
	"os"
	"slices"
	"strconv"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// isObserverNode checks whether THIS node is an observer node.
func isObserverNode(ctx context.Context, zc *zetacore.Client) (bool, error) {
	observers, err := zc.GetObserverList(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to get observers list")
	}

	operatorAddress := zc.GetKeys().GetOperatorAddress().String()

	if slices.Contains(observers, operatorAddress) {
		return true, nil
	}

	return false, nil
}

func isEnvFlagEnabled(flag string) bool {
	v, _ := strconv.ParseBool(os.Getenv(flag))
	return v
}

func btcChainIDsFromContext(app *zctx.AppContext) []int64 {
	var (
		btcChains   = app.FilterChains(zctx.Chain.IsBitcoin)
		btcChainIDs = make([]int64, len(btcChains))
	)

	for i, chain := range btcChains {
		btcChainIDs[i] = chain.ID()
	}

	return btcChainIDs
}

func resolveObserverPubKeyBech32(cfg config.Config, hotKeyPassword string) (string, error) {
	// Get observer's public key ("grantee pub key")
	_, granteePubKeyBech32, err := keys.GetKeyringKeybase(cfg, hotKeyPassword)
	if err != nil {
		return "", errors.Wrap(err, "unable to get keyring key base")
	}

	return granteePubKeyBech32, nil
}
