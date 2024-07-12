// Package tss provides the TSS signer functionalities for the zetaclient to sign transactions on external chains
package tss

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	tmcrypto "github.com/cometbft/cometbft/crypto"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gopeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	thorcommon "github.com/zeta-chain/go-tss/common"
	"github.com/zeta-chain/go-tss/keysign"
	"github.com/zeta-chain/go-tss/p2p"
	"github.com/zeta-chain/go-tss/tss"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/cosmos"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

const (
	// envFlagPostBlame is the environment flag to enable posting blame data to core
	envFlagPostBlame = "POST_BLAME"
)

// Key is a struct that holds the public key, bech32 pubkey, and address for the TSS
type Key struct {
	PubkeyInBytes  []byte
	PubkeyInBech32 string
	AddressInHex   string
}

// NewTSSKey creates a new TSS key
func NewTSSKey(pk string) (*Key, error) {
	TSSKey := &Key{
		PubkeyInBech32: pk,
	}
	pubkey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, pk)
	if err != nil {
		log.Error().Err(err).Msgf("GetPubKeyFromBech32 from %s", pk)
		return nil, fmt.Errorf("GetPubKeyFromBech32: %w", err)
	}

	decompresspubkey, err := crypto.DecompressPubkey(pubkey.Bytes())
	if err != nil {
		return nil, fmt.Errorf("NewTSS: DecompressPubkey error: %w", err)
	}

	TSSKey.PubkeyInBytes = crypto.FromECDSAPub(decompresspubkey)
	TSSKey.AddressInHex = crypto.PubkeyToAddress(*decompresspubkey).Hex()

	return TSSKey, nil
}

var _ interfaces.TSSSigner = (*TSS)(nil)

// TSS is a struct that holds the server and the keys for TSS
type TSS struct {
	Server          *tss.TssServer
	Keys            map[string]*Key // PubkeyInBech32 => TSSKey
	CurrentPubkey   string
	logger          zerolog.Logger
	Signers         []string
	ZetacoreClient  interfaces.ZetacoreClient
	KeysignsTracker *ConcurrentKeysignsTracker

	// TODO: support multiple Bitcoin network, not just one network
	// https://github.com/zeta-chain/node/issues/1397
	BitcoinChainID int64
}

// NewTSS creates a new TSS instance
func NewTSS(
	ctx context.Context,
	appContext *appcontext.AppContext,
	peer p2p.AddrList,
	privkey tmcrypto.PrivKey,
	preParams *keygen.LocalPreParams,
	client interfaces.ZetacoreClient,
	tssHistoricalList []observertypes.TSS,
	bitcoinChainID int64,
	tssPassword string,
	hotkeyPassword string,
) (*TSS, error) {
	logger := log.With().Str("module", "tss_signer").Logger()
	server, err := SetupTSSServer(peer, privkey, preParams, appContext.Config(), tssPassword)
	if err != nil {
		return nil, fmt.Errorf("SetupTSSServer error: %w", err)
	}

	newTss := TSS{
		Server:          server,
		Keys:            make(map[string]*Key),
		CurrentPubkey:   appContext.GetCurrentTssPubKey(),
		logger:          logger,
		ZetacoreClient:  client,
		KeysignsTracker: NewKeysignsTracker(logger),
		BitcoinChainID:  bitcoinChainID,
	}

	err = newTss.LoadTssFilesFromDirectory(appContext.Config().TssPath)
	if err != nil {
		return nil, err
	}

	_, pubkeyInBech32, err := keys.GetKeyringKeybase(appContext.Config(), hotkeyPassword)
	if err != nil {
		return nil, err
	}

	err = newTss.VerifyKeysharesForPubkeys(tssHistoricalList, pubkeyInBech32)
	if err != nil {
		client.GetLogger().Error().Err(err).Msg("VerifyKeysharesForPubkeys fail")
	}

	keygenRes, err := newTss.ZetacoreClient.GetKeyGen(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize metrics
	for _, key := range keygenRes.GranteePubkeys {
		metrics.TssNodeBlamePerPubKey.WithLabelValues(key).Inc()
	}
	metrics.NumActiveMsgSigns.Set(0)

	return &newTss, nil
}

// SetupTSSServer creates a new TSS server
// TODO(revamp): move to TSS server file
func SetupTSSServer(
	peer p2p.AddrList,
	privkey tmcrypto.PrivKey,
	preParams *keygen.LocalPreParams,
	cfg config.Config,
	tssPassword string,
) (*tss.TssServer, error) {
	bootstrapPeers := peer
	log.Info().Msgf("Peers AddrList %v", bootstrapPeers)

	tsspath := cfg.TssPath
	if len(tsspath) == 0 {
		log.Error().Msg("empty env TSSPATH")
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Error().Err(err).Msgf("cannot get UserHomeDir")
			return nil, err
		}
		tsspath = path.Join(homedir, ".Tss")
		log.Info().Msgf("create temporary TSSPATH: %s", tsspath)
	}

	IP := cfg.PublicIP
	if len(IP) == 0 {
		log.Info().Msg("empty public IP in config")
	}

	tssServer, err := tss.NewTss(
		bootstrapPeers,
		6668,
		privkey,
		"MetaMetaOpenTheDoor",
		tsspath,
		thorcommon.TssConfig{
			EnableMonitor:   true,
			KeyGenTimeout:   300 * time.Second, // must be shorter than constants.JailTimeKeygen
			KeySignTimeout:  30 * time.Second,  // must be shorter than constants.JailTimeKeysign
			PartyTimeout:    30 * time.Second,
			PreParamTimeout: 5 * time.Minute,
		},
		preParams, // use pre-generated pre-params if non-nil
		IP,        // for docker test
		tssPassword,
	)
	if err != nil {
		log.Error().Err(err).Msg("NewTSS error")
		return nil, fmt.Errorf("NewTSS error: %w", err)
	}

	err = tssServer.Start()
	if err != nil {
		log.Error().Err(err).Msg("tss server start error")
	}

	log.Info().Msgf("LocalID: %v", tssServer.GetLocalPeerID())
	if tssServer.GetLocalPeerID() == "" ||
		tssServer.GetLocalPeerID() == "0" ||
		tssServer.GetLocalPeerID() == "000000000000000000000000000000" ||
		tssServer.GetLocalPeerID() == gopeer.ID("").String() {
		log.Error().Msg("tss server start error")
		return nil, fmt.Errorf("tss server start error")
	}

	return tssServer, nil
}

// Pubkey returns the current pubkey
func (tss *TSS) Pubkey() []byte {
	return tss.Keys[tss.CurrentPubkey].PubkeyInBytes
}

// Sign signs a digest
// digest should be Hashes of some data
// NOTE: Specify optionalPubkey to use a different pubkey than the current pubkey set during keygen
func (tss *TSS) Sign(
	ctx context.Context,
	digest []byte,
	height uint64,
	nonce uint64,
	chainID int64,
	optionalPubKey string,
) ([65]byte, error) {
	H := digest
	log.Debug().Msgf("hash of digest is %s", H)

	tssPubkey := tss.CurrentPubkey
	if optionalPubKey != "" {
		tssPubkey = optionalPubKey
	}

	// #nosec G115 always in range
	keysignReq := keysign.NewRequest(
		tssPubkey,
		[]string{base64.StdEncoding.EncodeToString(H)},
		int64(height),
		nil,
		"0.14.0",
	)
	tss.KeysignsTracker.StartMsgSign()
	ksRes, err := tss.Server.KeySign(keysignReq)
	tss.KeysignsTracker.EndMsgSign()
	if err != nil {
		log.Warn().Msg("keysign fail")
	}

	if ksRes.Status == thorcommon.Fail {
		log.Warn().Msgf("keysign status FAIL posting blame to core, blaming node(s): %#v", ksRes.Blame.BlameNodes)

		// post blame data if enabled
		if IsEnvFlagEnabled(envFlagPostBlame) {
			digest := hex.EncodeToString(digest)
			index := observertypes.GetBlameIndex(chainID, nonce, digest, height)
			zetaHash, err := tss.ZetacoreClient.PostVoteBlameData(ctx, &ksRes.Blame, chainID, index)
			if err != nil {
				log.Error().Err(err).Msg("error sending blame data to core")
				return [65]byte{}, err
			}
			log.Info().Msgf("keysign posted blame data tx hash: %s", zetaHash)
		}

		// Increment Blame counter
		for _, node := range ksRes.Blame.BlameNodes {
			metrics.TssNodeBlamePerPubKey.WithLabelValues(node.Pubkey).Inc()
		}
	}
	signature := ksRes.Signatures

	// [{cyP8i/UuCVfQKDsLr1kpg09/CeIHje1FU6GhfmyMD5Q= D4jXTH3/CSgCg+9kLjhhfnNo3ggy9DTQSlloe3bbKAs= eY++Z2LwsuKG1JcghChrsEJ4u9grLloaaFZNtXI3Ujk= AA==}]
	// 32B msg hash, 32B R, 32B S, 1B RC
	log.Info().Msgf("signature of digest is... %v", signature)

	if len(signature) == 0 {
		log.Warn().Err(err).Msgf("signature has length 0")
		return [65]byte{}, fmt.Errorf("keysign fail: %s", err)
	}

	if !verifySignature(tssPubkey, signature, H) {
		log.Error().Err(err).Msgf("signature verification failure")
		return [65]byte{}, fmt.Errorf("signuature verification fail")
	}

	var sigbyte [65]byte
	_, err = base64.StdEncoding.Decode(sigbyte[:32], []byte(signature[0].R))
	if err != nil {
		log.Error().Err(err).Msg("decoding signature R")
		return [65]byte{}, fmt.Errorf("signuature verification fail")
	}

	_, err = base64.StdEncoding.Decode(sigbyte[32:64], []byte(signature[0].S))
	if err != nil {
		log.Error().Err(err).Msg("decoding signature S")
		return [65]byte{}, fmt.Errorf("signuature verification fail")
	}

	_, err = base64.StdEncoding.Decode(sigbyte[64:65], []byte(signature[0].RecoveryID))
	if err != nil {
		log.Error().Err(err).Msg("decoding signature RecoveryID")
		return [65]byte{}, fmt.Errorf("signuature verification fail")
	}

	return sigbyte, nil
}

// SignBatch is hash of some data
// digest should be batch of hashes of some data
func (tss *TSS) SignBatch(
	ctx context.Context,
	digests [][]byte,
	height uint64,
	nonce uint64,
	chainID int64,
) ([][65]byte, error) {
	tssPubkey := tss.CurrentPubkey
	digestBase64 := make([]string, len(digests))
	for i, digest := range digests {
		digestBase64[i] = base64.StdEncoding.EncodeToString(digest)
	}
	// #nosec G115 always in range
	keysignReq := keysign.NewRequest(tssPubkey, digestBase64, int64(height), nil, "0.14.0")

	tss.KeysignsTracker.StartMsgSign()
	ksRes, err := tss.Server.KeySign(keysignReq)
	tss.KeysignsTracker.EndMsgSign()
	if err != nil {
		log.Warn().Err(err).Msg("keysign fail")
	}

	if ksRes.Status == thorcommon.Fail {
		log.Warn().Msg("keysign status FAIL posting blame to core")

		// post blame data if enabled
		if IsEnvFlagEnabled(envFlagPostBlame) {
			digest := combineDigests(digestBase64)
			index := observertypes.GetBlameIndex(chainID, nonce, hex.EncodeToString(digest), height)
			zetaHash, err := tss.ZetacoreClient.PostVoteBlameData(ctx, &ksRes.Blame, chainID, index)
			if err != nil {
				log.Error().Err(err).Msg("error sending blame data to core")
				return [][65]byte{}, err
			}
			log.Info().Msgf("keysign posted blame data tx hash: %s", zetaHash)
		}

		// Increment Blame counter
		for _, node := range ksRes.Blame.BlameNodes {
			metrics.TssNodeBlamePerPubKey.WithLabelValues(node.Pubkey).Inc()
		}
	}

	signatures := ksRes.Signatures
	// [{cyP8i/UuCVfQKDsLr1kpg09/CeIHje1FU6GhfmyMD5Q= D4jXTH3/CSgCg+9kLjhhfnNo3ggy9DTQSlloe3bbKAs= eY++Z2LwsuKG1JcghChrsEJ4u9grLloaaFZNtXI3Ujk= AA==}]
	// 32B msg hash, 32B R, 32B S, 1B RC

	if len(signatures) != len(digests) {
		log.Warn().
			Err(err).
			Msgf("signature has length (%d) not equal to length of digests (%d)", len(signatures), len(digests))
		return [][65]byte{}, fmt.Errorf("keysign fail: %s", err)
	}

	//if !verifySignatures(tssPubkey, signatures, digests) {
	//	log.Error().Err(err).Msgf("signature verification failure")
	//	return [][65]byte{}, fmt.Errorf("signuature verification fail")
	//}
	pubkey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		log.Error().Msg("get pubkey from bech32 fail")
	}
	sigBytes := make([][65]byte, len(digests))
	for j, H := range digests {
		found := false
		D := base64.StdEncoding.EncodeToString(H)
		for _, signature := range signatures {
			if D == signature.Msg {
				found = true
				_, err = base64.StdEncoding.Decode(sigBytes[j][:32], []byte(signature.R))
				if err != nil {
					log.Error().Err(err).Msg("decoding signature R")
					return [][65]byte{}, fmt.Errorf("signuature verification fail")
				}
				_, err = base64.StdEncoding.Decode(sigBytes[j][32:64], []byte(signature.S))
				if err != nil {
					log.Error().Err(err).Msg("decoding signature S")
					return [][65]byte{}, fmt.Errorf("signuature verification fail")
				}
				_, err = base64.StdEncoding.Decode(sigBytes[j][64:65], []byte(signature.RecoveryID))
				if err != nil {
					log.Error().Err(err).Msg("decoding signature RecoveryID")
					return [][65]byte{}, fmt.Errorf("signuature verification fail")
				}
				sigPublicKey, err := crypto.SigToPub(H, sigBytes[j][:])
				if err != nil {
					log.Error().Err(err).Msg("SigToPub error in verify_signature")
					return [][65]byte{}, fmt.Errorf("signuature verification fail")
				}
				compressedPubkey := crypto.CompressPubkey(sigPublicKey)
				if !bytes.Equal(pubkey.Bytes(), compressedPubkey) {
					log.Warn().
						Msgf("%d-th pubkey %s recovered pubkey %s", j, pubkey.String(), hex.EncodeToString(compressedPubkey))
					return [][65]byte{}, fmt.Errorf("signuature verification fail")
				}
			}
		}
		if !found {
			log.Error().Err(err).Msg("signature not found")
			return [][65]byte{}, fmt.Errorf("signuature verification fail")
		}
	}

	return sigBytes, nil
}

// Validate validates the TSS
func (tss *TSS) Validate() error {
	evmAddress := tss.EVMAddress()
	blankAddress := ethcommon.Address{}

	if evmAddress == blankAddress {
		return fmt.Errorf("invalid evm address : %s", evmAddress.String())
	}

	if tss.BTCAddressWitnessPubkeyHash() == nil {
		return fmt.Errorf("invalid btc pub key hash : %s", tss.BTCAddress())
	}

	return nil
}

// EVMAddress generates an EVM address from pubkey
func (tss *TSS) EVMAddress() ethcommon.Address {
	addr, err := GetTssAddrEVM(tss.CurrentPubkey)
	if err != nil {
		log.Error().Err(err).Msg("getKeyAddr error")
		return ethcommon.Address{}
	}
	return addr
}

// BTCAddress generates a bech32 p2wpkh address from pubkey
func (tss *TSS) BTCAddress() string {
	addr, err := GetTssAddrBTC(tss.CurrentPubkey, tss.BitcoinChainID)
	if err != nil {
		log.Error().Err(err).Msg("getKeyAddr error")
		return ""
	}
	return addr
}

// BTCAddressWitnessPubkeyHash generates a bech32 p2wpkh address from pubkey
func (tss *TSS) BTCAddressWitnessPubkeyHash() *btcutil.AddressWitnessPubKeyHash {
	addrWPKH, err := getKeyAddrBTCWitnessPubkeyHash(tss.CurrentPubkey, tss.BitcoinChainID)
	if err != nil {
		log.Error().Err(err).Msg("BTCAddressPubkeyHash error")
		return nil
	}
	return addrWPKH
}

// PubKeyCompressedBytes returns the compressed bytes of the current pubkey
func (tss *TSS) PubKeyCompressedBytes() []byte {
	pubk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tss.CurrentPubkey)
	if err != nil {
		log.Error().Err(err).Msg("PubKeyCompressedBytes error")
		return nil
	}
	return pubk.Bytes()
}

// InsertPubKey adds a new key to the TSS keys map
func (tss *TSS) InsertPubKey(pk string) error {
	TSSKey, err := NewTSSKey(pk)
	if err != nil {
		return err
	}
	tss.Keys[pk] = TSSKey
	return nil
}

// VerifyKeysharesForPubkeys verifies the keyshares present on the node. It checks whether the node has TSS key shares for the TSS ceremonies it was part of.
func (tss *TSS) VerifyKeysharesForPubkeys(tssList []observertypes.TSS, granteePubKey32 string) error {
	for _, t := range tssList {
		if wasNodePartOfTss(granteePubKey32, t.TssParticipantList) {
			if _, ok := tss.Keys[t.TssPubkey]; !ok {
				return fmt.Errorf("pubkey %s not found in keyshare", t.TssPubkey)
			}
		}
	}
	return nil
}

// LoadTssFilesFromDirectory loads the TSS files at the directory specified by the `tssPath`
func (tss *TSS) LoadTssFilesFromDirectory(tssPath string) error {
	files, err := os.ReadDir(tssPath)
	if err != nil {
		fmt.Println("ReadDir error :", err.Error())
		return err
	}
	found := false

	var sharefiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(filepath.Base(file.Name()), "localstate") {
			sharefiles = append(sharefiles, file)
		}
	}

	if len(sharefiles) > 0 {
		sort.SliceStable(sharefiles, func(i, j int) bool {
			fi, err := sharefiles[i].Info()
			if err != nil {
				return false
			}
			fj, err := sharefiles[j].Info()
			if err != nil {
				return false
			}
			return fi.ModTime().After(fj.ModTime())
		})
		tss.logger.Info().Msgf("found %d localstate files", len(sharefiles))
		for _, localStateFile := range sharefiles {
			filename := filepath.Base(localStateFile.Name())
			filearray := strings.Split(filename, "-")
			if len(filearray) == 2 {
				log.Info().Msgf("Found stored Pubkey in local state: %s", filearray[1])
				pk := strings.TrimSuffix(filearray[1], ".json")

				err = tss.InsertPubKey(pk)
				if err != nil {
					log.Error().Err(err).Msg("InsertPubKey  in NewTSS fail")
				}
				tss.logger.Info().Msgf("registering TSS pubkey %s (eth hex %s)", pk, tss.Keys[pk].AddressInHex)
				found = true
			}
		}
	}

	if !found {
		log.Info().Msg("TSS Keyshare file NOT found")
	}
	return nil
}

// GetTssAddrBTC generates a bech32 p2wpkh address from pubkey
func GetTssAddrBTC(tssPubkey string, bitcoinChainID int64) (string, error) {
	addrWPKH, err := getKeyAddrBTCWitnessPubkeyHash(tssPubkey, bitcoinChainID)
	if err != nil {
		log.Fatal().Err(err)
		return "", err
	}

	return addrWPKH.EncodeAddress(), nil
}

// GetTssAddrEVM generates an EVM address from pubkey
func GetTssAddrEVM(tssPubkey string) (ethcommon.Address, error) {
	var keyAddr ethcommon.Address
	pubk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		log.Fatal().Err(err)
		return keyAddr, err
	}
	//keyAddrBytes := pubk.EVMAddress().Bytes()
	pubk.Bytes()
	decompresspubkey, err := crypto.DecompressPubkey(pubk.Bytes())
	if err != nil {
		log.Fatal().Err(err).Msg("decompress err")
		return keyAddr, err
	}

	keyAddr = crypto.PubkeyToAddress(*decompresspubkey)

	return keyAddr, nil
}

// TestKeysign tests the keysign
// it is called when a new TSS is generated to ensure the network works as expected
// TODO(revamp): move to a test package
func TestKeysign(tssPubkey string, tssServer *tss.TssServer) error {
	log.Info().Msg("trying keysign...")
	data := []byte("hello meta")
	H := crypto.Keccak256Hash(data)
	log.Info().Msgf("hash of data (hello meta) is %s", H)

	keysignReq := keysign.NewRequest(
		tssPubkey,
		[]string{base64.StdEncoding.EncodeToString(H.Bytes())},
		10,
		nil,
		"0.14.0",
	)
	ksRes, err := tssServer.KeySign(keysignReq)
	if err != nil {
		log.Warn().Msg("keysign fail")
	}

	signature := ksRes.Signatures
	// [{cyP8i/UuCVfQKDsLr1kpg09/CeIHje1FU6GhfmyMD5Q= D4jXTH3/CSgCg+9kLjhhfnNo3ggy9DTQSlloe3bbKAs= eY++Z2LwsuKG1JcghChrsEJ4u9grLloaaFZNtXI3Ujk= AA==}]
	// 32B msg hash, 32B R, 32B S, 1B RC
	log.Info().Msgf("signature of helloworld... %v", signature)

	if len(signature) == 0 {
		log.Info().Msgf("signature has length 0, skipping verify")
		return fmt.Errorf("signature has length 0")
	}

	verifySignature(tssPubkey, signature, H.Bytes())
	if verifySignature(tssPubkey, signature, H.Bytes()) {
		return nil
	}

	return fmt.Errorf("verify signature fail")
}

// IsEnvFlagEnabled checks if the environment flag is enabled
func IsEnvFlagEnabled(flag string) bool {
	value := os.Getenv(flag)
	return value == "true" || value == "1"
}

// verifySignature verifies the signature
// TODO(revamp): move to a test package
func verifySignature(tssPubkey string, signature []keysign.Signature, H []byte) bool {
	if len(signature) == 0 {
		log.Warn().Msg("verify_signature: empty signature array")
		return false
	}
	pubkey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		log.Error().Msg("get pubkey from bech32 fail")
	}

	// verify the signature of msg.
	var sigbyte [65]byte
	_, err = base64.StdEncoding.Decode(sigbyte[:32], []byte(signature[0].R))
	if err != nil {
		log.Error().Err(err).Msg("decoding signature R")
		return false
	}

	_, err = base64.StdEncoding.Decode(sigbyte[32:64], []byte(signature[0].S))
	if err != nil {
		log.Error().Err(err).Msg("decoding signature S")
		return false
	}

	_, err = base64.StdEncoding.Decode(sigbyte[64:65], []byte(signature[0].RecoveryID))
	if err != nil {
		log.Error().Err(err).Msg("decoding signature RecoveryID")
		return false
	}

	sigPublicKey, err := crypto.SigToPub(H, sigbyte[:])
	if err != nil {
		log.Error().Err(err).Msg("SigToPub error in verify_signature")
		return false
	}

	compressedPubkey := crypto.CompressPubkey(sigPublicKey)
	log.Info().Msgf("pubkey %s recovered pubkey %s", pubkey.String(), hex.EncodeToString(compressedPubkey))
	return bytes.Equal(pubkey.Bytes(), compressedPubkey)
}

// combineDigests combines the digests
func combineDigests(digestList []string) []byte {
	digestConcat := strings.Join(digestList[:], "")
	digestBytes := chainhash.DoubleHashH([]byte(digestConcat))
	return digestBytes.CloneBytes()
}

// wasNodePartOfTss checks if the node was part of the TSS
// it checks whether a pubkey is part of the list used to generate the TSS , Every TSS generated on the network has its own list of associated public keys
func wasNodePartOfTss(granteePubKey32 string, granteeList []string) bool {
	for _, grantee := range granteeList {
		if granteePubKey32 == grantee {
			return true
		}
	}
	return false
}

// getKeyAddrBTCWitnessPubkeyHash generates a bech32 p2wpkh address from pubkey
func getKeyAddrBTCWitnessPubkeyHash(tssPubkey string, chainID int64) (*btcutil.AddressWitnessPubKeyHash, error) {
	pubk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		return nil, err
	}

	bitcoinNetParams, err := chains.BitcoinNetParamsFromChainID(chainID)
	if err != nil {
		return nil, err
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubk.Bytes()), bitcoinNetParams)
	if err != nil {
		return nil, err
	}
	return addr, nil
}
