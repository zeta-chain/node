package scheduler

import (
	"time"

	cometbft "github.com/cometbft/cometbft/types"
)

// Opt Task option
type Opt func(task *Task, taskOpts *taskOpts)

// Name sets task name.
func Name(name string) Opt {
	return func(t *Task, _ *taskOpts) { t.name = name }
}

// GroupName sets task group. Otherwise, defaults to DefaultGroup.
func GroupName(group Group) Opt {
	return func(t *Task, _ *taskOpts) { t.group = group }
}

// LogFields augments Task's logger with some fields.
func LogFields(fields map[string]any) Opt {
	return func(_ *Task, opts *taskOpts) { opts.logFields = fields }
}

// Interval sets initial task interval.
func Interval(interval time.Duration) Opt {
	return func(_ *Task, opts *taskOpts) { opts.interval = normalizeInterval(interval) }
}

// Skipper sets task skipper function. If it returns true, the task is skipped.
func Skipper(skipper func() bool) Opt {
	return func(t *Task, _ *taskOpts) { t.skipper = skipper }
}

// IntervalUpdater sets interval updater function. Overrides Interval.
func IntervalUpdater(intervalUpdater func() time.Duration) Opt {
	return func(_ *Task, opts *taskOpts) {
		opts.interval = normalizeInterval(intervalUpdater())
		opts.intervalUpdater = intervalUpdater
	}
}

// BlockTicker makes Task to listen for new zeta blocks
// instead of using interval ticker. IntervalUpdater is ignored.
func BlockTicker(blocks <-chan cometbft.EventDataNewBlock) Opt {
	return func(_ *Task, opts *taskOpts) { opts.blockChan = blocks }
}
