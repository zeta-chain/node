package e2etests

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// EVM2RPCURL is the RPC URL for the additional EVM localnet
// Only this test currently uses a additional EVM localnet, and this test is only run locally
// Therefore, we hardcode RPC urls and addresses for simplicity
const EVM2RPCURL = "http://eth2:8545"

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
		common.CoinType_Gas,
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

	// deploy erc20 zrc20
	res, err := newRunner.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		adminAddr,
		newRunner.ERC20Addr.Hex(),
		common.SepoliaChain().ChainId,
		18,
		"USDT",
		"USDT",
		common.CoinType_ERC20,
		100000,
	))
	if err != nil {
		panic(err)
	}

	// fetch the erc20 zrc20 contract address and remove the quotes
	erc20zrc20Addr, err := fetchAttribute(res, "Contract")
	if err != nil {
		panic(err)
	}
	if !ethcommon.IsHexAddress(erc20zrc20Addr) {
		panic(fmt.Errorf("invalid contract address: %s", erc20zrc20Addr))
	}

	erc20ZRC20, err := zrc20.NewZRC20(ethcommon.HexToAddress(erc20zrc20Addr), newRunner.ZEVMClient)
	if err != nil {
		panic(err)
	}
	newRunner.ERC20ZRC20 = erc20ZRC20

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

	// restart ZetaClient to pick up the new chain
	r.Logger.Print("🔄 restarting ZetaClient to pick up the new chain")
	if err := restartZetaClient(); err != nil {
		panic(err)
	}

	// wait 10 set for the chain to start
	time.Sleep(10 * time.Second)

	// test cross-chain functionalities on the new network
	// we use a Go routine to manually mine blocks because Anvil network only mine blocks on tx by default
	// we need automatic block mining to get the necessary confirmations for the cross-chain functionalities
	stopMining, err := newRunner.AnvilMineBlocks(EVM2RPCURL, 3)
	if err != nil {
		panic(err)
	}

	// deposit Ethers and ERC20 on ZetaChain
	txEtherDeposit := newRunner.DepositEther(false)
	newRunner.WaitForMinedCCTX(txEtherDeposit)
	txERC20Deposit := newRunner.DepositERC20()
	newRunner.WaitForMinedCCTX(txERC20Deposit)

	// perform withdrawals on the new chain
	TestZetaWithdraw(newRunner, []string{"10000000000000000000"})
	TestEtherWithdraw(newRunner, []string{"50000000000000000"})
	TestERC20Withdraw(r, []string{"100000"})

	// finally try to deposit Zeta back
	TestZetaDeposit(newRunner, []string{"100000000000000000"})

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
		r.EVMAuth,
		r.ZEVMAuth,
		r.BtcRPCClient,
		runner.NewLogger(false, color.FgHiYellow, "admin-evm2"),
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

// fetchAttribute fetches the attribute from the tx response
// TODO: use tx server
func fetchAttribute(res *sdktypes.TxResponse, key string) (string, error) {
	var logs []messageLog
	err := json.Unmarshal([]byte(res.RawLog), &logs)
	if err != nil {
		return "", err
	}

	var attributes []string
	for _, log := range logs {
		for _, event := range log.Events {
			for _, attr := range event.Attributes {
				attributes = append(attributes, attr.Key)
				if strings.EqualFold(attr.Key, key) {
					address := attr.Value

					if len(address) < 2 {
						return "", fmt.Errorf("invalid address: %s", address)
					}

					// trim the quotes
					address = address[1 : len(address)-1]

					return address, nil
				}
			}
		}
	}

	return "", fmt.Errorf("attribute %s not found, attributes:  %+v", key, attributes)
}

type messageLog struct {
	Events []event `json:"events"`
}

type event struct {
	Type       string      `json:"type"`
	Attributes []attribute `json:"attributes"`
}

type attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
