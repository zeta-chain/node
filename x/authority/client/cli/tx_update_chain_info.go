package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/x/authority/types"
)

func CmdUpdateChainInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-chain-info [chain-info-json-file]",
		Short: "Update the chain info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Read the chain info from the file using os package and unmarshal it into the chain info variable
			var chainInfo types.ChainInfo
			chainInfoBytes, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			if err := chainInfo.Unmarshal(chainInfoBytes); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateChainInfo(
				clientCtx.GetFromAddress().String(),
				chainInfo,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func 