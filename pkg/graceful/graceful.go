// Package graceful contains tools for graceful shutdown.
// GS refers to the process of shutting down a system in a controlled manner, allowing it to complete ongoing tasks,
// release resources, and perform necessary cleanup operations before terminating.
// This ensures that the system stops functioning without causing data loss, corruption, or other issues.
package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Process represents "virtual" process that contains
// routines that can be started and stopped
// Note that ANY failure in starting a service will cause the process to shutdown
type Process struct {
	stop      <-chan os.Signal
	stopStack []func()

	timeout time.Duration
	mu      sync.Mutex
	stopped bool

	logger zerolog.Logger
}

// Service represents abstract service.
type Service interface {
	Start(ctx context.Context) error
	Stop()
}

// DefaultProcess is a process instance with some sane defaults.
var DefaultProcess = New(15*time.Second, log.Logger, NewSigChan(syscall.SIGINT, syscall.SIGTERM))

// New Process constructor.
func New(timeout time.Duration, logger zerolog.Logger, stop <-chan os.Signal) *Process {
	return &Process{
		stop:    stop,
		timeout: timeout,
		logger:  logger.With().Str("module", "graceful").Logger(),
	}
}

// AddService adds Service to the process.
func (p *Process) AddService(ctx context.Context, s Service) {
	p.AddStarter(ctx, s.Start)
	p.AddStopper(s.Stop)
}

// AddStarter runs a function that starts something.
// This is a blocking call for blocking .Start() services
func (p *Process) AddStarter(ctx context.Context, fn func(ctx context.Context) error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				p.logger.Error().Err(panicToErr(r, 10)).Msg("panic in service")
				p.ShutdownNow()
			}
		}()

		if err := fn(ctx); err != nil {
			p.logger.Error().Err(err).Msg("failed to start service")
			p.ShutdownNow()
		}
	}()
}

// AddStopper adds a function will be executed during shutdown.
func (p *Process) AddStopper(fn func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stopStack = append(p.stopStack, fn)
}

// WaitForShutdown blocks current routine until a shutdown signal is received
func (p *Process) WaitForShutdown() {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case sig := <-p.stop:
			p.logger.Info().Msgf("Received signal: %q", sig.String())
			p.ShutdownNow()
			return
		case <-t.C:
			// another goroutine already called ShutdownNow
			// safe to read w/o mutex
			if p.stopped {
				return
			}
		}
	}
}

// ShutdownNow invokes shutdown of all services.
func (p *Process) ShutdownNow() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// noop
	if p.stopped {
		return
	}

	defer func() {
		p.stopped = true
	}()

	p.logger.Info().Msg("Shutting down")

	deadline := time.After(p.timeout)
	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				p.logger.Error().Err(panicToErr(r, 10)).Msg("panic during shutdown")
			}

			// complete shutdown
			close(done)
		}()

		// stop services in the reverse order
		for i := len(p.stopStack) - 1; i >= 0; i-- {
			p.stopStack[i]()
		}
	}()

	select {
	case <-deadline:
		p.logger.Info().Msgf("Shutdown interrupted by timeout (%s)", p.timeout.String())
	case <-done:
		p.logger.Info().Msg("Shutdown completed")
	}
}

// panicToErr converts panic to error WITH exact line of panic.
// Note the offset should be determined empirically.
func panicToErr(panic any, offset int) error {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	line := ""

	if len(lines) > offset {
		line = strings.TrimSpace(lines[offset])
	}

	return fmt.Errorf("panic: %v at %s", panic, line)
}

// NewSigChan creates a new signal channel.
func NewSigChan(signals ...os.Signal) chan os.Signal {
	out := make(chan os.Signal, 1)
	signal.Notify(out, signals...)

	return out
}

func AddService(ctx context.Context, s Service) {
	DefaultProcess.AddService(ctx, s)
}

func AddStarter(ctx context.Context, fn func(ctx context.Context) error) {
	DefaultProcess.AddStarter(ctx, fn)
}

func AddStopper(fn func()) {
	DefaultProcess.AddStopper(fn)
}

func WaitForShutdown() {
	DefaultProcess.WaitForShutdown()
}

func ShutdownNow() {
	DefaultProcess.ShutdownNow()
}
