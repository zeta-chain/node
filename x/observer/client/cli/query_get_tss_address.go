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
		Use:   "get-tss-address [bitcoinChainId]]",
		Short: "Query current tss address",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryGetTssAddressRequest{}
			if len(args) == 1 {
				bitcoinChainId, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
				params.BitcoinChainId = bitcoinChainId
			}

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
		Use:   "get-historical-tss-address [finalizedZetaHeight] [bitcoinChainId]",
		Short: "Query tss address by finalized zeta height (for historical tss addresses)",
		Args:  cobra.ExactArgs(2),
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
			if len(args) == 2 {
				bitcoinChainId, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return err
				}
				params.BitcoinChainId = bitcoinChainId
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
