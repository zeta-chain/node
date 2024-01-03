package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdAddObserver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-observer [observer-address] [zetaclient-grantee-pubkey] [add_node_account_only]",
		Short: "Broadcast message add-observer",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			addNodeAccountOnly, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			fmt.Println("addNodeAccountOnly", addNodeAccountOnly)
			msg := types.NewMsgAddObserver(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
				addNodeAccountOnly,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
