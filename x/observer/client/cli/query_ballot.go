package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/observer/types"
)

// CmdBallotByIdentifier returns a command which queries a ballot by its identifier
func CmdBallotByIdentifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-ballot [ballot-identifier]",
		Short: "Query BallotByIdentifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqVoteIdentifier := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBallotByIdentifierRequest{
				BallotIdentifier: reqVoteIdentifier,
			}

			res, err := queryClient.BallotByIdentifier(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdAllBallots returns a command which queries all ballots
func CmdAllBallots() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-ballots",
		Short: "Query all ballots",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBallotsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Ballots(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	return cmd
}

func CmdBallotListForHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-ballots-for-height [height]",
		Short: "Query BallotListForHeight",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryBallotListForHeightRequest{
				Height: height,
			}

			res, err := queryClient.BallotListForHeight(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
