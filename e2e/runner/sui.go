package runner

import (
	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils/sui"
	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
)

// SUIDeposit calls Deposit on SUI
func (r *E2ERunner) SUIDeposit(
	receiver ethcommon.Address,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectId, err := r.splitSUI(signer, math.NewUint(100000))
	require.NoError(r, err)

	// create the tx
	resp, err := r.executeSUIDeposit(signer, coinObjectId, receiver)
	require.NoError(r, err)

	return resp
}

// executeSUIDeposit executes a deposit on the SUI contract
func (r *E2ERunner) executeSUIDeposit(
	signer *sui.SignerSecp256k1,
	coinObjectID string,
	receiver ethcommon.Address,
) (models.SuiTransactionBlockResponse, error) {
	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.GatewayPackageID,
		Module:          "gateway",
		Function:        "deposit",
		TypeArguments:   []any{string(zetasui.SUI)},
		Arguments:       []any{r.GatewayObjectID, coinObjectID, receiver.Hex()},
		GasBudget:       "5000000000",
	})
	if err != nil {
		return models.SuiTransactionBlockResponse{}, err
	}

	// sign the tx
	signature, err := signer.SignTransactionBlock(tx.TxBytes)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, err
	}
	resp, err := r.Clients.Sui.SuiExecuteTransactionBlock(r.Ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:     tx.TxBytes,
		Signature:   []string{signature},
		RequestType: "WaitForLocalExecution",
	})
	if err != nil {
		return models.SuiTransactionBlockResponse{}, err
	}

	return resp, nil
}

// splitSUI splits SUI coin and obtain a SUI coin object with the wanted balance
func (r *E2ERunner) splitSUI(signer *sui.SignerSecp256k1, balance math.Uint) (objId string, err error) {
	return r.findSUIWithBalanceAbove(signer.Address(), balance)

	// TODO: there is an issue splitting the coin -> fix this
	//// find the coin to split
	//originalCoin, err := r.findSUIWithBalanceAbove(signer.Address(), balance)
	//if err != nil {
	//	return "", err
	//}
	//
	//// split the coin using the PaySui API
	//tx, err := r.Clients.Sui.PaySui(r.Ctx, models.PaySuiRequest{
	//	Signer:      signer.Address(),
	//	SuiObjectId: []string{originalCoin},
	//	Recipient:   []string{signer.Address()},
	//	Amount:      []string{balance.String()},
	//	GasBudget:   "5000000000",
	//})
	//
	//// sign the split tx
	//_, err = signer.SignTransactionBlock(tx.TxBytes)
	//require.NoError(r, err)
	//
	//// find the split coin
	//return r.findSUIWithBalance(signer.Address(), balance)
}

func (r *E2ERunner) findSUIWithBalance(
	address string,
	balance math.Uint,
) (coinId string, err error) {
	return r.findSUI(address, balance, func(a, b math.Uint) bool {
		return a.Equal(b)
	})
}

func (r *E2ERunner) findSUIWithBalanceAbove(
	address string,
	balanceAbove math.Uint,
) (coinId string, err error) {
	return r.findSUI(address, balanceAbove, func(a, b math.Uint) bool {
		return a.GTE(b)
	})
}

type compFunc func(a, b math.Uint) bool

func (r *E2ERunner) findSUI(
	address string,
	balance math.Uint,
	comp compFunc,
) (coinId string, err error) {
	res, err := r.Clients.Sui.SuiXGetCoins(r.Ctx, models.SuiXGetCoinsRequest{
		Owner:    address,
		CoinType: string(zetasui.SUI),
	})
	if err != nil {
		return "", err
	}

	for _, data := range res.Data {
		coinBalance, err := math.ParseUint(data.Balance)
		if err != nil {
			return "", err
		}
		if comp(coinBalance, balance) {
			return data.CoinObjectId, nil
		}
	}

	return "", errors.New("no coin found")
}
