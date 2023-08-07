//go:build PRIVNET
// +build PRIVNET

package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/zeta-chain/zetacore/x/crosschain/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (sm *SmokeTest) TestCrosschainSwap() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Testing Bitcoin ERC20 crosschain swap...\n")
	// Firstly, deposit 1.15 BTC into Zeta for liquidity
	//sm.DepositBTC()
	// Secondly, deposit 1000.0 USDT into Zeta for liquidity
	LoudPrintf("Depositing 1000 USDT & 1.15 BTC for liquidity\n")

	txhash := sm.DepositERC20(big.NewInt(1e9), []byte{})
	WaitCctxMinedByInTxHash(txhash.Hex(), sm.cctxClient)

	sm.zevmAuth.GasLimit = 10000000
	tx, err := sm.UniswapV2Factory.CreatePair(sm.zevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
	usdtBtcPair, err := sm.UniswapV2Factory.GetPair(&bind.CallOpts{}, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT-BTC pair receipt txhash %s status %d pair addr %s\n", receipt.TxHash, receipt.Status, usdtBtcPair.Hex())

	tx, err = sm.USDTZRC20.Approve(sm.zevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("USDT ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	tx, err = sm.BTCZRC20.Approve(sm.zevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("BTC ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	// Add 100 USDT liq and 0.001 BTC
	bal, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on BTC ZRC20: %d\n", bal)
	bal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	tx, err = sm.UniswapV2Router.AddLiquidity(sm.zevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr, big.NewInt(1e8), big.NewInt(1e8), big.NewInt(1e8), big.NewInt(1e5), DeployerAddress, big.NewInt(time.Now().Add(10*time.Minute).Unix()))
	if err != nil {
		fmt.Printf("Error liq %s", err.Error())
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("Add liquidity receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	fmt.Printf("Funding contracts ZEVMSwapApp with gas ZRC20s; 1e7 ETH, 1e6 BTC\n")
	// Fund ZEVMSwapApp with gas ZRC20s
	tx, err = sm.ETHZRC20.Transfer(sm.zevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e7))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("  USDT ZRC20 transfer receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	bal1, _ := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.ZEVMSwapAppAddr)
	fmt.Printf("  ZEVMSwapApp ETHZRC20 balance %d", bal1)
	tx, err = sm.BTCZRC20.Transfer(sm.zevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e6))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("  BTC ZRC20 transfer receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	bal2, _ := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.ZEVMSwapAppAddr)
	fmt.Printf("  ZEVMSwapApp BTCZRC20 balance %d", bal2)

	// msg would be [ZEVMSwapAppAddr, memobytes]
	// memobytes is dApp specific; see the contracts/ZEVMSwapApp.sol for details
	msg := []byte{}
	msg = append(msg, sm.ZEVMSwapAppAddr.Bytes()...)
	memobytes, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.BTCZRC20Addr, []byte(BTCDeployerAddress.EncodeAddress()))

	if err != nil {
		panic(err)
	}
	fmt.Printf("memobytes(%d) %x\n", len(memobytes), memobytes)
	msg = append(msg, memobytes...)

	fmt.Printf("***** First test: USDT -> BTC\n")
	// Should deposit USDT for swap, swap for BTC and withdraw BTC
	txhash = sm.DepositERC20(big.NewInt(8e7), msg)
	cctx1 := WaitCctxMinedByInTxHash(txhash.Hex(), sm.cctxClient)

	_, err = sm.btcRPCClient.GenerateToAddress(10, BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	// cctx1 index acts like the inTxHash for the second cctx (the one that withdraws BTC)
	cctx2 := WaitCctxMinedByInTxHash(cctx1.Index, sm.cctxClient)
	_ = cctx2
	fmt.Printf("cctx2 outbound tx hash %s\n", cctx2.GetCurrentOutTxParam().OutboundTxHash)

	fmt.Printf("******* Second test: BTC -> USDT\n")
	utxos, err := sm.btcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	fmt.Printf("#utxos %d\n", len(utxos))
	//fmt.Printf("Unimplemented!\n")
	memo, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.USDTZRC20Addr, DeployerAddress.Bytes())
	if err != nil {
		panic(err)
	}
	memo = append(sm.ZEVMSwapAppAddr.Bytes(), memo...)
	fmt.Printf("memo length %d\n", len(memo))

	txid, err := SendToTSSFromDeployerWithMemo(BTCTSSAddress, 0.001, utxos[0:2], sm.btcRPCClient, memo)
	fmt.Printf("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation\n", txid)
	_, err = sm.btcRPCClient.GenerateToAddress(10, BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}

	cctx3 := WaitCctxMinedByInTxHash(txid.String(), sm.cctxClient)
	fmt.Printf("cctx3 index %s\n", cctx3.Index)
	fmt.Printf("  inboudn tx hash %s\n", cctx3.InboundTxParams.InboundTxObservedHash)
	fmt.Printf("  status %s\n", cctx3.CctxStatus.Status.String())
	fmt.Printf("  status msg: %s\n", cctx3.CctxStatus.StatusMessage)

	cctx4 := WaitCctxMinedByInTxHash(cctx3.Index, sm.cctxClient)
	fmt.Printf("cctx4 index %s\n", cctx4.Index)
	fmt.Printf("  outbound tx hash %s\n", cctx4.GetCurrentOutTxParam().OutboundTxHash)
	fmt.Printf("  status %s\n", cctx4.CctxStatus.Status.String())

	{
		fmt.Printf("******* Third test: BTC -> ETH with contract call reverted; should refund BTC\n")
		utxos, err := sm.btcRPCClient.ListUnspent()
		if err != nil {
			panic(err)
		}
		fmt.Printf("#utxos %d\n", len(utxos))
		// the following memo will result in a revert in the contract call as targetZRC20 is set to DeployerAddress
		// which is apparently not a ZRC20 contract; the UNISWAP call will revert
		memo, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, DeployerAddress, DeployerAddress.Bytes())
		if err != nil {
			panic(err)
		}
		memo = append(sm.ZEVMSwapAppAddr.Bytes(), memo...)
		fmt.Printf("memo length %d\n", len(memo))

		txid, err := SendToTSSFromDeployerWithMemo(BTCTSSAddress, 0.001, utxos[0:2], sm.btcRPCClient, memo)
		fmt.Printf("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation\n", txid)
		_, err = sm.btcRPCClient.GenerateToAddress(10, BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}

		cctx := WaitCctxMinedByInTxHash(txid.String(), sm.cctxClient)
		fmt.Printf("cctx3 index http://localhost:1317/zeta-chain/crosschain/cctx/%s\n", cctx.Index)
		fmt.Printf("  inboudn tx hash %s\n", cctx.InboundTxParams.InboundTxObservedHash)
		fmt.Printf("  status %s\n", cctx.CctxStatus.Status.String())
		fmt.Printf("  status msg: %s\n", cctx.CctxStatus.StatusMessage)
		if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
			panic(fmt.Sprintf("expected reverted status; got %s", cctx.CctxStatus.Status.String()))
		}
		outTxHash, err := chainhash.NewHashFromStr(cctx.GetCurrentOutTxParam().OutboundTxHash)
		if err != nil {
			panic(err)
		}
		txraw, err := sm.btcRPCClient.GetRawTransactionVerbose(outTxHash)
		if err != nil {
			panic(err)
		}
		fmt.Printf("out txid %s\n", txraw.Txid)
		for _, vout := range txraw.Vout {
			fmt.Printf("  vout %d\n", vout.N)
			fmt.Printf("  value %f\n", vout.Value)
			fmt.Printf("  scriptPubKey %s\n", vout.ScriptPubKey.Hex)
			fmt.Printf("  p2wpkh address: %s\n", ScriptPKToAddress(vout.ScriptPubKey.Hex))
		}
	}
}
