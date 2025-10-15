package scheduler

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	cometbft "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

func TestScheduler(t *testing.T) {
	t.Run("Basic case", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		exec := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		ts.scheduler.Register(ts.ctx, exec)
		time.Sleep(1500 * time.Millisecond)
		ts.scheduler.Stop()

		// ASSERT
		// Counter should be 2 because we invoke a task once on a start,
		// once after 1 second
		// and then at T=1.5s we stop the scheduler.
		assert.Equal(t, int32(2), counter)

		// Check logs
		assert.Contains(t, ts.logger.String(), "stopped scheduler task")
		assert.Contains(t, ts.logger.String(), `"task_group":"default"`)
	})

	t.Run("More opts", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		exec := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		ts.scheduler.Register(
			ts.ctx,
			exec,
			Name("counter-inc"),
			GroupName("my-custom-group"),
			Interval(300*time.Millisecond),
			LogFields(map[string]any{
				"blockchain": "doge",
				"validators": []string{"alice", "bob"},
			}),
		)

		time.Sleep(time.Second)
		ts.scheduler.Stop()

		// ASSERT
		// Counter should be 1 + 1000/300 = 4 (first run + interval runs)
		assert.Equal(t, int32(4), counter)

		// Also check that log fields are present
		assert.Contains(t, ts.logger.String(), `"task_name":"counter-inc","task_group":"my-custom-group"`)
		assert.Contains(t, ts.logger.String(), `"blockchain":"doge","validators":["alice","bob"]`)
	})

	t.Run("Task can stop itself", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		exec := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		task := ts.scheduler.Register(ts.ctx, exec, Interval(300*time.Millisecond))

		time.Sleep(time.Second)
		task.Stop()

		// ASSERT
		// Counter should be 1 + 1000/300 = 4 (first run + interval runs)
		assert.Equal(t, int32(4), counter)
	})

	t.Run("Skipper option", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		exec := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		const maxValue = 5

		// Skipper function that drops the task after reaching a certain counter value.
		skipper := func() bool {
			allowed := atomic.LoadInt32(&counter) < maxValue
			return !allowed
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		task := ts.scheduler.Register(ts.ctx, exec, Interval(50*time.Millisecond), Skipper(skipper))

		time.Sleep(time.Second)
		task.Stop()

		// ASSERT
		assert.Equal(t, int32(maxValue), counter)
	})

	t.Run("IntervalUpdater option", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		exec := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// Interval updater that increases the interval by 50ms on each counter increment.
		intervalUpdater := func() time.Duration {
			cnt := atomic.LoadInt32(&counter)
			if cnt == 0 {
				return time.Millisecond
			}

			return time.Duration(cnt) * 50 * time.Millisecond
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		task := ts.scheduler.Register(ts.ctx, exec, IntervalUpdater(intervalUpdater))

		time.Sleep(time.Second)
		task.Stop()

		// ASSERT
		assert.Equal(t, int32(6), counter)

		assert.Contains(t, ts.logger.String(), `"ticker.old_interval":0.001,"ticker.new_interval":0.05`)
		assert.Contains(t, ts.logger.String(), `"ticker.old_interval":0.05,"ticker.new_interval":0.1`)
		assert.Contains(t, ts.logger.String(), `"ticker.old_interval":0.1,"ticker.new_interval":0.15`)
		assert.Contains(t, ts.logger.String(), `"ticker.old_interval":0.15,"ticker.new_interval":0.2`)
		assert.Contains(t, ts.logger.String(), `"ticker.old_interval":0.2,"ticker.new_interval":0.25`)
		assert.Contains(t, ts.logger.String(), `"ticker.old_interval":0.25,"ticker.new_interval":0.3`)
	})

	t.Run("Multiple tasks in different groups", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		// Given multiple tasks
		var counterA, counterB, counterC int32

		// Two tasks for Alice
		taskAliceA := func(ctx context.Context) error {
			atomic.AddInt32(&counterA, 1)
			time.Sleep(60 * time.Millisecond)
			return nil
		}

		taskAliceB := func(ctx context.Context) error {
			atomic.AddInt32(&counterB, 1)
			time.Sleep(70 * time.Millisecond)
			return nil
		}

		// One task for Bob
		taskBobC := func(ctx context.Context) error {
			atomic.AddInt32(&counterC, 1)
			time.Sleep(80 * time.Millisecond)
			return nil
		}

		// ACT
		// Register all tasks with different intervals and groups
		ts.scheduler.Register(ts.ctx, taskAliceA, Interval(50*time.Millisecond), GroupName("alice"), Name("a"))
		ts.scheduler.Register(ts.ctx, taskAliceB, Interval(100*time.Millisecond), GroupName("alice"), Name("b"))
		ts.scheduler.Register(ts.ctx, taskBobC, Interval(200*time.Millisecond), GroupName("bob"), Name("c"))

		// Wait and then stop Alice's tasks
		time.Sleep(time.Second)
		ts.scheduler.StopGroup("alice")

		// ASSERT #1
		shutdownLogPattern := func(group, name string) string {
			const pattern = `"task_name":"%s","task_group":"%s",.*"message":"stopped scheduler task"`
			return fmt.Sprintf(pattern, name, group)
		}

		// Make sure Alice.A and Alice.B are stopped
		assert.Regexp(t, shutdownLogPattern("alice", "a"), ts.logger.String())
		assert.Regexp(t, shutdownLogPattern("alice", "b"), ts.logger.String())

		// But Bob.C is still running
		assert.NotRegexp(t, shutdownLogPattern("bob", "c"), ts.logger.String())

		// ACT #2
		time.Sleep(200 * time.Millisecond)
		ts.scheduler.StopGroup("bob")

		// ASSERT #2
		// Bob.C is not running
		assert.Regexp(t, shutdownLogPattern("bob", "c"), ts.logger.String())
	})

	t.Run("Block tick: tick is faster than the block", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		// Given a task that increments a counter by block height
		var counter int64

		task := func(ctx context.Context) error {
			// Note that ctx contains the block event
			blockEvent, ok := BlockFromContext(ctx)
			require.True(t, ok)

			atomic.AddInt64(&counter, blockEvent.Block.Height)
			time.Sleep(100 * time.Millisecond)
			return nil
		}

		// Given block ticker
		blockChan := ts.mockBlockChan(200*time.Millisecond, 0)

		// ACT
		// Register block
		ts.scheduler.Register(ts.ctx, task, BlockTicker(blockChan))
		time.Sleep(1200 * time.Millisecond)
		ts.scheduler.Stop()

		// ASSERT
		assert.Equal(t, int64(21), counter)
		assert.Contains(t, ts.logger.String(), "stopped scheduler task")
		assert.Contains(t, ts.logger.String(), `"task_type":"block_ticker"`)
	})

	t.Run("Block tick: tick is slower than the block", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		// Given a task that increments a counter on start
		// and then decrements before finish
		var counter int64

		exec := func(ctx context.Context) error {
			_, ok := BlockFromContext(ctx)
			require.True(t, ok)

			atomic.AddInt64(&counter, 1)
			time.Sleep(256 * time.Millisecond)
			atomic.AddInt64(&counter, -1)
			return nil
		}

		// Given block ticker
		blockChan := ts.mockBlockChan(100*time.Millisecond, 0)

		// ACT
		// Register block
		ts.scheduler.Register(ts.ctx, exec, BlockTicker(blockChan))
		time.Sleep(1200 * time.Millisecond)
		ts.scheduler.Stop()

		// ASSERT
		// zero indicates that Stop() waits for current iteration to finish (graceful shutdown)
		assert.Equal(t, int64(0), counter)
	})

	t.Run("Block tick: chan closes unexpectedly", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		// Given a task that increments a counter on start
		// and then decrements before finish
		var counter int64

		exec := func(ctx context.Context) error {
			_, ok := BlockFromContext(ctx)
			require.True(t, ok)

			atomic.AddInt64(&counter, 1)
			time.Sleep(200 * time.Millisecond)
			atomic.AddInt64(&counter, -1)
			return nil
		}

		// Given block ticker that closes after 3 blocks
		blockChan := ts.mockBlockChan(100*time.Millisecond, 3)

		// ACT
		// Register block
		ts.scheduler.Register(ts.ctx, exec, BlockTicker(blockChan), Name("block-tick"))

		// Wait for a while
		time.Sleep(1000 * time.Millisecond)

		// Stop the scheduler.
		// Note that actually the ticker is already stopped.
		ts.scheduler.Stop()

		// ASSERT
		// zero indicates that Stop() waits for current iteration to finish (graceful shutdown)
		assert.Equal(t, int64(0), counter)
		assert.Contains(t, ts.logger.String(), "Block channel closed")
	})
}

type testSuite struct {
	ctx       context.Context
	scheduler *Scheduler

	logger *testlog.Log
}

func newTestSuite(t *testing.T) *testSuite {
	logger := testlog.New(t)

	scheduler := New(logger.Logger, time.Second)
	t.Cleanup(scheduler.Stop)

	return &testSuite{
		ctx:       context.Background(),
		scheduler: scheduler,
		logger:    logger,
	}
}

// mockBlockChan mocks websocket blocks. Optionally halts after lastBlock.
func (ts *testSuite) mockBlockChan(interval time.Duration, lastBlock int64) chan cometbft.EventDataNewBlock {
	producer := make(chan cometbft.EventDataNewBlock)

	go func() {
		var blockNumber int64

		for {
			blockNumber++
			ts.logger.Info().Int64("block_number", blockNumber).Msg("Producing new block")

			header := cometbft.Header{
				ChainID: "zeta",
				Height:  blockNumber,
				Time:    time.Now(),
			}

			producer <- cometbft.EventDataNewBlock{
				Block: &cometbft.Block{Header: header},
			}

			if blockNumber > 0 && blockNumber == lastBlock {
				ts.logger.Info().Int64("block_number", blockNumber).Msg("Halting block producer")
				close(producer)
				return
			}

			time.Sleep(interval)
		}
	}()

	return producer
}
