package cli

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/cctx"
	"github.com/zeta-chain/node/cmd/zetatool/common"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	zetatoolcontext "github.com/zeta-chain/node/cmd/zetatool/context"
)

func NewTrackCCTXCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "track-cctx [inboundHash] [chain]",
		Short: "track a cross chain transaction",
		Long: `Track a cross chain transaction by its inbound hash and chain.

The chain argument can be either:
  - A chain ID (e.g., 7000, 1, 56)
  - A chain name (e.g., zeta_mainnet, eth_mainnet, bsc_mainnet)

Examples:
  zetatool track-cctx 0x1234... 7000
  zetatool track-cctx 0x1234... zeta_mainnet
  zetatool track-cctx 0x1234... eth_mainnet`,
		RunE: TrackCCTX,
		Args: cobra.ExactArgs(2),
	}
}

func TrackCCTX(cmd *cobra.Command, args []string) error {
	var (
		trackingDetailsList []cctx.TrackingDetails
		cctxTrackingDetails *cctx.TrackingDetails
		maxCCTXChainLength  = 5 // Maximum number of cctx chains to track
	)

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
	for i := 0; i < maxCCTXChainLength; i++ {
		chainID := inboundChainID
		chainHash := inboundHash

		if len(trackingDetailsList) > 0 {
			lastTrackingDetails := trackingDetailsList[len(trackingDetailsList)-1]
			if !lastTrackingDetails.IsInboundFinalized() {
				return nil
			}
			chainID = lastTrackingDetails.OutboundChain.ChainId
			chainHash = lastTrackingDetails.CCTXIdentifier
		}

		ctx, err := zetatoolcontext.NewContext(context.Background(), chainID, chainHash, configFile)
		if err != nil {
			return fmt.Errorf("failed to create context: %w", err)
		}

		cctxTrackingDetails, err = trackCCTX(ctx)
		if err != nil {
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
	cctxTrackingDetails := cctx.NewTrackingDetails()

	err := cctxTrackingDetails.CheckInbound(ctx)
	if err != nil {
		return cctxTrackingDetails, fmt.Errorf("failed to get ballot identifier: %w", err)
	}

	if cctxTrackingDetails.Status == cctx.Unknown || cctxTrackingDetails.CCTXIdentifier == "" {
		return cctxTrackingDetails, fmt.Errorf("unknown status")
	}

	cctxTrackingDetails.UpdateCCTXStatus(ctx)

	if cctxTrackingDetails.IsInboundFinalized() {
		cctxTrackingDetails.UpdateCCTXOutboundDetails(ctx)
	}

	if !cctxTrackingDetails.IsPendingOutbound() {
		return cctxTrackingDetails, nil
	}

	cctxTrackingDetails.UpdateHashListAndPendingStatus(ctx)

	if !cctxTrackingDetails.IsPendingConfirmation() {
		return cctxTrackingDetails, nil
	}

	err = cctxTrackingDetails.CheckOutbound(ctx)
	if err != nil {
		return cctxTrackingDetails, err
	}
	return cctxTrackingDetails, nil
}
