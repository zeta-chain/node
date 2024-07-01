package context

import (
	goctx "context"

	"github.com/pkg/errors"
)

type appContextKey struct{}

var ErrNotSet = errors.New("AppContext is not set in the context.Context")

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
