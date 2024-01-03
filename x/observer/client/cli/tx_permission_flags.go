package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdUpdateCrosschainFlags() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-crosschain-flags [is-inbound-enabled] [is-outbound-enabled]",
		Short: "Update crosschain flags",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			argIsInboundEnabled, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			arsIsOutboundEnabled, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateCrosschainFlags(clientCtx.GetFromAddress().String(), argIsInboundEnabled, arsIsOutboundEnabled)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
