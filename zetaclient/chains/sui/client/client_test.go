package client

import (
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

func TestCheckContainOwnedObject(t *testing.T) {
	tests := []struct {
		name    string
		input   []*models.SuiObjectResponse
		wantErr bool
	}{
		{
			name: "all immutable",
			input: []*models.SuiObjectResponse{
				{Data: &models.SuiObjectData{Owner: immutableOwner}},
				{Data: &models.SuiObjectData{Owner: immutableOwner}},
			},
			wantErr: false,
		},
		{
			name: "all shared",
			input: []*models.SuiObjectResponse{
				{Data: &models.SuiObjectData{Owner: map[string]interface{}{sharedOwner: map[string]interface{}{}}}},
			},
			wantErr: false,
		},
		{
			name: "mixed shared and immutable",
			input: []*models.SuiObjectResponse{
				{Data: &models.SuiObjectData{Owner: immutableOwner}},
				{Data: &models.SuiObjectData{Owner: map[string]interface{}{sharedOwner: nil}}},
			},
			wantErr: false,
		},
		{
			name: "missing data",
			input: []*models.SuiObjectResponse{
				{Data: nil},
			},
			wantErr: true,
		},
		{
			name: "unexpected string owner",
			input: []*models.SuiObjectResponse{
				{Data: &models.SuiObjectData{Owner: "SomeOtherOwner"}},
			},
			wantErr: true,
		},
		{
			name: "unknown owner type",
			input: []*models.SuiObjectResponse{
				{Data: &models.SuiObjectData{Owner: 123}},
			},
			wantErr: true,
		},
		{
			name: "map owner but not shared",
			input: []*models.SuiObjectResponse{
				{Data: &models.SuiObjectData{Owner: map[string]interface{}{"Owned": nil}}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckContainOwnedObject(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
