package ballots

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

func ZetaTokenVoteV1(
	event *zetaconnector.ZetaConnectorNonEthZetaSent,
	observationChain int64,
) *crosschaintypes.MsgVoteInbound {
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
