package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdAddObserver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-observer [observer-address] [zetaclient-grantee-address] [zetaclient-grantee-pubkey]",
		Short: "Broadcast message add-observer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgAddObserver(
				clientCtx.GetFromAddress().String(),
				args[1],
				args[2],
				args[3],
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
