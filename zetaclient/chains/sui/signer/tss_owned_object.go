package signer

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const tssOwnedObjectTTL = 5 * time.Minute

// tssOwnedObject represents an object that is owned by TSS address on Sui.
// There are two objects currently owned by TSS: WithdrawCap and MessageContext.
type tssOwnedObject struct {
	objectID  string
	mu        sync.RWMutex
	fetchedAt time.Time
}

func (wc *tssOwnedObject) valid() bool {
	wc.mu.RLock()
	defer wc.mu.RUnlock()

	if wc.objectID == "" {
		return false
	}

	return time.Since(wc.fetchedAt) < tssOwnedObjectTTL
}

func (wc *tssOwnedObject) set(objectID string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	wc.objectID = objectID
	wc.fetchedAt = time.Now()
}

// getWithdrawCapIDCached getWithdrawCapID with tssOwnedObjectTTL cache.
func (s *Signer) getWithdrawCapIDCached(ctx context.Context) (string, error) {
	if s.withdrawCap.valid() {
		return s.withdrawCap.objectID, nil
	}

	s.Logger().Std.Info().Msg("withdrawCap cache expired; fetching new objectID")

	objectID, err := s.getWithdrawCapID(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get withdraw cap ID")
	}

	s.withdrawCap.set(objectID)

	s.Logger().Std.Info().Str("sui_object_id", objectID).Msg("withdrawCap objectID fetched")

	return objectID, nil
}

// getWithdrawCapID returns the objectID of the WithdrawCap. Should belong to TSS address on Sui.
func (s *Signer) getWithdrawCapID(ctx context.Context) (string, error) {
	owner := s.TSS().PubKey().AddressSui()
	structType := s.gateway.WithdrawCapType()

	objectID, err := s.suiClient.GetOwnedObjectID(ctx, owner, structType)
	if err != nil {
		return "", errors.Wrap(err, "unable to get owned object ID")
	}

	return objectID, nil
}

// TODO: https://github.com/zeta-chain/node/issues/4066
// uncomment below helper functions used for authenticated call
// getMessageContextIDCached getMessageContextID with tssOwnedObjectTTL cache.
// func (s *Signer) getMessageContextIDCached(ctx context.Context) (string, error) {
// 	if s.messageContext.valid() {
// 		return s.messageContext.objectID, nil
// 	}

//	s.Logger().Std.Info().Msg("messageContext cache expired, fetching new objectID")

// 	objectID, err := s.getMessageContextID(ctx)
// 	if err != nil {
// 		return "", errors.Wrap(err, "unable to get message context ID")
// 	}

// 	s.messageContext.set(objectID)

//	s.Logger().Std.Info().Str("sui_object_id", objectID).Msg("messageContext objectID fetched")

// 	return objectID, nil
// }

// getMessageContextID returns the objectID of the MessageContext. Should belong to TSS address on Sui.
// func (s *Signer) getMessageContextID(ctx context.Context) (string, error) {
// 	owner := s.TSS().PubKey().AddressSui()
// 	structType := s.gateway.MessageContextType()

// 	objectID, err := s.client.GetOwnedObjectID(ctx, owner, structType)
// 	if err != nil {
// 		return "", errors.Wrap(err, "unable to get owned object ID")
// 	}

// 	return objectID, nil
// }
