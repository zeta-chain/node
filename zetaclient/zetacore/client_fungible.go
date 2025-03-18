package zetacore

import (
	"context"

	"cosmossdk.io/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/crypto"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// GetForeignCoinsFromAsset returns the foreign coin for a given asset for a given chain ID
func (c *Client) GetForeignCoinsFromAsset(
	ctx context.Context,
	chainID int64,
	assetAddress ethcommon.Address,
) (fungibletypes.ForeignCoins, error) {
	// convert asset to checksum address or empty string (for Gas asset)
	assetString := assetAddress.Hex()
	if crypto.IsEmptyAddress(assetAddress) {
		assetString = ""
	}

	request := &fungibletypes.QueryGetForeignCoinsFromAssetRequest{
		ChainId: chainID,
		Asset:   assetString,
	}

	resp, err := c.Fungible.ForeignCoinsFromAsset(ctx, request)
	if err != nil {
		return fungibletypes.ForeignCoins{}, errors.Wrapf(
			err,
			"unable to get foreign coins for asset %s on chain %d",
			assetString,
			chainID,
		)
	}

	return resp.ForeignCoins, nil
}
