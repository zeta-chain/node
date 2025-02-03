package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/cctx"
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

	cctxDetails, err := trackCCTX(ctx)
	if err != nil {
		return fmt.Errorf("failed to track cctx: %w", err)
	}
	if cmd.Flag(config.FlagDebug).Changed {
		log.Info().Msg(cctxDetails.DebugPrint())
		return nil
	}
	log.Info().Msg(cctxDetails.Print())
	return nil
}

func trackCCTX(ctx *zetatoolcontext.Context) (*cctx.TrackingDetails, error) {
	var (
		cctxDetails = cctx.NewCCTXDetails()
		err         error
	)
	// Get the ballot identifier for the inbound transaction and confirm that cctx status in atleast either PendingInboundConfirmation or PendingInboundVoting
	err = cctxDetails.CheckInbound(ctx)
	if err != nil {
		return cctxDetails, fmt.Errorf("failed to get ballot identifier: %v", err)
	}
	// Reject unknown status, as it is not valid
	if cctxDetails.Status == cctx.Unknown || cctxDetails.CCTXIdentifier == "" {
		return cctxDetails, fmt.Errorf("unknown status")
	}

	// At this point, we have confirmed the inbound hash is valid, and it was sent to valid address.
	// Update cctx status from zetacore.This copies the status from zetacore to the cctx details.The cctx status can only be `PendingInboundVoting` or `PendingInboundConfirmation` at this point
	cctxDetails.UpdateCCTXStatus(ctx)

	// The cctx details now have status from zetacore, we have not tried to a get more granular status from the outbound chain yet.
	// If it's not pending, we can just return here.
	if !cctxDetails.IsPendingOutbound() {
		return cctxDetails, nil
	}

	// update outbound details, this does not transition any status.
	cctxDetails.UpdateCCTXOutboundDetails(ctx)

	// Update tx hash list from outbound tracker
	// If the tracker is found, it means the outbound is broadcast, but we are waiting for the confirmations
	// If the tracker is not found, it means the outbound is not broadcast yet, we are wwaiting for the tss to sign the outbound
	cctxDetails.UpdateHashListAndPendingStatus(ctx)

	// If its not pending confirmation, we can return here, it means the outbound is not broadcast yet its pending tss signing
	if !cctxDetails.IsPendingConfirmation() {
		return cctxDetails, nil
	}

	// Check outbound tx, we are waiting for the outbound tx to be confirmed
	err = cctxDetails.CheckOutbound(ctx)
	if err != nil {
		return cctxDetails, err
	}
	return cctxDetails, nil
}
