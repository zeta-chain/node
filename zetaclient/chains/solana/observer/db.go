package observer

import (
	"github.com/pkg/errors"
)

// LoadDB open sql database and load data into Solana observer
func (ob *Observer) LoadDB(dbPath string) error {
	if dbPath == "" {
		return errors.New("empty db path")
	}

	// open database
	err := ob.OpenDB(dbPath, "")
	if err != nil {
		return errors.Wrapf(err, "error OpenDB for chain %d", ob.Chain().ChainId)
	}

	ob.Observer.LoadLastTxScanned(ob.Logger().Chain)

	return nil
}

// LoadLastTxScanned loads the last scanned tx from the database.
func (ob *Observer) LoadLastTxScanned() error {
	ob.Observer.LoadLastTxScanned(ob.Logger().Chain)
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from tx %s", ob.Chain().ChainId, ob.LastTxScanned())

	return nil
}
