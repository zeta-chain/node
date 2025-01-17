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

func NewGetInboundBallotCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "get-ballot [inboundHash] [chainID]",
		Short: "fetch ballot identifier from the inbound hash",
		RunE:  GetInboundBallot,
		Args:  cobra.ExactArgs(2),
	}
}

func GetInboundBallot(cmd *cobra.Command, args []string) error {
	inboundHash := args[0]
	inboundChainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chain id")
	}
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s , err %w", config.FlagConfig, err)
	}

	return GetBallotIdentifier(inboundHash, inboundChainID, configFile)
}

func GetBallotIdentifier(inboundHash string, inboundChainID int64, configFile string) error {
	observationChain, found := chains.GetChainFromChainID(inboundChainID, []chains.Chain{})
	if !found {
		return fmt.Errorf("chain not supported,chain id: %d", inboundChainID)
	}

	cfg, err := config.GetConfig(observationChain, configFile)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	zetacoreClient, err := zetacorerpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return fmt.Errorf("failed to create zetacore client: %w", err)
	}

	ctx := context.Background()
	ballotIdentifierMessage := ""

	if observationChain.IsEVMChain() {
		ballotIdentifierMessage, err = evmInboundBallotIdentifier(
			ctx,
			*cfg,
			zetacoreClient,
			inboundHash,
			observationChain,
			cfg.ZetaChainID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get inbound ballot for evm chain %d, %w",
				observationChain.ChainId,
				err,
			)
		}
	}

	if observationChain.IsBitcoinChain() {
		ballotIdentifierMessage, err = btcInboundBallotIdentifier(
			*cfg,
			zetacoreClient,
			inboundHash,
			observationChain,
			cfg.ZetaChainID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get inbound ballot for bitcoin chain %d, %w",
				observationChain.ChainId,
				err,
			)
		}
	}

	if observationChain.IsSolanaChain() {
		ballotIdentifierMessage, err = solanaInboundBallotIdentifier(
			ctx,
			*cfg,
			zetacoreClient,
			inboundHash,
			observationChain,
			cfg.ZetaChainID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get inbound ballot for solana chain %d, %w",
				observationChain.ChainId,
				err,
			)
		}
	}

	log.Info().Msgf(ballotIdentifierMessage)
	return nil
}
