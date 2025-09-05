package runner

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
)

// SuiGetGatewaySUIBalance retrieves the SUI balance of the Sui gateway by its ID.
func (r *E2ERunner) SuiGetGatewaySUIBalance() (*big.Int, error) {
	return suiGetGatewaySUIBalance(r.Ctx, r.Clients.Sui, r.SuiGateway.ObjectID())
}

func suiGetGatewaySUIBalance(ctx context.Context, client sui.ISuiAPI, gatewayID string) (*big.Int, error) {
	// get gateway object
	gatewayObject, err := client.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: gatewayID,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// extract the vault bag object ID from the gateway object
	vaults := gatewayObject.Data.Content.Fields["vaults"].(map[string]interface{})
	vaultsFields := vaults["fields"].(map[string]interface{})
	bagID := vaultsFields["id"].(map[string]interface{})["id"].(string)

	if bagID == "" {
		return nil, errors.New("vault ID is empty")
	}

	// get the Sui Dynamic Field for the vault bag object
	res, err := client.SuiXGetDynamicField(ctx, models.SuiXGetDynamicFieldRequest{
		ObjectId: bagID,
	})
	if err != nil {
		return nil, err
	}

	// extract the vault object for the SUI token from the dynamic fields
	suiVaultID := ""
	for _, field := range res.Data {
		if strings.Contains(field.ObjectType, "2::sui::SUI") {
			suiVaultID = field.ObjectId
			break
		}
	}
	if suiVaultID == "" {
		return nil, errors.New("SUI vault not found in the gateway object")
	}

	// get the SUI vault object
	suiVaultObject, err := client.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: suiVaultID,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// extract the balance from the SUI vault object
	balance := suiVaultObject.Data.Content.SuiMoveObject.Fields["value"].(map[string]interface{})["fields"].(map[string]interface{})["balance"].(string)
	if balance == "" {
		return nil, errors.New("balance is empty in the SUI vault object")
	}

	// convert string to int64
	balanceInt, err := strconv.ParseInt(balance, 10, 64)
	if err != nil {
		return nil, errors.New("failed to parse balance: " + err.Error())
	}

	return big.NewInt(balanceInt), nil
}
