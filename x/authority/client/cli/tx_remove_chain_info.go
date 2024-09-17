package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/node/x/authority/types"
)

func CmdRemoveChainInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-chain-info [chain-id]",
		Short: "Remove the chain info for the specified chain id",
		Long: "Remove the chain info for the specified chain id. The chain info will be removed from the chain info store" +
			`Example:
$ zetacored tx authority remove-chain-info 42`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveChainInfo(
				clientCtx.GetFromAddress().String(),
				chainID)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
