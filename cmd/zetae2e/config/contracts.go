package config

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/zevm/wzeta.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/zevm/zetaconnectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/v1/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/contracts/contextapp"
	"github.com/zeta-chain/node/e2e/contracts/erc20"
	"github.com/zeta-chain/node/e2e/contracts/zevmswap"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/pkg/contracts/testdappv2"
)

// setContractsFromConfig get EVM contracts from config
func setContractsFromConfig(r *runner.E2ERunner, conf config.Config) error {
	var err error

	// set Solana contracts
	if c := conf.Contracts.Solana.GatewayProgramID; c != "" {
		r.GatewayProgram = solana.MustPublicKeyFromBase58(c)
	}

	// set EVM contracts
	if c := conf.Contracts.EVM.ZetaEthAddr; c != "" {
		r.ZetaEthAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ZetaEthAddr: %w", err)
		}
		r.ZetaEth, err = zetaeth.NewZetaEth(r.ZetaEthAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}

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

	if c := conf.Contracts.EVM.CustodyAddr; c != "" {
		r.ERC20CustodyAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid CustodyAddr: %w", err)
		}
		r.ERC20Custody, err = erc20custody.NewERC20Custody(r.ERC20CustodyAddr, r.EVMClient)
		if err != nil {
			return err
		}
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

	// set ZEVM contracts
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

	if c := conf.Contracts.ZEVM.ETHZRC20Addr; c != "" {
		r.ETHZRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ETHZRC20Addr: %w", err)
		}
		r.ETHZRC20, err = zrc20.NewZRC20(r.ETHZRC20Addr, r.ZEVMClient)
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

	if c := conf.Contracts.ZEVM.ContextAppAddr; c != "" {
		r.ContextAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ContextAppAddr: %w", err)
		}
		r.ContextApp, err = contextapp.NewContextApp(r.ContextAppAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.ZEVM.TestDappAddr; c != "" {
		r.ZevmTestDAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ZevmTestDappAddr: %w", err)
		}
	}

	if c := conf.Contracts.EVM.TestDappAddr; c != "" {
		r.EvmTestDAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid EvmTestDappAddr: %w", err)
		}
	}

	// v2 contracts

	if c := conf.Contracts.EVM.Gateway; c != "" {
		r.GatewayEVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid GatewayAddr: %w", err)
		}
		r.GatewayEVM, err = gatewayevm.NewGatewayEVM(r.GatewayEVMAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.EVM.ERC20CustodyNew; c != "" {
		r.ERC20CustodyV2Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ERC20CustodyV2Addr: %w", err)
		}
		r.ERC20CustodyV2, err = erc20custodyv2.NewERC20Custody(r.ERC20CustodyV2Addr, r.EVMClient)
		if err != nil {
			return err
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
		r.TestDAppV2ZEVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2ZEVMAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}

	return nil
}
