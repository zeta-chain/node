// Package drain builds the signed drain payload from resolved per-chain inputs.
// The logic is network-agnostic: callers (the zetatool subcommand and the e2e test)
// fetch balances, gas, nonces and UTXOs however they like, then hand them here so the
// partitioning, fee math and signing live in one deterministic place.
package drain

import (
	"crypto/ecdsa"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/draintx"
	"github.com/zeta-chain/node/pkg/migration"
)

// MaxInputsPerTx mirrors signer.MaxNoOfInputsPerTx: a BTC outbound may spend at most
// 20 inputs, so a larger UTXO set is swept with multiple disjoint transactions.
const MaxInputsPerTx = 20

// EVMInput is the resolved state needed to build a single EVM drain tx.
type EVMInput struct {
	ChainID        int64
	To             string // hardcoded safe receiver
	Balance        sdkmath.Uint
	MedianGasPrice sdkmath.Uint
	Nonce          uint64
}

// UTXO is a single unspent output pinned into a BTC sweep.
type UTXO struct {
	TxID       string
	Vout       uint32
	AmountSats int64
}

// BTCInput is the resolved state needed to build the BTC drain sweeps.
type BTCInput struct {
	ChainID int64
	To      btcutil.Address // hardcoded safe receiver
	FeeRate int64           // sat/vB
	UTXOs   []UTXO
}

// GenerateEVMTx builds a fully-resolved EVM drain tx from a resolved input.
func GenerateEVMTx(in EVMInput) (draintx.EVMTx, error) {
	amount, gasPrice, gasLimit, err := migration.ComputeEVMMigration(in.Balance, in.MedianGasPrice)
	if err != nil {
		return draintx.EVMTx{}, err
	}
	return draintx.EVMTx{
		ChainID:  in.ChainID,
		To:       in.To,
		Nonce:    in.Nonce,
		Amount:   amount.String(),
		GasPrice: gasPrice.String(),
		GasLimit: gasLimit,
	}, nil
}

// GenerateBTCTxs partitions the UTXO set into disjoint groups of at most MaxInputsPerTx
// and builds one independent sweep per group. Each sweep spends its pinned inputs into a
// single output to the receiver, minus the fee.
func GenerateBTCTxs(in BTCInput) ([]draintx.BTCTx, error) {
	if len(in.UTXOs) == 0 {
		return nil, nil
	}

	to := in.To.EncodeAddress()
	txs := make([]draintx.BTCTx, 0, (len(in.UTXOs)+MaxInputsPerTx-1)/MaxInputsPerTx)

	for start := 0; start < len(in.UTXOs); start += MaxInputsPerTx {
		end := start + MaxInputsPerTx
		if end > len(in.UTXOs) {
			end = len(in.UTXOs)
		}
		group := in.UTXOs[start:end]

		var totalSats int64
		inputs := make([]draintx.BTCInput, len(group))
		for i, u := range group {
			totalSats += u.AmountSats
			inputs[i] = draintx.BTCInput{TxID: u.TxID, Vout: u.Vout, AmountSats: u.AmountSats}
		}

		outputSats, feeSats, err := migration.ComputeBTCMigration(totalSats, in.FeeRate)
		if err != nil {
			return nil, err
		}

		// a group whose swept output would be below dust is unspendable; reject rather than
		// build a tx that the network will not relay.
		if outputSats < constant.BTCWithdrawalDustAmount {
			return nil, errors.Errorf(
				"btc sweep group output %d below dust %d", outputSats, constant.BTCWithdrawalDustAmount,
			)
		}

		txs = append(txs, draintx.BTCTx{
			ChainID:    in.ChainID,
			To:         to,
			OutputSats: outputSats,
			FeeSats:    feeSats,
			Inputs:     inputs,
		})
	}

	return txs, nil
}

// BuildPayload assembles the txs into a payload and signs it with priv.
func BuildPayload(
	triggerHeight int64,
	seq uint64,
	final bool,
	evmTxs []draintx.EVMTx,
	btcTxs []draintx.BTCTx,
	priv *ecdsa.PrivateKey,
) (draintx.Payload, error) {
	p := draintx.Payload{
		TriggerZetaHeight: triggerHeight,
		Seq:               seq,
		Final:             final,
		EVMTxs:            evmTxs,
		BTCTxs:            btcTxs,
	}
	if err := p.Sign(priv); err != nil {
		return draintx.Payload{}, err
	}
	return p, nil
}
