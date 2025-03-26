package signer

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/block-vision/sui-go-sdk/models"
)

// checkObjectsAreShared checks if the provided object ID list represents Sui shared objects
// when doing a withdrawAndCall, only share objects can be used
func checkObjectsAreShared(ctx context.Context, client interface {
	SuiMultiGetObjects(ctx context.Context, req models.SuiMultiGetObjectsRequest) ([]*models.SuiObjectResponse, error)
}, objectIDs []string) error {
	// no object can eventually be provided
	if len(objectIDs) == 0 {
		return nil
	}

	res, err := client.SuiMultiGetObjects(ctx, models.SuiMultiGetObjectsRequest{
		ObjectIds: objectIDs,
		Options: models.SuiObjectDataOptions{
			ShowOwner: true,
		},
	})
	if err != nil {
		return errors.Wrap(err, "unable to get objects")
	}

	// should always be the case, we add this check as a extra safety measure to ensure an object is not skipped
	if len(res) != len(objectIDs) {
		return fmt.Errorf("expected %d objects, but got %d", len(objectIDs), len(res))
	}

	for i, obj := range res {
		if obj.Data == nil {
			return fmt.Errorf("object %d is missing data", i)
		}

		switch owner := obj.Data.Owner.(type) {
		case string:
			if owner != "Immutable" {
				return fmt.Errorf("object %d has unexpected string owner: %s", i, owner)
			}
			// Immutable is valid, continue
		case map[string]interface{}:
			if _, isShared := owner["Shared"]; !isShared {
				return fmt.Errorf("object %d is not shared or immutable: owner = %+v", i, owner)
			}
			// Shared is valid, continue
		default:
			return fmt.Errorf("object %d has unknown owner type: %+v", i, obj.Data.Owner)
		}
	}

	return nil
}
