package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdGetTssFundsMigrator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-tss-funds-migrator [chain-id]",
		Short: "show the tss funds migrator for a chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTssFundsMigratorInfoRequest{}

			res, err := queryClient.TssFundsMigratorInfo(context.Background(), params)
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

func CmdGetAllTssFundsMigrator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tss-funds-migrator",
		Short: "list all tss funds migrators",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTssFundsMigratorInfoAllRequest{}

			res, err := queryClient.TssFundsMigratorInfoAll(context.Background(), params)
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
