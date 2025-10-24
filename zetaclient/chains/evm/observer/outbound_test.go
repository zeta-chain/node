package observer

import (
	"context"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"
)

func Test_IsOutboundProcessed(t *testing.T) {
	// load archived outbound receipt that contains ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chainID := chains.Ethereum.ChainId
	nonce := uint64(9718)
	outboundHash := "0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f"
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(
		t,
		TestDataDir,
		chainID,
		outboundHash,
		coin.CoinType_Zeta,
		testutils.EventZetaReceived,
	)
	cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
		t,
		TestDataDir,
		chainID,
		nonce,
		testutils.EventZetaReceived,
	)

	ctx := context.Background()

	t.Run("should post vote and return true if outbound is processed", func(t *testing.T) {
		// create evm observer and set outbound and receipt
		ob := newTestSuite(t)
		ob.setTxNReceipt(nonce, receipt, outbound)

		// post outbound vote
		continueKeysign, err := ob.VoteOutboundIfConfirmed(ctx, cctx)
		require.NoError(t, err)
		require.False(t, continueKeysign)
	})
	t.Run("should post vote and return true on restricted address", func(t *testing.T) {
		// load cctx and modify sender address to arbitrary address
		// Note: other tests cases will fail if we use the original sender address because the
		// compliance config is globally set and will impact other tests when running in parallel
		cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		cctx.InboundParams.Sender = sample.EthAddress().Hex()

		// create evm observer and set outbound and receipt
		ob := newTestSuite(t)
		ob.setTxNReceipt(nonce, receipt, outbound)

		// modify compliance config to restrict sender address
		cfg := config.Config{
			ComplianceConfig: config.ComplianceConfig{},
		}
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.InboundParams.Sender}
		config.SetRestrictedAddressesFromConfig(cfg)

		// post outbound vote
		continueKeysign, err := ob.VoteOutboundIfConfirmed(ctx, cctx)
		require.NoError(t, err)
		require.False(t, continueKeysign)
	})
	t.Run("should return false if outbound is not confirmed", func(t *testing.T) {
		// create evm observer and DO NOT set outbound as confirmed
		ob := newTestSuite(t)
		continueKeysign, err := ob.VoteOutboundIfConfirmed(ctx, cctx)
		require.NoError(t, err)
		require.True(t, continueKeysign)
	})
	t.Run("should fail if unable to parse ZetaReceived event", func(t *testing.T) {
		// create evm observer and set outbound and receipt
		ob := newTestSuite(t)
		ob.setTxNReceipt(nonce, receipt, outbound)

		// set connector contract address to an arbitrary address to make event parsing fail
		chainParamsNew := ob.ChainParams()
		chainParamsNew.ConnectorContractAddress = sample.EthAddress().Hex()
		ob.SetChainParams(chainParamsNew)
		continueKeysign, err := ob.VoteOutboundIfConfirmed(ctx, cctx)
		require.Error(t, err)
		require.True(t, continueKeysign)
	})
}

func Test_IsOutboundProcessed_ContractError(t *testing.T) {
	// Note: this test is skipped because it will cause CI failure.
	// The only way to replicate a contract error is to use an invalid ABI.
	// See the code: https://github.com/ethereum/go-ethereum/blob/v1.10.26/accounts/abi/bind/base.go#L97
	// The ABI is hardcoded in the protocol-contracts package and initialized the 1st time it binds the contract.
	// Any subsequent modification to the ABI will not work and therefor fail the unit test.
	t.Skip("uncomment this line to run this test separately, otherwise it will fail CI")

	// load archived outbound receipt that contains ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chainID := chains.Ethereum.ChainId
	nonce := uint64(9718)
	outboundHash := "0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f"
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(
		t,
		TestDataDir,
		chainID,
		outboundHash,
		coin.CoinType_Zeta,
		testutils.EventZetaReceived,
	)
	cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
		t,
		TestDataDir,
		chainID,
		nonce,
		testutils.EventZetaReceived,
	)

	ctx := context.Background()

	t.Run("should fail if unable to get connector/custody contract", func(t *testing.T) {
		// create evm observer and set outbound and receipt
		ob := newTestSuite(t)
		ob.setTxNReceipt(nonce, receipt, outbound)
		abiConnector := zetaconnector.ZetaConnectorNonEthMetaData.ABI
		abiCustody := erc20custody.ERC20CustodyMetaData.ABI

		// set invalid connector ABI
		zetaconnector.ZetaConnectorNonEthMetaData.ABI = "invalid abi"
		continueKeysign, err := ob.VoteOutboundIfConfirmed(ctx, cctx)
		zetaconnector.ZetaConnectorNonEthMetaData.ABI = abiConnector // reset connector ABI
		require.ErrorContains(t, err, "error getting zeta connector")
		require.True(t, continueKeysign)

		// set invalid custody ABI
		erc20custody.ERC20CustodyMetaData.ABI = "invalid abi"
		continueKeysign, err = ob.VoteOutboundIfConfirmed(ctx, cctx)
		require.ErrorContains(t, err, "error getting erc20 custody")
		require.True(t, continueKeysign)
		erc20custody.ERC20CustodyMetaData.ABI = abiCustody // reset custody ABI
	})
}

func Test_PostVoteOutbound(t *testing.T) {
	// Note: outbound of Gas/ERC20 token can also be used for this test
	// load archived cctx, outbound and receipt for a ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chain := chains.Ethereum
	nonce := uint64(9718)
	coinType := coin.CoinType_Zeta
	cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
		t,
		TestDataDir,
		chain.ChainId,
		nonce,
		testutils.EventZetaReceived,
	)

	ctx := context.Background()

	t.Run("post vote outbound successfully", func(t *testing.T) {
		// the amount and status to be used for vote
		receiveValue := cctx.GetCurrentOutboundParam().Amount.BigInt()
		receiveStatus := chains.ReceiveStatus_success

		// create evm client using mock zetacore client and post outbound vote
		ob := newTestSuite(t)
		ob.postVoteOutbound(
			ctx,
			cctx.Index,
			receipt,
			outbound,
			receiveValue,
			receiveStatus,
			nonce,
			coinType,
			zerolog.Nop(),
		)
	})
}

func Test_ParseZetaReceived(t *testing.T) {
	// load archived outbound receipt that contains ZetaReceived event
	// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
	chainID := chains.Ethereum.ChainId
	nonce := uint64(9718)
	outboundHash := "0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f"
	connector := mocks.MockConnectorNonEth(t, chainID)
	connectorAddress := testutils.ConnectorAddresses[chainID]
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(
		t,
		TestDataDir,
		chainID,
		outboundHash,
		coin.CoinType_Zeta,
		testutils.EventZetaReceived,
	)

	t.Run("should parse ZetaReceived event from archived outbound receipt", func(t *testing.T) {
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.NoError(t, err)
		require.NotNil(t, receivedLog)
		require.Nil(t, revertedLog)
	})
	t.Run("should fail on connector address mismatch", func(t *testing.T) {
		// use an arbitrary address to make validation fail
		fakeConnectorAddress := sample.EthAddress()
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, fakeConnectorAddress, connector)
		require.ErrorContains(t, err, "error validating ZetaReceived event")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on receiver address mismatch", func(t *testing.T) {
		// load cctx and set receiver address to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.GetCurrentOutboundParam().Receiver = sample.EthAddress().Hex()
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "receiver address mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on amount mismatch", func(t *testing.T) {
		// load cctx and set amount to an arbitrary wrong value
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeAmount := sample.UintInRange(0, fakeCctx.GetCurrentOutboundParam().Amount.Uint64()-1)
		fakeCctx.GetCurrentOutboundParam().Amount = fakeAmount
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "amount mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on cctx index mismatch", func(t *testing.T) {
		cctx.Index = sample.Hash().Hex() // use an arbitrary index
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "cctx index mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail if no event found in receipt", func(t *testing.T) {
		// load receipt and remove ZetaReceived event from logs
		receipt := testutils.LoadEVMOutboundReceipt(
			t,
			TestDataDir,
			chainID,
			outboundHash,
			coin.CoinType_Zeta,
			testutils.EventZetaReceived,
		)
		receipt.Logs = receipt.Logs[:1] // the 2nd log is ZetaReceived event
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "no ZetaReceived/ZetaReverted event")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
}

func Test_ParseZetaReverted(t *testing.T) {
	// load archived outbound receipt that contains ZetaReverted event
	chainID := chains.GoerliLocalnet.ChainId
	nonce := uint64(14)
	outboundHash := "0x1487e6a31dd430306667250b72bf15b8390b73108b69f3de5c1b2efe456036a7"
	connector := mocks.MockConnectorNonEth(t, chainID)
	connectorAddress := testutils.ConnectorAddresses[chainID]
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(
		t,
		TestDataDir,
		chainID,
		outboundHash,
		coin.CoinType_Zeta,
		testutils.EventZetaReverted,
	)

	t.Run("should parse ZetaReverted event from archived outbound receipt", func(t *testing.T) {
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.NoError(t, err)
		require.Nil(t, receivedLog)
		require.NotNil(t, revertedLog)
	})
	t.Run("should fail on connector address mismatch", func(t *testing.T) {
		// use an arbitrary address to make validation fail
		fakeConnectorAddress := sample.EthAddress()
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, fakeConnectorAddress, connector)
		require.ErrorContains(t, err, "error validating ZetaReverted event")
		require.Nil(t, receivedLog)
		require.Nil(t, revertedLog)
	})
	t.Run("should fail on receiver address mismatch", func(t *testing.T) {
		// load cctx and set receiver address to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.InboundParams.Sender = sample.EthAddress().Hex() // the receiver is the sender for reverted ccxt
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "receiver address mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on amount mismatch", func(t *testing.T) {
		// load cctx and set amount to an arbitrary wrong value
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeAmount := sample.UintInRange(0, fakeCctx.GetCurrentOutboundParam().Amount.Uint64()-1)
		fakeCctx.GetCurrentOutboundParam().Amount = fakeAmount
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(fakeCctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "amount mismatch")
		require.Nil(t, revertedLog)
		require.Nil(t, receivedLog)
	})
	t.Run("should fail on cctx index mismatch", func(t *testing.T) {
		cctx.Index = sample.Hash().Hex() // use an arbitrary index to make validation fail
		receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
		require.ErrorContains(t, err, "cctx index mismatch")
		require.Nil(t, receivedLog)
		require.Nil(t, revertedLog)
	})
}

func Test_ParseERC20WithdrawnEvent(t *testing.T) {
	// load archived outbound receipt that contains ERC20 Withdrawn event
	chainID := chains.Ethereum.ChainId
	nonce := uint64(8014)
	outboundHash := "0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146"
	custody := mocks.MockERC20Custody(t, chainID)
	custodyAddress := testutils.CustodyAddresses[chainID]
	cctx := testutils.LoadCctxByNonce(t, chainID, nonce)
	receipt := testutils.LoadEVMOutboundReceipt(
		t,
		TestDataDir,
		chainID,
		outboundHash,
		coin.CoinType_ERC20,
		testutils.EventERC20Withdraw,
	)

	t.Run("should parse ERC20 Withdrawn event from archived outbound receipt", func(t *testing.T) {
		withdrawn, err := parseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
		require.NoError(t, err)
		require.NotNil(t, withdrawn)
	})
	t.Run("should fail on erc20 custody address mismatch", func(t *testing.T) {
		// use an arbitrary address to make validation fail
		fakeCustodyAddress := sample.EthAddress()
		withdrawn, err := parseAndCheckWithdrawnEvent(cctx, receipt, fakeCustodyAddress, custody)
		require.ErrorContains(t, err, "error validating Withdrawn event")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail on receiver address mismatch", func(t *testing.T) {
		// load cctx and set receiver address to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.GetCurrentOutboundParam().Receiver = sample.EthAddress().Hex()
		withdrawn, err := parseAndCheckWithdrawnEvent(fakeCctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "receiver address mismatch")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail on asset mismatch", func(t *testing.T) {
		// load cctx and set asset to an arbitrary address
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeCctx.InboundParams.Asset = sample.EthAddress().Hex()
		withdrawn, err := parseAndCheckWithdrawnEvent(fakeCctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "asset mismatch")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail on amount mismatch", func(t *testing.T) {
		// load cctx and set amount to an arbitrary wrong value
		fakeCctx := testutils.LoadCctxByNonce(t, chainID, nonce)
		fakeAmount := sample.UintInRange(0, fakeCctx.GetCurrentOutboundParam().Amount.Uint64()-1)
		fakeCctx.GetCurrentOutboundParam().Amount = fakeAmount
		withdrawn, err := parseAndCheckWithdrawnEvent(fakeCctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "amount mismatch")
		require.Nil(t, withdrawn)
	})
	t.Run("should fail if no Withdrawn event found in receipt", func(t *testing.T) {
		// load receipt and remove Withdrawn event from logs
		receipt := testutils.LoadEVMOutboundReceipt(
			t,
			TestDataDir,
			chainID,
			outboundHash,
			coin.CoinType_ERC20,
			testutils.EventERC20Withdraw,
		)
		receipt.Logs = receipt.Logs[:1] // the 2nd log is Withdrawn event
		withdrawn, err := parseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
		require.ErrorContains(t, err, "no ERC20 Withdrawn event")
		require.Nil(t, withdrawn)
	})
}

func Test_FilterTSSOutbound(t *testing.T) {
	// load archived evm block
	// https://etherscan.io/block/19363323
	chain := chains.Ethereum
	chainID := chain.ChainId

	// load archived evm block
	// https://etherscan.io/block/19363323
	blockNumber := uint64(19363323)
	block := testutils.LoadEVMBlock(t, TestDataDir, chainID, blockNumber, true)

	// the outbound to be tested
	outboundNonce := uint64(7260)
	outboundHash := ethcommon.HexToHash("0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3")
	tx, receipt := testutils.LoadEVMOutboundNReceipt(t, TestDataDir, chainID, outboundHash.Hex(), coin.CoinType_Gas)

	ctx := context.Background()

	t.Run("should filter TSS outbound", func(t *testing.T) {
		// create evm observer for testing
		ob := newTestSuite(t)

		confirmations := ob.chainParams.OutboundConfirmationSafe()

		// create mock evm client with preloaded block, tx and receipt
		ob.evmMock.On("TransactionByHash", mock.Anything, outboundHash).Return(tx, false, nil)
		ob.evmMock.On("TransactionReceipt", mock.Anything, outboundHash).Return(receipt, nil)

		ob.BlockCache().Add(blockNumber, block)
		ob.WithLastBlock(blockNumber + confirmations - 1)

		// filter TSS outbound
		ob.filterTSSOutbound(ctx, blockNumber, blockNumber)

		// tx should be confirmed after filtering
		found := ob.isTxConfirmed(outboundNonce)
		require.True(t, found)

		// retrieve tx and receipt
		receipt, tx = ob.getTxNReceipt(outboundNonce)
		require.NotNil(t, tx)
		require.NotNil(t, receipt)
		require.Equal(t, outboundHash, tx.Hash())
		require.Equal(t, outboundNonce, tx.Nonce())
	})

	t.Run("should filter nothing on RPC error", func(t *testing.T) {
		ob := newTestSuite(t)

		ob.evmMock.On("BlockByNumberCustom", mock.Anything, mock.Anything).Return(nil, errors.New("rpc error"))

		// filter TSS outbound
		ob.filterTSSOutbound(ctx, blockNumber, blockNumber)

		// tx should be confirmed after filtering
		found := ob.isTxConfirmed(outboundNonce)
		require.False(t, found)
	})
}

// TODO: create mocks for gateway and ERC20Custody and uncomment these tests
// https://github.com/zeta-chain/node/issues/2669
//
//func Test_ParseOutboundReceivedValue(t *testing.T) {
//	chainID := chains.Ethereum.ChainId
//	connector, connectorAddr, custody, custodyAddr := getContractsByChainID(t, chainID)
//
//	t.Run("should parse and check ZetaReceived event from archived outbound receipt", func(t *testing.T) {
//		// load archived outbound receipt that contains ZetaReceived event
//		// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
//		nonce := uint64(9718)
//		coinType := coin.CoinType_Zeta
//		cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
//			t,
//			TestDataDir,
//			chainID,
//			nonce,
//			testutils.EventZetaReceived,
//		)
//		params := cctx.GetCurrentOutboundParam()
//		value, status, err := observer.parseOutboundReceivedValue(
//			cctx,
//			receipt,
//			outbound,
//			coinType,
//			connectorAddr,
//			connector,
//			custodyAddr,
//			custody,
//		)
//		require.NoError(t, err)
//		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
//		require.Equal(t, chains.ReceiveStatus_success, status)
//	})
//	t.Run("should parse and check ZetaReverted event from archived outbound receipt", func(t *testing.T) {
//		// load archived outbound receipt that contains ZetaReverted event
//		// use local network tx: 0x1487e6a31dd430306667250b72bf15b8390b73108b69f3de5c1b2efe456036a7
//		localChainID := chains.GoerliLocalnet.ChainId
//		nonce := uint64(14)
//		coinType := coin.CoinType_Zeta
//		connectorLocal, connectorAddrLocal, custodyLocal, custodyAddrLocal := getContractsByChainID(t, localChainID)
//		cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
//			t,
//			TestDataDir,
//			localChainID,
//			nonce,
//			testutils.EventZetaReverted,
//		)
//		params := cctx.GetCurrentOutboundParam()
//		value, status, err := observer.parseOutboundReceivedValue(
//			cctx, receipt, outbound, coinType, connectorAddrLocal, connectorLocal, custodyAddrLocal, custodyLocal)
//		require.NoError(t, err)
//		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
//		require.Equal(t, chains.ReceiveStatus_success, status)
//	})
//	t.Run("should parse and check ERC20 Withdrawn event from archived outbound receipt", func(t *testing.T) {
//		// load archived outbound receipt that contains ERC20 Withdrawn event
//		// https://etherscan.io/tx/0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146
//		nonce := uint64(8014)
//		coinType := coin.CoinType_ERC20
//		cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
//			t,
//			TestDataDir,
//			chainID,
//			nonce,
//			testutils.EventERC20Withdraw,
//		)
//		params := cctx.GetCurrentOutboundParam()
//		value, status, err := observer.parseOutboundReceivedValue(
//			cctx,
//			receipt,
//			outbound,
//			coinType,
//			connectorAddr,
//			connector,
//			custodyAddr,
//			custody,
//		)
//		require.NoError(t, err)
//		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
//		require.Equal(t, chains.ReceiveStatus_success, status)
//	})
//	t.Run("nothing to parse if coinType is Gas", func(t *testing.T) {
//		// load archived outbound receipt of Gas token transfer
//		// https://etherscan.io/tx/0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3
//		nonce := uint64(7260)
//		coinType := coin.CoinType_Gas
//		cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(t, TestDataDir, chainID, nonce, "")
//		params := cctx.GetCurrentOutboundParam()
//		value, status, err := observer.parseOutboundReceivedValue(
//			cctx,
//			receipt,
//			outbound,
//			coinType,
//			connectorAddr,
//			connector,
//			custodyAddr,
//			custody,
//		)
//		require.NoError(t, err)
//		require.True(t, params.Amount.BigInt().Cmp(value) == 0)
//		require.Equal(t, chains.ReceiveStatus_success, status)
//	})
//	t.Run("should fail on unknown coin type", func(t *testing.T) {
//		// load archived outbound receipt that contains ZetaReceived event
//		// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
//		nonce := uint64(9718)
//		coinType := coin.FungibleTokenCoinType(5) // unknown coin type
//		cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
//			t,
//			TestDataDir,
//			chainID,
//			nonce,
//			testutils.EventZetaReceived,
//		)
//		value, status, err := observer.parseOutboundReceivedValue(
//			cctx,
//			receipt,
//			outbound,
//			coinType,
//			connectorAddr,
//			connector,
//			custodyAddr,
//			custody,
//		)
//		require.ErrorContains(t, err, "unknown coin type")
//		require.Nil(t, value)
//		require.Equal(t, chains.ReceiveStatus_failed, status)
//	})
//	t.Run("should fail if unable to parse ZetaReceived event", func(t *testing.T) {
//		// load archived outbound receipt that contains ZetaReceived event
//		// https://etherscan.io/tx/0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f
//		nonce := uint64(9718)
//		coinType := coin.CoinType_Zeta
//		cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
//			t,
//			TestDataDir,
//			chainID,
//			nonce,
//			testutils.EventZetaReceived,
//		)
//
//		// use an arbitrary address to make event parsing fail
//		fakeConnectorAddress := sample.EthAddress()
//		value, status, err := observer.parseOutboundReceivedValue(
//			cctx,
//			receipt,
//			outbound,
//			coinType,
//			fakeConnectorAddress,
//			connector,
//			custodyAddr,
//			custody,
//		)
//		require.Error(t, err)
//		require.Nil(t, value)
//		require.Equal(t, chains.ReceiveStatus_failed, status)
//	})
//	t.Run("should fail if unable to parse ERC20 Withdrawn event", func(t *testing.T) {
//		// load archived outbound receipt that contains ERC20 Withdrawn event
//		// https://etherscan.io/tx/0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146
//		nonce := uint64(8014)
//		coinType := coin.CoinType_ERC20
//		cctx, outbound, receipt := testutils.LoadEVMCctxNOutboundNReceipt(
//			t,
//			TestDataDir,
//			chainID,
//			nonce,
//			testutils.EventERC20Withdraw,
//		)
//
//		// use an arbitrary address to make event parsing fail
//		fakeCustodyAddress := sample.EthAddress()
//		value, status, err := observer.parseOutboundReceivedValue(
//			cctx,
//			receipt,
//			outbound,
//			coinType,
//			connectorAddr,
//			connector,
//			fakeCustodyAddress,
//			custody,
//		)
//		require.Error(t, err)
//		require.Nil(t, value)
//		require.Equal(t, chains.ReceiveStatus_failed, status)
//	})
//}
