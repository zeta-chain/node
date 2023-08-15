package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdUpdatePermissionFlags() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-permission-flags [is-inbound-enabled] [is-outbound-enabled]",
		Short: "Update PermissionFlags",
		Args:  cobra.ExactArgs(1),
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
			msg := types.NewMsgUpdatePermissionFlags(clientCtx.GetFromAddress().String(), argIsInboundEnabled, arsIsOutboundEnabled)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
