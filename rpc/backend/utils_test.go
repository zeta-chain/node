package backend

import (
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func mookProofs(num int, withData bool) *crypto.ProofOps {
	var proofOps *crypto.ProofOps
	if num > 0 {
		proofOps = new(crypto.ProofOps)
		for i := 0; i < num; i++ {
			proof := crypto.ProofOp{}
			if withData {
				proof.Data = []byte("\n\031\n\003KEY\022\005VALUE\032\013\010\001\030\001 \001*\003\000\002\002")
			}
			proofOps.Ops = append(proofOps.Ops, proof)
		}
	}
	return proofOps
}

func (s *TestSuite) TestGetHexProofs() {
	defaultRes := []string{""}
	testCases := []struct {
		name  string
		proof *crypto.ProofOps
		exp   []string
	}{
		{
			"no proof provided",
			mookProofs(0, false),
			defaultRes,
		},
		{
			"no proof data provided",
			mookProofs(1, false),
			defaultRes,
		},
		{
			"valid proof provided",
			mookProofs(1, true),
			[]string{"0x0a190a034b4559120556414c55451a0b0801180120012a03000202"},
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.Require().Equal(tc.exp, GetHexProofs(tc.proof))
		})
	}
}

// createMockLog creates a mock EVM log with the given parameters
func createMockLog(address common.Address, txIndex, logIndex uint, topic common.Hash) *evmtypes.Log {
	return &evmtypes.Log{
		Address:     address.Hex(),
		Topics:      []string{topic.Hex()},
		Data:        []byte{},
		BlockNumber: 100,
		TxHash:      common.BytesToHash([]byte("txhash")).Hex(),
		TxIndex:     uint64(txIndex),
		BlockHash:   common.BytesToHash([]byte("blockhash")).Hex(),
		Index:       uint64(logIndex),
		Removed:     false,
	}
}

// createTxLogEvent creates a mock tx log event with the given logs
func createTxLogEvent(logs []*evmtypes.Log) abci.Event {
	attrs := make([]abci.EventAttribute, 0, len(logs))
	for _, log := range logs {
		logJSON, _ := json.Marshal(log)
		attrs = append(attrs, abci.EventAttribute{
			Key:   evmtypes.AttributeKeyTxLog,
			Value: string(logJSON),
		})
	}
	return abci.Event{
		Type:       evmtypes.EventTypeTxLog,
		Attributes: attrs,
	}
}

// createTxResult creates a mock transaction result with the given events
func createTxResult(events []abci.Event) *abci.ExecTxResult {
	return &abci.ExecTxResult{
		Code:   0,
		Events: events,
	}
}

func (s *TestSuite) TestGetLogsFromBlockResults() {
	addr1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	addr2 := common.HexToAddress("0x2222222222222222222222222222222222222222")
	topic1 := common.HexToHash("0x49f492222906ac486c3c1401fa545626df1f0c0e5a77a05597ea2ed66af9850d")
	topic2 := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	testCases := []struct {
		name               string
		blockRes           *cmtrpctypes.ResultBlockResults
		expectedTxCount    int
		expectedLogCount   int
		expectedTxIndices  []uint // expected TxIndex for each log group
		expectedLogIndices []uint // expected Index (log index) for each log
		expectError        bool
	}{
		{
			name: "empty block results",
			blockRes: &cmtrpctypes.ResultBlockResults{
				TxsResults: []*abci.ExecTxResult{},
			},
			expectedTxCount:    0,
			expectedLogCount:   0,
			expectedTxIndices:  []uint{},
			expectedLogIndices: []uint{},
			expectError:        false,
		},
		{
			name: "single transaction with single log",
			blockRes: &cmtrpctypes.ResultBlockResults{
				TxsResults: []*abci.ExecTxResult{
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr1, 0, 0, topic1), // original indices are 0, 0
						}),
					}),
				},
			},
			expectedTxCount:    1,
			expectedLogCount:   1,
			expectedTxIndices:  []uint{0},
			expectedLogIndices: []uint{0},
			expectError:        false,
		},
		{
			name: "single transaction with multiple logs",
			blockRes: &cmtrpctypes.ResultBlockResults{
				TxsResults: []*abci.ExecTxResult{
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr1, 0, 0, topic1), // original: txIndex=0, logIndex=0
							createMockLog(addr1, 0, 1, topic2), // original: txIndex=0, logIndex=1
						}),
					}),
				},
			},
			expectedTxCount:    1,
			expectedLogCount:   2,
			expectedTxIndices:  []uint{0, 0}, // both logs belong to tx 0
			expectedLogIndices: []uint{0, 1}, // sequential log indices
			expectError:        false,
		},
		{
			name: "multiple transactions with duplicate indices - reported bug",
			blockRes: &cmtrpctypes.ResultBlockResults{
				TxsResults: []*abci.ExecTxResult{
					// First Cosmos transaction with EVM log event
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr1, 0, 0, topic1), // BSC setGasPrice - original: txIndex=0, logIndex=0
						}),
					}),
					// Second Cosmos transaction with EVM log event
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr2, 0, 0, topic1), // Base setGasPrice - original: txIndex=0, logIndex=0
						}),
					}),
				},
			},
			expectedTxCount:    2,
			expectedLogCount:   2,
			expectedTxIndices:  []uint{0, 1}, // should be re-indexed to 0, 1
			expectedLogIndices: []uint{0, 1}, // should be re-indexed to 0, 1
			expectError:        false,
		},
		{
			name: "multiple transactions with multiple logs each",
			blockRes: &cmtrpctypes.ResultBlockResults{
				TxsResults: []*abci.ExecTxResult{
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr1, 0, 0, topic1),
							createMockLog(addr1, 0, 0, topic2), // duplicate logIndex in original
						}),
					}),
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr2, 0, 0, topic1), // duplicate txIndex and logIndex
							createMockLog(addr2, 0, 0, topic2), // duplicate txIndex and logIndex
						}),
					}),
				},
			},
			expectedTxCount:    2,
			expectedLogCount:   4,
			expectedTxIndices:  []uint{0, 0, 1, 1}, // logs 0,1 belong to tx 0; logs 2,3 belong to tx 1
			expectedLogIndices: []uint{0, 1, 2, 3}, // globally unique log indices
			expectError:        false,
		},
		{
			name: "transaction without log events",
			blockRes: &cmtrpctypes.ResultBlockResults{
				TxsResults: []*abci.ExecTxResult{
					createTxResult([]abci.Event{
						{Type: "other_event", Attributes: []abci.EventAttribute{}},
					}),
				},
			},
			expectedTxCount:    0,
			expectedLogCount:   0,
			expectedTxIndices:  []uint{},
			expectedLogIndices: []uint{},
			expectError:        false,
		},
		{
			name: "logs already have unique indices - no reindexing needed",
			blockRes: &cmtrpctypes.ResultBlockResults{
				TxsResults: []*abci.ExecTxResult{
					// First transaction with properly indexed logs
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr1, 0, 0, topic1), // Already correct: txIndex=0, logIndex=0
							createMockLog(addr1, 0, 1, topic2), // Already correct: txIndex=0, logIndex=1
						}),
					}),
					// Second transaction with properly indexed logs
					createTxResult([]abci.Event{
						createTxLogEvent([]*evmtypes.Log{
							createMockLog(addr2, 1, 2, topic1), // Already correct: txIndex=1, logIndex=2
							createMockLog(addr2, 1, 3, topic2), // Already correct: txIndex=1, logIndex=3
						}),
					}),
				},
			},
			expectedTxCount:    2,
			expectedLogCount:   4,
			expectedTxIndices:  []uint{0, 0, 1, 1}, // Original indices preserved (after future zetacore fix)
			expectedLogIndices: []uint{0, 1, 2, 3}, // Original indices preserved
			expectError:        false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			logs, err := GetLogsFromBlockResults(tc.blockRes)

			if tc.expectError {
				s.Require().Error(err)
				return
			}

			s.Require().NoError(err)
			s.Require().Equal(tc.expectedTxCount, len(logs), "unexpected number of transaction log groups")

			// Flatten logs to check indices
			var allLogs []*evmtypes.Log
			for _, txLogs := range logs {
				for _, log := range txLogs {
					allLogs = append(allLogs, &evmtypes.Log{
						TxIndex: uint64(log.TxIndex),
						Index:   uint64(log.Index),
					})
				}
			}

			s.Require().Equal(tc.expectedLogCount, len(allLogs), "unexpected total number of logs")

			for i, log := range allLogs {
				if i < len(tc.expectedTxIndices) {
					s.Require().Equal(tc.expectedTxIndices[i], uint(log.TxIndex),
						"log %d has wrong TxIndex", i)
				}
				if i < len(tc.expectedLogIndices) {
					s.Require().Equal(tc.expectedLogIndices[i], uint(log.Index),
						"log %d has wrong Index (logIndex)", i)
				}
			}

			seenLogIndices := make(map[uint]bool)
			for _, log := range allLogs {
				s.Require().False(seenLogIndices[uint(log.Index)],
					"duplicate logIndex found: %d", log.Index)
				seenLogIndices[uint(log.Index)] = true
			}
		})
	}
}

func (s *TestSuite) TestNeedsReindexing() {
	addr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	topic := common.HexToHash("0x49f492222906ac486c3c1401fa545626df1f0c0e5a77a05597ea2ed66af9850d")

	testCases := []struct {
		name     string
		logs     [][]*ethtypes.Log
		expected bool
	}{
		{
			name:     "empty logs",
			logs:     [][]*ethtypes.Log{},
			expected: false,
		},
		{
			name: "single log - no duplicates",
			logs: [][]*ethtypes.Log{
				{createEthLog(addr, 0, 0, topic)},
			},
			expected: false,
		},
		{
			name: "multiple logs with unique indices",
			logs: [][]*ethtypes.Log{
				{createEthLog(addr, 0, 0, topic), createEthLog(addr, 0, 1, topic)},
				{createEthLog(addr, 1, 2, topic), createEthLog(addr, 1, 3, topic)},
			},
			expected: false,
		},
		{
			name: "duplicate log indices - needs reindexing",
			logs: [][]*ethtypes.Log{
				{createEthLog(addr, 0, 0, topic)},
				{createEthLog(addr, 0, 0, topic)}, // Duplicate logIndex=0
			},
			expected: true,
		},
		{
			name: "duplicate within same tx group",
			logs: [][]*ethtypes.Log{
				{createEthLog(addr, 0, 0, topic), createEthLog(addr, 0, 0, topic)}, // Duplicate
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := needsReindexing(tc.logs)
			s.Require().Equal(tc.expected, result)
		})
	}
}

// createEthLog creates an ethtypes.Log for testing
func createEthLog(address common.Address, txIndex, logIndex uint, topic common.Hash) *ethtypes.Log {
	return &ethtypes.Log{
		Address:     address,
		Topics:      []common.Hash{topic},
		Data:        []byte{},
		BlockNumber: 100,
		TxHash:      common.BytesToHash([]byte("txhash")),
		TxIndex:     txIndex,
		BlockHash:   common.BytesToHash([]byte("blockhash")),
		Index:       logIndex,
		Removed:     false,
	}
}

// TestGetLogsFromBlockResults_IndexUniqueness specifically tests the original bug scenario
func (s *TestSuite) TestGetLogsFromBlockResults_IndexUniqueness() {
	addr := common.HexToAddress("0x91d18e54daf4f677cb28167158d6dd21f6ab3921")
	topic := common.HexToHash("0x49f492222906ac486c3c1401fa545626df1f0c0e5a77a05597ea2ed66af9850d")

	// Create two logs with identical indices reported issue
	log1 := createMockLog(addr, 0, 0, topic) // BSC gas price update
	log2 := createMockLog(addr, 0, 0, topic) // Base gas price update

	blockRes := &cmtrpctypes.ResultBlockResults{
		TxsResults: []*abci.ExecTxResult{
			createTxResult([]abci.Event{createTxLogEvent([]*evmtypes.Log{log1})}),
			createTxResult([]abci.Event{createTxLogEvent([]*evmtypes.Log{log2})}),
		},
	}

	logs, err := GetLogsFromBlockResults(blockRes)
	s.Require().NoError(err)
	s.Require().Equal(2, len(logs), "expected 2 transaction log groups")

	s.Require().Equal(uint(0), logs[0][0].TxIndex, "first log should have TxIndex 0")
	s.Require().Equal(uint(0), logs[0][0].Index, "first log should have logIndex 0")

	s.Require().Equal(uint(1), logs[1][0].TxIndex, "second log should have TxIndex 1")
	s.Require().Equal(uint(1), logs[1][0].Index, "second log should have logIndex 1")

	type logKey struct {
		txIndex  uint
		logIndex uint
	}
	seen := make(map[logKey]bool)
	for _, txLogs := range logs {
		for _, log := range txLogs {
			key := logKey{txIndex: log.TxIndex, logIndex: log.Index}
			s.Require().False(seen[key], "duplicate (txIndex=%d, logIndex=%d) found", log.TxIndex, log.Index)
			seen[key] = true
		}
	}
}
