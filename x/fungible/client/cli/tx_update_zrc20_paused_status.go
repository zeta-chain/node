package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/fungible/types"
)

func CmdPauseZRC20() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause-zrc20 [contractAddress1, contractAddress2, ...]",
		Short:   "Broadcast message PauseZRC20",
		Example: `zetacored tx fungible pause-zrc20 "0xece40cbB54d65282c4623f141c4a8a0bE7D6AdEc, 0xece40cbB54d65282c4623f141c4a8a0bEjgksncf" `,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddressList := strings.Split(strings.TrimSpace(args[0]), ",")

			msg := types.NewMsgPauseZRC20(
				clientCtx.GetFromAddress().String(),
				contractAddressList,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUnpauseZRC20() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unpause-zrc20 [contractAddress1, contractAddress2, ...]",
		Short:   "Broadcast message UnpauseZRC20",
		Example: `zetacored tx fungible unpause-zrc20 "0xece40cbB54d65282c4623f141c4a8a0bE7D6AdEc, 0xece40cbB54d65282c4623f141c4a8a0bEjgksncf" `,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddressList := strings.Split(strings.TrimSpace(args[0]), ",")

			msg := types.NewMsgUnpauseZRC20(
				clientCtx.GetFromAddress().String(),
				contractAddressList,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
