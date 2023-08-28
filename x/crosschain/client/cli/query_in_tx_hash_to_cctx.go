package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdListInTxHashToCctx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-in-tx-hash-to-cctx",
		Short: "list all inTxHashToCctx",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllInTxHashToCctxRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.InTxHashToCctxAll(context.Background(), params)
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

func CmdShowInTxHashToCctx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-in-tx-hash-to-cctx [in-tx-hash]",
		Short: "shows a inTxHashToCctx",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argInTxHash := args[0]

			params := &types.QueryGetInTxHashToCctxRequest{
				InTxHash: argInTxHash,
			}

			res, err := queryClient.InTxHashToCctx(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdInTxHashToCctxData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "in-tx-hash-to-cctx-data [in-tx-hash]",
		Short: "query a cctx data from a in tx hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argInTxHash := args[0]

			params := &types.QueryInTxHashToCctxDataRequest{
				InTxHash: argInTxHash,
			}

			res, err := queryClient.InTxHashToCctxData(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
