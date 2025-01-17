package inbound

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	zetacorerpc "github.com/zeta-chain/node/pkg/rpc"
)

func NewFetchInboundBallotCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "get-ballot [inboundHash] [chainID]",
		Short: "fetch ballot identifier from the inbound hash",
		RunE:  GetInboundBallot,
	}
}

func GetInboundBallot(cmd *cobra.Command, args []string) error {
	cobra.ExactArgs(2)

	inboundHash := args[0]
	inboundChainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chain id")
	}
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s , err %s", config.FlagConfig, err.Error())
	}

	return GetBallotIdentifier(inboundHash, inboundChainID, configFile)
}

func GetBallotIdentifier(inboundHash string, inboundChainID int64, configFile string) error {
	observationChain, found := chains.GetChainFromChainID(inboundChainID, []chains.Chain{})
	if !found {
		return fmt.Errorf("chain not supported,chain id : %d", inboundChainID)
	}

	cfg, err := config.GetConfig(observationChain, configFile)
	if err != nil {
		return fmt.Errorf("failed to get config, %s", err.Error())
	}

	zetacoreClient, err := zetacorerpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return fmt.Errorf("failed to create zetacore client, %s", err.Error())
	}

	ctx := context.Background()
	ballotIdentifier := ""

	if observationChain.IsEVMChain() {
		ballotIdentifier, err = evmInboundBallotIdentifier(
			ctx,
			*cfg,
			zetacoreClient,
			inboundHash,
			observationChain,
			cfg.ZetaChainID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get inbound ballot for evm chain %d, %s",
				observationChain.ChainId,
				err.Error(),
			)
		}
	}

	if observationChain.IsBitcoinChain() {
		ballotIdentifier, err = btcInboundBallotIdentifier(
			*cfg,
			zetacoreClient,
			inboundHash,
			observationChain,
			cfg.ZetaChainID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get inbound ballot for bitcoin chain %d, %s",
				observationChain.ChainId,
				err.Error(),
			)
		}
	}

	if observationChain.IsSolanaChain() {
		ballotIdentifier, err = solanaInboundBallotIdentifier(
			ctx,
			*cfg,
			zetacoreClient,
			inboundHash,
			observationChain,
			cfg.ZetaChainID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get inbound ballot for solana chain %d, %s",
				observationChain.ChainId,
				err.Error(),
			)
		}
	}

	log.Info().Msgf("Ballot Identifier: %s", ballotIdentifier)
	return nil
}
