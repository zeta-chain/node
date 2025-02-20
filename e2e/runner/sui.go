package runner

import (
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
)

// SUIDeposit calls Deposit on SUI
func (r *E2ERunner) SUIDeposit(
	receiver ethcommon.Address,
) {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	signerAddr := signer.Address()

	r.Logger.Print("SUI Address:")
	r.Logger.Print(signerAddr)

	// retrieve SUI coin object to deposit
	// TODO: split coin with an amountto not deposit all tokens
	coinObjectId, err := r.filterOwnedObject(signerAddr, string(zetasui.SUI))
	require.NoError(r, err)

	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signerAddr,
		PackageObjectId: r.GatewayPackageID,
		Module:          "gateway",
		Function:        "deposit",
		TypeArguments:   []any{string(zetasui.SUI)},
		Arguments:       []any{r.GatewayObjectID, coinObjectId, receiver.Hex()},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	// sign the tx
	resp, err := signer.SignTransactionBlock(tx.TxBytes)
	require.NoError(r, err)

	r.Logger.Print("tx response")
	r.Logger.Print(resp)
}

// filterOwnedObject finds one object owned by the address and has the type of typeName
func (r *E2ERunner) filterOwnedObject(address, typeName string) (objId string, err error) {
	// see https://docs.sui.io/sui-api-ref#suix_getownedobjects
	suiObjectResponseQuery := models.SuiObjectResponseQuery{
		Filter: models.SuiObjectDataFilter{
			"StructType": typeName,
		},
		Options: models.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
			ShowBcs:     true,
			ShowOwner:   true,
		},
	}

	resp, err := r.Clients.Sui.SuiXGetOwnedObjects(r.Ctx, models.SuiXGetOwnedObjectsRequest{
		Address: address,
		Query:   suiObjectResponseQuery,
		Limit:   50,
	})
	if err != nil {
		return "", err
	}

	for _, data := range resp.Data {
		if data.Data.Type == typeName {
			objId = data.Data.ObjectId
			return objId, nil
		}
	}

	return "", fmt.Errorf("no object of type %s found", typeName)
}
