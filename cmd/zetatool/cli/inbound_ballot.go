package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/node/cmd/zetatool/cctx"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	zetacontext "github.com/zeta-chain/node/cmd/zetatool/context"
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

	ctx, err := zetacontext.NewContext(context.Background(), inboundChainID, inboundHash, configFile)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	cctxDetails := cctx.NewCCTXDetails()

	err = cctxDetails.CheckInbound(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ballot identifier: %w", err)
	}
	if cctxDetails.Status == cctx.PendingInboundConfirmation {
		log.Print("Ballot Identifier: %s, warning the inbound hash might not be confirmed yet", cctxDetails.CCCTXIdentifier)
		return nil
	}
	log.Print("Ballot Identifier: ", cctxDetails.CCCTXIdentifier)
	return nil
}
