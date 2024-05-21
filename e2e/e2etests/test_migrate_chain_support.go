package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/txserver"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// EVM2RPCURL is the RPC URL for the additional EVM localnet
// Only this test currently uses a additional EVM localnet, and this test is only run locally
// Therefore, we hardcode RPC urls and addresses for simplicity
const EVM2RPCURL = "http://eth2:8545"

// EVM2ChainID is the chain ID for the additional EVM localnet
// We set Sepolia testnet although the value is not important, only used to differentiate
var EVM2ChainID = chains.SepoliaChain.ChainId

func TestMigrateChainSupport(r *runner.E2ERunner, _ []string) {
	// deposit most of the ZETA supply on ZetaChain
	zetaAmount := big.NewInt(1e18)
	zetaAmount = zetaAmount.Mul(zetaAmount, big.NewInt(20_000_000_000)) // 20B Zeta
	r.DepositZetaWithAmount(r.DeployerAddress, zetaAmount)

	// do an ethers withdraw on the previous chain (0.01eth) for some interaction
	TestEtherWithdraw(r, []string{"10000000000000000"})

	// create runner for the new EVM and set it up
	newRunner, err := configureEVM2(r)
	if err != nil {
		panic(err)
	}
	newRunner.SetupEVM(false, false)

	// mint some ERC20
	newRunner.MintERC20OnEvm(10000)

	// we deploy connectorETH in this test to simulate a new "canonical" chain emitting ZETA
	// to represent the ZETA already existing on ZetaChain we manually send the minted ZETA to the connector
	newRunner.SendZetaOnEvm(newRunner.ConnectorEthAddr, 20_000_000_000)

	// update the chain params to set up the chain
	chainParams := getNewEVMChainParams(newRunner)
	adminAddr, err := newRunner.ZetaTxServer.GetAccountAddressFromName(utils.FungibleAdminName)
	if err != nil {
		panic(err)
	}
	_, err = newRunner.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, observertypes.NewMsgUpdateChainParams(
		adminAddr,
		chainParams,
	))
	if err != nil {
		panic(err)
	}

	// setup the gas token
	if err != nil {
		panic(err)
	}
	_, err = newRunner.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		adminAddr,
		"",
		chainParams.ChainId,
		18,
		"Sepolia ETH",
		"sETH",
		coin.CoinType_Gas,
		100000,
	))
	if err != nil {
		panic(err)
	}

	// set the gas token in the runner
	ethZRC20Addr, err := newRunner.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(chainParams.ChainId))
	if err != nil {
		panic(err)
	}
	if (ethZRC20Addr == ethcommon.Address{}) {
		panic("eth zrc20 not found")
	}
	newRunner.ETHZRC20Addr = ethZRC20Addr
	ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, newRunner.ZEVMClient)
	if err != nil {
		panic(err)
	}
	newRunner.ETHZRC20 = ethZRC20

	// set the chain nonces for the new chain
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, observertypes.NewMsgResetChainNonces(
		adminAddr,
		chainParams.ChainId,
		0,
		0,
	))
	if err != nil {
		panic(err)
	}

	// deactivate the previous chain
	chainParams = observertypes.GetDefaultGoerliLocalnetChainParams()
	chainParams.IsSupported = false
	_, err = newRunner.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, observertypes.NewMsgUpdateChainParams(
		adminAddr,
		chainParams,
	))
	if err != nil {
		panic(err)
	}

	// restart ZetaClient to pick up the new chain
	r.Logger.Print("🔄 restarting ZetaClient to pick up the new chain")
	if err := restartZetaClient(); err != nil {
		panic(err)
	}

	// wait 10 set for the chain to start
	time.Sleep(10 * time.Second)

	// emitting a withdraw with the previous chain should fail
	txWithdraw, err := r.ETHZRC20.Withdraw(r.ZEVMAuth, r.DeployerAddress.Bytes(), big.NewInt(10000000000000000))
	if err == nil {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txWithdraw, r.Logger, r.ReceiptTimeout)
		if receipt.Status == 1 {
			panic("withdraw should have failed on the previous chain")
		}
	}

	// test cross-chain functionalities on the new network
	// we use a Go routine to manually mine blocks because Anvil network only mine blocks on tx by default
	// we need automatic block mining to get the necessary confirmations for the cross-chain functionalities
	stopMining, err := newRunner.AnvilMineBlocks(EVM2RPCURL, 3)
	if err != nil {
		panic(err)
	}

	// deposit Ethers and ERC20 on ZetaChain
	etherAmount := big.NewInt(1e18)
	etherAmount = etherAmount.Mul(etherAmount, big.NewInt(10))
	txEtherDeposit := newRunner.DepositEtherWithAmount(false, etherAmount)
	newRunner.WaitForMinedCCTX(txEtherDeposit)

	// perform withdrawals on the new chain
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(10))
	newRunner.DepositAndApproveWZeta(amount)
	tx := newRunner.WithdrawZeta(amount, true)
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			crosschaintypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		))
	}

	TestEtherWithdraw(newRunner, []string{"50000000000000000"})

	// finally try to deposit Zeta back
	TestZetaDeposit(newRunner, []string{"100000000000000000"})

	// ERC20 test

	// whitelist erc20 zrc20
	newRunner.Logger.Info("whitelisting ERC20 on new network")
	res, err := newRunner.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, crosschaintypes.NewMsgWhitelistERC20(
		adminAddr,
		newRunner.ERC20Addr.Hex(),
		chains.SepoliaChain.ChainId,
		"USDT",
		"USDT",
		18,
		100000,
	))
	if err != nil {
		panic(err)
	}

	// retrieve zrc20 and cctx from event
	whitelistCCTXIndex, err := txserver.FetchAttributeFromTxResponse(res, "whitelist_cctx_index")
	if err != nil {
		panic(err)
	}

	erc20zrc20Addr, err := txserver.FetchAttributeFromTxResponse(res, "zrc20_address")
	if err != nil {
		panic(err)
	}

	// wait for the whitelist cctx to be mined
	newRunner.WaitForMinedCCTXFromIndex(whitelistCCTXIndex)

	// set erc20 zrc20 contract address
	if !ethcommon.IsHexAddress(erc20zrc20Addr) {
		panic(fmt.Errorf("invalid contract address: %s", erc20zrc20Addr))
	}
	erc20ZRC20, err := zrc20.NewZRC20(ethcommon.HexToAddress(erc20zrc20Addr), newRunner.ZEVMClient)
	if err != nil {
		panic(err)
	}
	newRunner.ERC20ZRC20 = erc20ZRC20

	// deposit ERC20 on ZetaChain
	txERC20Deposit := newRunner.DepositERC20()
	newRunner.WaitForMinedCCTX(txERC20Deposit)

	// stop mining
	stopMining()
}

// configureEVM2 takes a runner and configures it to use the additional EVM localnet
func configureEVM2(r *runner.E2ERunner) (*runner.E2ERunner, error) {
	// initialize a new runner with previous runner values
	newRunner := runner.NewE2ERunner(
		r.Ctx,
		"admin-evm2",
		r.CtxCancel,
		r.DeployerAddress,
		r.DeployerPrivateKey,
		r.FungibleAdminMnemonic,
		r.EVMClient,
		r.ZEVMClient,
		r.CctxClient,
		r.ZetaTxServer,
		r.FungibleClient,
		r.AuthClient,
		r.BankClient,
		r.ObserverClient,
		r.LightclientClient,
		r.EVMAuth,
		r.ZEVMAuth,
		r.BtcRPCClient,
		runner.NewLogger(true, color.FgHiYellow, "admin-evm2"),
	)

	// All existing fields of the runner are the same except for the RPC URL and client for EVM
	ewvmClient, evmAuth, err := getEVMClient(newRunner.Ctx, EVM2RPCURL, r.DeployerPrivateKey)
	if err != nil {
		return nil, err
	}
	newRunner.EVMClient = ewvmClient
	newRunner.EVMAuth = evmAuth

	// Copy the ZetaChain contract addresses from the original runner
	if err := newRunner.CopyAddressesFrom(r); err != nil {
		return nil, err
	}

	// reset evm contracts to ensure they are re-initialized
	newRunner.ZetaEthAddr = ethcommon.Address{}
	newRunner.ZetaEth = nil
	newRunner.ConnectorEthAddr = ethcommon.Address{}
	newRunner.ConnectorEth = nil
	newRunner.ERC20CustodyAddr = ethcommon.Address{}
	newRunner.ERC20Custody = nil
	newRunner.ERC20Addr = ethcommon.Address{}
	newRunner.ERC20 = nil

	return newRunner, nil
}

// getEVMClient get evm client from rpc and private key
func getEVMClient(ctx context.Context, rpc, privKey string) (*ethclient.Client, *bind.TransactOpts, error) {
	evmClient, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, nil, err
	}

	chainid, err := evmClient.ChainID(ctx)
	if err != nil {
		return nil, nil, err
	}
	deployerPrivkey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, nil, err
	}
	evmAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		return nil, nil, err
	}

	return evmClient, evmAuth, nil
}

// getNewEVMChainParams returns the chain params for the new EVM chain
func getNewEVMChainParams(r *runner.E2ERunner) *observertypes.ChainParams {
	// goerli local as base
	chainParams := observertypes.GetDefaultGoerliLocalnetChainParams()

	// set the chain id to the new chain id
	chainParams.ChainId = EVM2ChainID

	// set contracts
	chainParams.ConnectorContractAddress = r.ConnectorEthAddr.Hex()
	chainParams.Erc20CustodyContractAddress = r.ERC20CustodyAddr.Hex()
	chainParams.ZetaTokenContractAddress = r.ZetaEthAddr.Hex()

	// set supported
	chainParams.IsSupported = true

	return chainParams
}

// restartZetaClient restarts the Zeta client
func restartZetaClient() error {
	sshCommandFilePath := "/work/restart-zetaclientd.sh"
	cmd := exec.Command("/bin/sh", sshCommandFilePath)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error restarting ZetaClient: %s - %s", err.Error(), output)
	}
	return nil
}
