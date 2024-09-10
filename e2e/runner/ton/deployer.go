package ton

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/wallet"
)

type Deployer struct {
	client *liteapi.Client
	wallet *wallet.Wallet
}

func NewDeployer(client *liteapi.Client, cfg Faucet) (*Deployer, error) {
	version := wallet.V3R2
	if cfg.WalletVersion != "V3R2" {
		return nil, fmt.Errorf("unsupported wallet version %q", cfg.WalletVersion)
	}

	pk, err := wallet.SeedToPrivateKey(cfg.Mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "invalid mnemonic")
	}

	w, err := wallet.New(pk, version, client, wallet.WithSubWalletID(uint32(cfg.SubWalletId)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet")
	}

	return &Deployer{client: client, wallet: &w}, nil
}

func (d *Deployer) Wallet() *wallet.Wallet {
	return d.wallet
}

func (d *Deployer) GetBalance(ctx context.Context) (math.Uint, error) {
	b, err := d.wallet.GetBalance(ctx)

	return math.NewUint(b), err
}

func (d *Deployer) Deploy(ctx context.Context, code, state *boc.Cell) (string, error) {
	// todo
	return "", nil
}
