//go:build drain

package drain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"
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
	mu        sync.Mutex
	signCalls int
	bcastTx   *eth.Transaction
	lastTo    ethcommon.Address
	lastAmt   *big.Int
	lastNonce uint64
	lastHt    uint64
}

func (m *mockEVMSigner) Chain() chains.Chain { return m.chain }
func (m *mockEVMSigner) SignDrainTx(_ context.Context, to ethcommon.Address, amount, gasPrice *big.Int, gasLimit, nonce, height uint64) (*eth.Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.signCalls++
	m.lastTo, m.lastAmt, m.lastNonce, m.lastHt = to, amount, nonce, height
	return eth.NewTx(&eth.LegacyTx{To: &to, Value: amount, GasPrice: gasPrice, Gas: gasLimit, Nonce: nonce}), nil
}
func (m *mockEVMSigner) BroadcastDrainTx(_ context.Context, tx *eth.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bcastTx = tx
	return nil
}

type mockBTCSigner struct {
	chain     chains.Chain
	mu        sync.Mutex
	signedTx  *wire.MsgTx
	inAmounts []int64
	height    uint64
	broadcast bool
}

func (m *mockBTCSigner) Chain() chains.Chain { return m.chain }
func (m *mockBTCSigner) SignTx(_ context.Context, tx *wire.MsgTx, inputAmounts []int64, height, _ uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.signedTx, m.inAmounts, m.height = tx, inputAmounts, height
	return nil
}
func (m *mockBTCSigner) Broadcast(_ context.Context, _ *wire.MsgTx) error {
	m.mu.Lock()
	defer m.mu.Unlock()
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

func evmResolver(m map[int64]EVMSigner) func(int64) (EVMSigner, bool) {
	return func(id int64) (EVMSigner, bool) { s, ok := m[id]; return s, ok }
}

func btcResolver(m map[int64]BTCSigner) func(int64) (BTCSigner, bool) {
	return func(id int64) (BTCSigner, bool) { s, ok := m[id]; return s, ok }
}

func btcInputs() []draintx.BTCInput {
	return []draintx.BTCInput{
		{TxID: "3ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 0, AmountSats: 45_000_000},
		{TxID: "4ba58f8f2f3f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f2f", Vout: 1, AmountSats: 45_010_000},
	}
}

func signedPayload(t *testing.T, priv *ecdsa.PrivateKey, final bool, triggerHeight int64, btcTo string) draintx.Payload {
	in := btcInputs()
	var total int64
	for _, i := range in {
		total += i.AmountSats
	}
	fee := int64(10_000)
	p := draintx.Payload{
		TriggerZetaHeight: triggerHeight,
		Seq:               1,
		Final:             final,
		EVMTxs: []draintx.EVMTx{
			{ChainID: chains.Ethereum.ChainId, To: evmReceiverHex, Nonce: 5, Amount: "1000", GasPrice: "250000", GasLimit: 21000},
		},
		BTCTxs: []draintx.BTCTx{
			{ChainID: chains.BitcoinRegtest.ChainId, To: btcTo, OutputSats: total - fee, FeeSats: fee, Inputs: in},
		},
	}
	require.NoError(t, p.Sign(priv))
	return p
}

func newTestPoller(cfg Config) *Poller {
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
	p := newTestPoller(Config{Window: 5})
	tests := []struct {
		current, trigger int64
		fire, missed     bool
	}{
		{99, 100, false, false},
		{100, 100, true, false},
		{104, 100, true, false},
		{105, 100, false, true},
		{200, 100, false, true},
	}
	for _, tc := range tests {
		fire, missed := p.readyToFire(tc.current, tc.trigger)
		require.Equal(t, tc.fire, fire, "fire c=%d t=%d", tc.current, tc.trigger)
		require.Equal(t, tc.missed, missed, "missed c=%d t=%d", tc.current, tc.trigger)
	}
}

func TestStepRejectsNonFinal(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(Config{
		Fetcher:          mockFetcher{signedPayload(t, priv, false, 100, recv.EncodeAddress())},
		Height:           mockHeight(100),
		PubKey:           ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
		ResolveBTCSigner: btcResolver(nil),
	})

	var active *activeDrain
	p.step(context.Background(), &active)
	require.Nil(t, active)
	require.Zero(t, evm.signCalls)
}

func TestStepRejectsBadSignature(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	other, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(Config{
		Fetcher:          mockFetcher{signedPayload(t, priv, true, 100, recv.EncodeAddress())},
		Height:           mockHeight(100),
		PubKey:           ethcrypto.CompressPubkey(&other.PublicKey),
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
		ResolveBTCSigner: btcResolver(nil),
	})

	var active *activeDrain
	p.step(context.Background(), &active)
	require.Nil(t, active)
	require.Zero(t, evm.signCalls)
}

func TestStepMissedWindow(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(Config{
		Fetcher:          mockFetcher{signedPayload(t, priv, true, 100, recv.EncodeAddress())},
		Height:           mockHeight(1000),
		PubKey:           ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
		ResolveBTCSigner: btcResolver(nil),
	})

	var active *activeDrain
	p.step(context.Background(), &active)
	require.Nil(t, active)
	require.Zero(t, evm.signCalls)
	// a missed payload is marked handled so it isn't reconsidered
	require.EqualValues(t, 100, p.lastFiredHeight)
}

func TestStepFiresAndCompletes(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	btc := &mockBTCSigner{chain: chains.BitcoinRegtest}
	p := newTestPoller(Config{
		Fetcher:          mockFetcher{signedPayload(t, priv, true, 100, recv.EncodeAddress())},
		Height:           mockHeight(100),
		PubKey:           ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
		ResolveBTCSigner: btcResolver(map[int64]BTCSigner{chains.BitcoinRegtest.ChainId: btc}),
	})

	var active *activeDrain
	p.step(context.Background(), &active)
	require.Equal(t, 1, evm.signCalls)
	require.EqualValues(t, 100, evm.lastHt)
	require.True(t, btc.broadcast)
	require.EqualValues(t, 100, btc.height)
	// once complete, the active drain resets so the poller keeps polling for the next payload
	require.Nil(t, active)
	require.EqualValues(t, 100, p.lastFiredHeight)
}

func TestStepRefiresAtHigherHeight(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	btc := &mockBTCSigner{chain: chains.BitcoinRegtest}
	p := newTestPoller(Config{
		Fetcher:          mockFetcher{signedPayload(t, priv, true, 100, recv.EncodeAddress())},
		Height:           mockHeight(100),
		PubKey:           ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
		ResolveBTCSigner: btcResolver(map[int64]BTCSigner{chains.BitcoinRegtest.ChainId: btc}),
	})

	ctx := context.Background()
	var active *activeDrain

	// fires at H=100 and completes
	p.step(ctx, &active)
	require.Nil(t, active)
	require.EqualValues(t, 100, p.lastFiredHeight)
	require.Equal(t, 1, evm.signCalls)

	// same-height payload is already handled -> ignored, nothing fires
	evm.signCalls = 0
	p.step(ctx, &active)
	require.Nil(t, active)
	require.Zero(t, evm.signCalls)
	require.EqualValues(t, 100, p.lastFiredHeight)

	// operator republishes the remaining chains at a higher height -> fires again
	p.Fetcher = mockFetcher{signedPayload(t, priv, true, 200, recv.EncodeAddress())}
	p.Height = mockHeight(200)
	p.step(ctx, &active)
	require.EqualValues(t, 200, p.lastFiredHeight)
	require.Equal(t, 1, evm.signCalls)
}

func TestStepFailsClosedOnMissingSigner(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	// BTC signer missing: the whole payload must fail closed (EVM must not fire either)
	p := newTestPoller(Config{
		Fetcher:          mockFetcher{signedPayload(t, priv, true, 100, recv.EncodeAddress())},
		Height:           mockHeight(100),
		PubKey:           ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
		ResolveBTCSigner: btcResolver(nil),
	})

	ctx := context.Background()
	var active *activeDrain

	// nothing fires, and lastFiredHeight stays 0 so a retry can still fire within the window
	p.step(ctx, &active)
	require.Nil(t, active)
	require.Zero(t, evm.signCalls)
	require.Zero(t, p.lastFiredHeight)

	// BTC signer comes up within the window -> the payload now fires
	btc := &mockBTCSigner{chain: chains.BitcoinRegtest}
	p.ResolveBTCSigner = btcResolver(map[int64]BTCSigner{chains.BitcoinRegtest.ChainId: btc})
	p.step(ctx, &active)
	require.EqualValues(t, 100, p.lastFiredHeight)
	require.Equal(t, 1, evm.signCalls)
	require.True(t, btc.broadcast)
}

func TestExecuteEVM(t *testing.T) {
	recv := btcReceiver(t)
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(Config{
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
	})

	t.Run("happy path signs to receiver at height", func(t *testing.T) {
		tx := draintx.EVMTx{ChainID: chains.Ethereum.ChainId, To: evmReceiverHex, Nonce: 5, Amount: "1000", GasPrice: "250000", GasLimit: 21000}
		require.NoError(t, p.executeEVM(context.Background(), tx, 100))
		require.Equal(t, 1, evm.signCalls)
		require.Equal(t, ethcommon.HexToAddress(evmReceiverHex), evm.lastTo)
		require.EqualValues(t, 100, evm.lastHt)
		require.NotNil(t, evm.bcastTx)
	})

	t.Run("rejects wrong receiver", func(t *testing.T) {
		evm.signCalls, evm.bcastTx = 0, nil
		tx := draintx.EVMTx{ChainID: chains.Ethereum.ChainId, To: "0x2222222222222222222222222222222222222222", Amount: "1000", GasPrice: "250000", GasLimit: 21000}
		require.Error(t, p.executeEVM(context.Background(), tx, 100))
		require.Zero(t, evm.signCalls)
	})
}

func TestExecuteEVMRejectsZeroReceiver(t *testing.T) {
	evm := &mockEVMSigner{chain: chains.Ethereum}
	p := newTestPoller(Config{
		EVMReceiver:      ethcommon.Address{}, // zero
		ResolveEVMSigner: evmResolver(map[int64]EVMSigner{chains.Ethereum.ChainId: evm}),
	})
	tx := draintx.EVMTx{ChainID: chains.Ethereum.ChainId, To: "0x0000000000000000000000000000000000000000", Amount: "1000", GasPrice: "250000", GasLimit: 21000}
	require.Error(t, p.executeEVM(context.Background(), tx, 100))
	require.Zero(t, evm.signCalls)
}

func TestExecuteBTC(t *testing.T) {
	recv := btcReceiver(t)
	btc := &mockBTCSigner{chain: chains.BitcoinRegtest}
	p := newTestPoller(Config{
		BTCReceiver:      recv,
		ResolveBTCSigner: btcResolver(map[int64]BTCSigner{chains.BitcoinRegtest.ChainId: btc}),
	})
	in := btcInputs()
	var total int64
	for _, i := range in {
		total += i.AmountSats
	}

	t.Run("happy path builds sweep and signs at height", func(t *testing.T) {
		tx := draintx.BTCTx{ChainID: chains.BitcoinRegtest.ChainId, To: recv.EncodeAddress(), OutputSats: total - 10_000, FeeSats: 10_000, Inputs: in}
		require.NoError(t, p.executeBTC(context.Background(), tx, 100))
		require.NotNil(t, btc.signedTx)
		require.Len(t, btc.signedTx.TxIn, 2)
		require.Len(t, btc.signedTx.TxOut, 1)
		require.EqualValues(t, total-10_000, btc.signedTx.TxOut[0].Value)
		require.EqualValues(t, 100, btc.height)
		require.True(t, btc.broadcast)
	})

	t.Run("rejects wrong receiver", func(t *testing.T) {
		btc.signedTx, btc.broadcast = nil, false
		tx := draintx.BTCTx{ChainID: chains.BitcoinRegtest.ChainId, To: "bcrt1qwrong", OutputSats: total - 10_000, FeeSats: 10_000, Inputs: in}
		require.Error(t, p.executeBTC(context.Background(), tx, 100))
		require.Nil(t, btc.signedTx)
	})

	t.Run("rejects fee/amount mismatch (burn-to-miners)", func(t *testing.T) {
		btc.signedTx = nil
		tx := draintx.BTCTx{ChainID: chains.BitcoinRegtest.ChainId, To: recv.EncodeAddress(), OutputSats: 1, FeeSats: 10_000, Inputs: in}
		require.Error(t, p.executeBTC(context.Background(), tx, 100))
		require.Nil(t, btc.signedTx)
	})
}

func TestValidateBTCFee(t *testing.T) {
	in := btcInputs()
	var total int64
	for _, i := range in {
		total += i.AmountSats
	}

	require.NoError(t, validateBTCFee(draintx.BTCTx{OutputSats: total - 10_000, FeeSats: 10_000, Inputs: in}))
	require.Error(t, validateBTCFee(draintx.BTCTx{OutputSats: total, FeeSats: 10_000, Inputs: in}))        // inconsistent
	require.Error(t, validateBTCFee(draintx.BTCTx{OutputSats: total / 2, FeeSats: total / 2, Inputs: in})) // excessive fee
	require.Error(t, validateBTCFee(draintx.BTCTx{OutputSats: 1, FeeSats: total - 1, Inputs: in}))         // dust output
}

func TestBuildSweep(t *testing.T) {
	recv := btcReceiver(t)
	in := btcInputs()

	tx, amounts, err := buildSweep(recv, in, 90_000_000)
	require.NoError(t, err)
	require.Len(t, tx.TxIn, 2)
	require.EqualValues(t, rbfSequenceNum, tx.TxIn[0].Sequence)
	require.Len(t, tx.TxOut, 1)
	require.EqualValues(t, 90_000_000, tx.TxOut[0].Value)
	require.Equal(t, []int64{45_000_000, 45_010_000}, amounts)

	// deterministic: identical inputs -> identical tx
	tx2, _, err := buildSweep(recv, in, 90_000_000)
	require.NoError(t, err)
	require.Equal(t, tx.TxHash(), tx2.TxHash())

	_, _, err = buildSweep(recv, nil, 1)
	require.Error(t, err)
}

func TestRunStopsOnContextCancel(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	recv := btcReceiver(t)
	// non-final payload so the poller never fires and just keeps polling until cancel
	p := newTestPoller(Config{
		Fetcher:          mockFetcher{signedPayload(t, priv, false, 100, recv.EncodeAddress())},
		Height:           mockHeight(100),
		PubKey:           ethcrypto.CompressPubkey(&priv.PublicKey),
		EVMReceiver:      ethcommon.HexToAddress(evmReceiverHex),
		BTCReceiver:      recv,
		ResolveEVMSigner: evmResolver(nil),
		ResolveBTCSigner: btcResolver(nil),
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { p.Run(ctx); close(done) }()

	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not stop on context cancel")
	}
}
