package tss

import (
	"context"
	"fmt"
	"os"
	"path"
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
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/conversion"
	"gitlab.com/thorchain/tss/go-tss/tss"

	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// SetupProps represents options for Setup.
type SetupProps struct {
	Config          config.Config
	Zetacore        *zetacore.Client
	HotKeyPassword  string
	BitcoinChainIDs []int64
	PostBlame       bool
}

// Setup beefy function that does all the logic for bootstrapping tss-server, tss signer,
// generating TSS key is needed, etc...
func Setup(ctx context.Context, p SetupProps, logger zerolog.Logger) (*Service, error) {
	logger = logger.With().Str(logs.FieldModule, "tss_setup").Logger()

	// 0. Resolve Hot Private Key
	hotPrivateKey, err := p.Zetacore.GetKeys().GetPrivateKey(p.HotKeyPassword)
	switch {
	case err != nil:
		return nil, errors.Wrap(err, "unable to get hot private key")
	case len(hotPrivateKey.Bytes()) != 32:
		return nil, fmt.Errorf("hot privateKey: expect 32 bytes, got %d bytes", len(hotPrivateKey.Bytes()))
	}

	p.Zetacore.GetKeys().GetKeybase()

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
		logger.Warn().Msg("No bootstrap peers provided")
	} else {
		logger.Info().Interface("bootstrap_peers", bootstrapPeers).Msgf("Bootstrap peers")
	}

	// 2. Resolve pre-params. We want to enforce pre-params file existence
	tssPreParams, err := ResolvePreParamsFromPath(p.Config.PreParamsPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve TSS pre-params. Use `zetaclient tss gen-pre-params`")
	}

	logger.Info().Msg("Pre-params file resolved")

	// 3. Prepare whitelist of peers
	tssKeygen, err := p.Zetacore.GetKeyGen(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get TSS keygen")
	}

	logger.Info().Msg("Fetched TSS keygen info")

	whitelistedPeers := make([]peer.ID, len(tssKeygen.GranteePubkeys))
	for i, pk := range tssKeygen.GranteePubkeys {
		whitelistedPeers[i], err = conversion.Bech32PubkeyToPeerID(pk)
		if err != nil {
			return nil, errors.Wrap(err, pk)
		}
	}

	logger.Info().Interface("whitelisted_peers", whitelistedPeers).Msg("Resolved whitelist peers")

	// 4.
	// 	err = newTss.LoadTssFilesFromDirectory(app.Config().TssPath)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	_, pubkeyInBech32, err := keys.GetKeyringKeybase(app.Config(), hotkeyPassword)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	err = newTss.VerifyKeysharesForPubkeys(tssHistoricalList, pubkeyInBech32)
	//	if err != nil {
	//		client.GetLogger().Error().Err(err).Msg("VerifyKeysharesForPubkeys fail")
	//	}

	// todo bump numbers
	// 4. Bootstrap go-tss TSS server
	tssServer, err := NewTSSServer(
		bootstrapPeers,
		whitelistedPeers,
		tssPreParams,
		hotPrivateKeyECDSA,
		p.Config,
		p.HotKeyPassword,
		logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to start TSS server")
	}

	logger.Info().Msg("TSS server started")

	// 5. Perform key generation (if needed)
	if err = KeygenCeremony(ctx, tssServer, p.Zetacore, logger); err != nil {
		return nil, errors.Wrap(err, "unable to perform keygen ceremony")
	}

	// 6. Get tss & tss history from zetacore
	tssInfo, err := p.Zetacore.GetTSS(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get TSS from zetacore")
	}

	logger.Info().Msg("Got TSS info from zetacore")

	historicalTSSInfo, err := p.Zetacore.GetTSSHistory(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get TSS history")
	}

	// 7. Verify key shared for public keys
	logger.Info().Msg("Got historical TSS info from zetacore. Verifying key shares...")
	if err = verifyKeySharesForPubKeys(p, historicalTSSInfo, logger); err != nil {
		return nil, errors.Wrap(err, "unable to verify key shares for pub keys")
	}

	logger.Info().Msg("Key shared verified")

	// 8. Optionally test key signing
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
		WithMetrics(ctx, p.Zetacore, &Metrics{
			ActiveMsgsSigns:    metrics.NumActiveMsgSigns,
			SignLatency:        metrics.SignLatency,
			NodeBlamePerPubKey: metrics.TssNodeBlamePerPubKey,
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create TSS service")
	}

	logger.Info().Msg("TSS service created")

	if err = validateAddresses(service, p.BitcoinChainIDs, logger); err != nil {
		return nil, errors.Wrap(err, "unable to validate tss addresses")
	}

	logger.Info().Msg("TSS addresses validated")

	// todo health checks

	return service, nil
}

// NewTSSServer creates a new tss.TssServer (go-tss) instance for key signing.
// - bootstrapPeers are used to discover other peers
// - whitelistPeers are the only peers that are allowed in p2p key signing.
// - preParams are the TSS pre-params required for key generation
func NewTSSServer(
	bootstrapPeers []multiaddr.Multiaddr,
	whitelistPeers []peer.ID,
	preParams *keygen.LocalPreParams,
	privateKey crypto.PrivKey,
	cfg config.Config,
	tssPassword string,
	logger zerolog.Logger,
) (*tss.TssServer, error) {
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

	tssConfig := tsscommon.TssConfig{
		EnableMonitor:   true,              // enables prometheus metrics
		KeyGenTimeout:   300 * time.Second, // must be shorter than constants.JailTimeKeygen
		KeySignTimeout:  30 * time.Second,  // must be shorter than constants.JailTimeKeysign
		PartyTimeout:    30 * time.Second,
		PreParamTimeout: 5 * time.Minute,
	}

	tssServer, err := tss.NewTss(
		bootstrapPeers,
		Port,
		privateKey,
		tssPath,
		tssConfig,
		preParams,
		cfg.PublicIP,
		tssPassword,
		whitelistPeers,
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

	logger.Info().Msgf("TSS local peer ID is %s", tssServer.GetLocalPeerID())

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

	tssPath = path.Join(home, ".tss")
	logger.Warn().Msgf("TSS path is empty, falling back to %s", tssPath)

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

	pubKeys, err := ParsePubKeysFromPath(tssPath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to parse public keys from path")
	}

	pubKeysSet := make(map[string]PubKey, len(pubKeys))
	for _, k := range pubKeys {
		pubKeysSet[k.Bech32String()] = k
	}

	// Get observer's public key ("grantee pub key")
	_, granteePubKeyBech32, err := keys.GetKeyringKeybase(p.Config, p.HotKeyPassword)
	if err != nil {
		return errors.Wrap(err, "unable to get keyring key base")
	}

	wasPartOfTSS := func(grantees []string) bool {
		return slices.Contains(grantees, granteePubKeyBech32)
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

	logger.Info().Str("evm", evm.String()).Msg("EVM address")

	// validate TSS BTC address for each btc chain
	for _, chainID := range btcChainIDs {
		addr, err := service.PubKey().AddressBTC(chainID)
		if err != nil {
			return fmt.Errorf("unable to derive BTC address for chain %d", chainID)
		}

		logger.Info().Int64("chain_id", chainID).Str("addr", addr.String()).Msg("BTC address")
	}

	return nil
}
