// Package scheduler provides a background task scheduler that allows for the registration,
// execution, and management of periodic tasks. Tasks can be grouped, named, and configured
// with various options such as custom intervals, log fields, and skip conditions.
//
// The scheduler supports dynamic interval updates and can gracefully stop tasks either
// individually or by group.
package scheduler

import (
	"context"
	"sync"
	"time"

	cometbft "github.com/cometbft/cometbft/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/ticker"
)

// Scheduler represents background task scheduler.
type Scheduler struct {
	tasks  map[uuid.UUID]*Task
	mu     sync.RWMutex
	logger zerolog.Logger
}

// Executable arbitrary function that can be executed.
type Executable func(ctx context.Context) error

// Group represents Task group. Tasks can be grouped for easier management.
type Group string

// DefaultGroup is the default task group.
const DefaultGroup = Group("default")

// Task represents scheduler's task.
type Task struct {
	// ref to the Scheduler is required
	scheduler *Scheduler

	// naming stuff
	id    uuid.UUID
	group Group
	name  string

	exec Executable

	// represents interval ticker and its options
	ticker          *ticker.Ticker
	interval        time.Duration
	intervalUpdater func() time.Duration
	skipper         func() bool

	// zeta block ticker (also supports skipper)
	blockChan       <-chan cometbft.EventDataNewBlock
	blockChanTicker *blockTicker

	// logging
	logFields map[string]any
	logger    zerolog.Logger
}

// New Scheduler instance.
func New(logger zerolog.Logger) *Scheduler {
	return &Scheduler{
		tasks:  make(map[uuid.UUID]*Task),
		logger: logger.With().Str("module", "scheduler").Logger(),
	}
}

// Register registers and starts new Task in the background
func (s *Scheduler) Register(ctx context.Context, exec Executable, opts ...Opt) *Task {
	id := uuid.New()
	task := &Task{
		scheduler: s,
		id:        id,
		group:     DefaultGroup,
		name:      id.String(),
		exec:      exec,
		interval:  time.Second,
	}
	for _, opt := range opts {
		opt(task)
	}

	task.logger = newTaskLogger(task, s.logger)

	task.startTicker(ctx)

	s.mu.Lock()
	s.tasks[id] = task
	s.mu.Unlock()

	return task
}

// Stop stops all tasks.
func (s *Scheduler) Stop() {
	s.StopGroup("")
}

// StopGroup stops all tasks in the group.
func (s *Scheduler) StopGroup(group Group) {
	var selectedTasks []*Task

	s.mu.RLock()

	// Filter desired tasks
	for _, task := range s.tasks {
		// "" is for wildcard i.e. all groups
		if group == "" || task.group == group {
			selectedTasks = append(selectedTasks, task)
		}
	}

	s.mu.RUnlock()

	if len(selectedTasks) == 0 {
		return
	}

	// Stop all selected tasks concurrently
	var wg sync.WaitGroup
	wg.Add(len(selectedTasks))

	for _, task := range selectedTasks {
		go func(task *Task) {
			defer wg.Done()
			task.Stop()
		}(task)
	}

	wg.Wait()
}

// Stop stops the task and offloads it from the scheduler.
func (t *Task) Stop() {
	start := time.Now()

	// delete task from scheduler
	defer func() {
		t.scheduler.mu.Lock()
		delete(t.scheduler.tasks, t.id)
		t.scheduler.mu.Unlock()

		timeTakenMS := time.Since(start).Milliseconds()
		t.logger.Info().Int64("time_taken_ms", timeTakenMS).Msg("Stopped scheduler task")
	}()

	t.logger.Info().Msg("Stopping scheduler task")

	if t.isIntervalTicker() {
		t.ticker.StopBlocking()
		return
	}

	t.blockChanTicker.Stop()
}

func (t *Task) isIntervalTicker() bool {
	return t.blockChan == nil
}

func (t *Task) startTicker(ctx context.Context) {
	t.logger.Info().Msg("Starting scheduler task")

	if t.isIntervalTicker() {
		t.ticker = ticker.New(t.interval, t.invokeByInterval, ticker.WithLogger(t.logger, t.name))
		bg.Work(ctx, t.ticker.Start, bg.WithLogger(t.logger))

		return
	}

	t.blockChanTicker = newBlockTicker(t.invoke, t.blockChan, t.logger)

	bg.Work(ctx, t.blockChanTicker.Start, bg.WithLogger(t.logger))
}

// invokeByInterval a ticker.Task wrapper of invoke.
func (t *Task) invokeByInterval(ctx context.Context, tt *ticker.Ticker) error {
	if err := t.invoke(ctx); err != nil {
		t.logger.Error().Err(err).Msg("task failed")
	}

	if t.intervalUpdater != nil {
		// noop if interval is not changed
		tt.SetInterval(t.intervalUpdater())
	}

	return nil
}

// invoke executes a given Task with logging & telemetry.
func (t *Task) invoke(ctx context.Context) error {
	// skip tick
	if t.skipper != nil && t.skipper() {
		return nil
	}

	t.logger.Debug().Msg("Invoking task")

	err := t.exec(ctx)

	// todo metrics (TBD)
	//   - duration (time taken)
	//   - outcome (skip, err, ok)
	//   - bump invocation counter
	//   - "last invoked at" timestamp (?)
	//   - chain_id
	//   - metrics cardinality: "task_group (?)" "task_name", "status", "chain_id"

	return err
}

func newTaskLogger(task *Task, logger zerolog.Logger) zerolog.Logger {
	logOpts := logger.With().
		Str("task.name", task.name).
		Str("task.group", string(task.group))

	if len(task.logFields) > 0 {
		logOpts = logOpts.Fields(task.logFields)
	}

	taskType := "interval_ticker"
	if task.blockChanTicker != nil {
		taskType = "block_ticker"
	}

	return logOpts.Str("task.type", taskType).Logger()
}
