//go:build PRIVNET

package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// This is for privnet/testnet only
func (k Keeper) BlockOneDeploySystemContracts(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// setup uniswap v2 factory
	uniswapV2Factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployUniswapV2Factory")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("UniswapV2Factory", uniswapV2Factory.String()),
		),
	)

	// setup WZETA contract
	wzeta, err := k.DeployWZETA(ctx)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployWZetaContract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployWZetaContract", wzeta.String()),
		),
	)

	router, err := k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployUniswapV2Router02")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployUniswapV2Router02", router.String()),
		),
	)

	connector, err := k.DeployConnectorZEVM(ctx, wzeta)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployConnectorZEVM")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployConnectorZEVM", connector.String()),
		),
	)
	ctx.Logger().Info("Deployed Connector ZEVM at " + connector.String())

	SystemContractAddress, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, router)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to SystemContractAddress")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("SystemContractAddress", SystemContractAddress.String()),
		),
	)

	// set the system contract
	system, _ := k.GetSystemContract(ctx)
	system.SystemContract = SystemContractAddress.String()
	// FIXME: remove unnecessary SetGasPrice and setupChainGasCoinAndPool
	k.SetSystemContract(ctx, system)
	//err = k.SetGasPrice(ctx, big.NewInt(1337), big.NewInt(1))
	if err != nil {
		return err
	}

	_, err = k.setupChainGasCoinAndPool(ctx, common.GoerliChain().ChainId, "ETH", "gETH", 18)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
	}

	_, err = k.setupChainGasCoinAndPool(ctx, common.BtcRegtestChain().ChainId, "BTC", "tBTC", 8)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
	}

	//FIXME: clean up and config the above based on localnet/testnet/mainnet

	// for localnet only: USDT ZRC20
	USDTAddr := "0xff3135df4F2775f4091b81f4c7B6359CfA07862a"
	USDTZRC20Addr, err := k.DeployZRC20Contract(ctx, "USDT", "USDT", uint8(6), common.GoerliChain().ChainId, common.CoinType_ERC20, USDTAddr, big.NewInt(90_000))
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployZRC20Contract USDT")
	}
	ctx.Logger().Info("Deployed USDT ZRC20 at " + USDTZRC20Addr.String())
	// for localnet only: ZEVM Swap App
	//ZEVMSwapAddress, err := k.DeployZEVMSwapApp(ctx, router, SystemContractAddress)
	//if err != nil {
	//	return sdkerrors.Wrapf(err, "failed to deploy ZEVMSwapApp")
	//}
	//ctx.Logger().Info("Deployed ZEVM Swap App at " + ZEVMSwapAddress.String())
	fmt.Println("Successfully deployed contracts")
	return nil
}

func (k Keeper) TestUpdateSystemContractAddress(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to GetWZetaContractAddress")
	}
	uniswapV2Factory, err := k.GetUniswapv2FacotryAddress(ctx)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to GetUniswapv2FacotryAddress")
	}
	router, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to GetUniswapV2Router02Address")
	}

	SystemContractAddress, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, router)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeploySystemContract")
	}
	creator := k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_deploy_fungible_coin)
	msg := types.NewMessageUpdateSystemContract(creator, SystemContractAddress.Hex())
	_, err = k.UpdateSystemContract(ctx, msg)
	k.Logger(ctx).Info("System contract updated", "new address", SystemContractAddress.String())
	return err
}

func (k Keeper) TestUpdateZRC20WithdrawFee(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	foreignCoins := k.GetAllForeignCoins(ctx)
	creator := k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_deploy_fungible_coin)

	for _, foreignCoin := range foreignCoins {
		msg := types.NewMsgUpdateZRC20WithdrawFee(creator, foreignCoin.Zrc20ContractAddress, sdk.NewUint(uint64(foreignCoin.ForeignChainId)))
		_, err := k.UpdateZRC20WithdrawFee(ctx, msg)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	NewZRC20Bytecode = "0x608060405234801561001057600080fd5b50600436106101425760003560e01c806385e1f4d0116100b8578063c835d7cc1161007c578063c835d7cc146103a5578063d9eeebed146103c1578063dd62ed3e146103e0578063eddeb12314610410578063f2441b321461042c578063f687d12a1461044a57610142565b806385e1f4d0146102eb57806395d89b4114610309578063a3413d0314610327578063a9059cbb14610345578063c70126261461037557610142565b8063313ce5671161010a578063313ce567146102015780633ce4a5bc1461021f57806342966c681461023d57806347e7ef241461026d5780634d8943bb1461029d57806370a08231146102bb57610142565b806306fdde0314610147578063091d278814610165578063095ea7b31461018357806318160ddd146101b357806323b872dd146101d1575b600080fd5b61014f610466565b60405161015c9190611c04565b60405180910390f35b61016d6104f8565b60405161017a9190611c26565b60405180910390f35b61019d600480360381019061019891906118c5565b6104fe565b6040516101aa9190611b52565b60405180910390f35b6101bb61051c565b6040516101c89190611c26565b60405180910390f35b6101eb60048036038101906101e69190611872565b610526565b6040516101f89190611b52565b60405180910390f35b61020961061e565b6040516102169190611c41565b60405180910390f35b610227610635565b6040516102349190611ad7565b60405180910390f35b6102576004803603810190610252919061198e565b61064d565b6040516102649190611b52565b60405180910390f35b610287600480360381019061028291906118c5565b610662565b6040516102949190611b52565b60405180910390f35b6102a56107ce565b6040516102b29190611c26565b60405180910390f35b6102d560048036038101906102d091906117d8565b6107d4565b6040516102e29190611c26565b60405180910390f35b6102f361081d565b6040516103009190611c26565b60405180910390f35b610311610841565b60405161031e9190611c04565b60405180910390f35b61032f6108d3565b60405161033c9190611be9565b60405180910390f35b61035f600480360381019061035a91906118c5565b6108f7565b60405161036c9190611b52565b60405180910390f35b61038f600480360381019061038a9190611932565b610915565b60405161039c9190611b52565b60405180910390f35b6103bf60048036038101906103ba91906117d8565b610a6b565b005b6103c9610b5e565b6040516103d7929190611b29565b60405180910390f35b6103fa60048036038101906103f59190611832565b610dcb565b6040516104079190611c26565b60405180910390f35b61042a6004803603810190610425919061198e565b610e52565b005b610434610f0c565b6040516104419190611ad7565b60405180910390f35b610464600480360381019061045f919061198e565b610f30565b005b60606006805461047590611e8a565b80601f01602080910402602001604051908101604052809291908181526020018280546104a190611e8a565b80156104ee5780601f106104c3576101008083540402835291602001916104ee565b820191906000526020600020905b8154815290600101906020018083116104d157829003601f168201915b5050505050905090565b60015481565b600061051261050b610fea565b8484610ff2565b6001905092915050565b6000600554905090565b60006105338484846111ab565b6000600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600061057e610fea565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050828110156105f5576040517f10bad14700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61061285610601610fea565b858461060d9190611d9a565b610ff2565b60019150509392505050565b6000600860009054906101000a900460ff16905090565b73735b14bb79463307aacbed86daf3322b1e6226ab81565b60006106593383611407565b60019050919050565b600073735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015610700575060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15610737576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61074183836115bf565b8273ffffffffffffffffffffffffffffffffffffffff167f67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab373735b14bb79463307aacbed86daf3322b1e6226ab60405160200161079e9190611abc565b604051602081830303815290604052846040516107bc929190611b6d565b60405180910390a26001905092915050565b60025481565b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60606007805461085090611e8a565b80601f016020809104026020016040519081016040528092919081815260200182805461087c90611e8a565b80156108c95780601f1061089e576101008083540402835291602001916108c9565b820191906000526020600020905b8154815290600101906020018083116108ac57829003601f168201915b5050505050905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b600061090b610904610fea565b84846111ab565b6001905092915050565b6000806000610922610b5e565b915091508173ffffffffffffffffffffffffffffffffffffffff166323b872dd3373735b14bb79463307aacbed86daf3322b1e6226ab846040518463ffffffff1660e01b815260040161097793929190611af2565b602060405180830381600087803b15801561099157600080fd5b505af11580156109a5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109c99190611905565b6109ff576040517f0a7cd6d600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a093385611407565b3373ffffffffffffffffffffffffffffffffffffffff167f9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955868684600254604051610a579493929190611b9d565b60405180910390a260019250505092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610ae4576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507fd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae81604051610b539190611ad7565b60405180910390a150565b60008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630be155477f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b8152600401610bdd9190611c26565b60206040518083038186803b158015610bf557600080fd5b505afa158015610c09573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c2d9190611805565b9050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610c96576040517f78fff39600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d7fd7afb7f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b8152600401610d129190611c26565b60206040518083038186803b158015610d2a57600080fd5b505afa158015610d3e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d6291906119bb565b90506000811415610d9f576040517fe661aed000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060025460015483610db29190611d40565b610dbc9190611cea565b90508281945094505050509091565b6000600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610ecb576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806002819055507fef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f81604051610f019190611c26565b60405180910390a150565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610fa9576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806001819055507fff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a81604051610fdf9190611c26565b60405180910390a150565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611059576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156110c0576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258360405161119e9190611c26565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611212576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611279576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050818110156112f7576040517ffe382aa700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81816113039190611d9a565b600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546113959190611cea565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516113f99190611c26565b60405180910390a350505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561146e576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050818110156114ec576040517ffe382aa700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81816114f89190611d9a565b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816005600082825461154d9190611d9a565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516115b29190611c26565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611626576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600560008282546116389190611cea565b9250508190555080600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461168e9190611cea565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516116f39190611c26565b60405180910390a35050565b600061171261170d84611c81565b611c5c565b90508281526020810184848401111561172e5761172d611fd2565b5b611739848285611e48565b509392505050565b60008135905061175081612013565b92915050565b60008151905061176581612013565b92915050565b60008151905061177a8161202a565b92915050565b600082601f83011261179557611794611fcd565b5b81356117a58482602086016116ff565b91505092915050565b6000813590506117bd81612041565b92915050565b6000815190506117d281612041565b92915050565b6000602082840312156117ee576117ed611fdc565b5b60006117fc84828501611741565b91505092915050565b60006020828403121561181b5761181a611fdc565b5b600061182984828501611756565b91505092915050565b6000806040838503121561184957611848611fdc565b5b600061185785828601611741565b925050602061186885828601611741565b9150509250929050565b60008060006060848603121561188b5761188a611fdc565b5b600061189986828701611741565b93505060206118aa86828701611741565b92505060406118bb868287016117ae565b9150509250925092565b600080604083850312156118dc576118db611fdc565b5b60006118ea85828601611741565b92505060206118fb858286016117ae565b9150509250929050565b60006020828403121561191b5761191a611fdc565b5b60006119298482850161176b565b91505092915050565b6000806040838503121561194957611948611fdc565b5b600083013567ffffffffffffffff81111561196757611966611fd7565b5b61197385828601611780565b9250506020611984858286016117ae565b9150509250929050565b6000602082840312156119a4576119a3611fdc565b5b60006119b2848285016117ae565b91505092915050565b6000602082840312156119d1576119d0611fdc565b5b60006119df848285016117c3565b91505092915050565b6119f181611dce565b82525050565b611a08611a0382611dce565b611eed565b82525050565b611a1781611de0565b82525050565b6000611a2882611cb2565b611a328185611cc8565b9350611a42818560208601611e57565b611a4b81611fe1565b840191505092915050565b611a5f81611e36565b82525050565b6000611a7082611cbd565b611a7a8185611cd9565b9350611a8a818560208601611e57565b611a9381611fe1565b840191505092915050565b611aa781611e1f565b82525050565b611ab681611e29565b82525050565b6000611ac882846119f7565b60148201915081905092915050565b6000602082019050611aec60008301846119e8565b92915050565b6000606082019050611b0760008301866119e8565b611b1460208301856119e8565b611b216040830184611a9e565b949350505050565b6000604082019050611b3e60008301856119e8565b611b4b6020830184611a9e565b9392505050565b6000602082019050611b676000830184611a0e565b92915050565b60006040820190508181036000830152611b878185611a1d565b9050611b966020830184611a9e565b9392505050565b60006080820190508181036000830152611bb78187611a1d565b9050611bc66020830186611a9e565b611bd36040830185611a9e565b611be06060830184611a9e565b95945050505050565b6000602082019050611bfe6000830184611a56565b92915050565b60006020820190508181036000830152611c1e8184611a65565b905092915050565b6000602082019050611c3b6000830184611a9e565b92915050565b6000602082019050611c566000830184611aad565b92915050565b6000611c66611c77565b9050611c728282611ebc565b919050565b6000604051905090565b600067ffffffffffffffff821115611c9c57611c9b611f9e565b5b611ca582611fe1565b9050602081019050919050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600082825260208201905092915050565b6000611cf582611e1f565b9150611d0083611e1f565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611d3557611d34611f11565b5b828201905092915050565b6000611d4b82611e1f565b9150611d5683611e1f565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611d8f57611d8e611f11565b5b828202905092915050565b6000611da582611e1f565b9150611db083611e1f565b925082821015611dc357611dc2611f11565b5b828203905092915050565b6000611dd982611dff565b9050919050565b60008115159050919050565b6000819050611dfa82611fff565b919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b6000611e4182611dec565b9050919050565b82818337600083830152505050565b60005b83811015611e75578082015181840152602081019050611e5a565b83811115611e84576000848401525b50505050565b60006002820490506001821680611ea257607f821691505b60208210811415611eb657611eb5611f6f565b5b50919050565b611ec582611fe1565b810181811067ffffffffffffffff82111715611ee457611ee3611f9e565b5b80604052505050565b6000611ef882611eff565b9050919050565b6000611f0a82611ff2565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b60008160601b9050919050565b600381106120105761200f611f40565b5b50565b61201c81611dce565b811461202757600080fd5b50565b61203381611de0565b811461203e57600080fd5b50565b61204a81611e1f565b811461205557600080fd5b5056fea2646970667358221220edaf9ed98354e71aa84b95b4433f47537dd491d72f649020c367c23ec482327064736f6c63430008070033"
)

func (k Keeper) TestUpdateContractBytecode(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	creator := k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_deploy_fungible_coin)
	_ = creator
	////types.NewMsgUpdateContractBytecode(creator, )
	//addr, err := k.DeployZRC20Contract(ctx, "USDT", "USDT", uint8(6), common.GoerliChain().ChainId, common.CoinType_ERC20, "0xff3135df4F2775f4091b81f4c7B6359CfA07862a", big.NewInt(90_000))
	//if err != nil {
	//	panic(err)
	//}
	//if addr == (ethcommon.Address{}) {
	//	panic("empty address")
	//}
	//k.Logger(ctx).Info("Deployed USDT ZRC20 at " + addr.String())
	bytecode, err := hex.DecodeString(NewZRC20Bytecode[2:])
	_ = bytecode
	if err != nil {
		panic(err)
	}
	//msg := types.NewMsgUpdateContractBytecode(creator, addr, bytecode)
	//_, err = k.UpdateContractBytecode(ctx, msg)
	//if err != nil {
	//	panic(err)
	//}
	//k.Logger(ctx).Info("UpdateContractBytecode")
	//data, err := k.QueryZRC20Data(ctx, addr)
	//if err != nil {
	//	panic(err)
	//}
	//k.Logger(ctx).Info("QueryZRC20Data", "data", data)
	//zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	//if err != nil || zrc20ABI == nil {
	//	panic(err)
	//}

	//res, err := k.CallEVM(ctx, *zrc20ABI, types.ModuleAddressEVM, addr, big.NewInt(0), big.NewInt(90000), true,
	//	"increaseAllowance", addr, big.NewInt(1000000000000000000))
	//if err == nil {
	//	panic("increaseAllowance should fail")
	//}
	//k.Logger(ctx).Info("CallEVM", "res", res, "err", err, "revert code", res.Revert())
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil || zrc20ABI == nil {
		panic("zrc20ABI is nil")
	}
	coins := k.GetAllForeignCoins(ctx)
	for _, coin := range coins {
		zrc20Addr := ethcommon.HexToAddress(coin.Zrc20ContractAddress)
		//msg := types.NewMsgUpdateContractBytecode(creator, zrc20Addr, bytecode)
		//_, err = k.UpdateContractBytecode(ctx, msg)
		//if err != nil {
		//	panic(err)
		//}
		//_ = msg
		k.Logger(ctx).Info("UpdateContractBytecode", "zrc20Addr", zrc20Addr.String(), "coin", coin.Name)
		data, err := k.QueryZRC20Data(ctx, zrc20Addr)
		if err != nil {
			panic(err)
		}
		k.Logger(ctx).Info("QueryZRC20Data", "data", data)
		res, err := k.CallEVM(ctx, *zrc20ABI, types.ModuleAddressEVM, zrc20Addr, big.NewInt(0), big.NewInt(90000), false,
			"SYSTEM_CONTRACT_ADDRESS")
		if err != nil {
			k.Logger(ctx).Error("CallEVM SYSTEM_CONTRACT_ADDRESS()", "err", err, "revert code", res.Revert())
			panic(err)
		}
		{
			out, err := zrc20ABI.Unpack("SYSTEM_CONTRACT_ADDRESS", res.Ret)
			if err != nil {
				panic(err)
			}
			if len(out) != 1 {
				panic("output length is not 1")
			}
			out0 := *abi.ConvertType(out[0], new(ethcommon.Address)).(*ethcommon.Address)
			k.Logger(ctx).Info("CallEVM SYSTEM_CONTRACT_ADDRESS", "out0", out0.String())
		}

		res, err = k.CallEVM(ctx, *zrc20ABI, types.ModuleAddressEVM, zrc20Addr, big.NewInt(0), big.NewInt(90000), false,
			"withdrawGasFee")
		if err != nil {
			k.Logger(ctx).Error("CallEVM withdrawGasFee", "err", err, "revert code", res.Revert())
			panic(err)
		}
		{
			out, err := zrc20ABI.Unpack("withdrawGasFee", res.Ret)
			if err != nil {
				panic(err)
			}
			if len(out) != 2 {
				panic("output length is not 2")
			}
			out0 := *abi.ConvertType(out[0], new(ethcommon.Address)).(*ethcommon.Address)
			out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
			k.Logger(ctx).Info("CallEVM", "out0", out0.String(), "out1", out1.String())
		}
	}
	return nil
}
