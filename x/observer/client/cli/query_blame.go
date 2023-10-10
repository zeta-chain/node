package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"strconv"
)

func CmdBlameByIdentifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-blame [blame-identifier]",
		Short: "Query BlameByIdentifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			Identifier := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBlameByIdentifierRequest{
				BlameIdentifier: Identifier,
			}

			res, err := queryClient.BlameByIdentifier(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetAllBlameRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-blame",
		Short: "Query AllBlameRecords",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllBlameRecordsRequest{}

			res, err := queryClient.GetAllBlameRecords(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetBlameByChainAndNonce() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-blame-by-msg [chainId] [nonce]",
		Short: "Query AllBlameRecords",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			chainId := args[0]
			nonce := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			chain, err := strconv.Atoi(chainId)
			if err != nil {
				return err
			}
			nonceInt, err := strconv.Atoi(nonce)
			if err != nil {
				return err
			}
			params := &types.QueryBlameByChainAndNonceRequest{
				ChainId: uint32(chain),
				Nonce:   uint64(nonceInt),
			}

			res, err := queryClient.BlameByChainAndNonce(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
