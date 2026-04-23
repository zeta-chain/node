package observer

import (
	"encoding/hex"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnectornative.sol"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

type v2OutboundParser func(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
) (*big.Int, chains.ReceiveStatus, error)

func TestV2OutboundParsersSkipForgedLogs(t *testing.T) {
	txHash := sample.Hash()
	receiver := sample.EthAddress()
	asset := sample.EthAddress()
	amount := big.NewInt(42)
	payload := []byte{0xde, 0xad, 0xbe, 0xef}

	gatewayAddr := sample.EthAddress()
	forgedGatewayAddr := sample.EthAddress()
	connectorAddr := sample.EthAddress()
	forgedConnectorAddr := sample.EthAddress()
	custodyAddr := sample.EthAddress()
	forgedCustodyAddr := sample.EthAddress()

	gateway := mustNewGatewayEVM(t, gatewayAddr)
	connector := mustNewZetaConnectorNative(t, connectorAddr)
	custody := mustNewERC20Custody(t, custodyAddr)

	tests := []struct {
		name    string
		cctx    *crosschaintypes.CrossChainTx
		receipt *ethtypes.Receipt
		parse   v2OutboundParser
	}{
		{
			name: "parseAndCheckGatewayExecuted",
			cctx: newV2OutboundCCTX(t, receiver, asset, amount, payload),
			receipt: newReceipt(
				txHash,
				makeGatewayExecutedLog(t, forgedGatewayAddr, txHash, receiver, amount, payload),
				makeGatewayExecutedLog(t, gatewayAddr, txHash, receiver, amount, payload),
			),
			parse: func(cctx *crosschaintypes.CrossChainTx, receipt *ethtypes.Receipt) (*big.Int, chains.ReceiveStatus, error) {
				return parseAndCheckGatewayExecuted(cctx, receipt, gatewayAddr, gateway)
			},
		},
		{
			name: "parseAndCheckGatewayReverted",
			cctx: newV2OutboundCCTX(t, receiver, asset, amount, payload),
			receipt: newReceipt(
				txHash,
				makeGatewayRevertedLog(t, forgedGatewayAddr, txHash, receiver, asset, amount, payload),
				makeGatewayRevertedLog(t, gatewayAddr, txHash, receiver, asset, amount, payload),
			),
			parse: func(cctx *crosschaintypes.CrossChainTx, receipt *ethtypes.Receipt) (*big.Int, chains.ReceiveStatus, error) {
				return parseAndCheckGatewayReverted(cctx, receipt, gatewayAddr, gateway)
			},
		},
		{
			name: "parseAndCheckZetaConnectorWithdraw",
			cctx: newV2OutboundCCTX(t, receiver, asset, amount, payload),
			receipt: newReceipt(
				txHash,
				makeZetaConnectorWithdrawLog(t, forgedConnectorAddr, txHash, receiver, amount),
				makeZetaConnectorWithdrawLog(t, connectorAddr, txHash, receiver, amount),
			),
			parse: func(cctx *crosschaintypes.CrossChainTx, receipt *ethtypes.Receipt) (*big.Int, chains.ReceiveStatus, error) {
				return parseAndCheckZetaConnectorWithdraw(cctx, receipt, connectorAddr, connector)
			},
		},
		{
			name: "parseAndCheckERC20CustodyWithdraw",
			cctx: newV2OutboundCCTX(t, receiver, asset, amount, payload),
			receipt: newReceipt(
				txHash,
				makeERC20CustodyWithdrawLog(t, forgedCustodyAddr, txHash, receiver, asset, amount),
				makeERC20CustodyWithdrawLog(t, custodyAddr, txHash, receiver, asset, amount),
			),
			parse: func(cctx *crosschaintypes.CrossChainTx, receipt *ethtypes.Receipt) (*big.Int, chains.ReceiveStatus, error) {
				return parseAndCheckERC20CustodyWithdraw(cctx, receipt, custodyAddr, custody)
			},
		},
		{
			name: "parseAndCheckERC20CustodyWithdrawAndCall",
			cctx: newV2OutboundCCTX(t, receiver, asset, amount, payload),
			receipt: newReceipt(
				txHash,
				makeERC20CustodyWithdrawAndCallLog(t, forgedCustodyAddr, txHash, receiver, asset, amount, payload),
				makeERC20CustodyWithdrawAndCallLog(t, custodyAddr, txHash, receiver, asset, amount, payload),
			),
			parse: func(cctx *crosschaintypes.CrossChainTx, receipt *ethtypes.Receipt) (*big.Int, chains.ReceiveStatus, error) {
				return parseAndCheckERC20CustodyWithdrawAndCall(cctx, receipt, custodyAddr, custody)
			},
		},
		{
			name: "parseAndZetaConnectorWithdrawAndCall",
			cctx: newV2OutboundCCTX(t, receiver, asset, amount, payload),
			receipt: newReceipt(
				txHash,
				makeZetaConnectorWithdrawAndCallLog(t, forgedConnectorAddr, txHash, receiver, amount, payload),
				makeZetaConnectorWithdrawAndCallLog(t, connectorAddr, txHash, receiver, amount, payload),
			),
			parse: func(cctx *crosschaintypes.CrossChainTx, receipt *ethtypes.Receipt) (*big.Int, chains.ReceiveStatus, error) {
				return parseAndZetaConnectorWithdrawAndCall(cctx, receipt, connectorAddr, connector)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAmount, gotStatus, err := tt.parse(tt.cctx, tt.receipt)
			require.NoError(t, err)
			require.Equal(t, chains.ReceiveStatus_success, gotStatus)
			require.Zero(t, gotAmount.Cmp(amount))
		})
	}
}

func newV2OutboundCCTX(
	t *testing.T,
	receiver ethcommon.Address,
	asset ethcommon.Address,
	amount *big.Int,
	payload []byte,
) *crosschaintypes.CrossChainTx {
	t.Helper()

	cctx := sample.CrossChainTxV2(t, "issue-4567")
	cctx.InboundParams = &crosschaintypes.InboundParams{Asset: asset.Hex()}
	cctx.OutboundParams = []*crosschaintypes.OutboundParams{{
		Receiver:    receiver.Hex(),
		Amount:      sdkmath.NewUintFromBigInt(amount),
		CallOptions: &crosschaintypes.CallOptions{},
	}}
	cctx.RelayedMessage = hex.EncodeToString(payload)

	return cctx
}

func newReceipt(txHash ethcommon.Hash, logs ...*ethtypes.Log) *ethtypes.Receipt {
	return &ethtypes.Receipt{
		Status: ethtypes.ReceiptStatusSuccessful,
		TxHash: txHash,
		Logs:   logs,
	}
}

func makeGatewayExecutedLog(
	t *testing.T,
	emitter ethcommon.Address,
	txHash ethcommon.Hash,
	destination ethcommon.Address,
	amount *big.Int,
	data []byte,
) *ethtypes.Log {
	t.Helper()

	event := mustGatewayABI(t).Events["Executed"]
	return newEventLog(t, emitter, txHash, event, []ethcommon.Hash{addressTopic(destination)}, amount, data)
}

func makeGatewayRevertedLog(
	t *testing.T,
	emitter ethcommon.Address,
	txHash ethcommon.Hash,
	to ethcommon.Address,
	token ethcommon.Address,
	amount *big.Int,
	data []byte,
) *ethtypes.Log {
	t.Helper()

	event := mustGatewayABI(t).Events["Reverted"]
	return newEventLog(
		t,
		emitter,
		txHash,
		event,
		[]ethcommon.Hash{addressTopic(to), addressTopic(token)},
		amount,
		data,
		gatewayevm.RevertContext{
			Sender:        sample.EthAddress(),
			Asset:         token,
			Amount:        amount,
			RevertMessage: data,
		},
	)
}

func makeZetaConnectorWithdrawLog(
	t *testing.T,
	emitter ethcommon.Address,
	txHash ethcommon.Hash,
	to ethcommon.Address,
	amount *big.Int,
) *ethtypes.Log {
	t.Helper()

	event := mustZetaConnectorNativeABI(t).Events["Withdrawn"]
	return newEventLog(t, emitter, txHash, event, []ethcommon.Hash{addressTopic(to)}, amount)
}

func makeERC20CustodyWithdrawLog(
	t *testing.T,
	emitter ethcommon.Address,
	txHash ethcommon.Hash,
	to ethcommon.Address,
	token ethcommon.Address,
	amount *big.Int,
) *ethtypes.Log {
	t.Helper()

	event := mustERC20CustodyABI(t).Events["Withdrawn"]
	return newEventLog(t, emitter, txHash, event, []ethcommon.Hash{addressTopic(to), addressTopic(token)}, amount)
}

func makeERC20CustodyWithdrawAndCallLog(
	t *testing.T,
	emitter ethcommon.Address,
	txHash ethcommon.Hash,
	to ethcommon.Address,
	token ethcommon.Address,
	amount *big.Int,
	data []byte,
) *ethtypes.Log {
	t.Helper()

	event := mustERC20CustodyABI(t).Events["WithdrawnAndCalled"]
	return newEventLog(t, emitter, txHash, event, []ethcommon.Hash{addressTopic(to), addressTopic(token)}, amount, data)
}

func makeZetaConnectorWithdrawAndCallLog(
	t *testing.T,
	emitter ethcommon.Address,
	txHash ethcommon.Hash,
	to ethcommon.Address,
	amount *big.Int,
	data []byte,
) *ethtypes.Log {
	t.Helper()

	event := mustZetaConnectorNativeABI(t).Events["WithdrawnAndCalled"]
	return newEventLog(t, emitter, txHash, event, []ethcommon.Hash{addressTopic(to)}, amount, data)
}

func newEventLog(
	t *testing.T,
	emitter ethcommon.Address,
	txHash ethcommon.Hash,
	event abi.Event,
	indexedTopics []ethcommon.Hash,
	nonIndexed ...interface{},
) *ethtypes.Log {
	t.Helper()

	data, err := event.Inputs.NonIndexed().Pack(nonIndexed...)
	require.NoError(t, err)

	topics := []ethcommon.Hash{event.ID}
	topics = append(topics, indexedTopics...)

	return &ethtypes.Log{
		Address: emitter,
		Topics:  topics,
		Data:    data,
		TxHash:  txHash,
	}
}

func addressTopic(address ethcommon.Address) ethcommon.Hash {
	return ethcommon.BytesToHash(ethcommon.LeftPadBytes(address.Bytes(), 32))
}

func mustNewGatewayEVM(t *testing.T, address ethcommon.Address) *gatewayevm.GatewayEVM {
	t.Helper()

	gateway, err := gatewayevm.NewGatewayEVM(address, &ethclient.Client{})
	require.NoError(t, err)

	return gateway
}

func mustNewERC20Custody(t *testing.T, address ethcommon.Address) *erc20custody.ERC20Custody {
	t.Helper()

	custody, err := erc20custody.NewERC20Custody(address, &ethclient.Client{})
	require.NoError(t, err)

	return custody
}

func mustNewZetaConnectorNative(t *testing.T, address ethcommon.Address) *zetaconnectornative.ZetaConnectorNative {
	t.Helper()

	connector, err := zetaconnectornative.NewZetaConnectorNative(address, &ethclient.Client{})
	require.NoError(t, err)

	return connector
}

func mustGatewayABI(t *testing.T) *abi.ABI {
	t.Helper()

	parsed, err := gatewayevm.GatewayEVMMetaData.GetAbi()
	require.NoError(t, err)
	require.NotNil(t, parsed)

	return parsed
}

func mustERC20CustodyABI(t *testing.T) *abi.ABI {
	t.Helper()

	parsed, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	require.NoError(t, err)
	require.NotNil(t, parsed)

	return parsed
}

func mustZetaConnectorNativeABI(t *testing.T) *abi.ABI {
	t.Helper()

	parsed, err := zetaconnectornative.ZetaConnectorNativeMetaData.GetAbi()
	require.NoError(t, err)
	require.NotNil(t, parsed)

	return parsed
}
