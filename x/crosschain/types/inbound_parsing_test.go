package types_test

import (
	"testing"
)

//func newEvent(t *testing.T) {
//	parsedABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
//	require.NoError(t, err)
//
//	// Create the parameters for the event
//	sender := common.HexToAddress("0x000000000000000000000000000000000000000a")
//	chainID := big.NewInt(1)
//	receiver := []byte("receiverAddress")
//	zrc20 := common.HexToAddress("0x0000000000000000000000000000000000000020")
//	value := big.NewInt(400)
//	gasFee := big.NewInt(100)
//	protocolFlatFee := big.NewInt(10)
//	message := []byte("TestMessage")
//	callOptions := gatewayzevm.CallOptions{}
//	revertOptions := gatewayzevm.RevertOptions{}
//
//	// Pack the event data (non-indexed fields)
//	data, err := parsedABI.Pack("Withdraw", sender, chainID, receiver, zrc20, value, gasFee, protocolFlatFee, message, callOptions, revertOptions)
//	require.NoError(t, err)
//
//	// Automatically get the event signature
//	event := parsedABI.Events["Withdraw"]
//	eventSignatureHash := event.ID // Automatically retrieves the Keccak-256 hash of the event signature
//
//	// Generate the log automatically, including indexed topics
//	log := ethtypes.Log{
//		Address: common.HexToAddress("0xContractAddress"), // Contract address
//		Topics: []common.Hash{
//			eventSignatureHash,                 // Automatically retrieved event signature hash
//			common.BytesToHash(sender.Bytes()), // Indexed field: sender address
//		},
//		Data: data, // Packed data for non-indexed fields
//	}
//
//	// Log for debugging purposes
//	t.Logf("Generated log: %+v", log)
//}

func TestParseGatewayEvent(t *testing.T) {

}

func TestParseGatewayWithdrawalEvent(t *testing.T) {

}

func TestParseGatewayCallEvent(t *testing.T) {

}

func TestParseGatewayWithdrawAndCallEvent(t *testing.T) {

}

func TestNewWithdrawalInbound(t *testing.T) {

}

func TestNewCallInbound(t *testing.T) {

}

func TestNewWithdrawAndCallInbound(t *testing.T) {

}
