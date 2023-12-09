package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func CmdUpdateContractBytecode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-contract-bytecode [contractAddress] [newBytecodeAddress]",
		Short: "Broadcast message UpdateContractBytecode",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddress := args[0]

			newBytecodeAddress := args[1]

			msg := types.NewMsgUpdateContractBytecode(
				clientCtx.GetFromAddress().String(),
				contractAddress,
				newBytecodeAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
