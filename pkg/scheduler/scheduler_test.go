package scheduler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	t.Run("Basic case", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		task := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		ts.scheduler.Register(ts.ctx, task)
		time.Sleep(1500 * time.Millisecond)
		ts.scheduler.Stop()

		// ASSERT
		// Counter should be 2 because we invoke a task once on a start,
		// once after 1 second (default interval),
		// and then at T=1.5s we stop the scheduler.
		assert.Equal(t, int32(2), counter)

		// Check logs
		assert.Contains(t, ts.logBuffer.String(), "Stopped task")
		assert.Contains(t, ts.logBuffer.String(), `"task.group":"default"`)
	})

	t.Run("More opts", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		task := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		ts.scheduler.Register(
			ts.ctx,
			task,
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
		assert.Contains(t, ts.logBuffer.String(), `"task.name":"counter-inc","task.group":"my-custom-group"`)
		assert.Contains(t, ts.logBuffer.String(), `"blockchain":"doge","validators":["alice","bob"]`)
	})

	t.Run("Definition can also stop itself", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		task := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		def := ts.scheduler.Register(ts.ctx, task, Interval(300*time.Millisecond))

		time.Sleep(time.Second)
		def.Stop()

		// ASSERT
		// Counter should be 1 + 1000/300 = 4 (first run + interval runs)
		assert.Equal(t, int32(4), counter)
	})

	t.Run("Skipper option", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		task := func(ctx context.Context) error {
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
		def := ts.scheduler.Register(ts.ctx, task, Interval(50*time.Millisecond), Skipper(skipper))

		time.Sleep(time.Second)
		def.Stop()

		// ASSERT
		assert.Equal(t, int32(maxValue), counter)
	})

	t.Run("IntervalUpdater option", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t)

		var counter int32

		task := func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// Interval updater that increases the interval by 50ms on each counter increment.
		intervalUpdater := func() time.Duration {
			return time.Duration(atomic.LoadInt32(&counter)) * 50 * time.Millisecond
		}

		// ACT
		// Register task and stop it after x1.5 interval.
		def := ts.scheduler.Register(ts.ctx, task, Interval(time.Millisecond), IntervalUpdater(intervalUpdater))

		time.Sleep(time.Second)
		def.Stop()

		// ASSERT
		assert.Equal(t, int32(6), counter)

		assert.Contains(t, ts.logBuffer.String(), `"ticker.old_interval":1,"ticker.new_interval":50`)
		assert.Contains(t, ts.logBuffer.String(), `"ticker.old_interval":50,"ticker.new_interval":100`)
		assert.Contains(t, ts.logBuffer.String(), `"ticker.old_interval":100,"ticker.new_interval":150`)
		assert.Contains(t, ts.logBuffer.String(), `"ticker.old_interval":150,"ticker.new_interval":200`)
		assert.Contains(t, ts.logBuffer.String(), `"ticker.old_interval":200,"ticker.new_interval":250`)
		assert.Contains(t, ts.logBuffer.String(), `"ticker.old_interval":250,"ticker.new_interval":300`)
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
			return fmt.Sprintf(
				`"task\.name":"%s","task\.group":"%s","time_taken_ms":.*"message":"Stopped task"`,
				name,
				group,
			)
		}

		// Make sure Alice.A and Alice.B are stopped
		assert.Regexp(t, shutdownLogPattern("alice", "a"), ts.logBuffer.String())
		assert.Regexp(t, shutdownLogPattern("alice", "b"), ts.logBuffer.String())

		// But Bob.C is still running
		assert.NotRegexp(t, shutdownLogPattern("bob", "c"), ts.logBuffer.String())

		// ACT #2
		time.Sleep(200 * time.Millisecond)
		ts.scheduler.StopGroup("bob")

		// ASSERT #2
		// Bob.C is not running
		assert.Regexp(t, shutdownLogPattern("bob", "c"), ts.logBuffer.String())
	})
}

type testSuite struct {
	ctx       context.Context
	scheduler *Scheduler

	logger    zerolog.Logger
	logBuffer *bytes.Buffer
}

func newTestSuite(t *testing.T) *testSuite {
	logBuffer := &bytes.Buffer{}
	logger := zerolog.New(io.MultiWriter(zerolog.NewTestWriter(t), logBuffer))

	return &testSuite{
		ctx:       context.Background(),
		scheduler: New(logger),
		logger:    logger,
		logBuffer: logBuffer,
	}
}
