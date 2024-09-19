package metrics

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/zeta-chain/node/pkg/chains"
	. "gopkg.in/check.v1"
)

type MetricsSuite struct {
	m *Metrics
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&MetricsSuite{})

func (ms *MetricsSuite) SetUpSuite(c *C) {
	m, err := NewMetrics()
	c.Assert(err, IsNil)
	m.Start()
	ms.m = m
}

// assert that the curried metric actually uses the same underlying storage
func (ms *MetricsSuite) TestCurryWith(c *C) {
	rpcTotalsC := RPCCount.MustCurryWith(prometheus.Labels{"host": "test"})
	rpcTotalsC.With(prometheus.Labels{"code": "400"}).Add(1.0)

	rpcCtr := testutil.ToFloat64(RPCCount.With(prometheus.Labels{"host": "test", "code": "400"}))
	c.Assert(rpcCtr, Equals, 1.0)

	RPCCount.Reset()
}

func (ms *MetricsSuite) Test_RPCCount(c *C) {
	GetFilterLogsPerChain.WithLabelValues("chain1").Inc()
	GetFilterLogsPerChain.WithLabelValues("chain2").Inc()
	GetFilterLogsPerChain.WithLabelValues("chain2").Inc()
	time.Sleep(1 * time.Second)

	chain1Ctr := testutil.ToFloat64(GetFilterLogsPerChain.WithLabelValues("chain1"))
	c.Assert(chain1Ctr, Equals, 1.0)

	httpClient, err := GetInstrumentedHTTPClient("http://127.0.0.1:8886/myauthuuid")
	c.Assert(err, IsNil)

	res, err := httpClient.Get("http://127.0.0.1:8886")
	c.Assert(err, IsNil)
	defer res.Body.Close()
	c.Assert(res.StatusCode, Equals, http.StatusOK)

	res, err = httpClient.Get("http://127.0.0.1:8886/metrics")
	c.Assert(err, IsNil)
	defer res.Body.Close()
	c.Assert(res.StatusCode, Equals, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	c.Assert(err, IsNil)
	metricsBody := string(body)
	c.Assert(strings.Contains(metricsBody, fmt.Sprintf("%s_%s", ZetaClientNamespace, "rpc_count")), Equals, true)

	// assert that rpc count is being incremented at all
	rpcCount := testutil.ToFloat64(RPCCount)
	c.Assert(rpcCount, Equals, 2.0)

	// assert that rpc count is being incremented correctly
	rpcCount = testutil.ToFloat64(RPCCount.With(prometheus.Labels{"host": "127.0.0.1:8886", "code": "200"}))
	c.Assert(rpcCount, Equals, 2.0)

	// assert that rpc count is not being incremented incorrectly
	rpcCount = testutil.ToFloat64(RPCCount.With(prometheus.Labels{"host": "127.0.0.1:8886", "code": "502"}))
	c.Assert(rpcCount, Equals, 0.0)
}

func (ms *MetricsSuite) Test_RelayerKeyBalance(c *C) {
	RelayerKeyBalance.WithLabelValues(chains.SolanaDevnet.Name).Set(2.1564)

	// assert that relayer key balance is being set correctly
	balance := testutil.ToFloat64(RelayerKeyBalance.WithLabelValues(chains.SolanaDevnet.Name))
	c.Assert(balance, Equals, 2.1564)
}
