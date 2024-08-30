package cli

import (
	"context"
	"strconv"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdShowInboundTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-inbound-tracker [chainID] [txHash]",
		Short: "shows an inbound tracker by chainID and txHash",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "unable to parse chain id from %q", args[0])
			}
			params := &types.QueryInboundTrackerRequest{
				ChainId: argChain,
				TxHash:  args[1],
			}
			res, err := queryClient.InboundTracker(context.Background(), params)
			if err != nil {
				return errors.Wrapf(
					err,
					"failed to fetch inbound tracker for chain %d and tx hash %s",
					argChain,
					args[1],
				)
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
func CmdListInboundTrackerByChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-inbound-tracker [chainId]",
		Short: "shows a list of inbound trackers by chainId",
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
		Short: "shows all inbound trackers",
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
