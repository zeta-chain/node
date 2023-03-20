package main

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"strconv"
	"strings"
)

func AddObserverAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-observer-list [observer-list.json] ",
		Short: "Add a list of observers to the observer mapper",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			observerInfo, err := types.ParsefileToObserverDetails(args[0])
			if err != nil {
				return err
			}
			var observerMapper []*types.ObserverMapper
			var grantAuthorizations []authz.GrantAuthorization
			observersforChain := map[int64][]string{}
			// Generate the grant authorizations and created observer list for chain
			for _, info := range observerInfo {
				grantAuthorizations = append(grantAuthorizations, generateGrants(info)...)
				for _, chain := range info.SupportedChainsList {
					observersforChain[chain] = append(observersforChain[chain], info.ObserverAddress)
				}
			}
			// Generate observer mappers for each chain
			for chainID, observers := range observersforChain {
				observers = removeDuplicate(observers)
				chain := common.GetChainFromChainID(chainID)
				mapper := types.ObserverMapper{
					ObserverChain: chain,
					ObserverList:  observers,
				}
				observerMapper = append(observerMapper, &mapper)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			zetaObserverGenState := types.GetGenesisStateFromAppState(cdc, appState)
			zetaObserverGenState.Observers = observerMapper
			var authzGenState authz.GenesisState
			if appState[authz.ModuleName] != nil {
				err := cdc.UnmarshalJSON(appState[authz.ModuleName], &authzGenState)
				if err != nil {
					panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
				}
			}

			authzGenState.Authorization = grantAuthorizations

			zetaObserverStateBz, err := json.Marshal(zetaObserverGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal Observer List into Genesis File: %w", err)
			}

			err = codectypes.UnpackInterfaces(authzGenState, cdc)
			if err != nil {
				return fmt.Errorf("failed to authz grants into upackeder: %w", err)
			}
			authZStateBz, err := cdc.MarshalJSON(&authzGenState)
			if err != nil {
				return fmt.Errorf("failed to authz grants into Genesis File: %w", err)
			}
			appState[types.ModuleName] = zetaObserverStateBz
			appState[authz.ModuleName] = authZStateBz
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}
	return cmd
}

func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func generateGrants(info types.ObserverInfoReader) (grants []authz.GrantAuthorization) {

	grants = append(append(append(grants, addStakingGrants(grants, info)...),
		addSpendingGrants(grants, info)...),
		addZetaClientGrants(grants, info)...)
	return grants
}

func addZetaClientGrants(grants []authz.GrantAuthorization, info types.ObserverInfoReader) []authz.GrantAuthorization {
	txTypes := crosschaintypes.GetAllAuthzTxTypes()
	for _, txType := range txTypes {
		auth, err := codectypes.NewAnyWithValue(authz.NewGenericAuthorization(txType))
		if err != nil {
			panic(err)
		}
		grants = append(grants, authz.GrantAuthorization{
			Granter:       info.ObserverAddress,
			Grantee:       info.ZetaClientGranteeAddress,
			Authorization: auth,
			Expiration:    nil,
		})
	}
	return grants
}

func addSpendingGrants(grants []authz.GrantAuthorization, info types.ObserverInfoReader) []authz.GrantAuthorization {
	spendMaxTokens, ok := sdk.NewIntFromString(info.SpendMaxTokens)
	if !ok {
		panic("Failed to parse spend max tokens")
	}
	spendAuth, err := codectypes.NewAnyWithValue(&banktypes.SendAuthorization{
		SpendLimit: sdk.NewCoins(sdk.NewCoin(config.BaseDenom, spendMaxTokens)),
	})
	if err != nil {
		panic(err)
	}
	grants = append(grants, authz.GrantAuthorization{
		Granter:       info.ObserverAddress,
		Grantee:       info.SpendGranteeAddress,
		Authorization: spendAuth,
		Expiration:    nil,
	})
	return grants
}

func addStakingGrants(grants []authz.GrantAuthorization, info types.ObserverInfoReader) []authz.GrantAuthorization {
	stakingMaxTokens, ok := sdk.NewIntFromString(info.StakingMaxTokens)
	if !ok {
		panic("Failed to parse staking max tokens")
	}
	alllowList := stakingtypes.StakeAuthorization_AllowList{AllowList: &stakingtypes.StakeAuthorization_Validators{Address: info.StakingValidatorAllowList}}

	stakingAuth, err := codectypes.NewAnyWithValue(&stakingtypes.StakeAuthorization{
		MaxTokens:         &sdk.Coin{Denom: config.BaseDenom, Amount: stakingMaxTokens},
		Validators:        &alllowList,
		AuthorizationType: stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE,
	})
	if err != nil {
		panic(err)
	}
	grants = append(grants, authz.GrantAuthorization{
		Granter:       info.ObserverAddress,
		Grantee:       info.StakingGranteeAddress,
		Authorization: stakingAuth,
		Expiration:    nil,
	})
	delAuth, err := codectypes.NewAnyWithValue(&stakingtypes.StakeAuthorization{
		MaxTokens:         &sdk.Coin{Denom: config.BaseDenom, Amount: stakingMaxTokens},
		Validators:        &alllowList,
		AuthorizationType: stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_UNDELEGATE,
	})
	if err != nil {
		panic(err)
	}
	grants = append(grants, authz.GrantAuthorization{
		Granter:       info.ObserverAddress,
		Grantee:       info.StakingGranteeAddress,
		Authorization: delAuth,
		Expiration:    nil,
	})
	reDelauth, err := codectypes.NewAnyWithValue(&stakingtypes.StakeAuthorization{
		MaxTokens:         &sdk.Coin{Denom: config.BaseDenom, Amount: stakingMaxTokens},
		Validators:        &alllowList,
		AuthorizationType: stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_REDELEGATE,
	})
	if err != nil {
		panic(err)
	}
	grants = append(grants, authz.GrantAuthorization{
		Granter:       info.ObserverAddress,
		Grantee:       info.StakingGranteeAddress,
		Authorization: reDelauth,
		Expiration:    nil,
	})
	return grants

}

// AddObserverAccountCmd Deprecated : Use AddObserverAccountsCmd instead
func AddObserverAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-observer [chainID] [comma separate list of address] ",
		Short: "Add a list of observers to the observer mapper",
		Long: `
           Chain Types :
					"Empty"     
					"Eth"       
					"ZetaChain" 
					"Btc"       
					"Polygon"   
					"BscMainnet"
					"Goerli"    
					"Mumbai"    
					"Ropsten"   
					"Ganache"   
					"Baobap"    
					"BscTestnet"
					"BTCTestnet"
			`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			chainID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			chain := common.GetChainFromChainID(chainID)
			observer := &types.ObserverMapper{
				Index:         "",
				ObserverChain: chain,
				ObserverList:  strings.Split(args[1], ","),
			}
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			zetaObserverGenState := types.GetGenesisStateFromAppState(cdc, appState)
			zetaObserverGenState.Observers = append(zetaObserverGenState.Observers, observer)

			zetaObserverStateBz, err := json.Marshal(zetaObserverGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal Observer List into Genesis File: %w", err)
			}
			appState[types.ModuleName] = zetaObserverStateBz
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}
	return cmd
}
