package mocks

import (
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/mock"
)

//go:generate mockery --name SuiClient --filename sui_client.go --case underscore --output .
type SuiClientMock interface {
	sui.ISuiAPI
}

func (m *SuiClient) MockSuiXQueryEvents(res models.PaginatedEventsResponse, err error) {
	m.On(
		"SuiXQueryEvents",
		mock.Anything,
		mock.Anything,
	).Return(res, err).Once()
}
