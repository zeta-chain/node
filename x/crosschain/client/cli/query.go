package cli

import (
	"fmt"

	// "strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group crosschain queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdListOutTxTracker(),
		CmdShowOutTxTracker(),
		CmdListGasPrice(),
		CmdShowGasPrice(),
		CmdListChainNonces(),
		CmdShowChainNonces(),
		CmdListSend(),
		CmdShowSend(),
		CmdLastZetaHeight(),
		CmdInTxHashToCctxData(),
		CmdListInTxHashToCctx(),
		CmdShowInTxHashToCctx(),
		CmdQueryParams(),
		CmdListPendingNonces(),
		CmdPendingCctx(),
		CmdListInTxTrackerByChain(),
		CmdListInTxTrackers(),
		CmdGetZetaAccounting(),
	)

	return cmd
}
