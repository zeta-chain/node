package rpc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/rpc/stream"
	"github.com/zeta-chain/node/server/config"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/client"
)

func newTestWebsocketServer() *websocketsServer {
	// dummy values for testing
	cfg := &config.Config{}
	cfg.JSONRPC.Address = "localhost:9999"   // not used
	cfg.JSONRPC.WsAddress = "localhost:9999" // not used
	cfg.TLS.CertificatePath = ""
	cfg.TLS.KeyPath = ""

	return &websocketsServer{
		rpcAddr:        cfg.JSONRPC.Address,
		wsAddr:         cfg.JSONRPC.WsAddress,
		certFile:       cfg.TLS.CertificatePath,
		keyFile:        cfg.TLS.KeyPath,
		api:            newPubSubAPI(client.Context{}, log.NewNopLogger(), &stream.RPCStream{}),
		logger:         log.NewNopLogger(),
		allowedOrigins: []string{"*"},
	}
}

func TestWebsocketPayloadLimit(t *testing.T) {
	srv := newTestWebsocketServer()

	ts := httptest.NewServer(srv)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	u.Scheme = "ws"

	dialer := websocket.Dialer{}
	conn, httpResp, err := dialer.Dial(u.String(), nil)
	require.NotNil(t, httpResp)
	require.NoError(t, err)

	defer conn.Close()

	// Send oversized message (2 MB)
	oversizedPayload := make([]byte, 2<<20)
	_ = conn.WriteMessage(websocket.TextMessage, oversizedPayload)

	// The connection should close
	_, _, readErr := conn.ReadMessage()
	require.Error(t, readErr, "expected connection to close on oversized message")
}

func TestCheckOrigin(t *testing.T) {
	logger := log.NewNopLogger()
	tests := []struct {
		name           string
		allowedOrigins []string
		requestOrigin  string
		expected       bool
	}{
		{
			name:           "empty allowed origins - should reject",
			allowedOrigins: []string{},
			requestOrigin:  "https://example.com",
			expected:       false,
		},
		{
			name:           "allowed origin - should accept",
			allowedOrigins: []string{"localhost", "127.0.0.1", "example.com"},
			requestOrigin:  "https://example.com",
			expected:       true,
		},
		{
			name:           "not allowed origin - should reject",
			allowedOrigins: []string{"localhost", "127.0.0.1"},
			requestOrigin:  "https://malicious.com",
			expected:       false,
		},
		{
			name:           "wildcard origin - should accept",
			allowedOrigins: []string{"*"},
			requestOrigin:  "https://anything.com",
			expected:       true,
		},
		{
			name:           "empty origin header - should accept",
			allowedOrigins: []string{"localhost"},
			requestOrigin:  "",
			expected:       true,
		},
		{
			name:           "localhost origin - should accept",
			allowedOrigins: []string{"localhost", "127.0.0.1"},
			requestOrigin:  "http://localhost:3000",
			expected:       true,
		},
		{
			name:           "127.0.0.1 origin - should accept",
			allowedOrigins: []string{"localhost", "127.0.0.1"},
			requestOrigin:  "http://127.0.0.1:8080",
			expected:       true,
		},
		{
			name:           "invalid origin URL - should reject",
			allowedOrigins: []string{"localhost"},
			requestOrigin:  "invalid-url",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &websocketsServer{
				allowedOrigins: tt.allowedOrigins,
				logger:         logger,
			}

			req := &http.Request{
				Header: make(http.Header),
			}

			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}

			result := server.checkOrigin(req)
			if result != tt.expected {
				t.Errorf("checkOrigin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeOriginForLogging(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal origin - should pass through",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "origin with newlines - should be stripped",
			input:    "https://example.com\nmalicious-log-entry",
			expected: "https://example.commalicious-log-entry",
		},
		{
			name:     "origin with carriage return - should be stripped",
			input:    "https://example.com\rmalicious-log-entry",
			expected: "https://example.commalicious-log-entry",
		},
		{
			name:     "origin with tab - should be stripped",
			input:    "https://example.com\tmalicious-log-entry",
			expected: "https://example.commalicious-log-entry",
		},
		{
			name:     "origin with control characters - should be stripped",
			input:    "https://example.com\x00\x1f\x7f",
			expected: "https://example.com",
		},
		{
			name:     "very long origin - should be truncated",
			input:    "https://example.com/" + strings.Repeat("a", 300),
			expected: "", // Will be checked separately below
		},
		{
			name:     "mostly non-printable characters - should use placeholder",
			input:    "\x00\x01\x02\x03",
			expected: "<sanitized-origin>",
		},
		{
			name:     "empty string - should use placeholder",
			input:    "",
			expected: "<sanitized-origin>",
		},
		{
			name:     "origin with unicode - should be stripped",
			input:    "https://example.com/测试",
			expected: "https://example.com/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeOriginForLogging(tt.input)

			// Special handling for the truncation test
			if tt.name == "very long origin - should be truncated" {
				if len(result) != 203 || !strings.HasSuffix(result, "...") || !strings.HasPrefix(result, "https://example.com/") {
					t.Errorf("sanitizeOriginForLogging() for long input: got length %d, want 203 with prefix and suffix", len(result))
				}
				return
			}

			if result != tt.expected {
				t.Errorf("sanitizeOriginForLogging(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
