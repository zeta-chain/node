package sui

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

func Test_WithdrawAndCallPTB_TokenAmount(t *testing.T) {
	event := WithdrawAndCallPTB{
		Amount: math.NewUint(100),
		Nonce:  1,
	}
	require.Equal(t, math.NewUint(100), event.TokenAmount())
}

func Test_WithdrawAndCallPTB_TxNonce(t *testing.T) {
	event := WithdrawAndCallPTB{
		Amount: math.NewUint(100),
		Nonce:  1,
	}
	require.Equal(t, uint64(1), event.TxNonce())
}

func Test_ExtractInitialSharedVersion(t *testing.T) {
	tests := []struct {
		name        string
		objData     models.SuiObjectData
		wantVersion uint64
		errMsg      string
	}{
		{
			name: "successful extraction",
			objData: models.SuiObjectData{
				Owner: map[string]any{
					"Shared": map[string]any{
						"initial_shared_version": float64(3),
					},
				},
			},
			wantVersion: 3,
		},
		{
			name: "invalid owner type",
			objData: models.SuiObjectData{
				Owner: "invalid",
			},
			wantVersion: 0,
			errMsg:      "invalid object owner type string",
		},
		{
			name: "missing shared object",
			objData: models.SuiObjectData{
				Owner: map[string]any{
					"Owned": map[string]any{},
				},
			},
			wantVersion: 0,
			errMsg:      "missing shared object",
		},
		{
			name: "invalid shared object type",
			objData: models.SuiObjectData{
				Owner: map[string]any{
					"Shared": "invalid",
				},
			},
			wantVersion: 0,
			errMsg:      "invalid shared object type string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := ExtractInitialSharedVersion(tt.objData)
			if tt.errMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantVersion, version)
		})
	}
}

func Test_parseWithdrawAndCallPTB(t *testing.T) {
	const (
		txHash    = "AvpjaWZmHEJqj6AewiyK3bziGE2RPH5URub9nMx1sNL7"
		packageID = "0xa0464ffe0ffb12b2a474e0669e15aeb0b1e2b31a1865cca83f47b42c4707a550"
		amountStr = "100"
		nonceStr  = "2"
	)

	gw := &Gateway{
		packageID: packageID,
	}

	tests := []struct {
		name     string
		response models.SuiTransactionBlockResponse
		want     WithdrawAndCallPTB
		errMsg   string
	}{
		{
			name:     "valid transaction block",
			response: createPTBResponse(txHash, packageID, amountStr, nonceStr),
			want: WithdrawAndCallPTB{
				MoveCall: MoveCall{
					PackageID:  packageID,
					Module:     GatewayModule,
					Function:   FuncWithdrawImpl,
					ArgIndexes: ptbWithdrawImplArgIndexes,
				},
				Amount: math.NewUint(100),
				Nonce:  2,
			},
		},
		{
			name: "invalid number of inputs",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				res.Transaction.Data.Transaction.Inputs = res.Transaction.Data.Transaction.Inputs[:4]
				return res
			}(),
			errMsg: "invalid number of inputs",
		},
		{
			name: "unable to parse withdraw_impl",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				res.Transaction.Data.Transaction.Transactions[0] = "invalid"
				return res
			}(),
			errMsg: "unable to parse withdraw_impl command",
		},
		{
			name: "invalid package ID",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				moveCall := res.Transaction.Data.Transaction.Transactions[0].(map[string]any)["MoveCall"].(map[string]any)
				moveCall["package"] = "wrong_package"
				return res
			}(),
			errMsg: "invalid package id",
		},
		{
			name: "invalid module name",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				moveCall := res.Transaction.Data.Transaction.Transactions[0].(map[string]any)["MoveCall"].(map[string]any)
				moveCall["module"] = "wrong_module"
				return res
			}(),
			errMsg: "invalid module name",
		},
		{
			name: "invalid function name",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				moveCall := res.Transaction.Data.Transaction.Transactions[0].(map[string]any)["MoveCall"].(map[string]any)
				moveCall["function"] = "wrong_function"
				return res
			}(),
			errMsg: "invalid function name",
		},
		{
			name: "invalid argument indexes",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				moveCall := res.Transaction.Data.Transaction.Transactions[0].(map[string]any)["MoveCall"].(map[string]any)
				arguments := moveCall["arguments"].([]any)
				arguments[0] = map[string]any{"Input": float64(5)} // Change index to make it invalid
				return res
			}(),
			errMsg: "invalid argument indexes",
		},
		{
			name: "invalid amount format",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				res.Transaction.Data.Transaction.Inputs[1] = models.SuiCallArg{
					"value": "invalid_number",
				}
				return res
			}(),
			errMsg: "unable to parse amount",
		},
		{
			name: "invalid nonce format",
			response: func() models.SuiTransactionBlockResponse {
				res := createPTBResponse(txHash, packageID, amountStr, nonceStr)
				res.Transaction.Data.Transaction.Inputs[2] = models.SuiCallArg{
					"value": "invalid_nonce",
				}
				return res
			}(),
			errMsg: "unable to parse nonce",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event, content, err := gw.parseWithdrawAndCallPTB(test.response)
			if test.errMsg != "" {
				require.ErrorContains(t, err, test.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, txHash, event.TxHash)
			require.Zero(t, event.EventIndex)
			require.Equal(t, WithdrawAndCallEvent, event.EventType)

			withdrawCallPTB, ok := content.(WithdrawAndCallPTB)
			require.True(t, ok)
			require.Equal(t, test.want, withdrawCallPTB)
		})
	}
}

func createPTBResponse(txHash, packageID, amount, nonce string) models.SuiTransactionBlockResponse {
	return models.SuiTransactionBlockResponse{
		Digest: txHash,
		Transaction: models.SuiTransactionBlock{
			Data: models.SuiTransactionBlockData{
				Transaction: models.SuiTransactionBlockKind{
					Inputs: []models.SuiCallArg{
						{
							"initialSharedVersion": "3",
							"objectId":             "0xb3630c3eba7b1211c12604a4ceade7a5c0811c4a5eb55af227f9943fcef0e24c",
						},
						{
							"type":  "pure",
							"value": amount,
						},
						{
							"type":  "pure",
							"value": nonce,
						},
						{
							"type":  "pure",
							"value": "1000",
						},
						{
							"digest":   "26kdbzHiCt4nBFbMA6DjaazjRw98d6BRbe5gy81Wr1Aj",
							"objectId": "0x48a52371089644d30703f726ca5c30cf76e85347b263549e092ccf22ac059c6c",
						},
					},
					Transactions: []any{
						map[string]any{
							"MoveCall": map[string]any{
								"arguments": []any{
									map[string]any{"Input": float64(0)},
									map[string]any{"Input": float64(1)},
									map[string]any{"Input": float64(2)},
									map[string]any{"Input": float64(3)},
									map[string]any{"Input": float64(4)},
								},
								"function": FuncWithdrawImpl,
								"module":   GatewayModule,
								"package":  packageID,
							},
						},
						map[string]any{
							"TransferObjects": map[string]any{},
						},
						map[string]any{
							"MoveCall": map[string]any{
								"arguments": []any{},
								"function":  FuncSetMessageContext,
								"module":    GatewayModule,
								"package":   packageID,
							},
						},
						map[string]any{
							"MoveCall": map[string]any{
								"arguments": []any{},
								"function":  FuncOnCall,
								"module":    ModuleConnected,
								"package":   "target_package_id",
							},
						},
						map[string]any{
							"MoveCall": map[string]any{
								"arguments": []any{},
								"function":  FuncResetMessageContext,
								"module":    GatewayModule,
								"package":   packageID,
							},
						},
					},
				},
			},
		},
	}
}
