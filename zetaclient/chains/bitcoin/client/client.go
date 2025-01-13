// Package client represents BTC RPC client.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cosmossdk.io/errors"
	types "github.com/btcsuite/btcd/btcjson"
	chains "github.com/btcsuite/btcd/chaincfg"
	"github.com/rs/zerolog"
	"github.com/tendermint/btcd/btcjson"
	"github.com/tendermint/btcd/chaincfg"

	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
)

type Client struct {
	hostURL string
	client  *http.Client
	config  config.BTCConfig
	params  chains.Params
	logger  zerolog.Logger
}

type Opt func(c *Client)

type rawResponse struct {
	Result json.RawMessage   `json:"result"`
	Error  *btcjson.RPCError `json:"error"`
}

const (
	// v1 means "no batch mode"
	rpcVersion = types.RpcVersion1

	// rpc command id. as we don't send batch requests, it's always 1
	commandID = uint64(1)
)

func WithHTTP(httpClient *http.Client) Opt {
	return func(c *Client) { c.client = httpClient }
}

// New Client constructor
func New(cfg config.BTCConfig, logger zerolog.Logger, opts ...Opt) (*Client, error) {
	params, err := resolveParams(cfg.RPCParams)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve chain params")
	}

	c := &Client{
		hostURL: normalizeHostURL(cfg.RPCHost, true),
		client:  defaultHTTPClient(),
		params:  params,
		logger:  logger.With().Str(logs.FieldModule, "btc_client").Logger(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// send sends RPC command to the server via http post request
func (c *Client) sendCommand(ctx context.Context, cmd any) (json.RawMessage, error) {
	method, reqBody, err := c.marshalCmd(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal cmd")
	}

	// ps: we can add retry logic if needed

	req, err := c.newRequest(ctx, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create http request for %q", method)
	}

	out, err := c.sendRequest(req, method)
	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "%q failed", method)
	case out.Error != nil:
		return nil, errors.Wrapf(out.Error, "got rpc error for %q", method)
	}

	return out.Result, nil
}

func (c *Client) sendRequest(req *http.Request, method string) (out rawResponse, err error) {
	// todo prometheus metrics (chain_id, method, latency, ok/fail)

	c.logger.Debug().Str("rpc.method", method).Msg("Sending request")
	start := time.Now()

	defer func() {
		c.logger.Debug().
			Str("rpc.method", method).
			Dur("rpc.duration", time.Since(start)).
			Err(err).
			Msg("Sent request")
	}()

	res, err := c.client.Do(req)
	if err != nil {
		return rawResponse{}, errors.Wrap(err, "unable to send the request")
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return rawResponse{}, errors.Wrap(err, "unable to read response body")
	}

	if err = json.Unmarshal(resBody, &out); err != nil {
		return rawResponse{}, errors.Wrapf(err, "unable to unmarshal rpc response (%s)", resBody)
	}

	return out, nil
}

func (c *Client) marshalCmd(cmd any) (string, []byte, error) {
	methodName, err := types.CmdMethod(cmd)
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to resolve method")
	}

	body, err := types.MarshalCmd(rpcVersion, commandID, cmd)
	if err != nil {
		return "", nil, errors.Wrapf(err, "unable to marshal cmd %q", methodName)
	}

	return methodName, body, nil
}

func unmarshal[T any](raw json.RawMessage) (T, error) {
	var tt T

	if err := json.Unmarshal(raw, &tt); err != nil {
		return tt, errors.Wrapf(err, "unable to unmarshal to '%T' (%s)", tt, raw)
	}

	return tt, nil
}

func unmarshalPtr[T any](raw json.RawMessage) (*T, error) {
	tt, err := unmarshal[T](raw)
	if err != nil {
		return nil, err
	}

	return &tt, nil
}

func (c *Client) newRequest(ctx context.Context, body []byte) (*http.Request, error) {
	payload := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.hostURL, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.config.RPCPassword != "" || c.config.RPCUsername != "" {
		req.SetBasicAuth(c.config.RPCUsername, c.config.RPCPassword)
	}

	return req, nil
}

func resolveParams(name string) (chains.Params, error) {
	switch name {
	case chains.MainNetParams.Name:
		return chains.MainNetParams, nil
	case chains.TestNet3Params.Name:
		return chains.TestNet3Params, nil
	case chaincfg.RegressionNetParams.Name:
		return chains.RegressionNetParams, nil
	case chaincfg.SimNetParams.Name:
		return chains.SimNetParams, nil
	default:
		return chains.Params{}, fmt.Errorf("unknown chain params %q", name)
	}
}

func normalizeHostURL(host string, disableHTTPS bool) string {
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		return host
	}

	protocol := "http"
	if !disableHTTPS {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s", protocol, host)
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}
}
