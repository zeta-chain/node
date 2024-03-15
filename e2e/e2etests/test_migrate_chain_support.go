package e2etests

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// EVM2RPCURL is the RPC URL for the additional EVM localnet
// Only this test currently uses a additional EVM localnet, and this test is only run locally
// Therefore, we hardcode RPC urls and addresses for simplicity
const EVM2RPCURL = "http://0.0.0.0:8546"

// EVM2ChainID is the chain ID for the additional EVM localnet
// We set Sepolia testnet although the value is not important, only used to differentiate
var EVM2ChainID = common.SepoliaChain().ChainId

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
	newRunner.SetupEVM(false)
	newRunner.MintERC20OnEvm(10000)

	// update the chain params to set up the chain

	// deposit Ethers and ERC20 on ZetaChain
	txEtherDeposit := newRunner.DepositEther(false)
	txERC20Deposit := newRunner.DepositERC20()
	newRunner.WaitForMinedCCTX(txEtherDeposit)
	newRunner.WaitForMinedCCTX(txERC20Deposit)

	// withdraw Zeta, Ethers and ERC20 to the new chain
	TestZetaWithdraw(r, []string{"1000000000000000000"})
	TestEtherWithdraw(r, []string{"10000000000000000"})
	TestERC20Withdraw(r, []string{"1000000000000000000"})
}

// configureEVM2 takes a runner and configures it to use the additional EVM localnet
func configureEVM2(r *runner.E2ERunner) (*runner.E2ERunner, error) {
	var newRunner runner.E2ERunner

	err := newRunner.CopyAddressesFrom(r)
	if err != nil {
		return nil, err
	}

	// All existing fields of the runner are kept except for the RPC URL and client for EVM
	client, auth, err := getEVMClient(newRunner.Ctx, EVM2RPCURL, r.DeployerPrivateKey)
	if err != nil {
		return nil, err
	}
	newRunner.EVMClient = client
	newRunner.EVMAuth = auth

	// reset evm contracts to ensure they are re-initialized
	newRunner.ZetaEthAddr = ethcommon.Address{}
	newRunner.ZetaEth = nil
	newRunner.ConnectorEthAddr = ethcommon.Address{}
	newRunner.ConnectorEth = nil
	newRunner.ERC20CustodyAddr = ethcommon.Address{}
	newRunner.ERC20Custody = nil
	newRunner.ERC20Addr = ethcommon.Address{}
	newRunner.ERC20 = nil

	return &newRunner, nil
}

// getEVMClient get evm client
// TODO: put this logic in common with the one in cmd/zetae2e/config/clients.go
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
