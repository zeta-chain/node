package signer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/solana/signer"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_NewSigner(t *testing.T) {
	// test parameters
	chain := chains.SolanaDevnet
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.GatewayAddresses[chain.ChainId]

	tests := []struct {
		name        string
		chain       chains.Chain
		chainParams observertypes.ChainParams
		solClient   interfaces.SolanaRPCClient
		tss         interfaces.TSSSigner
		relayerKey  *keys.RelayerKey
		ts          *metrics.TelemetryServer
		logger      base.Logger
		errMessage  string
	}{
		{
			name:        "should create solana signer successfully with relayer key",
			chain:       chain,
			chainParams: *chainParams,
			solClient:   nil,
			tss:         nil,
			relayerKey: &keys.RelayerKey{
				PrivateKey: "3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ",
			},
			ts:     nil,
			logger: base.DefaultLogger(),
		},
		{
			name:        "should create solana signer successfully without relayer key",
			chainParams: *chainParams,
			solClient:   nil,
			tss:         nil,
			relayerKey:  nil,
			ts:          nil,
			logger:      base.DefaultLogger(),
		},
		{
			name: "should fail to create solana signer with invalid gateway address",
			chainParams: func() observertypes.ChainParams {
				cp := *chainParams
				cp.GatewayAddress = "invalid"
				return cp
			}(),
			solClient:  nil,
			tss:        nil,
			relayerKey: nil,
			ts:         nil,
			logger:     base.DefaultLogger(),
			errMessage: "cannot parse gateway address",
		},
		{
			name:        "should fail to create solana signer with invalid relayer key",
			chainParams: *chainParams,
			solClient:   nil,
			tss:         nil,
			relayerKey: &keys.RelayerKey{
				PrivateKey: "3EMjCcCJg53fMEGVj13", // too short
			},
			ts:         nil,
			logger:     base.DefaultLogger(),
			errMessage: "unable to construct solana private key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := signer.NewSigner(tt.chain, tt.chainParams, tt.solClient, tt.tss, tt.relayerKey, tt.ts, tt.logger)
			if tt.errMessage != "" {
				require.ErrorContains(t, err, tt.errMessage)
				require.Nil(t, s)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, s)
		})
	}
}

func Test_SetRelayerBalanceMetrics(t *testing.T) {
	// test parameters
	chain := chains.SolanaDevnet
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.GatewayAddresses[chain.ChainId]
	relayerKey := &keys.RelayerKey{
		PrivateKey: "3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ",
	}
	ctx := context.Background()

	// mock solana client with RPC error
	mckClient := mocks.NewSolanaRPCClient(t)
	mckClient.On("GetBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("rpc error"))

	// create signer and set relayer balance metrics
	s, err := signer.NewSigner(chain, *chainParams, mckClient, nil, relayerKey, nil, base.DefaultLogger())
	require.NoError(t, err)
	s.SetRelayerBalanceMetrics(ctx)

	// assert that relayer key balance metrics is not set (due to RPC error)
	balance := testutil.ToFloat64(metrics.RelayerKeyBalance.WithLabelValues(chain.Name))
	require.Equal(t, 0.0, balance)

	// mock solana client with balance
	mckClient = mocks.NewSolanaRPCClient(t)
	mckClient.On("GetBalance", mock.Anything, mock.Anything, mock.Anything).Return(&rpc.GetBalanceResult{
		Value: 123400000,
	}, nil)

	// create signer and set relayer balance metrics again
	s, err = signer.NewSigner(chain, *chainParams, mckClient, nil, relayerKey, nil, base.DefaultLogger())
	require.NoError(t, err)
	s.SetRelayerBalanceMetrics(ctx)

	// assert that relayer key balance metrics is set correctly
	balance = testutil.ToFloat64(metrics.RelayerKeyBalance.WithLabelValues(chain.Name))
	require.Equal(t, 0.1234, balance)
}
