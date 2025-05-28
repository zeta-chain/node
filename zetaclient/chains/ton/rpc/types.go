package rpc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type BlockIDExt struct {
	Workchain int    `json:"workchain"`
	Seqno     uint32 `json:"seqno"`
	Shard     string `json:"shard"`
	RootHash  string `json:"root_hash"`
	FileHash  string `json:"file_hash"`
}

type MasterchainInfo struct {
	Last BlockIDExt `json:"last"`
}

type BlockHeader struct {
	ID            BlockIDExt `json:"id"`
	MinRefMcSeqno uint32     `json:"min_ref_mc_seqno"`
	GenUtime      uint32     `json:"gen_utime"`
}
type Account struct {
	ID         ton.AccountID
	Status     tlb.AccountStatus
	Balance    uint64
	LastTxHash tlb.Bits256
	LastTxLT   uint64
}

func (acc *Account) UnmarshalJSON(data []byte) error {
	items := gjson.GetManyBytes(
		data,
		"address.account_address",
		"balance",
		"last_transaction_id.lt",
		"last_transaction_id.hash",
		"account_state.@type",
		"account_state.frozen_hash",
	)

	var (
		addrRaw       = items[0].String()
		balanceRaw    = items[1].String()
		ltRaw         = items[2].String()
		hashRaw       = items[3].String()
		stateRaw      = items[4].String()
		frozenHashRaw = items[5].String()
	)

	id, err := ton.ParseAccountID(addrRaw)
	if err != nil {
		return errors.Wrapf(err, "unable to parse account id from %q", addrRaw)
	}

	acc.ID = id

	if balanceRaw != "-1" {
		acc.Balance = items[1].Uint()
	}

	switch {
	case ltRaw == "0" && strings.Contains(stateRaw, "uninit"):
		acc.Status = tlb.AccountNone
		return nil
	case strings.Contains(stateRaw, "uninit"):
		acc.Status = tlb.AccountUninit
	case stateRaw == "raw.accountState":
		acc.Status = tlb.AccountActive
	case frozenHashRaw != "":
		acc.Status = tlb.AccountFrozen
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hashRaw)
	if err != nil {
		return errors.Wrapf(err, "unable to decode last tx hash from %q", hashRaw)
	}

	copy(acc.LastTxHash[:], hashBytes)
	acc.LastTxLT = items[2].Uint()

	return nil
}

// takes base64 encoded BOC and decodes it into v
func unmarshalFromBase64(b64 string, v any) error {
	cells, err := boc.DeserializeBocBase64(b64)
	switch {
	case err != nil:
		return errors.Wrapf(err, "unable to deserialize boc from %q", b64)
	case len(cells) == 0:
		return errors.Errorf("expected at least one cell, got 0")
	default:
		return tlb.Unmarshal(cells[0], v)
	}
}

type rpcRequest struct {
	Jsonrpc string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
	ID      string         `json:"id"`
}

func newRPCRequest(method string, params map[string]any) rpcRequest {
	if params == nil {
		params = make(map[string]any)
	}

	return rpcRequest{
		Jsonrpc: "2.0",
		ID:      "1",
		Method:  method,
		Params:  params,
	}
}

func (r *rpcRequest) asBody() (io.Reader, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal rpc request")
	}

	return bytes.NewReader(body), nil
}

type rpcResponse struct {
	Success bool            `json:"ok"`
	Result  json.RawMessage `json:"result"`
	Error   string          `json:"error"`
	Code    int             `json:"code"`
}
