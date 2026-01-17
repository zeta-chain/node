package sample

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONSamples(t *testing.T) {
	var (
		gatewayID = ton.MustParseAccountID("0:997d889c815aeac21c47f86ae0e38383efc3c3463067582f6263ad48c5a1485b")
		gw        = toncontracts.NewGateway(gatewayID)
	)

	t.Run("Donate", func(t *testing.T) {
		// ARRANGE
		d := toncontracts.Donation{
			Sender: GenerateTONAccountID(),
			Amount: sdkmath.NewUint(100_000_000),
		}

		tx := TONTransaction(t, TONDonateProps(t, gatewayID, d))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)

		d2, err := parsedTX.Donation()
		require.NoError(t, err)

		require.Equal(t, int(d.Amount.Uint64()), int(d2.Amount.Uint64()))
		require.Equal(t, d.Sender.ToRaw(), d2.Sender.ToRaw())
	})

	t.Run("Deposit", func(t *testing.T) {
		// ARRANGE
		d := toncontracts.Deposit{
			Sender:    GenerateTONAccountID(),
			Amount:    sdkmath.NewUint(200_000_000),
			Recipient: EthAddress(),
		}

		tx := TONTransaction(t, TONDepositProps(t, gatewayID, d))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)

		d2, err := parsedTX.Deposit()
		require.NoError(t, err)

		require.Equal(t, int(d.Amount.Uint64()), int(d2.Amount.Uint64()))
		require.Equal(t, d.Sender.ToRaw(), d2.Sender.ToRaw())
		require.Equal(t, d.Recipient.Hex(), d2.Recipient.Hex())
	})

	t.Run("Deposit and call", func(t *testing.T) {
		// ARRANGE
		d := toncontracts.DepositAndCall{
			Deposit: toncontracts.Deposit{
				Sender:    GenerateTONAccountID(),
				Amount:    sdkmath.NewUint(300_000_000),
				Recipient: EthAddress(),
			},
			CallData: []byte("Evidently, the most known and used kind of dictionaries in TON is hashmap."),
		}

		tx := TONTransaction(t, TONDepositAndCallProps(t, gatewayID, d))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)

		d2, err := parsedTX.DepositAndCall()
		require.NoError(t, err)

		require.Equal(t, int(d.Amount.Uint64()), int(d2.Amount.Uint64()))
		require.Equal(t, d.Sender.ToRaw(), d2.Sender.ToRaw())
		require.Equal(t, d.Recipient.Hex(), d2.Recipient.Hex())
		require.Equal(t, d.CallData, d2.CallData)
	})

}
