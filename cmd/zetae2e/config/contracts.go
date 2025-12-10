package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/protocol-contracts/pkg/coreregistry.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/wzeta.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zetaconnectornative.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/zetaconnectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zetaeth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/contracts/erc20"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	"github.com/zeta-chain/node/e2e/contracts/zevmswap"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/pkg/contracts/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/node/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/node/x/observer/types"
)

func chainParamsBySelector(
	chainParams []*types.ChainParams,
	selector func(chainID int64, additionalChains []chains.Chain) bool,
) *types.ChainParams {
	for _, chainParam := range chainParams {
		if selector(chainParam.ChainId, nil) {
			return chainParam
		}
	}
	return nil
}

func chainParamsByChainID(chainParams []*types.ChainParams, id int64) *types.ChainParams {
	for _, chainParam := range chainParams {
		if chainParam.ChainId == id {
			return chainParam
		}
	}
	return nil
}

func foreignCoinByChainID(
	foreignCoins []fungibletypes.ForeignCoins,
	id int64,
	coinType coin.CoinType,
) *fungibletypes.ForeignCoins {
	for _, fCoin := range foreignCoins {
		if fCoin.ForeignChainId == id && fCoin.CoinType == coinType {
			return &fCoin
		}
	}
	return nil
}

func setContractsGatewayEVM(r *runner.E2ERunner, params *types.ChainParams) error {
	r.GatewayEVMAddr = common.HexToAddress(params.GatewayAddress)
	if r.GatewayEVMAddr == (common.Address{}) {
		return nil
	}
	gatewayCode, err := r.EVMClient.CodeAt(r.Ctx, r.GatewayEVMAddr, nil)
	if err != nil || len(gatewayCode) == 0 {
		r.Logger.Print("❓ no code at EVM gateway address (%s)", r.GatewayEVMAddr)
		return nil
	}
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(r.GatewayEVMAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ZetaEthAddr = common.HexToAddress(params.ZetaTokenContractAddress)
	r.ZetaEth, err = zetaeth.NewZetaEth(r.ZetaEthAddr, r.EVMClient)
	if err != nil {
		return err
	}

	r.ConnectorEthAddr = common.HexToAddress(params.ConnectorContractAddress)
	r.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(r.ConnectorEthAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20CustodyAddr = common.HexToAddress(params.Erc20CustodyContractAddress)
	r.ERC20Custody, err = erc20custody.NewERC20Custody(r.ERC20CustodyAddr, r.EVMClient)
	if err != nil {
		return err
	}
	return nil
}

// setContractsFromConfig get EVM contracts from config
func setContractsFromConfig(r *runner.E2ERunner, conf config.Config) error {
	var err error

	chainParams, err := r.Clients.Zetacore.GetChainParams(r.Ctx)
	require.NoError(r, err, "get chain params")

	solChainParams := chainParamsBySelector(chainParams, chains.IsSolanaChain)

	// set Solana contracts
	if c := conf.Contracts.Solana.GatewayProgramID; c != "" {
		r.GatewayProgram = solana.MustPublicKeyFromBase58(c.String())
	} else if solChainParams != nil && solChainParams.GatewayAddress != "" {
		r.GatewayProgram = solana.MustPublicKeyFromBase58(solChainParams.GatewayAddress)
	}

	if c := conf.Contracts.Solana.SPLAddr; c != "" {
		r.SPLAddr = solana.MustPublicKeyFromBase58(c.String())
	}

	if c := conf.Contracts.Solana.ConnectedProgramID; c != "" {
		r.ConnectedProgram = solana.MustPublicKeyFromBase58(c.String())
	}

	if c := conf.Contracts.Solana.ConnectedSPLProgramID; c != "" {
		r.ConnectedSPLProgram = solana.MustPublicKeyFromBase58(c.String())
	}

	// set TON contracts
	if c := conf.Contracts.TON.GatewayAccountID; c != "" {
		r.TONGateway = ton.MustParseAccountID(c.String())
	}

	// set Sui contracts
	suiPackageID := conf.Contracts.Sui.GatewayPackageID
	suiGatewayID := conf.Contracts.Sui.GatewayObjectID

	if suiPackageID != "" && suiGatewayID != "" {
		r.SuiGateway = sui.NewGateway(suiPackageID.String(), suiGatewayID.String())
	}
	if c := conf.Contracts.Sui.GatewayUpgradeCap; c != "" {
		r.SuiGatewayUpgradeCap = c.String()
	}
	if c := conf.Contracts.Sui.FungibleTokenCoinType; c != "" {
		r.SuiTokenCoinType = c.String()
	}
	if c := conf.Contracts.Sui.FungibleTokenTreasuryCap; c != "" {
		r.SuiTokenTreasuryCap = c.String()
	}
	r.SuiExample = conf.Contracts.Sui.Example

	// set EVM contracts
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err, "get evm chain ID")
	evmChainParams := chainParamsByChainID(chainParams, evmChainID.Int64())

	if evmChainParams == nil {
		return fmt.Errorf("no EVM chain params found for chain ID %d", evmChainID.Int64())
	}

	err = setContractsGatewayEVM(r, evmChainParams)
	if err != nil {
		return fmt.Errorf("setContractsGatewayEVM: %w", err)
	}

	if c := conf.Contracts.EVM.ERC20; c != "" {
		r.ERC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ERC20: %w", err)
		}
		r.ERC20, err = erc20.NewERC20(r.ERC20Addr, r.EVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.EVM.ConnectorNative; c != "" {
		r.ConnectorNativeAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ConnectorNativeAddr: %w", err)
		}
		r.ConnectorNative, err = zetaconnectornative.NewZetaConnectorNative(r.ConnectorNativeAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}
	// Overwrite using contract addresses from config
	if c := conf.Contracts.EVM.ConnectorEthAddr; c != "" {
		r.ConnectorEthAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ConnectorEthAddr: %w", err)
		}
		r.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(r.ConnectorEthAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}
	// set ZEVM contracts
	foreignCoins, err := r.Clients.Zetacore.Fungible.ForeignCoinsAll(
		r.Ctx,
		&fungibletypes.QueryAllForeignCoinsRequest{},
	)
	if err != nil {
		return err
	}

	ethForeignCoin := foreignCoinByChainID(foreignCoins.ForeignCoins, evmChainID.Int64(), coin.CoinType_Gas)
	if ethForeignCoin != nil {
		r.ETHZRC20Addr = common.HexToAddress(ethForeignCoin.Zrc20ContractAddress)
		r.ETHZRC20, err = zrc20.NewZRC20(r.ETHZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.SystemContractAddr; c != "" {
		r.SystemContractAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SystemContractAddr: %w", err)
		}
		r.SystemContract, err = systemcontract.NewSystemContract(r.SystemContractAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.ERC20ZRC20Addr; c != "" {
		r.ERC20ZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ERC20ZRC20Addr: %w", err)
		}
		r.ERC20ZRC20, err = zrc20.NewZRC20(r.ERC20ZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.BTCZRC20Addr; c != "" {
		r.BTCZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid BTCZRC20Addr: %w", err)
		}
		r.BTCZRC20, err = zrc20.NewZRC20(r.BTCZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.SOLZRC20Addr; c != "" {
		r.SOLZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SOLZRC20Addr: %w", err)
		}
		r.SOLZRC20, err = zrc20.NewZRC20(r.SOLZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.SPLZRC20Addr; c != "" {
		r.SPLZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SPLZRC20Addr: %w", err)
		}
		r.SPLZRC20, err = zrc20.NewZRC20(r.SPLZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.TONZRC20Addr; c != "" {
		r.TONZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid TONZRC20Addr: %w", err)
		}
		r.TONZRC20, err = zrc20.NewZRC20(r.TONZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.SUIZRC20Addr; c != "" {
		r.SUIZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SUIZRC20Addr: %w", err)
		}
		r.SUIZRC20, err = zrc20.NewZRC20(r.SUIZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.SuiTokenZRC20Addr; c != "" {
		r.SuiTokenZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SuiTokenZRC20Addr: %w", err)
		}
		r.SuiTokenZRC20, err = zrc20.NewZRC20(r.SuiTokenZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.UniswapFactoryAddr; c != "" {
		r.UniswapV2FactoryAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid UniswapFactoryAddr: %w", err)
		}
		r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.UniswapRouterAddr; c != "" {
		r.UniswapV2RouterAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid UniswapRouterAddr: %w", err)
		}
		r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.ConnectorZEVMAddr; c != "" {
		r.ConnectorZEVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ConnectorZEVMAddr: %w", err)
		}
		r.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(r.ConnectorZEVMAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.WZetaAddr; c != "" {
		r.WZetaAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid WZetaAddr: %w", err)
		}
		r.WZeta, err = wzeta.NewWETH9(r.WZetaAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.ZEVMSwapAppAddr; c != "" {
		r.ZEVMSwapAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ZEVMSwapAppAddr: %w", err)
		}
		r.ZEVMSwapApp, err = zevmswap.NewZEVMSwapApp(r.ZEVMSwapAppAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.EVM.TestDappAddr; c != "" {
		r.EvmTestDAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid EvmTestDappAddr: %w", err)
		}
	}

	if c := conf.Contracts.ZEVM.TestDappAddr; c != "" {
		r.ZevmTestDAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ZevmTestDappAddr: %w", err)
		}
	}

	if c := conf.Contracts.EVM.TestDAppV2Addr; c != "" {
		r.TestDAppV2EVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid TestDAppV2Addr: %w", err)
		}
		r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2EVMAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.Gateway; c != "" {
		r.GatewayZEVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid GatewayAddr: %w", err)
		}
		r.GatewayZEVM, err = gatewayzevm.NewGatewayZEVM(r.GatewayZEVMAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.TestDAppV2Addr; c != "" {
		r.TestDAppV2ZEVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid TestDAppV2Addr: %w", err)
		}
		r.TestDAppV2ZEVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2ZEVMAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.CoreRegistry; c != "" {
		r.CoreRegistryAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid CoreRegistryAddr: %w", err)
		}
		r.CoreRegistry, err = coreregistry.NewCoreRegistry(r.CoreRegistryAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	return nil
}
