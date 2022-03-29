package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

var _ = strconv.Itoa(0)

func CmdTxListRich() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx-list-rich [last]",
		Short: "Query txListRich",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reqLast := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			last, _ := strconv.ParseUint(reqLast, 10, 64)
			params := &types.QueryTxListRichRequest{
				Last: last,
			}

			res, err := queryClient.TxListRich(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
