package runner

import (
	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils/sui"
	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
)

// SUIDeposit calls Deposit on Sui
func (r *E2ERunner) SUIDeposit(
	receiver ethcommon.Address,
	amount math.Uint,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectId := r.splitSUI(signer, amount)

	// create the tx
	return r.executeDeposit(signer, string(zetasui.SUI), coinObjectId, receiver)
}

// SUIDepositAndCall calls DepositAndCall on Sui
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
	return r.executeDepositAndCall(signer, string(zetasui.SUI), coinObjectId, receiver, payload)
}

// SuiFungibleTokenDeposit calls Deposit with fungible token on Sui
func (r *E2ERunner) SuiFungibleTokenDeposit(
	receiver ethcommon.Address,
	amount math.Uint,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectId := r.splitSUI(signer, amount)

	// create the tx
	return r.executeDeposit(signer, r.SuiTokenCoinType, coinObjectId, receiver)
}

// SuiFungibleTokenDepositAndCall calls DepositAndCall with fungible token on Sui
func (r *E2ERunner) SuiFungibleTokenDepositAndCall(
	receiver ethcommon.Address,
	amount math.Uint,
	payload []byte,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectId := r.splitUSDC(signer, amount)

	// create the tx
	return r.executeDepositAndCall(signer, r.SuiTokenCoinType, coinObjectId, receiver, payload)
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

// executeDeposit executes a deposit on the SUI contract
func (r *E2ERunner) executeDeposit(
	signer *sui.SignerSecp256k1,
	coinType string,
	coinObjectID string,
	receiver ethcommon.Address,
) models.SuiTransactionBlockResponse {
	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.GatewayPackageID,
		Module:          "gateway",
		Function:        "deposit",
		TypeArguments:   []any{coinType},
		Arguments:       []any{r.GatewayObjectID, coinObjectID, receiver.Hex()},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	return r.executeSuiTx(signer, tx)
}

// executeSUIDepositAndCall executes a depositAndCall on the SUI contract
func (r *E2ERunner) executeDepositAndCall(
	signer *sui.SignerSecp256k1,
	coinType string,
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
		TypeArguments:   []any{coinType},
		Arguments:       []any{r.GatewayObjectID, coinObjectID, receiver.Hex(), payload},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	return r.executeSuiTx(signer, tx)
}

// splitUSDC splits USDC coin and obtain a USDC coin object with the wanted balance
func (r *E2ERunner) splitUSDC(signer *sui.SignerSecp256k1, balance math.Uint) (objId string) {
	// find the coin to split
	originalCoin := r.findSuiCoinWithBalanceAbove(signer.Address(), balance, r.SuiTokenCoinType)

	tx, err := r.Clients.Sui.SplitCoin(r.Ctx, models.SplitCoinRequest{
		Signer:       signer.Address(),
		CoinObjectId: originalCoin,
		SplitAmounts: []string{balance.String()},
		GasBudget:    "5000000000",
	})

	require.NoError(r, err)
	r.executeSuiTx(signer, tx)

	// find the split coin
	return r.findSuiCoinWithBalance(signer.Address(), balance, r.SuiTokenCoinType)
}

// splitSUI splits SUI coin and obtain a SUI coin object with the wanted balance
func (r *E2ERunner) splitSUI(signer *sui.SignerSecp256k1, balance math.Uint) (objId string) {
	// find the coin to split
	originalCoin := r.findSuiCoinWithBalanceAbove(signer.Address(), balance, string(zetasui.SUI))

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
	return r.findSuiCoinWithBalance(signer.Address(), balance, string(zetasui.SUI))
}

func (r *E2ERunner) findSuiCoinWithBalance(
	address string,
	balance math.Uint,
	coinType string,
) (coinId string) {
	return r.findSuiCoin(address, balance, coinType, func(a, b math.Uint) bool {
		return a.Equal(b)
	})
}

func (r *E2ERunner) findSuiCoinWithBalanceAbove(
	address string,
	balanceAbove math.Uint,
	coinType string,
) (coinId string) {
	return r.findSuiCoin(address, balanceAbove, coinType, func(a, b math.Uint) bool {
		return a.GTE(b)
	})
}

type compFunc func(a, b math.Uint) bool

func (r *E2ERunner) findSuiCoin(
	address string,
	balance math.Uint,
	coinType string,
	comp compFunc,
) (coinId string) {
	res, err := r.Clients.Sui.SuiXGetCoins(r.Ctx, models.SuiXGetCoinsRequest{
		Owner:    address,
		CoinType: coinType,
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
