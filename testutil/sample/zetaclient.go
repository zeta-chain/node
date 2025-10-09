package sample

import (
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/types"
)

const (
	// These are sample restricted addresses for e2e tests.
	RestrictedEVMAddressTest = "0x8a81Ba8eCF2c418CAe624be726F505332DF119C6"
	RestrictedBtcAddressTest = "bcrt1qzp4gt6fc7zkds09kfzaf9ln9c5rvrzxmy6qmpp"
	RestrictedSolAddressTest = "9fA4vYZfCa9k9UHjnvYCk4YoipsooapGciKMgaTBw9UH"
	RestrictedSuiAddressTest = "0x14454c46e2ac4603adaa15df30e5dbf7662c3177db4b83c326bed5663d25d1bd"
	RestrictedTonAddressTest = "0:fffffbd865df68188ea84d6615086c26a7b5912a60bc55fded2cdb029b67cdef"
	RevertAddressZEVM        = "0x4c40813A6a726FE9353a4A691ecFe2D8641914F7"
)

// InboundEvent returns a sample InboundEvent.
func InboundEvent(chainID int64, sender string, receiver string, amount uint64, memo []byte) *types.InboundEvent {
	r := newRandFromSeed(chainID)

	return &types.InboundEvent{
		SenderChainID: chainID,
		Sender:        sender,
		Receiver:      receiver,
		TxOrigin:      sender,
		Amount:        amount,
		Memo:          memo,
		BlockNumber:   r.Uint64(),
		TxHash:        StringRandom(r, 32),
		Index:         0,
		CoinType:      coin.CoinType(r.Intn(100)),
		Asset:         StringRandom(r, 32),
	}
}

// ComplianceConfig returns a sample compliance config
func ComplianceConfig() config.ComplianceConfig {
	return config.ComplianceConfig{
		RestrictedAddresses: []string{
			RestrictedEVMAddressTest,
			RestrictedBtcAddressTest,
			RestrictedSolAddressTest,
			RestrictedSuiAddressTest,
			RestrictedTonAddressTest,
		},
	}
}

// FeatureFlags returns a sample feature flags
func FeatureFlags() config.FeatureFlags {
	return config.FeatureFlags{
		EnableMultipleCalls:            true,
		EnableSolanaAddressLookupTable: true,
	}
}
