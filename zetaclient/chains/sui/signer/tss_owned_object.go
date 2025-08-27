package signer

import (
	"context"
	"sync"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/contracts/sui"
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

// getMessageContextIDCached getMessageContextID with tssOwnedObjectTTL cache.
func (s *Signer) getMessageContextIDCached(ctx context.Context) (string, error) {
	if s.messageContext.valid() {
		return s.messageContext.objectID, nil
	}

	s.Logger().Std.Info().Msg("MessageContext cache expired, fetching new objectID")

	objectID, err := s.getMessageContextID(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get message context ID")
	}

	s.messageContext.set(objectID)

	s.Logger().Std.Info().Str("sui.object_id", objectID).Msg("MessageContext objectID fetched")

	return objectID, nil
}

// getMessageContextID returns the objectID of the active MessageContext. Should belong to TSS address on Sui.
func (s *Signer) getMessageContextID(ctx context.Context) (string, error) {
	nameJSON, err := sui.ActiveMessageContextDynamicFieldName()
	if err != nil {
		return "", errors.Wrap(err, "unable to get dynamic field name")
	}

	response, err := s.client.SuiXGetDynamicFieldObject(ctx, models.SuiXGetDynamicFieldObjectRequest{
		ObjectId: s.gateway.ObjectID(),
		DynamicFieldName: models.DynamicFieldObjectName{
			Type:  "vector<u8>",
			Value: nameJSON,
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "unable to get message context dynamic field object")
	}

	if response.Data == nil || response.Data.Content == nil {
		return "", errors.New("dynamic field object data is nil")
	}

	messageContextID, err := sui.ParseDynamicFieldValueStr(*response.Data.Content)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse message context ID")
	}

	return messageContextID, nil
}
