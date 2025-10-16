// Package zetacore provides the client to interact with zetacore node via GRPC.
package zetacore

import (
	"slices"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// BroadcastMode is the broadcast mode for the client
type BroadcastMode string

const (
	Single              BroadcastMode = "single"
	Multi               BroadcastMode = "multi"
	maxSelectedMessages               = 2
)

// PoolMsg represents a message in the message pool
type PoolMsg struct {
	Msg      sdktypes.Msg
	Type     string
	Digest   string
	GasLimit uint64
	GasPrice sdkmath.LegacyDec
	AddedAt  time.Time
}

// MessagePool is a pool of cached messages for broadcasting
type MessagePool struct {
	mu sync.Mutex

	typedMessageMap map[string]map[string]PoolMsg
}

func NewMessagePool() *MessagePool {
	return &MessagePool{
		mu:              sync.Mutex{},
		typedMessageMap: make(map[string]map[string]PoolMsg),
	}
}

// AddMessage adds a message to the message pool
func (mp *MessagePool) AddMessage(msg sdktypes.Msg, msgType, msgDigest string, gasPrice sdkmath.LegacyDec, gasLimit uint64) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	zetacoreMsg := PoolMsg{
		Msg:      msg,
		Type:     msgType,
		Digest:   msgDigest,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		AddedAt:  time.Now(),
	}

	if _, found := mp.typedMessageMap[msgType]; !found {
		mp.typedMessageMap[msgType] = make(map[string]PoolMsg)
	}

	poolMsg, found := mp.typedMessageMap[msgType][msgDigest]
	if !found {
		mp.typedMessageMap[msgType][msgDigest] = zetacoreMsg
	} else {
		poolMsg.GasPrice = gasPrice
		poolMsg.GasLimit = gasLimit
	}
}

// GetMultipleMessages gets multiple messages from the message pool
func (mp *MessagePool) GetMultipleMessages(msgType string, maxCount int) []PoolMsg {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	msgMap, found := mp.typedMessageMap[msgType]
	if !found {
		return nil
	}

	// flatten the messages into a single list
	msgList := make([]PoolMsg, 0, len(msgMap))
	for _, m := range msgMap {
		msgList = append(msgList, m)
	}

	// sort the messages by timestamp by addedAt, ascending
	slices.SortStableFunc(msgList, func(a, b PoolMsg) int {
		return a.AddedAt.Compare(b.AddedAt)
	})

	// select up to maxCount messages oldest to newest
	if len(msgList) < maxCount {
		return msgList
	}

	return msgList[:maxCount]
}
