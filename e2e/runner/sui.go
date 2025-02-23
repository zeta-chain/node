package runner

import (
	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils/sui"
	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
)

// SUIDeposit calls Deposit on SUI
func (r *E2ERunner) SUIDeposit(
	receiver ethcommon.Address,
	amount math.Uint,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectId := r.splitSUI(signer, amount)

	// create the tx
	return r.executeSUIDeposit(signer, coinObjectId, receiver)
}

// SUIDepositAndCall calls DepositAndCall on SUI
func (r *E2ERunner) SUIDepositAndCall(
	receiver ethcommon.Address,
	amount math.Uint,
	payload []byte,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectId := r.splitSUI(signer, amount)

	// create the tx
	return r.executeSUIDepositAndCall(signer, coinObjectId, receiver, payload)
}

// MintSuiUSDC mints FakeUSDC on Sui to a receiver
// this function requires the signer to be the owner of the trasuryCap
func (r *E2ERunner) MintSuiUSDC(
	amount,
	receiver string,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.GatewayPackageID,
		Module:          "fake_usdc",
		Function:        "mint",
		Arguments:       []any{r.SuiTokenTreasuryCap, amount, receiver},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	return r.executeSuiTx(signer, tx)
}

// executeSUIDeposit executes a deposit on the SUI contract
func (r *E2ERunner) executeSUIDeposit(
	signer *sui.SignerSecp256k1,
	coinObjectID string,
	receiver ethcommon.Address,
) models.SuiTransactionBlockResponse {
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
	require.NoError(r, err)

	return r.executeSuiTx(signer, tx)
}

// executeSUIDepositAndCall executes a depositAndCall on the SUI contract
func (r *E2ERunner) executeSUIDepositAndCall(
	signer *sui.SignerSecp256k1,
	coinObjectID string,
	receiver ethcommon.Address,
	payload []byte,
) models.SuiTransactionBlockResponse {
	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.GatewayPackageID,
		Module:          "gateway",
		Function:        "deposit_and_call",
		TypeArguments:   []any{string(zetasui.SUI)},
		Arguments:       []any{r.GatewayObjectID, coinObjectID, receiver.Hex(), payload},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	return r.executeSuiTx(signer, tx)
}

// splitSUI splits SUI coin and obtain a SUI coin object with the wanted balance
func (r *E2ERunner) splitSUI(signer *sui.SignerSecp256k1, balance math.Uint) (objId string) {
	// find the coin to split
	originalCoin := r.findSUIWithBalanceAbove(signer.Address(), balance)

	// split the coin using the PaySui API
	tx, err := r.Clients.Sui.PaySui(r.Ctx, models.PaySuiRequest{
		Signer:      signer.Address(),
		SuiObjectId: []string{originalCoin},
		Recipient:   []string{signer.Address()},
		Amount:      []string{balance.String()},
		GasBudget:   "5000000000",
	})
	require.NoError(r, err)

	r.executeSuiTx(signer, tx)

	// find the split coin
	return r.findSUIWithBalance(signer.Address(), balance)
}

func (r *E2ERunner) findSUIWithBalance(
	address string,
	balance math.Uint,
) (coinId string) {
	return r.findSUI(address, balance, func(a, b math.Uint) bool {
		return a.Equal(b)
	})
}

func (r *E2ERunner) findSUIWithBalanceAbove(
	address string,
	balanceAbove math.Uint,
) (coinId string) {
	return r.findSUI(address, balanceAbove, func(a, b math.Uint) bool {
		return a.GTE(b)
	})
}

type compFunc func(a, b math.Uint) bool

func (r *E2ERunner) findSUI(
	address string,
	balance math.Uint,
	comp compFunc,
) (coinId string) {
	res, err := r.Clients.Sui.SuiXGetCoins(r.Ctx, models.SuiXGetCoinsRequest{
		Owner:    address,
		CoinType: string(zetasui.SUI),
	})
	require.NoError(r, err)

	for _, data := range res.Data {
		coinBalance, err := math.ParseUint(data.Balance)
		require.NoError(r, err)

		if comp(coinBalance, balance) {
			return data.CoinObjectId
		}
	}

	require.FailNow(r, "coin not found")
	return ""
}

func (r *E2ERunner) executeSuiTx(
	signer *sui.SignerSecp256k1,
	tx models.TxnMetaData,
) models.SuiTransactionBlockResponse {
	// sign the tx
	signature, err := signer.SignTransactionBlock(tx.TxBytes)
	require.NoError(r, err, "sign transaction")

	resp, err := r.Clients.Sui.SuiExecuteTransactionBlock(r.Ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:     tx.TxBytes,
		Signature:   []string{signature},
		RequestType: "WaitForLocalExecution",
	})
	require.NoError(r, err)

	return resp
}
