package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

var _ = strconv.Itoa(0)

func CmdRemoveFromWatchList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-from-out-tx-tracker [chain] [nonce]",
		Short: "Broadcast message RemoveFromWatchlist",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain := args[0]
			argNonce, _ := strconv.ParseInt(args[1], 10, 64)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveFromOutTxTracker(
				clientCtx.GetFromAddress().String(),
				argChain,
				uint64(argNonce),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
