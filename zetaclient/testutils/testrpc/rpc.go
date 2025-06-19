package testrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// RPCHandler is a function that handles a JSON RPC request.
type RPCHandler func(params map[string]any) (any, error)

// Server represents JSON RPC mock with a "real" HTTP server allocated (httptest)
type Server struct {
	t *testing.T

	// method => args_key => handler
	// This allows for having different mocks for the same method based on the input provided
	// Kinda similar to mockery.
	handlers map[string]map[string]RPCHandler
	name     string
}

// New constructs Server.
func New(t *testing.T, name string) (*Server, string) {
	var (
		handlers = make(map[string]map[string]RPCHandler)
		rpc      = &Server{t, handlers, name}
		testWeb  = httptest.NewServer(http.HandlerFunc(rpc.httpHandler))
	)

	t.Cleanup(testWeb.Close)

	return rpc, testWeb.URL
}

// On registers a handler for a given method and optional parameters.
// If params is provided, it registers a parameter-specific handler.
func (s *Server) On(method string, call RPCHandler, params ...any) {
	if s.handlers[method] == nil {
		s.handlers[method] = make(map[string]RPCHandler)
	}

	paramKey := paramsKeyFromArray(s.t, params)
	s.handlers[method][paramKey] = call
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
	methodHandlers, ok := s.handlers[req.Method]
	if !ok {
		return rpcResponse{Error: errors.New("method not found")}
	}

	// build param key
	paramKey := paramsKeyFromMap(s.t, req.Params)

	// look for parameter-specific handler
	call, ok := methodHandlers[paramKey]
	if !ok {
		// look for default handler
		call, ok = methodHandlers[""]
		if !ok {
			return rpcResponse{Error: errors.New("no handler found")}
		}
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

// paramsKeyFromArray creates a key for the given parameter array.
func paramsKeyFromArray(t *testing.T, params []any) string {
	if len(params) == 0 {
		return ""
	}

	// marshal params to string
	paramBytes, err := json.Marshal(params)
	require.NoError(t, err)

	return string(paramBytes)
}

// paramsKeyFromMap creates a key for the given parameter map.
func paramsKeyFromMap(t *testing.T, params map[string]any) string {
	if len(params) == 0 {
		return ""
	}

	// convert map to ordered array
	// example: {"0": "foo", "1": "bar", ...}
	paramsOrdered := make([]any, len(params))
	for i := range paramsOrdered {
		param, ok := params[fmt.Sprintf("%d", i)]
		require.True(t, ok, "param %d not found", i)
		paramsOrdered[i] = param
	}

	return paramsKeyFromArray(t, paramsOrdered)
}
