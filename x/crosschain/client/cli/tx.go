package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
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
		CmdAddToWatchList(),
		CmdCreateTSSVoter(),
		CmdGasPriceVoter(),
		CmdNonceVoter(),
		CmdCCTXOutboundVoter(),
		CmdCCTXInboundVoter(),
		CmdRemoveFromWatchList(),
		CmdUpdateTss(),
		CmdMigrateTssFunds(),
		CmdAddToInTxTracker(),
	)

	return cmd
}
