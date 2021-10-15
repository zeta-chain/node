package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group metacore queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	cmd.AddCommand(CmdListSend())
	cmd.AddCommand(CmdShowSend())

	cmd.AddCommand(CmdListSendVoter())
	cmd.AddCommand(CmdShowSendVoter())

	cmd.AddCommand(CmdListTxoutConfirmation())
	cmd.AddCommand(CmdShowTxoutConfirmation())

	cmd.AddCommand(CmdListTxout())
	cmd.AddCommand(CmdShowTxout())

	cmd.AddCommand(CmdListNodeAccount())
	cmd.AddCommand(CmdShowNodeAccount())

	cmd.AddCommand(CmdLastMetaHeight())

	cmd.AddCommand(CmdListTxinVoter())
	cmd.AddCommand(CmdShowTxinVoter())

	cmd.AddCommand(CmdListTxin())
	cmd.AddCommand(CmdShowTxin())

	return cmd
}
