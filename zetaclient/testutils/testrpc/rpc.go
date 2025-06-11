package testrpc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// Server represents JSON RPC mock with a "real" HTTP server allocated (httptest)
type Server struct {
	t        *testing.T
	handlers map[string]func(params map[string]any) (any, error)
	name     string
}

// New constructs Server.
func New(t *testing.T, name string) (*Server, string) {
	var (
		handlers = make(map[string]func(params map[string]any) (any, error))
		rpc      = &Server{t, handlers, name}
		testWeb  = httptest.NewServer(http.HandlerFunc(rpc.httpHandler))
	)

	t.Cleanup(testWeb.Close)

	return rpc, testWeb.URL
}

// On registers a handler for a given method.
func (s *Server) On(method string, call func(params map[string]any) (any, error)) {
	s.handlers[method] = call
}

// example: {"jsonrpc":"1.0","method":"ping","params":{},"id":1}
// also supports array params: {"jsonrpc":"1.0","method":"ping","params":[1,2,3],"id":1}.
// the latter would be casted to map["$idx"]any{"0": "foo" ,...}
type rpcRequest struct {
	Method string         `json:"method"`
	Params map[string]any `json:"params"`
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

	// Decode request
	raw, err := io.ReadAll(r.Body)
	require.NoError(s.t, err)

	req, err := parseRequest(raw)
	require.NoError(s.t, err, "unable to unmarshal request for %s [%s]", s.name, string(raw))

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

func parseRequest(raw []byte) (rpcRequest, error) {
	items := gjson.GetManyBytes(raw, "method", "params")

	if !items[0].Exists() || items[0].Type != gjson.String {
		return rpcRequest{}, errors.New("method string is expected")
	}

	req := rpcRequest{
		Method: items[0].String(),
		Params: map[string]any{},
	}

	// .params is optional
	if !items[1].Exists() {
		return req, nil
	}

	items[1].ForEach(func(key, value gjson.Result) bool {
		req.Params[key.String()] = value.Value()

		return true
	})

	return req, nil
}
