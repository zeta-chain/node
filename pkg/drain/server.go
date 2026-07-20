package drain

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/draintx"
)

// PayloadServer serves the latest drain payload over HTTP. It is the local, no-S3 variant
// of the object-store distribution used by the e2e test: Publish swaps the served payload
// atomically, so a later final supersedes earlier drafts. The zetaclient polls its URL.
type PayloadServer struct {
	mu       sync.RWMutex
	payload  []byte
	server   *http.Server
	listener net.Listener
}

// NewPayloadServer creates an empty PayloadServer.
func NewPayloadServer() *PayloadServer {
	return &PayloadServer{}
}

// Publish sets the payload served to subsequent requests.
func (s *PayloadServer) Publish(p draintx.Payload) error {
	b, err := json.Marshal(p)
	if err != nil {
		return errors.Wrap(err, "unable to marshal payload")
	}
	s.mu.Lock()
	s.payload = b
	s.mu.Unlock()
	return nil
}

// ServeHTTP writes the current payload, or 503 if none has been published yet.
func (s *PayloadServer) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	b := s.payload
	s.mu.RUnlock()

	if b == nil {
		http.Error(w, "no payload published", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

// Start binds addr (use ":0" for an ephemeral port) and serves in the background.
func (s *PayloadServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "unable to listen")
	}
	s.listener = listener
	s.server = &http.Server{Handler: s, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = s.server.Serve(listener) }()
	return nil
}

// URL returns the base URL of the running server.
func (s *PayloadServer) URL() string {
	if s.listener == nil {
		return ""
	}
	return "http://" + s.listener.Addr().String()
}

// Close shuts the server down.
func (s *PayloadServer) Close() error {
	if s.server == nil {
		return nil
	}
	return s.server.Close()
}

// Generator produces a payload for the given sequence and finality.
type Generator func(ctx context.Context, seq uint64, final bool) (draintx.Payload, error)

// RunCron republishes fresh draft payloads on an interval, then publishes exactly one
// final payload once isFinalTime reports true, and stops. Drafts are for monitoring only;
// clients only ever sign the final. This is the thin wrapper around the generator.
func RunCron(
	ctx context.Context,
	interval time.Duration,
	gen Generator,
	publish func(draintx.Payload) error,
	isFinalTime func(ctx context.Context) (bool, error),
) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var seq uint64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			final, err := isFinalTime(ctx)
			if err != nil {
				return errors.Wrap(err, "unable to check final time")
			}
			p, err := gen(ctx, seq, final)
			if err != nil {
				return errors.Wrap(err, "unable to generate payload")
			}
			if err := publish(p); err != nil {
				return errors.Wrap(err, "unable to publish payload")
			}
			seq++
			if final {
				return nil
			}
		}
	}
}
