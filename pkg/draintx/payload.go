// Package draintx defines the off-chain, signed transaction payload used by the
// emergency TSS drain: the generator computes fully-resolved EVM and BTC txs, signs
// the payload with an operator key, and each zetaclient verifies it against a
// baked-in public key before signing. The client makes no decisions — it copies the
// pinned values and signs.
package draintx

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// domainTag is a fixed prefix mixed into every signed digest so a drain signature can never be
// mistaken for a signature over any other message.
const domainTag = "ZETADRAIN"

// Payload is the full, byte-final drain instruction shared by every signer.
type Payload struct {
	// TriggerZetaHeight is the zeta block height at which clients fire. It doubles as
	// the go-tss leader-election blockHeight, so it must be identical across nodes.
	TriggerZetaHeight int64 `json:"trigger_zeta_height"`
	// Seq is a monotonic version for observability only.
	Seq uint64 `json:"seq"`
	// Final gates signing: clients only ever sign a payload marked final.
	Final bool `json:"final"`
	// Network binds the payload to a single zeta network (mainnet/testnet/localnet) so a client
	// armed for one network rejects a payload built for another.
	Network string `json:"network"`
	// EVMTxs holds one fully-resolved native transfer per EVM chain.
	EVMTxs []EVMTx `json:"evm_txs"`
	// BTCTxs holds one sweep per disjoint group of <=20 UTXOs.
	BTCTxs []BTCTx `json:"btc_txs"`
	// Signature is the 0x-prefixed hex secp256k1 signature over the canonical bytes.
	Signature string `json:"signature"`
}

// EVMTx is a fully-resolved native transfer on an EVM chain.
type EVMTx struct {
	ChainID  int64  `json:"chain_id"`
	To       string `json:"to"`
	Nonce    uint64 `json:"nonce"`
	Amount   string `json:"amount"`    // wei, decimal string
	GasPrice string `json:"gas_price"` // wei, decimal string
	GasLimit uint64 `json:"gas_limit"`
}

// BTCTx is a fully-resolved Bitcoin sweep spending the pinned inputs into a single output.
type BTCTx struct {
	ChainID    int64      `json:"chain_id"`
	To         string     `json:"to"`
	OutputSats int64      `json:"output_sats"`
	FeeSats    int64      `json:"fee_sats"`
	Inputs     []BTCInput `json:"inputs"`
}

// BTCInput is a pinned UTXO the client spends verbatim.
type BTCInput struct {
	TxID       string `json:"txid"`
	Vout       uint32 `json:"vout"`
	AmountSats int64  `json:"amount_sats"`
}

// canonicalBytes returns a deterministic byte encoding of the payload excluding the
// signature. Every field is written in a fixed order; strings are length-prefixed so
// no concatenation is ambiguous.
func (p Payload) canonicalBytes() []byte {
	var b bytes.Buffer

	writeString(&b, domainTag)
	writeString(&b, p.Network)
	writeInt64(&b, p.TriggerZetaHeight)
	writeUint64(&b, p.Seq)
	writeBool(&b, p.Final)

	writeUint64(&b, uint64(len(p.EVMTxs)))
	for _, tx := range p.EVMTxs {
		writeInt64(&b, tx.ChainID)
		writeString(&b, tx.To)
		writeUint64(&b, tx.Nonce)
		writeString(&b, tx.Amount)
		writeString(&b, tx.GasPrice)
		writeUint64(&b, tx.GasLimit)
	}

	writeUint64(&b, uint64(len(p.BTCTxs)))
	for _, tx := range p.BTCTxs {
		writeInt64(&b, tx.ChainID)
		writeString(&b, tx.To)
		writeInt64(&b, tx.OutputSats)
		writeInt64(&b, tx.FeeSats)
		writeUint64(&b, uint64(len(tx.Inputs)))
		for _, in := range tx.Inputs {
			writeString(&b, in.TxID)
			writeUint32(&b, in.Vout)
			writeInt64(&b, in.AmountSats)
		}
	}

	return b.Bytes()
}

// digest is the keccak256 hash of the canonical bytes, i.e. what gets signed.
func (p Payload) digest() []byte {
	return ethcrypto.Keccak256(p.canonicalBytes())
}

// Sign signs the payload with priv and stores the signature.
func (p *Payload) Sign(priv *ecdsa.PrivateKey) error {
	sig, err := ethcrypto.Sign(p.digest(), priv)
	if err != nil {
		return errors.Wrap(err, "unable to sign drain payload")
	}
	p.Signature = hexutil.Encode(sig)
	return nil
}

// Verify checks the payload signature against the expected public key, which may be
// in 33-byte compressed or 65-byte uncompressed form.
func (p Payload) Verify(pubBytes []byte) error {
	sig, err := hexutil.Decode(p.Signature)
	if err != nil {
		return errors.Wrap(err, "invalid signature encoding")
	}

	expected, err := parsePubKey(pubBytes)
	if err != nil {
		return errors.Wrap(err, "invalid expected public key")
	}

	recovered, err := ethcrypto.SigToPub(p.digest(), sig)
	if err != nil {
		return errors.Wrap(err, "unable to recover public key")
	}

	if !recovered.Equal(expected) {
		return errors.New("drain payload signature does not match expected public key")
	}
	return nil
}

func parsePubKey(b []byte) (*ecdsa.PublicKey, error) {
	if len(b) == 33 {
		return ethcrypto.DecompressPubkey(b)
	}
	return ethcrypto.UnmarshalPubkey(b)
}

func writeString(b *bytes.Buffer, s string) {
	writeUint64(b, uint64(len(s)))
	b.WriteString(s)
}

func writeInt64(b *bytes.Buffer, v int64) {
	// #nosec G115 two's-complement round-trip is intentional
	writeUint64(b, uint64(v))
}

func writeUint64(b *bytes.Buffer, v uint64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], v)
	b.Write(buf[:])
}

func writeUint32(b *bytes.Buffer, v uint32) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], v)
	b.Write(buf[:])
}

func writeBool(b *bytes.Buffer, v bool) {
	if v {
		b.WriteByte(1)
		return
	}
	b.WriteByte(0)
}
