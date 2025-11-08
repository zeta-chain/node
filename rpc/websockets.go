package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	rpcfilters "github.com/zeta-chain/node/rpc/namespaces/ethereum/eth/filters"
	"github.com/zeta-chain/node/rpc/stream"
	"github.com/zeta-chain/node/server/config"
)

const (
	maxMessageSize = 1 << 20 // 1 MiB is the max message size for the websocket server
)

type WebsocketsServer interface {
	Start()
}

type SubscriptionResponseJSON struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      float64     `json:"id"`
}

type SubscriptionNotification struct {
	Jsonrpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  *SubscriptionResult `json:"params"`
}

type SubscriptionResult struct {
	Subscription rpc.ID      `json:"subscription"`
	Result       interface{} `json:"result"`
}

type ErrorResponseJSON struct {
	Jsonrpc string            `json:"jsonrpc"`
	Error   *ErrorMessageJSON `json:"error"`
	ID      *big.Int          `json:"id"`
}

type ErrorMessageJSON struct {
	Code    *big.Int `json:"code"`
	Message string   `json:"message"`
}

type websocketsServer struct {
	rpcAddr        string // listen address of rest-server
	wsAddr         string // listen address of ws server
	certFile       string
	keyFile        string
	allowedOrigins []string // allowed origins for WebSocket connections
	api            *pubSubAPI
	logger         log.Logger
}

func NewWebsocketsServer(
	clientCtx client.Context,
	logger log.Logger,
	stream *stream.RPCStream,
	cfg *config.Config,
) WebsocketsServer {
	logger = logger.With("api", "websocket-server")
	return &websocketsServer{
		rpcAddr:        cfg.JSONRPC.Address,
		wsAddr:         cfg.JSONRPC.WsAddress,
		certFile:       cfg.TLS.CertificatePath,
		keyFile:        cfg.TLS.KeyPath,
		allowedOrigins: cfg.JSONRPC.WSOrigins,
		api:            newPubSubAPI(clientCtx, logger, stream),
		logger:         logger,
	}
}

func (s *websocketsServer) Start() {
	ws := mux.NewRouter()
	ws.Handle("/", s)

	go func() {
		var err error
		if s.certFile == "" || s.keyFile == "" {
			//#nosec G114 -- http functions have no support for timeouts
			err = http.ListenAndServe(s.wsAddr, ws)
		} else {
			//#nosec G114 -- http functions have no support for timeouts
			err = http.ListenAndServeTLS(s.wsAddr, s.certFile, s.keyFile, ws)
		}

		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			s.logger.Error("failed to start HTTP server for WS", "error", err.Error())
		}
	}()
}

// sanitizeOriginForLogging sanitizes the origin header to prevent log injection attacks
func sanitizeOriginForLogging(origin string) string {
	// Limit length to prevent log flooding
	if len(origin) > 200 {
		origin = origin[:200] + "..."
	}

	// Remove or replace dangerous characters that could be used for log injection
	// Replace newlines, carriage returns, and other control characters
	sanitized := regexp.MustCompile(`[\r\n\t\x00-\x1f\x7f-\x9f]`).ReplaceAllString(origin, "")

	// Additional safety: only allow printable ASCII and common URL characters
	sanitized = regexp.MustCompile(`[^\x20-\x7E]`).ReplaceAllString(sanitized, "")

	// If the result is empty or too different from original, use a safe placeholder
	if sanitized == "" || len(sanitized) < len(origin)/2 {
		return "<sanitized-origin>"
	}

	// Replace newlines, carriage returns, and other control characters
	sanitized = strings.ReplaceAll(sanitized, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")

	// Escape the input to prevent HTML injection
	sanitized = html.EscapeString(sanitized)

	return sanitized
}

// checkOrigin validates the Origin header of incoming WebSocket upgrade requests
func (s *websocketsServer) checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	sanitizedOrigin := sanitizeOriginForLogging(origin)

	// If no allowed origins are configured, reject all requests for security
	if len(s.allowedOrigins) == 0 {
		s.logger.Debug("websocket connection rejected: no allowed origins configured", "origin", sanitizedOrigin)
		return false
	}

	// Allow requests without an Origin header (e.g., from server-side clients)
	if origin == "" {
		return true
	}

	// Parse the origin URL to get the host
	originURL, err := url.Parse(origin)
	if err != nil {
		s.logger.Debug(
			"websocket connection rejected: invalid origin URL",
			"origin",
			sanitizedOrigin,
			"error",
			err.Error(),
		)
		return false
	}

	originHost := originURL.Hostname()

	// Check if the origin host is in the allowed list
	for _, allowedOrigin := range s.allowedOrigins {
		if originHost == allowedOrigin || allowedOrigin == "*" {
			return true
		}
	}

	s.logger.Debug(
		"websocket connection rejected: origin not allowed",
		"origin",
		sanitizedOrigin,
		"allowed",
		s.allowedOrigins,
	)
	return false
}

func (s *websocketsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: s.checkOrigin,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Debug("websocket upgrade failed", "error", err.Error())
		return
	}

	conn.SetReadLimit(maxMessageSize)

	ws := &wsConn{
		mux:  new(sync.Mutex),
		conn: conn,
	}

	s.readLoop(ws)
}

func (s *websocketsServer) sendErrResponse(wsConn *wsConn, msg string) {
	res := &ErrorResponseJSON{
		Jsonrpc: "2.0",
		Error: &ErrorMessageJSON{
			Code:    big.NewInt(-32600),
			Message: msg,
		},
		ID: nil,
	}

	_ = wsConn.WriteJSON(res) // #nosec G703
}

type wsConn struct {
	conn *websocket.Conn
	mux  *sync.Mutex
}

func (w *wsConn) WriteJSON(v interface{}) error {
	w.mux.Lock()
	defer w.mux.Unlock()

	return w.conn.WriteJSON(v)
}

func (w *wsConn) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()

	return w.conn.Close()
}

func (w *wsConn) ReadMessage() (messageType int, p []byte, err error) {
	// not protected by write mutex

	return w.conn.ReadMessage()
}

func (s *websocketsServer) readLoop(wsConn *wsConn) {
	// subscriptions of current connection
	subscriptions := make(map[rpc.ID]context.CancelFunc)
	defer func() {
		// cancel all subscriptions when connection closed
		// #nosec G705
		for _, unsubFn := range subscriptions {
			unsubFn()
		}
	}()

readLoop:
	for {
		_, mb, err := wsConn.ReadMessage()
		if err != nil {
			_ = wsConn.Close() // #nosec G703
			s.logger.Error("read message error, breaking read loop", "error", err.Error())
			return
		}

		if isBatch(mb) {
			if err := s.tcpGetAndSendResponse(wsConn, mb); err != nil {
				s.sendErrResponse(wsConn, err.Error())
			}
			continue
		}

		var msg map[string]interface{}
		if err = json.Unmarshal(mb, &msg); err != nil {
			s.sendErrResponse(wsConn, err.Error())
			continue
		}

		// check if method == eth_subscribe or eth_unsubscribe
		method, ok := msg["method"].(string)
		if !ok {
			// otherwise, call the usual rpc server to respond
			if err := s.tcpGetAndSendResponse(wsConn, mb); err != nil {
				s.sendErrResponse(wsConn, err.Error())
			}

			continue
		}

		var connID float64
		switch id := msg["id"].(type) {
		case string:
			connID, err = strconv.ParseFloat(id, 64)
		case float64:
			connID = id
		default:
			err = fmt.Errorf("unknown type")
		}
		if err != nil {
			s.sendErrResponse(
				wsConn,
				fmt.Errorf("invalid type for connection ID: %T", msg["id"]).Error(),
			)
			continue
		}

		switch method {
		case "eth_subscribe":
			params, ok := s.getParamsAndCheckValid(msg, wsConn)
			if !ok {
				continue
			}

			subID := rpc.NewID()
			unsubFn, err := s.api.subscribe(wsConn, subID, params)
			if err != nil {
				s.sendErrResponse(wsConn, err.Error())
				continue
			}
			subscriptions[subID] = unsubFn

			res := &SubscriptionResponseJSON{
				Jsonrpc: "2.0",
				ID:      connID,
				Result:  subID,
			}

			if err := wsConn.WriteJSON(res); err != nil {
				s.logger.Error("error writing subscription response", "error", err.Error())
				break readLoop
			}
		case "eth_unsubscribe":
			params, ok := s.getParamsAndCheckValid(msg, wsConn)
			if !ok {
				continue
			}

			id, ok := params[0].(string)
			if !ok {
				s.sendErrResponse(wsConn, "invalid parameters")
				continue
			}

			subID := rpc.ID(id)
			unsubFn, ok := subscriptions[subID]
			if ok {
				delete(subscriptions, subID)
				unsubFn()
			}

			res := &SubscriptionResponseJSON{
				Jsonrpc: "2.0",
				ID:      connID,
				Result:  ok,
			}

			if err := wsConn.WriteJSON(res); err != nil {
				s.logger.Error("error writing unsubscribe response", "error", err.Error())
				break readLoop
			}
		default:
			// otherwise, call the usual rpc server to respond
			if err := s.tcpGetAndSendResponse(wsConn, mb); err != nil {
				s.sendErrResponse(wsConn, err.Error())
			}
		}
	}
}

// tcpGetAndSendResponse sends error response to client if params is invalid
func (s *websocketsServer) getParamsAndCheckValid(msg map[string]interface{}, wsConn *wsConn) ([]interface{}, bool) {
	params, ok := msg["params"].([]interface{})
	if !ok {
		s.sendErrResponse(wsConn, "invalid parameters")
		return nil, false
	}

	if len(params) == 0 {
		s.sendErrResponse(wsConn, "empty parameters")
		return nil, false
	}

	return params, true
}

// tcpGetAndSendResponse connects to the rest-server over tcp, posts a JSON-RPC request, and sends the response
// to the client over websockets
func (s *websocketsServer) tcpGetAndSendResponse(wsConn *wsConn, mb []byte) error {
	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://"+s.rpcAddr, bytes.NewBuffer(mb))
	if err != nil {
		return errors.Wrap(err, "Could not build request")
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Could not perform request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body from response")
	}

	var wsSend interface{}
	err = json.Unmarshal(body, &wsSend)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal rest-server response")
	}

	return wsConn.WriteJSON(wsSend)
}

// pubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
type pubSubAPI struct {
	events    *stream.RPCStream
	logger    log.Logger
	clientCtx client.Context
}

// newPubSubAPI creates an instance of the ethereum PubSub API.
func newPubSubAPI(clientCtx client.Context, logger log.Logger, stream *stream.RPCStream) *pubSubAPI {
	logger = logger.With("module", "websocket-client")
	return &pubSubAPI{
		events:    stream,
		logger:    logger,
		clientCtx: clientCtx,
	}
}

func (api *pubSubAPI) subscribe(wsConn *wsConn, subID rpc.ID, params []interface{}) (context.CancelFunc, error) {
	method, ok := params[0].(string)
	if !ok {
		return nil, errors.New("invalid parameters")
	}

	switch method {
	case "newHeads":
		// TODO: handle extra params
		return api.subscribeNewHeads(wsConn, subID)
	case "logs":
		if len(params) > 1 {
			return api.subscribeLogs(wsConn, subID, params[1])
		}
		return api.subscribeLogs(wsConn, subID, nil)
	case "newPendingTransactions":
		return api.subscribePendingTransactions(wsConn, subID)
	case "syncing":
		return api.subscribeSyncing(wsConn, subID)
	default:
		return nil, errors.Errorf("unsupported method %s", method)
	}
}

func (api *pubSubAPI) subscribeNewHeads(wsConn *wsConn, subID rpc.ID) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	//nolint: errcheck
	go api.events.HeaderStream().Subscribe(ctx, func(headers []stream.RPCHeader, _ int) error {
		for _, header := range headers {
			// write to ws conn
			res := &SubscriptionNotification{
				Jsonrpc: "2.0",
				Method:  "eth_subscription",
				Params: &SubscriptionResult{
					Subscription: subID,
					Result:       header.EthHeader,
				},
			}

			if err := wsConn.WriteJSON(res); err != nil {
				api.logger.Error("error writing header, will drop peer", "error", err.Error())

				try(func() {
					if err != websocket.ErrCloseSent {
						_ = wsConn.Close()
					}
				}, api.logger, "closing websocket peer sub")
				return err
			}
		}
		return nil
	})

	return cancel, nil
}

func try(fn func(), l log.Logger, desc string) {
	defer func() {
		if x := recover(); x != nil {
			if err, ok := x.(error); ok {
				// debug.PrintStack()
				l.Debug("panic during "+desc, "error", err.Error())
				return
			}

			l.Debug(fmt.Sprintf("panic during %s: %+v", desc, x))
			return
		}
	}()

	fn()
}

func (api *pubSubAPI) subscribeLogs(wsConn *wsConn, subID rpc.ID, extra interface{}) (context.CancelFunc, error) {
	crit := filters.FilterCriteria{}

	if extra != nil {
		params, ok := extra.(map[string]interface{})
		if !ok {
			err := errors.New("invalid criteria")
			api.logger.Debug("invalid criteria", "type", fmt.Sprintf("%T", extra))
			return nil, err
		}

		if params["address"] != nil {
			switch address := params["address"].(type) {
			case string:
				crit.Addresses = []common.Address{common.HexToAddress(address)}
			case []any:
				for _, addr := range address {
					address, ok := addr.(string)
					if !ok {
						return nil, errors.New("invalid address")
					}

					crit.Addresses = append(crit.Addresses, common.HexToAddress(address))
				}
			default:
				return nil, errors.New("invalid addresses; must be address or array of addresses")
			}
		}

		if params["topics"] != nil {
			topics, ok := params["topics"].([]interface{})
			if !ok {
				err := errors.Errorf("invalid topics: %s", topics)
				api.logger.Error("invalid topics", "type", fmt.Sprintf("%T", topics))
				return nil, err
			}

			crit.Topics = make([][]common.Hash, len(topics))

			addCritTopic := func(topicIdx int, topic interface{}) error {
				tstr, ok := topic.(string)
				if !ok {
					err := errors.Errorf("invalid topic: %s", topic)
					api.logger.Error("invalid topic", "type", fmt.Sprintf("%T", topic))
					return err
				}

				crit.Topics[topicIdx] = []common.Hash{common.HexToHash(tstr)}
				return nil
			}

			for topicIdx, subtopics := range topics {
				if subtopics == nil {
					continue
				}

				// in case we don't have list, but a single topic value
				if topic, ok := subtopics.(string); ok {
					if err := addCritTopic(topicIdx, topic); err != nil {
						return nil, err
					}

					continue
				}

				// in case we actually have a list of subtopics
				subtopicsList, ok := subtopics.([]interface{})
				if !ok {
					err := errors.New("invalid subtopics")
					api.logger.Error("invalid subtopic", "type", fmt.Sprintf("%T", subtopics))
					return nil, err
				}

				subtopicsCollect := make([]common.Hash, len(subtopicsList))
				for idx, subtopic := range subtopicsList {
					tstr, ok := subtopic.(string)
					if !ok {
						err := errors.Errorf("invalid subtopic: %s", subtopic)
						api.logger.Error("invalid subtopic", "type", fmt.Sprintf("%T", subtopic))
						return nil, err
					}

					subtopicsCollect[idx] = common.HexToHash(tstr)
				}

				crit.Topics[topicIdx] = subtopicsCollect
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	//nolint: errcheck
	go api.events.LogStream().Subscribe(ctx, func(txLogs []*ethtypes.Log, _ int) error {
		logs := rpcfilters.FilterLogs(txLogs, crit.FromBlock, crit.ToBlock, crit.Addresses, crit.Topics)
		if len(logs) == 0 {
			return nil
		}

		for _, ethLog := range logs {
			res := &SubscriptionNotification{
				Jsonrpc: "2.0",
				Method:  "eth_subscription",
				Params: &SubscriptionResult{
					Subscription: subID,
					Result:       ethLog,
				},
			}

			err := wsConn.WriteJSON(res)
			if err != nil {
				api.logger.Error("error writing header, will drop peer", "error", err.Error())
				try(func() {
					if err != websocket.ErrCloseSent {
						_ = wsConn.Close()
					}
				}, api.logger, "closing websocket peer sub")

				return err
			}
		}
		return nil
	})

	return cancel, nil
}

func (api *pubSubAPI) subscribePendingTransactions(wsConn *wsConn, subID rpc.ID) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	//nolint: errcheck
	go api.events.PendingTxStream().Subscribe(ctx, func(items []common.Hash, _ int) error {
		for _, hash := range items {
			// write to ws conn
			res := &SubscriptionNotification{
				Jsonrpc: "2.0",
				Method:  "eth_subscription",
				Params: &SubscriptionResult{
					Subscription: subID,
					Result:       hash,
				},
			}

			err := wsConn.WriteJSON(res)
			if err != nil {
				api.logger.Debug("error writing header, will drop peer", "error", err.Error())

				try(func() {
					if err != websocket.ErrCloseSent {
						_ = wsConn.Close()
					}
				}, api.logger, "closing websocket peer sub")
				return err
			}
		}
		return nil
	})

	return cancel, nil
}

func (api *pubSubAPI) subscribeSyncing(_ *wsConn, _ rpc.ID) (context.CancelFunc, error) {
	return nil, errors.New("syncing subscription is not implemented")
}

// copy from github.com/ethereum/go-ethereum/rpc/json.go
// isBatch returns true when the first non-whitespace characters is '['
func isBatch(raw []byte) bool {
	for _, c := range raw {
		// skip insignificant whitespace (http://www.ietf.org/rfc/rfc4627.txt)
		if c == 0x20 || c == 0x09 || c == 0x0a || c == 0x0d {
			continue
		}
		return c == '['
	}
	return false
}
