package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

var _ = strconv.Itoa(0)

func CmdAllObserverMappers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-observer",
		Short: "Query All Observer Mappers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryAllObserverMappersRequest{}
			res, err := queryClient.AllObserverMappers(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdObserversByChainAndType() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-observer [observation-chain] [observation-type]",
		Short: "Query ObserversByChainAndType , Use common.chain for querying",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqObservationChain := args[0]
			reqObservationType := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryObserversByChainAndTypeRequest{
				ObservationChain: reqObservationChain,
				ObservationType:  reqObservationType,
			}

			res, err := queryClient.ObserversByChainAndType(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
