package cli

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdVoteInbound() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vote-inbound [sender] [senderChainID] [txOrigin] [receiver] [receiverChainID] [amount] [message" +
			"] [inboundHash] [inBlockHeight] [coinType] [asset] [eventIndex]",
		Short: "Broadcast message to vote an inbound",
		Args:  cobra.ExactArgs(12),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSender := args[0]
			argsSenderChain, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			argsTxOrigin := args[2]
			argsReceiver := args[3]
			argsReceiverChain, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			amount := math.NewUintFromString(args[5])

			argsMessage := args[6]
			argsInboundHash := args[7]

			argsInBlockHeight, err := strconv.ParseUint(args[8], 10, 64)
			if err != nil {
				return err
			}

			coinType, ok := coin.CoinType_value[args[9]]
			if !ok {
				return fmt.Errorf("wrong coin type %s", args[9])
			}
			argsCoinType := coin.CoinType(coinType)

			argsAsset := args[10]

			// parse argsp[11] to uint type and not uint64
			argsEventIndex, err := strconv.ParseUint(args[11], 10, 32)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgVoteInbound(
				clientCtx.GetFromAddress().String(),
				argsSender,
				argsSenderChain,
				argsTxOrigin,
				argsReceiver,
				argsReceiverChain,
				amount,
				argsMessage,
				argsInboundHash,
				argsInBlockHeight,
				250_000,
				argsCoinType,
				argsAsset,
				uint(argsEventIndex),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
