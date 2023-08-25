package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func CmdGasStabilityPoolAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-stability-pool-address [chain-id]",
		Short: "query the address of a gas stability pool for a chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			chainID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetGasStabilityPoolAddress{
				ChainId: chainID,
			}

			res, err := queryClient.GasStabilityPoolAddress(context.Background(), params)
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

func CmdGasStabilityPoolBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-stability-pool-balance [chain-id]",
		Short: "query the balance of a gas stability pool for a chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			chainID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetGasStabilityPoolBalance{
				ChainId: chainID,
			}

			res, err := queryClient.GasStabilityPoolBalance(context.Background(), params)
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
