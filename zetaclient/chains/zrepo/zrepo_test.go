package zrepo

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	crosschain "github.com/zeta-chain/node/x/crosschain/types"
	observer "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func TestNew(t *testing.T) {
	t.Run("invalid client", func(t *testing.T) {
		repo := New(nil, chains.Ethereum, mode.StandardMode)
		require.Nil(t, repo)
	})

	t.Run("dry-mode", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.DryMode)
		require.NotNil(t, repo)

		_, ok := repo.client.(*dryZetacoreClient)
		require.True(t, ok)
	})
}

func TestDryMode(t *testing.T) {
	t.Run("PostOutboundTracker", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.DryMode)

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)
		const nonce = 1234
		const txHash = "some hash"

		zhash, err := repo.PostOutboundTracker(context.Background(), logger, nonce, txHash)
		require.NoError(t, err)
		require.Empty(t, zhash)
		require.Contains(t, buffer.String(), "skipping outbound tracker")
	})

	t.Run("VoteGasPrice", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.DryMode)

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)
		const gasPrice = 100000
		const priorityFee = 200
		const block = 12345

		zhash, err := repo.VoteGasPrice(context.Background(), logger, gasPrice, priorityFee, block)
		require.NoError(t, err)
		require.Empty(t, zhash)
		require.Contains(t, buffer.String(), "skipping gas price vote")
	})

	t.Run("VoteInbound", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.DryMode)

		cctxNotFoundErr := grpcstatus.Error(grpccodes.InvalidArgument, "anything")
		sourceChainID := chains.Ethereum.ChainId
		targetChainID := chains.ZetaChainMainnet.ChainId
		voteInboundMsg := sample.InboundVote(coin.CoinType_Gas, sourceChainID, targetChainID)
		const retryGasLimit = 100000

		client.MockGetCctxByHash("", cctxNotFoundErr)

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)

		ballot, err := repo.VoteInbound(context.Background(), logger, &voteInboundMsg, retryGasLimit, nil)
		require.NoError(t, err)
		require.Empty(t, ballot)
		require.Contains(t, buffer.String(), "skipping inbound vote")
	})

	t.Run("VoteOutbound", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.DryMode)

		ctx := context.Background()
		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)
		const gasLimit = 10000
		const retryGasLimit = 100000
		msg := sample.OutboundVote(t)

		zhash, ballot, err := repo.VoteOutbound(ctx, logger, gasLimit, retryGasLimit, &msg)
		require.NoError(t, err)
		require.Empty(t, zhash)
		require.Empty(t, ballot)
		require.Contains(t, buffer.String(), "skipping outbound vote")
	})
}

func TestVoteInbound(t *testing.T) {
	cctxNotFoundErr := grpcstatus.Error(grpccodes.InvalidArgument, "anything")
	ballotNotFoundErr := grpcstatus.Error(grpccodes.NotFound, "anything")
	sourceChainID := chains.Ethereum.ChainId
	targetChainID := chains.ZetaChainMainnet.ChainId
	msg := sample.InboundVote(coin.CoinType_Gas, sourceChainID, targetChainID)
	const retryGasLimit = 100000

	t.Run("ok", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.StandardMode)

		mockBallot := "some ballot"

		client.MockGetCctxByHash("", cctxNotFoundErr)
		client.WithPostVoteInbound(sample.ZetaIndex(t), mockBallot)

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)

		ballot, err := repo.VoteInbound(context.Background(), logger, &msg, retryGasLimit, nil)
		require.NoError(t, err)
		require.Equal(t, mockBallot, ballot)
		require.Contains(t, buffer.String(), "posted inbound vote")
		client.AssertNumberOfCalls(t, "PostVoteInbound", 1)
	})

	t.Run("already voted", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.StandardMode)

		mockBallot := "some ballot"

		client.MockGetCctxByHash("", cctxNotFoundErr)
		client.WithPostVoteInbound("", mockBallot)

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)

		ballot, err := repo.VoteInbound(context.Background(), logger, &msg, retryGasLimit, nil)
		require.NoError(t, err)
		require.Equal(t, mockBallot, ballot)
		require.Contains(t, buffer.String(), "already voted on the inbound")
		client.AssertNumberOfCalls(t, "PostVoteInbound", 1)
	})

	t.Run("invalid input message", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.StandardMode)

		client.MockGetCctxByHash("", cctxNotFoundErr)

		msg := msg                                                       // copying
		msg.Message = strings.Repeat("1", crosschain.MaxMessageLength+1) // long mock message

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)

		ballot, err := repo.VoteInbound(context.Background(), logger, &msg, retryGasLimit, nil)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
		require.Contains(t, buffer.String(), "invalid vote-inbound message")
	})

	t.Run("CCTX already exists but the ballot does not", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.StandardMode)

		client.MockGetCctxByHash("anything", nil)
		client.MockGetBallotByID(msg.Digest(), ballotNotFoundErr)

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)

		ballot, err := repo.VoteInbound(context.Background(), logger, &msg, retryGasLimit, nil)
		require.NoError(t, err)
		require.Equal(t, ballot, msg.Digest())
		require.Contains(t, buffer.String(), "not voting on inbound; CCTX exists but the ballot does not")
	})

	t.Run("vote on finalized ballot", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)
		repo := New(client, chains.Ethereum, mode.StandardMode)

		client.MockGetCctxByHash("anything", nil)
		client.MockGetBallotByID(msg.Digest(), nil)
		client.WithPostVoteInbound(sample.ZetaIndex(t), msg.Digest())

		var buffer bytes.Buffer
		logger := zerolog.New(&buffer)

		ballot, err := repo.VoteInbound(context.Background(), logger, &msg, retryGasLimit, nil)
		require.NoError(t, err)
		require.Equal(t, ballot, msg.Digest())
		require.Contains(t, buffer.String(), "posted inbound vote")
		client.AssertNumberOfCalls(t, "PostVoteInbound", 1)
	})
}

// ------------------------------------------------------------------------------------------------
// Tests of auxiliary functions
// ------------------------------------------------------------------------------------------------

func TestCheckCode(t *testing.T) {
	notFoundErr := grpcstatus.Error(grpccodes.NotFound, "anything")
	invalidErr := errors.New("not a GRPC error")

	testCases := []struct {
		name         string
		err          error
		code         grpccodes.Code
		expectedErrs []error
	}{{
		name:         "same code",
		err:          notFoundErr,
		code:         grpccodes.NotFound,
		expectedErrs: nil,
	}, {
		name:         "different code",
		err:          notFoundErr,
		code:         grpccodes.InvalidArgument,
		expectedErrs: []error{notFoundErr},
	}, {
		name:         "invalid error",
		err:          invalidErr,
		code:         grpccodes.InvalidArgument,
		expectedErrs: []error{ErrNotRPCError, invalidErr},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := checkCode(testCase.err, testCase.code)
			for _, expectedErr := range testCase.expectedErrs {
				require.ErrorIs(t, err, expectedErr)
			}
		})
	}
}

// TestExists tests the cctxExists and ballotExists functions (and the exists function implicitly).
func TestExists(t *testing.T) {
	client := mocks.NewZetacoreClient(t)
	repo := New(client, chains.Ethereum, mode.StandardMode)

	invalidErr := errors.New("not a valid GRPC error")

	const validHash1 = "this hash belongs to a CCTX"
	const validHash2 = "this hash belongs to a CCTX, but the client does not return the CCTX"
	const invalidHash1 = "this hash does not belong to a CCTX, and it triggers a valid error"
	const invalidHash2 = "this hash does not belong to a CCTX, and it triggers an invalid error"
	validCCTXErr := grpcstatus.Error(grpccodes.InvalidArgument, "a valid GRPC error (CCTX)")
	cctx := &crosschain.CrossChainTx{}
	client.
		On("GetCctxByHash", mock.Anything, validHash1).Return(cctx, nil).
		On("GetCctxByHash", mock.Anything, validHash2).Return(nil, validCCTXErr).
		On("GetCctxByHash", mock.Anything, invalidHash1).Return(nil, validCCTXErr).
		On("GetCctxByHash", mock.Anything, invalidHash2).Return(nil, invalidErr)

	const validID1 = "this ID belongs to a ballot"
	const validID2 = "this ID belongs to a ballot, but the client does not return the ballot"
	const invalidID1 = "this ID does not belong to a ballot, and it triggers a valid error"
	const invalidID2 = "this ID does not belong to a ballot, and it triggers an invalid error"
	validIDErr := grpcstatus.Error(grpccodes.NotFound, "a valid GRPC error (ID)")
	ballot := &observer.QueryBallotByIdentifierResponse{}
	client.
		On("GetBallotByID", mock.Anything, validID1).Return(ballot, nil).
		On("GetBallotByID", mock.Anything, validID2).Return(nil, validIDErr).
		On("GetBallotByID", mock.Anything, invalidID1).Return(nil, validIDErr).
		On("GetBallotByID", mock.Anything, invalidID2).Return(nil, invalidErr)

	testCases := []struct {
		name           string
		hash           string
		id             string
		expectedExists bool
		expectedErrs   []error
	}{{
		name:           "CCTX exists",
		hash:           validHash1,
		expectedExists: true,
		expectedErrs:   nil,
	}, {
		name:           "CCTX does not exist",
		hash:           invalidHash1,
		expectedExists: false,
		expectedErrs:   nil,
	}, {
		name:           "client does not return the CCTX",
		hash:           validHash2,
		expectedExists: false,
		expectedErrs:   nil,
	}, {
		name:           "client error (GetCCTXByHash)",
		hash:           invalidHash2,
		expectedExists: false,
		expectedErrs:   []error{ErrClient, ErrClientGetCCTXByHash, invalidErr},
	}, {
		name:           "ballot exists",
		id:             validID1,
		expectedExists: true,
		expectedErrs:   nil,
	}, {
		name:           "ballot does not exist",
		id:             invalidID1,
		expectedExists: false,
		expectedErrs:   nil,
	}, {
		name:           "client does not return the ballot",
		id:             validID2,
		expectedExists: false,
		expectedErrs:   nil,
	}, {
		name:           "client error (GetBallotByID)",
		id:             invalidID2,
		expectedExists: false,
		expectedErrs:   []error{ErrClient, ErrClientGetBallotByID, invalidErr},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			require := require.New(t)
			if testCase.hash != "" && testCase.id != "" {
				panic("invalid test case")
			}

			var exists bool
			var err error
			if testCase.hash != "" {
				exists, err = repo.CCTXExists(context.Background(), testCase.hash)
			} else if testCase.id != "" {
				exists, err = repo.ballotExists(context.Background(), testCase.id)
			} else {
				panic("unreachable")
			}

			require.Equal(testCase.expectedExists, exists)
			for _, expectedErr := range testCase.expectedErrs {
				require.ErrorIs(err, expectedErr)
			}
		})
	}
}
