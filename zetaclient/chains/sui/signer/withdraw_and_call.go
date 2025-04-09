package signer

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
)

const (
	typeSeparator    = "::"
	funcWithdrawImpl = "withdraw_impl"
	funcOnCall       = "on_call"
	moduleConnected  = "connected"
)

// withdrawAndCallPTB builds unsigned withdraw and call PTB Sui transaction
// it chains the following calls:
// 1. withdraw_impl on gateway
// 2. gas budget coin transfer to TSS
// 3. on_call on target contract
// The function returns a TxnMetaData object with tx bytes, the other fields are ignored
func withdrawAndCallPTB(
	signerAddrStr,
	gatewayPackageIDStr,
	gatewayModule,
	gatewayObjectIDStr,
	withdrawCapIDStr,
	coinTypeStr,
	amountStr,
	nonceStr,
	gasBudgetStr,
	receiver string,
	cp zetasui.CallPayload,
) (models.TxnMetaData, error) {
	ptb := suiptb.NewTransactionDataTransactionBuilder()

	// Parse arguments
	packageID, err := sui.PackageIdFromHex(gatewayPackageIDStr)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse package ID: %w", err)
	}

	coinType, err := parseTypeString(coinTypeStr)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse coin type: %w", err)
	}

	gatewayObjectID, err := sui.ObjectIdFromHex(gatewayObjectIDStr)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse gateway object ID: %w", err)
	}
	gatewayObject, err := ptb.Obj(suiptb.ObjectArg{
		SharedObject: &suiptb.SharedObjectArg{
			Id:                   gatewayObjectID,
			InitialSharedVersion: 0,
			Mutable:              true,
		},
	})
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to create object argument: %w", err)
	}

	amountUint64, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse amount: %w", err)
	}
	amount, err := ptb.Pure(amountUint64)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to create pure argument: %w", err)
	}

	nonceUint64, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse nonce: %w", err)
	}
	nonce, err := ptb.Pure(nonceUint64)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to create pure argument: %w", err)
	}

	gasBudgetUint64, err := strconv.ParseUint(gasBudgetStr, 10, 64)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse gas budget: %w", err)
	}
	gasBudget, err := ptb.Pure(gasBudgetUint64)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to create pure argument: %w", err)
	}

	withdrawCapID, err := sui.ObjectIdFromHex(withdrawCapIDStr)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse withdraw cap ID: %w", err)
	}
	withdrawCap, err := ptb.Obj(suiptb.ObjectArg{
		ImmOrOwnedObject: &sui.ObjectRef{
			ObjectId: withdrawCapID,
			Version:  0,
			Digest:   nil,
		},
	})
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to create object argument: %w", err)
	}

	// Move call for withdraw_impl and get its command index
	cmdIndex := uint16(len(ptb.Commands))
	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:  packageID,
			Module:   gatewayModule,
			Function: funcWithdrawImpl,
			TypeArguments: []sui.TypeTag{
				{Struct: coinType},
			},
			Arguments: []suiptb.Argument{
				gatewayObject,
				amount,
				nonce,
				gasBudget,
				withdrawCap,
			},
		},
	})

	// Create arguments to access the two results from the withdraw_impl call
	withdrawnCoinsArg := suiptb.Argument{
		NestedResult: &suiptb.NestedResult{
			Cmd:    cmdIndex,
			Result: 0, // First result (main coins)
		},
	}

	budgetCoinsArg := suiptb.Argument{
		NestedResult: &suiptb.NestedResult{
			Cmd:    cmdIndex,
			Result: 1, // Second result (budget coins)
		},
	}

	// Transfer the budget coins to the TSS address
	tssAddrArg, err := ptb.Pure(signerAddrStr)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to create pure address argument: %w", err)
	}

	// Transfer budget coins to the TSS address
	ptb.Command(suiptb.Command{
		TransferObjects: &suiptb.ProgrammableTransferObjects{
			Objects: []suiptb.Argument{budgetCoinsArg},
			Address: tssAddrArg,
		},
	})

	// Extract argument for on_call

	// The receiver in the cctx is used as target package ID
	targetPackageID, err := sui.PackageIdFromHex(receiver)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse target package ID: %w", err)
	}

	// Convert call payload type arguments in addition to the withdrawn coin type
	onCallTypeArgs := make([]sui.TypeTag, 0, len(cp.TypeArgs)+1)
	onCallTypeArgs = append(onCallTypeArgs, sui.TypeTag{Struct: coinType})
	for _, typeArg := range cp.TypeArgs {
		typeStruct, err := parseTypeString(typeArg)
		if err != nil {
			return models.TxnMetaData{}, fmt.Errorf("failed to parse type argument: %w", err)
		}
		onCallTypeArgs = append(onCallTypeArgs, sui.TypeTag{Struct: typeStruct})
	}

	// Build the args for on_call, contains withdrawns coins + payload objects + message
	onCallArgs := make([]suiptb.Argument, 0, len(cp.ObjectIDs)+1)
	onCallArgs = append(onCallArgs, withdrawnCoinsArg)

	// Add the payload objects, objects are all shared
	for _, objectID := range cp.ObjectIDs {
		objectIDParsed, err := sui.ObjectIdFromHex(objectID)
		if err != nil {
			return models.TxnMetaData{}, fmt.Errorf("failed to parse object ID: %w", err)
		}
		objectArg, err := ptb.Obj(suiptb.ObjectArg{
			SharedObject: &suiptb.SharedObjectArg{
				Id:                   objectIDParsed,
				InitialSharedVersion: 0,
				Mutable:              true,
			},
		})
		if err != nil {
			return models.TxnMetaData{}, fmt.Errorf("failed to create object argument: %w", err)
		}
		onCallArgs = append(onCallArgs, objectArg)
	}

	// Add any additional message arguments
	messageArg, err := ptb.Pure(cp.Message)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to create pure message argument: %w", err)
	}
	onCallArgs = append(onCallArgs, messageArg)

	// Call the target contract on_call
	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       targetPackageID,
			Module:        moduleConnected,
			Function:      funcOnCall,
			TypeArguments: onCallTypeArgs,
			Arguments:     onCallArgs,
		},
	})

	// Finish building the PTB
	pt := ptb.Finish()

	// Get the signer address
	signerAddr, err := sui.AddressFromHex(signerAddrStr)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to parse signer address: %w", err)
	}

	// TODO: get coin object for gas payment

	txData := suiptb.NewTransactionData(
		signerAddr,
		pt,
		[]*sui.ObjectRef{},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		return models.TxnMetaData{}, fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	// Encode the transaction bytes to base64
	return models.TxnMetaData{
		TxBytes: base64.StdEncoding.EncodeToString(txBytes),
	}, nil
}

func parseTypeString(t string) (*sui.StructTag, error) {
	parts := strings.Split(t, typeSeparator)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid type string: %s", t)
	}

	address, err := sui.AddressFromHex(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid address: %s", parts[0])
	}

	module := parts[1]
	name := parts[2]

	return &sui.StructTag{
		Address: address,
		Module:  module,
		Name:    name,
	}, nil
}
