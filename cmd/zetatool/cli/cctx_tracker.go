package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/node/cmd/zetatool/ballot"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	zetatoolcontext "github.com/zeta-chain/node/cmd/zetatool/context"
)

func NewTrackCCTXCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "track-cctx [inboundHash] [chainID]",
		Short: "track a cross chain transaction",
		RunE:  TrackCCTX,
		Args:  cobra.ExactArgs(2),
	}
}

func TrackCCTX(cmd *cobra.Command, args []string) error {
	inboundHash := args[0]
	inboundChainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chain id")
	}
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s , err %w", config.FlagConfig, err)
	}

	ctx, err := zetatoolcontext.NewContext(context.Background(), inboundChainID, inboundHash, configFile)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	cctx, err := ballot.GetBallotIdentifier(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ballot identifier: %w", err)
	}

	return nil
}
