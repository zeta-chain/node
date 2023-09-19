package common

import (
	"errors"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zeta-chain/zetacore/common/ethereum"
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
func NewEthereumProof(proof ethereum.Proof) *Proof {
	return &Proof{
		Proof: &Proof_EthereumProof{
			EthereumProof: &proof,
		},
	}
}

// Verify verifies the proof against the header
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
	default:
		return nil, errors.New("unrecognized proof type")
	}
}
