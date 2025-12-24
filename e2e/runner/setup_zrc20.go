package runner

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/pkg/chains"
)

// SetupZRC20 setup ZRC20 for the ZEVM
func (r *E2ERunner) SetupZRC20(zrc20Deployment txserver.ZRC20Deployment) {
	r.Logger.Print("⚙️ deploying ZRC20s on ZEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("System contract deployments took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and ZRC20 contracts on ZetaChain
	deployedZRC20Addresses, err := r.ZetaTxServer.DeployZRC20s(
		zrc20Deployment,
		r.skipChainOperations,
	)
	require.NoError(r, err)

	// Set ERC20ZRC20Addr
	r.ERC20ZRC20Addr = deployedZRC20Addresses.ERC20ZRC20Addr
	r.ERC20ZRC20, err = zrc20.NewZRC20(r.ERC20ZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	// Set SPLZRC20Addr if set
	if deployedZRC20Addresses.SPLZRC20Addr != (ethcommon.Address{}) {
		r.SPLZRC20Addr = deployedZRC20Addresses.SPLZRC20Addr
		r.SPLZRC20, err = zrc20.NewZRC20(r.SPLZRC20Addr, r.ZEVMClient)
		require.NoError(r, err)
	}

	// set ZRC20 contracts
	r.SetupETHZRC20()
	r.SetupBTCZRC20()
	r.SetupSOLZRC20()
	r.SetupTONZRC20()
	r.ActivateChainsOnRegistry()
}

func (r *E2ERunner) ActivateChainsOnRegistry() {
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)
	// Activate ETH
	_, err = r.CoreRegistry.ChangeChainStatus(r.ZEVMAuth, evmChainID, r.ETHZRC20Addr, []byte{}, true)
	require.NoError(r, err)
}

// SetupETHZRC20 sets up the ETH ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupETHZRC20() {
	ethZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.GoerliLocalnet.ChainId),
	)
	require.NoError(r, err)
	require.NotEqual(r, ethcommon.Address{}, ethZRC20Addr, "eth zrc20 not found")

	r.ETHZRC20Addr = ethZRC20Addr
	ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	r.ETHZRC20 = ethZRC20
}

// SetupBTCZRC20 sets up the BTC ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupBTCZRC20() {
	BTCZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.BitcoinRegtest.ChainId),
	)
	require.NoError(r, err)
	r.BTCZRC20Addr = BTCZRC20Addr
	r.Logger.Info("BTCZRC20Addr: %s", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)
	r.BTCZRC20 = BTCZRC20
}

// SetupSOLZRC20 sets up the SOL ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupSOLZRC20() {
	// set SOLZRC20 address by chain ID
	SOLZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.SolanaLocalnet.ChainId),
	)
	require.NoError(r, err)

	// set SOLZRC20 address
	r.SOLZRC20Addr = SOLZRC20Addr
	r.Logger.Info("SOLZRC20Addr: %s", SOLZRC20Addr.Hex())

	// set SOLZRC20 contract
	SOLZRC20, err := zrc20.NewZRC20(SOLZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)
	r.SOLZRC20 = SOLZRC20
}

// SetupTONZRC20 sets up the TON ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupTONZRC20() {
	chainID := chains.TONLocalnet.ChainId

	// noop
	if r.skipChainOperations(chainID) {
		return
	}

	TONZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(chainID))
	require.NoError(r, err)

	r.TONZRC20Addr = TONZRC20Addr
	r.Logger.Info("TON ZRC20 address: %s", TONZRC20Addr.Hex())

	TONZRC20, err := zrc20.NewZRC20(TONZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	r.TONZRC20 = TONZRC20
}

// SetupSUIZRC20 sets up the SUI ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupSUIZRC20() {
	chainID := chains.SuiLocalnet.ChainId

	// noop
	if r.skipChainOperations(chainID) {
		return
	}

	SUIZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(chainID))
	require.NoError(r, err)

	r.SUIZRC20Addr = SUIZRC20Addr
	r.Logger.Info("SUI ZRC20 address: %s", SUIZRC20Addr.Hex())

	SUIZRC20, err := zrc20.NewZRC20(SUIZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	r.SUIZRC20 = SUIZRC20
}
