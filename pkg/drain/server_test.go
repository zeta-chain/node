package drain_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/draintx"
)

func fetch(t *testing.T, url string) (draintx.Payload, int) {
	resp, err := http.Get(url) //nolint:noctx
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		return draintx.Payload{}, resp.StatusCode
	}
	var p draintx.Payload
	require.NoError(t, json.Unmarshal(body, &p))
	return p, resp.StatusCode
}

func TestPayloadServerPublishAndServe(t *testing.T) {
	// ARRANGE
	srv := drain.NewPayloadServer()
	require.NoError(t, srv.Start("127.0.0.1:0"))
	defer srv.Close()

	// no payload yet -> 503
	_, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusServiceUnavailable, status)

	// ACT: publish a draft, then a final that supersedes it
	require.NoError(t, srv.Publish(draintx.Payload{Seq: 1, Final: false, TriggerZetaHeight: 100}))
	draft, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusOK, status)
	require.False(t, draft.Final)
	require.EqualValues(t, 1, draft.Seq)

	require.NoError(t, srv.Publish(draintx.Payload{Seq: 2, Final: true, TriggerZetaHeight: 100}))
	final, status := fetch(t, srv.URL())

	// ASSERT
	require.Equal(t, http.StatusOK, status)
	require.True(t, final.Final)
	require.EqualValues(t, 2, final.Seq)
}

func TestRunCronPublishesDraftsThenFinal(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	srv := drain.NewPayloadServer()
	require.NoError(t, srv.Start("127.0.0.1:0"))
	defer srv.Close()

	gen := func(_ context.Context, seq uint64, final bool) (draintx.Payload, error) {
		return drain.BuildPayload(100, seq, final, nil, nil, priv)
	}
	// become final on the 3rd tick
	ticks := 0
	isFinalTime := func(context.Context) (bool, error) {
		ticks++
		return ticks >= 3, nil
	}

	// ACT
	err = drain.RunCron(context.Background(), time.Millisecond, gen, srv.Publish, isFinalTime)

	// ASSERT: cron stops after publishing the final; the served payload is final
	require.NoError(t, err)
	final, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusOK, status)
	require.True(t, final.Final)
	require.EqualValues(t, 2, final.Seq) // seq 0,1 drafts; seq 2 final
	require.NoError(t, final.Verify(ethcrypto.CompressPubkey(&priv.PublicKey)))
}
