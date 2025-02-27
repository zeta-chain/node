package signer

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const withdrawCapTTL = 5 * time.Minute

// withdrawCap represents WithdrawCap (capability) object
// that is required as a "permission" to withdraw funds.
// Should belong to TSS address on Sui.
type withdrawCap struct {
	objectID  string
	mu        sync.RWMutex
	fetchedAt time.Time
}

func (wc *withdrawCap) valid() bool {
	wc.mu.RLock()
	defer wc.mu.RUnlock()

	if wc.objectID == "" {
		return false
	}

	return time.Since(wc.fetchedAt) < withdrawCapTTL
}

func (wc *withdrawCap) set(objectID string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	wc.objectID = objectID
	wc.fetchedAt = time.Now()
}

// getWithdrawCapIDCached getWithdrawCapID with withdrawCapTTL cache.
func (s *Signer) getWithdrawCapIDCached(ctx context.Context) (string, error) {
	if s.withdrawCap.valid() {
		return s.withdrawCap.objectID, nil
	}

	s.Logger().Std.Info().Msg("WithdrawCap cache expired, fetching new objectID")

	objectID, err := s.getWithdrawCapID(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get withdraw cap ID")
	}

	s.withdrawCap.set(objectID)

	s.Logger().Std.Info().Str("sui.object_id", objectID).Msg("WithdrawCap objectID fetched")

	return objectID, nil
}

// getWithdrawCapID returns the objectID of the WithdrawCap. Should belong to TSS address on Sui.
func (s *Signer) getWithdrawCapID(ctx context.Context) (string, error) {
	owner := s.TSS().PubKey().AddressSui()
	structType := s.gateway.WithdrawCapType()

	objectID, err := s.client.GetOwnedObjectID(ctx, owner, structType)
	if err != nil {
		return "", errors.Wrap(err, "unable to get owned object ID")
	}

	return objectID, nil
}
