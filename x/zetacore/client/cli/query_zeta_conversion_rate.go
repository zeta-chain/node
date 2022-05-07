package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func CmdListZetaConversionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-zeta-conversion-rate",
		Short: "list all zetaConversionRate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllZetaConversionRateRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ZetaConversionRateAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowZetaConversionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-zeta-conversion-rate [index]",
		Short: "shows a zetaConversionRate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]

			params := &types.QueryGetZetaConversionRateRequest{
				Index: argIndex,
			}

			res, err := queryClient.ZetaConversionRate(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
