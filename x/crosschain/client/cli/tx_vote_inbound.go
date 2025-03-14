package cli

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdVoteInbound() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vote-inbound [sender] [senderChainID] [txOrigin] [receiver] [receiverChainID] [amount] [message" +
			"] [inboundHash] [inBlockHeight] [coinType] [asset] [eventIndex] [protocolContractVersion] [isArbitraryCall] [confirmationMode] [inboundStatus]",
		Short: "Broadcast message to vote an inbound",
		Example: `zetacored tx crosschain vote-inbound 0xfa233D806C8EB69548F3c4bC0ABb46FaD4e2EB26 8453 "" 0xfa233D806C8EB69548F3c4bC0ABb46FaD4e2EB26 7000 1000000 "" ` +
			`0x66b59ad844404e91faa9587a3061e2f7af36f7a7a1a0afaca3a2efd811bc9463 26170791 Gas 0x0000000000000000000000000000000000000000 587 V2 FALSE SAFE SUCCESS`,

		Args: cobra.ExactArgs(16),
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

			protocolContractVersion, err := parseProtocolContractVersion(args[12])
			if err != nil {
				return err
			}

			isArbitraryCall, err := strconv.ParseBool(args[13])
			if err != nil {
				return err
			}

			confirmationMode, ok := types.ConfirmationMode_value[args[14]]
			if !ok {
				return fmt.Errorf("wrong confirmation mode %s", args[14])
			}
			argsConfirmationMode := types.ConfirmationMode(confirmationMode)

			inboundStatus, ok := types.InboundStatus_value[args[15]]
			if !ok {
				return fmt.Errorf("wrong inbound status %s", args[15])
			}
			argsInboundStatus := types.InboundStatus(inboundStatus)

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
				argsEventIndex,
				protocolContractVersion,
				isArbitraryCall,
				argsInboundStatus,
				argsConfirmationMode,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseProtocolContractVersion(version string) (types.ProtocolContractVersion, error) {
	switch version {
	case "V1":
		return types.ProtocolContractVersion_V1, nil
	case "V2":
		return types.ProtocolContractVersion_V2, nil
	default:
		return types.ProtocolContractVersion_V1, fmt.Errorf(
			"invalid protocol contract version, specify either V1 or V2",
		)
	}
}
