package klaytn

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type KlaytnClient struct {
	c *rpc.Client
}

type RPCHeader struct {
	Hash   *common.Hash `json:"hash"`
	Number *hexutil.Big `json:"number"`
	Time   *hexutil.Big `json:"timestamp"`
}

type RPCTransaction struct {
	From  *common.Address `json:"from"`
	Input hexutil.Bytes   `json:"input"`
	To    *common.Address `json:"to"`
	Hash  common.Hash     `json:"hash"`
	Value *hexutil.Big    `json:"value"`
}

type RPCBlock struct {
	Hash         *common.Hash     `json:"hash"`
	Transactions []RPCTransaction `json:"transactions"`
}

func Dial(url string) (*KlaytnClient, error) {
	c, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &KlaytnClient{c}, nil
}

func (ec *KlaytnClient) BlockByNumber(ctx context.Context, number *big.Int) (*RPCBlock, error) {
	return ec.getBlock(ctx, "klay_getBlockByNumber", toBlockNumArg(number), true)
}

func (ec *KlaytnClient) getBlock(ctx context.Context, method string, args ...interface{}) (*RPCBlock, error) {
	var raw json.RawMessage
	err := ec.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, errors.New("not found")
	}

	var block RPCBlock
	if err := json.Unmarshal(raw, &block); err != nil {
		return nil, err
	}

	return &block, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}
