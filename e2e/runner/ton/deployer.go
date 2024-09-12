package ton

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

// Deployer represents a wrapper around ton Wallet with some helpful methods.
type Deployer struct {
	wallet.Wallet
	blockchain blockchain
}

type blockchain interface {
	GetSeqno(ctx context.Context, account ton.AccountID) (uint32, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
	WaitForBlocks(ctx context.Context) error
}

// NewDeployer deployer constructor.
func NewDeployer(client blockchain, cfg Faucet) (*Deployer, error) {
	// this is a bit outdated, but we can't change it (it's created by my-local-ton)
	const version = wallet.V3R2
	if cfg.WalletVersion != "V3R2" {
		return nil, fmt.Errorf("unsupported wallet version %q", cfg.WalletVersion)
	}

	if cfg.WorkChain != 0 {
		return nil, fmt.Errorf("unsupported workchain id %d", cfg.WorkChain)
	}

	pk, err := wallet.SeedToPrivateKey(cfg.Mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "invalid mnemonic")
	}

	// #nosec G115 always in range
	w, err := wallet.New(pk, version, client, wallet.WithSubWalletID(uint32(cfg.SubWalletId)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet")
	}

	return &Deployer{Wallet: w, blockchain: client}, nil
}

func (d *Deployer) Seqno(ctx context.Context) (uint32, error) {
	return d.blockchain.GetSeqno(ctx, d.GetAddress())
}

// GetBalanceOf returns the balance of the given account.
func (d *Deployer) GetBalanceOf(ctx context.Context, id ton.AccountID) (math.Uint, error) {
	if err := d.waitForAccountActivation(ctx, id); err != nil {
		return math.Uint{}, errors.Wrap(err, "failed to wait for account activation")
	}

	state, err := d.blockchain.GetAccountState(ctx, id)
	if err != nil {
		return math.Uint{}, errors.Wrapf(err, "failed to get account %s state", id.ToRaw())
	}

	balance := uint64(state.Account.Account.Storage.Balance.Grams)

	return math.NewUint(balance), nil
}

// Fund sends the given amount of coins to the recipient. Returns tx hash and error.
func (d *Deployer) Fund(ctx context.Context, recipient ton.AccountID, amount math.Uint) (ton.Bits256, error) {
	msg := wallet.SimpleTransfer{
		Amount:  UintToCoins(amount),
		Address: recipient,
	}

	return d.send(ctx, msg, true)
}

// Deploy deploys AccountInit with the given amount of coins. Returns tx hash and error.
func (d *Deployer) Deploy(ctx context.Context, account *AccountInit, amount math.Uint) error {
	msg := wallet.Message{
		Amount:  UintToCoins(amount),
		Address: account.ID,
		Code:    account.Code,
		Data:    account.Data,
		Mode:    1, // pay gas fees separately
	}

	if _, err := d.send(ctx, msg, true); err != nil {
		return err
	}

	return d.waitForAccountActivation(ctx, account.ID)
}

func (d *Deployer) CreateWallet(ctx context.Context, amount math.Uint) (*wallet.Wallet, error) {
	seed := wallet.RandomSeed()

	accInit, w, err := ConstructWalletFromSeed(seed, d.blockchain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct wallet")
	}

	if err := d.Deploy(ctx, accInit, amount); err != nil {
		return nil, errors.Wrap(err, "failed to deploy wallet")
	}

	// Double-check the balance
	b, err := w.GetBalance(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get balance")
	}

	if b == 0 {
		return nil, fmt.Errorf("balance of %s is zero", w.GetAddress().ToRaw())
	}

	return w, nil
}

func (d *Deployer) send(ctx context.Context, message wallet.Sendable, waitForBlocks bool) (ton.Bits256, error) {
	// 2-3 blocks
	const maxWaitingTime = 18 * time.Second

	seqno, err := d.Seqno(ctx)
	if err != nil {
		return ton.Bits256{}, errors.Wrap(err, "failed to get seqno")
	}

	// Note that message hash IS NOT a tra hash.
	// It's not possible to get TX hash after tx sending
	msgHash, err := d.Wallet.SendV2(ctx, 0, message)
	if err != nil {
		return msgHash, errors.Wrap(err, "failed to send message")
	}

	if err := d.waitForNextSeqno(ctx, seqno, maxWaitingTime); err != nil {
		return msgHash, errors.Wrap(err, "failed to wait for confirmation")
	}

	if waitForBlocks {
		return msgHash, d.blockchain.WaitForBlocks(ctx)
	}

	return msgHash, nil
}

func (d *Deployer) waitForNextSeqno(ctx context.Context, oldSeqno uint32, timeout time.Duration) error {
	t := time.Now()

	for ; time.Since(t) < timeout; time.Sleep(timeout / 10) {
		newSeqno, err := d.Seqno(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get seqno")
		}

		if newSeqno > oldSeqno {
			return nil
		}
	}

	return errors.New("waiting confirmation timeout")
}

func (d *Deployer) waitForAccountActivation(ctx context.Context, account ton.AccountID) error {
	const interval = 5 * time.Second

	for i := 0; i < 10; i++ {
		state, err := d.blockchain.GetAccountState(ctx, account)
		if err != nil {
			return err
		}

		if state.Account.Status() == tlb.AccountActive {
			return nil
		}

		time.Sleep(interval)
	}

	return fmt.Errorf("account %s is not active", account.ToRaw())
}
