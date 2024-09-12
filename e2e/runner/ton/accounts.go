package ton

import (
	_ "embed"
	"encoding/json"
	"fmt"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"golang.org/x/crypto/ed25519"
)

const workchainID = 0

// https://github.com/zeta-chain/protocol-contracts-ton
// `make compile`
//
//go:embed gateway.compiled.json
var tonGatewayCodeJSON []byte

type AccountInit struct {
	Code      *boc.Cell
	Data      *boc.Cell
	StateInit *tlb.StateInit
	ID        ton.AccountID
}

// ConstructWalletFromSeed constructs wallet AccountInit from seed.
// Used for wallets deployment.
func ConstructWalletFromSeed(seed string, client blockchain) (*AccountInit, *wallet.Wallet, error) {
	pk, err := wallet.SeedToPrivateKey(seed)
	if err != nil {
		return nil, nil, errors.Wrap(err, "invalid mnemonic")
	}

	return ConstructWalletFromPrivateKey(pk, client)
}

func ConstructWalletFromPrivateKey(pk ed25519.PrivateKey, client blockchain) (*AccountInit, *wallet.Wallet, error) {
	const version = wallet.V5R1

	w, err := wallet.New(pk, version, client)
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

func ConstructGatewayAccount(tss eth.Address) (*AccountInit, error) {
	code, err := getGatewayCode()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get TON Gateway code")
	}

	data, err := buildGatewayData(tss)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build TON Gateway data")
	}

	return ConstructAccount(code, data)
}

func getGatewayCode() (*boc.Cell, error) {
	var code struct {
		Hex string `json:"hex"`
	}

	if err := json.Unmarshal(tonGatewayCodeJSON, &code); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal TON Gateway code")
	}

	cells, err := boc.DeserializeBocHex(code.Hex)
	if err != nil {
		return nil, errors.Wrap(err, "unable to deserialize TON Gateway code")
	}

	if len(cells) != 1 {
		return nil, errors.New("invalid cells count")
	}

	return cells[0], nil
}

// buildGatewayState returns TON Gateway initial state data cell
func buildGatewayData(tss eth.Address) (*boc.Cell, error) {
	const evmAddressBits = 20 * 8

	tssSlice := boc.NewBitString(evmAddressBits)
	if err := tssSlice.WriteBytes(tss.Bytes()); err != nil {
		return nil, errors.Wrap(err, "unable to convert TSS address to ton slice")
	}

	var (
		zeroCoins = tlb.Coins(0)
		enc       = &tlb.Encoder{}
		cell      = boc.NewCell()
	)

	err := errCollect(
		cell.WriteBit(true),             // deposits_enabled
		zeroCoins.MarshalTLB(cell, enc), // total_locked
		zeroCoins.MarshalTLB(cell, enc), // fees
		cell.WriteUint(0, 32),           // seqno
		cell.WriteBitString(tssSlice),   // tss_address
	)

	if err != nil {
		return nil, errors.Wrap(err, "unable to write TON Gateway state cell")
	}

	return cell, nil
}

func errCollect(errs ...error) error {
	for i, err := range errs {
		if err != nil {
			return errors.Wrapf(err, "error at index %d", i)
		}
	}

	return nil
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
