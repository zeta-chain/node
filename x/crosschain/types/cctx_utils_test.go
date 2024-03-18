package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestValidateAddressForChain(t *testing.T) {
	require.Error(t, types.ValidateAddressForChain("0x123", common.GoerliChain().ChainId))
	require.Error(t, types.ValidateAddressForChain("", common.GoerliChain().ChainId))
	require.Error(t, types.ValidateAddressForChain("%%%%", common.GoerliChain().ChainId))
	require.NoError(t, types.ValidateAddressForChain("0x792c127Fa3AC1D52F904056Baf1D9257391e7D78", common.GoerliChain().ChainId))
	require.Error(t, types.ValidateAddressForChain("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", common.BtcMainnetChain().ChainId))
	require.Error(t, types.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", common.BtcMainnetChain().ChainId))
	require.Error(t, types.ValidateAddressForChain("", common.BtcRegtestChain().ChainId))
	require.NoError(t, types.ValidateAddressForChain("bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", common.BtcMainnetChain().ChainId))
	require.NoError(t, types.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", common.BtcRegtestChain().ChainId))
	require.NoError(t, types.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", common.ZetaChainMainnet().ChainId))
	require.NoError(t, types.ValidateAddressForChain("0x792c127Fa3AC1D52F904056Baf1D9257391e7D78", common.ZetaChainMainnet().ChainId))
}

func TestValidateZetaIndex(t *testing.T) {
	require.NoError(t, types.ValidateZetaIndex("0x84bd5c9922b63c52d8a9ca686e0a57ff978150b71be0583514d01c27aa341910"))
	require.NoError(t, types.ValidateZetaIndex(sample.ZetaIndex(t)))
	require.Error(t, types.ValidateZetaIndex("0"))
	require.Error(t, types.ValidateZetaIndex("0x70e967acFcC17c3941E87562161406d41676FD83"))
}

func TestValidateHashForChain(t *testing.T) {
	require.NoError(t, types.ValidateHashForChain("0x84bd5c9922b63c52d8a9ca686e0a57ff978150b71be0583514d01c27aa341910", common.GoerliChain().ChainId))
	require.Error(t, types.ValidateHashForChain("", common.GoerliChain().ChainId))
	require.Error(t, types.ValidateHashForChain("a0fa5a82f106fb192e4c503bfa8d54b2de20a821e09338094ab825cc9b275059", common.GoerliChain().ChainId))
	require.NoError(t, types.ValidateHashForChain("15b7880f5d236e857a5e8f043ce9d56f5ef01e1c3f2a786baf740fc0bb7a22a3", common.BtcMainnetChain().ChainId))
	require.NoError(t, types.ValidateHashForChain("a0fa5a82f106fb192e4c503bfa8d54b2de20a821e09338094ab825cc9b275059", common.BtcTestNetChain().ChainId))
	require.Error(t, types.ValidateHashForChain("0x84bd5c9922b63c52d8a9ca686e0a57ff978150b71be0583514d01c27aa341910", common.BtcMainnetChain().ChainId))
}

func TestInboundTxParams_Validate(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	inTxParams := sample.InboundTxParamsValidChainID(r)
	inTxParams.Sender = ""
	require.ErrorContains(t, inTxParams.Validate(), "sender cannot be empty")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.SenderChainId = 1000
	require.ErrorContains(t, inTxParams.Validate(), "invalid sender chain id 1000")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.SenderChainId = common.GoerliChain().ChainId
	inTxParams.Sender = "0x123"
	require.ErrorContains(t, inTxParams.Validate(), "invalid address 0x123")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.SenderChainId = common.GoerliChain().ChainId
	inTxParams.TxOrigin = "0x123"
	require.ErrorContains(t, inTxParams.Validate(), "invalid address 0x123")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.Amount = sdkmath.Uint{}
	require.ErrorContains(t, inTxParams.Validate(), "amount cannot be nil")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.InboundTxObservedHash = "12"
	require.ErrorContains(t, inTxParams.Validate(), "hash must be a valid ethereum hash 12")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.InboundTxObservedHash = sample.Hash().String()
	inTxParams.InboundTxBallotIndex = "12"
	require.ErrorContains(t, inTxParams.Validate(), "invalid index hash 12")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.InboundTxObservedHash = sample.Hash().String()
	inTxParams.InboundTxBallotIndex = sample.ZetaIndex(t)
	require.NoError(t, inTxParams.Validate())
}

func TestOutboundTxParams_Validate(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	outTxParams := sample.OutboundTxParamsValidChainID(r)
	outTxParams.Receiver = ""
	require.ErrorContains(t, outTxParams.Validate(), "receiver cannot be empty")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.ReceiverChainId = 1000
	require.ErrorContains(t, outTxParams.Validate(), "invalid receiver chain id 1000")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.Receiver = "0x123"
	require.ErrorContains(t, outTxParams.Validate(), "invalid address 0x123")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.Amount = sdkmath.Uint{}
	require.ErrorContains(t, outTxParams.Validate(), "amount cannot be nil")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.OutboundTxBallotIndex = "12"
	require.ErrorContains(t, outTxParams.Validate(), "invalid index hash 12")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.OutboundTxBallotIndex = sample.ZetaIndex(t)
	outTxParams.OutboundTxHash = sample.Hash().String()
	require.NoError(t, outTxParams.Validate())
}

func TestCrossChainTx_Validate(t *testing.T) {
	cctx := sample.CrossChainTx(t, "foo")
	cctx.InboundTxParams = nil
	require.ErrorContains(t, cctx.Validate(), "inbound tx params cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundTxParams = nil
	require.ErrorContains(t, cctx.Validate(), "outbound tx params cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.CctxStatus = nil
	require.ErrorContains(t, cctx.Validate(), "cctx status cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundTxParams = make([]*types.OutboundTxParams, 3)
	require.ErrorContains(t, cctx.Validate(), "outbound tx params cannot be more than 2")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.Index = "0"
	require.ErrorContains(t, cctx.Validate(), "invalid index hash 0")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.InboundTxParams = sample.InboundTxParamsValidChainID(rand.New(rand.NewSource(42)))
	cctx.InboundTxParams.SenderChainId = 1000
	require.ErrorContains(t, cctx.Validate(), "invalid sender chain id 1000")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParamsValidChainID(rand.New(rand.NewSource(42)))}
	cctx.InboundTxParams = sample.InboundTxParamsValidChainID(rand.New(rand.NewSource(42)))
	cctx.InboundTxParams.InboundTxObservedHash = sample.Hash().String()
	cctx.InboundTxParams.InboundTxBallotIndex = sample.ZetaIndex(t)
	cctx.OutboundTxParams[0].ReceiverChainId = 1000
	require.ErrorContains(t, cctx.Validate(), "invalid receiver chain id 1000")
}

func TestCrossChainTx_GetCurrentOutTxParam(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundTxParams = []*types.OutboundTxParams{}
	require.Equal(t, &types.OutboundTxParams{}, cctx.GetCurrentOutTxParam())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[0], cctx.GetCurrentOutTxParam())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r), sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[1], cctx.GetCurrentOutTxParam())
}

func TestCrossChainTx_IsCurrentOutTxRevert(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundTxParams = []*types.OutboundTxParams{}
	require.False(t, cctx.IsCurrentOutTxRevert())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r)}
	require.False(t, cctx.IsCurrentOutTxRevert())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r), sample.OutboundTxParams(r)}
	require.True(t, cctx.IsCurrentOutTxRevert())
}

func TestCrossChainTx_OriginalDestinationChainID(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundTxParams = []*types.OutboundTxParams{}
	require.Equal(t, int64(-1), cctx.OriginalDestinationChainID())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[0].ReceiverChainId, cctx.OriginalDestinationChainID())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r), sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[0].ReceiverChainId, cctx.OriginalDestinationChainID())
}

func TestOutboundTxParams_GetGasPrice(t *testing.T) {
	// #nosec G404 - random seed is not used for security purposes
	r := rand.New(rand.NewSource(42))
	outTxParams := sample.OutboundTxParams(r)

	outTxParams.OutboundTxGasPrice = "42"
	gasPrice, err := outTxParams.GetGasPrice()
	require.NoError(t, err)
	require.EqualValues(t, uint64(42), gasPrice)

	outTxParams.OutboundTxGasPrice = "invalid"
	_, err = outTxParams.GetGasPrice()
	require.Error(t, err)
}
