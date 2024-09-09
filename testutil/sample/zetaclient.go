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
		},
	}
}
