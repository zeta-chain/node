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
	sm.ZevmAuth.GasLimit = 10000000

	// TODO: move into setup and skip it if already initialized
	// https://github.com/zeta-chain/node-private/issues/88
	// it is kept as is for now to be consistent with the old implementation
	// if the tx fails due to already initialized, it will be ignored
	txCreatePair, err := sm.UniswapV2Factory.CreatePair(sm.ZevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	if err != nil {
		panic(err)
	}
	txUSDTApprove, err := sm.USDTZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	txBTCApprove, err := sm.BTCZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}

	// Fund ZEVMSwapApp with gas ZRC20s
	txTransferETH, err := sm.ETHZRC20.Transfer(sm.ZevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e7))
	if err != nil {
		panic(err)
	}
	txTransferBTC, err := sm.BTCZRC20.Transfer(sm.ZevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e6))
	if err != nil {
		panic(err)
	}

	if receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, txCreatePair, sm.Logger); receipt.Status != 1 {
		panic("create pair failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, txUSDTApprove, sm.Logger); receipt.Status != 1 {
		panic("usdt approve failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, txBTCApprove, sm.Logger); receipt.Status != 1 {
		panic("btc approve failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, txTransferETH, sm.Logger); receipt.Status != 1 {
		panic("ETH ZRC20 transfer failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, txTransferBTC, sm.Logger); receipt.Status != 1 {
		panic("BTC ZRC20 transfer failed")
	}

	// Add 100 USDT liq and 0.001 BTC
	txAddLiquidity, err := sm.UniswapV2Router.AddLiquidity(
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
		panic(fmt.Sprintf("Error liq %s", err.Error()))
	}

	if receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, txAddLiquidity, sm.Logger); receipt.Status != 1 {
		panic("add liq receipt status is not 1")
	}

	// msg would be [ZEVMSwapAppAddr, memobytes]
	// memobytes is dApp specific; see the contracts/ZEVMSwapApp.sol for details
	msg := []byte{}
	msg = append(msg, sm.ZEVMSwapAppAddr.Bytes()...)
	memobytes, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.BTCZRC20Addr, []byte(sm.BTCDeployerAddress.EncodeAddress()))

	if err != nil {
		panic(err)
	}
	sm.Logger.Info("memobytes(%d) %x", len(memobytes), memobytes)
	msg = append(msg, memobytes...)

	sm.Logger.Info("***** First test: USDT -> BTC")
	// Should deposit USDT for swap, swap for BTC and withdraw BTC
	txHash := sm.DepositERC20WithAmountAndMessage(big.NewInt(8e7), msg)
	cctx1 := utils.WaitCctxMinedByInTxHash(txHash.Hex(), sm.CctxClient, sm.Logger)

	// check the cctx status
	if cctx1.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected outbound mined status; got %s, message: %s", cctx1.CctxStatus.Status.String(), cctx1.CctxStatus.StatusMessage))
	}

	// mine 10 blocks to confirm the outbound tx
	_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	stop := sm.MineBlocks()

	// cctx1 index acts like the inTxHash for the second cctx (the one that withdraws BTC)
	cctx2 := utils.WaitCctxMinedByInTxHash(cctx1.Index, sm.CctxClient, sm.Logger)

	// check the cctx status
	if cctx2.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected outbound mined status; got %s, message: %s",
			cctx2.CctxStatus.Status.String(),
			cctx2.CctxStatus.StatusMessage),
		)
	}

	sm.Logger.Info("cctx2 outbound tx hash %s", cctx2.GetCurrentOutTxParam().OutboundTxHash)

	sm.Logger.Info("******* Second test: BTC -> USDT")
	utxos, err := sm.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("#utxos %d", len(utxos))
	sm.Logger.Info("memo address %s", sm.USDTZRC20Addr)
	memo, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.USDTZRC20Addr, sm.DeployerAddress.Bytes())
	if err != nil {
		panic(err)
	}
	memo = append(sm.ZEVMSwapAppAddr.Bytes(), memo...)
	sm.Logger.Info("memo length %d", len(memo))

	txId, err := sm.SendToTSSFromDeployerWithMemo(
		sm.BTCTSSAddress,
		0.01,
		utxos[0:2],
		sm.BtcRPCClient,
		memo,
		sm.BTCDeployerAddress,
	)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation", txId)
	_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}

	cctx3 := utils.WaitCctxMinedByInTxHash(txId.String(), sm.CctxClient, sm.Logger)
	sm.Logger.Info("cctx3 index %s", cctx3.Index)
	sm.Logger.Info("  inbound tx hash %s", cctx3.InboundTxParams.InboundTxObservedHash)
	sm.Logger.Info("  status %s", cctx3.CctxStatus.Status.String())
	sm.Logger.Info("  status msg: %s", cctx3.CctxStatus.StatusMessage)

	cctx4 := utils.WaitCctxMinedByInTxHash(cctx3.Index, sm.CctxClient, sm.Logger)
	sm.Logger.Info("cctx4 index %s", cctx4.Index)
	sm.Logger.Info("  outbound tx hash %s", cctx4.GetCurrentOutTxParam().OutboundTxHash)
	sm.Logger.Info("  status %s", cctx4.CctxStatus.Status.String())

	{
		sm.Logger.Info("******* Third test: BTC -> ETH with contract call reverted; should refund BTC")
		utxos, err := sm.BtcRPCClient.ListUnspent()
		if err != nil {
			panic(err)
		}
		sm.Logger.Info("#utxos %d", len(utxos))
		// the following memo will result in a revert in the contract call as targetZRC20 is set to DeployerAddress
		// which is apparently not a ZRC20 contract; the UNISWAP call will revert
		memo, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.DeployerAddress, sm.DeployerAddress.Bytes())
		if err != nil {
			panic(err)
		}
		memo = append(sm.ZEVMSwapAppAddr.Bytes(), memo...)
		sm.Logger.Info("memo length %d", len(memo))

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
		sm.Logger.Info("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation", txid)
		_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}

		cctx := utils.WaitCctxMinedByInTxHash(txid.String(), sm.CctxClient, sm.Logger)
		sm.Logger.Info("cctx3 index http://localhost:1317/zeta-chain/crosschain/cctx/%s", cctx.Index)
		sm.Logger.Info("  inboudn tx hash %s", cctx.InboundTxParams.InboundTxObservedHash)
		sm.Logger.Info("  status %s", cctx.CctxStatus.Status.String())
		sm.Logger.Info("  status msg: %s", cctx.CctxStatus.StatusMessage)

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
		sm.Logger.Info("out txid %s", txraw.Txid)
		for _, vout := range txraw.Vout {
			sm.Logger.Info("  vout %d", vout.N)
			sm.Logger.Info("  value %f", vout.Value)
			sm.Logger.Info("  scriptPubKey %s", vout.ScriptPubKey.Hex)
			sm.Logger.Info("  p2wpkh address: %s", utils.ScriptPKToAddress(vout.ScriptPubKey.Hex))
		}
	}

	// stop mining
	stop <- struct{}{}
}
