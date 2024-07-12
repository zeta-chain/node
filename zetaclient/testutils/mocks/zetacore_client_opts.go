package mocks

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	keyinterfaces "github.com/zeta-chain/zetacore/zetaclient/keys/interfaces"
)

var errSomethingIsWrong = errors.New("oopsie")

// Note that this is NOT codegen but a handwritten mock improvement.

func (_m *ZetacoreClient) WithKeys(keys keyinterfaces.ObserverKeys) *ZetacoreClient {
	_m.On("GetKeys").Maybe().Return(keys)

	return _m
}

func (_m *ZetacoreClient) WithZetaChain() *ZetacoreClient {
	_m.On("Chain").Maybe().Return(chains.ZetaChainMainnet)

	return _m
}

func (_m *ZetacoreClient) WithPostVoteOutbound(zetaTxHash string, ballotIndex string) *ZetacoreClient {
	_m.On("PostVoteOutbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(zetaTxHash, ballotIndex, nil)

	return _m
}

func (_m *ZetacoreClient) WithPostVoteInbound(zetaTxHash string, ballotIndex string) *ZetacoreClient {
	_m.On("PostVoteInbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(zetaTxHash, ballotIndex, nil)

	return _m
}

func (_m *ZetacoreClient) WithRateLimiterFlags(flags *crosschaintypes.RateLimiterFlags) *ZetacoreClient {
	on := _m.On("GetRateLimiterFlags", mock.Anything).Maybe()
	if flags != nil {
		on.Return(*flags, nil)
	} else {
		on.Return(crosschaintypes.RateLimiterFlags{}, errSomethingIsWrong)
	}

	return _m
}

func (_m *ZetacoreClient) WithRateLimiterInput(in *crosschaintypes.QueryRateLimiterInputResponse) *ZetacoreClient {
	on := _m.On("GetRateLimiterInput", mock.Anything, mock.Anything).Maybe()
	if in != nil {
		on.Return(in, nil)
	} else {
		on.Return(nil, errSomethingIsWrong)
	}

	return _m
}

func (_m *ZetacoreClient) WithPendingCctx(chainID int64, cctxs []*crosschaintypes.CrossChainTx) *ZetacoreClient {
	totalPending := uint64(len(cctxs))

	_m.On("ListPendingCCTX", mock.Anything, chainID).Maybe().Return(cctxs, totalPending, nil)

	return _m
}
