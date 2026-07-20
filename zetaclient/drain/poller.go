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
	"github.com/zeta-chain/node/pkg/draintx"
)

// btcKeysignNonce is a fixed keysign nonce for BTC sweeps. The go-tss ceremony matches on
// digests + height, not nonce (nonce is cosmetic), so a constant is safe and deterministic.
const btcKeysignNonce = 0

// rbfSequenceNum opts the first input into full-RBF, mirroring the production BTC signer.
const rbfSequenceNum uint32 = 1

// EVMSigner is the subset of the EVM signer the poller drives.
type EVMSigner interface {
	Chain() chains.Chain
	SignDrainTx(ctx context.Context, to ethcommon.Address, amount, gasPrice *big.Int, gasLimit, nonce, height uint64) (*eth.Transaction, error)
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
	Fetcher      Fetcher
	Height       HeightProvider
	PubKey       []byte // baked-in operator public key
	EVMReceiver  ethcommon.Address
	BTCReceiver  btcutil.Address
	EVMSigners   map[int64]EVMSigner
	BTCSigners   map[int64]BTCSigner
	Window       int64 // blocks after H during which a node may still fire ("ignored if missed")
	PollInterval time.Duration
	Logger       zerolog.Logger
}

// Poller polls the drain endpoint and fires the drain once at the trigger height.
type Poller struct {
	Config
}

// New creates a Poller.
func New(cfg Config) *Poller {
	return &Poller{Config: cfg}
}

// Run polls until the drain fires, is missed, or the context is cancelled.
func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.PollInterval)
	defer ticker.Stop()

	p.Logger.Info().Msg("drain poller started")
	for {
		select {
		case <-ctx.Done():
			p.Logger.Info().Msg("drain poller stopped")
			return
		case <-ticker.C:
			if p.tick(ctx) {
				return
			}
		}
	}
}

// tick runs one poll iteration. It returns true when the poller is done (fired or missed).
func (p *Poller) tick(ctx context.Context) (done bool) {
	payload, err := p.Fetcher.Fetch(ctx)
	if err != nil {
		p.Logger.Warn().Err(err).Msg("unable to fetch drain payload")
		return false
	}

	if err := payload.Verify(p.PubKey); err != nil {
		p.Logger.Error().Err(err).Msg("drain payload signature verification failed")
		return false
	}

	if !payload.Final {
		p.Logger.Debug().Uint64("seq", payload.Seq).Msg("ignoring non-final drain payload")
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

	p.execute(ctx, payload)
	return true
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

// execute signs and broadcasts every tx in the payload.
func (p *Poller) execute(ctx context.Context, payload draintx.Payload) {
	height := payload.TriggerZetaHeight
	p.Logger.Warn().Int64("trigger_height", height).Msg("firing drain")

	for _, tx := range payload.EVMTxs {
		if err := p.executeEVM(ctx, tx, height); err != nil {
			p.Logger.Error().Err(err).Int64("chain", tx.ChainID).Msg("evm drain tx failed")
		}
	}
	for _, tx := range payload.BTCTxs {
		if err := p.executeBTC(ctx, tx, height); err != nil {
			p.Logger.Error().Err(err).Int64("chain", tx.ChainID).Msg("btc drain tx failed")
		}
	}
}

func (p *Poller) executeEVM(ctx context.Context, tx draintx.EVMTx, height int64) error {
	signer, ok := p.EVMSigners[tx.ChainID]
	if !ok {
		return errors.Errorf("no evm signer for chain %d", tx.ChainID)
	}

	// security anchor: the tx may only ever send to the compiled-in receiver.
	if !strings.EqualFold(tx.To, p.EVMReceiver.Hex()) {
		return errors.Errorf("evm receiver mismatch: payload %s, expected %s", tx.To, p.EVMReceiver.Hex())
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
	signer, ok := p.BTCSigners[tx.ChainID]
	if !ok {
		return errors.Errorf("no btc signer for chain %d", tx.ChainID)
	}

	// security anchor: the sweep may only ever send to the compiled-in receiver.
	if tx.To != p.BTCReceiver.EncodeAddress() {
		return errors.Errorf("btc receiver mismatch: payload %s, expected %s", tx.To, p.BTCReceiver.EncodeAddress())
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
