package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

var _ = strconv.Itoa(0)

func CmdDepoistERC20() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depoist-erc-20 [home-erc-20-contract-address] [recipient-address]",
		Short: "Broadcast message DepoistERC20",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argHomeERC20ContractAddress := args[0]
			argRecipientAddress := args[1]
			argAmount := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDepoistERC20(
				clientCtx.GetFromAddress().String(),
				argHomeERC20ContractAddress,
				argRecipientAddress,
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
