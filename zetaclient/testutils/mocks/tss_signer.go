package mocks

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

// TestPrivateKey is a random private key for testing
var TestPrivateKey *ecdsa.PrivateKey

// init generates a random private key for testing
func init() {
	var err error
	TestPrivateKey, err = crypto.GenerateKey()
	if err != nil {
		fmt.Println(err.Error())
	}
}

var _ interfaces.TSSSigner = (*TSS)(nil)

// TSS is a mock of TSS signer for testing
type TSS struct {
	paused bool

	// set evmAddress/btcAddress if just want to mock EVMAddress()/BTCAddress()
	chain      chains.Chain
	evmAddress string
	btcAddress string

	// set PrivKey if you want to use a specific private key
	PrivKey *ecdsa.PrivateKey
}

func NewMockTSS(chain chains.Chain, evmAddress string, btcAddress string) *TSS {
	return &TSS{
		paused:     false,
		chain:      chain,
		evmAddress: evmAddress,
		btcAddress: btcAddress,
		PrivKey:    TestPrivateKey,
	}
}

func NewTSSMainnet() *TSS {
	return NewMockTSS(chains.BitcoinMainnet, testutils.TSSAddressEVMMainnet, testutils.TSSAddressBTCMainnet)
}

func NewTSSAthens3() *TSS {
	return NewMockTSS(chains.BscTestnet, testutils.TSSAddressEVMAthens3, testutils.TSSAddressBTCAthens3)
}

func NewGeneratedTSS(t *testing.T, chain chains.Chain) *TSS {
	pk, err := crypto.GenerateKey()
	require.NoError(t, err)

	btcPub, err := btcec.ParsePubKey(crypto.FromECDSAPub(&pk.PublicKey))
	require.NoError(t, err)

	btcAddress, err := btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(btcPub.SerializeCompressed()),
		&chaincfg.TestNet3Params,
	)

	require.NoError(t, err)

	return &TSS{
		paused:     false,
		chain:      chain,
		evmAddress: crypto.PubkeyToAddress(pk.PublicKey).Hex(),
		btcAddress: btcAddress.String(),
		PrivKey:    pk,
	}
}

// WithPrivKey sets the private key for the TSS
func (s *TSS) WithPrivKey(privKey *ecdsa.PrivateKey) *TSS {
	s.PrivKey = privKey
	return s
}

// Sign uses test key unrelated to any tss key in production
func (s *TSS) Sign(_ context.Context, data []byte, _ uint64, _ uint64, _ int64, _ string) ([65]byte, error) {
	// return error if tss is paused
	if s.paused {
		return [65]byte{}, fmt.Errorf("tss is paused")
	}

	signature, err := crypto.Sign(data, s.PrivKey)
	if err != nil {
		return [65]byte{}, err
	}
	var sigbyte [65]byte
	_ = copy(sigbyte[:], signature[:65])

	return sigbyte, nil
}

// SignBatch uses test key unrelated to any tss key in production
func (s *TSS) SignBatch(_ context.Context, _ [][]byte, _ uint64, _ uint64, _ int64) ([][65]byte, error) {
	// return error if tss is paused
	if s.paused {
		return nil, fmt.Errorf("tss is paused")
	}

	// mock not implemented yet
	return nil, fmt.Errorf("not implemented")
}

func (s *TSS) Pubkey() []byte {
	publicKeyBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	return publicKeyBytes
}

func (s *TSS) EVMAddress() ethcommon.Address {
	// force use evmAddress if set
	if s.evmAddress != "" {
		return ethcommon.HexToAddress(s.evmAddress)
	}
	return crypto.PubkeyToAddress(s.PrivKey.PublicKey)
}

func (s *TSS) EVMAddressList() []ethcommon.Address {
	return []ethcommon.Address{s.EVMAddress()}
}

func (s *TSS) BTCAddress(_ int64) (*btcutil.AddressWitnessPubKeyHash, error) {
	// return error if tss is paused
	if s.paused {
		return nil, fmt.Errorf("tss is paused")
	}

	// force use static btcAddress if set
	if s.btcAddress != "" {
		net, err := chains.GetBTCChainParams(s.chain.ChainId)
		if err != nil {
			return nil, err
		}
		addr, err := btcutil.DecodeAddress(s.btcAddress, net)
		if err != nil {
			return nil, err
		}
		return addr.(*btcutil.AddressWitnessPubKeyHash), nil
	}
	// if privkey is set, use it to generate a segwit address
	if s.PrivKey != nil {
		pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
		pk, err := btcec.ParsePubKey(pkBytes)
		if err != nil {
			fmt.Printf("error parsing pubkey: %v", err)
			return nil, err
		}

		// witness program: https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki#Witness_program
		// The HASH160 of the public key must match the 20-byte witness program.
		addrWPKH, err := btcutil.NewAddressWitnessPubKeyHash(
			btcutil.Hash160(pk.SerializeCompressed()),
			&chaincfg.TestNet3Params,
		)
		if err != nil {
			fmt.Printf("error NewAddressWitnessPubKeyHash: %v", err)
			return nil, err
		}

		return addrWPKH, nil
	}
	return nil, nil
}

// PubKeyCompressedBytes returns 33B compressed pubkey
func (s *TSS) PubKeyCompressedBytes() []byte {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	pk, err := btcec.ParsePubKey(pkBytes)
	if err != nil {
		fmt.Printf("error parsing pubkey: %v", err)
		return nil
	}
	return pk.SerializeCompressed()
}

// ----------------------------------------------------------------------------
// methods to control the mock for testing
// ----------------------------------------------------------------------------
func (s *TSS) Pause() {
	s.paused = true
}

func (s *TSS) Unpause() {
	s.paused = false
}
