package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

func CreateAuthzSigner(granter string, grantee sdk.AccAddress) {
	authz.SetupAuthZSignerList(granter, grantee)
}

func CreateZetacoreClient(cfg config.Config, hotkeyPassword string, logger zerolog.Logger) (*zetacore.Client, error) {
	hotKey := cfg.AuthzHotkey
	if cfg.HsmMode {
		hotKey = cfg.HsmHotKey
	}

	chainIP := cfg.ZetaCoreURL

	kb, _, err := keys.GetKeyringKeybase(cfg, hotkeyPassword)
	if err != nil {
		return nil, err
	}

	granterAddreess, err := sdk.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, err
	}

	k := keys.NewKeysWithKeybase(kb, granterAddreess, cfg.AuthzHotkey, hotkeyPassword)

	client, err := zetacore.NewClient(k, chainIP, hotKey, cfg.ChainID, cfg.HsmMode, logger)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// TODO
// 	// create Solana chain observer
//	solChain, solConfig, enabled := appContext.GetSolanaChainAndConfig()
//	if enabled {
//		rpcClient := solrpc.New(solConfig.Endpoint)
//		if rpcClient == nil {
//			// should never happen
//			logger.Std.Error().Msg("solana create Solana client error")
//			return observerMap, nil
//		}
//
//		observer, err := solanaobserver.NewObserver(
//			solChain,
//			rpcClient,
//			*solChainParams,
//			zetacoreClient,
//			tss,
//			dbpath,
//			logger,
//			ts,
//		)
//		if err != nil {
//			logger.Std.Error().Err(err).Msg("NewObserver error for solana chain")
//		} else {
//			observerMap[solChainParams.ChainId] = observer
//		}
//	}
