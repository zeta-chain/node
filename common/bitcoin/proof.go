package bitcoin

import (
	"errors"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

const BitcoinBlockHeaderLen = 80

// Merkle is a wrapper around "github.com/btcsuite/btcd/blockchain" merkle tree.
// Additionally, it provides a method to generate a merkle proof for a given transaction.
type Merkle struct {
	tree []*chainhash.Hash
}

func NewMerkle(txns []*btcutil.Tx) *Merkle {
	return &Merkle{
		tree: blockchain.BuildMerkleTreeStore(txns, false),
	}
}

// BuildMerkleProof builds merkle proof for a given transaction index in block.
func (m *Merkle) BuildMerkleProof(txIndex int) ([]byte, uint, error) {
	if len(m.tree) <= 0 {
		return nil, 0, errors.New("merkle tree is empty")
	}

	// len(m.tree) + 1 must be a power of 2. E.g. 2, 4, 8, 16, 32, 64, 128, 256, ...
	N := len(m.tree) + 1
	if N&(N-1) != 0 {
		return nil, 0, errors.New("merkle tree is not full")
	}

	// Ensure the provided txIndex points to a valid leaf node.
	if txIndex >= N/2 || m.tree[txIndex] == nil {
		return nil, 0, errors.New("transaction index is invalid")
	}
	path := make([]byte, 0)
	var siblingIndexes uint

	// Find intermediate nodes on the path to the root buttom-up.
	nodeIndex := txIndex
	nodesOnLevel := N / 2
	for nodesOnLevel > 1 {
		var flag uint
		var sibling *chainhash.Hash

		if nodeIndex%2 == 1 {
			flag = 1 // left sibling
			sibling = m.tree[nodeIndex-1]
		} else {
			flag = 0 // right sibling
			if m.tree[nodeIndex+1] == nil {
				sibling = m.tree[nodeIndex] // When there is no right sibling, self hash is used.
			} else {
				sibling = m.tree[nodeIndex+1]
			}
		}

		// Append the sibling and flag to the proof.
		path = append(path, sibling[:]...)
		siblingIndexes |= flag << (len(path)/32 - 1)

		// Go up one level to the parent node.
		nodeIndex = N - nodesOnLevel + (nodeIndex%nodesOnLevel)/2
		nodesOnLevel /= 2
	}

	return path, siblingIndexes, nil
}
