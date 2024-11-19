package tss

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	thorcommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keysign"

	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// KeySigner signs messages using TSS (subset of go-tss)
type KeySigner interface {
	KeySign(req keysign.Request) (keysign.Response, error)
}

// Service TSS service
type Service struct {
	zetacore      interfaces.ZetacoreClient
	tss           KeySigner
	currentPubKey PubKey

	postBlame bool
	metrics   *Metrics

	logger zerolog.Logger
}

// Metrics Prometheus metrics for the TSS service.
type Metrics struct {
	ActiveMsgsSigns    prometheus.Gauge
	SignLatency        *prometheus.HistogramVec
	NodeBlamePerPubKey *prometheus.CounterVec
}

type serviceConfig struct {
	postBlame bool
	metrics   *Metrics
}

// Opt Service option.
type Opt func(cfg *serviceConfig, logger zerolog.Logger) error

// WithPostBlame configures the TSS service to post blame in case of failed key signatures.
func WithPostBlame(postBlame bool) Opt {
	return func(cfg *serviceConfig, _ zerolog.Logger) error {
		cfg.postBlame = postBlame
		return nil
	}
}

// WithMetrics registers Prometheus metrics for the TSS service.
// Otherwise, no metrics will be collected.
func WithMetrics(ctx context.Context, zetacore interfaces.ZetacoreClient, m *Metrics) Opt {
	return func(cfg *serviceConfig, _ zerolog.Logger) error {
		keygen, err := zetacore.GetKeyGen(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get keygen (WithMetrics)")
		}

		m.ActiveMsgsSigns.Set(0)
		m.SignLatency.Reset()
		m.NodeBlamePerPubKey.Reset()

		for _, granteeBech32 := range keygen.GranteePubkeys {
			m.NodeBlamePerPubKey.WithLabelValues(granteeBech32).Inc()
		}

		cfg.metrics = m

		return nil
	}
}

var noopMetrics = Metrics{
	ActiveMsgsSigns:    prometheus.NewGauge(prometheus.GaugeOpts{Name: "noop"}),
	SignLatency:        prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "noop"}, []string{"result"}),
	NodeBlamePerPubKey: prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop"}, []string{"pubkey"}),
}

// NewService Service constructor.
// TODO LRU cache
func NewService(
	keySigner KeySigner,
	tssPubKeyBech32 string,
	zc interfaces.ZetacoreClient,
	logger zerolog.Logger,
	opts ...Opt,
) (*Service, error) {
	logger = logger.With().Str(logs.FieldModule, "tss_service").Logger()

	cfg := serviceConfig{
		metrics:   &noopMetrics,
		postBlame: false,
	}

	for _, opt := range opts {
		if err := opt(&cfg, logger); err != nil {
			return nil, errors.Wrap(err, "failed to apply tss config option")
		}
	}

	// Represents the current TSS public key.
	// FWIW, based on this, we can derive EVM / BTC addresses.
	currentPubKey, err := NewPubKeyFromBech32(tssPubKeyBech32)
	if err != nil {
		return nil, errors.Wrap(err, "invalid tss pub key")
	}

	return &Service{
		tss:           keySigner,
		currentPubKey: currentPubKey,
		zetacore:      zc,
		postBlame:     cfg.postBlame,
		metrics:       cfg.metrics,
		logger:        logger,
	}, nil
}

// PubKey returns current TSS PubKey.
func (s *Service) PubKey() PubKey {
	return s.currentPubKey
}

// Sign signs msg digest (hash). Returns signature in the format of R (32B), S (32B), V (1B).
func (s *Service) Sign(ctx context.Context, digest []byte, height, nonce uint64, chainID int64) ([65]byte, error) {
	sigs, err := s.SignBatch(ctx, [][]byte{digest}, height, nonce, chainID)
	if err != nil {
		return [65]byte{}, err
	}

	return sigs[0], nil
}

// SignBatch signs msgs digests (hash). Returns list of signatures in the format of R (32B), S (32B), V (1B).
func (s *Service) SignBatch(
	ctx context.Context,
	digests [][]byte,
	height, nonce uint64,
	chainID int64,
) ([][65]byte, error) {
	if len(digests) == 0 {
		return nil, errors.New("empty digests list")
	}

	// todo check cache for digest & block height & chainID -> return signature (LRU cache)

	digestsBase64 := make([]string, len(digests))
	for i, digest := range digests {
		digestsBase64[i] = base64EncodeString(digest)
	}

	tssPubKeyBech32 := s.PubKey().Bech32String()

	// #nosec G115 always in range
	req := keysign.NewRequest(
		tssPubKeyBech32,
		digestsBase64,
		int64(height),
		nil,
		Version,
	)

	res, err := s.sign(req)
	switch {
	case err != nil:
		// unexpected error (not related to failed key sign)
		return nil, errors.Wrap(err, "unable to perform a key sign")
	case res.Status == thorcommon.Fail:
		return nil, s.blameFailure(ctx, req, res, digests, height, nonce, chainID)
	case res.Status != thorcommon.Success:
		return nil, fmt.Errorf("keysign fail: status %d", res.Status)
	case len(res.Signatures) == 0:
		return nil, fmt.Errorf("keysign fail: signature list is empty")
	case len(res.Signatures) != len(digests):
		return nil, fmt.Errorf("keysign fail: signature list length mismatch")
	}

	signatures := make([][65]byte, len(res.Signatures))
	for i, sigResponse := range res.Signatures {
		signatures[i], err = VerifySignature(sigResponse, tssPubKeyBech32, digests[i])
		if err != nil {
			return nil, fmt.Errorf("unable to verify signature: %w (#%d)", err, i)
		}
	}

	// todo sig save to LRU cache (chain-id + digest). We need LRU per EACH chain

	return signatures, nil
}

var (
	signLabelsSuccess = prometheus.Labels{"result": "success"}
	signLabelsError   = prometheus.Labels{"result": "error"}
)

// sign sends TSS key sign request to the underlying go-tss and registers metrics
func (s *Service) sign(req keysign.Request) (res keysign.Response, err error) {
	// metrics start
	messagesCount, start := float64(len(req.Messages)), time.Now()
	s.metrics.ActiveMsgsSigns.Add(messagesCount)

	// metrics finish
	defer func() {
		s.metrics.ActiveMsgsSigns.Sub(messagesCount)

		latency := time.Since(start).Seconds()
		if err == nil && res.Status == thorcommon.Success {
			s.metrics.SignLatency.With(signLabelsSuccess).Observe(latency)
		} else {
			s.metrics.SignLatency.With(signLabelsError).Observe(latency)
		}
	}()

	return s.tss.KeySign(req)
}

func (s *Service) blameFailure(
	ctx context.Context,
	req keysign.Request,
	res keysign.Response,
	digests [][]byte,
	height uint64,
	nonce uint64,
	chainID int64,
) error {
	errFailure := errors.Errorf("keysign failed: %s", res.Blame.FailReason)
	lf := keysignLogFields(req, height, nonce, chainID)

	s.logger.Error().Err(errFailure).
		Fields(lf).
		Interface("keysign.fail_blame", res.Blame).
		Msg("Keysign failed")

	// register blame metrics
	for _, node := range res.Blame.BlameNodes {
		s.metrics.NodeBlamePerPubKey.WithLabelValues(node.Pubkey).Inc()
	}

	if !s.postBlame {
		return errFailure
	}

	var digest []byte
	if len(req.Messages) > 1 {
		digest = combineDigests(req.Messages)
	} else {
		digest = digests[0]
	}

	digestHex := hex.EncodeToString(digest)
	index := observertypes.GetBlameIndex(chainID, nonce, digestHex, height)
	zetaHash, err := s.zetacore.PostVoteBlameData(ctx, &res.Blame, chainID, index)
	if err != nil {
		return errors.Wrap(err, "unable to post blame data for failed keysign")
	}

	s.logger.Info().
		Fields(lf).
		Str("keygen.blame_tx_hash", zetaHash).
		Msg("Posted blame data to zetacore")

	return errFailure
}

func keysignLogFields(req keysign.Request, height, nonce uint64, chainID int64) map[string]any {
	return map[string]any{
		"keysign.chain_id":     chainID,
		"keysign.block_height": height,
		"keysign.nonce":        nonce,
		"keysign.request":      req,
	}
}
