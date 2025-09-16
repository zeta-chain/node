package e2etests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/contracts/erc20"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestWhitelistAsset tests the whitelist asset functionality
func TestWhitelistAsset(r *runner.E2ERunner, _ []string) {
	// Deploy a new ERC20 on the new EVM chain
	r.Logger.Info("Deploying new ERC20 contract")
	erc20Addr, txERC20, _, err := erc20.DeployERC20(r.EVMAuth, r.EVMClient, "NEWERC20", "NEWERC20", 6)
	require.NoError(r, err)

	// wait for the ERC20 to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txERC20, r.Logger, r.ReceiptTimeout)
	require.Equal(r, ethtypes.ReceiptStatusSuccessful, receipt.Status)

	// ERC20 test

	// whitelist erc20 zrc20
	r.Logger.Info("whitelisting ERC20 on new network")
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, crosschaintypes.NewMsgWhitelistAsset(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		erc20Addr.Hex(),
		chains.GoerliLocalnet.ChainId,
		"NEWERC20",
		"NEWERC20",
		6,
		100000,
		sdkmath.NewUintFromString("100000000000000000000000000"),
	))
	require.NoError(r, err)

	event, ok := txserver.EventOfType[*crosschaintypes.EventERC20Whitelist](res.Events)
	require.True(r, ok, "no EventERC20Whitelist in %s", res.TxHash)
	erc20zrc20Addr := event.Zrc20Address
	whitelistCCTXIndex := event.WhitelistCctxIndex

	err = r.ZetaTxServer.InitializeLiquidityCaps(erc20zrc20Addr)
	require.NoError(r, err)

	// ensure CCTX created
	resCCTX, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: whitelistCCTXIndex})
	require.NoError(r, err)

	cctx := resCCTX.CrossChainTx
	r.Logger.CCTX(*cctx, "whitelist_cctx")

	// wait for the whitelist cctx to be mined
	r.WaitForMinedCCTXFromIndex(whitelistCCTXIndex)

	// save old ERC20 attribute to set it back after the test
	oldERC20Addr := r.ERC20Addr
	oldERC20 := r.ERC20
	oldERC20ZRC20Addr := r.ERC20ZRC20Addr
	oldERC20ZRC20 := r.ERC20ZRC20
	defer func() {
		r.ERC20Addr = oldERC20Addr
		r.ERC20 = oldERC20
		r.ERC20ZRC20Addr = oldERC20ZRC20Addr
		r.ERC20ZRC20 = oldERC20ZRC20
	}()

	// set erc20 and zrc20 in runner
	require.True(r, ethcommon.IsHexAddress(erc20zrc20Addr), "invalid contract address: %s", erc20zrc20Addr)
	erc20zrc20AddrHex := ethcommon.HexToAddress(erc20zrc20Addr)
	erc20ZRC20, err := zrc20.NewZRC20(erc20zrc20AddrHex, r.ZEVMClient)
	require.NoError(r, err)
	r.ERC20ZRC20Addr = erc20zrc20AddrHex
	r.ERC20ZRC20 = erc20ZRC20

	erc20ERC20, err := erc20.NewERC20(erc20Addr, r.EVMClient)
	require.NoError(r, err)
	r.ERC20Addr = erc20Addr
	r.ERC20 = erc20ERC20

	// get balance
	balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.Account.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("ERC20 balance: %s", balance.String())

	// run deposit and withdraw ERC20 test
	txHash := r.LegacyDepositERC20WithAmountAndMessage(r.EVMAddress(), balance, []byte{})
	r.WaitForMinedCCTX(txHash)

	// approve 1 unit of the gas token to cover the gas fee
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ERC20ZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	tx = r.LegacyWithdrawERC20(balance)
	r.WaitForMinedCCTX(tx.Hash())
}
