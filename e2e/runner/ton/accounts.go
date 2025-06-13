package ton

import (
	"encoding/hex"
	"fmt"
	"strings"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"golang.org/x/crypto/ed25519"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

const workchainID = 0

type AccountInit struct {
	Code      *boc.Cell
	Data      *boc.Cell
	StateInit *tlb.StateInit
	ID        ton.AccountID
}

func PrivateKeyFromHex(raw string) (ed25519.PrivateKey, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(raw, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode private key")
	}

	if len(b) != ed25519.SeedSize {
		return nil, errors.New("invalid private key length")
	}

	return ed25519.NewKeyFromSeed(b), nil
}

// ConstructWalletFromSeed constructs wallet AccountInit from seed.
// Used for wallets deployment.
func ConstructWalletFromSeed(seed string, client *Client) (*AccountInit, *wallet.Wallet, error) {
	if seed == "" {
		return nil, nil, errors.New("seed is empty")
	}

	pk, err := wallet.SeedToPrivateKey(seed)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "invalid mnemonic")
	}

	return ConstructWalletFromPrivateKey(pk, client)
}

func ConstructWalletFromPrivateKey(pk ed25519.PrivateKey, client *Client) (*AccountInit, *wallet.Wallet, error) {
	const version = wallet.V5R1

	w, err := wallet.New(pk, version, client.tongoAdapter())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create wallet")
	}

	stateInit, err := w.StateInit()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get state init for wallet")
	}

	var (
		code = stateInit.Code.Value.Value
		data = stateInit.Data.Value.Value
	)

	accInit, err := ConstructAccount(&code, &data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to construct account")
	}

	// If state init and internal tongo's wallet.stateInit are the same,
	// then addresses should match.
	if accInit.ID.String() != w.GetAddress().String() {
		return nil, nil, errors.New("account init doesn't match to created wallet")
	}

	return accInit, &w, nil
}

func ConstructAccount(code, data *boc.Cell) (*AccountInit, error) {
	stateInit := generateStateInit(code, data)

	id, err := generateAddress(workchainID, stateInit)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate address")
	}

	return &AccountInit{
		Code:      code,
		Data:      data,
		StateInit: stateInit,
		ID:        id,
	}, nil
}

// ConstructGatewayAccount constructs gateway AccountInit.
// - Authority is the address of the gateway "admin".
// - TSS is the EVM address of TSS.
// - Deposits are enabled by default.
func ConstructGatewayAccount(authority ton.AccountID, tss eth.Address) (*AccountInit, error) {
	return ConstructAccount(
		toncontracts.GatewayCode(),
		toncontracts.GatewayStateInit(authority, tss, true),
	)
}

// copied from tongo wallets_common.go
func generateStateInit(code, data *boc.Cell) *tlb.StateInit {
	return &tlb.StateInit{
		Code: tlb.Maybe[tlb.Ref[boc.Cell]]{
			Exists: true,
			Value:  tlb.Ref[boc.Cell]{Value: *code},
		},
		Data: tlb.Maybe[tlb.Ref[boc.Cell]]{
			Exists: true,
			Value:  tlb.Ref[boc.Cell]{Value: *data},
		},
	}
}

// copied from tongo wallets_common.go
func generateAddress(workchain int32, stateInit *tlb.StateInit) (ton.AccountID, error) {
	stateCell := boc.NewCell()
	if err := tlb.Marshal(stateCell, stateInit); err != nil {
		return ton.AccountID{}, fmt.Errorf("can not marshal wallet state: %v", err)
	}

	h, err := stateCell.Hash()
	if err != nil {
		return ton.AccountID{}, err
	}

	var hash tlb.Bits256
	copy(hash[:], h[:])

	return ton.AccountID{Workchain: workchain, Address: hash}, nil
}
