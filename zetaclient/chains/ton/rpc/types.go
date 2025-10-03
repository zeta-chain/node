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
	ID      ton.AccountID
	Status  tlb.AccountStatus
	Balance uint64

	Code *boc.Cell
	Data *boc.Cell

	LastTxHash tlb.Bits256
	LastTxLT   uint64
}

func (acc *Account) UnmarshalJSON(data []byte) error {
	items := gjson.GetManyBytes(
		data,
		"state",
		"balance",
		"code",
		"data",
		"frozen_hash",
		"last_transaction_id.lt",
		"last_transaction_id.hash",
	)

	var (
		state      = items[0].String()
		balance    = items[1].Uint()
		codeRaw    = items[2].String()
		dataRaw    = items[3].String()
		frozenHash = items[4].String()
		lt         = items[5].Uint()
		hashRaw    = items[6].String()
		err        error
	)

	acc.Balance = balance

	if codeRaw != "" {
		acc.Code, err = cellFromBase64(codeRaw)
		if err != nil {
			return errors.Wrapf(err, "code")
		}
	}

	if dataRaw != "" {
		acc.Data, err = cellFromBase64(dataRaw)
		if err != nil {
			return errors.Wrapf(err, "data")
		}
	}

	switch {
	case frozenHash != "":
		acc.Status = tlb.AccountFrozen
	case strings.Contains(state, "uninit"):
		acc.Status = tlb.AccountUninit
	case state == "active":
		acc.Status = tlb.AccountActive
	default:
		return errors.New("unable to parse account status")
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hashRaw)
	if err != nil {
		return errors.Wrapf(err, "unable to decode last tx hash from %q", hashRaw)
	}

	copy(acc.LastTxHash[:], hashBytes)
	acc.LastTxLT = lt

	return nil
}

// ToShardAccount partially converts Account to tongo's tlb.ShardAccount
func (acc *Account) ToShardAccount() tlb.ShardAccount {
	if acc.Status == tlb.AccountNone {
		return tlb.ShardAccount{
			Account: tlb.Account{SumType: "AccountNone"},
		}
	}

	return tlb.ShardAccount{
		Account: tlb.Account{
			SumType: "Account",
			Account: tlb.ExistedAccount{
				Addr:        acc.ID.ToMsgAddress(),
				StorageStat: tlb.StorageInfo{},
				Storage: tlb.AccountStorage{
					State:       tlbAccountState(acc),
					LastTransLt: acc.LastTxLT,
					Balance: tlb.CurrencyCollection{
						Grams: tlb.Grams(acc.Balance),
					},
				},
			},
		},
		LastTransHash: acc.LastTxHash,
		LastTransLt:   acc.LastTxLT,
	}
}

// takes base64 encoded BOC and decodes it into v
func unmarshalFromBase64(b64 string, v any) error {
	cell, err := cellFromBase64(b64)
	if err != nil {
		return err
	}

	return tlb.Unmarshal(cell, v)
}

func cellFromBase64(b64 string) (*boc.Cell, error) {
	if b64 == "" {
		return nil, errors.New("empty boc")
	}

	cells, err := boc.DeserializeBocBase64(b64)
	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "unable to deserialize boc from %q", b64)
	case len(cells) == 0:
		return nil, errors.Errorf("expected at least one cell, got 0 (raw: %q)", b64)
	}

	return cells[0], nil
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

func tlbAccountState(a *Account) tlb.AccountState {
	switch a.Status {
	case tlb.AccountActive:
		return tlb.AccountState{
			SumType: "AccountActive",
			AccountActive: struct {
				StateInit tlb.StateInit
			}{
				StateInit: tlb.StateInit{
					Code: wrapCell(a.Code),
					Data: wrapCell(a.Data),
				},
			},
		}
	case tlb.AccountFrozen:
		return tlb.AccountState{SumType: "AccountFrozen"}
	case tlb.AccountUninit:
		return tlb.AccountState{SumType: "AccountUninit"}
	default:
		// should not happen
		return tlb.AccountState{}
	}
}

func wrapCell(v *boc.Cell) tlb.Maybe[tlb.Ref[boc.Cell]] {
	if v == nil {
		return tlb.Maybe[tlb.Ref[boc.Cell]]{}
	}

	return tlb.Maybe[tlb.Ref[boc.Cell]]{
		Exists: true,
		Value:  tlb.Ref[boc.Cell]{Value: *v},
	}
}
