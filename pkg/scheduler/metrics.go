package scheduler

import (
	"time"

	"github.com/zeta-chain/node/zetaclient/metrics"
)

// Note that currently the hard-coded "global" metrics are used.
func recordMetrics(task *Task, startedAt time.Time, err error, skipped bool) {
	var status string
	switch {
	case skipped:
		status = "skipped"
	case err != nil:
		status = "failed"
	default:
		status = "ok"
	}

	var (
		group = string(task.group)
		name  = task.name
		dur   = time.Since(startedAt).Seconds()
	)

	metrics.SchedulerTaskInvocationCounter.WithLabelValues(status, group, name).Inc()
	metrics.SchedulerTaskExecutionDuration.WithLabelValues(status, group, name).Observe(dur)
}
