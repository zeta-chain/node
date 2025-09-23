package repo

import (
	// "context"
	"errors"
	// eth "github.com/ethereum/go-ethereum/common"
	// "github.com/tonkeeper/tongo/ton"
	//
	// "github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

var (
	ErrInvalidOutboundSigner = errors.New("outbound signer is not TSS")
)

// // TODO
// func (repo *TONRepo) TSSSignerAddress() eth.Address {
// 	return repo.TSSSigner.PubKey().AddressEVM()
// }
//
// // TODO
// func (repo *TONRepo) PostOutboundTracker(ctx context.Context,
// 	chainID int64,
// 	nonce uint64,
// 	tx ton.Transaction,
// ) error {
// 	hash := rpc.TransactionToHashString(tx)
// 	_, err := repo.ZetacoreClient.PostOutboundTracker(ctx, chainID, nonce, hash)
// 	return err
// }
