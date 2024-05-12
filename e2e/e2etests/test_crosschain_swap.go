package e2etests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestCrosschainSwap(r *runner.E2ERunner, _ []string) {
	r.ZEVMAuth.GasLimit = 10000000

	// TODO: move into setup and skip it if already initialized
	// https://github.com/zeta-chain/node-private/issues/88
	// it is kept as is for now to be consistent with the old implementation
	// if the tx fails due to already initialized, it will be ignored
	_, err := r.UniswapV2Factory.CreatePair(r.ZEVMAuth, r.ERC20ZRC20Addr, r.BTCZRC20Addr)
	if err != nil {
		r.Logger.Print("ℹ️create pair error")
	}
	txERC20ZRC20Approve, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	txBTCApprove, err := r.BTCZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}

	// Fund ZEVMSwapApp with gas ZRC20s
	txTransferETH, err := r.ETHZRC20.Transfer(r.ZEVMAuth, r.ZEVMSwapAppAddr, big.NewInt(1e7))
	if err != nil {
		panic(err)
	}
	txTransferBTC, err := r.BTCZRC20.Transfer(r.ZEVMAuth, r.ZEVMSwapAppAddr, big.NewInt(1e6))
	if err != nil {
		panic(err)
	}

	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txERC20ZRC20Approve, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("ZRC20 ERC20 approve failed")
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txBTCApprove, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("btc approve failed")
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txTransferETH, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("ETH ZRC20 transfer failed")
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txTransferBTC, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("BTC ZRC20 transfer failed")
	}

	// Add 100 erc20 zrc20 liq and 0.001 BTC
	txAddLiquidity, err := r.UniswapV2Router.AddLiquidity(
		r.ZEVMAuth,
		r.ERC20ZRC20Addr,
		r.BTCZRC20Addr,
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e5),
		r.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		panic(fmt.Sprintf("Error liq %s", err.Error()))
	}

	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txAddLiquidity, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("add liq receipt status is not 1")
	}

	// msg would be [ZEVMSwapAppAddr, memobytes]
	// memobytes is dApp specific; see the contracts/ZEVMSwapApp.sol for details
	msg := []byte{}
	msg = append(msg, r.ZEVMSwapAppAddr.Bytes()...)
	memobytes, err := r.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, r.BTCZRC20Addr, []byte(r.BTCDeployerAddress.EncodeAddress()))

	if err != nil {
		panic(err)
	}
	r.Logger.Info("memobytes(%d) %x", len(memobytes), memobytes)
	msg = append(msg, memobytes...)

	r.Logger.Info("***** First test: ERC20 -> BTC")
	// Should deposit ERC20 for swap, swap for BTC and withdraw BTC
	txHash := r.DepositERC20WithAmountAndMessage(r.DeployerAddress, big.NewInt(8e7), msg)
	cctx1 := utils.WaitCctxMinedByInTxHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	// check the cctx status
	if cctx1.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected outbound mined status; got %s, message: %s", cctx1.CctxStatus.Status.String(), cctx1.CctxStatus.StatusMessage))
	}

	// mine 10 blocks to confirm the outbound tx
	_, err = r.BtcRPCClient.GenerateToAddress(10, r.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	stop := r.MineBlocks()

	// cctx1 index acts like the inTxHash for the second cctx (the one that withdraws BTC)
	cctx2 := utils.WaitCctxMinedByInTxHash(r.Ctx, cctx1.Index, r.CctxClient, r.Logger, r.CctxTimeout)

	// check the cctx status
	if cctx2.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected outbound mined status; got %s, message: %s",
			cctx2.CctxStatus.Status.String(),
			cctx2.CctxStatus.StatusMessage),
		)
	}

	r.Logger.Info("cctx2 outbound tx hash %s", cctx2.GetCurrentOutTxParam().OutboundTxHash)

	r.Logger.Info("******* Second test: BTC -> ERC20ZRC20")
	utxos, err := r.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	r.Logger.Info("#utxos %d", len(utxos))
	r.Logger.Info("memo address %s", r.ERC20ZRC20Addr)
	memo, err := r.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, r.ERC20ZRC20Addr, r.DeployerAddress.Bytes())
	if err != nil {
		panic(err)
	}
	memo = append(r.ZEVMSwapAppAddr.Bytes(), memo...)
	r.Logger.Info("memo length %d", len(memo))

	txID, err := r.SendToTSSFromDeployerWithMemo(
		r.BTCTSSAddress,
		0.01,
		utxos[0:2],
		r.BtcRPCClient,
		memo,
		r.BTCDeployerAddress,
	)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation", txID)
	_, err = r.BtcRPCClient.GenerateToAddress(10, r.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}

	cctx3 := utils.WaitCctxMinedByInTxHash(r.Ctx, txID.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx3.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected outbound mined status; got %s, message: %s",
			cctx3.CctxStatus.Status.String(),
			cctx3.CctxStatus.StatusMessage),
		)
	}
	r.Logger.Info("cctx3 index %s", cctx3.Index)
	r.Logger.Info("  inbound tx hash %s", cctx3.InboundTxParams.InboundTxObservedHash)
	r.Logger.Info("  status %s", cctx3.CctxStatus.Status.String())
	r.Logger.Info("  status msg: %s", cctx3.CctxStatus.StatusMessage)

	cctx4 := utils.WaitCctxMinedByInTxHash(r.Ctx, cctx3.Index, r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx4.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected outbound mined status; got %s, message: %s",
			cctx3.CctxStatus.Status.String(),
			cctx3.CctxStatus.StatusMessage),
		)
	}
	r.Logger.Info("cctx4 index %s", cctx4.Index)
	r.Logger.Info("  outbound tx hash %s", cctx4.GetCurrentOutTxParam().OutboundTxHash)
	r.Logger.Info("  status %s", cctx4.CctxStatus.Status.String())

	{
		r.Logger.Info("******* Third test: BTC -> ETH with contract call reverted; should refund BTC")
		utxos, err := r.BtcRPCClient.ListUnspent()
		if err != nil {
			panic(err)
		}
		r.Logger.Info("#utxos %d", len(utxos))
		// the following memo will result in a revert in the contract call as targetZRC20 is set to DeployerAddress
		// which is apparently not a ZRC20 contract; the UNISWAP call will revert
		memo, err := r.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, r.DeployerAddress, r.DeployerAddress.Bytes())
		if err != nil {
			panic(err)
		}
		memo = append(r.ZEVMSwapAppAddr.Bytes(), memo...)
		r.Logger.Info("memo length %d", len(memo))

		amount := 0.1
		txid, err := r.SendToTSSFromDeployerWithMemo(
			r.BTCTSSAddress,
			amount,
			utxos, //[0:2],
			r.BtcRPCClient,
			memo,
			r.BTCDeployerAddress,
		)
		if err != nil {
			panic(err)
		}
		r.Logger.Info("Sent BTC to TSS txid %s; now mining 10 blocks for confirmation", txid)
		_, err = r.BtcRPCClient.GenerateToAddress(10, r.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}

		cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, txid.String(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.Info("cctx3 index %s", cctx.Index)
		r.Logger.Info("  inbound tx hash %s", cctx.InboundTxParams.InboundTxObservedHash)
		r.Logger.Info("  status %s", cctx.CctxStatus.Status.String())
		r.Logger.Info("  status msg: %s", cctx.CctxStatus.StatusMessage)

		if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
			panic(fmt.Sprintf("expected reverted status; got %s", cctx.CctxStatus.Status.String()))
		}
		outTxHash, err := chainhash.NewHashFromStr(cctx.GetCurrentOutTxParam().OutboundTxHash)
		if err != nil {
			panic(err)
		}
		txraw, err := r.BtcRPCClient.GetRawTransactionVerbose(outTxHash)
		if err != nil {
			panic(err)
		}
		r.Logger.Info("out txid %s", txraw.Txid)
		for _, vout := range txraw.Vout {
			r.Logger.Info("  vout %d", vout.N)
			r.Logger.Info("  value %f", vout.Value)
			r.Logger.Info("  scriptPubKey %s", vout.ScriptPubKey.Hex)
			r.Logger.Info("  p2wpkh address: %s", utils.ScriptPKToAddress(vout.ScriptPubKey.Hex, r.BitcoinParams))
		}
	}

	// stop mining
	stop <- struct{}{}
}
