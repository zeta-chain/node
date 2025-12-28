package chains

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	zetacontext "github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/crypto"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

func resolveRPC(chain chains.Chain, cfg *config.Config) string {
	return map[chains.Network]string{
		chains.Network_eth:     cfg.EthereumRPC,
		chains.Network_base:    cfg.BaseRPC,
		chains.Network_polygon: cfg.PolygonRPC,
		chains.Network_bsc:     cfg.BscRPC,
	}[chain.Network]
}

func GetEvmClient(ctx *zetacontext.Context, chain chains.Chain) (*ethclient.Client, error) {
	evmRRC := resolveRPC(chain, ctx.GetConfig())
	if evmRRC == "" {
		return nil, fmt.Errorf("rpc not found for chain %d network %s", chain.ChainId, chain.Network)
	}
	rpcClient, err := ethrpc.DialHTTP(evmRRC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to eth rpc: %w", err)
	}
	return ethclient.NewClient(rpcClient), nil
}

func GetEvmTx(
	ctx *zetacontext.Context,
	evmClient *ethclient.Client,
	inboundHash string,
	chain chains.Chain,
) (*ethtypes.Transaction, *ethtypes.Receipt, error) {
	goCtx := ctx.GetContext()
	// Fetch transaction from the inbound
	hash := ethcommon.HexToHash(inboundHash)
	tx, isPending, err := evmClient.TransactionByHash(goCtx, hash)
	if err != nil {
		return nil, nil, fmt.Errorf("tx not found on chain: %w,chainID: %d", err, chain.ChainId)
	}
	if isPending {
		return nil, nil, fmt.Errorf("tx is still pending on chain: %d", chain.ChainId)
	}
	receipt, err := evmClient.TransactionReceipt(goCtx, hash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt: %w, tx hash: %s", err, inboundHash)
	}
	return tx, receipt, nil
}

func ZetaTokenVoteV1(
	event *zetaconnector.ZetaConnectorNonEthZetaSent,
	observationChain int64,
) *crosschaintypes.MsgVoteInbound {
	// note that this is most likely zeta chain
	destChain, found := chains.GetChainFromChainID(event.DestinationChainId.Int64(), []chains.Chain{})
	if !found {
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
		uint64(event.Raw.Index),
		crosschaintypes.InboundStatus_SUCCESS,
	)
}

func Erc20VoteV1(
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
		uint64(event.Raw.Index),
		crosschaintypes.InboundStatus_SUCCESS,
	)
}

func GasVoteV1(
	tx *ethtypes.Transaction,
	sender ethcommon.Address,
	blockNumber uint64,
	senderChainID int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	message := string(tx.Data())
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
		crosschaintypes.InboundStatus_SUCCESS,
	)
}

func DepositInboundVoteV2(event *gatewayevm.GatewayEVMDeposited,
	senderChainID int64,
	zetacoreChainID int64) *crosschaintypes.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	// to maintain compatibility with previous gateway version, deposit event with a non-empty payload is considered as a call
	isCrossChainCall := len(event.Payload) > 0

	return crosschaintypes.NewMsgVoteInbound(
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
		uint64(event.Raw.Index),
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
		crosschaintypes.WithCrossChainCall(isCrossChainCall),
	)
}

func DepositAndCallInboundVoteV2(event *gatewayevm.GatewayEVMDepositedAndCalled,
	senderChainID int64,
	zetacoreChainID int64) *crosschaintypes.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	return crosschaintypes.NewMsgVoteInbound(
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
		uint64(event.Raw.Index),
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
		crosschaintypes.WithCrossChainCall(true),
	)
}

func CallInboundVoteV2(event *gatewayevm.GatewayEVMCalled,
	senderChainID int64,
	zetacoreChainID int64) *crosschaintypes.MsgVoteInbound {
	return crosschaintypes.NewMsgVoteInbound(
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
		uint64(event.Raw.Index),
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
	)
}

// GetEVMBalance fetches the native token balance for an address on an EVM chain
func GetEVMBalance(ctx context.Context, rpcURL string, address ethcommon.Address) (*big.Int, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.BalanceAt(ctx, address, nil)
}

// FormatEVMBalance converts wei to ETH with 9 decimal places
func FormatEVMBalance(wei *big.Int) string {
	if wei == nil {
		return "0.000000000"
	}

	weiFloat := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetInt(big.NewInt(params.Ether))
	eth := new(big.Float).Quo(weiFloat, divisor)

	return eth.Text('f', 9)
}
