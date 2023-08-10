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
	"encoding/base64"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func NewProof() *Proof {
	return &Proof{
		Proof: make(map[string][]byte),
	}
}

func encodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

func (m *Proof) Put(key []byte, value []byte) error {
	m.Proof[encodeKey(key)] = value
	return nil
}

func (m *Proof) Delete(key []byte) error {
	_, exists := m.Proof[encodeKey(key)]
	if !exists {
		return errors.New("key not found")
	}
	delete(m.Proof, encodeKey(key))
	return nil
}

func (m *Proof) Has(key []byte) (bool, error) {
	_, exists := m.Proof[encodeKey(key)]
	return exists, nil
}

func (m *Proof) Get(key []byte) ([]byte, error) {
	value, exists := m.Proof[encodeKey(key)]
	if !exists {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (m *Proof) Verify(rootHash common.Hash, key int) ([]byte, error) {
	var indexBuf []byte
	indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(key))
	return trie.VerifyProof(rootHash, indexBuf, m)
}

type Trie struct {
	trie *trie.Trie
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
	var indexBuf []byte
	indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(txIndex))
	proof := NewProof()
	t.trie.Prove(indexBuf, 0, proof)
	return proof, nil
}

// NewTrie builds a trie from a DerivableList. The DerivableList must be types.Transactions
// or types.Receipts.
func NewTrie(list types.DerivableList, hasher *trie.Trie) Trie {
	hasher.Reset()

	valueBuf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(valueBuf)

	// StackTrie requires values to be inserted in increasing hash order, which is not the
	// order that `list` provides hashes in. This insertion sequence ensures that the
	// order is correct.
	var indexBuf []byte
	for i := 1; i < list.Len() && i <= 0x7f; i++ {
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
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		value := encodeForDerive(list, i, valueBuf)
		hasher.Update(indexBuf, value)
	}
	return Trie{hasher}
}
