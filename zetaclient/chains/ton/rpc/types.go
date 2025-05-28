package rpc

import (
	"encoding/base64"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
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

type AccountState uint8

const (
	AccountStateInvalid AccountState = iota
	AccountStateNotExists
	AccountStateUninit
	AccountStateActive
	AccountStateFrozen
)

type Account struct {
	ID         ton.AccountID
	State      AccountState
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
		acc.State = AccountStateNotExists
		return nil
	case strings.Contains(stateRaw, "uninit"):
		acc.State = AccountStateUninit
	case stateRaw == "raw.accountState":
		acc.State = AccountStateActive
	case frozenHashRaw != "":
		acc.State = AccountStateFrozen
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hashRaw)
	if err != nil {
		return errors.Wrapf(err, "unable to decode last tx hash from %q", hashRaw)
	}

	copy(acc.LastTxHash[:], hashBytes)
	acc.LastTxLT = items[2].Uint()

	return nil
}
