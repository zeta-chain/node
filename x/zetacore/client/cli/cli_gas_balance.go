package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"strconv"
)

func CmdListGasBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-gas-balance",
		Short: "list all GasBalance",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllGasBalanceRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.GasBalanceAll(context.Background(), params)
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

func CmdShowGasBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-gas-balance [index]",
		Short: "shows a GasBalance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetGasBalanceRequest{
				Index: args[0],
			}

			res, err := queryClient.GasBalance(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// Transaction CLI /////////////////////////

func CmdGasBalanceVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-balance-voter [chain] [balance] [blockNumber]",
		Short: "Broadcast message gasBalanceVoter",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsChain := (args[0])
			argsBalance := (args[1])
			argsBlockNumber, _ := strconv.Atoi(args[2])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgGasBalanceVoter(clientCtx.GetFromAddress().String(), (argsChain), (argsBalance), uint64(argsBlockNumber))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
