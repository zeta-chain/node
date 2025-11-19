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
	var (
		trackingDetailsList []cctx.TrackingDetails
		cctxTrackingDetails *cctx.TrackingDetails
		maxCCTXChainLength  = 5 // Maximum number of cctx chains to track
	)

	inboundHash := args[0]
	inboundChainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chain id")
	}
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s , err %w", config.FlagConfig, err)
	}
	for i := 0; i < maxCCTXChainLength; i++ {
		chainID := inboundChainID
		chainHash := inboundHash

		// if len of tx is greated than 0, we have already tracked the cctx and we are trying to continue tracking the chain
		// In thin case we should use the cctx identifier from the last cctx as the inbound hash and outbound chain id as theinbound chain id
		if len(trackingDetailsList) > 0 {
			lastTrackingDetails := trackingDetailsList[len(trackingDetailsList)-1]
			// The last inbound was not finalized, this means that next cctx has not been created, yet we can return.
			// There is no need to log the error here as we are just trying to track the cctx
			if !lastTrackingDetails.IsInboundFinalized() {
				return nil
			}
			// Update the chain id and hash to the last cctx details
			chainID = lastTrackingDetails.OutboundChain.ChainId
			chainHash = lastTrackingDetails.CCTXIdentifier
		}

		// Create a new context based on the chain id and hash
		ctx, err := zetatoolcontext.NewContext(context.Background(), chainID, chainHash, configFile)
		if err != nil {
			return fmt.Errorf("failed to create context: %w", err)
		}

		// fetch the cctx details and status
		cctxTrackingDetails, err = trackCCTX(ctx)
		if err != nil {
			// The error can be caused by two reasons
			// 1. We have reached the end of the cctx error chain, we can return
			// 2. There was an error in tracking the cctx.

			// If debug flag is set, we log everything
			// If the length of the tracking details is 0, it means that we have not tracked any cctx yet, we log the error
			if cmd.Flag(config.FlagDebug).Changed || len(trackingDetailsList) == 0 {
				log.Error().Msgf("failed to track cctx : %v", err)
				if cctxTrackingDetails != nil {
					log.Info().Msg(cctxTrackingDetails.DebugPrint())
				}
			}

			return nil
		}

		log.Info().Msg(cctxTrackingDetails.Print())
		trackingDetailsList = append(trackingDetailsList, *cctxTrackingDetails)
	}
	return nil
}

func trackCCTX(ctx *zetatoolcontext.Context) (*cctx.TrackingDetails, error) {
	var (
		cctxTrackingDetails = cctx.NewTrackingDetails()
		err                 error
	)
	// Get the ballot identifier for the inbound transaction and confirm that cctx status is at least either PendingInboundConfirmation or PendingInboundVoting
	err = cctxTrackingDetails.CheckInbound(ctx)
	if err != nil {
		return cctxTrackingDetails, fmt.Errorf("failed to get ballot identifier: %w", err)
	}
	// Reject unknown status, as it is not valid
	if cctxTrackingDetails.Status == cctx.Unknown || cctxTrackingDetails.CCTXIdentifier == "" {
		return cctxTrackingDetails, fmt.Errorf("unknown status")
	}

	// At this point, we have confirmed the inbound hash is valid, and it was sent to valid address.
	// After this we attach error messages to the message field as we already have some details about the cctx which can be printed
	// Update cctx status from zetacore.This copies the status from zetacore to the cctx details.The cctx status can only be `PendingInboundVoting` or `PendingInboundConfirmation` at this point
	cctxTrackingDetails.UpdateCCTXStatus(ctx)

	// UpdateCCTXStatus can return without updating the cctx details if the cctx is not found.That is fine which just means that the cctx is not yet created
	// If the inbound is finalized, we can update the outbound details as they would now be available
	if cctxTrackingDetails.IsInboundFinalized() {
		cctxTrackingDetails.UpdateCCTXOutboundDetails(ctx)
	}

	// The cctx details now have status from zetacore, we have not tried to a get more granular status from the outbound chain yet.
	// If it's not pending, we can just return here.
	if !cctxTrackingDetails.IsPendingOutbound() {
		return cctxTrackingDetails, nil
	}

	// Update tx hash list from outbound tracker
	// If the tracker is found, it means the outbound is broadcast, but we are waiting for the confirmations
	// If the tracker is not found, it means the outbound is not broadcast yet, we are wwaiting for the tss to sign the outbound
	cctxTrackingDetails.UpdateHashListAndPendingStatus(ctx)

	// If its not pending confirmation, we can return here, it means the outbound is not broadcast yet its pending tss signing
	if !cctxTrackingDetails.IsPendingConfirmation() {
		return cctxTrackingDetails, nil
	}

	// Check outbound tx, we are waiting for the outbound tx to be confirmed
	err = cctxTrackingDetails.CheckOutbound(ctx)
	if err != nil {
		return cctxTrackingDetails, err
	}
	return cctxTrackingDetails, nil
}
