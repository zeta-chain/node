package solana

import (
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/token"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func DeserializePdaInfo(pdaInfo *solrpc.GetAccountInfoResult) (PdaInfo, error) {
	pda := PdaInfo{}
	err := borsh.Deserialize(&pda, pdaInfo.Bytes())
	if err != nil {
		return PdaInfo{}, err
	}

	return pda, nil
}

func DeserializeMintAccountInfo(mintInfo *solrpc.GetAccountInfoResult) (token.Mint, error) {
	var mint token.Mint
	// Account{}.Data.GetBinary() returns the *decoded* binary data
	// regardless the original encoding (it can handle them all).
	err := bin.NewBinDecoder(mintInfo.Value.Data.GetBinary()).Decode(&mint)
	if err != nil {
		return token.Mint{}, err
	}

	return mint, err
}
