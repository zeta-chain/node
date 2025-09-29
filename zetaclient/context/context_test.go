package context_test

import (
	goctx "context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/context"
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
	app := context.New(config.New(false), nil, zerolog.Nop())

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
		app  = context.New(config.New(false), nil, zerolog.Nop())
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

func TestCopyWithTimeout(t *testing.T) {
	// ARRANGE
	var (
		app     = context.New(config.New(false), nil, zerolog.Nop())
		ctx1    = context.WithAppContext(goctx.Background(), app)
		timeout = 500 * time.Millisecond
	)

	// ACT
	ctx2, cancel := context.CopyWithTimeout(ctx1, goctx.Background(), timeout)
	defer cancel()

	// ASSERT
	// Verify that AppContext is copied correctly
	app2, err := context.FromContext(ctx2)
	assert.NoError(t, err)
	assert.NotNil(t, app2)
	assert.Equal(t, app, app2)

	// Verify that timeout is working
	start := time.Now()
	<-ctx2.Done()
	elapsed := time.Since(start)

	// The context should be cancelled after approximately the timeout duration
	assert.True(t, elapsed >= timeout, "context should be not cancelled too early")
	assert.True(t, elapsed < timeout*2, "context should not be cancelled too late")
	assert.ErrorIs(t, ctx2.Err(), goctx.DeadlineExceeded)
}
