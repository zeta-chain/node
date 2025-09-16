package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	zetacrypto "github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/pkg/gas"
)

const (
	// TssMigrationGasMultiplierEVM is multiplied to the median gas price to get the gas price for the tss migration .
	// This is done to avoid the tss migration tx getting stuck in the mempool
	TssMigrationGasMultiplierEVM = "2.5"

	// TSSMigrationBufferAmountEVM is the buffer amount added to the gas price for the tss migration transaction
	TSSMigrationBufferAmountEVM = "2100000000"

	// ERC20CustodyMigrationGasMultiplierEVM is multiplied to the median gas price to get the gas price for the erc20 custody migration
	// NOTE: this is a integer type unlike type above because the message logic is slightly different and an integer is needed
	ERC20CustodyMigrationGasMultiplierEVM = 2

	// AssetCustodyWhitelistGasMultiplierEVM is multiplied to the median gas price to get the gas price for the asset custody whitelist
	AssetCustodyWhitelistGasMultiplierEVM = 2

	// ERC20CustodyPausingGasMultiplierEVM is multiplied to the median gas price to get the gas price for the erc20 custody pausing
	ERC20CustodyPausingGasMultiplierEVM = 2
)

// MigrateERC20CustodyFundsCmdCCTX returns a CCTX allowing to migrate ERC20 custody funds
func MigrateERC20CustodyFundsCmdCCTX(
	creator string,
	erc20Address string,
	custodyContractAddress string,
	newCustodyContractAddress string,
	chainID int64,
	amount sdkmath.Uint,
	gasPrice string,
	priorityFee string,
	tssPubKey string,
	currentNonce uint64,
) CrossChainTx {
	indexString := GetERC20CustodyMigrationCCTXIndexString(tssPubKey, currentNonce, chainID, erc20Address)
	hash := crypto.Keccak256Hash([]byte(indexString))

	return newCmdCCTX(
		creator,
		hash.Hex(),
		fmt.Sprintf(
			"%s:%s,%s,%s",
			constant.CmdMigrateERC20CustodyFunds,
			newCustodyContractAddress,
			erc20Address,
			amount.String(),
		),
		creator,
		hash.Hex(),
		custodyContractAddress,
		chainID,
		sdkmath.NewUint(0),
		100_000,
		gasPrice,
		priorityFee,
		tssPubKey,
	)
}

// GetERC20CustodyMigrationCCTXIndexString returns the index string of the CCTX for migrating ERC20 custody funds
func GetERC20CustodyMigrationCCTXIndexString(
	tssPubKey string,
	nonce uint64,
	chainID int64,
	erc20Address string,
) string {
	return fmt.Sprintf("%s-%s-%d-%d-%s", constant.CmdMigrateERC20CustodyFunds, tssPubKey, nonce, chainID, erc20Address)
}

// UpdateERC20CustodyPauseStatusCmdCCTX returns a CCTX allowing to update the pause status of the ERC20 custody contract
func UpdateERC20CustodyPauseStatusCmdCCTX(
	creator string,
	custodyContractAddress string,
	chainID int64,
	pause bool,
	gasPrice string,
	priorityFee string,
	tssPubKey string,
	currentNonce uint64,
) CrossChainTx {
	indexString := GetERC20CustodyPausingCmdCCTXIndexString(tssPubKey, currentNonce, chainID)
	hash := crypto.Keccak256Hash([]byte(indexString))

	params := constant.OptionUnpause
	if pause {
		params = constant.OptionPause
	}

	return newCmdCCTX(
		creator,
		hash.Hex(),
		fmt.Sprintf("%s:%s", constant.CmdUpdateERC20CustodyPauseStatus, params),
		creator,
		hash.Hex(),
		custodyContractAddress,
		chainID,
		sdkmath.NewUint(0),
		100_000,
		gasPrice,
		priorityFee,
		tssPubKey,
	)
}

// GetERC20CustodyPausingCmdCCTXIndexString returns the index string of the CCTX for updating the pause status of the ERC20 custody contract
func GetERC20CustodyPausingCmdCCTXIndexString(
	tssPubKey string,
	nonce uint64,
	chainID int64,
) string {
	return fmt.Sprintf("%s-%s-%d-%d", constant.CmdUpdateERC20CustodyPauseStatus, tssPubKey, nonce, chainID)
}

// WhitelistAssetCmdCCTX returns a CCTX allowing to whitelist an asset token on an external chain
func WhitelistAssetCmdCCTX(
	creator string,
	zrc20Address ethcommon.Address,
	assetAddress string,
	custodyContractAddress string,
	chainID int64,
	gasPrice string,
	priorityFee string,
	tssPubKey string,
) CrossChainTx {
	// calculate the cctx index
	// we use the deployed zrc20 contract address to generate a unique index
	// since other parts of the system may use the zrc20 for the index, we add a message specific suffix
	hash := crypto.Keccak256Hash(zrc20Address.Bytes(), []byte("WhitelistAsset"))

	return newCmdCCTX(
		creator,
		hash.Hex(),
		fmt.Sprintf("%s:%s", constant.CmdWhitelistAsset, assetAddress),
		creator,
		hash.Hex(),
		custodyContractAddress,
		chainID,
		sdkmath.NewUint(0),
		100_000,
		gasPrice,
		priorityFee,
		tssPubKey,
	)
}

// MigrateFundCmdCCTX returns a CCTX allowing to migrate funds from the current TSS to the new TSS
func MigrateFundCmdCCTX(
	blockHeight int64,
	creator string,
	inboundHash string,
	chainID int64,
	amount sdkmath.Uint,
	medianGasPrice sdkmath.Uint,
	priorityFee sdkmath.Uint,
	currentTSSPubKey string,
	newTSSPubKey string,
	additionalStaticChainInfo []chains.Chain,
) (CrossChainTx, error) {
	var (
		sender      string
		receiver    string
		gasLimit    uint64
		gasPrice    string
		finalAmount sdkmath.Uint
	)

	// set sender, receiver, gas limit, gas price and final amount based on the chain
	switch {
	case chains.IsEVMChain(chainID, additionalStaticChainInfo):
		ethAddressOld, err := zetacrypto.GetTSSAddrEVM(currentTSSPubKey)
		if err != nil {
			return CrossChainTx{}, err
		}
		ethAddressNew, err := zetacrypto.GetTSSAddrEVM(newTSSPubKey)
		if err != nil {
			return CrossChainTx{}, err
		}
		sender = ethAddressOld.String()
		receiver = ethAddressNew.String()
		gasLimit = gas.EVMSend
		gasPriceUint, err := gas.MultiplyGasPrice(medianGasPrice, TssMigrationGasMultiplierEVM)
		if err != nil {
			return CrossChainTx{}, err
		}
		evmFee := sdkmath.NewUint(gasLimit).
			Mul(gasPriceUint).
			Add(sdkmath.NewUintFromString(TSSMigrationBufferAmountEVM))
		if evmFee.GT(amount) {
			return CrossChainTx{}, errorsmod.Wrap(
				ErrInsufficientFundsTssMigration,
				fmt.Sprintf(
					"insufficient funds to pay for gas fee, amount: %s, gas fee: %s, chainid: %d",
					amount.String(),
					evmFee.String(),
					chainID,
				),
			)
		}
		gasPrice = gasPriceUint.String()
		finalAmount = amount.Sub(evmFee)
	case chains.IsBitcoinChain(chainID, additionalStaticChainInfo):
		bitcoinNetParams, err := chains.BitcoinNetParamsFromChainID(chainID)
		if err != nil {
			return CrossChainTx{}, err
		}
		btcAddressOld, err := zetacrypto.GetTSSAddrBTC(currentTSSPubKey, bitcoinNetParams)
		if err != nil {
			return CrossChainTx{}, err
		}
		btcAddressNew, err := zetacrypto.GetTSSAddrBTC(newTSSPubKey, bitcoinNetParams)
		if err != nil {
			return CrossChainTx{}, err
		}
		sender = btcAddressOld
		receiver = btcAddressNew
		gasLimit = 1_000_000
		gasPrice = medianGasPrice.MulUint64(2).String()
		finalAmount = amount
	default:
		return CrossChainTx{}, errorsmod.Wrap(ErrUnsupportedChain, fmt.Sprintf("chain %d is not supported", chainID))
	}

	indexString := GetTssMigrationCCTXIndexString(currentTSSPubKey, newTSSPubKey, chainID, amount, blockHeight)
	hash := crypto.Keccak256Hash([]byte(indexString))

	return newCmdCCTX(
		creator,
		hash.Hex(),
		fmt.Sprintf("%s:%s", constant.CmdMigrateTssFunds, "Funds Migrator Admin Cmd"),
		sender,
		inboundHash,
		receiver,
		chainID,
		finalAmount,
		gasLimit,
		gasPrice,
		priorityFee.MulUint64(2).String(),
		currentTSSPubKey,
	), nil
}

// GetTssMigrationCCTXIndexString returns the index string of the CCTX for migrating funds from the current TSS to the new TSS
func GetTssMigrationCCTXIndexString(
	currentTssPubkey,
	newTssPubkey string,
	chainID int64,
	amount sdkmath.Uint,
	height int64,
) string {
	return fmt.Sprintf("%s-%s-%d-%s-%d", currentTssPubkey, newTssPubkey, chainID, amount.String(), height)
}

// newCmdCCTX returns a new CCTX for admin cmd with the given parameters
func newCmdCCTX(
	creator string,
	index string,
	relayedMessage,
	sender string,
	inboundHash string,
	receiver string,
	chainID int64,
	amount sdkmath.Uint,
	gasLimit uint64,
	medianGasPrice string,
	priorityFee string,
	tssPubKey string,
) CrossChainTx {
	return CrossChainTx{
		Creator:        creator,
		Index:          index,
		RelayedMessage: relayedMessage,
		CctxStatus: &Status{
			Status: CctxStatus_PendingOutbound,
		},
		InboundParams: &InboundParams{
			Sender:       sender,
			CoinType:     coin.CoinType_Cmd,
			ObservedHash: inboundHash,
			// irrelevant to observer voting, set it to success by default
			Status: InboundStatus_SUCCESS,
			// any inbound initiated from ZetaChain is deemed safely confirmed
			ConfirmationMode: ConfirmationMode_SAFE,
		},
		OutboundParams: []*OutboundParams{
			{
				Receiver:        receiver,
				ReceiverChainId: chainID,
				CoinType:        coin.CoinType_Cmd,
				Amount:          amount,
				CallOptions: &CallOptions{
					GasLimit: gasLimit,
				},
				GasPrice:       medianGasPrice,
				GasPriorityFee: priorityFee,
				TssPubkey:      tssPubKey,
				// use SAFE confirmation mode as default value.
				// zetaclient should ALWAYS use SAFE confirmation mode to confirm a CMD tx.
				ConfirmationMode: ConfirmationMode_SAFE,
			},
		},
	}
}
