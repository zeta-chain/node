// Package client implements a Bitcoin RPC with support for context, logging, and metrics.
//
// Portions of this package in `./commands.go` are derived from or inspired by btcd,
// which is licensed under the ISC License.
//
// # ISC License
//
// Copyright (c) 2013-2024 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	types "github.com/btcsuite/btcd/btcjson"
	chains "github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tendermint/btcd/chaincfg"

	pkgchains "github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// Client Bitcoin RPC client
type Client struct {
	hostURL    string
	client     *http.Client
	clientName string
	isRegnet   bool
	config     config.BTCConfig
	params     chains.Params
	logger     zerolog.Logger
}

type Opt func(c *Client)

type rawResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *types.RPCError `json:"error"`
}

const (
	// v1 means "no batch mode"
	rpcVersion = types.RpcVersion1

	// rpc command id. as we don't send batch requests, it's always 1
	commandID = uint64(1)
)

var _ client = (*Client)(nil)

func WithHTTP(httpClient *http.Client) Opt {
	return func(c *Client) { c.client = httpClient }
}

// New Client constructor
func New(cfg config.BTCConfig, chainID int64, logger zerolog.Logger, opts ...Opt) (*Client, error) {
	params, err := resolveParams(cfg.RPCParams)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve chain params")
	}

	var (
		clientName = fmt.Sprintf("btc:%d", chainID)
		isRegnet   = pkgchains.IsBitcoinRegnet(chainID)
	)

	c := &Client{
		hostURL:    normalizeHostURL(cfg.RPCHost, true),
		client:     defaultHTTPClient(),
		config:     cfg,
		params:     params,
		clientName: clientName,
		isRegnet:   isRegnet,
		logger: logger.With().
			Str(logs.FieldModule, "btc_client").
			Int64(logs.FieldChain, chainID).
			Logger(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// NetParams returns the bitcoin network parameters
func (c *Client) NetParams() *chains.Params {
	return &c.params
}

// send sends RPC command to the server via http post request
func (c *Client) sendCommand(ctx context.Context, cmd any) (json.RawMessage, error) {
	method, reqBody, err := c.marshalCmd(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal cmd")
	}

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

func (c *Client) sendRequest(req *http.Request, method string) (out rawResponse, err error) {
	c.logger.Debug().Str("rpc.method", method).Msg("Sending request")
	start := time.Now()

	defer func() {
		c.recordMetrics(method, start, out, err)
		c.logger.Debug().Err(err).
			Str("rpc.method", method).
			Float64("rpc.duration", time.Since(start).Seconds()).
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

	if res.StatusCode != http.StatusOK {
		return rawResponse{}, errors.Errorf("unexpected status code %d (%s)", res.StatusCode, resBody)
	}

	if err = json.Unmarshal(resBody, &out); err != nil {
		return rawResponse{}, errors.Wrapf(err, "unable to unmarshal rpc response (%s)", resBody)
	}

	return out, nil
}

func (c *Client) recordMetrics(method string, start time.Time, out rawResponse, err error) {
	dur := time.Since(start).Seconds()

	status := "ok"
	if err != nil || out.Error != nil {
		status = "failed"
	}

	metrics.RPCClientCounter.WithLabelValues(status, c.clientName, method).Inc()
	metrics.RPCClientDuration.WithLabelValues(c.clientName).Observe(dur)
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

func unmarshalHex(raw json.RawMessage) ([]byte, error) {
	str, err := unmarshal[string](raw)
	if err != nil {
		return nil, err
	}

	return hex.DecodeString(str)
}

func resolveParams(name string) (chains.Params, error) {
	const regNetAlias = "regnet"

	switch name {
	case chains.MainNetParams.Name:
		return chains.MainNetParams, nil
	case chains.TestNet3Params.Name:
		return chains.TestNet3Params, nil
	case chaincfg.RegressionNetParams.Name, regNetAlias:
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
