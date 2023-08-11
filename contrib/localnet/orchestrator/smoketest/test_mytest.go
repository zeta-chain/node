//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (sm *SmokeTest) TestMyTest() {
	LoudPrintf("Test ERC20 whitelist\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	res, err := sm.observerClient.GetCoreParamsForChain(context.Background(), &observertypes.QueryGetCoreParamsForChainRequest{
		ChainID: int64(1337),
	})
	if err != nil {
		panic(err)
	}
	custodyAddr := ethcommon.HexToAddress(res.CoreParams.Erc20CustodyContractAddress)
	if custodyAddr == (ethcommon.Address{}) {
		panic("custody address is empty")
	}
	custody, err := erc20custody.NewERC20Custody(custodyAddr, sm.goerliClient)
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
