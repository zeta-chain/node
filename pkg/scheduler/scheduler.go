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
	definitions map[uuid.UUID]*Definition
	mu          sync.RWMutex
	logger      zerolog.Logger
}

// Task represents scheduler's task
type Task func(ctx context.Context) error

// Group represents Definition group.
// Definitions can be grouped for easier management.
type Group string

// DefaultGroup is the default group for definitions.
const DefaultGroup = Group("default")

// Definition represents a configuration of a Task
type Definition struct {
	// ref to the Scheduler is required
	scheduler *Scheduler

	// naming stuff
	id    uuid.UUID
	group Group
	name  string

	// arbitrary function that will be invoked by the scheduler
	task Task

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
		definitions: make(map[uuid.UUID]*Definition),
		logger:      logger.With().Str("module", "scheduler").Logger(),
	}
}

// Register registers and starts new task in the background
func (s *Scheduler) Register(ctx context.Context, task Task, opts ...Opt) *Definition {
	id := uuid.New()
	def := &Definition{
		scheduler: s,
		id:        id,
		group:     DefaultGroup,
		name:      id.String(),
		task:      task,
		interval:  time.Second,
	}
	for _, opt := range opts {
		opt(def)
	}

	def.logger = newDefinitionLogger(def, s.logger)

	def.startTicker(ctx)

	s.mu.Lock()
	s.definitions[id] = def
	s.mu.Unlock()

	return def
}

// Stop stops all tasks.
func (s *Scheduler) Stop() {
	s.StopGroup("")
}

// StopGroup stops all tasks in the group.
func (s *Scheduler) StopGroup(group Group) {
	var selectedDefs []*Definition

	s.mu.RLock()

	// Filter desired definitions
	for _, def := range s.definitions {
		// "" is for wildcard i.e. all groups
		if group == "" || def.group == group {
			selectedDefs = append(selectedDefs, def)
		}
	}

	s.mu.RUnlock()

	if len(selectedDefs) == 0 {
		return
	}

	// Stop all selected tasks concurrently
	var wg sync.WaitGroup
	wg.Add(len(selectedDefs))

	for _, def := range selectedDefs {
		go func(def *Definition) {
			defer wg.Done()
			def.Stop()
		}(def)
	}

	wg.Wait()
}

// Stop stops the task and offloads it from the scheduler.
func (d *Definition) Stop() {
	start := time.Now()

	// delete definition from scheduler
	defer func() {
		d.scheduler.mu.Lock()
		delete(d.scheduler.definitions, d.id)
		d.scheduler.mu.Unlock()

		timeTakenMS := time.Since(start).Milliseconds()
		d.logger.Info().Int64("time_taken_ms", timeTakenMS).Msg("Stopped scheduler task")
	}()

	d.logger.Info().Msg("Stopping scheduler task")

	if d.isIntervalTicker() {
		d.ticker.StopBlocking()
		return
	}

	d.blockChanTicker.Stop()
}

func (d *Definition) isIntervalTicker() bool {
	return d.blockChan == nil
}

func (d *Definition) startTicker(ctx context.Context) {
	d.logger.Info().Msg("Starting scheduler task")

	if d.isIntervalTicker() {
		d.ticker = ticker.New(d.interval, d.invokeByInterval, ticker.WithLogger(d.logger, d.name))
		bg.Work(ctx, d.ticker.Start, bg.WithLogger(d.logger))

		return
	}

	d.blockChanTicker = newBlockTicker(d.invoke, d.blockChan, d.logger)

	bg.Work(ctx, d.blockChanTicker.Start, bg.WithLogger(d.logger))
}

// invokeByInterval a ticker.Task wrapper of invoke.
func (d *Definition) invokeByInterval(ctx context.Context, t *ticker.Ticker) error {
	if err := d.invoke(ctx); err != nil {
		d.logger.Error().Err(err).Msg("task failed")
	}

	if d.intervalUpdater != nil {
		// noop if interval is not changed
		t.SetInterval(d.intervalUpdater())
	}

	return nil
}

// invoke executes a given Task with logging & telemetry.
func (d *Definition) invoke(ctx context.Context) error {
	// skip tick
	if d.skipper != nil && d.skipper() {
		return nil
	}

	d.logger.Debug().Msg("Invoking task")

	err := d.task(ctx)

	// todo metrics (TBD)
	//   - duration (time taken)
	//   - outcome (skip, err, ok)
	//   - bump invocation counter
	//   - "last invoked at" timestamp (?)
	//   - chain_id
	//   - metrics cardinality: "task_group (?)" "task_name", "status", "chain_id"

	return err
}

func newDefinitionLogger(def *Definition, logger zerolog.Logger) zerolog.Logger {
	logOpts := logger.With().
		Str("task.name", def.name).
		Str("task.group", string(def.group))

	if len(def.logFields) > 0 {
		logOpts = logOpts.Fields(def.logFields)
	}

	taskType := "interval_ticker"
	if def.blockChanTicker != nil {
		taskType = "block_ticker"
	}

	return logOpts.Str("task.type", taskType).Logger()
}
