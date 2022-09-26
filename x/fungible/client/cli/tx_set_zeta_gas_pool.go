package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

var _ = strconv.Itoa(0)

func CmdSetZetaGasPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-zeta-gas-pool [address] [pool-type] [chain]",
		Short: "Broadcast message setZetaGasPool",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAddress := args[0]
			argPoolType := args[1]
			argChain := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetZetaGasPool(
				clientCtx.GetFromAddress().String(),
				argAddress,
				argPoolType,
				argChain,
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
