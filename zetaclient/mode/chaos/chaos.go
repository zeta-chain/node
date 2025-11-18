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

	ErrNotChaosMode      = errors.New("not in chaos mode")
	ErrReadPercentages   = errors.New("failed to read chaos percentages")
	ErrParsePercentages  = errors.New("failed to parse chaos percentages")
	ErrInvalidPercentage = errors.New("invalid percentage")
)

// Source is the base chaos object from which all chaos interface implementations inherit.
// It is safe for concurrent use by multiple goroutines.
type Source struct {
	mu      sync.Mutex
	profile map[string](map[string]int) // map[itfc](map[mthd]percentage)
	seed    int64
	rand    *rand.Rand
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

	// Read the file with the fail profile.
	data, err := os.ReadFile(config.ChaosProfilePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadPercentages, err)
	}

	// Parse the file with the fail profile.
	profile := make(map[string](map[string]int))
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParsePercentages, err)
	}

	// Validate percentages.
	for itfc, mthds := range profile {
		for mthd, percentage := range mthds {
			if percentage < 0 || percentage > 100 {
				return nil, fmt.Errorf("%w for method %q in %q", ErrInvalidPercentage, mthd, itfc)
			}
		}
	}

	return &Source{profile: profile, seed: seed, rand: rand}, nil
}

// shouldFail determines whether a method should fail based on its failure percentage.
//
// It generates a random integer in the range [0, 100). If the generated number is less than the
// method's fail percentage (stored in the percentages map), shouldFail returns an error; otherwise,
// it returns nil.
func (source *Source) shouldFail(itfc, mthd string) error {
	source.mu.Lock()
	defer source.mu.Unlock()
	n := source.rand.Intn(100)
	p := source.profile[itfc][mthd]
	if n < p {
		return fmt.Errorf("%w (seed: %d): %s.%s (%d < %d)", ErrChaos, source.seed, itfc, mthd, n, p)
	}
	return nil
}

// SetSelf forwards the call to the underlying client if it has the method to enable chaos mode to intercept internal calls
func (self *chaosZetacoreClient) SetSelf(queryTxResulter interface{}) {
	type selfSetter interface {
		SetSelf(interface{})
	}
	if setter, ok := self.client.(selfSetter); ok {
		setter.SetSelf(queryTxResulter)
	}
}
