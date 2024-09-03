// Forked from here and extended to test ethermint txs https://github.com/cosmos/cosmos-sdk/blob/v0.47.10/baseapp/abci_utils_test.go
// TODO: remove this once cosmos is upgraded: https://github.com/zeta-chain/node/issues/2156
package mempool_test

import (
	"bytes"
	"math/big"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtsecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	baseapptestutil "github.com/cosmos/cosmos-sdk/baseapp/testutil"
	"github.com/cosmos/cosmos-sdk/baseapp/testutil/mock"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"

	zetamempool "github.com/zeta-chain/node/pkg/mempool"
)

type testValidator struct {
	consAddr sdk.ConsAddress
	tmPk     cmtprotocrypto.PublicKey
	privKey  cmtsecp256k1.PrivKey
}

func newTestValidator() testValidator {
	privkey := cmtsecp256k1.GenPrivKey()
	pubkey := privkey.PubKey()
	tmPk := cmtprotocrypto.PublicKey{
		Sum: &cmtprotocrypto.PublicKey_Secp256K1{
			Secp256K1: pubkey.Bytes(),
		},
	}

	return testValidator{
		consAddr: sdk.ConsAddress(pubkey.Address()),
		tmPk:     tmPk,
		privKey:  privkey,
	}
}

type ABCIUtilsTestSuite struct {
	suite.Suite

	vals [3]testValidator
	ctx  sdk.Context
}

func NewABCIUtilsTestSuite(t *testing.T) *ABCIUtilsTestSuite {
	t.Helper()
	// create 3 validators
	s := &ABCIUtilsTestSuite{
		vals: [3]testValidator{
			newTestValidator(),
			newTestValidator(),
			newTestValidator(),
		},
	}

	// create context
	s.ctx = sdk.Context{}.WithConsensusParams(&cmtproto.ConsensusParams{})
	return s
}

func TestABCIUtilsTestSuite(t *testing.T) {
	suite.Run(t, NewABCIUtilsTestSuite(t))
}

func (s *ABCIUtilsTestSuite) TestCustomProposalHandler_NoopMempool() {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	baseapptestutil.RegisterInterfaces(cdc.InterfaceRegistry())
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	var (
		secret1 = []byte("secret1")
		secret2 = []byte("secret2")
	)
	type testTx struct {
		tx       sdk.Tx
		priority int64
		bz       []byte
	}
	testTxs := []testTx{
		{tx: buildMsg(s.T(), txConfig, []byte(`0`), [][]byte{secret1}, []uint64{1}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`12345678910`), [][]byte{secret1}, []uint64{2}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`22`), [][]byte{secret1}, []uint64{3}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`32`), [][]byte{secret2}, []uint64{1}), priority: 8},
	}
	ctrl := gomock.NewController(s.T())
	app := mock.NewMockProposalTxVerifier(ctrl)
	mp := mempool.NoOpMempool{}
	customProposalHandler := zetamempool.NewCustomProposalHandler(mp, app)
	prepareProposalHandler := customProposalHandler.PrepareProposalHandler()
	processProposalHandler := customProposalHandler.ProcessProposalHandler()

	req := abci.RequestPrepareProposal{
		MaxTxBytes: 111 + 112,
	}
	for _, v := range testTxs {
		app.EXPECT().PrepareProposalVerifyTx(v.tx).Return(v.bz, nil).AnyTimes()
		s.NoError(mp.Insert(s.ctx.WithPriority(v.priority), v.tx))
		req.Txs = append(req.Txs, v.bz)
	}

	// since it is noop mempool, all txs returned in fifo order
	resp := prepareProposalHandler(s.ctx, req)
	for i, tx := range resp.Txs {
		require.Equal(s.T(), tx, testTxs[i].bz)
	}

	// verify that proposal will be processed
	for i, v := range req.Txs {
		app.EXPECT().ProcessProposalVerifyTx(v).Return(testTxs[i].tx, nil).AnyTimes()
	}
	resp2 := processProposalHandler(s.ctx, abci.RequestProcessProposal{
		Txs: resp.Txs,
	})
	s.Require().Equal(resp2.Status, abci.ResponseProcessProposal_ACCEPT)
}

func (s *ABCIUtilsTestSuite) TestCustomProposalHandler_PriorityNonceMempoolTxSelection() {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	baseapptestutil.RegisterInterfaces(cdc.InterfaceRegistry())
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	var (
		secret1 = []byte("secret1")
		secret2 = []byte("secret2")
		secret3 = []byte("secret3")
		secret4 = []byte("secret4")
		secret5 = []byte("secret5")
		secret6 = []byte("secret6")
	)

	type testTx struct {
		tx       sdk.Tx
		priority int64
		bz       []byte
		size     int
	}

	testTxs := []testTx{
		// skip same-sender non-sequential sequence and then add others txs
		{tx: buildMsg(s.T(), txConfig, []byte(`0`), [][]byte{secret1}, []uint64{1}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`12345678910`), [][]byte{secret1}, []uint64{2}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`22`), [][]byte{secret1}, []uint64{3}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`32`), [][]byte{secret2}, []uint64{1}), priority: 8},
		// skip multi-signers msg non-sequential sequence
		{tx: buildMsg(s.T(), txConfig, []byte(`4`), [][]byte{secret1, secret2}, []uint64{3, 3}), priority: 10},
		{
			tx:       buildMsg(s.T(), txConfig, []byte(`52345678910`), [][]byte{secret1, secret3}, []uint64{4, 3}),
			priority: 10,
		},
		{tx: buildMsg(s.T(), txConfig, []byte(`62`), [][]byte{secret1, secret4}, []uint64{5, 3}), priority: 8},
		{tx: buildMsg(s.T(), txConfig, []byte(`72`), [][]byte{secret3, secret5}, []uint64{4, 3}), priority: 8},
		{tx: buildMsg(s.T(), txConfig, []byte(`82`), [][]byte{secret2, secret6}, []uint64{4, 3}), priority: 8},
		// only the first tx is added
		{tx: buildMsg(s.T(), txConfig, []byte(`9`), [][]byte{secret3, secret4}, []uint64{3, 3}), priority: 10},
		{
			tx:       buildMsg(s.T(), txConfig, []byte(`1052345678910`), [][]byte{secret1, secret2}, []uint64{4, 4}),
			priority: 8,
		},
		{tx: buildMsg(s.T(), txConfig, []byte(`11`), [][]byte{secret1, secret2}, []uint64{5, 5}), priority: 8},
		// no txs added
		{tx: buildMsg(s.T(), txConfig, []byte(`1252345678910`), [][]byte{secret1}, []uint64{3}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`13`), [][]byte{secret1}, []uint64{5}), priority: 10},
		{tx: buildMsg(s.T(), txConfig, []byte(`14`), [][]byte{secret1}, []uint64{6}), priority: 8},
		// (eth tx) skip same-sender non-sequential sequence and then add others txs
		{tx: buildEthMsg(s.T(), txConfig, []byte(`0`), [][]byte{secret1}, []uint64{1}), priority: 10},
		{tx: buildEthMsg(s.T(), txConfig, []byte(`12345678910`), [][]byte{secret1}, []uint64{2}), priority: 10},
		{tx: buildEthMsg(s.T(), txConfig, []byte(`22`), [][]byte{secret1}, []uint64{3}), priority: 10},
		{tx: buildEthMsg(s.T(), txConfig, []byte(`32`), [][]byte{secret2}, []uint64{1}), priority: 8},
		// (eth tx) no txs added
		{tx: buildEthMsg(s.T(), txConfig, []byte(`1252345678910`), [][]byte{secret1}, []uint64{3}), priority: 10},
		{tx: buildEthMsg(s.T(), txConfig, []byte(`13`), [][]byte{secret1}, []uint64{5}), priority: 10},
		{tx: buildEthMsg(s.T(), txConfig, []byte(`14`), [][]byte{secret1}, []uint64{6}), priority: 8},
		// large cosmos tx to test with eth txs
		{
			tx: buildMsg(
				s.T(),
				txConfig,
				[]byte(
					`1252345678910125234567891012523456789101252345678910125234567891012523456789101252345678910125234567891012523456789101252345678912343`,
				),
				[][]byte{secret2},
				[]uint64{1},
			),
			priority: 8,
		},
		// same as tx2, but with same nonce as tx1, used in skip same-sender and nonce and add same-sender and nonce with less size
		{tx: buildMsg(s.T(), txConfig, []byte(`22`), [][]byte{secret1}, []uint64{2}), priority: 10},
	}

	for i := range testTxs {
		bz, err := txConfig.TxEncoder()(testTxs[i].tx)
		s.Require().NoError(err)
		testTxs[i].bz = bz
		testTxs[i].size = int(cmttypes.ComputeProtoSizeForTxs([]cmttypes.Tx{bz}))
	}

	s.Require().Equal(testTxs[0].size, 111)
	s.Require().Equal(testTxs[1].size, 121)
	s.Require().Equal(testTxs[2].size, 112)
	s.Require().Equal(testTxs[3].size, 112)
	s.Require().Equal(testTxs[4].size, 195)
	s.Require().Equal(testTxs[5].size, 205)
	s.Require().Equal(testTxs[6].size, 196)
	s.Require().Equal(testTxs[7].size, 196)
	s.Require().Equal(testTxs[8].size, 196)

	// ethermint txs
	s.Require().Equal(testTxs[15].size, 247)
	s.Require().Equal(testTxs[16].size, 257)
	s.Require().Equal(testTxs[17].size, 248)
	s.Require().Equal(testTxs[18].size, 248)

	s.Require().Equal(testTxs[19].size, 259)
	s.Require().Equal(testTxs[20].size, 248)
	s.Require().Equal(testTxs[21].size, 248)

	s.Require().Equal(testTxs[22].size, 248)

	testCases := map[string]struct {
		ctx         sdk.Context
		txInputs    []testTx
		req         abci.RequestPrepareProposal
		handler     sdk.PrepareProposalHandler
		expectedTxs []int
	}{
		"skip same-sender non-sequential sequence and then add others txs": {
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[0], testTxs[1], testTxs[2], testTxs[3]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 111 + 112,
			},
			expectedTxs: []int{0, 3},
		},
		"skip same-sender and nonce and add same-sender and nonce with less size": {
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[0], testTxs[1], testTxs[23]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 111 + 112,
			},
			expectedTxs: []int{0, 23},
		},
		"(eth tx) skip same-sender non-sequential sequence and then add others txs": {
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[15], testTxs[16], testTxs[17], testTxs[18]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 247 + 248,
			},
			expectedTxs: []int{15, 18},
		},
		"(eth tx + cosmos tx) skip same-sender non-sequential sequence and then add others txs": {
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[15], testTxs[16], testTxs[17], testTxs[22]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 247 + 248,
			},
			expectedTxs: []int{15, 22},
		},
		"skip multi-signers msg non-sequential sequence": {
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[4], testTxs[5], testTxs[6], testTxs[7], testTxs[8]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 195 + 196,
			},
			expectedTxs: []int{4, 8},
		},
		"only the first tx is added": {
			// Because tx 10 is invalid, tx 11 can't be valid as they have higher sequence numbers.
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[9], testTxs[10], testTxs[11]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 195 + 196,
			},
			expectedTxs: []int{9},
		},
		"no txs added": {
			// Becasuse the first tx was deemed valid but too big, the next expected valid sequence is tx[0].seq (3), so
			// the rest of the txs fail because they have a seq of 4.
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[12], testTxs[13], testTxs[14]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 112,
			},
			expectedTxs: []int{},
		},
		"(eth tx) no txs added": {
			// Becasuse the first tx was deemed valid but too big, the next expected valid sequence is tx[0].seq (3), so
			// the rest of the txs fail because they have a seq of 4.
			ctx:      s.ctx,
			txInputs: []testTx{testTxs[19], testTxs[20], testTxs[21]},
			req: abci.RequestPrepareProposal{
				MaxTxBytes: 248,
			},
			expectedTxs: []int{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctrl := gomock.NewController(s.T())
			app := mock.NewMockProposalTxVerifier(ctrl)
			mp := zetamempool.NewPriorityMempool()
			customProposalHandler := zetamempool.NewCustomProposalHandler(mp, app)
			prepareProposalHandler := customProposalHandler.PrepareProposalHandler()
			processProposalHandler := customProposalHandler.ProcessProposalHandler()

			for _, v := range tc.txInputs {
				app.EXPECT().PrepareProposalVerifyTx(v.tx).Return(v.bz, nil).AnyTimes()
				s.NoError(mp.Insert(s.ctx.WithPriority(v.priority), v.tx))
				tc.req.Txs = append(tc.req.Txs, v.bz)
			}

			resp := prepareProposalHandler(tc.ctx, tc.req)
			respTxIndexes := []int{}
			for _, tx := range resp.Txs {
				for i, v := range testTxs {
					if bytes.Equal(tx, v.bz) {
						respTxIndexes = append(respTxIndexes, i)
					}
				}
			}

			s.Require().EqualValues(tc.expectedTxs, respTxIndexes)

			// verify that proposal will be processed
			req := abci.RequestProcessProposal{
				Txs: resp.Txs,
			}
			for i, v := range req.Txs {
				app.EXPECT().ProcessProposalVerifyTx(v).Return(tc.txInputs[i].tx, nil).AnyTimes()
			}
			resp2 := processProposalHandler(tc.ctx, req)
			s.Require().Equal(resp2.Status, abci.ResponseProcessProposal_ACCEPT)
		})
	}
}

func buildMsg(t *testing.T, txConfig client.TxConfig, value []byte, secrets [][]byte, nonces []uint64) sdk.Tx {
	t.Helper()
	builder := txConfig.NewTxBuilder()
	_ = builder.SetMsgs(
		&baseapptestutil.MsgKeyValue{Value: value},
	)
	require.Equal(t, len(secrets), len(nonces))
	signatures := make([]signingtypes.SignatureV2, 0)
	for index, secret := range secrets {
		nonce := nonces[index]
		privKey := secp256k1.GenPrivKeyFromSecret(secret)
		pubKey := privKey.PubKey()
		signatures = append(signatures, signingtypes.SignatureV2{
			PubKey:   pubKey,
			Sequence: nonce,
			Data:     &signingtypes.SingleSignatureData{},
		})
	}
	setTxSignatureWithSecret(t, builder, signatures...)
	return builder.GetTx()
}

func buildEthMsg(t *testing.T, txConfig client.TxConfig, value []byte, secrets [][]byte, nonces []uint64) sdk.Tx {
	t.Helper()
	msgEthereumTx := evmtypes.NewTx(
		big.NewInt(1),
		nonces[0],
		nil,
		nil,
		0,
		nil,
		nil,
		nil,
		value,
		nil,
	)

	privKey := secp256k1.GenPrivKeyFromSecret(secrets[0])
	addr := sdk.AccAddress(privKey.PubKey().Address())
	msgEthereumTx.From = common.BytesToAddress(addr.Bytes()).Hex()
	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
	require.NoError(t, err)
	builder, ok := txConfig.NewTxBuilder().(authtx.ExtensionOptionsTxBuilder)
	require.True(t, ok)

	builder.SetExtensionOptions(option)
	builder.SetMsgs(msgEthereumTx)

	return builder.GetTx()
}

func setTxSignatureWithSecret(t *testing.T, builder client.TxBuilder, signatures ...signingtypes.SignatureV2) {
	t.Helper()
	err := builder.SetSignatures(
		signatures...,
	)
	require.NoError(t, err)
}
