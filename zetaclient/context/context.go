package context

import (
	goctx "context"
	"time"

	"github.com/pkg/errors"
)

type appContextKey struct{}

var ErrNotSet = errors.New("unable to get AppContext from context.Context")

// WithAppContext applied AppContext to standard Go context.Context.
func WithAppContext(ctx goctx.Context, app *AppContext) goctx.Context {
	return goctx.WithValue(ctx, appContextKey{}, app)
}

// FromContext extracts AppContext from context.Context
func FromContext(ctx goctx.Context) (*AppContext, error) {
	app, ok := ctx.Value(appContextKey{}).(*AppContext)
	if !ok || app == nil {
		return nil, ErrNotSet
	}

	return app, nil
}

// Copy copies AppContext from one context to another (is present).
// This is useful when you want to drop timeouts and deadlines from the context
// (e.g. run something in another goroutine).
func Copy(from, to goctx.Context) goctx.Context {
	app, err := FromContext(from)
	if err != nil {
		return to
	}

	return WithAppContext(to, app)
}

// CopyWithTimeout copies AppContext from one context to another and adds a timeout.
// This is useful when you want to run something in another goroutine with a timeout.
func CopyWithTimeout(from, to goctx.Context, timeout time.Duration) (goctx.Context, goctx.CancelFunc) {
	app, err := FromContext(from)
	if err != nil {
		// If no AppContext found, just add timeout to the target context
		return goctx.WithTimeout(to, timeout)
	}

	// Create a new context with timeout
	ctxWithTimeout, cancel := goctx.WithTimeout(to, timeout)

	// Copy the AppContext to the new context
	return WithAppContext(ctxWithTimeout, app), cancel
}
