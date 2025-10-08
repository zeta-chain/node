package sample

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/testdata"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func RateLimiterFlags() types.RateLimiterFlags {
	r := Rand()

	return types.RateLimiterFlags{
		Enabled: true,
		Window:  r.Int63(),
		Rate:    sdkmath.NewUint(r.Uint64()),
		Conversions: []types.Conversion{
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
		},
	}
}

func RateLimiterFlagsFromRand(r *rand.Rand) types.RateLimiterFlags {
	return types.RateLimiterFlags{
		Enabled: true,
		Window:  r.Int63(),
		Rate:    sdkmath.NewUint(r.Uint64()),
		Conversions: []types.Conversion{
			{
				Zrc20: EthAddressFromRand(r).Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddressFromRand(r).Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddressFromRand(r).Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddressFromRand(r).Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
			{
				Zrc20: EthAddressFromRand(r).Hex(),
				Rate:  sdkmath.LegacyNewDec(r.Int63()),
			},
		},
	}
}

// CustomAssetRate creates a custom asset rate with the given parameters
func CustomAssetRate(
	chainID int64,
	asset string,
	decimals uint32,
	coinType coin.CoinType,
	rate sdkmath.LegacyDec,
) types.AssetRate {
	return types.AssetRate{
		ChainId:  chainID,
		Asset:    strings.ToLower(asset),
		Decimals: decimals,
		CoinType: coinType,
		Rate:     rate,
	}
}

func OutboundTracker(t *testing.T, index string) types.OutboundTracker {
	r := newRandFromStringSeed(t, index)

	return types.OutboundTracker{
		Index:   index,
		ChainId: r.Int63(),
		Nonce:   r.Uint64(),
	}
}

func InboundTracker(t *testing.T, index string) types.InboundTracker {
	r := newRandFromStringSeed(t, index)

	return types.InboundTracker{
		ChainId:  r.Int63(),
		CoinType: coin.CoinType_Zeta,
		TxHash:   Hash().Hex(),
	}
}

func GasPrice(t *testing.T, index string) *types.GasPrice {
	r := newRandFromStringSeed(t, index)

	return &types.GasPrice{
		Creator:     AccAddress(),
		Index:       index,
		ChainId:     r.Int63(),
		Signers:     []string{AccAddress(), AccAddress()},
		BlockNums:   []uint64{r.Uint64(), r.Uint64()},
		Prices:      []uint64{r.Uint64(), r.Uint64()},
		MedianIndex: 0,
	}
}

func GasPriceWithChainID(t *testing.T, chainID int64) types.GasPrice {
	r := newRandFromStringSeed(t, fmt.Sprintf("%d", chainID))

	return types.GasPrice{
		Creator:     AccAddress(),
		ChainId:     chainID,
		Signers:     []string{AccAddress(), AccAddress()},
		BlockNums:   []uint64{r.Uint64(), r.Uint64()},
		Prices:      []uint64{r.Uint64(), r.Uint64()},
		MedianIndex: 0,
	}
}

func GasPriceFromRand(r *rand.Rand, chainID int64) *types.GasPrice {
	var price uint64
	for price == 0 {
		maxGasPrice := uint64(1000 * 1e9) // 1000 Gwei
		price = uint64(1e9) + r.Uint64()%maxGasPrice
	}
	// Select priority fee between 0 and price
	priorityFee := r.Uint64() % price

	// Set priority fee to 0 for Bitcoin chain as it does not have priority fee
	if chains.IsBitcoinChain(chainID, []chains.Chain{}) {
		priorityFee = 0
	}

	return &types.GasPrice{
		Creator:      "",
		ChainId:      chainID,
		Signers:      []string{AccAddressFromRand(r)},
		BlockNums:    []uint64{r.Uint64()},
		Prices:       []uint64{price},
		MedianIndex:  0,
		PriorityFees: []uint64{priorityFee},
	}
}

func InboundParams(r *rand.Rand) *types.InboundParams {
	return &types.InboundParams{
		Sender:                 EthAddress().String(),
		SenderChainId:          r.Int63(),
		TxOrigin:               EthAddress().String(),
		CoinType:               coin.CoinType_Gas,
		Asset:                  StringRandom(r, 32),
		Amount:                 sdkmath.NewUint(uint64(r.Int63())),
		ObservedHash:           StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		BallotIndex:            StringRandom(r, 32),
		FinalizedZetaHeight:    r.Uint64(),
		Status:                 types.InboundStatus_SUCCESS,
		ConfirmationMode:       ConfirmationModeFromRand(r),
	}
}

func InboundParamsValidChainID(r *rand.Rand) *types.InboundParams {
	return &types.InboundParams{
		Sender:                 EthAddress().String(),
		SenderChainId:          chains.Goerli.ChainId,
		TxOrigin:               EthAddress().String(),
		Asset:                  StringRandom(r, 32),
		Amount:                 sdkmath.NewUint(uint64(r.Int63())),
		ObservedHash:           StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		BallotIndex:            StringRandom(r, 32),
		FinalizedZetaHeight:    r.Uint64(),
		Status:                 InboundStatusFromRand(r),
		ConfirmationMode:       ConfirmationModeFromRand(r),
	}
}

func OutboundParams(r *rand.Rand) *types.OutboundParams {
	return &types.OutboundParams{
		Receiver:        EthAddress().String(),
		ReceiverChainId: r.Int63(),
		CoinType:        coin.CoinType(r.Intn(100)),
		Amount:          sdkmath.NewUint(uint64(r.Int63())),
		TssNonce:        uint64(r.Uint32()), // using r.Uint32() can avoid overflow when dealing with `PendingNonces`
		CallOptions: &types.CallOptions{
			GasLimit: r.Uint64(),
		},
		GasPrice:               sdkmath.NewUint(uint64(r.Int63())).String(),
		Hash:                   StringRandom(r, 32),
		BallotIndex:            StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		GasUsed:                r.Uint64(),
		EffectiveGasPrice:      sdkmath.NewInt(r.Int63()),
		ConfirmationMode:       ConfirmationModeFromRand(r),
	}
}

func OutboundParamsValidChainID(r *rand.Rand) *types.OutboundParams {
	return &types.OutboundParams{
		Receiver:        EthAddress().String(),
		ReceiverChainId: chains.Goerli.ChainId,
		Amount:          sdkmath.NewUint(uint64(r.Int63())),
		TssNonce:        r.Uint64(),
		CallOptions: &types.CallOptions{
			GasLimit: r.Uint64(),
		},
		GasPrice:               sdkmath.NewUint(uint64(r.Int63())).String(),
		Hash:                   StringRandom(r, 32),
		BallotIndex:            StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		GasUsed:                r.Uint64(),
		EffectiveGasPrice:      sdkmath.NewInt(r.Int63()),
		ConfirmationMode:       ConfirmationModeFromRand(r),
	}
}

func Status(t *testing.T, index string) *types.Status {
	r := newRandFromStringSeed(t, index)

	createdAt := r.Int63()

	return &types.Status{
		Status:              types.CctxStatus(r.Intn(100)),
		StatusMessage:       String(),
		ErrorMessage:        String(),
		CreatedTimestamp:    createdAt,
		LastUpdateTimestamp: createdAt,
	}
}

func GetCctxIndexFromString(index string) string {
	return crypto.Keccak256Hash([]byte(index)).String()
}

func CrossChainTx(t *testing.T, index string) *types.CrossChainTx {
	r := newRandFromStringSeed(t, index)
	return &types.CrossChainTx{
		Creator:                 AccAddress(),
		Index:                   GetCctxIndexFromString(index),
		ZetaFees:                sdkmath.NewUint(uint64(r.Int63())),
		RelayedMessage:          StringRandom(r, 32),
		CctxStatus:              Status(t, index),
		InboundParams:           InboundParams(r),
		OutboundParams:          []*types.OutboundParams{OutboundParams(r), OutboundParams(r)},
		ProtocolContractVersion: types.ProtocolContractVersion_V1,
		RevertOptions:           types.NewEmptyRevertOptions(),
	}
}

func CrossChainTxV2(t *testing.T, index string) *types.CrossChainTx {
	r := newRandFromStringSeed(t, index)

	return &types.CrossChainTx{
		Creator:                 AccAddress(),
		Index:                   GetCctxIndexFromString(index),
		ZetaFees:                sdkmath.NewUint(uint64(r.Int63())),
		RelayedMessage:          StringRandom(r, 32),
		CctxStatus:              Status(t, index),
		InboundParams:           InboundParams(r),
		OutboundParams:          []*types.OutboundParams{OutboundParams(r), OutboundParams(r)},
		ProtocolContractVersion: types.ProtocolContractVersion_V2,
		RevertOptions:           types.NewEmptyRevertOptions(),
	}
}

// CustomCctxsInBlockRange create 1 cctx per block in block range [lowBlock, highBlock] (inclusive)
func CustomCctxsInBlockRange(
	t *testing.T,
	lowBlock uint64,
	highBlock uint64,
	senderChainID int64,
	receiverChainID int64,
	coinType coin.CoinType,
	asset string,
	amount uint64,
	status types.CctxStatus,
) (cctxs []*types.CrossChainTx) {
	// create 1 cctx per block
	for i := lowBlock; i <= highBlock; i++ {
		nonce := i - 1
		cctx := CrossChainTx(t, fmt.Sprintf("%d-%d", receiverChainID, nonce))
		cctx.CctxStatus.Status = status
		cctx.InboundParams.SenderChainId = senderChainID
		cctx.InboundParams.CoinType = coinType
		cctx.InboundParams.Asset = asset
		cctx.InboundParams.ObservedExternalHeight = i
		cctx.GetCurrentOutboundParam().ReceiverChainId = receiverChainID
		cctx.GetCurrentOutboundParam().Amount = sdkmath.NewUint(amount)
		cctx.GetCurrentOutboundParam().TssNonce = nonce
		cctxs = append(cctxs, cctx)
	}
	return cctxs
}

func LastBlockHeight(t *testing.T, index string) *types.LastBlockHeight {
	r := newRandFromStringSeed(t, index)

	return &types.LastBlockHeight{
		Creator:            AccAddress(),
		Index:              index,
		Chain:              StringRandom(r, 32),
		LastInboundHeight:  r.Uint64(),
		LastOutboundHeight: r.Uint64(),
	}
}

func InboundHashToCctx(t *testing.T, inboundHash string) types.InboundHashToCctx {
	r := newRandFromStringSeed(t, inboundHash)

	return types.InboundHashToCctx{
		InboundHash: inboundHash,
		CctxIndex:   []string{StringRandom(r, 32), StringRandom(r, 32)},
	}
}

func ZetaAccounting(t *testing.T, index string) types.ZetaAccounting {
	r := newRandFromStringSeed(t, index)
	return types.ZetaAccounting{
		AbortedZetaAmount: sdkmath.NewUint(uint64(r.Int63())),
	}
}

// InboundVote creates a sample inbound vote message
func InboundVote(coinType coin.CoinType, from, to int64) types.MsgVoteInbound {
	return types.MsgVoteInbound{
		Creator:            Bech32AccAddress().String(),
		Sender:             EthAddress().String(),
		SenderChainId:      Chain(from).ChainId,
		Receiver:           EthAddress().String(),
		ReceiverChain:      Chain(to).ChainId,
		Amount:             UintInRange(10000000, 1000000000),
		Message:            base64.StdEncoding.EncodeToString(Bytes()),
		InboundBlockHeight: Uint64InRange(1, 10000),
		CallOptions: &types.CallOptions{
			GasLimit: 1000000000,
		},
		InboundHash: Hash().String(),
		CoinType:    coinType,
		TxOrigin:    EthAddress().String(),
		Asset:       "",
		EventIndex:  EventIndex(),
	}
}

func OutboundVote(t *testing.T) types.MsgVoteOutbound {
	cctx := CrossChainTx(t, EthAddress().String())
	return types.MsgVoteOutbound{
		CctxHash:                          cctx.Index,
		OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
		OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
		Status:                            chains.ReceiveStatus_success,
		Creator:                           cctx.Creator,
		ObservedOutboundHash:              common.BytesToHash(EthAddress().Bytes()).String(),
		ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
		ObservedOutboundBlockHeight:       cctx.GetCurrentOutboundParam().ObservedExternalHeight,
		ObservedOutboundEffectiveGasPrice: cctx.GetCurrentOutboundParam().EffectiveGasPrice,
		ObservedOutboundGasUsed:           cctx.GetCurrentOutboundParam().GasUsed,
		CoinType:                          cctx.InboundParams.CoinType,
		ConfirmationMode:                  cctx.GetCurrentOutboundParam().ConfirmationMode,
	}
}

// InboundVoteFromRand creates a simulated inbound vote message. This function uses the provided source of randomness to generate the vote
func InboundVoteFromRand(from, to int64, r *rand.Rand, asset string) types.MsgVoteInbound {
	coinType := CoinTypeFromRand(r)
	_, _, memo := MemoFromRand(r)

	return types.MsgVoteInbound{
		Creator:            "",
		Sender:             EthAddressFromRand(r).String(),
		SenderChainId:      from,
		Receiver:           EthAddressFromRand(r).String(),
		ReceiverChain:      to,
		Amount:             sdkmath.NewUint(r.Uint64()),
		Message:            memo,
		InboundBlockHeight: r.Uint64(),
		CallOptions: &types.CallOptions{
			GasLimit: 1000000000,
		},
		InboundHash:             common.BytesToHash(RandomBytes(r)).String(),
		CoinType:                coinType,
		TxOrigin:                EthAddressFromRand(r).String(),
		Asset:                   asset,
		EventIndex:              r.Uint64(),
		ProtocolContractVersion: ProtocolVersionFromRand(r),
	}
}

func ProtocolVersionFromRand(r *rand.Rand) types.ProtocolContractVersion {
	versions := []types.ProtocolContractVersion{types.ProtocolContractVersion_V1, types.ProtocolContractVersion_V2}
	return versions[r.Intn(len(versions))]
}

func CoinTypeFromRand(r *rand.Rand) coin.CoinType {
	coinTypes := []coin.CoinType{coin.CoinType_Gas, coin.CoinType_ERC20, coin.CoinType_Zeta}
	coinType := coinTypes[r.Intn(len(coinTypes))]
	return coinType
}

func MemoFromRand(r *rand.Rand) (common.Address, []byte, string) {
	randomMemo := common.BytesToAddress([]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78}).
		Hex()
	randomData := []byte(StringRandom(r, 10))
	memoHex := hex.EncodeToString(append(common.FromHex(randomMemo), randomData...))
	return common.HexToAddress(randomMemo), randomData, memoHex
}

func ConfirmationModeFromRand(r *rand.Rand) types.ConfirmationMode {
	types := []types.ConfirmationMode{types.ConfirmationMode_SAFE, types.ConfirmationMode_FAST}
	return types[r.Intn(len(types))]
}

func InboundStatusFromRand(r *rand.Rand) types.InboundStatus {
	statuses := []types.InboundStatus{
		types.InboundStatus_SUCCESS,
		types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE,
		types.InboundStatus_INVALID_RECEIVER_ADDRESS,
	}
	return statuses[r.Intn(len(statuses))]
}

func CCTXfromRand(r *rand.Rand,
	creator string,
	index string,
	to int64,
	from int64,
	tssPubkey string,
	asset string,
) types.CrossChainTx {
	coinType := CoinTypeFromRand(r)

	amount := sdkmath.NewUint(uint64(r.Int63()))
	inbound := &types.InboundParams{
		Sender:                 EthAddressFromRand(r).String(),
		SenderChainId:          from,
		TxOrigin:               EthAddressFromRand(r).String(),
		CoinType:               coinType,
		Asset:                  asset,
		Amount:                 amount,
		ObservedHash:           StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		BallotIndex:            StringRandom(r, 32),
		FinalizedZetaHeight:    r.Uint64(),
		Status:                 InboundStatusFromRand(r),
		ConfirmationMode:       ConfirmationModeFromRand(r),
	}

	outbound := &types.OutboundParams{
		Receiver:        EthAddressFromRand(r).String(),
		ReceiverChainId: to,
		CoinType:        coinType,
		Amount:          sdkmath.NewUint(uint64(r.Int63())),
		TssNonce:        0,
		TssPubkey:       tssPubkey,
		CallOptions: &types.CallOptions{
			GasLimit: r.Uint64(),
		},
		GasPrice:               sdkmath.NewUint(uint64(r.Int63())).String(),
		Hash:                   StringRandom(r, 32),
		BallotIndex:            StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		GasUsed:                100,
		EffectiveGasPrice:      sdkmath.NewInt(r.Int63()),
		EffectiveGasLimit:      100,
		ConfirmationMode:       ConfirmationModeFromRand(r),
	}

	cctx := types.CrossChainTx{
		Creator:        creator,
		Index:          index,
		ZetaFees:       sdkmath.NewUint(1),
		RelayedMessage: base64.StdEncoding.EncodeToString(RandomBytes(r)),
		CctxStatus: &types.Status{
			IsAbortRefunded: false,
			Status:          types.CctxStatus_PendingOutbound,
		},
		InboundParams:           inbound,
		OutboundParams:          []*types.OutboundParams{outbound},
		ProtocolContractVersion: ProtocolVersionFromRand(r),
	}
	return cctx
}

func ZRC20Withdrawal(to []byte, value *big.Int) *zrc20.ZRC20Withdrawal {
	return &zrc20.ZRC20Withdrawal{
		From:            EthAddress(),
		To:              to,
		Value:           value,
		GasFee:          big.NewInt(Int64InRange(100000, 10000000)),
		ProtocolFlatFee: big.NewInt(Int64InRange(100000, 10000000)),
	}
}

func readZetaReceipt(t *testing.T, name string) ethtypes.Receipt {
	receipt, err := testdata.ReadZetaReceipt(name)
	require.NoError(t, err)
	return receipt
}

// InvalidZRC20WithdrawToExternalReceipt is a receipt for a invalid ZRC20 withdrawal to an external address
// receiver is 1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3
func InvalidZRC20WithdrawToExternalReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "invalid_zrc20_withdraw_to_external")
}

// ValidZrc20WithdrawToETHReceipt is a receipt for a ZRC20 withdrawal to an eth address
func ValidZrc20WithdrawToETHReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "zrc20_withdraw_to_eth")
}

// ValidZRC20WithdrawToBTCReceipt is a receipt for a ZRC20 withdrawal to a BTC address
// receiver is bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu
func ValidZRC20WithdrawToBTCReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "zrc20_withdraw_to_btc")
}

// ValidZetaSentDestinationExternalReceipt is a receipt for a Zeta sent to an external destination
func ValidZetaSentDestinationExternalReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "zeta_sent_destination_external")
}

// ValidGatewayWithdrawToSOLChainReceipt is a receipt for a gateway withdraw to a SOL address
// receiver is 9fA4vYZfCa9k9UHjnvYCk4YoipsooapGciKMgaTBw9UH
func ValidGatewayWithdrawToSOLChainReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "gateway_withdraw_to_sol")
}

// InvalidGatewayWithdrawToSOLChainReceipt is a receipt for a invalid gateway withdraw to a SOL address
// receiver is in wrong non solana format
func InvalidGatewayWithdrawToSOLChainReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "invalid_gateway_withdraw_to_sol")
}

// ValidGatewayWithdrawAndCallToSOLChainReceipt is a receipt for a gateway withdraw and call to a SOL program
// receiver is 9fA4vYZfCa9k9UHjnvYCk4YoipsooapGciKMgaTBw9UH
func ValidGatewayWithdrawAndCallToSOLChainReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "gateway_withdraw_and_call_to_sol")
}

// ValidGatewayCallToSOLChainReceipt is a receipt for a gateway call SOL program
// receiver is 9fA4vYZfCa9k9UHjnvYCk4YoipsooapGciKMgaTBw9UH
func ValidGatewayCallToSOLChainReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "gateway_call_to_sol")
}

// InvalidGatewayCallToSOLChainReceipt is a receipt for a gateway call SOL program
// receiver is in wrong non solana format
func InvalidGatewayCallToSOLChainReceipt(t *testing.T) (receipt ethtypes.Receipt) {
	return readZetaReceipt(t, "invalid_gateway_call_to_sol")
}
