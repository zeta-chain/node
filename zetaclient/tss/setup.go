package tss

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tsscommon "github.com/zeta-chain/go-tss/common"
	"github.com/zeta-chain/go-tss/conversion"
	"github.com/zeta-chain/go-tss/tss"

	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// SetupProps represents options for Setup.
type SetupProps struct {
	Config              config.Config
	Zetacore            Zetacore
	GranteePubKeyBech32 string
	HotKeyPassword      string
	TSSKeyPassword      string
	BitcoinChainIDs     []int64
	PostBlame           bool
	Telemetry           Telemetry
}

// Setup beefy function that does all the logic for bootstrapping tss-server, tss signer,
// generating TSS key is needed, etc...
func Setup(ctx context.Context, p SetupProps, logger zerolog.Logger) (*Service, error) {
	setupLogger := logger.With().Str(logs.FieldModule, logs.ModNameTssSetup).Logger()

	// 0. Resolve Hot Private Key
	hotPrivateKey, err := p.Zetacore.GetKeys().GetPrivateKey(p.HotKeyPassword)
	switch {
	case err != nil:
		return nil, errors.Wrap(err, "unable to get hot private key")
	case len(hotPrivateKey.Bytes()) != 32:
		return nil, fmt.Errorf("hot privateKey: expect 32 bytes, got %d bytes", len(hotPrivateKey.Bytes()))
	}

	hotPrivateKeyECDSA := secp256k1.PrivKey(hotPrivateKey.Bytes()[:32])

	// 1. Parse bootstrap peer if provided
	var bootstrapPeers []multiaddr.Multiaddr
	if p.Config.Peer != "" {
		bp, err := MultiAddressFromString(p.Config.Peer)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse bootstrap peers (%s)", p.Config.Peer)
		}
		bootstrapPeers = bp
	}

	if len(bootstrapPeers) == 0 {
		setupLogger.Warn().Msg("no bootstrap peers provided")
	} else {
		setupLogger.Info().Interface("bootstrap_peers", bootstrapPeers).Send()
	}

	// 2. Resolve pre-params. We want to enforce pre-params file existence
	tssPreParams, err := ResolvePreParamsFromPath(p.Config.PreParamsPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve TSS pre-params; use `zetaclient tss gen-pre-params`")
	}

	setupLogger.Info().Msg("resolved pre-params file")

	// 3. Prepare whitelist of peers
	tssKeygen, err := p.Zetacore.GetKeyGen(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get TSS keygen")
	}

	setupLogger.Info().Msg("fetched TSS keygen info")

	whitelistedPeers := make([]peer.ID, len(tssKeygen.GranteePubkeys))
	for i, pk := range tssKeygen.GranteePubkeys {
		whitelistedPeers[i], err = conversion.Bech32PubkeyToPeerID(pk)
		if err != nil {
			return nil, errors.Wrap(err, pk)
		}
	}

	setupLogger.Info().
		Any("whitelisted_peers", whitelistedPeers).
		Msg("resolved whitelist peers")

	// 4. Bootstrap go-tss TSS server
	tssServer, err := NewServer(
		bootstrapPeers,
		whitelistedPeers,
		tssPreParams,
		hotPrivateKeyECDSA,
		p.Config,
		p.TSSKeyPassword,
		logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to start TSS server")
	}

	if p.Telemetry != nil {
		p.Telemetry.SetP2PID(tssServer.GetLocalPeerID())
	}

	setupLogger.Info().Msg("started TSS server")

	// 5. Perform key generation (if needed)
	tssInfo, err := KeygenCeremony(ctx, tssServer, p.Zetacore, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to perform keygen ceremony")
	}

	historicalTSSInfo, err := p.Zetacore.GetTSSHistory(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get TSS history")
	}

	// 6. Verify key shares
	setupLogger.Info().Msg("got historical TSS info from zetacore; verifying key shares...")
	if err = verifyKeySharesForPubKeys(p, historicalTSSInfo, logger); err != nil {
		return nil, errors.Wrap(err, "unable to verify key shares for pub keys")
	}

	setupLogger.Info().Msg("key shares verified")

	// 7. Optionally test key signing
	if p.Config.TestTssKeysign {
		if err = TestKeySign(tssServer, tssInfo.TssPubkey, logger); err != nil {
			return nil, errors.Wrap(err, "unable to test key signing")
		}
	}

	// 8. Setup TSS zetaclient service (wrapper around go-tss TssServer)
	service, err := NewService(
		tssServer,
		tssInfo.TssPubkey,
		p.Zetacore,
		logger,
		WithPostBlame(p.PostBlame),
		WithRateLimit(p.Config.TSSMaxPendingSignatures),
		WithMetrics(ctx, p.Zetacore, &Metrics{
			ActiveMsgsSigns:    metrics.NumActiveMsgSigns,
			SignLatency:        metrics.SignLatency,
			NodeBlamePerPubKey: metrics.TSSNodeBlamePerPubKey,
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create TSS service")
	}

	setupLogger.Info().Msg("TSS service created")

	// 9. Ensure that TSS has valid EVM and BTC addresses
	if err = validateAddresses(service, p.BitcoinChainIDs, logger); err != nil {
		return nil, errors.Wrap(err, "unable to validate tss addresses")
	}

	setupLogger.Info().Msg("validated TSS addresses; starting healthcheck worker")

	healthCheckProps := HealthcheckProps{
		Telemetry:               p.Telemetry,
		WhitelistPeers:          whitelistedPeers,
		NumConnectedPeersMetric: metrics.NumConnectedPeers,
	}

	// 10. Start healthcheck worker
	if err = HealthcheckWorker(ctx, tssServer, healthCheckProps, logger); err != nil {
		return nil, errors.Wrap(err, "unable to start healthcheck worker")
	}

	return service, nil
}

// NewServer creates a new tss.TssServer (go-tss) instance for key signing.
// - bootstrapPeers are used to discover other peers
// - whitelistPeers are the only peers that are allowed in p2p key signing.
// - preParams are the TSS pre-params required for key generation
func NewServer(
	bootstrapPeers []multiaddr.Multiaddr,
	whitelistPeers []peer.ID,
	preParams *keygen.LocalPreParams,
	privateKey crypto.PrivKey,
	cfg config.Config,
	tssPassword string,
	logger zerolog.Logger,
) (*tss.Server, error) {
	switch {
	case len(whitelistPeers) == 0 && len(bootstrapPeers) == 0:
		return nil, errors.New("no bootstrap peers and whitelist peers provided")
	case preParams == nil:
		return nil, errors.New("pre-params are nil")
	case tssPassword == "":
		return nil, errors.New("tss password is empty")
	case privateKey == nil:
		return nil, errors.New("private key is nil")
	case cfg.PublicIP == "":
		logger.Warn().Msg("public IP is empty")
	}

	tssPath, err := resolveTSSPath(cfg.TssPath, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve TSS path")
	}

	networkConfig := tss.NetworkConfig{
		TssConfig: tsscommon.TssConfig{
			EnableMonitor:   true,              // enables prometheus metrics
			KeyGenTimeout:   300 * time.Second, // must be shorter than constants.JailTimeKeygen
			KeySignTimeout:  30 * time.Second,  // must be shorter than constants.JailTimeKeysign
			PartyTimeout:    30 * time.Second,
			PreParamTimeout: 5 * time.Minute,
		},
		ExternalIP:       cfg.PublicIP,
		Port:             Port,
		BootstrapPeers:   bootstrapPeers,
		WhitelistedPeers: whitelistPeers,
	}

	tssServer, err := tss.New(
		networkConfig,
		tssPath,
		privateKey,
		tssPassword,
		preParams,
		logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create TSS server")
	}

	// fyi: actually it does nothing, just logs "starting the tss servers"
	if err = tssServer.Start(); err != nil {
		return nil, errors.Wrap(err, "unable to start TSS server")
	}

	if isEmptyPeerID(tssServer.GetLocalPeerID()) {
		return nil, fmt.Errorf("local peer ID is empty, aborting")
	}

	logger.Info().Str("id", tssServer.GetLocalPeerID()).Msg("TSS local peer ID")

	return tssServer, nil
}

func resolveTSSPath(tssPath string, logger zerolog.Logger) (string, error) {
	// noop
	if tssPath != "" {
		return tssPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "unable to get user home dir")
	}

	tssPath = filepath.Join(home, ".tss")
	logger.Warn().
		Str("to", tssPath).
		Msg("TSS path is empty, falling back")

	return tssPath, nil
}

// not sure regarding this function, but the logic is the same
// as in the original code (for backward compatibility)
func isEmptyPeerID(id string) bool {
	return id == "" || id == "0" || id == "000000000000000000000000000000" || id == peer.ID("").String()
}

// verifyKeySharesForPubKeys ensures that observer&signer has the correct key shares
func verifyKeySharesForPubKeys(p SetupProps, history []observertypes.TSS, logger zerolog.Logger) error {
	// Parse bech32 public keys from tssPath (i.e. zetapub*...)
	tssPath, err := resolveTSSPath(p.Config.TssPath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to resolve TSS path")
	}

	pubKeys, err := ParsePubKeysFromPath(logger, tssPath)
	if err != nil {
		return errors.Wrap(err, "unable to parse public keys from path")
	}

	pubKeysSet := make(map[string]PubKey, len(pubKeys))
	for _, k := range pubKeys {
		pubKeysSet[k.Bech32String()] = k
	}

	wasPartOfTSS := func(grantees []string) bool {
		return slices.Contains(grantees, p.GranteePubKeyBech32)
	}

	for _, tssEntry := range history {
		if !wasPartOfTSS(tssEntry.TssParticipantList) {
			continue
		}

		if _, ok := pubKeysSet[tssEntry.TssPubkey]; !ok {
			return fmt.Errorf("pubkey %q not found in keyshare", tssEntry.TssPubkey)
		}
	}

	return nil
}

// validateAddresses ensures that TSS has valid EVM and BTC addresses.
func validateAddresses(service *Service, btcChainIDs []int64, logger zerolog.Logger) error {
	evm := service.PubKey().AddressEVM()
	if evm == (ethcommon.Address{}) {
		return fmt.Errorf("blank tss evm address is empty")
	}

	logger.Info().
		Stringer("address", evm).
		Msg("EVM address")

	// validate TSS BTC address for each btc chain
	for _, chainID := range btcChainIDs {
		addr, err := service.PubKey().AddressBTC(chainID)
		if err != nil {
			return fmt.Errorf("unable to derive BTC address for chain %d", chainID)
		}

		logger.Info().
			Int64(logs.FieldChain, chainID).
			Stringer("address", addr).
			Msg("BTC address")
	}

	return nil
}
