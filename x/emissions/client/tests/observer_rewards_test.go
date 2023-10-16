package querytests

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	cli3 "github.com/zeta-chain/zetacore/x/emissions/client/cli"
	emmisonstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	observerCli "github.com/zeta-chain/zetacore/x/observer/client/cli"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (s *CliTestSuite) TestObserverRewards() {
	emissionPool := "800000000000000000000azeta"
	val := s.network.Validators[0]

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli3.CmdListPoolAddresses(), []string{"--output", "json"})
	s.Require().NoError(err)
	resPools := emmisonstypes.QueryListPoolAddressesResponse{}
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resPools))
	txArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(10))).String()),
	}

	// Fund the emission pool to start the emission process
	sendArgs := []string{val.Address.String(),
		resPools.EmissionModuleAddress, emissionPool}
	args := append(sendArgs, txArgs...)
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewSendTxCmd(), args)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// Collect parameter values and build assertion map for the randomised ballot set created
	emissionFactors := emmisonstypes.QueryGetEmmisonsFactorsResponse{}
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli3.CmdGetEmmisonsFactors(), []string{"--output", "json"})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &emissionFactors))
	emissionParams := emmisonstypes.QueryParamsResponse{}
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli3.CmdQueryParams(), []string{"--output", "json"})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &emissionParams))
	observerParams := observerTypes.QueryParamsResponse{}
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, observerCli.CmdQueryParams(), []string{"--output", "json"})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &observerParams))
	_, err = s.network.WaitForHeight(s.ballots[0].BallotCreationHeight + observerParams.Params.BallotMaturityBlocks)
	s.Require().NoError(err)
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli3.CmdGetEmmisonsFactors(), []string{"--output", "json"})
	resFactorsNewBlocks := emmisonstypes.QueryGetEmmisonsFactorsResponse{}
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resFactorsNewBlocks))
	// Duration factor is calculated in the same block,so we need to query based from the committed state at which the distribution is done
	// Would be cleaner to use `--height` flag, but it is not supported by the ExecTestCLICmd function yet
	emissionFactors.DurationFactor = resFactorsNewBlocks.DurationFactor
	asertValues := CalculateObserverRewards(s.ballots, emissionParams.Params.ObserverEmissionPercentage, emissionFactors.ReservesFactor, emissionFactors.BondFactor, emissionFactors.DurationFactor)

	// Assert withdrawable rewards for each validator
	resAvailable := emmisonstypes.QueryShowAvailableEmissionsResponse{}
	for i := 0; i < len(s.network.Validators); i++ {
		out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli3.CmdShowAvailableEmissions(), []string{s.network.Validators[i].Address.String(), "--output", "json"})
		s.Require().NoError(err)
		s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resAvailable))
		s.Require().Equal(sdk.NewCoin(config.BaseDenom, asertValues[s.network.Validators[i].Address.String()]).String(), resAvailable.Amount, "Validator %s has incorrect withdrawable rewards", s.network.Validators[i].Address.String())
	}

}

func CalculateObserverRewards(ballots []*observerTypes.Ballot, observerEmissionPercentage, reservesFactor, bondFactor, durationFactor string) map[string]sdkmath.Int {
	calculatedDistributer := map[string]sdkmath.Int{}
	blockRewards := sdk.MustNewDecFromStr(reservesFactor).Mul(sdk.MustNewDecFromStr(bondFactor)).Mul(sdk.MustNewDecFromStr(durationFactor))
	observerRewards := sdk.MustNewDecFromStr(observerEmissionPercentage).Mul(blockRewards).TruncateInt()
	rewardsDistributer := map[string]int64{}
	totalRewardsUnits := int64(0)
	// BuildRewardsDistribution has a separate unit test
	for _, ballot := range ballots {
		totalRewardsUnits = totalRewardsUnits + ballot.BuildRewardsDistribution(rewardsDistributer)
	}
	rewardPerUnit := observerRewards.Quo(sdk.NewInt(totalRewardsUnits))
	for address, units := range rewardsDistributer {
		if units == 0 {
			calculatedDistributer[address] = sdk.ZeroInt()
			continue
		}
		if units < 0 {
			calculatedDistributer[address] = sdk.ZeroInt()
			continue
		}
		if units > 0 {
			calculatedDistributer[address] = rewardPerUnit.Mul(sdkmath.NewInt(units))
		}
	}
	return calculatedDistributer
}
