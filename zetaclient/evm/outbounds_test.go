package evm_test

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/evm"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
)

// getContractsByChainID is a helper func to get contracts and addresses by chainID
func getContractsByChainID(chainID int64) (*zetaconnector.ZetaConnectorNonEth, ethcommon.Address, *erc20custody.ERC20Custody, ethcommon.Address) {
	connector := stub.MockConnectorNonEth(chainID)
	connectorAddress := testutils.ConnectorAddresses[chainID]
	custody := stub.MockERC20Custody(chainID)
	custodyAddress := testutils.CustodyAddresses[chainID]
	return connector, connectorAddress, custody, custodyAddress
}

func Test_IsOutboundProcessed(t *testing.T) {
	// load archived outtx receipt that contains ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chain := chains.EthChain
	chainID := chains.EthChain.ChainId
	nonce := uint64(9718)
	chainParam := stub.MockChainParams(chain.ChainId, 1)
	outtxHash := "0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f"
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(t, chainID, outtxHash, coin.CoinType_Zeta, testutils.EventZetaReceived)
	cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, testutils.EventZetaReceived)
	t.Run("should post vote and return true if outtx is processed", func(t *testing.T) {
		// create evm client and set outtx and receipt
		client := MockEVMClient(t, chain, nil, nil, nil, nil, 1, chainParam)
		client.SetTxNReceipt(nonce, receipt, outtx)
		// post outbound vote
		isIncluded, isConfirmed, err := client.IsOutboundProcessed(cctx, zerolog.Logger{})
		require.NoError(t, err)
		require.True(t, isIncluded)
		require.True(t, isConfirmed)
	})
	t.Run("should post vote and return true on restricted address", func(t *testing.T) {
		// load cctx and modify sender address to arbitrary address
		// Note: other tests cases will fail if we use the original sender address because the
		// compliance config is globally set and will impact other tests when running in parallel
		cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		cctx.InboundParams.Sender = sample.EthAddress().Hex()

		// create evm client and set outtx and receipt
		client := MockEVMClient(t, chain, nil, nil, nil, nil, 1, chainParam)
		client.SetTxNReceipt(nonce, receipt, outtx)

		// modify compliance config to restrict sender address
		cfg := config.Config{
			ComplianceConfig: config.ComplianceConfig{},
		}
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.InboundParams.Sender}
		config.LoadComplianceConfig(cfg)

		// post outbound vote
		isIncluded, isConfirmed, err := client.IsOutboundProcessed(cctx, zerolog.Logger{})
		require.NoError(t, err)
		require.True(t, isIncluded)
		require.True(t, isConfirmed)
	})
	t.Run("should return false if outtx is not confirmed", func(t *testing.T) {
		// create evm client and DO NOT set outtx as confirmed
		client := MockEVMClient(t, chain, nil, nil, nil, nil, 1, chainParam)
		isIncluded, isConfirmed, err := client.IsOutboundProcessed(cctx, zerolog.Logger{})
		require.NoError(t, err)
		require.False(t, isIncluded)
		require.False(t, isConfirmed)
	})
	t.Run("should fail if unable to parse ZetaReceived event", func(t *testing.T) {
		// create evm client and set outtx and receipt
		client := MockEVMClient(t, chain, nil, nil, nil, nil, 1, chainParam)
		client.SetTxNReceipt(nonce, receipt, outtx)

		// set connector contract address to an arbitrary address to make event parsing fail
		chainParamsNew := client.GetChainParams()
		chainParamsNew.ConnectorContractAddress = sample.EthAddress().Hex()
		client.SetChainParams(chainParamsNew)
		isIncluded, isConfirmed, err := client.IsOutboundProcessed(cctx, zerolog.Logger{})
		require.Error(t, err)
		require.False(t, isIncluded)
		require.False(t, isConfirmed)
	})
}

func Test_IsOutboundProcessed_ContractError(t *testing.T) {
	// Note: this test is skipped because it will cause CI failure.
	// The only way to replicate a contract error is to use an invalid ABI.
	// See the code: https://github.com/ethereum/go-ethereum/blob/v1.10.26/accounts/abi/bind/base.go#L97
	// The ABI is hardcoded in the protocol-contracts package and initialized the 1st time it binds the contract.
	// Any subsequent modification to the ABI will not work and therefor fail the unit test.
	t.Skip("uncomment this line to run this test separately, otherwise it will fail CI")

	// load archived outtx receipt that contains ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chain := chains.EthChain
	chainID := chains.EthChain.ChainId
	nonce := uint64(9718)
	chainParam := stub.MockChainParams(chain.ChainId, 1)
	outtxHash := "0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f"
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(t, chainID, outtxHash, coin.CoinType_Zeta, testutils.EventZetaReceived)
	cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, testutils.EventZetaReceived)

	t.Run("should fail if unable to get connector/custody contract", func(t *testing.T) {
		// create evm client and set outtx and receipt
		client := MockEVMClient(t, chain, nil, nil, nil, nil, 1, chainParam)
		client.SetTxNReceipt(nonce, receipt, outtx)
		abiConnector := zetaconnector.ZetaConnectorNonEthMetaData.ABI
		abiCustody := erc20custody.ERC20CustodyMetaData.ABI

		// set invalid connector ABI
		zetaconnector.ZetaConnectorNonEthMetaData.ABI = "invalid abi"
		isIncluded, isConfirmed, err := client.IsOutboundProcessed(cctx, zerolog.Logger{})
		zetaconnector.ZetaConnectorNonEthMetaData.ABI = abiConnector // reset connector ABI
		require.ErrorContains(t, err, "error getting zeta connector")
		require.False(t, isIncluded)
		require.False(t, isConfirmed)

		// set invalid custody ABI
		erc20custody.ERC20CustodyMetaData.ABI = "invalid abi"
		isIncluded, isConfirmed, err = client.IsOutboundProcessed(cctx, zerolog.Logger{})
		require.ErrorContains(t, err, "error getting erc20 custody")
		require.False(t, isIncluded)
		require.False(t, isConfirmed)
		erc20custody.ERC20CustodyMetaData.ABI = abiCustody // reset custody ABI
	})
}

func Test_PostVoteOutbound(t *testing.T) {
	// Note: outtx of Gas/ERC20 token can also be used for this test
	// load archived cctx, outtx and receipt for a ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chain := chains.EthChain
	nonce := uint64(9718)
	coinType := coin.CoinType_Zeta
	cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chain.ChainId, nonce, testutils.EventZetaReceived)

	t.Run("post vote outbound successfully", func(t *testing.T) {
		// the amount and status to be used for vote
		receiveValue := cctx.GetCurrentOutboundParam().Amount.BigInt()
		receiveStatus := chains.ReceiveStatus_success

		// create evm client using mock zetaBridge and post outbound vote
		zetaBridge := stub.NewMockZetaCoreBridge()
		client := MockEVMClient(t, chain, nil, nil, zetaBridge, nil, 1, observertypes.ChainParams{})
		client.PostVoteOutbound(cctx.Index, receipt, outtx, receiveValue, receiveStatus, nonce, coinType, zerolog.Logger{})

		// pause the mock zetaBridge to simulate error posting vote
		zetaBridge.Pause()
		client.PostVoteOutbound(cctx.Index, receipt, outtx, receiveValue, receiveStatus, nonce, coinType, zerolog.Logger{})
	})
}

func Test_ParseZetaReceived(t *testing.T) {
	// load archived outtx receipt that contains ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chainID := chains.EthChain.ChainId
	nonce := uint64(9718)
	outtxHash := "0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f"
	connector := stub.MockConnectorNonEth(chainID)
	connectorAddress := testutils.ConnectorAddresses[chainID]
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(t, chainID, outtxHash, coin.CoinType_Zeta, testutils.EventZetaReceived)

	t.Run("should parse ZetaReceived event from archived outtx receipt", func(t *testing.T) {
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.NoError(t, err)
		require.NotNil(t, receivedLog)
		require.Nil(t, revertedLog)
	})
	t.Run("should fail on connector address mismatch", func(t *testing.T) {
		// use an arbitrary address to make validation fail
		fakeConnectorAddress := sample.EthAddress()
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(cctx, receipt, fakeConnectorAddress, connector)
		require.ErrorContains(t, err, "error validating ZetaReceived event")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on receiver address mismatch", func(t *testing.T) {
		// load cctx and set receiver address to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.GetCurrentOutboundParam().Receiver = sample.EthAddress().Hex()
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "receiver address mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on amount mismatch", func(t *testing.T) {
		// load cctx and set amount to an arbitrary wrong value
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeAmount := sample.UintInRange(0, fakeCctx.GetCurrentOutboundParam().Amount.Uint64()-1)
		fakeCctx.GetCurrentOutboundParam().Amount = fakeAmount
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "amount mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on cctx index mismatch", func(t *testing.T) {
		cctx.Index = sample.Hash().Hex() // use an arbitrary index
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "cctx index mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail if no event found in receipt", func(t *testing.T) {
		// load receipt and remove ZetaReceived event from logs
		receipt := testutils.LoadEVMOutboundReceipt(t, chainID, outtxHash, coin.CoinType_Zeta, testutils.EventZetaReceived)
		receipt.Logs = receipt.Logs[:1] // the 2nd log is ZetaReceived event
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "no ZetaReceived/ZetaReverted event")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
}

func Test_ParseZetaReverted(t *testing.T) {
	// load archived outtx receipt that contains ZetaReverted event
	chainID := chains.GoerliLocalnetChain.ChainId
	nonce := uint64(14)
	outtxHash := "0x1487e6a31dd430306667250b72bf15b8390b73108b69f3de5c1b2efe456036a7"
	connector := stub.MockConnectorNonEth(chainID)
	connectorAddress := testutils.ConnectorAddresses[chainID]
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(t, chainID, outtxHash, coin.CoinType_Zeta, testutils.EventZetaReverted)

	t.Run("should parse ZetaReverted event from archived outtx receipt", func(t *testing.T) {
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.NoError(t, err)
		require.Nil(t, receivedLog)
		require.NotNil(t, revertedLog)
	})
	t.Run("should fail on connector address mismatch", func(t *testing.T) {
		// use an arbitrary address to make validation fail
		fakeConnectorAddress := sample.EthAddress()
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(cctx, receipt, fakeConnectorAddress, connector)
		require.ErrorContains(t, err, "error validating ZetaReverted event")
		require.Nil(t, receivedLog)
		require.Nil(t, revertedLog)
	})
	t.Run("should fail on receiver address mismatch", func(t *testing.T) {
		// load cctx and set receiver address to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.InboundParams.Sender = sample.EthAddress().Hex() // the receiver is the sender for reverted ccxt
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "receiver address mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on amount mismatch", func(t *testing.T) {
		// load cctx and set amount to an arbitrary wrong value
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeAmount := sample.UintInRange(0, fakeCctx.GetCurrentOutboundParam().Amount.Uint64()-1)
		fakeCctx.GetCurrentOutboundParam().Amount = fakeAmount
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "amount mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on cctx index mismatch", func(t *testing.T) {
		cctx.Index = sample.Hash().Hex() // use an arbitrary index to make validation fail
		receivedLog, revertedLog, err := evm.ParseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "cctx index mismatch")
		require.Nil(t, receivedLog)
		require.Nil(t, revertedLog)
	})
}

func Test_ParseERC20WithdrawnEvent(t *testing.T) {
	// load archived outtx receipt that contains ERC20 Withdrawn event
	chainID := chains.EthChain.ChainId
	nonce := uint64(8014)
	outtxHash := "0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146"
	custody := stub.MockERC20Custody(chainID)
	custodyAddress := testutils.CustodyAddresses[chainID]
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(t, chainID, outtxHash, coin.CoinType_ERC20, testutils.EventERC20Withdraw)

	t.Run("should parse ERC20 Withdrawn event from archived outtx receipt", func(t *testing.T) {
		withdrawn, err := evm.ParseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
		require.NoError(t, err)
		require.NotNil(t, withdrawn)
	})
	t.Run("should fail on erc20 custody address mismatch", func(t *testing.T) {
		// use an arbitrary address to make validation fail
		fakeCustodyAddress := sample.EthAddress()
		withdrawn, err := evm.ParseAndCheckWithdrawnEvent(cctx, receipt, fakeCustodyAddress, custody)
		require.ErrorContains(t, err, "error validating Withdrawn event")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail on receiver address mismatch", func(t *testing.T) {
		// load cctx and set receiver address to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.GetCurrentOutboundParam().Receiver = sample.EthAddress().Hex()
		withdrawn, err := evm.ParseAndCheckWithdrawnEvent(fakeCctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "receiver address mismatch")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail on asset mismatch", func(t *testing.T) {
		// load cctx and set asset to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.InboundParams.Asset = sample.EthAddress().Hex()
		withdrawn, err := evm.ParseAndCheckWithdrawnEvent(fakeCctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "asset mismatch")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail on amount mismatch", func(t *testing.T) {
		// load cctx and set amount to an arbitrary wrong value
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeAmount := sample.UintInRange(0, fakeCctx.GetCurrentOutboundParam().Amount.Uint64()-1)
		fakeCctx.GetCurrentOutboundParam().Amount = fakeAmount
		withdrawn, err := evm.ParseAndCheckWithdrawnEvent(fakeCctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "amount mismatch")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail if no Withdrawn event found in receipt", func(t *testing.T) {
		// load receipt and remove Withdrawn event from logs
		receipt := testutils.LoadEVMOutboundReceipt(t, chainID, outtxHash, coin.CoinType_ERC20, testutils.EventERC20Withdraw)
		receipt.Logs = receipt.Logs[:1] // the 2nd log is Withdrawn event
		withdrawn, err := evm.ParseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "no ERC20 Withdrawn event")
		require.Nil(t, withdrawn)
	})
}

func Test_ParseOutboundReceivedValue(t *testing.T) {
	chainID := chains.EthChain.ChainId
	connector, connectorAddr, custody, custodyAddr := getContractsByChainID(chainID)

	t.Run("should parse and check ZetaReceived event from archived outtx receipt", func(t *testing.T) {
		// load archived outtx receipt that contains ZetaReceived event
		// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
		nonce := uint64(9718)
		coinType := coin.CoinType_Zeta
		cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, testutils.EventZetaReceived)
		params := cctx.GetCurrentOutboundParam()
		value, status, err := evm.ParseOutboundReceivedValue(cctx, receipt, outtx, coinType, connectorAddr, connector, custodyAddr, custody)
		require.NoError(t, err)
		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
		require.Equal(t, chains.ReceiveStatus_success, status)
	})
	t.Run("should parse and check ZetaReverted event from archived outtx receipt", func(t *testing.T) {
		// load archived outtx receipt that contains ZetaReverted event
		// use local network tx: 0x1487e6a31dd430306667250b72bf15b8390b73108b69f3de5c1b2efe456036a7
		localChainID := chains.GoerliLocalnetChain.ChainId
		nonce := uint64(14)
		coinType := coin.CoinType_Zeta
		connectorLocal, connectorAddrLocal, custodyLocal, custodyAddrLocal := getContractsByChainID(localChainID)
		cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, localChainID, nonce, testutils.EventZetaReverted)
		params := cctx.GetCurrentOutboundParam()
		value, status, err := evm.ParseOutboundReceivedValue(
			cctx, receipt, outtx, coinType, connectorAddrLocal, connectorLocal, custodyAddrLocal, custodyLocal)
		require.NoError(t, err)
		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
		require.Equal(t, chains.ReceiveStatus_success, status)
	})
	t.Run("should parse and check ERC20 Withdrawn event from archived outtx receipt", func(t *testing.T) {
		// load archived outtx receipt that contains ERC20 Withdrawn event
		// https://etherscan.io/tx/0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146
		nonce := uint64(8014)
		coinType := coin.CoinType_ERC20
		cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, testutils.EventERC20Withdraw)
		params := cctx.GetCurrentOutboundParam()
		value, status, err := evm.ParseOutboundReceivedValue(cctx, receipt, outtx, coinType, connectorAddr, connector, custodyAddr, custody)
		require.NoError(t, err)
		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
		require.Equal(t, chains.ReceiveStatus_success, status)
	})
	t.Run("nothing to parse if coinType is Gas", func(t *testing.T) {
		// load archived outtx receipt of Gas token transfer
		// https://etherscan.io/tx/0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3
		nonce := uint64(7260)
		coinType := coin.CoinType_Gas
		cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, "")
		params := cctx.GetCurrentOutboundParam()
		value, status, err := evm.ParseOutboundReceivedValue(cctx, receipt, outtx, coinType, connectorAddr, connector, custodyAddr, custody)
		require.NoError(t, err)
		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
		require.Equal(t, chains.ReceiveStatus_success, status)
	})
	t.Run("should fail on unknown coin type", func(t *testing.T) {
		// load archived outtx receipt that contains ZetaReceived event
		// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
		nonce := uint64(9718)
		coinType := coin.CoinType(5) // unknown coin type
		cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, testutils.EventZetaReceived)
		value, status, err := evm.ParseOutboundReceivedValue(cctx, receipt, outtx, coinType, connectorAddr, connector, custodyAddr, custody)
		require.ErrorContains(t, err, "unknown coin type")
		require.Nil(t, value)
		require.Equal(t, chains.ReceiveStatus_failed, status)
	})
	t.Run("should fail if unable to parse ZetaReceived event", func(t *testing.T) {
		// load archived outtx receipt that contains ZetaReceived event
		// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
		nonce := uint64(9718)
		coinType := coin.CoinType_Zeta
		cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, testutils.EventZetaReceived)

		// use an arbitrary address to make event parsing fail
		fakeConnectorAddress := sample.EthAddress()
		value, status, err := evm.ParseOutboundReceivedValue(cctx, receipt, outtx, coinType, fakeConnectorAddress, connector, custodyAddr, custody)
		require.Error(t, err)
		require.Nil(t, value)
		require.Equal(t, chains.ReceiveStatus_failed, status)
	})
	t.Run("should fail if unable to parse ERC20 Withdrawn event", func(t *testing.T) {
		// load archived outtx receipt that contains ERC20 Withdrawn event
		// https://etherscan.io/tx/0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146
		nonce := uint64(8014)
		coinType := coin.CoinType_ERC20
		cctx, outtx, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, chainID, nonce, testutils.EventERC20Withdraw)

		// use an arbitrary address to make event parsing fail
		fakeCustodyAddress := sample.EthAddress()
		value, status, err := evm.ParseOutboundReceivedValue(cctx, receipt, outtx, coinType, connectorAddr, connector, fakeCustodyAddress, custody)
		require.Error(t, err)
		require.Nil(t, value)
		require.Equal(t, chains.ReceiveStatus_failed, status)
	})
}
