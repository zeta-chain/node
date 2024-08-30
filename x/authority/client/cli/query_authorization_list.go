package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/authority/types"
)

// CmdAuthorizationsList shows the list of authorizations
func CmdAuthorizationsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-authorizations",
		Short: "lists all authorizations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AuthorizationList(context.Background(), &types.QueryAuthorizationListRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdAuthorization shows the authorization for a given message URL
func CmdAuthorization() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-authorization [msg-url]",
		Short: "shows the authorization for a given message URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			msgURL := args[0]
			res, err := queryClient.Authorization(context.Background(), &types.QueryAuthorizationRequest{
				MsgUrl: msgURL,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
