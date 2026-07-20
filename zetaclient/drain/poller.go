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
	"github.com/zeta-chain/node/pkg/draintx"
)

const (
	// btcKeysignNonce is a fixed keysign nonce for BTC sweeps. The go-tss ceremony matches
	// on digests + height, not nonce (nonce is cosmetic), so a constant is safe.
	btcKeysignNonce = 0
	// rbfSequenceNum opts the first input into full-RBF, mirroring the production BTC signer.
	rbfSequenceNum uint32 = 1
	// maxBTCFeeFraction bounds the sweep fee to at most 1/N of the total input, so a
	// malformed payload cannot burn the balance to miners.
	maxBTCFeeFraction = 4
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

// Run polls until the drain completes, the window elapses, or the context is cancelled.
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
			if p.step(ctx, &active) {
				return
			}
		}
	}
}

// step runs one poll iteration. It returns true when the poller is done (completed,
// missed, or window elapsed).
func (p *Poller) step(ctx context.Context, active **activeDrain) (done bool) {
	// before the drain is armed for a payload, look for an eligible final payload
	if *active == nil {
		payload, ok := p.fetchFinal(ctx)
		if !ok {
			return false
		}
		current, err := p.Height.GetBlockHeight(ctx)
		if err != nil {
			p.Logger.Warn().Err(err).Msg("unable to get zeta block height")
			return false
		}
		fire, missed := p.readyToFire(current, payload.TriggerZetaHeight)
		switch {
		case missed:
			p.Logger.Error().
				Int64("trigger_height", payload.TriggerZetaHeight).
				Int64("current_height", current).
				Msg("drain trigger height missed, ignoring")
			return true
		case !fire:
			p.Logger.Debug().
				Int64("trigger_height", payload.TriggerZetaHeight).
				Int64("current_height", current).
				Msg("waiting for drain trigger height")
			return false
		}
		p.Logger.Warn().Int64("trigger_height", payload.TriggerZetaHeight).Msg("firing drain")
		*active = newActiveDrain(payload)
	}

	a := *active
	p.attemptPending(ctx, a)
	if a.allDone() {
		p.logSummary(a, "drain complete")
		return true
	}

	current, err := p.Height.GetBlockHeight(ctx)
	if err != nil {
		// can't tell if the window elapsed; keep retrying
		return false
	}
	if current >= a.payload.TriggerZetaHeight+p.Window {
		p.logSummary(a, "drain window elapsed with unfinished txs")
		return true
	}
	return false
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
		return errors.Wrap(err, "broadcast")
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
		return errors.Wrap(err, "broadcast")
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
	case tx.FeeSats > totalIn/maxBTCFeeFraction:
		return errors.Errorf("btc fee %d exceeds 1/%d of inputs %d", tx.FeeSats, maxBTCFeeFraction, totalIn)
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
