package metaclient

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// Based on old router query code, currently WIP
func query_router_deposit(httpserver *TssHttpServer, chain string, bridge *metaclient.MetachainBridge) {
	var NODE string
	var ROUTER string
	switch chain {
	case "Polygon":
		NODE = POLYGON_NODE
		ROUTER = POLYGON_ROUTER
	case "Goerli":
		NODE = GOERLI_INFURA
		ROUTER = GOERLI_ROUTER
	case "BSCTestnet":
		NODE = BSC_NODE
		ROUTER = BSC_ROUTER
	}

	log.Info().Msgf("Starting monitoring deposit on %s", chain)
	cl, err := ethclient.Dial(NODE)
	if err != nil {
		log.Fatal().Err(err).Msgf("dial %s error", chain)
	}
	ctx := context.Background()

	// ticker to run every blocktime
	blockTime := metaclient.ChainNameToBlocktime[chain]
	ticker := time.NewTicker(time.Duration(blockTime) * time.Second)

	chainid, err := cl.ChainID(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("chainid error")
	}
	log.Info().Msgf("chainid: %s", chainid.String())

	router_address := ethcommon.HexToAddress(ROUTER)

	// 5533138
	var mostRecentBlock *big.Int

	// chain, _ := metacommon.NewChain("ETH")
	// // 	lastObserved, err := b.GetLastBlockObserved(chain)

	chainObj, _ := metacommon.NewChain(strings.ToUpper(chain))

	mostRecentBlockQuery, err := bridge.GetLastBlockObserved(chainObj)
	if err != nil {
		log.Warn().Err(err).Msg("last block observed error")
	}

	if mostRecentBlockQuery == 0 {
		// get most recent block number
		header, err := cl.HeaderByNumber(context.Background(), nil)
		bigIntTen := big.NewInt(int64(10))

		mostRecentBlock = big.NewInt(0).Sub(header.Number, bigIntTen) // hard coded as current block minus 10 if not found -- this can be tuned
		if err != nil {
			log.Fatal()
		}
	} else {
		mostRecentBlock = big.NewInt(int64(mostRecentBlockQuery))
	}

	fmt.Printf("Starting block height query: %d\n", mostRecentBlockQuery) // 5671744
	// big.NewInt(int64(5533267))

	fmt.Println("mostRecentBlock is ", mostRecentBlock)
	// go func() {
	for {
		select {
		case t := <-ticker.C:
			fmt.Println("Ticker at ", t)

			// get most recent block number
			header, err := cl.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Warn().Err(err).Msg("block header err")
				continue // if error, then continue loop
			}

			fmt.Printf("Most recent block: %s\n", header.Number.String()) // 5671744

			// to block = now, from block = most recently queried block
			query := ethereum.FilterQuery{
				Addresses: []ethcommon.Address{router_address},
				FromBlock: mostRecentBlock,
				ToBlock:   header.Number,
			}

			logs, err := cl.FilterLogs(context.Background(), query)
			if err != nil {
				log.Fatal()
				continue // keep going to next iteration
			}

			contractABI, err := abi.JSON(strings.NewReader(string(contracts.MainABI)))
			if err != nil {
				log.Fatal().Err(err).Msg("ABI contract err")
				continue // keep going to next iteration
			}

			logDepositSig := []byte("DepositMade(address,uint256,address,uint256,address,address,uint256)")
			logDepositSigHash := crypto.Keccak256Hash(logDepositSig)

			logWithdrawSig := []byte("WithdrawalMade(address,uint256,address,uint256")
			logWithdrawSigHash := crypto.Keccak256Hash(logWithdrawSig)

			// update most recent block
			mostRecentBlock = header.Number

			for _, vLog := range logs {
				// tx has already been processed, so we should continue
				if _, ok := bridge.ProcessedTransactions[vLog.TxHash.Hex()]; ok {
					continue
				}
				switch vLog.Topics[0].Hex() {
				case logDepositSigHash.Hex():
					returnVal, err := contractABI.Unpack("DepositMade", vLog.Data)

					if err != nil {
						log.Warn().Err(err).Msg("unpack err")
						continue // keep going to next iteration
					}
					fmt.Printf("Source Address: %s\n", ethcommon.HexToAddress(vLog.Topics[1].Hex()))
					fmt.Printf("Destination Address: %s\n", ethcommon.HexToAddress(vLog.Topics[2].Hex()))
					fmt.Printf("Destination Unit: %d\n", returnVal[0])
					fmt.Printf("MTP Burned: %d\n", returnVal[1])
					fmt.Printf("Source Token: %s\n", returnVal[2])
					fmt.Printf("Destination Token: %s\n", returnVal[3])
					fmt.Printf("Amount Deposited: %d\n", returnVal[4])

					bigIntAmountDeposited := returnVal[4].(*big.Int)
					if bigIntAmountDeposited == nil || bigIntAmountDeposited.Int64() <= 0 {
						log.Warn().Msg("amount deposited is invalid")
						break // shouldn't keep executing if invalid deposit
					}
					uintAmountDeposited := bigIntAmountDeposited.Uint64()

					bigIntMTPBurned := returnVal[1].(*big.Int)
					uintMTPBurned := bigIntMTPBurned.Uint64()

					fmt.Println("Uint Deposited: ", uintAmountDeposited)
					fmt.Println("Uint burned: ", uintMTPBurned)

					stringSourceToken := fmt.Sprintf("%v", returnVal[2])
					// sourceTicker := metaclient.AddressToTicker[chain][stringSourceToken] // use token address now instead of ticker so we don't need to store lookup
					sourceAsset := fmt.Sprintf("%s.%s", stringSourceToken, chain)
					fmt.Println("Source Asset: ", sourceAsset)

					destChainId := fmt.Sprintf("%v", returnVal[0])
					destChainName := metaclient.ChainIdToName[destChainId]

					fmt.Println("destination chain name: ", destChainName)

					stringDestToken := fmt.Sprintf("%v", returnVal[3])
					// destTicker := metaclient.AddressToTicker[destChainName][stringDestToken]
					destAsset := fmt.Sprintf("%s.%s", stringDestToken, destChainName)

					fmt.Println("Dest Asset: ", destAsset)

					fmt.Println("Attempting to post transaction")
					// need to figure out mapping for chainIDs and assets
					err = bridge.PostTxIn(sourceAsset, uintAmountDeposited, uintMTPBurned, destAsset, vLog.Topics[2].Hex(), vLog.TxHash.Hex(), vLog.BlockNumber)

					if err != nil {
						log.Warn().Err(err).Msg("tx post in err")
					}

					bridge.ProcessedTransactions[vLog.TxHash.Hex()] = 0 // Pending --> 0 is enum value defined in metachain.go

					// c.Assert(err, IsNil)

					fmt.Println("Successfully posted transaction")

				case logWithdrawSigHash.Hex():
					returnVal, err := contractABI.Unpack("WithdrawalMade", vLog.Data)
					if err != nil {
						log.Warn().Err(err).Msg("unpack withdrawal err")
					}
					fmt.Println("Withdrawal made: ", returnVal)
				}
			}
		}
	}
	// }()
}
