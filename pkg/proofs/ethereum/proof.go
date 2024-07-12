// This file was adapted from go-ethereum. Here's the go-ethereum license reproduced:

// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package trie implements Merkle Patricia Tries.

package ethereum

import (
	"bytes"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func NewProof() *Proof {
	return &Proof{
		Keys:   make([][]byte, 0),
		Values: make([][]byte, 0),
	}
}

func (m *Proof) Put(key []byte, value []byte) error {
	for i := 0; i < len(m.Keys); i++ {
		if bytes.Equal(m.Keys[i], key) {
			m.Values[i] = value
			return nil
		}
	}
	m.Keys = append(m.Keys, key)
	m.Values = append(m.Values, value)

	return nil
}

func (m *Proof) Delete(key []byte) error {
	found := false
	index := -1
	for i := 0; i < len(m.Keys); i++ {
		if bytes.Equal(m.Keys[i], key) {
			found = true
			index = i
			break
		}
	}
	if !found {
		return errors.New("key not found")
	}
	copy(m.Keys[index:len(m.Keys)-1], m.Keys[index+1:])
	copy(m.Values[index:len(m.Values)-1], m.Values[index+1:])
	m.Keys = m.Keys[:len(m.Keys)-1]
	m.Values = m.Values[:len(m.Values)-1]

	return nil
}

func (m *Proof) Has(key []byte) (bool, error) {
	for i := 0; i < len(m.Keys); i++ {
		if bytes.Equal(m.Keys[i], key) {
			return true, nil
		}
	}
	return false, nil
}

func (m *Proof) Get(key []byte) ([]byte, error) {
	found := false
	index := -1
	for i := 0; i < len(m.Keys); i++ {
		if bytes.Equal(m.Keys[i], key) {
			found = true
			index = i
			break
		}
	}
	if !found {
		return nil, errors.New("key not found")
	}

	return m.Values[index], nil
}

// Verify verifies the proof against the given root hash and key.
// Typically, the rootHash is from a trusted source (e.g. a trusted block header),
// and the key is the index of the transaction in the block.
func (m *Proof) Verify(rootHash common.Hash, key int) ([]byte, error) {
	if key < 0 {
		return nil, errors.New("key not found")
	}
	var indexBuf []byte
	// #nosec G115 range is valid
	indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(key))
	return trie.VerifyProof(rootHash, indexBuf, m)
}

type Trie struct {
	*trie.Trie
}

var encodeBufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func encodeForDerive(list types.DerivableList, i int, buf *bytes.Buffer) []byte {
	buf.Reset()
	list.EncodeIndex(i, buf)
	// It's really unfortunate that we need to do perform this copy.
	// StackTrie holds onto the values until Hash is called, so the values
	// written to it must not alias.
	return common.CopyBytes(buf.Bytes())
}

func (t *Trie) GenerateProof(txIndex int) (*Proof, error) {
	if txIndex < 0 {
		return nil, errors.New("transaction index out of range")
	}
	var indexBuf []byte
	// #nosec G115 checked as non-negative
	indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(txIndex))
	proof := NewProof()
	err := t.Prove(indexBuf, 0, proof)
	if err != nil {
		return nil, err
	}
	return proof, nil
}

// NewTrie builds a trie from a DerivableList. The DerivableList must be types.Transactions
// or types.Receipts.
func NewTrie(list types.DerivableList) Trie {
	hasher := new(trie.Trie)
	hasher.Reset()

	valueBuf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(valueBuf)

	// StackTrie requires values to be inserted in increasing hash order, which is not the
	// order that `list` provides hashes in. This insertion sequence ensures that the
	// order is correct.
	var indexBuf []byte
	for i := 1; i < list.Len() && i <= 0x7f; i++ {
		// #nosec G115 iterator
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		value := encodeForDerive(list, i, valueBuf)
		hasher.Update(indexBuf, value)
	}
	if list.Len() > 0 {
		indexBuf = rlp.AppendUint64(indexBuf[:0], 0)
		value := encodeForDerive(list, 0, valueBuf)
		hasher.Update(indexBuf, value)
	}
	for i := 0x80; i < list.Len(); i++ {
		// #nosec G115 iterator
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		value := encodeForDerive(list, i, valueBuf)
		hasher.Update(indexBuf, value)
	}
	return Trie{hasher}
}
