package types

import (
	"fmt"
	// this line is used by starport scaffolding # ibc/genesistype/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # ibc/genesistype/default
		ZetaConversionRateList: []ZetaConversionRate{},
		OutTxTrackerList:       []OutTxTracker{},
		// this line is used by starport scaffolding # genesis/types/default
		Keygen:              nil,
		TSSVoterList:        []*TSSVoter{},
		TSSList:             []*TSS{},
		InTxList:            []*InTx{},
		TxList:              &TxList{Tx: []*Tx{}},
		GasBalanceList:      []*GasBalance{},
		GasPriceList:        []*GasPrice{},
		ChainNoncesList:     []*ChainNonces{},
		LastBlockHeightList: []*LastBlockHeight{},
		ReceiveList:         []*Receive{},
		SendList:            []*Send{},
		NodeAccountList:     []*NodeAccount{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # ibc/genesistype/validate

	// Check for duplicated index in zetaConversionRate
	zetaConversionRateIndexMap := make(map[string]struct{})

	for _, elem := range gs.ZetaConversionRateList {
		index := string(ZetaConversionRateKey(elem.Index))
		if _, ok := zetaConversionRateIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for zetaConversionRate")
		}
		zetaConversionRateIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in outTxTracker
	outTxTrackerIndexMap := make(map[string]struct{})

	for _, elem := range gs.OutTxTrackerList {
		index := string(OutTxTrackerKey(elem.Index))
		if _, ok := outTxTrackerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for outTxTracker")
		}
		outTxTrackerIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated index in tSSVoter
	tSSVoterIndexMap := make(map[string]bool)

	for _, elem := range gs.TSSVoterList {
		if _, ok := tSSVoterIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for tSSVoter")
		}
		tSSVoterIndexMap[elem.Index] = true
	}
	// Check for duplicated index in tSS
	tSSIndexMap := make(map[string]bool)

	for _, elem := range gs.TSSList {
		if _, ok := tSSIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for tSS")
		}
		tSSIndexMap[elem.Index] = true
	}
	// Check for duplicated index in inTx
	inTxIndexMap := make(map[string]bool)

	for _, elem := range gs.InTxList {
		if _, ok := inTxIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for inTx")
		}
		inTxIndexMap[elem.Index] = true
	}
	// Check for duplicated index in gasBalance
	gasBalanceIndexMap := make(map[string]bool)

	for _, elem := range gs.GasBalanceList {
		if _, ok := gasBalanceIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for gasBalance")
		}
		gasBalanceIndexMap[elem.Index] = true
	}
	// Check for duplicated index in gasPrice
	gasPriceIndexMap := make(map[string]bool)

	for _, elem := range gs.GasPriceList {
		if _, ok := gasPriceIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for gasPrice")
		}
		gasPriceIndexMap[elem.Index] = true
	}
	// Check for duplicated index in chainNonces
	chainNoncesIndexMap := make(map[string]bool)

	for _, elem := range gs.ChainNoncesList {
		if _, ok := chainNoncesIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for chainNonces")
		}
		chainNoncesIndexMap[elem.Index] = true
	}
	// Check for duplicated index in lastBlockHeight
	lastBlockHeightIndexMap := make(map[string]bool)

	for _, elem := range gs.LastBlockHeightList {
		if _, ok := lastBlockHeightIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for lastBlockHeight")
		}
		lastBlockHeightIndexMap[elem.Index] = true
	}
	// Check for duplicated index in receive
	receiveIndexMap := make(map[string]bool)

	for _, elem := range gs.ReceiveList {
		if _, ok := receiveIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for receive")
		}
		receiveIndexMap[elem.Index] = true
	}
	// Check for duplicated index in send
	sendIndexMap := make(map[string]bool)

	for _, elem := range gs.SendList {
		if _, ok := sendIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for send")
		}
		sendIndexMap[elem.Index] = true
	}

	// Check for duplicated index in nodeAccount
	nodeAccountIndexMap := make(map[string]bool)

	for _, elem := range gs.NodeAccountList {
		if _, ok := nodeAccountIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for nodeAccount")
		}
		nodeAccountIndexMap[elem.Index] = true
	}

	return nil
}
