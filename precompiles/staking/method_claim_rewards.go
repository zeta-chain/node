package staking

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/precompiles/bank"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// claimRewards claims all the rewards for a delegator from a validator.
// As F1 Cosmos distribution scheme implements an all or nothing withdrawal, the precompile will
// withdraw all the rewards for the delegator, filter ZRC20 and unlock them to the delegator EVM address.
func (c *Contract) claimRewards(
	ctx sdk.Context,
	evm *vm.EVM,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, &precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		}
	}

	delegatorAddr, validatorAddr, err := unpackClaimRewardsArgs(args)
	if err != nil {
		return nil, err
	}

	// Get delegator Cosmos address.
	delegatorCosmosAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, delegatorAddr)
	if err != nil {
		return nil, err
	}

	// Get validator Cosmos address.
	validatorCosmosAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return nil, err
	}

	// Withdraw all the delegation rewards.
	// The F1 Cosmos distribution scheme implements an all or nothing withdrawal.
	// The coins could be of multiple denomination, and a mix of ZRC20 and Cosmos coins.
	coins, err := c.distributionKeeper.WithdrawDelegationRewards(ctx, delegatorCosmosAddr, validatorCosmosAddr)
	if err != nil {
		return nil, precompiletypes.ErrUnexpected{
			When: "WithdrawDelegationRewards",
			Got:  err.Error(),
		}
	}

	// For all the ZRC20 coins withdrawed:
	// - Check the amount to unlock is valid.
	// - Burn the Cosmos coins.
	// - Unlock the ZRC20 coins.
	for _, coin := range coins {
		// Filter out invalid coins.
		if !coin.IsValid() || !coin.Amount.IsPositive() || !precompiletypes.CoinIsZRC20(coin.Denom) {
			continue
		}

		// Notice that instead of returning errors we just skip the coin. This is because there might be
		// more than one ZRC20 coin in the delegation rewards, and we want to unlock as many as possible.
		// Coins are locked in the bank precompile, so it should be possible to unlock them afterwards.
		var (
			zrc20Addr   = common.HexToAddress(strings.TrimPrefix(coin.Denom, config.ZRC20DenomPrefix))
			zrc20Amount = coin.Amount.BigInt()
		)

		// Check if bank address has enough ZRC20 balance.
		// This check is also made inside UnlockZRC20, but repeat it here to avoid burning the coins.
		if err := c.fungibleKeeper.CheckZRC20Balance(ctx, zrc20Addr, bank.ContractAddress, zrc20Amount); err != nil {
			ctx.Logger().Error(
				"Claimed invalid amount of ZRC20 Validator Rewards",
				"Total", zrc20Amount,
				"Denom", precompiletypes.ZRC20ToCosmosDenom(zrc20Addr),
			)

			continue
		}

		coinSet := sdk.NewCoins(coin)

		// Send the coins to the fungible module to burn them.
		if err := c.bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorCosmosAddr, fungibletypes.ModuleName, coinSet); err != nil {
			continue
		}

		if err := c.bankKeeper.BurnCoins(ctx, fungibletypes.ModuleName, coinSet); err != nil {
			return nil, &precompiletypes.ErrUnexpected{
				When: "BurnCoins",
				Got:  err.Error(),
			}
		}

		// Finally, unlock the ZRC20 coins.
		if err := c.fungibleKeeper.UnlockZRC20(ctx, zrc20Addr, delegatorAddr, bank.ContractAddress, zrc20Amount); err != nil {
			return nil, &precompiletypes.ErrUnexpected{
				When: "UnlockZRC20",
				Got:  err.Error(),
			}
		}

		// Emit an event per ZRC20 coin unlocked.
		// This keeps events as granular and deterministic as possible.
		if err := c.addClaimRewardsLog(ctx, evm.StateDB, delegatorAddr, zrc20Addr, validatorCosmosAddr, zrc20Amount); err != nil {
			return nil, &precompiletypes.ErrUnexpected{
				When: "AddClaimRewardLog",
				Got:  err.Error(),
			}
		}

		ctx.Logger().Debug(
			"Claimed ZRC20 rewards",
			"Delegator", delegatorCosmosAddr,
			"Denom", precompiletypes.ZRC20ToCosmosDenom(zrc20Addr),
			"Amount", coin.Amount,
		)
	}

	return method.Outputs.Pack(true)
}

func unpackClaimRewardsArgs(args []interface{}) (delegator common.Address, validator string, err error) {
	delegator, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, "", &precompiletypes.ErrInvalidAddr{
			Got: delegator.String(),
		}
	}

	validator, ok = args[1].(string)
	if !ok {
		return common.Address{}, "", &precompiletypes.ErrInvalidAddr{
			Got: validator,
		}
	}

	return delegator, validator, nil
}
