package cli

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/lightclient/types"
)

func CmdDisableVerificationFlags() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-header-verification [list of chain-id]",
		Short: "Disable header verification for the list of chains separated by comma",
		Long: `Provide a list of chain ids separated by comma to disable block header verification for the specified chain ids.

  				Example:
                    To disable verification flags for chain ids 1 and 56
					zetacored tx lightclient disable-header-verification "1,56"
				`,
		Args: cobra.ExactArgs(1),
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

			msg := types.NewMsgDisableHeaderVerification(clientCtx.GetFromAddress().String(), chainIDList)
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
