package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/observer/types"
)

func CmdDisableFastConfirmation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable-fast-confirmation [chain-id]",
		Short:   "Disable fast confirmation for the given chain ID",
		Example: `zetacored tx observer disable-fast-confirmation 1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// get chainID as int64
			chainID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgDisableFastConfirmation(clientCtx.GetFromAddress().String(), chainID)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
