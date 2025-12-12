package backend

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	abci "github.com/cometbft/cometbft/abci/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/evm/crypto/ethsecp256k1"
	"github.com/cosmos/evm/indexer"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/node/rpc/backend/mocks"
	rpctypes "github.com/zeta-chain/node/rpc/types"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/crypto"
)

func (s *TestSuite) TestTraceTransaction() {
	msgEthereumTx, _ := s.buildEthereumTx()
	msgEthereumTx2, _ := s.buildEthereumTx()

	txHash := msgEthereumTx.AsTransaction().Hash()
	txHash2 := msgEthereumTx2.AsTransaction().Hash()

	priv, _ := ethsecp256k1.GenerateKey()
	from := common.BytesToAddress(priv.PubKey().Address().Bytes())
	armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
	_ = s.backend.ClientCtx.Keyring.ImportPrivKey("test_key", armor, "")

	ethSigner := ethtypes.LatestSigner(s.backend.ChainConfig())

	txEncoder := s.backend.ClientCtx.TxConfig.TxEncoder()

	msgEthereumTx.From = from.Bytes()
	_ = msgEthereumTx.Sign(ethSigner, s.signer)

	baseDenom := evmtypes.GetEVMCoinDenom()

	tx, _ := msgEthereumTx.BuildTx(s.backend.ClientCtx.TxConfig.NewTxBuilder(), baseDenom)
	txBz, _ := txEncoder(tx)

	msgEthereumTx2.From = from.Bytes()
	_ = msgEthereumTx2.Sign(ethSigner, s.signer)

	tx2, _ := msgEthereumTx2.BuildTx(s.backend.ClientCtx.TxConfig.NewTxBuilder(), baseDenom)
	txBz2, _ := txEncoder(tx2)

	testCases := []struct {
		name          string
		txHash        common.Hash
		registerMock  func()
		block         *types.Block
		responseBlock []*abci.ExecTxResult
		expResult     interface{}
		expPass       bool
	}{
		{
			"fail - tx not found",
			txHash,
			func() {},
			&types.Block{Header: types.Header{Height: 1}, Data: types.Data{Txs: []types.Tx{}}},
			[]*abci.ExecTxResult{
				{
					Code: 0,
					Events: []abci.Event{
						{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
							{Key: "ethereumTxHash", Value: txHash.Hex()},
							{Key: "txIndex", Value: "0"},
							{Key: "amount", Value: "1000"},
							{Key: "txGasUsed", Value: "21000"},
							{Key: "txHash", Value: ""},
							{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
						}},
					},
				},
			},
			nil,
			false,
		},
		{
			"fail - block not found",
			txHash,
			func() {
				// var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, 1)
			},
			&types.Block{Header: types.Header{Height: 1}, Data: types.Data{Txs: []types.Tx{txBz}}},
			[]*abci.ExecTxResult{
				{
					Code: 0,
					Events: []abci.Event{
						{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
							{Key: "ethereumTxHash", Value: txHash.Hex()},
							{Key: "txIndex", Value: "0"},
							{Key: "amount", Value: "1000"},
							{Key: "txGasUsed", Value: "21000"},
							{Key: "txHash", Value: ""},
							{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
						}},
					},
				},
			},
			map[string]interface{}{"test": "hello"},
			false,
		},
		{
			"pass - transaction found in a block with multiple transactions",
			txHash2, // tx1 is predecessor of tx2
			func() {
				var (
					QueryClient       = s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
					client            = s.backend.ClientCtx.Client.(*mocks.Client)
					height      int64 = 1
				)
				_, err := RegisterBlockMultipleTxs(client, height, []types.Tx{txBz, txBz2})
				s.Require().NoError(err)
				RegisterTraceTransactionWithPredecessors(
					QueryClient,
					msgEthereumTx2,
					[]*evmtypes.MsgEthereumTx{msgEthereumTx},
				)
				RegisterConsensusParams(client, height)
				txResults := []*abci.ExecTxResult{
					{
						Code: 0,
						Events: []abci.Event{
							{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
								{Key: "ethereumTxHash", Value: txHash.Hex()},
								{Key: "txIndex", Value: "0"},
								{Key: "amount", Value: "1000"},
								{Key: "txGasUsed", Value: "21000"},
								{Key: "txHash", Value: ""},
								{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
							}},
						},
					},
					{
						Code: 0,
						Events: []abci.Event{
							{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
								{Key: "ethereumTxHash", Value: txHash2.Hex()},
								{Key: "txIndex", Value: "1"},
								{Key: "amount", Value: "1000"},
								{Key: "txGasUsed", Value: "21000"},
								{Key: "txHash", Value: ""},
								{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
							}},
						},
					},
				}

				RegisterBlockResultsWithTxResults(client, 1, txResults)
			},
			&types.Block{
				Header: types.Header{Height: 1, ChainID: ChainID.ChainID},
				Data:   types.Data{Txs: []types.Tx{txBz, txBz2}},
			},
			[]*abci.ExecTxResult{
				{
					Code: 0,
					Events: []abci.Event{
						{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
							{Key: "ethereumTxHash", Value: txHash.Hex()},
							{Key: "txIndex", Value: "0"},
							{Key: "amount", Value: "1000"},
							{Key: "txGasUsed", Value: "21000"},
							{Key: "txHash", Value: ""},
							{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
						}},
					},
				},
				{
					Code: 0,
					Events: []abci.Event{
						{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
							{Key: "ethereumTxHash", Value: txHash2.Hex()},
							{Key: "txIndex", Value: "1"},
							{Key: "amount", Value: "1000"},
							{Key: "txGasUsed", Value: "21000"},
							{Key: "txHash", Value: ""},
							{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
						}},
					},
				},
			},
			map[string]interface{}{"test": "hello"},
			true,
		},
		{
			"pass - transaction found",
			txHash,
			func() {
				var (
					QueryClient       = s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
					client            = s.backend.ClientCtx.Client.(*mocks.Client)
					height      int64 = 1
				)
				_, err := RegisterBlock(client, height, txBz)
				s.Require().NoError(err)
				RegisterTraceTransaction(QueryClient, msgEthereumTx)
				RegisterConsensusParams(client, height)
				txResults := []*abci.ExecTxResult{
					{
						Code: 0,
						Events: []abci.Event{
							{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
								{Key: "ethereumTxHash", Value: txHash.Hex()},
								{Key: "txIndex", Value: "0"},
								{Key: "amount", Value: "1000"},
								{Key: "txGasUsed", Value: "21000"},
								{Key: "txHash", Value: ""},
								{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
							}},
						},
					},
				}
				RegisterBlockResultsWithTxResults(client, 1, txResults)
			},
			&types.Block{Header: types.Header{Height: 1}, Data: types.Data{Txs: []types.Tx{txBz}}},
			[]*abci.ExecTxResult{
				{
					Code: 0,
					Events: []abci.Event{
						{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
							{Key: "ethereumTxHash", Value: txHash.Hex()},
							{Key: "txIndex", Value: "0"},
							{Key: "amount", Value: "1000"},
							{Key: "txGasUsed", Value: "21000"},
							{Key: "txHash", Value: ""},
							{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
						}},
					},
				},
			},
			map[string]interface{}{"test": "hello"},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock()

			db := dbm.NewMemDB()
			s.backend.Indexer = indexer.NewKVIndexer(db, log.NewNopLogger(), s.backend.ClientCtx)

			err := s.backend.Indexer.IndexBlock(tc.block, tc.responseBlock)
			s.Require().NoError(err)
			txResult, err := s.backend.TraceTransaction(tc.txHash, nil)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().Equal(tc.expResult, txResult)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestTraceBlock() {
	msgEthTx, bz := s.buildEthereumTx()
	emptyBlock := types.MakeBlock(1, []types.Tx{}, nil, nil)
	emptyBlock.ChainID = ChainID.ChainID
	filledBlock := types.MakeBlock(1, []types.Tx{bz}, nil, nil)
	filledBlock.ChainID = ChainID.ChainID
	resBlockEmpty := tmrpctypes.ResultBlock{Block: emptyBlock, BlockID: emptyBlock.LastBlockID}
	resBlockFilled := tmrpctypes.ResultBlock{Block: filledBlock, BlockID: filledBlock.LastBlockID}

	testCases := []struct {
		name            string
		registerMock    func()
		expTraceResults []*evmtypes.TxTraceResult
		resBlock        *tmrpctypes.ResultBlock
		config          *rpctypes.TraceConfig
		expPass         bool
	}{
		{
			"pass - no transaction returning empty array",
			func() {},
			[]*evmtypes.TxTraceResult{},
			&resBlockEmpty,
			&rpctypes.TraceConfig{},
			true,
		},
		{
			"fail - cannot unmarshal data",
			func() {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterTraceBlock(QueryClient, []*evmtypes.MsgEthereumTx{msgEthTx})
				RegisterConsensusParams(client, 1)
				RegisterBlockResults(client, 1)
			},
			[]*evmtypes.TxTraceResult{},
			&resBlockFilled,
			&rpctypes.TraceConfig{},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock()

			traceResults, err := s.backend.TraceBlock(1, tc.config, tc.resBlock)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().Equal(tc.expTraceResults, traceResults)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
