package metaclient

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	thorcommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-peerstore/addr"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/tss/go-tss/conversion"
	"gitlab.com/thorchain/tss/go-tss/keysign"
	"gitlab.com/thorchain/tss/go-tss/tss"
	"os"
	"time"
)

var testPubKeys = []string{
	"metapub1addwnpepqtdklw8tf3anjz7nn5fly3uvq2e67w2apn560s4smmrt9e3x52nt2y5225d",
	"metapub1addwnpepqtspqyy6gk22u37ztra4hq3hdakc0w0k60sfy849mlml2vrpfr0wv0zdn8e",
	"metapub1addwnpepq2ryyje5zr09lq7gqptjwnxqsy2vcdngvwd6z7yt5yjcnyj8c8cn5qahgje",
	"metapub1addwnpepqfjcw5l4ay5t00c32mmlky7qrppepxzdlkcwfs2fd5u73qrwna0vzfsq4aw",
}

var testPrivKeys = []string{
	"MjQ1MDc2MmM4MjU5YjRhZjhhNmFjMmI0ZDBkNzBkOGE1ZTBmNDQ5NGI4NzM4OTYyM2E3MmI0OWMzNmE1ODZhNw==",
	"YmNiMzA2ODU1NWNjMzk3NDE1OWMwMTM3MDU0NTNjN2YwMzYzZmVhZDE5NmU3NzRhOTMwOWIxN2QyZTQ0MzdkNg==",
	"ZThiMDAxOTk2MDc4ODk3YWE0YThlMjdkMWY0NjA1MTAwZDgyNDkyYzdhNmMwZWQ3MDBhMWIyMjNmNGMzYjVhYg==",
	"ZTc2ZjI5OTIwOGVlMDk2N2M3Yzc1MjYyODQ0OGUyMjE3NGJiOGRmNGQyZmVmODg0NzQwNmUzYTk1YmQyODlmNA==",
}

type TSS struct {
	Server *tss.TssServer
	PubkeyInBytes []byte
	PubkeyInBech32 string
	AddressInHex string
}

func (tss *TSS) Pubkey() []byte {
	return tss.PubkeyInBytes
}

// digest should be Keccak256 Hash of some data
func (tss *TSS) Sign(digest []byte) ([65]byte, error) {

	H := digest
	log.Debug().Msgf("hash of digest is %s", H)

	tssPubkey := tss.PubkeyInBech32
	keysignReq := keysign.NewRequest(tssPubkey, []string{base64.StdEncoding.EncodeToString(H)}, 10, testPubKeys[:2], "0.14.0")
	ks_res, err := tss.Server.KeySign(keysignReq)
	if err != nil {
		log.Warn().Msg("keysign fail")
	}
	signature := ks_res.Signatures
	// [{cyP8i/UuCVfQKDsLr1kpg09/CeIHje1FU6GhfmyMD5Q= D4jXTH3/CSgCg+9kLjhhfnNo3ggy9DTQSlloe3bbKAs= eY++Z2LwsuKG1JcghChrsEJ4u9grLloaaFZNtXI3Ujk= AA==}]
	// 32B msg hash, 32B R, 32B S, 1B RC
	log.Info().Msgf("signature of digest is... %v", signature)

	if len(signature) == 0 {
		log.Warn().Err(err).Msgf("signature has length 0")
		return [65]byte{}, fmt.Errorf("keysign fail: %s", err)
	}
	if !verify_signature(tssPubkey, signature, H) {
		log.Error().Err(err).Msgf("signature verification failure")
		return [65]byte{}, fmt.Errorf("signuature verification fail")
	}
	var sigbyte [65]byte
	base64.StdEncoding.Decode(sigbyte[:32], []byte(signature[0].R))
	base64.StdEncoding.Decode(sigbyte[32:64], []byte(signature[0].S))
	base64.StdEncoding.Decode(sigbyte[64:65], []byte(signature[0].RecoveryID))

	return sigbyte, nil
}

func (tss *TSS) Address() ethcommon.Address {
	addr, err := getKeyAddr(tss.PubkeyInBech32)
	if err != nil {
		log.Error().Err(err).Msg("getKeyAddr error")
		return ethcommon.Address{}
	}
	return addr
}

func getKeyAddr(tssPubkey string) (ethcommon.Address, error ){
	var keyAddr ethcommon.Address
	pubk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		log.Fatal().Err(err)
		return keyAddr, err
	}
	//keyAddrBytes := pubk.Address().Bytes()
	pubk.Bytes()
	decompresspubkey, err := crypto.DecompressPubkey(pubk.Bytes())
	if err != nil {
		log.Fatal().Err(err).Msg("decompress err")
		return keyAddr, err
	}

	keyAddr = crypto.PubkeyToAddress(*decompresspubkey)
	//keyAddr = ethcommon.BytesToAddress(keyAddrBytes)
	log.Info().Msgf("keyAddr: %s", keyAddr.String())

	return keyAddr, nil
}

func NewTSS(peer addr.AddrList) (*TSS, error) {
	server, _, err  := SetupTSSServer(peer, "")
	if err != nil {
		return nil, fmt.Errorf("SetupTSSServer error: %w", err)
	}
	tss := TSS{
		Server: server,
	}
	tsspath := os.Getenv(("TSSPATH"))
	if len(tsspath) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("UserHomeDir")
			return nil, err
		}
		tsspath = filepath.Join(home, ".tss")
	}
	files, err := os.ReadDir(tsspath)
	if err != nil {
		return nil, err
	}
	found := false
	sharefiles := []os.DirEntry{}
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(filepath.Base(file.Name()), "localstate"){
			sharefiles = append(sharefiles, file)
		}
	}
	if len(sharefiles) > 0 {
		sort.SliceStable(sharefiles, func(i, j int) bool {
			fi, _ := sharefiles[i].Info()
			fj, _ := sharefiles[j].Info()
			return fi.ModTime().After(fj.ModTime())
		})
		localStateFile := sharefiles[0]
		filename := filepath.Base(localStateFile.Name())
		filearray := strings.Split(filename, "-")
		if len(filearray) == 2 {
			log.Info().Msgf("Found stored Pubkey in local state: %s", filearray[1])
			tss.PubkeyInBech32 = strings.TrimSuffix(filearray[1], ".json")
			found = true
		}
	}
	if !found {
		log.Info().Msgf("Local state not found; new Keygen starting...")
		var req keygen.Request
		req = keygen.NewRequest(testPubKeys[:2], 10, "0.14.0")
		res, err := server.Keygen(req)
		if err != nil {
			log.Fatal().Msg("keygen fail")
			return nil, fmt.Errorf("Keygen error: %w", err)
		}

		log.Info().Msgf("pubkey: %s", res.PubKey)
		tss.PubkeyInBech32 = res.PubKey
	}
	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, tss.PubkeyInBech32)
	if err != nil {
		log.Error().Err(err).Msgf("GetPubKeyFromBech32 from %s", tss.PubkeyInBech32)
		return nil, fmt.Errorf("GetPubKeyFromBech32: %w", err)
	}
	decompresspubkey, err := crypto.DecompressPubkey(pubkey.Bytes())
	tss.PubkeyInBytes = crypto.FromECDSAPub(decompresspubkey)
	log.Info().Msgf("pubkey.Bytes() gives %d Bytes", len(tss.PubkeyInBytes))

	tss.AddressInHex = crypto.PubkeyToAddress(*decompresspubkey).Hex()
	log.Info().Msgf("TSS Address ETH %s", tss.AddressInHex)

	return &tss, nil
}

func SetupTSSServer(peer addr.AddrList, tssAddr string) (*tss.TssServer, *TssHttpServer, error) {
	bootstrapPeers := peer

	log.Info().Msgf("Peers AddrList %v", bootstrapPeers)

	nodeIdxStr := os.Getenv("IDX")
	nodeIdx, err := strconv.Atoi(nodeIdxStr)
	if nodeIdxStr == "" || err != nil || nodeIdx < 0 || nodeIdx >= 4 {
		return nil, nil, fmt.Errorf("cannot get privkey from env IDX: %w", err)
	}
	priKeyBytes := testPrivKeys[nodeIdx]
	log.Debug().Msgf("test privkey is %s", priKeyBytes)
	priKey, err := conversion.GetPriKey(priKeyBytes)
	if err != nil {
		log.Err(err).Msg("GetPriKey")
		return nil, nil, err
	}
	tsspath := os.Getenv("TSSPATH")
	if len(tsspath) == 0 {
		log.Err(err).Msg("empty env TSSPATH")
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Error().Err(err).Msgf("cannot get UserHomeDir")
			return nil, nil, err
		}
		tsspath = path.Join(homedir, ".tss")
		log.Info().Msgf("create temporary TSSPATH: %s", tsspath)
	}
	IP := os.Getenv("MYIP")
	if len(IP) == 0 {
		log.Err(err).Msg("empty env MYIP")
	}
	tssServer, err := tss.NewTss(
		bootstrapPeers,
		6668,
		priKey,
		"MetaMetaOpenTheDoor",
		tsspath,
		thorcommon.TssConfig{
			EnableMonitor:   true,
			KeyGenTimeout:   300 * time.Second, // must be shorter than constants.JailTimeKeygen
			KeySignTimeout:  60 * time.Second,  // must be shorter than constants.JailTimeKeysign
			PartyTimeout:    45 * time.Second,
			PreParamTimeout: 5 * time.Minute,
		},
		nil, // don't set to precomputed values
		IP,  // for docker test
	)

	tssServer.Start()

	s := NewTssHttpServer(tssAddr, tssServer)
	go func() {
		if err := s.Start(); err != nil {
			fmt.Println(err)
		}
	}()

	log.Info().Msgf("LocalID: %v", tssServer.GetLocalPeerID())
	return tssServer, s, nil
}

func TestKeysign(tssPubkey string, tssServer *tss.TssServer) {
	log.Info().Msg("trying keysign...")
	data := []byte("hello meta")
	H := crypto.Keccak256Hash(data)
	log.Info().Msgf("hash of data (hello meta) is %s", H)

	keysignReq := keysign.NewRequest(tssPubkey, []string{base64.StdEncoding.EncodeToString(H.Bytes())}, 10, testPubKeys, "0.14.0")
	ks_res, err := tssServer.KeySign(keysignReq)
	if err != nil {
		log.Warn().Msg("keysign fail")
	}
	signature := ks_res.Signatures
	// [{cyP8i/UuCVfQKDsLr1kpg09/CeIHje1FU6GhfmyMD5Q= D4jXTH3/CSgCg+9kLjhhfnNo3ggy9DTQSlloe3bbKAs= eY++Z2LwsuKG1JcghChrsEJ4u9grLloaaFZNtXI3Ujk= AA==}]
	// 32B msg hash, 32B R, 32B S, 1B RC
	log.Info().Msgf("signature of helloworld... %v", signature)

	if len(signature) == 0 {
		log.Info().Msgf("signature has length, skipping verify", signature)
	} else {
		verify_signature(tssPubkey, signature, H.Bytes())
	}
}


func TestKeygen(tssServer *tss.TssServer) keygen.Response {
	// check if we already have LocalState persisted in files
	var req keygen.Request
	req = keygen.NewRequest(testPubKeys[:2], 10, "0.13.0")
	res, err := tssServer.Keygen(req)
	if err != nil {
		log.Fatal().Msg("keygen fail")
	}

	log.Info().Msgf("pubkey: %s", res.PubKey)

	return res
}


func verify_signature(tssPubkey string, signature []keysign.Signature, H []byte) bool {
	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		log.Fatal().Msg("get pubkey from bech32 fail")
	}
	// verify the signature of msg.
	var sigbyte [65]byte
	base64.StdEncoding.Decode(sigbyte[:32], []byte(signature[0].R))
	base64.StdEncoding.Decode(sigbyte[32:64], []byte(signature[0].S))
	base64.StdEncoding.Decode(sigbyte[64:65], []byte(signature[0].RecoveryID))
	sigPublicKey, err := crypto.SigToPub(H, sigbyte[:])
	compressedPubkey := crypto.CompressPubkey(sigPublicKey)
	log.Debug().Msgf("pubkey %s recovered pubkey %s", pubkey.String(), hex.EncodeToString(compressedPubkey))
	return bytes.Compare(pubkey.Bytes(), compressedPubkey) == 0
}
