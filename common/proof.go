package common

import (
	"errors"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zeta-chain/zetacore/common/ethereum"
)

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
		return proof.EthereumProof.Verify(ethHeader.TxHash, txIndex)
	default:
		return nil, errors.New("unrecognized proof type")
	}
}
