package testrpc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// Server represents JSON RPC mock with a "real" HTTP server allocated (httptest)
type Server struct {
	t        *testing.T
	handlers map[string]func(params []any) (any, error)
	name     string
}

// New constructs Server.
func New(t *testing.T, name string) (*Server, string) {
	var (
		handlers = make(map[string]func(params []any) (any, error))
		rpc      = &Server{t, handlers, name}
		testWeb  = httptest.NewServer(http.HandlerFunc(rpc.httpHandler))
	)

	t.Cleanup(testWeb.Close)

	return rpc, testWeb.URL
}

// On registers a handler for a given method.
func (s *Server) On(method string, call func(params []any) (any, error)) {
	s.handlers[method] = call
}

// example: {"jsonrpc":"1.0","method":"ping","params":[],"id":1}
type rpcRequest struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}

// example: {"result":0,"error":null,"id":"curltest"}
type rpcResponse struct {
	Result any   `json:"result"`
	Error  error `json:"error"`
}

// handler is a simple HTTP handler that returns 200 OK.
// Later we can add any logic here.
func (s *Server) httpHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure method matches
	require.Equal(s.t, http.MethodPost, r.Method)

	var req rpcRequest

	// Decode request
	raw, err := io.ReadAll(r.Body)
	require.NoError(s.t, err)
	require.NoError(s.t, json.Unmarshal(raw, &req), "unable to unmarshal request")

	// Process request
	res := s.rpcHandler(req)

	// Encode response
	response, err := json.Marshal(res)
	require.NoError(s.t, err, "unable to marshal response")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	require.NoError(s.t, err, "unable to write response")

	s.t.Logf("%s RPC: incoming request: %+v; response: %+v", s.name, req, res)
}

func (s *Server) rpcHandler(req rpcRequest) rpcResponse {
	call, ok := s.handlers[req.Method]
	if !ok {
		return rpcResponse{Error: errors.New("method not found")}
	}

	res, err := call(req.Params)

	return rpcResponse{Result: res, Error: err}
}
