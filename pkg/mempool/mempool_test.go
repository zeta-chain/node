package mempool_test

import (
	"fmt"
	"math/big"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
)

// testPubKey is a dummy implementation of PubKey used for testing.
type testPubKey struct {
	address sdk.AccAddress
}

func (t testPubKey) Reset() { panic("not implemented") }

func (t testPubKey) String() string { panic("not implemented") }

func (t testPubKey) ProtoMessage() { panic("not implemented") }

func (t testPubKey) Address() cryptotypes.Address { return t.address.Bytes() }

func (t testPubKey) Bytes() []byte { panic("not implemented") }

func (t testPubKey) VerifySignature(msg []byte, sig []byte) bool { panic("not implemented") }

func (t testPubKey) Equals(key cryptotypes.PubKey) bool { panic("not implemented") }

func (t testPubKey) Type() string { panic("not implemented") }

type testTxDetailsGetter interface {
	GetID() int
	GetPriority() int64
	GetNonce() uint64
	GetAddress() sdk.AccAddress
}

// testTx is a dummy implementation of cosmos Tx used for testing.
type testTx struct {
	id       int
	priority int64
	nonce    uint64
	address  sdk.AccAddress
	// useful for debugging
	strAddress string
}

func (tx testTx) GetID() int                 { return tx.id }
func (tx testTx) GetPriority() int64         { return tx.priority }
func (tx testTx) GetNonce() uint64           { return tx.nonce }
func (tx testTx) GetAddress() sdk.AccAddress { return tx.address }

func (tx testTx) GetSigners() []sdk.AccAddress { panic("not implemented") }

func (tx testTx) GetPubKeys() ([]cryptotypes.PubKey, error) { panic("not implemented") }

func (tx testTx) GetSignaturesV2() (res []txsigning.SignatureV2, err error) {
	res = append(res, txsigning.SignatureV2{
		PubKey:   testPubKey{address: tx.address},
		Data:     nil,
		Sequence: tx.nonce,
	})

	return res, nil
}

// testTx is a dummy implementation of unsigned cosmos Tx used for testing.
type testUnsignedTx struct {
	id       int
	priority int64
	nonce    uint64
	address  sdk.AccAddress
}

func (tx testUnsignedTx) GetID() int                 { return tx.id }
func (tx testUnsignedTx) GetPriority() int64         { return tx.priority }
func (tx testUnsignedTx) GetNonce() uint64           { return tx.nonce }
func (tx testUnsignedTx) GetAddress() sdk.AccAddress { return tx.address }

func (tx testUnsignedTx) GetSigners() []sdk.AccAddress { panic("not implemented") }

func (tx testUnsignedTx) GetPubKeys() ([]cryptotypes.PubKey, error) { panic("not implemented") }

func (tx testUnsignedTx) GetSignaturesV2() (res []txsigning.SignatureV2, err error) {
	return res, nil
}

var (
	_ sdk.Tx                  = (*testTx)(nil)
	_ sdk.Tx                  = (*testUnsignedTx)(nil)
	_ signing.SigVerifiableTx = (*testTx)(nil)
	_ signing.SigVerifiableTx = (*testUnsignedTx)(nil)
	_ cryptotypes.PubKey      = (*testPubKey)(nil)
)

func (tx testTx) GetMsgs() []sdk.Msg { return nil }

func (tx testTx) ValidateBasic() error { return nil }

func (tx testTx) String() string {
	return fmt.Sprintf("tx a: %s, p: %d, n: %d", tx.address, tx.priority, tx.nonce)
}

func (tx testUnsignedTx) GetMsgs() []sdk.Msg { return nil }

func (tx testUnsignedTx) ValidateBasic() error { return nil }

func (tx testUnsignedTx) String() string {
	return fmt.Sprintf("tx a: %s, p: %d, n: %d", tx.address, tx.priority, tx.nonce)
}

// testEthTx is a dummy implementation of ethermint Tx used for testing.
type testEthTx struct {
	id              int
	priority        int64
	nonce           uint64
	address         sdk.AccAddress
	extensionOption *codectypes.Any
	msgs            []sdk.Msg
}

func (tx testEthTx) GetID() int                 { return tx.id }
func (tx testEthTx) GetPriority() int64         { return tx.priority }
func (tx testEthTx) GetNonce() uint64           { return tx.nonce }
func (tx testEthTx) GetAddress() sdk.AccAddress { return tx.address }

func (tx testEthTx) GetExtensionOptions() []*codectypes.Any {
	return []*codectypes.Any{tx.extensionOption}
}

func (tx testEthTx) GetNonCriticalExtensionOptions() []*codectypes.Any { panic("not implemented") }

func (tx testEthTx) GetSigners() []sdk.AccAddress { panic("not implemented") }

func (tx testEthTx) GetPubKeys() ([]cryptotypes.PubKey, error) { panic("not implemented") }

// testEthTx is a dummy implementation of unsigned ethermint Tx used for testing.
type testUnsignedEthTx struct {
	id              int
	priority        int64
	nonce           uint64
	address         sdk.AccAddress
	extensionOption *codectypes.Any
	msgs            []sdk.Msg
}

func (tx testUnsignedEthTx) GetID() int                 { return tx.id }
func (tx testUnsignedEthTx) GetPriority() int64         { return tx.priority }
func (tx testUnsignedEthTx) GetNonce() uint64           { return tx.nonce }
func (tx testUnsignedEthTx) GetAddress() sdk.AccAddress { return tx.address }

func (tx testUnsignedEthTx) GetExtensionOptions() []*codectypes.Any {
	return []*codectypes.Any{tx.extensionOption}
}

func (tx testUnsignedEthTx) GetNonCriticalExtensionOptions() []*codectypes.Any {
	panic("not implemented")
}

func (tx testUnsignedEthTx) GetSigners() []sdk.AccAddress { panic("not implemented") }

func (tx testUnsignedEthTx) GetPubKeys() ([]cryptotypes.PubKey, error) { panic("not implemented") }

var (
	_ sdk.Tx                         = (*testEthTx)(nil)
	_ authante.HasExtensionOptionsTx = (*testEthTx)(nil)
	_ sdk.Tx                         = (*testUnsignedEthTx)(nil)
	_ authante.HasExtensionOptionsTx = (*testUnsignedEthTx)(nil)
)

func (tx testEthTx) GetMsgs() []sdk.Msg { return tx.msgs }

func (tx testEthTx) ValidateBasic() error { return nil }

func (tx testEthTx) String() string {
	return fmt.Sprintf("tx a: %s, p: %d, n: %d", tx.address, tx.priority, tx.nonce)
}

func (tx testUnsignedEthTx) GetMsgs() []sdk.Msg { return tx.msgs }

func (tx testUnsignedEthTx) ValidateBasic() error { return nil }

func (tx testUnsignedEthTx) String() string {
	return fmt.Sprintf("tx a: %s, p: %d, n: %d", tx.address, tx.priority, tx.nonce)
}

type txSpec struct {
	i int
	p int
	n int
	a sdk.AccAddress
}

func (tx txSpec) String() string {
	return fmt.Sprintf("[tx i: %d, a: %s, p: %d, n: %d]", tx.i, tx.a, tx.p, tx.n)
}

func (s *MempoolTestSuite) buildMockEthTx(id int, priority int64, from string, nonce uint64) testEthTx {
	msg := evmtypes.NewTx(big.NewInt(1), nonce, nil, big.NewInt(1), 0, nil, nil, nil, nil, nil)
	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
	require.NoError(s.T(), err)
	msg.From = from
	return testEthTx{
		id:              id,
		priority:        priority,
		nonce:           nonce,
		msgs:            []sdk.Msg{msg},
		extensionOption: option,
		address:         common.HexToAddress(from).Bytes(),
	}
}

func (s *MempoolTestSuite) buildInvalidMockEthTx(id int, priority int64, from string, nonce uint64) testUnsignedEthTx {
	msg := evmtypes.NewTx(big.NewInt(1), nonce, nil, big.NewInt(1), 0, nil, nil, nil, nil, nil)
	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
	require.NoError(s.T(), err)
	msg.From = from
	return testUnsignedEthTx{
		id:              id,
		priority:        priority,
		nonce:           nonce,
		msgs:            []sdk.Msg{}, // empty msgs
		extensionOption: option,
		address:         common.HexToAddress(from).Bytes(),
	}
}

func fetchTxs(iterator mempool.Iterator, maxBytes int64) []sdk.Tx {
	const txSize = 1
	var (
		txs      []sdk.Tx
		numBytes int64
	)
	for iterator != nil {
		if numBytes += txSize; numBytes > maxBytes {
			break
		}
		txs = append(txs, iterator.Tx())
		i := iterator.Next()
		iterator = i
	}
	return txs
}

type MempoolTestSuite struct {
	suite.Suite
	numTxs      int
	numAccounts int
	iterations  int
	mempool     mempool.Mempool
}

func (s *MempoolTestSuite) resetMempool() {
	s.iterations = 0
	s.mempool = mempool.NewSenderNonceMempool()
}

func (s *MempoolTestSuite) SetupTest() {
	s.numTxs = 1000
	s.numAccounts = 100
	s.resetMempool()
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}
