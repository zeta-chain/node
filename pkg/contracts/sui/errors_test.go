package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MoveAbortFromExecutionError(t *testing.T) {
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
		name     string
		errorMsg string
		want     bool
	}{
		{
			name:     "retryable MoveAbort",
			errorMsg: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, 3) in command 0",
			want:     true,
		},
		{
			name:     "non-retryable MoveAbort",
			errorMsg: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier(\"gateway\") }, function: 11, instruction: 37, function_name: Some(\"withdraw_impl\") }, 4) in command 0",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsRetryableExecutionError(tt.errorMsg)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
