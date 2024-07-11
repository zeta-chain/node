package context_test

import (
	goctx "context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
)

func TestFromContext(t *testing.T) {
	// ARRANGE #1
	// Given default go ctx
	ctx := goctx.Background()

	// ACT #1
	// Extract App
	_, err := context.FromContext(ctx)

	// ASSERT #1
	assert.ErrorIs(t, err, context.ErrNotSet)

	// ARRANGE #2
	// Given basic app
	app := context.New(config.New(false), zerolog.Nop())

	// That is included in the ctx
	ctx = context.WithAppContext(ctx, app)

	// ACT #2
	app2, err := context.FromContext(ctx)

	// ASSERT #2
	assert.NoError(t, err)
	assert.NotNil(t, app2)
	assert.Equal(t, app, app2)
	assert.NotEmpty(t, app.Config())
}

func TestCopy(t *testing.T) {
	// ARRANGE
	var (
		app  = context.New(config.New(false), zerolog.Nop())
		ctx1 = context.WithAppContext(goctx.Background(), app)
	)

	// ACT
	ctx2 := context.Copy(ctx1, goctx.Background())

	// ASSERT
	app2, err := context.FromContext(ctx2)
	assert.NoError(t, err)
	assert.NotNil(t, app2)
	assert.Equal(t, app, app2)
}
