//go:build drain

// Package drain implements the emergency TSS native-fund drain poller. It is compiled
// only under the `drain` build tag and is off by default: operators run a dedicated
// drain build during the drain window and upgrade away afterwards.
//
// The poller makes no decisions about the transaction. It fetches an operator-signed,
// fully-resolved payload, verifies it against a baked-in public key, asserts every tx
// sends only to the compiled-in safe receiver, and signs the pinned values at the pinned
// trigger height. A compromised endpoint can change when funds move, never where.
package drain

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	ethcommon "github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	pkgdrain "github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/draintx"
)

const (
	// btcKeysignNonce is a fixed keysign nonce for BTC sweeps. The go-tss ceremony matches
	// on digests + height, not nonce (nonce is cosmetic), so a constant is safe.
	btcKeysignNonce = 0
	// rbfSequenceNum opts the first input into full-RBF, mirroring the production BTC signer.
	rbfSequenceNum uint32 = 1
)

// EVMSigner is the subset of the EVM signer the poller drives.
type EVMSigner interface {
	Chain() chains.Chain
	SignDrainTx(
		ctx context.Context,
		to ethcommon.Address,
		amount, gasPrice *big.Int,
		gasLimit, nonce, height uint64,
	) (*eth.Transaction, error)
	BroadcastDrainTx(ctx context.Context, tx *eth.Transaction) error
}

// BTCSigner is the subset of the Bitcoin signer the poller drives.
type BTCSigner interface {
	Chain() chains.Chain
	SignTx(ctx context.Context, tx *wire.MsgTx, inputAmounts []int64, height, nonce uint64) error
	Broadcast(ctx context.Context, tx *wire.MsgTx) error
}

// HeightProvider returns the current zeta block height.
type HeightProvider interface {
	GetBlockHeight(ctx context.Context) (int64, error)
}

// Fetcher fetches the current drain payload from the operator endpoint.
type Fetcher interface {
	Fetch(ctx context.Context) (draintx.Payload, error)
}

// Config configures a Poller.
type Config struct {
	Fetcher Fetcher
	Height  HeightProvider
	PubKey  []byte // baked-in operator public key
	// EVMReceiver/BTCReceiver are the compiled-in security anchors; every tx must send here.
	EVMReceiver ethcommon.Address
	BTCReceiver btcutil.Address
	// ResolveEVMSigner/ResolveBTCSigner resolve the live per-chain signer at fire time, so
	// a signer that bootstraps late is still picked up on a retry.
	ResolveEVMSigner func(chainID int64) (EVMSigner, bool)
	ResolveBTCSigner func(chainID int64) (BTCSigner, bool)
	Window           int64 // blocks after H during which a node may still fire/retry
	PollInterval     time.Duration
	Logger           zerolog.Logger
}

// Poller polls the drain endpoint and fires the drain at the trigger height, retrying any
// failed txs across the firing window.
type Poller struct {
	Config
	// lastFiredHeight is the trigger height of the most recent payload the poller acted on
	// (fired or gave up as missed). A fetched payload at or below it is ignored, so an old
	// payload can never re-fire; the operator retries a partial drain by republishing the
	// remaining chains at a NEW, higher trigger height.
	lastFiredHeight int64
}

// New creates a Poller.
func New(cfg Config) *Poller {
	return &Poller{Config: cfg}
}

// activeDrain tracks the in-progress drain and which txs still need broadcasting.
type activeDrain struct {
	payload draintx.Payload
	mu      sync.Mutex
	evm     []*evmItem
	btc     []*btcItem
}

type evmItem struct {
	tx   draintx.EVMTx
	done bool
}

type btcItem struct {
	tx   draintx.BTCTx
	done bool
}

func newActiveDrain(p draintx.Payload) *activeDrain {
	a := &activeDrain{payload: p}
	for _, tx := range p.EVMTxs {
		a.evm = append(a.evm, &evmItem{tx: tx})
	}
	for _, tx := range p.BTCTxs {
		a.btc = append(a.btc, &btcItem{tx: tx})
	}
	return a
}

func (a *activeDrain) allDone() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, it := range a.evm {
		if !it.done {
			return false
		}
	}
	for _, it := range a.btc {
		if !it.done {
			return false
		}
	}
	return true
}

// Run polls until the context is cancelled. It never exits after a single payload: once a
// drain completes or its window elapses it resets and keeps polling, so a fresh payload at a
// higher trigger height (an operator retrying the remaining chains) fires on its own.
func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.PollInterval)
	defer ticker.Stop()

	p.Logger.Info().Msg("drain poller started")

	var active *activeDrain
	for {
		select {
		case <-ctx.Done():
			p.Logger.Info().Msg("drain poller stopped")
			return
		case <-ticker.C:
			p.step(ctx, &active)
		}
	}
}

// step runs one poll iteration: arm on an eligible newer payload, then push its pending txs.
// When the active drain finishes or its window elapses it resets *active to nil so the poller
// keeps looking for the next payload.
func (p *Poller) step(ctx context.Context, active **activeDrain) {
	if *active == nil {
		p.maybeArm(ctx, active)
		if *active == nil {
			return
		}
	}

	a := *active
	p.attemptPending(ctx, a)
	if a.allDone() {
		p.logSummary(a, "drain complete")
		*active = nil
		return
	}

	current, err := p.Height.GetBlockHeight(ctx)
	if err != nil {
		// can't tell if the window elapsed; keep retrying
		return
	}
	if current >= a.payload.TriggerZetaHeight+p.Window {
		p.logSummary(a, "drain window elapsed with unfinished txs")
		*active = nil
	}
}

// maybeArm fetches the current final payload and, if it is newer than the last one handled and
// inside its firing window with all signers ready, arms it into *active.
func (p *Poller) maybeArm(ctx context.Context, active **activeDrain) {
	payload, ok := p.fetchFinal(ctx)
	if !ok {
		return
	}
	if payload.TriggerZetaHeight <= p.lastFiredHeight {
		p.Logger.Debug().
			Int64("trigger_height", payload.TriggerZetaHeight).
			Int64("last_fired_height", p.lastFiredHeight).
			Msg("ignoring already-handled drain payload")
		return
	}
	current, err := p.Height.GetBlockHeight(ctx)
	if err != nil {
		p.Logger.Warn().Err(err).Msg("unable to get zeta block height")
		return
	}
	fire, missed := p.readyToFire(current, payload.TriggerZetaHeight)
	switch {
	case missed:
		// past the window: mark it handled so it isn't reconsidered; operator must reset higher.
		p.Logger.Error().
			Int64("trigger_height", payload.TriggerZetaHeight).
			Int64("current_height", current).
			Msg("drain trigger height missed, ignoring")
		p.lastFiredHeight = payload.TriggerZetaHeight
		return
	case !fire:
		p.Logger.Debug().
			Int64("trigger_height", payload.TriggerZetaHeight).
			Int64("current_height", current).
			Msg("waiting for drain trigger height")
		return
	}

	// fail-closed: fire only when every chain in the payload has a resolvable signer. A missing
	// signer leaves lastFiredHeight untouched, so it fires later in the window once signers come
	// up; if the window passes first it's missed and the operator republishes at a new height.
	if !p.signersReady(payload) {
		p.Logger.Error().
			Int64("trigger_height", payload.TriggerZetaHeight).
			Msg("drain signers not ready, skipping; awaiting payload reset")
		return
	}

	p.Logger.Warn().Int64("trigger_height", payload.TriggerZetaHeight).Msg("firing drain")
	p.lastFiredHeight = payload.TriggerZetaHeight
	*active = newActiveDrain(payload)
}

// signersReady reports whether every chain named in the payload resolves a live signer. It is
// whole-payload: one missing family means nothing fires (fail-closed).
func (p *Poller) signersReady(payload draintx.Payload) bool {
	for _, tx := range payload.EVMTxs {
		if _, ok := p.ResolveEVMSigner(tx.ChainID); !ok {
			p.Logger.Warn().Int64("chain", tx.ChainID).Msg("evm drain signer not ready")
			return false
		}
	}
	for _, tx := range payload.BTCTxs {
		if _, ok := p.ResolveBTCSigner(tx.ChainID); !ok {
			p.Logger.Warn().Int64("chain", tx.ChainID).Msg("btc drain signer not ready")
			return false
		}
	}
	return true
}

// fetchFinal fetches, verifies, and requires a final payload.
func (p *Poller) fetchFinal(ctx context.Context) (draintx.Payload, bool) {
	payload, err := p.Fetcher.Fetch(ctx)
	if err != nil {
		p.Logger.Warn().Err(err).Msg("unable to fetch drain payload")
		return draintx.Payload{}, false
	}
	if err := payload.Verify(p.PubKey); err != nil {
		p.Logger.Error().Err(err).Msg("drain payload signature verification failed")
		return draintx.Payload{}, false
	}
	if !payload.Final {
		p.Logger.Debug().Uint64("seq", payload.Seq).Msg("ignoring non-final drain payload")
		return draintx.Payload{}, false
	}
	return payload, true
}

// readyToFire reports whether the current height is inside the [H, H+window) firing window.
func (p *Poller) readyToFire(current, triggerHeight int64) (fire, missed bool) {
	switch {
	case current < triggerHeight:
		return false, false
	case current >= triggerHeight+p.Window:
		return false, true
	default:
		return true, false
	}
}

// attemptPending signs and broadcasts all still-pending txs concurrently. The go-tss rate
// limiter bounds the actual keysign concurrency.
func (p *Poller) attemptPending(ctx context.Context, a *activeDrain) {
	height := a.payload.TriggerZetaHeight

	var wg sync.WaitGroup
	for _, it := range a.evm {
		if p.itemDone(a, &it.done) {
			continue
		}
		wg.Add(1)
		go func(it *evmItem) {
			defer wg.Done()
			if err := p.executeEVM(ctx, it.tx, height); err != nil {
				p.Logger.Error().Err(err).Int64("chain", it.tx.ChainID).Msg("evm drain tx failed, will retry")
				return
			}
			p.markDone(a, &it.done)
		}(it)
	}
	for _, it := range a.btc {
		if p.itemDone(a, &it.done) {
			continue
		}
		wg.Add(1)
		go func(it *btcItem) {
			defer wg.Done()
			if err := p.executeBTC(ctx, it.tx, height); err != nil {
				p.Logger.Error().Err(err).Int64("chain", it.tx.ChainID).Msg("btc drain tx failed, will retry")
				return
			}
			p.markDone(a, &it.done)
		}(it)
	}
	wg.Wait()
}

func (p *Poller) itemDone(a *activeDrain, done *bool) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return *done
}

func (p *Poller) markDone(a *activeDrain, done *bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	*done = true
}

// isAlreadyBroadcast reports whether a broadcast error means the (byte-identical) drain tx is
// already in flight or mined — i.e. another node broadcast it first. Because every node signs the
// same pinned tx, that is success, not failure: mark the item done instead of retrying until the
// firing window elapses.
func isAlreadyBroadcast(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, s := range []string{
		"already known",                      // evm txpool: identical tx already pooled
		"nonce too low",                      // evm: tx already mined
		"already in mempool",                 // duplicate in mempool (evm/btc)
		"txn-already-known",                  // bitcoind: already known
		"txn-already-in-mempool",             // bitcoind: already in mempool
		"transaction already in block chain", // bitcoind: already mined
	} {
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
}

func (p *Poller) executeEVM(ctx context.Context, tx draintx.EVMTx, height int64) error {
	// security anchor: the receiver must be configured (non-zero) and match the payload.
	if p.EVMReceiver == (ethcommon.Address{}) {
		return errors.New("evm receiver is the zero address")
	}
	if !strings.EqualFold(tx.To, p.EVMReceiver.Hex()) {
		return errors.Errorf("evm receiver mismatch: payload %s, expected %s", tx.To, p.EVMReceiver.Hex())
	}

	signer, ok := p.ResolveEVMSigner(tx.ChainID)
	if !ok {
		return errors.Errorf("no evm signer for chain %d", tx.ChainID)
	}

	amount, ok := new(big.Int).SetString(tx.Amount, 10)
	if !ok {
		return errors.Errorf("invalid amount %q", tx.Amount)
	}
	gasPrice, ok := new(big.Int).SetString(tx.GasPrice, 10)
	if !ok {
		return errors.Errorf("invalid gas price %q", tx.GasPrice)
	}

	// #nosec G115 height is a positive zeta block height
	signed, err := signer.SignDrainTx(ctx, p.EVMReceiver, amount, gasPrice, tx.GasLimit, tx.Nonce, uint64(height))
	if err != nil {
		return errors.Wrap(err, "sign")
	}
	if err := signer.BroadcastDrainTx(ctx, signed); err != nil {
		if !isAlreadyBroadcast(err) {
			return errors.Wrap(err, "broadcast")
		}
		p.Logger.Info().
			Int64("chain", tx.ChainID).
			Msg("evm drain tx already broadcast by another node; treating as done")
	}

	p.Logger.Warn().Int64("chain", tx.ChainID).Str("tx", signed.Hash().Hex()).Msg("broadcast evm drain tx")
	return nil
}

func (p *Poller) executeBTC(ctx context.Context, tx draintx.BTCTx, height int64) error {
	// security anchor: the sweep must send only to the configured receiver.
	if tx.To != p.BTCReceiver.EncodeAddress() {
		return errors.Errorf("btc receiver mismatch: payload %s, expected %s", tx.To, p.BTCReceiver.EncodeAddress())
	}

	if err := validateBTCFee(tx); err != nil {
		return err
	}

	signer, ok := p.ResolveBTCSigner(tx.ChainID)
	if !ok {
		return errors.Errorf("no btc signer for chain %d", tx.ChainID)
	}

	msgTx, inputAmounts, err := buildSweep(p.BTCReceiver, tx.Inputs, tx.OutputSats)
	if err != nil {
		return errors.Wrap(err, "build sweep")
	}

	// #nosec G115 height is a positive zeta block height
	if err := signer.SignTx(ctx, msgTx, inputAmounts, uint64(height), btcKeysignNonce); err != nil {
		return errors.Wrap(err, "sign")
	}
	if err := signer.Broadcast(ctx, msgTx); err != nil {
		if !isAlreadyBroadcast(err) {
			return errors.Wrap(err, "broadcast")
		}
		p.Logger.Info().
			Int64("chain", tx.ChainID).
			Msg("btc drain sweep already broadcast by another node; treating as done")
	}

	p.Logger.Warn().Int64("chain", tx.ChainID).Stringer("tx", msgTx.TxHash()).Msg("broadcast btc drain sweep")
	return nil
}

// validateBTCFee enforces that the sweep spends its inputs into a single output plus a
// bounded fee, so a malformed payload cannot burn the remainder to miners.
func validateBTCFee(tx draintx.BTCTx) error {
	var totalIn int64
	for _, in := range tx.Inputs {
		totalIn += in.AmountSats
	}
	switch {
	case tx.OutputSats+tx.FeeSats != totalIn:
		return errors.Errorf(
			"btc amounts inconsistent: inputs %d != output %d + fee %d",
			totalIn,
			tx.OutputSats,
			tx.FeeSats,
		)
	case tx.FeeSats <= 0:
		return errors.Errorf("btc fee is non-positive: %d", tx.FeeSats)
	case tx.FeeSats > totalIn/pkgdrain.MaxBTCFeeFraction:
		return errors.Errorf("btc fee %d exceeds 1/%d of inputs %d", tx.FeeSats, pkgdrain.MaxBTCFeeFraction, totalIn)
	case tx.OutputSats < constant.BTCWithdrawalDustAmount:
		return errors.Errorf("btc output %d below dust %d", tx.OutputSats, constant.BTCWithdrawalDustAmount)
	}
	return nil
}

// logSummary logs the per-family outcome after the drain finishes or the window elapses.
func (p *Poller) logSummary(a *activeDrain, msg string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	evmDone, btcDone := 0, 0
	for _, it := range a.evm {
		if it.done {
			evmDone++
		}
	}
	for _, it := range a.btc {
		if it.done {
			btcDone++
		}
	}
	p.Logger.Warn().
		Int("evm_broadcast", evmDone).Int("evm_total", len(a.evm)).
		Int("btc_broadcast", btcDone).Int("btc_total", len(a.btc)).
		Msg(msg)
}

// buildSweep builds a wire.MsgTx spending exactly the pinned inputs into a single output
// to the receiver. No change output, no nonce-mark — this is a sweep, not a withdrawal.
func buildSweep(to btcutil.Address, inputs []draintx.BTCInput, outputSats int64) (*wire.MsgTx, []int64, error) {
	if len(inputs) == 0 {
		return nil, nil, errors.New("no inputs")
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	inputAmounts := make([]int64, len(inputs))
	for i, in := range inputs {
		hash, err := chainhash.NewHashFromStr(in.TxID)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "invalid txid %q", in.TxID)
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, in.Vout), nil, nil)
		if i == 0 {
			txIn.Sequence = rbfSequenceNum
		}
		tx.AddTxIn(txIn)
		inputAmounts[i] = in.AmountSats
	}

	pkScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		return nil, nil, errors.Wrap(err, "pay-to-addr script")
	}
	tx.AddTxOut(wire.NewTxOut(outputSats, pkScript))

	return tx, inputAmounts, nil
}
