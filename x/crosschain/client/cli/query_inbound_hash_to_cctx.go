package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdListInboundHashToCctx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-inbound-hash-to-cctx",
		Short: "list all inboundHashToCctx",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllInboundHashToCctxRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.InboundHashToCctxAll(context.Background(), params)
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

func CmdShowInboundHashToCctx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-inbound-hash-to-cctx [inbound-hash]",
		Short: "shows a inboundHashToCctx",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argInboundHash := args[0]

			params := &types.QueryGetInboundHashToCctxRequest{
				InboundHash: argInboundHash,
			}

			res, err := queryClient.InboundHashToCctx(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdInboundHashToCctxData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inbound-hash-to-cctx-data [inbound-hash]",
		Short: "query a cctx data from a inbound hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argInboundHash := args[0]

			params := &types.QueryInboundHashToCctxDataRequest{
				InboundHash: argInboundHash,
			}

			res, err := queryClient.InboundHashToCctxData(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
