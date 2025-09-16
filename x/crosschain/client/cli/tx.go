package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/crosschain/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdAddOutboundTracker(),
		CmdVoteGasPrice(),
		CmdVoteInbound(),
		CmdVoteOutbound(),
		CmdRemoveOutboundTracker(),
		CmdUpdateTss(),
		CmdMigrateTssFunds(),
		CmdAddInboundTracker(),
		CmdWhitelistAsset(),
		CmdAbortStuckCCTX(),
		CmdRefundAborted(),
		CmdRemoveInboundTracker(),
	)

	return cmd
}
