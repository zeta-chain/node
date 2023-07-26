package cli

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/zeta-chain/zetacore/x/observer/client/cli"
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

	cmd.AddCommand(CmdAddToWatchList())
	cmd.AddCommand(CmdCreateTSSVoter())
	cmd.AddCommand(CmdGasPriceVoter())
	cmd.AddCommand(CmdNonceVoter())
	cmd.AddCommand(CmdCCTXOutboundVoter())
	cmd.AddCommand(CmdCCTXInboundVoter())
	cmd.AddCommand(CmdRemoveFromWatchList())
	cmd.AddCommand(cli.CmdUpdatePermissionFlags())
	cmd.AddCommand(cli.CmdUpdateKeygen())
	// this line is used by starport scaffolding # 1

	return cmd
}
