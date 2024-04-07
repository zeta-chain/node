package sample

import (
	"github.com/zeta-chain/zetacore/pkg/proofs"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
)

func BlockHeader(blockHash []byte) proofs.BlockHeader {
	return proofs.BlockHeader{
		Height:     42,
		Hash:       blockHash,
		ParentHash: Hash().Bytes(),
		ChainId:    42,
		Header:     proofs.HeaderData{},
	}
}

func ChainState(chainID int64) lightclienttypes.ChainState {
	return lightclienttypes.ChainState{
		ChainId:         chainID,
		LatestHeight:    42,
		EarliestHeight:  42,
		LatestBlockHash: Hash().Bytes(),
	}
}
