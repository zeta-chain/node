package observer

import (
	"github.com/pkg/errors"

	solanarpc "github.com/zeta-chain/zetacore/zetaclient/chains/solana/rpc"
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

	// load last scanned tx
	err = ob.LoadLastTxScanned()

	return err
}

// LoadLastTxScanned loads the last scanned tx from the database.
func (ob *Observer) LoadLastTxScanned() error {
	ob.Observer.LoadLastTxScanned(ob.Logger().Chain)

	// when last scanned tx is absent in the database, the observer will scan from the 1st signature for the gateway address.
	// this is useful when bootstrapping the Solana observer
	if ob.LastTxScanned() == "" {
		firstSigature, err := solanarpc.GetFirstSignatureForAddress(
			ob.solClient,
			ob.gatewayID,
			solanarpc.DefaultPageLimit,
		)
		if err != nil {
			return err
		}
		ob.WithLastTxScanned(firstSigature.String())
	}
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from tx %s", ob.Chain().ChainId, ob.LastTxScanned())

	return nil
}
