package common

import "fmt"

type Chain interface {
	Start()
	ID() uint16
	Name() string
}

var ALL_CHAINS []Chain

func FindChainByID(id uint16) (Chain, error) {
	for _, chain := range ALL_CHAINS {
		if chain.ID() == id {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("Not listening for chain with ID: %d", id)
}

func FindChainByName(name string) (Chain, error) {
	for _, chain := range ALL_CHAINS {
		if chain.Name() == name {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find chain: %s", name)
}
