package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func CmdShowHeaderEnabledChains() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-header-enabled-chains",
		Short: "Show the verification flags",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryHeaderEnabledChainsRequest{}

			res, err := queryClient.HeaderEnabledChains(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
