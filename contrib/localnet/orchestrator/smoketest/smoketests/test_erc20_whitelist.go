package smoketests

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestWhitelistERC20(sm *runner.SmokeTestRunner) {
	utils.LoudPrintf("Test ERC20 whitelist\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	res, err := sm.ObserverClient.GetCoreParamsForChain(context.Background(), &observertypes.QueryGetCoreParamsForChainRequest{
		ChainId: int64(1337),
	})
	if err != nil {
		panic(err)
	}
	custodyAddr := ethcommon.HexToAddress(res.CoreParams.Erc20CustodyContractAddress)
	if custodyAddr == (ethcommon.Address{}) {
		panic("custody address is empty")
	}
	custody, err := erc20custody.NewERC20Custody(custodyAddr, sm.GoerliClient)
	if err != nil {
		panic(err)
	}
	iter, err := custody.FilterWhitelisted(&bind.FilterOpts{
		Start:   0,
		End:     nil,
		Context: context.Background(),
	}, []ethcommon.Address{})
	if err != nil {
		panic(err)
	}
	for iter.Next() {
		fmt.Printf("whitelisted: %s\n", iter.Event.Asset.Hex())
	}
}
