package scheduler

import (
	"time"

	cometbft "github.com/cometbft/cometbft/types"
)

// Opt Definition option
type Opt func(*Definition)

// Name sets task name.
func Name(name string) Opt {
	return func(d *Definition) { d.name = name }
}

// GroupName sets task group. Otherwise, defaults to DefaultGroup.
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

// BlockTicker makes Definition to listen for new zeta blocks instead of using interval ticker.
// IntervalUpdater is ignored.
func BlockTicker(blocks <-chan cometbft.EventDataNewBlock) Opt {
	return func(d *Definition) { d.blockChan = blocks }
}
