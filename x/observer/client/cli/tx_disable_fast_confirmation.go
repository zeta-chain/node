package cli

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/observer/types"
)

func CmdDisableFastConfirmation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable-fast-confirmation [list of chain-id]",
		Short:   "Disable fast confirmation for the list of chains separated by comma; empty list will disable all chains",
		Example: `zetacored tx observer disable-fast-confirmation "1,56"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainIDs := strings.Split(strings.TrimSpace(args[0]), ",")
			var chainIDList []int64
			for _, chainID := range chainIDs {
				chainIDInt, err := strconv.ParseInt(chainID, 10, 64)
				if err != nil {
					return err
				}
				chainIDList = append(chainIDList, chainIDInt)
			}

			msg := types.NewMsgDisableFastConfirmation(clientCtx.GetFromAddress().String(), chainIDList)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
