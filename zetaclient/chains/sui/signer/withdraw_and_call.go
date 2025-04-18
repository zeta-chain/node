package signer

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pkg/errors"

	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
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
	gatewayModule string,
	gatewayObjRef,
	suiCoinObjRef,
	withdrawCapObjRef sui.ObjectRef,
	onCallObjectRefs []sui.ObjectRef,
	coinTypeStr,
	amountStr,
	nonceStr,
	gasBudgetStr,
	receiver string,
	cp zetasui.CallPayload,
) (tx models.TxnMetaData, err error) {
	ptb := suiptb.NewTransactionDataTransactionBuilder()

	// Parse arguments
	signerAddr, err := sui.AddressFromHex(signerAddrStr)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to parse signer address %s", signerAddrStr)
	}

	gatewayPackageID, err := sui.PackageIdFromHex(gatewayPackageIDStr)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to parse package ID %s", gatewayPackageIDStr)
	}

	coinType, err := parseTypeString(coinTypeStr)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to parse coin type %s", coinTypeStr)
	}

	gatewayObject, err := ptb.Obj(suiptb.ObjectArg{
		SharedObject: &suiptb.SharedObjectArg{
			Id:                   gatewayObjRef.ObjectId,
			InitialSharedVersion: gatewayObjRef.Version,
			Mutable:              true,
		},
	})
	if err != nil {
		return tx, errors.Wrap(err, "failed to create gateway object argument")
	}

	amountUint64, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to parse amount %s", amountStr)
	}
	amount, err := ptb.Pure(amountUint64)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to create amount argument")
	}

	nonceUint64, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to parse nonce %s", nonceStr)
	}
	nonce, err := ptb.Pure(nonceUint64)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to create nonce argument")
	}

	gasBudgetUint64, err := strconv.ParseUint(gasBudgetStr, 10, 64)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to parse gas budget %s", gasBudgetStr)
	}
	gasBudget, err := ptb.Pure(gasBudgetUint64)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to create gas budget argument")
	}

	withdrawCap, err := ptb.Obj(suiptb.ObjectArg{ImmOrOwnedObject: &withdrawCapObjRef})
	if err != nil {
		return tx, errors.Wrapf(err, "failed to create withdraw cap object argument")
	}

	// Move call for withdraw_impl and get its command index (0)
	cmdIndex := uint16(len(ptb.Commands))
	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:  gatewayPackageID,
			Module:   gatewayModule,
			Function: zetasui.FuncWithdrawImpl,
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

	// Transfer gas budget coins to the TSS address
	tssAddrArg, err := ptb.Pure(signerAddr)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to create tss address argument")
	}

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
		return tx, errors.Wrapf(err, "failed to parse target package ID %s", receiver)
	}

	// Build the type arguments for on_call in order: [withdrawn coin type, ... payload type arguments]
	onCallTypeArgs := make([]sui.TypeTag, 0, len(cp.TypeArgs)+1)
	onCallTypeArgs = append(onCallTypeArgs, sui.TypeTag{Struct: coinType})
	for _, typeArg := range cp.TypeArgs {
		typeStruct, err := parseTypeString(typeArg)
		if err != nil {
			return tx, errors.Wrapf(err, "failed to parse type argument %s", typeArg)
		}
		onCallTypeArgs = append(onCallTypeArgs, sui.TypeTag{Struct: typeStruct})
	}

	// Build the args for on_call: [withdrawns coins + payload objects + message]
	onCallArgs := make([]suiptb.Argument, 0, len(cp.ObjectIDs)+1)
	onCallArgs = append(onCallArgs, withdrawnCoinsArg)

	// Add the payload objects, objects are all shared
	for _, onCallObjectRef := range onCallObjectRefs {
		objectArg, err := ptb.Obj(suiptb.ObjectArg{
			SharedObject: &suiptb.SharedObjectArg{
				Id:                   onCallObjectRef.ObjectId,
				InitialSharedVersion: onCallObjectRef.Version,
				Mutable:              true,
			},
		})
		if err != nil {
			return tx, errors.Wrapf(err, "failed to create object argument: %v", onCallObjectRef)
		}
		onCallArgs = append(onCallArgs, objectArg)
	}

	// Add any additional message arguments
	messageArg, err := ptb.Pure(cp.Message)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to create message argument: %x", cp.Message)
	}
	onCallArgs = append(onCallArgs, messageArg)

	// Call the target contract on_call
	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       targetPackageID,
			Module:        zetasui.ModuleConnected,
			Function:      zetasui.FuncOnCall,
			TypeArguments: onCallTypeArgs,
			Arguments:     onCallArgs,
		},
	})

	// Finish building the PTB
	pt := ptb.Finish()

	// Get the signer address
	txData := suiptb.NewTransactionData(
		signerAddr,
		pt,
		[]*sui.ObjectRef{
			&suiCoinObjRef,
		},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to marshal transaction data: %v", txData)
	}

	// Encode the transaction bytes to base64
	return models.TxnMetaData{
		TxBytes: base64.StdEncoding.EncodeToString(txBytes),
	}, nil
}

func parseTypeString(t string) (*sui.StructTag, error) {
	parts := strings.Split(t, zetasui.TypeSeparator)
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

// getWithdrawAndCallObjectRefs returns the SUI object references for withdraw and call
func (s *Signer) getWithdrawAndCallObjectRefs(
	ctx context.Context,
	gatewayID, withdrawCapID string,
	onCallObjectIDs []string,
) ([]sui.ObjectRef, error) {
	objectIDs := append([]string{gatewayID, withdrawCapID}, onCallObjectIDs...)

	// query objects in batch
	suiObjects, err := s.client.SuiMultiGetObjects(ctx, models.SuiMultiGetObjectsRequest{
		ObjectIds: objectIDs,
		Options: models.SuiObjectDataOptions{
			// show owner info in order to retrieve object initial shared version
			ShowOwner: true,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get SUI objects for %v", objectIDs)
	}

	// convert object data to object references
	objectRefs := make([]sui.ObjectRef, 0, len(objectIDs))

	for _, object := range suiObjects {
		objectID, err := sui.ObjectIdFromHex(object.Data.ObjectId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse SUI object ID for %s", object.Data.ObjectId)
		}

		objectVersion, err := strconv.ParseUint(object.Data.Version, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse SUI object version for %s", object.Data.ObjectId)
		}

		// must use initial version for shared object, not the current version
		// withdraw cap is not a shared object, so we must use current version
		if object.Data.ObjectId != withdrawCapID {
			objectVersion, err = zetasui.ExtractInitialSharedVersion(*object.Data)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to extract initial shared version for %s", object.Data.ObjectId)
			}
		}

		objectDigest, err := sui.NewBase58(object.Data.Digest)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse SUI object digest for %s", object.Data.ObjectId)
		}

		objectRefs = append(objectRefs, sui.ObjectRef{
			ObjectId: objectID,
			Version:  objectVersion,
			Digest:   objectDigest,
		})
	}

	return objectRefs, nil
}

// getTSSSuiCoinObjectRef returns the latest SUI coin object reference for the TSS address
// Note: the SUI object may change over time, so we need to get the latest object
func (s *Signer) getTSSSuiCoinObjectRef(ctx context.Context) (sui.ObjectRef, error) {
	coins, err := s.client.SuiXGetAllCoins(ctx, models.SuiXGetAllCoinsRequest{
		Owner: s.TSS().PubKey().AddressSui(),
	})
	if err != nil {
		return sui.ObjectRef{}, errors.Wrap(err, "unable to get TSS coins")
	}

	// locate the SUI coin object under TSS account
	var suiCoin *models.CoinData
	for _, coin := range coins.Data {
		if zetasui.IsSUIType(zetasui.CoinType(coin.CoinType)) {
			suiCoin = &coin
			break
		}
	}
	if suiCoin == nil {
		return sui.ObjectRef{}, errors.New("SUI coin not found")
	}

	// convert coin data to object ref
	suiCoinID, err := sui.ObjectIdFromHex(suiCoin.CoinObjectId)
	if err != nil {
		return sui.ObjectRef{}, fmt.Errorf("failed to parse SUI coin ID: %w", err)
	}
	suiCoinVersion, err := strconv.ParseUint(suiCoin.Version, 10, 64)
	if err != nil {
		return sui.ObjectRef{}, fmt.Errorf("failed to parse SUI coin version: %w", err)
	}
	suiCoinDigest, err := sui.NewBase58(suiCoin.Digest)
	if err != nil {
		return sui.ObjectRef{}, fmt.Errorf("failed to parse SUI coin digest: %w", err)
	}

	return sui.ObjectRef{
		ObjectId: suiCoinID,
		Version:  suiCoinVersion,
		Digest:   suiCoinDigest,
	}, nil
}
