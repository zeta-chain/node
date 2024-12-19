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

type Task func(ctx context.Context) error

type Group string

const DefaultGroup = Group("default")

type Definition struct {
	scheduler *Scheduler

	id     uuid.UUID
	group  Group
	name   string
	task   Task
	ticker *ticker.Ticker

	interval        time.Duration
	intervalUpdater func() time.Duration
	skipper         func() bool

	logFields map[string]any
	logger    zerolog.Logger

	// todo block subscriber (on zeta-chain new block)
}

// New Scheduler instance.
func New(logger zerolog.Logger) *Scheduler {
	return &Scheduler{
		definitions: make(map[uuid.UUID]*Definition),
		logger:      logger.With().Str("module", "scheduler").Logger(),
	}
}

// Opt Definition option
type Opt func(*Definition)

// Name sets task name.
func Name(name string) Opt {
	return func(d *Definition) { d.name = name }
}

func GroupName(group Group) Opt {
	return func(d *Definition) { d.group = group }
}

// LogFields augments definition logger with some fields.
func LogFields(fields map[string]any) Opt {
	return func(d *Definition) { d.logFields = fields }
}

// Interval sets initial task interval.
func Interval(interval time.Duration) Opt {
	return func(d *Definition) { d.interval = interval }
}

// Skipper sets task skipper function
func Skipper(skipper func() bool) Opt {
	return func(d *Definition) { d.skipper = skipper }
}

// IntervalUpdater sets interval updater function.
func IntervalUpdater(intervalUpdater func() time.Duration) Opt {
	return func(d *Definition) { d.intervalUpdater = intervalUpdater }
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

	logOpts := s.logger.With().
		Str("task.name", def.name).
		Str("task.group", string(def.group))

	if len(def.logFields) > 0 {
		logOpts = logOpts.Fields(def.logFields)
	}

	def.logger = logOpts.Logger()

	defTicker := def.provisionTicker(task)

	bgTask := func(ctx context.Context) error {
		return defTicker.Run(ctx)
	}

	s.mu.Lock()
	s.definitions[id] = def
	s.mu.Unlock()

	// Run async worker
	bg.Work(ctx, bgTask, bg.WithLogger(def.logger), bg.WithName(def.name))

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
	d.logger.Info().Msg("Stopping scheduler task")
	d.ticker.StopBlocking()
	d.logger.Info().Dur("time_taken", time.Since(start)).Msg("Stopped scheduler task")

	// delete definition from scheduler
	d.scheduler.mu.Lock()
	delete(d.scheduler.definitions, d.id)
	d.scheduler.mu.Unlock()
}

func (d *Definition) provisionTicker(task Task) *ticker.Ticker {
	d.ticker = ticker.New(
		d.interval,
		d.tickerTask(task),
		ticker.WithLogger(d.logger, d.name),
	)

	return d.ticker
}

// tickerTask wraps Task to be executed by ticker.Ticker
func (d *Definition) tickerTask(task Task) ticker.Task {
	// todo metrics
	//   - duration
	//   - outcome (skip, err, ok)
	//   - bump invocation counter

	return func(ctx context.Context, t *ticker.Ticker) error {
		// skip tick
		if d.skipper != nil && d.skipper() {
			return nil
		}

		err := task(ctx)

		if err != nil {
			d.logger.Error().Err(err).Msg("task failed")
			return nil
		}

		if d.intervalUpdater != nil {
			// noop if interval is not changed
			t.SetInterval(d.intervalUpdater())
		}

		return nil
	}
}
