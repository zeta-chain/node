package ton

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

// Deployer represents a wrapper around ton Wallet with some helpful methods.
type Deployer struct {
	wallet.Wallet
	client *Client
}

// NewDeployer deployer constructor.
func NewDeployer(client *Client, cfg Faucet) (*Deployer, error) {
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
	w, err := wallet.New(pk, version, client.tongoAdapter(), wallet.WithSubWalletID(uint32(cfg.SubWalletId)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet")
	}

	return &Deployer{Wallet: w, client: client}, nil
}

// Fund sends the given amount of coins to the recipient. Returns tx hash and error.
func (d *Deployer) Fund(ctx context.Context, recipient ton.AccountID, amount math.Uint) (ton.Bits256, error) {
	msg := wallet.SimpleTransfer{
		Amount:  toncontracts.UintToCoins(amount),
		Address: recipient,
	}

	return d.send(ctx, msg, true)
}

// Deploy deploys AccountInit with the given amount of coins. Returns tx hash and error.
func (d *Deployer) Deploy(ctx context.Context, account *AccountInit, amount math.Uint) error {
	msg := wallet.Message{
		Amount:  toncontracts.UintToCoins(amount),
		Address: account.ID,
		Code:    account.Code,
		Data:    account.Data,
		Mode:    toncontracts.SendFlagSeparateFees,
	}

	if _, err := d.send(ctx, msg, true); err != nil {
		return errors.Wrapf(err, "unable to deploy account %q", account.ID.ToRaw())
	}

	return d.client.WaitForAccountActivation(ctx, account.ID)
}

func (d *Deployer) send(ctx context.Context, message wallet.Sendable, waitForBlocks bool) (ton.Bits256, error) {
	// 2-3 blocks
	const maxWaitingTime = 18 * time.Second

	id := d.GetAddress()

	seqno, err := d.client.GetSeqno(ctx, id)
	if err != nil {
		return ton.Bits256{}, errors.Wrap(err, "failed to get seqno")
	}

	// Note that message hash IS NOT a tx hash.
	// It's not possible to get tx hash right after tx sending
	msgHash, err := d.Wallet.SendV2(ctx, 0, message)
	if err != nil {
		return msgHash, errors.Wrap(err, "failed to send message")
	}

	if err := d.client.WaitForNextSeqno(ctx, id, seqno, maxWaitingTime); err != nil {
		return msgHash, errors.Wrap(err, "failed to wait for confirmation")
	}

	if waitForBlocks {
		return msgHash, d.client.WaitForBlocks(ctx)
	}

	return msgHash, nil
}
