package drain_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync/atomic"
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

// failOnCronError is the RunCron onError hook for tests that expect every tick to succeed.
func failOnCronError(t *testing.T) func(error) {
	return func(err error) { t.Errorf("unexpected cron error: %v", err) }
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
		return drain.BuildPayload(100, seq, final, drain.NetworkLocalnet, nil, nil, priv)
	}
	// become final on the 3rd tick
	ticks := 0
	isFinalTime := func(context.Context) (bool, error) {
		ticks++
		return ticks >= 3, nil
	}

	// ACT
	err = drain.RunCron(context.Background(), time.Millisecond, gen, srv.Publish, isFinalTime, failOnCronError(t))

	// ASSERT: cron stops after publishing the final; the served payload is final
	require.NoError(t, err)
	final, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusOK, status)
	require.True(t, final.Final)
	require.EqualValues(t, 2, final.Seq) // seq 0,1 drafts; seq 2 final
	require.NoError(t, final.Verify(ethcrypto.CompressPubkey(&priv.PublicKey)))
}

func TestRunCronPublishesImmediately(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	srv := drain.NewPayloadServer()
	require.NoError(t, srv.Start("127.0.0.1:0"))
	defer srv.Close()

	var published int32
	publish := func(p draintx.Payload) error {
		atomic.AddInt32(&published, 1)
		return srv.Publish(p)
	}
	gen := func(_ context.Context, seq uint64, final bool) (draintx.Payload, error) {
		return drain.BuildPayload(100, seq, final, drain.NetworkLocalnet, nil, nil, priv)
	}
	isFinalTime := func(context.Context) (bool, error) { return false, nil }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ACT: a long interval means the only publish that can happen is the immediate first one
	done := make(chan error, 1)
	go func() { done <- drain.RunCron(ctx, time.Hour, gen, publish, isFinalTime, failOnCronError(t)) }()

	// ASSERT: a draft is served without waiting a full interval
	require.Eventually(t, func() bool {
		_, status := fetch(t, srv.URL())
		return status == http.StatusOK
	}, time.Second, 5*time.Millisecond)

	draft, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusOK, status)
	require.False(t, draft.Final)
	require.EqualValues(t, 0, draft.Seq)
	require.EqualValues(t, 1, atomic.LoadInt32(&published)) // exactly the immediate publish, no tick

	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
}

func TestRunCronImmediateFinal(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	srv := drain.NewPayloadServer()
	require.NoError(t, srv.Start("127.0.0.1:0"))
	defer srv.Close()

	var published int32
	publish := func(p draintx.Payload) error {
		atomic.AddInt32(&published, 1)
		return srv.Publish(p)
	}
	gen := func(_ context.Context, seq uint64, final bool) (draintx.Payload, error) {
		return drain.BuildPayload(100, seq, final, drain.NetworkLocalnet, nil, nil, priv)
	}
	// already past the freeze window: the immediate publish is the single final, no tick needed
	isFinalTime := func(context.Context) (bool, error) { return true, nil }

	// ACT
	err = drain.RunCron(context.Background(), time.Hour, gen, publish, isFinalTime, failOnCronError(t))

	// ASSERT: cron stops after the one final publish
	require.NoError(t, err)
	final, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusOK, status)
	require.True(t, final.Final)
	require.EqualValues(t, 0, final.Seq) // final published immediately as seq 0
	require.EqualValues(t, 1, atomic.LoadInt32(&published))
}

func TestRunCronSurvivesDraftError(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	srv := drain.NewPayloadServer()
	require.NoError(t, srv.Start("127.0.0.1:0"))
	defer srv.Close()

	gen := func(_ context.Context, seq uint64, final bool) (draintx.Payload, error) {
		return drain.BuildPayload(100, seq, final, drain.NetworkLocalnet, nil, nil, priv)
	}
	// zetacore is down for the immediate publish and the first tick, then a draft lands and the
	// freeze window opens on the 4th check
	checks := 0
	isFinalTime := func(context.Context) (bool, error) {
		checks++
		if checks <= 2 {
			return false, errors.New("zetacore unavailable")
		}
		return checks >= 4, nil
	}
	var reported []error
	onError := func(err error) { reported = append(reported, err) }

	// ACT
	err = drain.RunCron(context.Background(), time.Millisecond, gen, srv.Publish, isFinalTime, onError)

	// ASSERT: the draft failures were reported but not fatal, and the final still published
	require.NoError(t, err)
	require.Len(t, reported, 2)
	final, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusOK, status)
	require.True(t, final.Final)
	require.EqualValues(t, 1, final.Seq) // seq 0 the one draft that landed; seq 1 the final
	require.NoError(t, final.Verify(ethcrypto.CompressPubkey(&priv.PublicKey)))
}

func TestRunCronRetriesFailedFinal(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	srv := drain.NewPayloadServer()
	require.NoError(t, srv.Start("127.0.0.1:0"))
	defer srv.Close()

	gen := func(_ context.Context, seq uint64, final bool) (draintx.Payload, error) {
		return drain.BuildPayload(100, seq, final, drain.NetworkLocalnet, nil, nil, priv)
	}
	// already past the freeze window, so every attempt is a final
	isFinalTime := func(context.Context) (bool, error) { return true, nil }
	// the first two finals fail to publish, the third lands
	var attempts int32
	publish := func(p draintx.Payload) error {
		if atomic.AddInt32(&attempts, 1) <= 2 {
			return errors.New("publish failed")
		}
		return srv.Publish(p)
	}
	var reported []error
	onError := func(err error) { reported = append(reported, err) }

	// ACT
	err = drain.RunCron(context.Background(), time.Millisecond, gen, publish, isFinalTime, onError)

	// ASSERT: a failed final is retried instead of giving up, and seq only advances on success
	require.NoError(t, err)
	require.EqualValues(t, 3, atomic.LoadInt32(&attempts))
	require.Len(t, reported, 2)
	final, status := fetch(t, srv.URL())
	require.Equal(t, http.StatusOK, status)
	require.True(t, final.Final)
	require.EqualValues(t, 0, final.Seq)
}
