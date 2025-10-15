package legacy

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/example"
	testcontract "github.com/zeta-chain/node/e2e/contracts/reverter"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// TestEtherDepositAndCall tests deposit of ethers calling a example contract
func TestEtherDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount
	value := utils.ParseBigInt(r, args[0])

	r.Logger.Info("Deploying example contract")
	exampleAddr, txDeploy, exampleContract, err := example.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	r.Logger.Info("Example contract deployed")

	// preparing tx
	evmClient := r.EVMClient
	gasLimit := uint64(23000)
	gasPrice, err := evmClient.SuggestGasPrice(r.Ctx)
	require.NoError(r, err)

	nonce, err := evmClient.PendingNonceAt(r.Ctx, r.EVMAddress())
	require.NoError(r, err)

	data := append(exampleAddr.Bytes(), []byte("hello sailors")...)
	tx := ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := evmClient.NetworkID(r.Ctx)
	require.NoError(r, err)

	deployerPrivkey, err := r.Account.PrivateKey()
	require.NoError(r, err)

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	require.NoError(r, err)

	r.Logger.Info("Sending a cross-chain call to example contract")
	err = evmClient.SendTransaction(r.Ctx, signedTx)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)

	// Checking example contract has been called, bar value should be set to amount
	utils.WaitAndVerifyExampleContractCall(r, exampleContract, value, r.EVMAddress().Bytes())
	r.Logger.Info("Cross-chain call succeeded")

	r.Logger.Info("Deploying reverter contract")
	reverterAddr, txDeploy, _, err := testcontract.DeployReverter(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	r.Logger.Info("Example reverter deployed")

	// preparing tx for reverter
	gasPrice, err = evmClient.SuggestGasPrice(r.Ctx)
	require.NoError(r, err)

	nonce, err = evmClient.PendingNonceAt(r.Ctx, r.EVMAddress())
	require.NoError(r, err)

	data = append(reverterAddr.Bytes(), []byte("hello sailors")...)
	tx = ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	signedTx, err = ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	require.NoError(r, err)

	r.Logger.Info("Sending a cross-chain call to reverter contract")
	err = evmClient.SendTransaction(r.Ctx, signedTx)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

	r.Logger.Info("Cross-chain call to reverter reverted")

	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo)
}
