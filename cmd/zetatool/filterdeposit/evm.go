package filterdeposit

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/cmd/zetatool/config"
)

const (
	TopicsDeposited = 2
	TopicsZetaSent  = 3
	DonationMessage = "I am rich!"
)

var evmCmd = &cobra.Command{
	Use:   "eth",
	Short: "Filter inbound eth deposits",
	Run:   FilterEVMTransactions,
}

func init() {
	Cmd.AddCommand(evmCmd)
}

// FilterEVMTransactions is a command that queries an EVM explorer and Contracts for inbound transactions that qualify
// for cross chain transactions.
func FilterEVMTransactions(cmd *cobra.Command, _ []string) {
	configFile, err := cmd.Flags().GetString(config.Flag)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.GetConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	list := GetEthHashList(cfg)
	CheckForCCTX(list, cfg)
}

// GetEthHashList is a helper function querying total inbound txns by segments of blocks in ranges defined by the config
func GetEthHashList(cfg *config.Config) []Deposit {
	startBlock := cfg.EvmStartBlock
	client, err := ethclient.Dial(cfg.EthRPC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection successful")

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	latestBlock := header.Number.Uint64()
	fmt.Println("latest Block: ", latestBlock)

	endBlock := startBlock + cfg.EvmMaxRange
	deposits := make([]Deposit, 0)
	segment := 0
	for startBlock < latestBlock {
		fmt.Printf("adding segment: %d, startblock: %d\n", segment, startBlock)
		deposits = append(deposits, GetHashListSegment(client, startBlock, endBlock, cfg)...)
		startBlock = endBlock
		endBlock = endBlock + cfg.EvmMaxRange
		if endBlock > latestBlock {
			endBlock = latestBlock
		}
		segment++
	}
	return deposits
}

// GetHashListSegment queries and filters deposits for a given range
func GetHashListSegment(client *ethclient.Client, startBlock uint64, endBlock uint64, cfg *config.Config) []Deposit {
	deposits := make([]Deposit, 0)

	connectorAddress := common.HexToAddress(cfg.ConnectorAddress)
	connectorContract, err := zetaconnector.NewZetaConnectorNonEth(connectorAddress, client)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
	erc20CustodyAddress := common.HexToAddress(cfg.CustodyAddress)
	erc20CustodyContract, err := erc20custody.NewERC20Custody(erc20CustodyAddress, client)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}

	custodyIter, err := erc20CustodyContract.FilterDeposited(&bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: context.TODO(),
	}, []common.Address{})
	if err != nil {
		fmt.Println("error loading filter: ", err.Error())
		return deposits
	}

	connectorIter, err := connectorContract.FilterZetaSent(&bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: context.TODO(),
	}, []common.Address{}, []*big.Int{})
	if err != nil {
		fmt.Println("error loading filter: ", err.Error())
		return deposits
	}

	// ********************** Get ERC20 Custody Deposit events
	for custodyIter.Next() {
		// sanity check tx event
		err := CheckEvmTxLog(&custodyIter.Event.Raw, erc20CustodyAddress, "", TopicsDeposited)
		if err == nil {
			//fmt.Println("adding deposits")
			deposits = append(deposits, Deposit{
				TxID:   custodyIter.Event.Raw.TxHash.Hex(),
				Amount: custodyIter.Event.Amount.Uint64(),
			})
		}
	}

	// ********************** Get Connector ZetaSent events
	for connectorIter.Next() {
		// sanity check tx event
		err := CheckEvmTxLog(&connectorIter.Event.Raw, connectorAddress, "", TopicsZetaSent)
		if err == nil {
			//fmt.Println("adding deposits")
			deposits = append(deposits, Deposit{
				TxID:   connectorIter.Event.Raw.TxHash.Hex(),
				Amount: connectorIter.Event.ZetaValueAndGas.Uint64(),
			})
		}
	}

	//********************** Get Transactions sent directly to TSS address
	tssDeposits, err := getTSSDeposits(cfg.TssAddressEVM, startBlock, endBlock)
	if err != nil {
		fmt.Printf("getTSSDeposits returned err: %s", err.Error())
	}
	deposits = append(deposits, tssDeposits...)

	return deposits
}

// getTSSDeposits more specifically queries and filters deposits based on direct transfers the TSS address.
func getTSSDeposits(tssAddress string, startBlock uint64, endBlock uint64) ([]Deposit, error) {
	client := etherscan.New(etherscan.Mainnet, "S3AVTNXDJQZQQUVXJM4XVIPBRYECGK88VX")
	deposits := make([]Deposit, 0)

	// #nosec G701 these block numbers need to be *int for this particular client package
	startInt := int(startBlock)
	// #nosec G701
	endInt := int(endBlock)
	txns, err := client.NormalTxByAddress(tssAddress, &startInt, &endInt, 0, 0, true)
	if err != nil {
		return deposits, err
	}

	fmt.Println("getTSSDeposits - Num of transactions: ", len(txns))

	for _, tx := range txns {
		if tx.To == tssAddress {
			if strings.Compare(tx.Input, DonationMessage) == 0 {
				continue // skip donation tx
			}
			if tx.TxReceiptStatus != "1" {
				continue
			}
			//fmt.Println("getTSSDeposits - adding Deposit")
			deposits = append(deposits, Deposit{
				TxID:   tx.Hash,
				Amount: tx.Value.Int().Uint64(),
			})
		}
	}

	return deposits, nil
}

// CheckEvmTxLog is a helper function used to validate receipts, logic is taken from zetaclient.
func CheckEvmTxLog(vLog *ethtypes.Log, wantAddress common.Address, wantHash string, wantTopics int) error {
	if vLog.Removed {
		return fmt.Errorf("log is removed, chain reorg?")
	}
	if vLog.Address != wantAddress {
		return fmt.Errorf("log emitter address mismatch: want %s got %s", wantAddress.Hex(), vLog.Address.Hex())
	}
	if vLog.TxHash.Hex() == "" {
		return fmt.Errorf("log tx hash is empty: %d %s", vLog.BlockNumber, vLog.TxHash.Hex())
	}
	if wantHash != "" && vLog.TxHash.Hex() != wantHash {
		return fmt.Errorf("log tx hash mismatch: want %s got %s", wantHash, vLog.TxHash.Hex())
	}
	if len(vLog.Topics) != wantTopics {
		return fmt.Errorf("number of topics mismatch: want %d got %d", wantTopics, len(vLog.Topics))
	}
	return nil
}
