package scheduler

import (
	"time"

	cometbft "github.com/cometbft/cometbft/types"
)

// Opt Task option
type Opt func(task *Task)

// Name sets task name.
func Name(name string) Opt {
	return func(d *Task) { d.name = name }
}

// GroupName sets task group. Otherwise, defaults to DefaultGroup.
func GroupName(group Group) Opt {
	return func(d *Task) { d.group = group }
}

// LogFields augments Task's logger with some fields.
func LogFields(fields map[string]any) Opt {
	return func(d *Task) { d.logFields = fields }
}

// Interval sets initial task interval.
func Interval(interval time.Duration) Opt {
	return func(d *Task) { d.interval = interval }
}

// Skipper sets task skipper function
func Skipper(skipper func() bool) Opt {
	return func(d *Task) { d.skipper = skipper }
}

// IntervalUpdater sets interval updater function.
func IntervalUpdater(intervalUpdater func() time.Duration) Opt {
	return func(d *Task) { d.intervalUpdater = intervalUpdater }
}

// BlockTicker makes Definition to listen for new zeta blocks instead of using interval ticker.
// IntervalUpdater is ignored.
func BlockTicker(blocks <-chan cometbft.EventDataNewBlock) Opt {
	return func(d *Task) { d.blockChan = blocks }
}
