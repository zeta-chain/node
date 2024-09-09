package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"gitlab.com/thorchain/tss/go-tss/blame"

	"github.com/zeta-chain/node/x/observer/types"
)

func CmdVoteBlame() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-blame [chain-id] [index] [failure-reason] [node-list]",
		Short: "Broadcast message vote-blame",
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

			var nodes []blame.Node
			err = json.Unmarshal(dst, &nodes)
			if err != nil {
				return err
			}
			blameNodes := ConvertNodes(nodes)
			blameInfo := types.Blame{
				Index:         index,
				FailureReason: failureReason,
				Nodes:         blameNodes,
			}

			msg := types.NewMsgVoteBlameMsg(clientCtx.GetFromAddress().String(), int64(chainID), blameInfo)

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
		RunE: func(_ *cobra.Command, args []string) (err error) {
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
