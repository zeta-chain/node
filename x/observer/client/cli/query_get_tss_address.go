package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdGetTssAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-tss-address",
		Short: "Query current tss address",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetTssAddressRequest{}

			res, err := queryClient.GetTssAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetTssAddressByFinalizedZetaHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-historical-tss-address [finalizedZetaHeight]",
		Short: "Query tss address by finalized zeta height (for historical tss addresses)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			finalizedZetaHeight, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			params := &types.QueryGetTssAddressByFinalizedHeightRequest{
				FinalizedZetaHeight: finalizedZetaHeight,
			}

			res, err := queryClient.GetTssAddressByFinalizedHeight(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
