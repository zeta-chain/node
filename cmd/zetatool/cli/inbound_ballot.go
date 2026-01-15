package cli

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/cctx"
	"github.com/zeta-chain/node/cmd/zetatool/common"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	zetacontext "github.com/zeta-chain/node/cmd/zetatool/context"
)

func NewGetInboundBallotCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "get-ballot [inboundHash] [chain]",
		Short: "fetch ballot identifier from the inbound hash",
		Long: `Fetch ballot identifier from the inbound hash and chain.

The chain argument can be either:
  - A chain ID (e.g., 7000, 1, 56)
  - A chain name (e.g., zeta_mainnet, eth_mainnet, bsc_mainnet)

Examples:
  zetatool get-ballot 0x1234... 7000
  zetatool get-ballot 0x1234... zeta_mainnet
  zetatool get-ballot 0x1234... eth_mainnet`,
		RunE: GetInboundBallot,
		Args: cobra.ExactArgs(2),
	}
}

func GetInboundBallot(cmd *cobra.Command, args []string) error {
	inboundHash := args[0]
	inboundChain, err := common.ResolveChain(args[1])
	if err != nil {
		return fmt.Errorf("failed to resolve chain %q: %w", args[1], err)
	}
	inboundChainID := inboundChain.ChainId
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s , err %w", config.FlagConfig, err)
	}

	ctx, err := zetacontext.NewContext(context.Background(), inboundChainID, inboundHash, configFile)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	cctxTrackingDetails := cctx.NewTrackingDetails()

	err = cctxTrackingDetails.CheckInbound(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ballot identifier: %w", err)
	}
	if cctxTrackingDetails.Status == cctx.PendingInboundConfirmation {
		log.Printf(
			"Ballot Identifier: %s, warning the inbound hash might not be confirmed yet",
			cctxTrackingDetails.CCTXIdentifier,
		)
		return nil
	}
	log.Print("Ballot Identifier: ", cctxTrackingDetails.CCTXIdentifier)
	return nil
}
