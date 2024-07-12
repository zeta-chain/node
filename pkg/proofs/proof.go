package proofs

import (
	"bytes"
	"errors"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/zeta-chain/zetacore/pkg/proofs/bitcoin"
	"github.com/zeta-chain/zetacore/pkg/proofs/ethereum"
)

// ErrInvalidProof is a error type for invalid proofs embedding the underlying error
type ErrInvalidProof struct {
	Err error
}

func NewErrInvalidProof(err error) ErrInvalidProof {
	return ErrInvalidProof{
		Err: err,
	}
}

func (e ErrInvalidProof) Error() string {
	return e.Err.Error()
}

// IsErrorInvalidProof returns true if the error is an ErrInvalidProof
func IsErrorInvalidProof(err error) bool {
	return errors.As(err, &ErrInvalidProof{})
}

// NewEthereumProof returns a new Proof containing an Ethereum proof
func NewEthereumProof(proof *ethereum.Proof) *Proof {
	return &Proof{
		Proof: &Proof_EthereumProof{
			EthereumProof: proof,
		},
	}
}

// NewBitcoinProof returns a new Proof containing a Bitcoin proof
func NewBitcoinProof(txBytes []byte, path []byte, index uint) *Proof {
	return &Proof{
		Proof: &Proof_BitcoinProof{
			BitcoinProof: &bitcoin.Proof{
				TxBytes: txBytes,
				Path:    path,
				// #nosec G115 always in range
				Index: uint32(index),
			},
		},
	}
}

// Verify verifies the proof against the header
// Returns the verified tx in bytes if the verification is successful
func (p Proof) Verify(headerData HeaderData, txIndex int) ([]byte, error) {
	switch proof := p.Proof.(type) {
	case *Proof_EthereumProof:
		ethHeaderBytes := headerData.GetEthereumHeader()
		if ethHeaderBytes == nil {
			return nil, errors.New("can't verify ethereum proof against non-ethereum header")
		}
		var ethHeader ethtypes.Header
		err := rlp.DecodeBytes(ethHeaderBytes, &ethHeader)
		if err != nil {
			return nil, err
		}
		val, err := proof.EthereumProof.Verify(ethHeader.TxHash, txIndex)
		if err != nil {
			return nil, NewErrInvalidProof(err)
		}
		return val, nil
	case *Proof_BitcoinProof:
		btcHeaderBytes := headerData.GetBitcoinHeader()
		if len(btcHeaderBytes) != bitcoin.BitcoinBlockHeaderLen {
			return nil, errors.New("can't verify bitcoin proof against non-bitcoin header")
		}
		var btcHeader wire.BlockHeader
		if err := btcHeader.Deserialize(bytes.NewReader(btcHeaderBytes)); err != nil {
			return nil, err
		}
		tx, err := btcutil.NewTxFromBytes(proof.BitcoinProof.TxBytes)
		if err != nil {
			return nil, err
		}
		pass := bitcoin.Prove(*tx.Hash(), btcHeader.MerkleRoot, proof.BitcoinProof.Path, uint(proof.BitcoinProof.Index))
		if !pass {
			return nil, NewErrInvalidProof(errors.New("invalid bitcoin proof"))
		}
		return proof.BitcoinProof.TxBytes, nil
	default:
		return nil, errors.New("unrecognized proof type")
	}
}
