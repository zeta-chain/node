package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdListInboundTrackerByChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-inbound-tracker [chainId]",
		Short: "shows a list of inbound tracker by chainId",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params := &types.QueryAllInboundTrackerByChainRequest{
				ChainId:    argChain,
				Pagination: pageReq,
			}
			res, err := queryClient.InboundTrackerAllByChain(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	return cmd
}

func CmdListInboundTrackers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-all-inbound-trackers",
		Short: "shows all inboundTrackers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryAllInboundTrackersRequest{}
			res, err := queryClient.InboundTrackerAll(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
