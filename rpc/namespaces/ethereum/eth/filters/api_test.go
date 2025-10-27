package filters

import (
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func TestTimeoutLoop_PanicOnNilCancel(t *testing.T) {
	api := &PublicFilterAPI{
		filters:   make(map[rpc.ID]*filter),
		filtersMu: sync.Mutex{},
		deadline:  10 * time.Millisecond,
	}
	api.filters[rpc.NewID()] = &filter{
		typ:      filters.BlocksSubscription,
		deadline: time.NewTimer(0),
	}
	done := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("cancel panic")
			}
			close(done)
		}()
		api.timeoutLoop()
	}()
	panicked := false
	select {
	case <-done:
		panicked = true
	case <-time.After(100 * time.Millisecond):
	}
	require.False(t, panicked)
}
