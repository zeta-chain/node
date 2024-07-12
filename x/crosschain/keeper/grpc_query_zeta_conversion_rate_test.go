package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_ConvertGasToZeta(t *testing.T) {
	t.Run("should err if chain not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.ConvertGasToZeta(ctx, &types.QueryConvertGasToZetaRequest{
			ChainId: 987,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should err if median price not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.ConvertGasToZeta(ctx, &types.QueryConvertGasToZetaRequest{
			ChainId: 5,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should err if zrc20 not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
			Return(common.Address{}, errors.New("err"))

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     5,
			MedianIndex: 0,
			Prices:      []uint64{2},
		})

		res, err := k.ConvertGasToZeta(ctx, &types.QueryConvertGasToZetaRequest{
			ChainId:  5,
			GasLimit: "10",
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should err if uniswap2router not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
			Return(sample.EthAddress(), nil)

		fungibleMock.On("QueryUniswapV2RouterGetZetaAmountsIn", mock.Anything, mock.Anything, mock.Anything).
			Return(big.NewInt(0), errors.New("err"))

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     5,
			MedianIndex: 0,
			Prices:      []uint64{2},
		})

		res, err := k.ConvertGasToZeta(ctx, &types.QueryConvertGasToZetaRequest{
			ChainId:  5,
			GasLimit: "10",
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if all is set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
			Return(sample.EthAddress(), nil)

		fungibleMock.On("QueryUniswapV2RouterGetZetaAmountsIn", mock.Anything, mock.Anything, mock.Anything).
			Return(big.NewInt(5), nil)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     5,
			MedianIndex: 0,
			Prices:      []uint64{2},
		})

		res, err := k.ConvertGasToZeta(ctx, &types.QueryConvertGasToZetaRequest{
			ChainId:  5,
			GasLimit: "10",
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryConvertGasToZetaResponse{
			OutboundGasInZeta: "5",
			ProtocolFeeInZeta: types.GetProtocolFee().String(),
			// #nosec G115 always positive
			ZetaBlockHeight: uint64(ctx.BlockHeight()),
		}, res)
	})
}

func TestKeeper_ProtocolFee(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	res, err := k.ProtocolFee(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, &types.QueryMessagePassingProtocolFeeResponse{
		FeeInZeta: types.GetProtocolFee().String(),
	}, res)
}
