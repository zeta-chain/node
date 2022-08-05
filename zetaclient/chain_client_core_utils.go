package zetaclient

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

func (ob *ChainObserver) PostNonceIfNotRecorded() error {
	zetaClient := ob.zetaClient
	evmClient := ob.EvmClient
	tss := ob.Tss
	chain := ob.chain

	_, err := zetaClient.GetNonceByChain(chain)
	if err != nil { // if Nonce of Chain is not found in ZetaCore; report it
		nonce, err := evmClient.NonceAt(context.TODO(), tss.Address(), nil)
		if err != nil {
			log.Fatal().Err(err).Msg("NonceAt")
			return err
		}
		pendingNonce, err := evmClient.PendingNonceAt(context.TODO(), tss.Address())
		if err != nil {
			log.Fatal().Err(err).Msg("PendingNonceAt")
			return err
		}
		if pendingNonce != nonce {
			log.Fatal().Msgf("fatal: pending nonce %d != nonce %d", pendingNonce, nonce)
			return fmt.Errorf("pending nonce %d != nonce %d", pendingNonce, nonce)
		}
		if err != nil {
			log.Fatal().Err(err).Msg("NonceAt")
			return err
		}
		log.Debug().Msgf("signer %s Posting Nonce of chain %s of nonce %d", zetaClient.GetKeys().signerName, chain, nonce)
		_, err = zetaClient.PostNonce(chain, nonce)
		if err != nil {
			log.Fatal().Err(err).Msg("PostNonce")
			return err
		}
	}
	return nil
}

func (ob *ChainObserver) WatchGasPrice() {
	gasTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-gasTicker.C:
			err := ob.PostGasPrice()
			if err != nil {
				log.Err(err).Msg("PostGasPrice error on " + ob.chain.String())
				continue
			}
		case <-ob.stop:
			log.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *ChainObserver) PostGasPrice() error {
	// GAS PRICE
	gasPrice, err := ob.EvmClient.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}
	blockNum, err := ob.EvmClient.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}

	// SUPPLY
	var supply string // lockedAmount on ETH, totalSupply on other chains
	supply = "100"
	//if chainOb.chain == common.ETHChain {
	//	input, err := chainOb.connectorAbi.Pack("getLockedAmount")
	//	if err != nil {
	//		return fmt.Errorf("fail to getLockedAmount")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.ETH_MPI_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.connectorAbi.Unpack("getLockedAmount", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	lockedAmount := *connectorAbi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("ETH: block %d: lockedAmount %d\n", bn, lockedAmount)
	//	supply = lockedAmount.String()
	//
	//} else if chainOb.chain == common.BSCChain {
	//	input, err := chainOb.connectorAbi.Pack("totalSupply")
	//	if err != nil {
	//		return fmt.Errorf("fail to totalSupply")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.BSC_MPI_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.connectorAbi.Unpack("totalSupply", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	totalSupply := *connectorAbi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("BSC: block %d: totalSupply %d\n", bn, totalSupply)
	//	supply = totalSupply.String()
	//} else if chainOb.chain == common.POLYGONChain {
	//	input, err := chainOb.connectorAbi.Pack("totalSupply")
	//	if err != nil {
	//		return fmt.Errorf("fail to totalSupply")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.POLYGON_MPI_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.connectorAbi.Unpack("totalSupply", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	totalSupply := *connectorAbi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("BSC: block %d: totalSupply %d\n", bn, totalSupply)
	//	supply = totalSupply.String()
	//} else {
	//	log.Error().Msgf("chain not supported %s", chainOb.chain)
	//	return fmt.Errorf("unsupported chain %s", chainOb.chain)
	//}

	_, err = ob.zetaClient.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}

	//bal, err := chainOb.Client.BalanceAt(context.TODO(), chainOb.Tss.Address(), nil)
	//if err != nil {
	//	log.Err(err).Msg("BalanceAt:")
	//	return err
	//}
	//_, err = chainOb.zetaClient.PostGasBalance(chainOb.chain, bal.String(), blockNum)
	//if err != nil {
	//	log.Err(err).Msg("PostGasBalance:")
	//	return err
	//}
	return nil
}

// query ZetaCore about the last block that it has heard from a specific chain.
// return 0 if not existent.
func (ob *ChainObserver) getLastHeight() uint64 {
	lastheight, err := ob.zetaClient.GetLastBlockHeightByChain(ob.chain)
	if err != nil {
		log.Warn().Err(err).Msgf("getLastHeight")
		return 0
	}
	return lastheight.LastSendHeight
}

func (ob *ChainObserver) WatchExchangeRate() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			price, bn, err := ob.ZetaPriceQuerier.GetZetaPrice()
			if err != nil {
				log.Err(err).Msg("GetZetaExchangeRate error on " + ob.chain.String())
				continue
			}
			priceInHex := fmt.Sprintf("0x%x", price)

			_, err = ob.zetaClient.PostZetaConversionRate(ob.chain, priceInHex, bn)
			if err != nil {
				log.Err(err).Msg("PostZetaConversionRate error on " + ob.chain.String())
			}
		case <-ob.stop:
			log.Info().Msg("WatchExchangeRate stopped")
			return
		}
	}
}

// TODO : Call this function from shepard send or in a Separate Goroutine
//func (ob *ChainObserver) HandleReceipts(receipt types.Receipt,nonce int )  {
//	if receipt.Status == 1 {
//		logs := receipt.Logs
//		for _, vLog := range logs {
//			receivedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReceived(*vLog)
//			if err == nil {
//				log.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", inTXHash, ob.chain, vLog.TxHash.Hex())
//				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
//					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
//					sendHash := vLog.Topics[3].Hex()
//					//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
//					mMint := receivedLog.ZetaAmount.String()
//					zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
//						sendHash,
//						vLog.TxHash.Hex(),
//						vLog.BlockNumber,
//						mMint,
//						common.ReceiveStatus_Success,
//						ob.chain.String(),
//						nonce,
//					)
//					if err != nil {
//						log.Error().Err(err).Msg("error posting confirmation to meta core")
//						continue
//					}
//					log.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
//					return true, true, nil
//				} else {
//					log.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.LastBlock), ob.chain, nonce)
//					return true, false, nil
//				}
//			}
//			revertedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReverted(*vLog)
//			if err == nil {
//				log.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", inTXHash, ob.chain, vLog.TxHash.Hex())
//				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
//					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
//					sendhash := vLog.Topics[3].Hex()
//					mMint := revertedLog.ZetaAmount.String()
//					metaHash, err := ob.zetaClient.PostReceiveConfirmation(
//						sendhash,
//						vLog.TxHash.Hex(),
//						vLog.BlockNumber,
//						mMint,
//						common.ReceiveStatus_Success,
//						ob.chain.String(),
//						nonce,
//					)
//					if err != nil {
//						log.Err(err).Msg("error posting confirmation to meta core")
//						continue
//					}
//					log.Info().Msgf("Zeta tx hash: %s", metaHash)
//					return true, true, nil
//				} else {
//					log.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.LastBlock), ob.chain, nonce)
//					return true, false, nil
//				}
//			}
//		}
//	} else if receipt.Status == 0 {
//		log.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", inTXHash, ob.chain, receipt.TxHash.Hex())
//		zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, ob.chain.String(), nonce)
//		if err != nil {
//			log.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
//		}
//		log.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
//		return true, true, nil
//	}
//}
