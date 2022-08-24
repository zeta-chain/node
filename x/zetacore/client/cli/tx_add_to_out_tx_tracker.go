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

func CmdAddToWatchList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-to-watch-list [chain] [nonce] [tx-hash]",
		Short: "Broadcast message AddToWatchList",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain := args[0]
			argNonce, _ := strconv.ParseInt(args[1], 10, 64)
			argTxHash := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddToOutTxTracker(
				clientCtx.GetFromAddress().String(),
				argChain,
				uint64(argNonce),
				argTxHash,
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
