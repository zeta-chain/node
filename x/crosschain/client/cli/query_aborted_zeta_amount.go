package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdGetAbortedZetaAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-aborted-zeta-amount",
		Short: "Query current tss address",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryAbortedZetaAmountRequest{}
			res, err := queryClient.AbortedZetaAmount(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
