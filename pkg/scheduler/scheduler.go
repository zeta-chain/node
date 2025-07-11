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
)

// Scheduler represents background task scheduler.
type Scheduler struct {
	tasks           map[uuid.UUID]*Task
	mu              sync.RWMutex
	logger          zerolog.Logger
	defaultInterval time.Duration
}

// Executable arbitrary function that can be executed.
type Executable func(ctx context.Context) error

// Group represents Task group. Tasks can be grouped for easier management.
type Group string

// DefaultGroup is the default task group.
const DefaultGroup = Group("default")

// tickable ticker abstraction to support different implementations
type tickable interface {
	Start(ctx context.Context) error
	Stop()
}

// Task represents scheduler's task.
type Task struct {
	// ref to the Scheduler is required
	scheduler *Scheduler

	id    uuid.UUID
	group Group
	name  string

	exec Executable

	// ticker abstraction to support different implementations
	ticker  tickable
	skipper func() bool

	logger zerolog.Logger
}

type taskOpts struct {
	interval        time.Duration
	intervalUpdater func() time.Duration

	blockChan <-chan cometbft.EventDataNewBlock

	logFields map[string]any
}

// New Scheduler instance.
func New(logger zerolog.Logger, defaultInterval time.Duration) *Scheduler {
	if defaultInterval <= 0 {
		defaultInterval = time.Second * 10
	}

	return &Scheduler{
		tasks:           make(map[uuid.UUID]*Task),
		logger:          logger.With().Str("module", "scheduler").Logger(),
		defaultInterval: defaultInterval,
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
	}

	config := &taskOpts{
		interval: s.defaultInterval,
	}

	for _, opt := range opts {
		opt(task, config)
	}

	task.logger = newTaskLogger(task, config, s.logger)
	task.ticker = newTickable(task, config)

	task.logger.Info().Msgf("Starting scheduler task %s", task.name)
	bg.Work(ctx, task.ticker.Start, bg.WithLogger(task.logger))

	s.mu.Lock()
	s.tasks[id] = task
	s.mu.Unlock()

	return task
}

func (s *Scheduler) Tasks() map[uuid.UUID]*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copied := make(map[uuid.UUID]*Task, len(s.tasks))
	for k, v := range s.tasks {
		copied[k] = v
	}

	return copied
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

	s.logger.Info().
		Int("tasks", len(selectedTasks)).
		Str("group", string(group)).
		Msg("Stopping scheduler group")

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
	t.logger.Info().Msgf("Stopping scheduler task %s", t.name)
	start := time.Now()

	t.ticker.Stop()

	t.scheduler.mu.Lock()
	delete(t.scheduler.tasks, t.id)
	t.scheduler.mu.Unlock()

	timeTakenMS := time.Since(start).Milliseconds()
	t.logger.Info().Int64("time_taken_ms", timeTakenMS).Msg("Stopped scheduler task")
}

func (t *Task) Group() Group {
	return t.group
}

func (t *Task) Name() string {
	return t.name
}

// execute executes Task with additional logging and metrics.
func (t *Task) execute(ctx context.Context) error {
	startedAt := time.Now().UTC()

	// skip tick
	if t.skipper != nil && t.skipper() {
		recordMetrics(t, startedAt, nil, true)
		return nil
	}

	err := t.exec(ctx)

	recordMetrics(t, startedAt, err, false)

	return err
}

func newTaskLogger(task *Task, opts *taskOpts, logger zerolog.Logger) zerolog.Logger {
	logOpts := logger.With().
		Str("task.name", task.name).
		Str("task.group", string(task.group))

	if len(opts.logFields) > 0 {
		logOpts = logOpts.Fields(opts.logFields)
	}

	taskType := "interval_ticker"
	if opts.blockChan != nil {
		taskType = "block_ticker"
	}

	return logOpts.Str("task.type", taskType).Logger()
}

func newTickable(task *Task, opts *taskOpts) tickable {
	// Block-based ticker
	if opts.blockChan != nil {
		return newBlockTicker(task.execute, opts.blockChan, task.logger)
	}

	return newIntervalTicker(
		task.execute,
		opts.interval,
		opts.intervalUpdater,
		task.name,
		task.logger,
	)
}

// normalizeInterval ensures that the interval is positive to prevent panics.
func normalizeInterval(dur time.Duration) time.Duration {
	if dur > 0 {
		return dur
	}

	return time.Second
}
