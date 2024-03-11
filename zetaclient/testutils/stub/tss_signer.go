package stub

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

var TestPrivateKey *ecdsa.PrivateKey

var _ interfaces.TSSSigner = (*TSS)(nil)

func init() {
	var err error
	TestPrivateKey, err = crypto.GenerateKey()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// TSS is a mock of TSS signer for testing
type TSS struct {
	evmAddress string
	btcAddress string
}

func NewMockTSS(evmAddress string, btcAddress string) *TSS {
	return &TSS{
		evmAddress: evmAddress,
		btcAddress: btcAddress,
	}
}

func NewTSSMainnet() *TSS {
	return NewMockTSS(testutils.TSSAddressEVMMainnet, testutils.TSSAddressBTCMainnet)
}

func NewTSSAthens3() *TSS {
	return NewMockTSS(testutils.TSSAddressEVMAthens3, testutils.TSSAddressBTCAthens3)
}

// Sign uses test key unrelated to any tss key in production
func (s *TSS) Sign(data []byte, _ uint64, _ uint64, _ *common.Chain, _ string) ([65]byte, error) {
	signature, err := crypto.Sign(data, TestPrivateKey)
	if err != nil {
		return [65]byte{}, err
	}
	var sigbyte [65]byte
	_ = copy(sigbyte[:], signature[:65])

	return sigbyte, nil
}

// Pubkey uses the hardcoded private test key to generate the public key in bytes
func (s *TSS) Pubkey() []byte {
	publicKey := TestPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("error casting public key to ECDSA")
	}
	return crypto.FromECDSAPub(publicKeyECDSA)
}

func (s *TSS) EVMAddress() ethcommon.Address {
	return ethcommon.HexToAddress(s.evmAddress)
}

func (s *TSS) BTCAddress() string {
	return s.btcAddress
}

func (s *TSS) BTCAddressWitnessPubkeyHash() *btcutil.AddressWitnessPubKeyHash {
	return nil
}

func (s *TSS) PubKeyCompressedBytes() []byte {
	return []byte{}
}
