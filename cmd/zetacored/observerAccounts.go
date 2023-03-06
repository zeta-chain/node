package main

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
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
				grantAuthorizations = append(grantAuthorizations, generateGrants(info.ObserverAddress, info.ObserverGranteeAddress)...)
				for _, chain := range info.SupportedChainsList {
					observersforChain[chain] = append(observersforChain[chain], info.ObserverAddress)
				}
			}
			// Generate observer mappers for each chain
			for chainId, observers := range observersforChain {
				observers = removeDuplicate(observers)
				chain := common.GetChainFromChainID(chainId)
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

func generateGrants(granter, grantee string) (grants []authz.GrantAuthorization) {
	txTypes := crosschaintypes.GetAllAuthzTxTypes()
	for _, txType := range txTypes {
		auth, err := codectypes.NewAnyWithValue(authz.NewGenericAuthorization(txType))
		if err != nil {
			panic(err)
		}
		grants = append(grants, authz.GrantAuthorization{
			Granter:       granter,
			Grantee:       grantee,
			Authorization: auth,
			Expiration:    nil,
		})
	}
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
