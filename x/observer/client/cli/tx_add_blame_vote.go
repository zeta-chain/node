package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"os"
	"path/filepath"
	"strconv"
)

var _ = strconv.Itoa(0)

func CmdAddBlameVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-blame-vote [creator] [chain-id] [index] [failure-reason] [node-list]",
		Short: "Broadcast message add-blame-vote",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			println(args)
			creator := args[0]
			chainID, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			index := args[2]
			failureReason := args[3]
			nodeList := args[4]

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

			msg := types.NewMsgAddBlameVoteMsg(creator, int64(chainID), blameInfo)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			println("about to broadcast")
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdEncode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encode [file.json]",
		Short: "Encode a json string into hex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fp := args[0]
			file, err := filepath.Abs(fp)
			if err != nil {
				return err
			}
			file = filepath.Clean(file)
			input, err := os.ReadFile(file) // #nosec G304
			if err != nil {
				return err
			}
			fmt.Println("Hex encoded Node list: ", hex.EncodeToString(input))
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
