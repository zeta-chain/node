package metaclientd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestWatchRouter(t *testing.T) {
	t.Logf("Settting up test...")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Logf("UserHomeDir error")
		t.Fail()
	}
	t.Logf("user home dir: %s", homeDir)
	chainHomeFoler := filepath.Join(homeDir, ".metacore")
	t.Logf("chain home dir: %s", chainHomeFoler)

	signerName := "alice"
	signerPass := "password"
	kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get keyring keybase")
	}

	k := NewKeysWithKeybase(kb, signerName, signerPass)

	chainIP := "127.0.0.1"
	bridge, err := NewMetachainBridge(k, chainIP, "alice")
	if err != nil {
		t.Fail()
	}

	EthObserver := &ChainObserver{}
	EthObserver.InitChainObserver("Ethereum", bridge)
	//EthObserver.WatchRouter()
}
