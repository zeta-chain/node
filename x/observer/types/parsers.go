package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"io/ioutil"
	"path/filepath"
)

type ObserverMapperReader struct {
	Index             string   `json:"index"`
	ObserverChainName string   `json:"observerChainName"`
	ObserverChainId   int64    `json:"observerChainId"`
	ObservationType   string   `json:"observationType"`
	ObserverList      []string `json:"observerList"`
}

func ParsefileToObserverMapper(fp string) ([]*ObserverMapper, error) {
	var observers []ObserverMapperReader
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	file = filepath.Clean(file)
	input, err := ioutil.ReadFile(file) // #nosec G304
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &observers)
	if err != nil {
		return nil, err
	}

	observerMappers := make([]*ObserverMapper, len(observers))
	for i, readerValue := range observers {
		chain := &Chain{
			ChainName: ParseStringToObserverChain(readerValue.ObserverChainName),
			ChainId:   readerValue.ObserverChainId,
		}
		observationType := ParseStringToObservationType(readerValue.ObservationType)
		if observationType == 0 || chain.ChainName == 0 {
			return nil, errors.Wrap(ErrUnableToParseMapper, fmt.Sprintf("Chain %s | ObserVation %s", readerValue.ObserverChainName, readerValue.ObservationType))
		}
		observerMappers[i] = &ObserverMapper{
			Index:           readerValue.Index,
			ObserverChain:   chain,
			ObservationType: ParseStringToObservationType(readerValue.ObservationType),
			ObserverList:    readerValue.ObserverList,
		}
	}
	return observerMappers, nil
}

// This should not be needed after removing common.Chain
func ParseCommonChaintoObservationChain(chain string) Chain {
	return Chain{}
	//commonChain := common.Chain(chain)
	//switch commonChain {
	//// Mainnet Chains
	//case common.ZETAChain, common.Chain(strings.ToUpper(string(common.ZETAChain))):
	//	return ObserverChain_ZetaChain
	//case common.ETHChain, common.Chain(strings.ToUpper(string(common.ETHChain))):
	//	return ObserverChain_Eth
	//case common.BSCChain, common.Chain(strings.ToUpper(string(common.BSCChain))):
	//	return ObserverChain_BscMainnet
	//case common.POLYGONChain, common.Chain(strings.ToUpper(string(common.POLYGONChain))):
	//	return ObserverChain_Polygon
	//// Testnet Chains
	//case common.MumbaiChain, common.Chain(strings.ToUpper(string(common.MumbaiChain))):
	//	return ObserverChain_Mumbai
	//case common.BaobabChain, common.Chain(strings.ToUpper(string(common.BaobabChain))):
	//	return ObserverChain_Baobab
	//case common.RopstenChain, common.Chain(strings.ToUpper(string(common.RopstenChain))):
	//	return ObserverChain_Ropsten
	//case common.GoerliChain, common.Chain(strings.ToUpper(string(common.GoerliChain))):
	//	return ObserverChain_Goerli
	//case common.BSCTestnetChain, common.Chain(strings.ToUpper(string(common.BSCTestnetChain))):
	//	return ObserverChain_BscTestnet
	//}
	//return ObserverChain_Empty
}

func ConvertReceiveStatusToVoteType(status common.ReceiveStatus) VoteType {
	switch status {
	case common.ReceiveStatus_Success:
		return VoteType_SuccessObservation
	case common.ReceiveStatus_Failed:
		return VoteType_FailureObservation
	default:
		return VoteType_NotYetVoted
	}
}

func ParseStringToObserverChain(chain string) ChainName {
	c := ChainName_value[chain]
	return ChainName(c)
}

func ParseStringToObservationType(observationType string) ObservationType {
	c := ObservationType_value[observationType]
	return ObservationType(c)
}

func GetOperatorAddressFromAccAddress(accAddr string) (sdk.ValAddress, error) {
	accAddressBech32, err := sdk.AccAddressFromBech32(accAddr)
	if err != nil {
		return nil, err
	}
	valAddress := sdk.ValAddress(accAddressBech32)
	valAddressBech32, err := sdk.ValAddressFromBech32(valAddress.String())
	if err != nil {
		return nil, err
	}
	return valAddressBech32, nil
}

func GetAccAddressFromOperatorAddress(valAddress string) (sdk.AccAddress, error) {
	valAddressBech32, err := sdk.ValAddressFromBech32(valAddress)
	if err != nil {
		return nil, err
	}
	accAddress := sdk.AccAddress(valAddressBech32)
	accAddressBech32, err := sdk.AccAddressFromBech32(accAddress.String())
	if err != nil {
		return nil, err
	}
	return accAddressBech32, nil
}
