package requestctx

import "context"

type developerGuardBypassKey struct{}

// WithDeveloperGuardBypass marks the context so subscription and quota guards are skipped.
func WithDeveloperGuardBypass(ctx context.Context) context.Context {
	return context.WithValue(ctx, developerGuardBypassKey{}, true)
}

// HasDeveloperGuardBypass reports whether developer guard bypass is active for the request.
func HasDeveloperGuardBypass(ctx context.Context) bool {
	value, ok := ctx.Value(developerGuardBypassKey{}).(bool)
	return ok && value
}

// WithoutDeveloperGuardBypass enforces normal subscription policy for private modules.
func WithoutDeveloperGuardBypass(ctx context.Context) context.Context {
	return context.WithValue(ctx, developerGuardBypassKey{}, false)
}
