package config

import (
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/erc20"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

// getEVMContractsFromConfig get EVM contracts from config
func setEVMContractsFromConfig(r *runner.SmokeTestRunner, conf config.Config) error {
	var err error

	// set EVM contracts
	if c := conf.Contracts.EVM.ZetaEthAddress; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ZetaEthAddress: %s", c)
		}
		r.ZetaEthAddr = ethcommon.HexToAddress(c)
		r.ZetaEth, err = zetaeth.NewZetaEth(r.ZetaEthAddr, r.GoerliClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.EVM.ConnectorEthAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ConnectorEthAddr: %s", c)
		}
		r.ConnectorEthAddr = ethcommon.HexToAddress(c)
		r.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(r.ConnectorEthAddr, r.GoerliClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.EVM.CustodyAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid CustodyAddr: %s", c)
		}
		r.ERC20CustodyAddr = ethcommon.HexToAddress(c)
		r.ERC20Custody, err = erc20custody.NewERC20Custody(r.ERC20CustodyAddr, r.GoerliClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.EVM.USDT; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid USDT: %s", c)
		}
		r.USDTERC20Addr = ethcommon.HexToAddress(c)
		r.USDTERC20, err = erc20.NewUSDT(r.USDTERC20Addr, r.GoerliClient)
		if err != nil {
			return err
		}
	}

	return nil
}
