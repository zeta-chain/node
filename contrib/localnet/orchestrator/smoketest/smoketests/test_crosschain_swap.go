package smoketests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestCrosschainSwap(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	utils.LoudPrintf("Testing Bitcoin ERC20 crosschain swap...\n")
	// Firstly, deposit 1.15 BTC into Zeta for liquidity
	//sm.DepositBTC()
	// Secondly, deposit 1000.0 USDT into Zeta for liquidity
	utils.LoudPrintf("Depositing 1000 USDT & 1.15 BTC for liquidity\n")

	txhash := sm.DepositERC20(big.NewInt(1e9), []byte{})
	utils.WaitCctxMinedByInTxHash(txhash.Hex(), sm.CctxClient)

	sm.ZevmAuth.GasLimit = 10000000

	// TODO: move into setup and skip it if already initialized
	// https://github.com/zeta-chain/node-private/issues/88
	// it is kept as is for now to be consistent with the old implementation
	// if the tx fails due to already initialized, it will be ignored
	tx, err := sm.UniswapV2Factory.CreatePair(sm.ZevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	//fmt.Printf("USDT-BTC pair receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	usdtBtcPair, err := sm.UniswapV2Factory.GetPair(&bind.CallOpts{}, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT-BTC pair addr %s\n", usdtBtcPair.Hex())

	tx, err = sm.USDTZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("USDT ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	tx, err = sm.BTCZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("BTC ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	// Add 100 USDT liq and 0.001 BTC
	bal, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	bal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	tx, err = sm.UniswapV2Router.AddLiquidity(
		sm.ZevmAuth,
		sm.USDTZRC20Addr,
		sm.BTCZRC20Addr,
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e5),
		sm.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		fmt.Printf("Error liq %s", err.Error())
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("Add liquidity receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	fmt.Printf("Funding contracts ZEVMSwapApp with gas ZRC20s; 1e7 ETH, 1e6 BTC\n")
	// Fund ZEVMSwapApp with gas ZRC20s
	tx, err = sm.ETHZRC20.Transfer(sm.ZevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e7))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("  USDT ZRC20 transfer receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	bal1, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.ZEVMSwapAppAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  ZEVMSwapApp ETHZRC20 balance %d", bal1)
	tx, err = sm.BTCZRC20.Transfer(sm.ZevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e6))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("  BTC ZRC20 transfer receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	bal2, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.ZEVMSwapAppAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  ZEVMSwapApp BTCZRC20 balance %d", bal2)

	// msg would be [ZEVMSwapAppAddr, memobytes]
	// memobytes is dApp specific; see the contracts/ZEVMSwapApp.sol for details
	msg := []byte{}
	msg = append(msg, sm.ZEVMSwapAppAddr.Bytes()...)
	memobytes, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.BTCZRC20Addr, []byte(sm.BTCDeployerAddress.EncodeAddress()))

	if err != nil {
		panic(err)
	}
	fmt.Printf("memobytes(%d) %x\n", len(memobytes), memobytes)
	msg = append(msg, memobytes...)

	fmt.Printf("***** First test: USDT -> BTC\n")
	// Should deposit USDT for swap, swap for BTC and withdraw BTC
	txhash = sm.DepositERC20(big.NewInt(8e7), msg)
	cctx1 := utils.WaitCctxMinedByInTxHash(txhash.Hex(), sm.CctxClient)

	// check the cctx status
	if cctx1.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected outbound mined status; got %s, message: %s", cctx1.CctxStatus.Status.String(), cctx1.CctxStatus.StatusMessage))
	}

	_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	// cctx1 index acts like the inTxHash for the second cctx (the one that withdraws BTC)
	cctx2 := utils.WaitCctxMinedByInTxHash(cctx1.Index, sm.CctxClient)

	// check the cctx status
	if cctx2.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected outbound mined status; got %s, message: %s",
			cctx2.CctxStatus.Status.String(),
			cctx2.CctxStatus.StatusMessage),
		)
	}

	fmt.Printf("cctx2 outbound tx hash %s\n", cctx2.GetCurrentOutTxParam().OutboundTxHash)

	fmt.Printf("******* Second test: BTC -> USDT\n")
	utxos, err := sm.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	fmt.Printf("#utxos %d\n", len(utxos))
	//fmt.Printf("Unimplemented!\n")
	fmt.Printf("memo address %s\n", sm.USDTZRC20Addr)
	memo, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.USDTZRC20Addr, sm.DeployerAddress.Bytes())
	if err != nil {
		panic(err)
	}
	memo = append(sm.ZEVMSwapAppAddr.Bytes(), memo...)
	fmt.Printf("memo length %d\n", len(memo))

	txid, err := sm.SendToTSSFromDeployerWithMemo(
		sm.BTCTSSAddress,
		0.01,
		utxos[0:2],
		sm.BtcRPCClient,
		memo,
		sm.BTCDeployerAddress,
	)
	fmt.Printf("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation\n", txid)
	_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}

	cctx3 := utils.WaitCctxMinedByInTxHash(txid.String(), sm.CctxClient)
	fmt.Printf("cctx3 index %s\n", cctx3.Index)
	fmt.Printf("  inboudn tx hash %s\n", cctx3.InboundTxParams.InboundTxObservedHash)
	fmt.Printf("  status %s\n", cctx3.CctxStatus.Status.String())
	fmt.Printf("  status msg: %s\n", cctx3.CctxStatus.StatusMessage)

	cctx4 := utils.WaitCctxMinedByInTxHash(cctx3.Index, sm.CctxClient)
	fmt.Printf("cctx4 index %s\n", cctx4.Index)
	fmt.Printf("  outbound tx hash %s\n", cctx4.GetCurrentOutTxParam().OutboundTxHash)
	fmt.Printf("  status %s\n", cctx4.CctxStatus.Status.String())

	{
		fmt.Printf("******* Third test: BTC -> ETH with contract call reverted; should refund BTC\n")
		utxos, err := sm.BtcRPCClient.ListUnspent()
		if err != nil {
			panic(err)
		}
		fmt.Printf("#utxos %d\n", len(utxos))
		// the following memo will result in a revert in the contract call as targetZRC20 is set to DeployerAddress
		// which is apparently not a ZRC20 contract; the UNISWAP call will revert
		memo, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.DeployerAddress, sm.DeployerAddress.Bytes())
		if err != nil {
			panic(err)
		}
		memo = append(sm.ZEVMSwapAppAddr.Bytes(), memo...)
		fmt.Printf("memo length %d\n", len(memo))

		amount := 0.1
		txid, err := sm.SendToTSSFromDeployerWithMemo(
			sm.BTCTSSAddress,
			amount,
			utxos[0:2],
			sm.BtcRPCClient,
			memo,
			sm.BTCDeployerAddress,
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation\n", txid)
		_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}

		cctx := utils.WaitCctxMinedByInTxHash(txid.String(), sm.CctxClient)
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
		txraw, err := sm.BtcRPCClient.GetRawTransactionVerbose(outTxHash)
		if err != nil {
			panic(err)
		}
		fmt.Printf("out txid %s\n", txraw.Txid)
		for _, vout := range txraw.Vout {
			fmt.Printf("  vout %d\n", vout.N)
			fmt.Printf("  value %f\n", vout.Value)
			fmt.Printf("  scriptPubKey %s\n", vout.ScriptPubKey.Hex)
			fmt.Printf("  p2wpkh address: %s\n", utils.ScriptPKToAddress(vout.ScriptPubKey.Hex))
		}
	}
}
