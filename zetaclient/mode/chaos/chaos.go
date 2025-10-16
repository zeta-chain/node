// Package chaos provides chaos-wrappers for the TSS signer, for the zetacore client, and for
// the standard clients of the connected chains.
//
// A chaos-wrapper overrides methods from the underlying client that return at least one error,
// that is, they override methods that can fail. The wrapper methods may call their inner
// counterparts or, depending on configured failure percentages, return ErrChaos.
package chaos

//go:generate go run generate/main.go
//go:generate gofmt -w generated.go

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/zetaclient/config"
)

var (
	// ErrChaos is the error that gets returned by the wrapped methods when they fail.
	ErrChaos = errors.New("chaos error")

	ErrNotChaosMode     = errors.New("not in chaos mode")
	ErrReadPercentages  = errors.New("failed to read chaos percentages")
	ErrParsePercentages = errors.New("failed to parse chaos percentages")
)

// Source is the base chaos object from which all chaos interface implementations inherit.
// It is safe for concurrent use by multiple goroutines.
type Source struct {
	mu          sync.Mutex
	percentages map[string](map[string]int) // map[itfc](map[mthd]percentage)
	rand        *rand.Rand
}

// NewSource parses the universal configuration into the source chaos object.
func NewSource(logger zerolog.Logger, config config.Config) (*Source, error) {
	if !config.ClientMode.IsChaosMode() {
		return nil, ErrNotChaosMode
	}

	// Set a random seed in case one was not provided.
	seed := config.ChaosSeed
	if seed == 0 {
		seed = time.Now().UnixNano()
		logger.Info().Int64("seed", seed).Msg("using a random chaos seed")
	}
	rand := rand.New(rand.NewSource(seed)) // #nosec G404 -- This is intended

	// Read the file with the fail percentages.
	data, err := os.ReadFile(config.ChaosPercentagesPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadPercentages, err)
	}

	// Parse the file with the fail percentages.
	percentages := make(map[string](map[string]int))
	err = json.Unmarshal(data, &percentages)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParsePercentages, err)
	}

	return &Source{percentages: percentages, rand: rand}, nil
}

// shouldFail returns true if and only if a given method from a given interface should fail.
func (source *Source) shouldFail(itfc, mthd string) bool {
	source.mu.Lock()
	defer source.mu.Unlock()
	return source.rand.Intn(100) < source.percentages[itfc][mthd]
}
