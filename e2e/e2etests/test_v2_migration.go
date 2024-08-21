package e2etests

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/txserver"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
)

func TestV2Migration(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// Part 1: add new admin authorization
	r.Logger.Info("Part 1: Adding authorization for new v2 contracts")
	err := r.ZetaTxServer.AddAuthorization("/zetachain.zetacore.crosschain.MsgUpdateERC20CustodyPauseStatus")
	require.NoError(r, err)

	err = r.ZetaTxServer.AddAuthorization("/zetachain.zetacore.crosschain.MsgMigrateERC20CustodyFunds")
	require.NoError(r, err)

	err = r.ZetaTxServer.AddAuthorization("/zetachain.zetacore.fungible.MsgUpdateGatewayContract")
	require.NoError(r, err)

	// Part 2: deploy v2 contracts on EVM chain
	r.Logger.Info("Part 2: Deploying v2 contracts on EVM chain")
	r.SetupEVMV2()

	// Part 3: deploy gateway on ZetaChain
	r.Logger.Info("Part 3: Deploying Gateway ZEVM")
	r.SetZEVMContractsV2()

	// Part 4: upgrade all ZRC20s
	r.Logger.Info("Part 4: Upgrading ZRC20s")
	upgradeZRC20s(r)

	// Part 5: migrate ERC20 custody funds
	r.Logger.Info("Part 5: Migrating ERC20 custody funds")
	migrateERC20CustodyFunds(r)
}

// upgradeZRC20s upgrades all ZRC20s to the new version
func upgradeZRC20s(r *runner.E2ERunner) {
	// get chain IDs
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)
	btcChainID := r.GetBitcoinChainID()

	// upgrade ETH ZRC20
	r.Logger.Info("Upgrading ETH ZRC20")
	upgradeZRC20(r, r.ETHZRC20Addr, r.ETHZRC20, evmChainID, uint8(coin.CoinType_Gas))

	// upgrade ERC20 ZRC20
	r.Logger.Info("Upgrading ERC20 ZRC20")
	upgradeZRC20(r, r.ERC20ZRC20Addr, r.ERC20ZRC20, evmChainID, uint8(coin.CoinType_ERC20))

	// upgrade BTC ZRC20
	r.Logger.Info("Upgrading BTC ZRC20")
	upgradeZRC20(r, r.BTCZRC20Addr, r.BTCZRC20, big.NewInt(btcChainID), uint8(coin.CoinType_Gas))
}

// zrc20Caller is an interface to call ZRC20 functions
type zrc20Caller interface {
	Name(opts *bind.CallOpts) (string, error)
	Symbol(opts *bind.CallOpts) (string, error)
	Decimals(opts *bind.CallOpts) (uint8, error)
}

// upgradeZRC20 upgrades a ZRC20 to the new version
func upgradeZRC20(
	r *runner.E2ERunner,
	zrc20Addr common.Address,
	zrc20Caller zrc20Caller,
	chainID *big.Int,
	coinType uint8,
) {
	// deploy new ZRC20 version
	name, err := zrc20Caller.Name(&bind.CallOpts{})
	require.NoError(r, err)
	symbol, err := zrc20Caller.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	decimal, err := zrc20Caller.Decimals(&bind.CallOpts{})
	require.NoError(r, err)

	newZRC20Addr, newZRC20Tx, _, err := zrc20.DeployZRC20(
		r.ZEVMAuth,
		r.ZEVMClient,
		name,
		symbol,
		decimal,
		chainID,
		coinType,
		big.NewInt(100_000),
		r.SystemContractAddr,
		r.GatewayZEVMAddr,
	)
	require.NoError(r, err)

	// wait tx to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, newZRC20Tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, ethtypes.ReceiptStatusSuccessful, receipt.Status)

	// upgrade ZRC20 bytecode with the one of the new ZRC20
	codeHashRes, err := r.FungibleClient.CodeHash(r.Ctx, &fungibletypes.QueryCodeHashRequest{
		Address: newZRC20Addr.String(),
	})
	require.NoError(r, err)

	msg := fungibletypes.NewMsgUpdateContractBytecode(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		zrc20Addr.Hex(),
		codeHashRes.CodeHash,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)
}

func migrateERC20CustodyFunds(r *runner.E2ERunner) {
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// Part 1: pause the ERC20Custody v1
	r.Logger.Info("Pausing ERC20 custody v1 contract")
	msgPausing := crosschaintypes.NewMsgUpdateERC20CustodyPauseStatus(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		evmChainID.Int64(),
		true,
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgPausing)
	require.NoError(r, err)

	// fetch cctx index from tx response
	cctxIndex, err := txserver.FetchAttributeFromTxResponse(res, "cctx_index")
	require.NoError(r, err)

	cctxRes, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
	require.NoError(r, err)

	cctx := cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "pausing")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(cctxIndex)

	// Part 2: pause the ZRC20 ERC20
	msgPause := fungibletypes.NewMsgPauseZRC20(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		[]string{r.ERC20ZRC20Addr.Hex()},
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msgPause)
	require.NoError(r, err)

	// Part 3: migrate all funds of the ERC20
	balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	require.NoError(r, err)

	// ensure balance is not zero to ensure the test tests actual migration
	require.NotEqual(r, int64(0), balance.Int64())

	// send MigrateERC20CustodyFunds command
	msgMigration := crosschaintypes.NewMsgMigrateERC20CustodyFunds(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		evmChainID.Int64(),
		r.ERC20CustodyV2Addr.Hex(),
		r.ERC20Addr.Hex(),
		sdkmath.NewUintFromBigInt(balance),
	)
	res, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgMigration)
	require.NoError(r, err)

	// fetch cctx index from tx response
	cctxIndex, err = txserver.FetchAttributeFromTxResponse(res, "cctx_index")
	require.NoError(r, err)

	cctxRes, err = r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
	require.NoError(r, err)

	cctx = cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "migration")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(cctxIndex)

	// Part 4: unpause the ZRC20
	msgUnpause := fungibletypes.NewMsgUnpauseZRC20(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		[]string{r.ERC20ZRC20Addr.Hex()},
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgUnpause)
	require.NoError(r, err)
}
