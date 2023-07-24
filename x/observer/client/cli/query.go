package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group observer queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdBallotByIdentifier())

	cmd.AddCommand(CmdObserversByChainAndType())
	cmd.AddCommand(CmdAllObserverMappers())
	cmd.AddCommand(CmdGetSupportedChains())

	cmd.AddCommand(CmdGetCoreParamsForChain())

	cmd.AddCommand(CmdGetCoreParams())
	cmd.AddCommand(CmdListNodeAccount())
	cmd.AddCommand(CmdShowNodeAccount())
	cmd.AddCommand(CmdShowPermissionFlags())
	cmd.AddCommand(CmdShowKeygen())

	cmd.AddCommand(CmdShowObserverCount())
	cmd.AddCommand(CmdBlameByIdentifier())
	cmd.AddCommand(CmdGetAllBlameRecords())

	// this line is used by starport scaffolding # 1

	return cmd
}
