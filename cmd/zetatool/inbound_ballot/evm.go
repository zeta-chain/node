package inbound_ballot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/observer/types"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
)

func EvmInboundBallotIdentified(ctx context.Context,
	cfg config.Config,
	zetacoreClient rpc.Clients,
	inboundHash string,
	inboundChain chains.Chain,
	zetaChainID int64) (string, error) {
	rpcClient, err := ethrpc.DialHTTP(cfg.EthRPCURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to eth rpc %s", err.Error())
	}
	evmClient := ethclient.NewClient(rpcClient)

	res, err := zetacoreClient.Observer.GetTssAddress(context.Background(), &types.QueryGetTssAddressRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get tss address %s", err.Error())
	}
	tssEthAddress := res.GetEth()

	fmt.Println("ETH Address: ", tssEthAddress)

	hash := ethcommon.HexToHash(inboundHash)
	tx, isPending, err := evmClient.TransactionByHash(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("tx not found on chain %s, %d", err.Error(), inboundChain.ChainId)
	}
	if isPending {
		return "", fmt.Errorf("tx is still pending on chain %d", inboundChain.ChainId)
	}
	receipt, err := evmClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("failed to get receipt %s , tx hash %s", err.Error(), inboundHash)
	}
	chainParams, err := zetacoreClient.GetChainParamsForChainID(context.Background(), inboundChain.ChainId)
	if err != nil {
		return "", fmt.Errorf("failed to get chain params %s", err.Error())
	}

	switch tx.To().Hex() {
	case chainParams.ConnectorContractAddress:
		{
			// build inbound vote message and post vote
			addrConnector := ethcommon.HexToAddress(chainParams.ConnectorContractAddress)
			connector, err := zetaconnector.NewZetaConnectorNonEth(addrConnector, evmClient)
			if err != nil {
				return "", fmt.Errorf("failed to get connector contract %s", err.Error())
			}
			for _, log := range receipt.Logs {
				event, err := connector.ParseZetaSent(*log)
				if err == nil && event != nil {
					msg := zetaTokenVoteV1(event, inboundChain.ChainId)
					return msg.Digest(), nil
				}
			}
		}
	case chainParams.Erc20CustodyContractAddress:
		{
			addrCustody := ethcommon.HexToAddress(chainParams.Erc20CustodyContractAddress)
			custody, err := erc20custody.NewERC20Custody(addrCustody, evmClient)
			if err != nil {
				return "", fmt.Errorf("failed to get custody contract %s", err.Error())
			}
			sender, err := evmClient.TransactionSender(ctx, tx, receipt.BlockHash, receipt.TransactionIndex)
			if err != nil {
				return "", fmt.Errorf("failed to get tx sender %s", err.Error())
			}
			for _, log := range receipt.Logs {
				zetaDeposited, err := custody.ParseDeposited(*log)
				if err == nil && zetaDeposited != nil {
					msg := erc20VoteV1(zetaDeposited, sender, inboundChain.ChainId, zetaChainID)
					return msg.Digest(), nil
				}
			}
		}
	case tssEthAddress:
		{
			if receipt.Status != ethtypes.ReceiptStatusSuccessful {
				return "", fmt.Errorf("tx failed on chain %d", inboundChain.ChainId)
			}
			sender, err := evmClient.TransactionSender(ctx, tx, receipt.BlockHash, receipt.TransactionIndex)
			if err != nil {
				return "", fmt.Errorf("failed to get tx sender %s", err.Error())
			}
			msg := gasVoteV1(tx, sender, receipt.BlockNumber.Uint64(), inboundChain.ChainId, zetaChainID)
			return msg.Digest(), nil
		}
	case chainParams.GatewayAddress:
		{
			gatewayAddr := ethcommon.HexToAddress(chainParams.GatewayAddress)
			gateway, err := gatewayevm.NewGatewayEVM(gatewayAddr, evmClient)
			if err != nil {
				return "", fmt.Errorf("failed to get gateway contract %s", err.Error())
			}
			for _, log := range receipt.Logs {
				if log == nil || log.Address != gatewayAddr {
					continue
				}
				eventDeposit, err := gateway.ParseDeposited(*log)
				if err == nil {
					msg := depositInboundVoteV2(eventDeposit, inboundChain.ChainId, zetaChainID)
					return msg.Digest(), nil
				}
				eventDepositAndCall, err := gateway.ParseDepositedAndCalled(*log)
				if err == nil {
					msg := depositAndCallInboundVoteV2(eventDepositAndCall, inboundChain.ChainId, zetaChainID)
					return msg.Digest(), nil
				}
				eventCall, err := gateway.ParseCalled(*log)
				if err == nil {
					msg := callInboundVoteV2(eventCall, inboundChain.ChainId, zetaChainID)
					return msg.Digest(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no event found for tx %s", inboundHash)
}

func zetaTokenVoteV1(
	event *zetaconnector.ZetaConnectorNonEthZetaSent,
	observationChain int64,
) *crosschaintypes.MsgVoteInbound {
	// note that this is most likely zeta chain
	destChain, found := chains.GetChainFromChainID(event.DestinationChainId.Int64(), []chains.Chain{})
	if !found {
		fmt.Println("Not found")
		return nil
	}

	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
	sender := event.ZetaTxSenderAddress.Hex()
	message := base64.StdEncoding.EncodeToString(event.Message)

	return zetacore.GetInboundVoteMessage(
		sender,
		observationChain,
		event.SourceTxOriginAddress.Hex(),
		destAddr,
		destChain.ChainId,
		sdkmath.NewUintFromBigInt(event.ZetaValueAndGas),
		message,
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		event.DestinationGasLimit.Uint64(),
		coin.CoinType_Zeta,
		"",
		"",
		event.Raw.Index,
	)
}

func erc20VoteV1(
	event *erc20custody.ERC20CustodyDeposited,
	sender ethcommon.Address,
	observationChain int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	// donation check
	if bytes.Equal(event.Message, []byte(constant.DonationMessage)) {
		return nil
	}

	return zetacore.GetInboundVoteMessage(
		sender.Hex(),
		observationChain,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		zetacoreChainID,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coin.CoinType_ERC20,
		event.Asset.String(),
		"",
		event.Raw.Index,
	)
}

func gasVoteV1(
	tx *ethtypes.Transaction,
	sender ethcommon.Address,
	blockNumber uint64,
	senderChainID int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	message := string(tx.Data())
	// donation check
	// #nosec G703 err is already checked
	data, _ := hex.DecodeString(message)
	if bytes.Equal(data, []byte(constant.DonationMessage)) {
		return nil
	}

	return zetacore.GetInboundVoteMessage(
		sender.Hex(),
		senderChainID,
		sender.Hex(),
		sender.Hex(),
		zetacoreChainID,
		sdkmath.NewUintFromString(tx.Value().String()),
		message,
		tx.Hash().Hex(),
		blockNumber,
		90_000,
		coin.CoinType_Gas,
		"",
		"",
		0, // not a smart contract call
	)
}

func depositInboundVoteV2(event *gatewayevm.GatewayEVMDeposited,
	senderChainID int64,
	zetacoreChainID int64) crosschaintypes.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	// to maintain compatibility with previous gateway version, deposit event with a non-empty payload is considered as a call
	isCrossChainCall := false
	if len(event.Payload) > 0 {
		isCrossChainCall = true
	}

	return *crosschaintypes.NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChainID,
		"",
		event.Receiver.Hex(),
		zetacoreChainID,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		zetacore.PostVoteInboundCallOptionsGasLimit,
		coinType,
		event.Asset.Hex(),
		event.Raw.Index,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
		crosschaintypes.WithCrossChainCall(isCrossChainCall),
	)
}

func depositAndCallInboundVoteV2(event *gatewayevm.GatewayEVMDepositedAndCalled,
	senderChainID int64,
	zetacoreChainID int64) crosschaintypes.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	return *crosschaintypes.NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChainID,
		"",
		event.Receiver.Hex(),
		zetacoreChainID,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coinType,
		event.Asset.Hex(),
		event.Raw.Index,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
		crosschaintypes.WithCrossChainCall(true),
	)
}

func callInboundVoteV2(event *gatewayevm.GatewayEVMCalled,
	senderChainID int64,
	zetacoreChainID int64) crosschaintypes.MsgVoteInbound {
	return *crosschaintypes.NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChainID,
		"",
		event.Receiver.Hex(),
		zetacoreChainID,
		sdkmath.ZeroUint(),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		zetacore.PostVoteInboundCallOptionsGasLimit,
		coin.CoinType_NoAssetCall,
		"",
		event.Raw.Index,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
	)
}
