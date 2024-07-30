package runner

import (
	"math/big"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/prototypes/evm/erc20custodynew.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/prototypes/evm/gatewayevm.sol"

	"github.com/zeta-chain/zetacore/e2e/contracts/erc1967proxy"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/constant"
)

// SetupEVMV2 setup contracts on EVM with v2 contracts
func (r *E2ERunner) SetupEVMV2() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	r.Logger.Print("⚙️ setting up EVM v2 network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("EVM v2 setup took %s\n", time.Since(startTime))
	}()

	r.Logger.InfoLoud("Deploy Gateway and ERC20Custody ERC20\n")

	// donate to the TSS address to avoid account errors because deploying gas token ZRC20 will automatically mint
	// gas token on ZetaChain to initialize the pool
	txDonation, err := r.SendEther(r.TSSAddress, big.NewInt(101000000000000000), []byte(constant.DonationMessage))
	require.NoError(r, err)

	r.Logger.Info("Deploying Gateway EVM")
	gatewayEVMAddr, txGateway, _, err := gatewayevm.DeployGatewayEVM(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txGateway, "Gateway deployment failed")

	gatewayEVMABI, err := gatewayevm.GatewayEVMMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := gatewayEVMABI.Pack("initialize", r.TSSAddress, r.ZetaEthAddr)
	require.NoError(r, err)

	// Deploy the proxy contract
	proxyAddress, txProxy, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		gatewayEVMAddr,
		initializerData,
	)
	require.NoError(r, err)

	r.GatewayEVMAddr = proxyAddress
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(proxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway EVM contract address: %s, tx hash: %s", gatewayEVMAddr.Hex(), txGateway.Hash().Hex())

	r.Logger.Info("Deploying ERC20Custody contract")
	erc20CustodyNewAddr, txCustody, erc20CustodyNew, err := erc20custodynew.DeployERC20CustodyNew(
		r.EVMAuth,
		r.EVMClient,
		r.GatewayEVMAddr,
		r.TSSAddress,
	)
	require.NoError(r, err)

	r.ERC20CustodyNewAddr = erc20CustodyNewAddr
	r.ERC20CustodyNew = erc20CustodyNew
	r.Logger.Info(
		"ERC20CustodyNew contract address: %s, tx hash: %s",
		erc20CustodyNewAddr.Hex(),
		txCustody.Hash().Hex(),
	)

	// check contract deployment receipt
	ensureTxReceipt(txDonation, "EVM donation tx failed")
	ensureTxReceipt(txCustody, "ERC20CustodyNew deployment failed")
	ensureTxReceipt(txProxy, "Gateway proxy deployment failed")
}
