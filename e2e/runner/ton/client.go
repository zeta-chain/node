package ton

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Client struct {
	*liteapi.Client
}

// Status checks the health of the TON node
func (c *Client) Status(ctx context.Context) error {
	_, err := c.GetMasterchainInfo(ctx)
	return err
}

// GetBalanceOf returns the balance of a given account.
// wait=true waits for account activation.
func (c *Client) GetBalanceOf(ctx context.Context, id ton.AccountID, wait bool) (math.Uint, error) {
	if wait {
		if err := c.WaitForAccountActivation(ctx, id); err != nil {
			return math.Uint{}, errors.Wrap(err, "failed to wait for account activation")
		}
	}

	state, err := c.GetAccountState(ctx, id)
	if err != nil {
		return math.Uint{}, errors.Wrapf(err, "failed to get account %s state", id.ToRaw())
	}

	balance := uint64(state.Account.Account.Storage.Balance.Grams)

	return math.NewUint(balance), nil
}

func (c *Client) WaitForBlocks(ctx context.Context) error {
	const (
		blocksToWait = 3
		interval     = 3 * time.Second
	)

	block, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return err
	}

	waitFor := block.Last.Seqno + blocksToWait

	for {
		freshBlock, err := c.GetMasterchainInfo(ctx)
		if err != nil {
			return err
		}

		if waitFor < freshBlock.Last.Seqno {
			return nil
		}

		time.Sleep(interval)
	}
}

func (c *Client) WaitForAccountActivation(ctx context.Context, account ton.AccountID) error {
	const interval = 5 * time.Second

	for i := 0; i < 10; i++ {
		state, err := c.GetAccountState(ctx, account)
		if err != nil {
			return err
		}

		if state.Account.Status() == tlb.AccountActive {
			return nil
		}

		time.Sleep(interval)
	}

	return fmt.Errorf("account %q is not active; timed out", account.ToRaw())
}

func (c *Client) WaitForNextSeqno(
	ctx context.Context,
	id ton.AccountID,
	oldSeqno uint32,
	timeout time.Duration,
) error {
	t := time.Now()

	for ; time.Since(t) < timeout; time.Sleep(timeout / 10) {
		newSeqno, err := c.GetSeqno(ctx, id)
		if err != nil {
			return errors.Wrap(err, "failed to get seqno")
		}

		if newSeqno > oldSeqno {
			return nil
		}
	}

	return errors.New("waiting confirmation timeout")
}
