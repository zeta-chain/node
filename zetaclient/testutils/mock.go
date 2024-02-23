package testutils

import (
	"github.com/btcsuite/btcutil"
	"github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
)

var _ interfaces.TSSSigner = (*MockTSS)(nil)

// MockTSS is a mock of TSS signer for testing
type MockTSS struct {
	evmAddress string
	btcAddress string
}

func NewMockTSS(evmAddress string, btcAddress string) *MockTSS {
	return &MockTSS{
		evmAddress: evmAddress,
		btcAddress: btcAddress,
	}
}

func NewMockTSSMainnet() *MockTSS {
	return NewMockTSS(TSSAddressEVMMainnet, TSSAddressBTCMainnet)
}

func NewMockTSSAthens3() *MockTSS {
	return NewMockTSS(TSSAddressEVMAthens3, TSSAddressBTCAthens3)
}

func (s *MockTSS) Sign(_ []byte, _ uint64, _ uint64, _ *common.Chain, _ string) ([65]byte, error) {
	return [65]byte{}, nil
}

func (s *MockTSS) Pubkey() []byte {
	return []byte{}
}

func (s *MockTSS) EVMAddress() ethcommon.Address {
	return ethcommon.HexToAddress(s.evmAddress)
}

func (s *MockTSS) BTCAddress() string {
	return s.btcAddress
}

func (s *MockTSS) BTCAddressWitnessPubkeyHash() *btcutil.AddressWitnessPubKeyHash {
	return nil
}

func (s *MockTSS) PubKeyCompressedBytes() []byte {
	return []byte{}
}

func MockCoreBridge() *zetabridge.ZetaCoreBridge {
	bridge, err := zetabridge.NewZetaCoreBridge(
		&keys.Keys{OperatorAddress: types.AccAddress{}},
		"127.0.0.1",
		"",
		"zetachain_7000-1",
		false,
		nil)
	if err != nil {
		panic(err)
	}
	return bridge
}
