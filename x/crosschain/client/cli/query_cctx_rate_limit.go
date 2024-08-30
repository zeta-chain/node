package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdListPendingCCTXWithinRateLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list_pending_cctx_within_rate_limit",
		Short: "list all pending CCTX within rate limit",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ListPendingCctxWithinRateLimit(
				context.Background(), &types.QueryListPendingCctxWithinRateLimitRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
