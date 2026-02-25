// Copyright 2020 Indefinite Integral Incorporated

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//   http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bitcoin

// This file was adapted from Summa bitcoin-spv. Here are some modifications:
// - define 'Hash256Digest' as alias for 'chainhash.Hash'
// - keep only Prove() and dependent functions

import (
	"bytes"
	"crypto/sha256"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Hash256Digest = chainhash.Hash

// Prove checks the validity of a merkle proof
func Prove(txid Hash256Digest, merkleRoot Hash256Digest, intermediateNodes []byte, index uint) bool {
	// Shortcut the empty-block case
	if bytes.Equal(txid[:], merkleRoot[:]) && index == 0 && len(intermediateNodes) == 0 {
		return true
	}

	proof := make([]byte, 0, len(txid)+len(intermediateNodes)+len(merkleRoot))
	proof = append(proof, txid[:]...)
	proof = append(proof, intermediateNodes...)
	proof = append(proof, merkleRoot[:]...)

	return VerifyHash256Merkle(proof, index)
}

// Hash256 implements bitcoin's hash256 (double sha2)
func Hash256(in []byte) Hash256Digest {
	first := sha256.Sum256(in)
	second := sha256.Sum256(first[:])
	return Hash256Digest(second)
}

// Hash256MerkleStep concatenates and hashes two inputs for merkle proving
func Hash256MerkleStep(a []byte, b []byte) Hash256Digest {
	c := make([]byte, 0, len(a)+len(b))
	c = append(c, a...)
	c = append(c, b...)
	return Hash256(c)
}

// VerifyHash256Merkle checks a merkle inclusion proof's validity.
// Note that `index` is not a reliable indicator of location within a block.
func VerifyHash256Merkle(proof []byte, index uint) bool {
	var current Hash256Digest
	idx := index
	proofLength := len(proof)

	if proofLength%32 != 0 {
		return false
	}

	if proofLength == 32 {
		return true
	}

	if proofLength == 64 {
		return false
	}

	root := proof[proofLength-32:]

	cur := proof[:32:32]
	copy(current[:], cur)

	numSteps := (proofLength / 32) - 1

	for i := 1; i < numSteps; i++ {
		start := i * 32
		end := i*32 + 32
		next := proof[start:end:end]
		if idx%2 == 1 {
			current = Hash256MerkleStep(next, current[:])
		} else {
			current = Hash256MerkleStep(current[:], next)
		}
		idx >>= 1
	}

	return bytes.Equal(current[:], root)
}
