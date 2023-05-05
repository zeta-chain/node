package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func CmdUpdateZRC20ProtocolFlatFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-zrc20-protocol-flat-fee [zrc-20 address ] [flat fee]",
		Short: "Broadcast message UpdateZRC20ProtocolFlatFee",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argZRC20 := args[0]
			argAmount := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateZRC20ProtocolFlatFee(
				clientCtx.GetFromAddress().String(),
				argZRC20,
				argAmount,
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
