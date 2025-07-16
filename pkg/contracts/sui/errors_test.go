package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewMoveAbortFromExecutionError(t *testing.T) {
	tests := []struct {
		name     string
		errorMsg string
		want     MoveAbort
		wantErr  bool
	}{
		{
			name:     "valid MoveAbort error message",
			errorMsg: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, 3) in command 0",
			want: MoveAbort{
				Message: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, 3) in command 0",
				Code:    3,
			},
		},
		{
			name:     "invalid MoveAbort error message",
			errorMsg: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, ) in command 0",
			want:     MoveAbort{},
			wantErr:  true,
		},
		{
			name:     "invalid MoveAbort error code",
			errorMsg: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, -1) in command 0",
			want:     MoveAbort{},
			wantErr:  true,
		},
		{
			name:     "other execution error message",
			errorMsg: "InsufficientCoinBalance in command 0",
			want:     MoveAbort{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moveAbort, err := NewMoveAbortFromExecutionError(tt.errorMsg)
			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, moveAbort)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, moveAbort)
		})
	}
}

func Test_IsRetryableExecutionError(t *testing.T) {
	tests := []struct {
		name         string
		errorMsgExec string
		want         bool
		errMsg       string
	}{
		{
			name:         "retryable: MoveAbort from withdraw_impl",
			errorMsgExec: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, 3) in command 0",
			want:         true,
		},
		{
			name:         "non-retryable: MoveAbort from withdraw_impl",
			errorMsgExec: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, 4) in command 0",
			want:         false,
		},
		{
			name:         "non-retryable: MoveAbort from on_call",
			errorMsgExec: "MoveAbort(MoveLocation { module: ModuleId { address: 0d553a3393a41e2fd88eae86f4f0423b86f8d76d57d7a427442244e2d8919761, name: Identifier(\"connected\") }, function: 1, instruction: 7, function_name: Some(\"on_call\") }, 3) in command 3",
			want:         false,
		},
		{
			name:         "non-retryable: command index not present",
			errorMsgExec: "InsufficientGas", // command index is not present
			want:         false,
		},
		{
			name:         "non-retryable: command index out of range",
			errorMsgExec: "MoveAbort(..., 3) in command 65536", // command index is out of range
			want:         false,
		},
		{
			name:         "non-retryable: command index >= 5",
			errorMsgExec: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, 3) in command 5",
			want:         false,
			errMsg:       "invalid command index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsRetryableExecutionError(tt.errorMsgExec)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
