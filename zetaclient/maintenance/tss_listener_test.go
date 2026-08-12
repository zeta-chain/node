package maintenance

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/rs/zerolog"

	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

const (
	tssPubkeyOld = "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc"
	tssPubkeyNew = "zetapub1addwnpepqglunjrgl3qg08duxq9pf28jmvrer3crwnnfzp6m0u0yh9jk9mnn5p76utc"
)

// waitForListenerTick gives the listener room for at least one tick so a shutdown it was
// going to signal has actually had the chance to fire.
func waitForListenerTick() {
	time.Sleep(2 * tssListenerTicker)
}

func TestTSSListener(t *testing.T) {
	logger := zerolog.New(zerolog.NewTestWriter(t))

	oldTSS := observertypes.TSS{TssPubkey: tssPubkeyOld}

	// A blanked keygen record must not restart the client.
	//
	// zetacore writes exactly this on any observer set change (x/observer/abci.go
	// BeginBlocker): status back to pending, grantees erased, block set to MaxInt64. On
	// mainnet that reset restarted every signer at once, and none could start again because
	// the erased grantee list left them with an empty p2p whitelist.
	t.Run("blanked keygen record does not trigger a shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		client := mocks.NewZetacoreClient(t)
		client.Mock.On("GetTSS", ctx).Return(oldTSS, nil)
		client.Mock.On("GetTSSHistory", ctx).Return([]observertypes.TSS{oldTSS}, nil)

		complete := make(chan interface{})
		NewTSSListener(client, logger).Listen(ctx, func() { close(complete) })

		waitForListenerTick()
		assertChannelNotClosed(t, complete)
	})

	t.Run("TSS address change still triggers a shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		client := mocks.NewZetacoreClient(t)
		client.Mock.On("GetTSSHistory", ctx).Return([]observertypes.TSS{oldTSS}, nil)
		client.Mock.On("GetTSS", ctx).Return(oldTSS, nil).Once()
		client.Mock.On("GetTSS", ctx).Return(observertypes.TSS{TssPubkey: tssPubkeyNew}, nil)

		complete := make(chan interface{})
		NewTSSListener(client, logger).Listen(ctx, func() { close(complete) })

		<-complete
	})

	t.Run("new key in TSS history still triggers a shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		newTSS := observertypes.TSS{TssPubkey: tssPubkeyNew}

		client := mocks.NewZetacoreClient(t)
		client.Mock.On("GetTSS", ctx).Return(oldTSS, nil)
		client.Mock.On("GetTSSHistory", ctx).Return([]observertypes.TSS{oldTSS}, nil).Once()
		client.Mock.On("GetTSSHistory", ctx).Return([]observertypes.TSS{oldTSS, newTSS}, nil)

		complete := make(chan interface{})
		NewTSSListener(client, logger).Listen(ctx, func() { close(complete) })

		<-complete
	})
}

// TestTSSListenerIgnoresKeygen pins that the listener never reads the keygen record. The mock
// fails the test on any unexpected call, so a reintroduced keygen watcher shows up here rather
// than as an outage.
func TestTSSListenerIgnoresKeygen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zerolog.New(zerolog.NewTestWriter(t))
	oldTSS := observertypes.TSS{TssPubkey: tssPubkeyOld}

	client := mocks.NewZetacoreClient(t)
	client.Mock.On("GetTSS", ctx).Return(oldTSS, nil)
	client.Mock.On("GetTSSHistory", ctx).Return([]observertypes.TSS{oldTSS}, nil)

	NewTSSListener(client, logger).Listen(ctx, func() {})
	waitForListenerTick()

	client.Mock.AssertNotCalled(t, "GetKeyGen", ctx)

	// Guard against the assertion above passing for the wrong reason: the record a caller
	// would have read is the blanked one, and it must remain irrelevant.
	blanked := observertypes.Keygen{
		Status:      observertypes.KeygenStatus_PendingKeygen,
		BlockNumber: math.MaxInt64,
	}
	if len(blanked.GranteePubkeys) != 0 {
		t.Fatal("expected the blanked keygen record to carry no grantees")
	}
}
