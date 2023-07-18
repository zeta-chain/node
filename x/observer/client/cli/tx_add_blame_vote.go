package cli

import (
	"encoding/hex"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"strconv"
)

var _ = strconv.Itoa(0)

func CmdAddBlameVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-blame-vote [chain-id] [index] [failure-reason] [node-list]",
		Short: "Broadcast message add-blame-vote",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			chainID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			index := args[1]
			failureReason := args[2]
			nodeList := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			dst := make([]byte, hex.DecodedLen(len(nodeList)))
			_, err = hex.Decode(dst, []byte(nodeList))
			if err != nil {
				return err
			}
			var nodes []*types.Node
			err = json.Unmarshal(dst, &nodes)
			if err != nil {
				return err
			}
			blameInfo := &types.Blame{
				Index:         index,
				FailureReason: failureReason,
				Nodes:         nodes,
			}

			msg := types.NewMsgAddBlameVoteMsg(clientCtx.From, int64(chainID), blameInfo)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
