//go:build drain

package drain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	ethcommon "github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/draintx"
)

type mockEVMSigner struct {
	chain     chains.Chain
	signCalls int
	bcastTx   *eth.Transaction
	lastTo    ethcommon.Address
	lastAmt   *big.Int
	lastNonce uint64
	lastHt    uint64
}

func (m *mockEVMSigner) Chain() chains.Chain { return m.chain }
func (m *mockEVMSigner) SignDrainTx(_ context.Context, to ethcommon.Address, amount, gasPrice *big.Int, gasLimit, nonce, height uint64) (*eth.Transaction, error) {
	m.signCalls++
	m.lastTo, m.lastAmt, m.lastNonce, m.lastHt = to, amount, nonce, height
	return eth.NewTx(&eth.LegacyTx{To: &to, Value: amount, GasPrice: gasPrice, Gas: gasLimit, Nonce: nonce}), nil
}
func (m *mockEVMSigner) BroadcastDrainTx(_ context.Context, tx *eth.Transaction) error {
	m.bcastTx = tx
	return nil
}

type mockBTCSigner struct {
	chain     chains.Chain
	signedTx  *wire.MsgTx
	inAmounts []int64
	height    uint64
	broadcast bool
}

func (m *mockBTCSigner) Chain() chains.Chain { return m.chain }
func (m *mockBTCSigner) SignTx(_ context.Context, tx *wire.MsgTx, inputAmounts []int64, height, _ uint64) error {
	m.signedTx, m.inAmounts, m.height = tx, inputAmounts, height
	return nil
}
func (m *mockBTCSigner) Broadcast(_ context.Context, _ *wire.MsgTx) error {
	m.broadcast = true
	return nil
}

type mockHeight int64

func (h mockHeight) GetBlockHeight(context.Context) (int64, error) { return int64(h), nil }

type mockFetcher struct{ p draintx.Payload }

func (f mockFetcher) Fetch(context.Context) (draintx.Payload, error) { return f.p, nil }

const evmReceiverHex = "0x1111111111111111111111111111111111111111"

func btcReceiver(t *testing.T) btcutil.Address {
	addr, err := btcutil.NewAddressWitnessPubKeyHash(make([]byte, 20), &chaincfg.RegressionNetParams)
	require.NoError(t, err)
	return addr
}

func signedPayload(t *testing.T, priv *ecdsa.PrivateKey, final bool, triggerHeight int64, btcTo string) draintx.Payload {
	p := draintx.Payload{
		TriggerZetaHeight: triggerHeight,
		Seq:               1,
		Final:             final,
		EVMTxs: []draintx.EVMTx{
			{ChainID: chains.Ethereum.ChainId, To: evmReceiverHex, Nonce: 5, Amount: "1000", GasPrice: "250000", GasLimit: 21000},
		},
		BTCTxs: []draintx.BTCTx{
			{
				ChainID: chains.BitcoinRegtest.ChainId, To: btcTo, OutputSats: 90_000_000, FeeSats: 10_000,
				Inputs: []draintx.BTCInput{
					{TxID: "3ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 0, AmountSats: 45_000_000},
					{TxID: "4ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 1, AmountSats: 45_010_000},
				},
			},
		},
	}
	require.NoError(t, p.Sign(priv))
	return p
}

func newTestPoller(t *testing.T, cfg Config) *Poller {
	if cfg.PollInterval == 0 {
		cfg.PollInterval = time.Millisecond
	}
	if cfg.Window == 0 {
		cfg.Window = 5
	}
	cfg.Logger = zerolog.Nop()
	return New(cfg)
}

func TestReadyToFire(t *testing.T) {
	p := newTestPoller(t, Config{Window: 5})
	tests := []struct {
		current, trigger int64
		fire, missed     bool
	}{
		{99, 100, false, false}, // before window
		{100, 100, true, false}, // at trigger
		{104, 100, true, false}, // inside window
		{105, 100, false, true}, // window elapsed -> missed
		{200, 100, false, true}, // long past
	}
	for _, tc := range tests {
		fire, missed := p.readyToFire(tc.current, tc.trigger)
		require.Equal(t, tc.fire, fire, "fire for current=%d trigger=%d", tc.current, tc.trigger)
		require.Equal(t, tc.missed, missed, "missed for current=%d trigger=%d", tc.current, tc.trigger)
	}
}

func TestTickRejectsNonFinal(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(t, Config{
		Fetcher:     mockFetcher{signedPayload(t, priv, false, 100, recv.EncodeAddress())},
		Height:      mockHeight(100),
		PubKey:      ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver: ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver: recv,
		EVMSigners:  map[int64]EVMSigner{chains.Ethereum.ChainId: evm},
	})

	require.False(t, p.tick(context.Background()))
	require.Zero(t, evm.signCalls)
}

func TestTickRejectsBadSignature(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	other, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(t, Config{
		Fetcher:     mockFetcher{signedPayload(t, priv, true, 100, recv.EncodeAddress())},
		Height:      mockHeight(100),
		PubKey:      ethcrypto.CompressPubkey(&other.PublicKey), // wrong key
		EVMReceiver: ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver: recv,
		EVMSigners:  map[int64]EVMSigner{chains.Ethereum.ChainId: evm},
	})

	require.False(t, p.tick(context.Background()))
	require.Zero(t, evm.signCalls)
}

func TestTickMissedWindow(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(t, Config{
		Fetcher:     mockFetcher{signedPayload(t, priv, true, 100, recv.EncodeAddress())},
		Height:      mockHeight(1000), // way past
		PubKey:      ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver: ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver: recv,
		EVMSigners:  map[int64]EVMSigner{chains.Ethereum.ChainId: evm},
	})

	require.True(t, p.tick(context.Background())) // done (missed)
	require.Zero(t, evm.signCalls)
}

func TestExecuteEVM(t *testing.T) {
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(t, Config{
		EVMReceiver: ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver: recv,
		EVMSigners:  map[int64]EVMSigner{chains.Ethereum.ChainId: evm},
	})

	t.Run("happy path signs and broadcasts to receiver", func(t *testing.T) {
		tx := draintx.EVMTx{ChainID: chains.Ethereum.ChainId, To: evmReceiverHex, Nonce: 5, Amount: "1000", GasPrice: "250000", GasLimit: 21000}
		require.NoError(t, p.executeEVM(context.Background(), tx, 100))
		require.Equal(t, 1, evm.signCalls)
		require.Equal(t, ethcommon.HexToAddress(evmReceiverHex), evm.lastTo)
		require.Equal(t, "1000", evm.lastAmt.String())
		require.EqualValues(t, 5, evm.lastNonce)
		require.EqualValues(t, 100, evm.lastHt)
		require.NotNil(t, evm.bcastTx)
	})

	t.Run("rejects wrong receiver", func(t *testing.T) {
		evm.signCalls, evm.bcastTx = 0, nil
		tx := draintx.EVMTx{ChainID: chains.Ethereum.ChainId, To: "0x2222222222222222222222222222222222222222", Nonce: 5, Amount: "1000", GasPrice: "250000", GasLimit: 21000}
		require.Error(t, p.executeEVM(context.Background(), tx, 100))
		require.Zero(t, evm.signCalls)
		require.Nil(t, evm.bcastTx)
	})
}

func TestExecuteBTC(t *testing.T) {
	recv := btcReceiver(t)
	btc := &mockBTCSigner{chain: chains.BitcoinRegtest}
	p := newTestPoller(t, Config{
		BTCReceiver: recv,
		BTCSigners:  map[int64]BTCSigner{chains.BitcoinRegtest.ChainId: btc},
	})
	inputs := []draintx.BTCInput{
		{TxID: "3ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 0, AmountSats: 45_000_000},
		{TxID: "4ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 1, AmountSats: 45_010_000},
	}

	t.Run("happy path builds sweep and signs at height", func(t *testing.T) {
		tx := draintx.BTCTx{ChainID: chains.BitcoinRegtest.ChainId, To: recv.EncodeAddress(), OutputSats: 90_000_000, FeeSats: 10_000, Inputs: inputs}
		require.NoError(t, p.executeBTC(context.Background(), tx, 100))
		require.NotNil(t, btc.signedTx)
		require.Len(t, btc.signedTx.TxIn, 2)
		require.Len(t, btc.signedTx.TxOut, 1)
		require.EqualValues(t, 90_000_000, btc.signedTx.TxOut[0].Value)
		require.Equal(t, []int64{45_000_000, 45_010_000}, btc.inAmounts)
		require.EqualValues(t, 100, btc.height)
		require.True(t, btc.broadcast)
	})

	t.Run("rejects wrong receiver", func(t *testing.T) {
		btc.signedTx, btc.broadcast = nil, false
		tx := draintx.BTCTx{ChainID: chains.BitcoinRegtest.ChainId, To: "bcrt1qwrongwrongwrong", OutputSats: 90_000_000, FeeSats: 10_000, Inputs: inputs}
		require.Error(t, p.executeBTC(context.Background(), tx, 100))
		require.Nil(t, btc.signedTx)
		require.False(t, btc.broadcast)
	})
}

func TestBuildSweep(t *testing.T) {
	recv := btcReceiver(t)
	inputs := []draintx.BTCInput{
		{TxID: "3ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 0, AmountSats: 45_000_000},
		{TxID: "4ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 3, AmountSats: 45_010_000},
	}

	tx, amounts, err := buildSweep(recv, inputs, 90_000_000)
	require.NoError(t, err)
	require.Len(t, tx.TxIn, 2)
	require.EqualValues(t, rbfSequenceNum, tx.TxIn[0].Sequence)
	require.EqualValues(t, 3, tx.TxIn[1].PreviousOutPoint.Index)
	require.Len(t, tx.TxOut, 1)
	require.EqualValues(t, 90_000_000, tx.TxOut[0].Value)
	require.Equal(t, []int64{45_000_000, 45_010_000}, amounts)

	_, _, err = buildSweep(recv, nil, 1)
	require.Error(t, err)
}
