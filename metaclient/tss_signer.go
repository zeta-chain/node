package metaclient

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	thorcommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-peerstore/addr"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/tss/go-tss/conversion"
	"gitlab.com/thorchain/tss/go-tss/keysign"
	"gitlab.com/thorchain/tss/go-tss/tss"
	"os"
	"strings"
	"time"
)

//var testPubKeys = []string{
//	"thorpub1addwnpepqtdklw8tf3anjz7nn5fly3uvq2e67w2apn560s4smmrt9e3x52nt2svmmu3",
//	"thorpub1addwnpepqtspqyy6gk22u37ztra4hq3hdakc0w0k60sfy849mlml2vrpfr0wvm6uz09",
//	"thorpub1addwnpepq2ryyje5zr09lq7gqptjwnxqsy2vcdngvwd6z7yt5yjcnyj8c8cn559xe69",
//	"thorpub1addwnpepqfjcw5l4ay5t00c32mmlky7qrppepxzdlkcwfs2fd5u73qrwna0vzag3y4j",
//}
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

type addrList []maddr.Multiaddr

// String implement fmt.Stringer
func (al *addrList) String() string {
	addresses := make([]string, len(*al))
	for i, address := range *al {
		addresses[i] = address.String()
	}
	return strings.Join(addresses, ",")
}

// Set add the given value to addList
func (al *addrList) Set(value string) error {
	address, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, address)
	return nil
}

func SetupTSSServer(peer addr.AddrList, tssAddr string) (*tss.TssServer, *TssHttpServer, error) {
	bootstrapPeers := peer

	log.Info().Msgf("Peers AddrList %v", bootstrapPeers)

	//conversion.SetupBech32Prefix()

	// Read stdin for the private key
	//inBuf := bufio.NewReader(os.Stdin)
	//priKeyBytes, err := input.GetPassword("input node secret key:", inBuf)
	//if err != nil {
	//	log.Err(err).Msg("GetPassword fail")
	//	return nil, nil, err
	//}
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
		tsspath, err = os.MkdirTemp("", "tsspath")
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
		verify_signature(err, tssPubkey, signature, H)
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
	log.Info().Msgf("persist to")

	return res
}


func verify_signature(err error, tssPubkey string, signature []keysign.Signature, H ethcommon.Hash) {
	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		log.Fatal().Msg("get pubkey from bech32 fail")
	}
	log.Info().Msgf("tss pubkey %s, size in B: %d", hex.EncodeToString(pubkey.Bytes()),
		len(pubkey.Bytes()))
	// verify the signature of msg.
	var sigbyte [65]byte
	base64.StdEncoding.Decode(sigbyte[:32], []byte(signature[0].R))
	base64.StdEncoding.Decode(sigbyte[32:64], []byte(signature[0].S))
	base64.StdEncoding.Decode(sigbyte[64:65], []byte(signature[0].RecoveryID))
	sigPublicKey, err := crypto.SigToPub(H.Bytes(), sigbyte[:])
	log.Info().Msgf("tss pubkey recovered in bytes: %s", hex.EncodeToString(crypto.CompressPubkey(sigPublicKey)))
}
