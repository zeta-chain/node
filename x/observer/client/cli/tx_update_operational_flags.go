package cli

import (
	"encoding/json"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/observer/types"
)

const (
	fileFlag                  = "file"
	restartHeightFlag         = "restart-height"
	signerBlockTimeOffsetFlag = "signer-block-time-offset"
)

func CmdUpdateOperationalFlags() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-operational-flags",
		Short: "Broadcast message UpdateOperationalFlags",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var operationalFlags types.OperationalFlags

			flagSet := cmd.Flags()
			file, _ := flagSet.GetString(fileFlag)
			restartHeight, _ := flagSet.GetInt64(restartHeightFlag)
			signerBlockTimeOffset, _ := flagSet.GetDuration(signerBlockTimeOffsetFlag)

			if file != "" {
				input, err := os.ReadFile(file) // #nosec G304
				if err != nil {
					return err
				}
				err = json.Unmarshal(input, &operationalFlags)
				if err != nil {
					return err
				}
			} else {
				operationalFlags.RestartHeight = restartHeight
				operationalFlags.SignerBlockTimeOffset = &signerBlockTimeOffset
			}

			msg := types.NewMsgUpdateOperationalFlags(
				clientCtx.GetFromAddress().String(),
				operationalFlags,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(fileFlag, "", "Path to a JSON file containing OperationalFlags")
	cmd.Flags().Int64(restartHeightFlag, 0, "Height for a coordinated zetaclient restart")
	cmd.Flags().Duration(signerBlockTimeOffsetFlag, 0, "Offset from the zetacore block time to initiate signing")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
