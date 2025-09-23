package observer

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	cc "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/config"
)

func TestInbound(t *testing.T) {
	t.Run("No gateway provided", func(t *testing.T) {
		ts := newTestSuite(t)

		_, err := New(ts.baseObserver, ts.rpc, nil)
		require.Error(t, err)
	})

	t.Run("Donation", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given observer
		ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
		require.NoError(t, err)

		lastScanned := ts.SetupLastScannedTX(ts.gateway.AccountID())

		// Given mocked lite client calls
		donation := sample.TONDonation(t, ts.gateway.AccountID(), toncontracts.Donation{
			Sender: sample.GenerateTONAccountID(),
			Amount: tonCoins(t, "12"),
		})

		txs := []ton.Transaction{donation}

		ts.
			OnGetTransactionsSince(ts.gateway.AccountID(), lastScanned.Lt, txHash(lastScanned), txs, nil).
			Once()

		// ACT
		// Observe inbounds once
		err = ob.ObserveInbound(ts.ctx)

		// ASSERT
		assert.NoError(t, err)

		// nothing happened, but tx scanned
		lt, hash, err := encoder.DecodeTx(ob.LastTxScanned())
		assert.NoError(t, err)
		assert.Equal(t, donation.Lt, lt)
		assert.Equal(t, donation.Hash().Hex(), hash.Hex())
	})

	t.Run("Deposit", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given observer
		ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
		require.NoError(t, err)

		lastScanned := ts.SetupLastScannedTX(ts.gateway.AccountID())

		// Given mocked lite client calls
		deposit := toncontracts.Deposit{
			Sender:    sample.GenerateTONAccountID(),
			Amount:    tonCoins(t, "12"),
			Recipient: sample.EthAddress(),
		}

		depositTX := sample.TONDeposit(t, ts.gateway.AccountID(), deposit)
		txs := []ton.Transaction{depositTX}

		ts.
			OnGetTransactionsSince(ts.gateway.AccountID(), lastScanned.Lt, txHash(lastScanned), txs, nil).
			Once()

		ts.MockGetCctxByHash()

		// ACT
		// Observe inbounds once
		err = ob.ObserveInbound(ts.ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that cctx was sent to zetacore
		require.Len(t, ts.votesBag, 1)

		// Check CCTX
		cctx := ts.votesBag[0]

		assert.NotNil(t, cctx)

		assert.Equal(t, deposit.Sender.ToRaw(), cctx.Sender)
		assert.Equal(t, ts.chain.ChainId, cctx.SenderChainId)

		assert.Equal(t, "", cctx.Asset)
		assert.Equal(t, deposit.Amount.Uint64(), cctx.Amount.Uint64())
		assert.Equal(t, "", cctx.Message)
		assert.Equal(t, deposit.Recipient.Hex(), cctx.Receiver)
		assert.False(t, cctx.IsCrossChainCall)

		// Check hash & block height
		expectedHash := encoder.EncodeHash(depositTX.Lt, txHash(depositTX))
		assert.Equal(t, expectedHash, cctx.InboundHash)
		assert.Equal(t, uint64(0), cctx.InboundBlockHeight)
	})

	t.Run("Deposit and call", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given observer
		ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
		require.NoError(t, err)

		lastScanned := ts.SetupLastScannedTX(ts.gateway.AccountID())

		// Given mocked lite client calls
		const callData = "hey there"
		depositAndCall := toncontracts.DepositAndCall{
			Deposit: toncontracts.Deposit{
				Sender:    sample.GenerateTONAccountID(),
				Amount:    tonCoins(t, "4"),
				Recipient: sample.EthAddress(),
			},
			CallData: []byte(callData),
		}

		depositAndCallTX := sample.TONDepositAndCall(t, ts.gateway.AccountID(), depositAndCall)
		txs := []ton.Transaction{depositAndCallTX}

		ts.
			OnGetTransactionsSince(ts.gateway.AccountID(), lastScanned.Lt, txHash(lastScanned), txs, nil).
			Once()

		ts.MockGetCctxByHash()

		// ACT
		// Observe inbounds once
		err = ob.ObserveInbound(ts.ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that cctx was sent to zetacore
		require.Len(t, ts.votesBag, 1)

		// Check CCTX
		cctx := ts.votesBag[0]

		assert.NotNil(t, cctx)

		assert.Equal(t, depositAndCall.Sender.ToRaw(), cctx.Sender)
		assert.Equal(t, ts.chain.ChainId, cctx.SenderChainId)

		assert.Equal(t, "", cctx.Asset)
		assert.Equal(t, depositAndCall.Amount.Uint64(), cctx.Amount.Uint64())
		assert.Equal(t, hex.EncodeToString([]byte(callData)), cctx.Message)
		assert.Equal(t, depositAndCall.Recipient.Hex(), cctx.Receiver)
		assert.True(t, cctx.IsCrossChainCall)

		// Check hash & block height
		expectedHash := encoder.EncodeHash(depositAndCallTX.Lt, txHash(depositAndCallTX))
		assert.Equal(t, expectedHash, cctx.InboundHash)
		assert.Equal(t, uint64(0), cctx.InboundBlockHeight)
	})

	t.Run("Call", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given observer
		ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
		require.NoError(t, err)

		lastScanned := ts.SetupLastScannedTX(ts.gateway.AccountID())

		// Given mocked lite client calls
		const callData = "ping pong"
		call := toncontracts.Call{
			Sender:    sample.GenerateTONAccountID(),
			Recipient: sample.EthAddress(),
			CallData:  []byte(callData),
		}

		callTX := sample.TONCall(t, ts.gateway.AccountID(), call)
		txs := []ton.Transaction{callTX}

		ts.
			OnGetTransactionsSince(ts.gateway.AccountID(), lastScanned.Lt, txHash(lastScanned), txs, nil).
			Once()

		ts.MockGetCctxByHash()

		// ACT
		// Observe inbounds once
		err = ob.ObserveInbound(ts.ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that cctx was sent to zetacore
		require.Len(t, ts.votesBag, 1)

		// Check CCTX
		cctx := ts.votesBag[0]

		assert.NotNil(t, cctx)

		assert.Equal(t, call.Sender.ToRaw(), cctx.Sender)
		assert.Equal(t, ts.chain.ChainId, cctx.SenderChainId)

		assert.Equal(t, "", cctx.Asset)
		assert.Equal(t, uint64(0), cctx.Amount.Uint64())
		assert.Equal(t, hex.EncodeToString([]byte(callData)), cctx.Message)
		assert.Equal(t, call.Recipient.Hex(), cctx.Receiver)
		assert.Equal(t, coin.CoinType_NoAssetCall, cctx.CoinType)

		// Check hash & block height
		expectedHash := encoder.EncodeHash(callTX.Lt, txHash(callTX))
		assert.Equal(t, expectedHash, cctx.InboundHash)
		assert.Equal(t, uint64(0), cctx.InboundBlockHeight)
	})

	t.Run("Deposit restricted", func(t *testing.T) {
		// ARRANGE
		// Given restricted sender
		sender := sample.GenerateTONAccountID()

		// note this might be flaky because it's a global variable (*sad*)
		config.SetRestrictedAddressesFromConfig(config.Config{
			ComplianceConfig: config.ComplianceConfig{
				RestrictedAddresses: []string{sender.ToRaw()},
			},
		})

		// Given test suite
		ts := newTestSuite(t)

		// Given observer
		ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
		require.NoError(t, err)

		lastScanned := ts.SetupLastScannedTX(ts.gateway.AccountID())

		// Given mocked lite client calls
		deposit := toncontracts.Deposit{
			Sender:    sender,
			Amount:    tonCoins(t, "12"),
			Recipient: sample.EthAddress(),
		}

		depositTX := sample.TONDeposit(t, ts.gateway.AccountID(), deposit)
		txs := []ton.Transaction{depositTX}

		ts.
			OnGetTransactionsSince(ts.gateway.AccountID(), lastScanned.Lt, txHash(lastScanned), txs, nil).
			Once()

		// ACT
		// Observe inbounds once
		err = ob.ObserveInbound(ts.ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that NO cctx was sent && log contains entry for restricted address
		require.Len(t, ts.votesBag, 0)
		require.Contains(t, ts.logger.String(), "restricted address detected in inbound")
	})

	// Yep, it's possible to have withdrawals here because we scroll through all gateway's txs
	t.Run("Withdrawal", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given observer
		ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
		require.NoError(t, err)

		lastScanned := ts.SetupLastScannedTX(ts.gateway.AccountID())

		// Given mocked lite client calls
		withdrawal := toncontracts.Withdrawal{
			Recipient: ton.MustParseAccountID("EQB5A1PJBbnxwf0YrA_bgWKyfuIv8GywEcfIAXrs3oZyqc1_"),
			Amount:    toncontracts.Coins(5),
			Seqno:     0,
		}

		ts.sign(&withdrawal)

		withdrawalSigner, err := withdrawal.Signer()
		require.NoError(t, err)
		require.Equal(t, ob.TSS().PubKey().AddressEVM().Hex(), withdrawalSigner.Hex())

		withdrawalTX := sample.TONWithdrawal(t, ts.gateway.AccountID(), withdrawal)
		txs := []ton.Transaction{withdrawalTX}

		ts.
			OnGetTransactionsSince(ts.gateway.AccountID(), lastScanned.Lt, txHash(lastScanned), txs, nil).
			Once()

		// ACT
		err = ob.ObserveInbound(ts.ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that no votes were sent
		require.Len(t, ts.votesBag, 0)

		// But an outbound tracker was created
		require.Len(t, ts.trackerBag, 1)

		tracker := ts.trackerBag[0]

		assert.Equal(t, uint64(withdrawal.Seqno), tracker.nonce)
		assert.Equal(t, encoder.EncodeTx(withdrawalTX), tracker.hash)
	})

	t.Run("Multiple transactions", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given observer
		ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
		require.NoError(t, err)

		lastScanned := ts.SetupLastScannedTX(ts.gateway.AccountID())

		// Given several transactions
		withdrawal := toncontracts.Withdrawal{
			Recipient: ton.MustParseAccountID("EQB5A1PJBbnxwf0YrA_bgWKyfuIv8GywEcfIAXrs3oZyqc1_"),
			Amount:    toncontracts.Coins(5),
			Seqno:     1,
		}
		ts.sign(&withdrawal)

		txs := []ton.Transaction{
			// should be skipped
			sample.TONDonation(t, ts.gateway.AccountID(), toncontracts.Donation{
				Sender: sample.GenerateTONAccountID(),
				Amount: tonCoins(t, "1"),
			}),
			// should be voted
			sample.TONDeposit(t, ts.gateway.AccountID(), toncontracts.Deposit{
				Sender:    sample.GenerateTONAccountID(),
				Amount:    tonCoins(t, "3"),
				Recipient: sample.EthAddress(),
			}),
			// should be skipped (invalid inbound message)
			sample.TONTransaction(t, sample.TONTransactionProps{
				Account: ts.gateway.AccountID(),
				Input:   &tlb.Message{},
			}),
			// should be voted
			sample.TONDeposit(t, ts.gateway.AccountID(), toncontracts.Deposit{
				Sender:    sample.GenerateTONAccountID(),
				Amount:    tonCoins(t, "3"),
				Recipient: sample.EthAddress(),
			}),
			// a tracker should be added
			sample.TONWithdrawal(t, ts.gateway.AccountID(), withdrawal),
			// should be skipped (invalid inbound/outbound messages)
			sample.TONTransaction(t, sample.TONTransactionProps{
				Account: ts.gateway.AccountID(),
				Input:   &tlb.Message{},
				Output:  &tlb.Message{},
			}),
		}

		ts.
			OnGetTransactionsSince(ts.gateway.AccountID(), lastScanned.Lt, txHash(lastScanned), txs, nil).
			Once()

		ts.MockGetCctxByHash()

		// ACT
		// Observe inbounds once
		err = ob.ObserveInbound(ts.ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that cctx was sent to zetacore
		assert.Equal(t, 2, len(ts.votesBag))

		var (
			hash1 = encoder.EncodeHash(txs[1].Lt, txHash(txs[1]))
			hash2 = encoder.EncodeHash(txs[3].Lt, txHash(txs[3]))
		)

		assert.Equal(t, hash1, ts.votesBag[0].InboundHash)
		assert.Equal(t, hash2, ts.votesBag[1].InboundHash)

		// Check that last scanned tx points to the last tx in a list (even if it was skipped)
		var (
			lastTX          = txs[len(txs)-1]
			lastScannedHash = ob.LastTxScanned()
		)

		lastLT, lastHash, err := encoder.DecodeTx(lastScannedHash)
		assert.NoError(t, err)
		assert.Equal(t, lastTX.Lt, lastLT)
		assert.Equal(t, lastTX.Hash().Hex(), lastHash.Hex())

		// Check that a tracker was added
		assert.Len(t, ts.trackerBag, 1)
		tracker := ts.trackerBag[0]

		assert.Equal(t, uint64(withdrawal.Seqno), tracker.nonce)
		assert.Equal(t, encoder.EncodeTx(txs[4]), tracker.hash)
	})
}

func TestInboundTracker(t *testing.T) {
	// ARRANGE
	ts := newTestSuite(t)

	// Given observer
	ob, err := New(ts.baseObserver, ts.rpc, ts.gateway)
	require.NoError(t, err)

	// Given TON gateway transactions
	// should be voted
	deposit := toncontracts.Deposit{
		Sender:    sample.GenerateTONAccountID(),
		Amount:    toncontracts.Coins(5),
		Recipient: sample.EthAddress(),
	}

	txDeposit := sample.TONDeposit(t, ts.gateway.AccountID(), deposit)
	ts.MockGetTransaction(ts.gateway.AccountID(), txDeposit)

	// Should be skipped (I doubt anyone would vote for this gov proposal, but let’s still put up rail guards)
	txWithdrawal := sample.TONWithdrawal(t, ts.gateway.AccountID(), toncontracts.Withdrawal{
		Recipient: sample.GenerateTONAccountID(),
		Amount:    toncontracts.Coins(5),
		Seqno:     1,
	})
	ts.MockGetTransaction(ts.gateway.AccountID(), txWithdrawal)

	// Given inbound trackers from zetacore
	trackers := []cc.InboundTracker{
		ts.TxToInboundTracker(txDeposit),
		ts.TxToInboundTracker(txWithdrawal),
	}
	ts.MockGetCctxByHash()

	ts.OnGetInboundTrackersForChain(trackers).Once()

	// ACT
	err = ob.ProcessInboundTrackers(ts.ctx)

	// ARRANGE
	require.NoError(t, err)
	require.Len(t, ts.votesBag, 1)

	vote := ts.votesBag[0]
	assert.Equal(t, deposit.Amount, vote.Amount)
	assert.Equal(t, deposit.Sender.ToRaw(), vote.Sender)
	assert.Equal(t, "", vote.Message)
	assert.Equal(t, deposit.Recipient.Hex(), vote.Receiver)
}

func txHash(tx ton.Transaction) ton.Bits256 {
	return ton.Bits256(tx.Hash())
}
