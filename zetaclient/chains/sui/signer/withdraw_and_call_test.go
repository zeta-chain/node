package signer

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/contracts/sui"
)

func Test_withdrawAndCallPTB(t *testing.T) {
	// ARRANGE
	signerAddrStr := "0x19c693d9b288073f5fd9572dd696604844cabab4c1637436338e023a5d67bc7a"
	gatewayPackageIDStr := "0x68efbe57639ca63db3bb6d37410f2963a81f9ba0cb4cfda205cff0b5ebd134ad"
	gatewayModule := "gateway"
	gatewayObjectIDStr := "0x71e93003b3b18a19bee11047778c4dcde65111cf27c0d76b42833704cf848d4b"
	withdrawCapIDStr := "0x0e5e4e15b982a46c40e33c22d0a528af4b61fa404a849938b2afa02f552dce20"
	coinTypeStr := string(sui.SUI)
	amountStr := "1000000"
	nonceStr := "1"
	gasBudgetStr := "20000000"
	receiver := "0xdfe280f111900991daa5c25d73a5ae134022b81616d3325945ac157b5394b598"
	message, err := hex.DecodeString("0ee37c85c59aecc644a3114a4182eacdf3e49435e588978b702191c531c397d6")
	require.NoError(t, err)
	encodedMessage := hex.EncodeToString(message)
	require.Equal(t, "0ee37c85c59aecc644a3114a4182eacdf3e49435e588978b702191c531c397d6", encodedMessage)
	cp := sui.CallPayload{
		TypeArgs: []string{"0x10bb172f3bd83b6ecd9861aeb843d4dc0549d3ca2085492d6810c5d768dfe880::token::TOKEN"},
		ObjectIDs: []string{
			"0xcf63b648079a11f00e6b738ad0b2226e49444d04fa6f6d585dbb636f551f323f",
			"0x382b6d92ed22bf13d61a6786ef9af505608e24c3572cadaad8eefb66e1366fc2",
			"0x2d4eb032bda2aaac27bac5a72c0b821fda8c3c50be90869f024959f2b350b053",
			"0x398a4b74e3ee229bb3deaf9facd2e435bcb8e216292ffe6e9616f77611718094",
		},
		Message: message,
	}

	tx, err := withdrawAndCallPTB(
		signerAddrStr,
		gatewayPackageIDStr,
		gatewayModule,
		gatewayObjectIDStr,
		withdrawCapIDStr,
		coinTypeStr,
		amountStr,
		nonceStr,
		gasBudgetStr,
		receiver,
		cp,
	)
	require.NoError(t, err)
	fmt.Printf("tx: %v\n", tx)
}
